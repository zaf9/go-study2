# 日志系统 (logger)

该包提供基于 GoFrame glog 的统一日志初始化与管理，包含多实例配置、TraceId 支持和日志分割保留策略。

快速开始：

1. 编辑 `backend/configs/config.yaml` 中的 `logger` 配置段，指定实例（default/access/error/slow）和路径。
2. 在应用启动（HTTP 模式）调用：

```go
// 在 backend/main.go 中会自动加载并初始化日志系统
lcfg, err := logger.LoadConfig()
if err != nil {
    // 处理错误并退出
}
if err := logger.Initialize(lcfg); err != nil {
    // 处理错误并退出
}
```

3. 在代码中使用命名实例：

```go
g.Log("access").Info(ctx, "请求完成", g.Map{"path": "/api"})
g.Log("error").Error(ctx, "处理失败", err)
```

## 日志查询功能

日志系统提供基础的日志查询和分析功能，支持按 TraceID、时间范围、日志级别和关键字进行过滤。

### 查询 API

```go
import "go-study2/internal/infrastructure/logger"

// 读取日志文件
entries, err := logger.ReadLogFile("/path/to/app.log")
if err != nil {
    // 处理错误
}

// 按 TraceID 查询
result, err := logger.QueryByTraceID("/path/to/app.log", "trace-123")
if err != nil {
    // 处理错误
}
fmt.Printf("找到 %d 条匹配记录\n", result.Matched)

// 按时间范围查询
start := time.Now().Add(-1 * time.Hour)
end := time.Now()
result, err := logger.QueryByTimeRange("/path/to/app.log", start, end)

// 按日志级别查询
result, err := logger.QueryByLevel("/path/to/app.log", "ERROR")

// 按关键字查询
result, err := logger.QueryByKeyword("/path/to/app.log", "login")

// 列出所有日志文件
files, err := logger.ListLogFiles("/path/to/logs/dir")
```

### 查询结果

所有查询函数返回 `*QueryResult` 结构体：

```go
type QueryResult struct {
    Entries []LogEntry  // 匹配的日志条目
    Total   int         // 日志文件总条目数
    Matched int         // 匹配的条目数
}

type LogEntry struct {
    Timestamp time.Time // 日志时间戳
    Level     string    // 日志级别
    TraceID   string    // TraceID
    Message   string    // 日志消息
    Raw       string    // 原始日志行
}
```

### 性能注意事项

- 查询功能使用流式读取，适合中等大小的日志文件
- 对于大文件（>1GB），建议使用外部工具如 `grep`、`awk` 等
- 时间范围查询依赖日志格式中的时间戳解析

测试：

运行单元测试：

```bash
cd backend
go test ./internal/infrastructure/logger -run TestQuery
```

注：配置验证会在初始化阶段进行，若配置缺失或目录权限不足，应用将以 fail-fast 策略退出。
