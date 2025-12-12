# Data Model: 日志系统数据结构

**Feature**: 012-logging-system  
**Date**: 2025-12-12  
**Purpose**: 定义日志系统的核心数据结构、字段、关系和验证规则

## 概述

本文档定义了日志系统重构所需的所有数据结构,包括配置模型、日志记录模型、追踪上下文模型等。所有结构体遵循 Go 命名规范,字段使用 JSON tag 支持配置文件解析和日志序列化。

---

## 1. 配置模型

### 1.1 LoggerConfig - 日志系统配置

**用途**: 表示整个日志系统的配置,从 YAML 配置文件加载

**字段定义**:

```go
type LoggerConfig struct {
    Level      string                       `json:"level" yaml:"level"`           // 全局日志级别
    Stdout     bool                         `json:"stdout" yaml:"stdout"`         // 是否输出到终端
    TimeFormat string                       `json:"timeFormat" yaml:"timeFormat"` // 时间格式
    CtxKeys    []string                     `json:"ctxKeys" yaml:"ctxKeys"`       // 从 Context 提取的字段
    Instances  map[string]*InstanceConfig   `json:",inline" yaml:",inline"`       // 日志实例配置
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| Level | string | 是 | "all" | 全局日志级别: all/dev/prod/debug/info/notice/warning/error/critical |
| Stdout | bool | 否 | false | 是否同时输出到标准输出 |
| TimeFormat | string | 否 | "2006-01-02T15:04:05.000Z07:00" | 日志时间格式 (Go time layout) |
| CtxKeys | []string | 否 | ["TraceId"] | 从 Context 自动提取的字段名列表 |
| Instances | map[string]*InstanceConfig | 是 | - | 日志实例配置,key 为实例名称 (default/access/error/slow) |

**验证规则**:

```go
func (c *LoggerConfig) Validate() error {
    // 1. 验证日志级别
    validLevels := []string{"all", "dev", "prod", "debug", "info", "notice", "warning", "error", "critical"}
    if !contains(validLevels, c.Level) {
        return fmt.Errorf("无效的日志级别: %s, 有效值: %v", c.Level, validLevels)
    }
    
    // 2. 验证实例配置
    if len(c.Instances) == 0 {
        return errors.New("至少需要配置一个日志实例")
    }
    
    for name, instance := range c.Instances {
        if err := instance.Validate(name); err != nil {
            return err
        }
    }
    
    return nil
}
```

---

### 1.2 InstanceConfig - 日志实例配置

**用途**: 表示单个日志实例的配置 (如 app/access/error/slow)

**字段定义**:

```go
type InstanceConfig struct {
    Path               string `json:"path" yaml:"path"`                             // 日志文件路径
    File               string `json:"file" yaml:"file"`                             // 日志文件名模板
    Format             string `json:"format" yaml:"format"`                         // 日志格式: json/text
    Level              string `json:"level" yaml:"level"`                           // 实例日志级别 (覆盖全局)
    RotateSize         string `json:"rotateSize" yaml:"rotateSize"`                 // 按大小分割
    RotateExpire       string `json:"rotateExpire" yaml:"rotateExpire"`             // 日志过期时间
    RotateBackupLimit  int    `json:"rotateBackupLimit" yaml:"rotateBackupLimit"`   // 备份文件数量限制
    RotateBackupExpire string `json:"rotateBackupExpire" yaml:"rotateBackupExpire"` // 备份文件过期时间
    RotateCheckInterval string `json:"rotateCheckInterval" yaml:"rotateCheckInterval"` // 分割检查间隔
    StdoutPrint        bool   `json:"stdoutPrint" yaml:"stdoutPrint"`               // 是否同时输出到终端
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| Path | string | 是 | - | 日志文件目录路径,如 "./logs" |
| File | string | 是 | - | 日志文件名模板,支持 {Y-m-d} 等占位符 |
| Format | string | 否 | "json" | 日志格式: json (结构化) 或 text (纯文本) |
| Level | string | 否 | 继承全局 | 实例级别,覆盖全局 Level |
| RotateSize | string | 否 | - | 单文件大小限制,如 "100M", "500M" |
| RotateExpire | string | 否 | - | 日志文件过期时间,如 "7d", "30d" |
| RotateBackupLimit | int | 否 | 0 | 备份文件数量限制,0 表示不限制 |
| RotateBackupExpire | string | 否 | "30d" | 备份文件过期时间 |
| RotateCheckInterval | string | 否 | "1h" | 分割检查间隔 |
| StdoutPrint | bool | 否 | false | 是否同时输出到终端 |

**验证规则**:

```go
func (c *InstanceConfig) Validate(name string) error {
    // 1. 验证路径
    if c.Path == "" {
        return fmt.Errorf("实例 %s: path 不能为空", name)
    }
    
    // 2. 验证文件名
    if c.File == "" {
        return fmt.Errorf("实例 %s: file 不能为空", name)
    }
    
    // 3. 验证格式
    if c.Format != "" && c.Format != "json" && c.Format != "text" {
        return fmt.Errorf("实例 %s: 无效的格式 %s, 有效值: json, text", name, c.Format)
    }
    
    // 4. 检查目录权限
    if err := checkDirectoryPermission(c.Path); err != nil {
        return fmt.Errorf("实例 %s: %v", name, err)
    }
    
    return nil
}
```

**关系**: 属于 `LoggerConfig`,一对多关系

---

## 2. 日志记录模型

### 2.1 LogRecord - 日志记录 (JSON 格式)

**用途**: 表示一条日志记录的完整结构 (JSON 序列化后的格式)

**字段定义**:

```go
// 注意: 这是日志输出的 JSON 结构,不是 Go 结构体
// GoFrame glog 会自动生成这些字段
type LogRecord struct {
    // === 必须字段 (GoFrame 自动添加) ===
    Time    string `json:"Time"`    // 时间戳, 格式由 TimeFormat 配置
    Level   string `json:"Level"`   // 日志级别: DEBUG/INFO/NOTICE/WARNING/ERROR/CRITICAL
    TraceId string `json:"TraceId"` // 请求追踪 ID (从 Context 提取)
    File    string `json:"File"`    // 调用文件和行号, 如 "logger.go:123"
    Content string `json:"Content"` // 日志内容 (主消息)
    
    // === 业务扩展字段 (通过 g.Map 传递) ===
    Module   string      `json:"Module,omitempty"`   // 业务模块: user/order/payment 等
    Action   string      `json:"Action,omitempty"`   // 操作类型: login/create/update 等
    UserId   string      `json:"UserId,omitempty"`   // 用户 ID (从 Context 提取)
    Duration int64       `json:"Duration,omitempty"` // 耗时 (毫秒)
    Error    string      `json:"Error,omitempty"`    // 错误详情
    Extra    interface{} `json:"Extra,omitempty"`    // 扩展信息 (任意结构)
}
```

**字段说明**:

| 字段 | 类型 | 来源 | 说明 |
|------|------|------|------|
| Time | string | GoFrame 自动 | 时间戳,格式由 TimeFormat 配置 |
| Level | string | GoFrame 自动 | 日志级别 |
| TraceId | string | Context (ctxKeys) | 请求追踪 ID,用于关联同一请求的所有日志 |
| File | string | GoFrame 自动 | 调用位置,格式: "文件名:行号" |
| Content | string | 调用参数 | 日志主消息 |
| Module | string | 业务传递 | 业务模块标识 |
| Action | string | 业务传递 | 操作类型 |
| UserId | string | Context (ctxKeys) | 用户 ID |
| Duration | int64 | 业务传递 | 操作耗时 (毫秒) |
| Error | string | 业务传递 | 错误详情 (错误日志) |
| Extra | interface{} | 业务传递 | 扩展信息 (任意 JSON 结构) |

**JSON 示例**:

```json
{
  "Time": "2025-12-12T13:45:30.123+08:00",
  "Level": "INFO",
  "TraceId": "3f2504e0-4f89-11d3-9a0c-0305e82c3301",
  "File": "user_service.go:45",
  "Content": "用户登录成功",
  "Module": "user",
  "Action": "login",
  "UserId": "user123",
  "Duration": 234,
  "Extra": {
    "ip": "192.168.1.100",
    "device": "iPhone 13"
  }
}
```

---

### 2.2 AccessLogEntry - 访问日志条目

**用途**: 表示 HTTP 访问日志的专用结构

**字段定义**:

```go
// 注意: 这是通过 g.Map 传递给 g.Log("access") 的数据结构
type AccessLogEntry struct {
    Method     string `json:"method"`     // HTTP 方法: GET/POST/PUT/DELETE
    Path       string `json:"path"`       // 请求路径: /api/v1/users
    Query      string `json:"query"`      // 查询参数: ?page=1&size=10
    StatusCode int    `json:"status"`     // HTTP 状态码: 200/404/500
    Duration   int64  `json:"duration"`   // 请求耗时 (毫秒)
    IP         string `json:"ip"`         // 客户端 IP
    UserAgent  string `json:"userAgent"`  // User-Agent
    TraceId    string `json:"traceId"`    // 请求追踪 ID
    Size       int    `json:"size"`       // 响应体大小 (字节)
}
```

**使用示例**:

```go
g.Log("access").Info(ctx, "请求完成", g.Map{
    "method": "GET",
    "path": "/api/v1/users",
    "query": "page=1&size=10",
    "status": 200,
    "duration": 123,
    "ip": "192.168.1.100",
    "userAgent": "Mozilla/5.0...",
    "size": 1024,
})
```

---

### 2.3 SlowQueryEntry - 慢查询日志条目

**用途**: 表示数据库慢查询日志的专用结构

**字段定义**:

```go
type SlowQueryEntry struct {
    SQL      string        `json:"sql"`      // SQL 语句
    Args     []interface{} `json:"args"`     // SQL 参数
    Duration int64         `json:"duration"` // 执行耗时 (毫秒)
    Rows     int           `json:"rows"`     // 影响行数
    Error    string        `json:"error,omitempty"` // 错误信息 (如有)
}
```

**使用示例**:

```go
g.Log("slow").Warning(ctx, "慢查询", g.Map{
    "sql": "SELECT * FROM users WHERE status = ?",
    "args": []interface{}{1},
    "duration": 1234,
    "rows": 100,
})
```

---

## 3. 追踪上下文模型

### 3.1 TraceContext - 追踪上下文

**用途**: 表示请求追踪的上下文信息,在 Context 中传递

**字段定义**:

```go
// 注意: 这些字段通过 context.WithValue 存储在 Context 中
// 通过 ctxKeys 配置自动提取到日志
type TraceContext struct {
    TraceID   string // 请求追踪 ID (UUID v4 格式)
    RequestID string // 请求 ID (可选,用于兼容其他系统)
    UserID    string // 用户 ID (认证后设置)
}
```

**Context 操作**:

```go
// 设置 TraceID
ctx = gtrace.WithTraceID(ctx, traceID)

// 设置 UserID
ctx = context.WithValue(ctx, "UserId", userID)

// 获取 TraceID
traceID := gtrace.GetTraceID(ctx)

// 获取 UserID
userID, _ := ctx.Value("UserId").(string)
```

**生命周期**:
- **HTTP 请求**: 从请求开始 (中间件) 到响应结束
- **CLI 命令**: 从命令入口到命令完成
- **后台任务**: 从任务开始到任务结束

---

## 4. 辅助数据结构

### 4.1 LogLevel - 日志级别枚举

**用途**: 定义所有有效的日志级别

```go
const (
    LevelAll      = "all"      // 所有级别
    LevelDev      = "dev"      // 开发模式 (包含 DEBUG)
    LevelProd     = "prod"     // 生产模式 (INFO 及以上)
    LevelDebug    = "debug"    // 调试信息
    LevelInfo     = "info"     // 一般信息
    LevelNotice   = "notice"   // 需要注意的信息
    LevelWarning  = "warning"  // 警告信息
    LevelError    = "error"    // 错误信息
    LevelCritical = "critical" // 严重错误
)

var ValidLevels = []string{
    LevelAll, LevelDev, LevelProd, LevelDebug, LevelInfo,
    LevelNotice, LevelWarning, LevelError, LevelCritical,
}
```

---

### 4.2 LogFormat - 日志格式枚举

**用途**: 定义所有有效的日志格式

```go
const (
    FormatJSON = "json" // JSON 格式 (结构化)
    FormatText = "text" // 纯文本格式
)

var ValidFormats = []string{FormatJSON, FormatText}
```

---

## 5. 数据关系图

```text
LoggerConfig (1)
    │
    ├─── Instances (map) ──> InstanceConfig (N)
    │                            │
    │                            └─── Path, File, Format, Level, Rotate*
    │
    └─── Level, Stdout, TimeFormat, CtxKeys

Context (请求级别)
    │
    ├─── TraceID (gtrace)
    ├─── UserID (context.Value)
    └─── RequestID (context.Value)
         │
         └─── 自动提取到 ──> LogRecord.TraceId, LogRecord.UserId

LogRecord (日志输出)
    │
    ├─── 必须字段: Time, Level, TraceId, File, Content
    └─── 扩展字段: Module, Action, UserId, Duration, Error, Extra

AccessLogEntry (HTTP 访问日志)
    └─── Method, Path, Query, StatusCode, Duration, IP, UserAgent, TraceId, Size

SlowQueryEntry (慢查询日志)
    └─── SQL, Args, Duration, Rows, Error
```

---

## 6. 状态转换

### 6.1 配置加载流程

```text
[配置文件] 
    ↓ (加载)
[LoggerConfig] 
    ↓ (验证)
[验证通过?]
    ├─ 是 → [初始化日志实例]
    └─ 否 → [返回详细错误] → [应用退出]
```

### 6.2 TraceID 生命周期

```text
[HTTP 请求到达]
    ↓
[中间件: 提取或生成 TraceID]
    ↓
[注入到 Context]
    ↓
[传递到 Controller/Service/DAO]
    ↓
[所有日志自动包含 TraceID]
    ↓
[响应返回] → [TraceID 生命周期结束]
```

### 6.3 日志记录流程

```text
[业务代码调用 g.Log().Info(ctx, msg, data)]
    ↓
[GoFrame 提取 Context 中的 ctxKeys]
    ↓
[合并: 必须字段 + ctxKeys + data]
    ↓
[格式化为 JSON/Text]
    ↓
[写入日志文件 (异步/同步)]
    ↓
[检查分割条件] → [是否需要分割?]
        ├─ 是 → [创建新文件]
        └─ 否 → [继续写入当前文件]
```

---

## 7. 验证规则总结

### 7.1 配置验证规则

| 字段 | 验证规则 | 错误示例 |
|------|----------|----------|
| Level | 必须在 ValidLevels 中 | "invalid" → 错误 |
| Path | 必须非空且可写 | "" → 错误, "/root/logs" (无权限) → 错误 |
| File | 必须非空 | "" → 错误 |
| Format | 必须为 json 或 text | "xml" → 错误 |
| RotateSize | 必须为有效大小格式 | "abc" → 错误, "100X" → 错误 |
| RotateExpire | 必须为有效时长格式 | "abc" → 错误, "30x" → 错误 |

### 7.2 运行时验证规则

| 场景 | 验证规则 | 处理策略 |
|------|----------|----------|
| TraceID 丢失 | 检测 Context 中是否存在 TraceID | 生成新 ID + 记录 WARNING |
| 日志文件写入失败 | 检测文件写入错误 | 降级到 stderr + 记录错误 |
| 磁盘空间不足 | 检测写入错误 | 降级到 stderr + 记录错误 |
| 配置热更新 | 检测配置文件变化 | 重新加载配置 (GoFrame 自动) |

---

## 8. 性能考虑

### 8.1 字段大小限制

| 字段 | 最大大小 | 超出处理 |
|------|----------|----------|
| Content | 1KB | 截断 + "... (truncated)" |
| Extra | 5KB | 截断 + "... (truncated)" |
| SQL | 10KB | 截断 + "... (truncated)" |
| Error | 2KB | 截断 + "... (truncated)" |

### 8.2 索引字段 (日志查询优化)

建议对以下字段建立索引 (如使用 ELK/Loki):
- `TraceId`: 用于请求链路查询
- `Level`: 用于按级别过滤
- `Time`: 用于时间范围查询
- `Module`: 用于按模块过滤
- `UserId`: 用于按用户查询

---

**文档版本**: 1.0  
**最后更新**: 2025-12-12  
**审核状态**: 待 Code Review
