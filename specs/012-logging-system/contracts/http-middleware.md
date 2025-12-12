# Contract: HTTP Middleware

**Feature**: 012-logging-system  
**Date**: 2025-12-12  
**Purpose**: 定义 HTTP 中间件的接口规范和实现要求

## 概述

本文档定义了 HTTP 日志中间件的完整规范,包括访问日志中间件、Panic 恢复中间件、数据库日志 Handler 的接口、执行顺序、Context 传递和错误处理机制。

---

## 1. 中间件注册和执行顺序

### 1.1 中间件注册

```go
package main

import (
    "github.com/gogf/gf/v2/net/ghttp"
    "go-study2/backend/internal/app/http_server/middleware"
)

func main() {
    s := g.Server()
    
    // 全局中间件 (按注册顺序执行)
    s.Use(middleware.PanicRecovery)  // 1. 最外层: Panic 恢复
    s.Use(middleware.AccessLog)      // 2. 访问日志 + TraceID 注入
    s.Use(middleware.Authentication) // 3. 认证 (可选)
    s.Use(middleware.Authorization)  // 4. 授权 (可选)
    
    // 路由注册
    s.Group("/api/v1", func(group *ghttp.RouterGroup) {
        // ...
    })
    
    s.Run()
}
```

### 1.2 执行顺序

```text
请求流程:
PanicRecovery (前) 
  → AccessLog (前) 
    → Authentication (前) 
      → Authorization (前) 
        → Handler 
      ← Authorization (后) 
    ← Authentication (后) 
  ← AccessLog (后) 
← PanicRecovery (后)
```

**关键点**:
- Panic Recovery 必须在最外层,确保所有 Panic 都被捕获
- Access Log 必须在第二层,确保所有请求都被记录 (即使 Panic)
- TraceID 注入在 Access Log 中完成,确保后续中间件和 Handler 都能访问

---

## 2. 访问日志中间件

### 2.1 接口定义

```go
package middleware

import (
    "time"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/os/gtrace"
    "github.com/gogf/gf/v2/frame/g"
    "go-study2/backend/internal/infrastructure/logger"
)

// AccessLog 访问日志中间件
func AccessLog(r *ghttp.Request) {
    // 1. 提取或生成 TraceID
    traceID := extractOrGenerateTraceID(r)
    
    // 2. 注入到 Context
    ctx := gtrace.WithTraceID(r.Context(), traceID)
    r.SetCtx(ctx)
    
    // 3. 记录请求开始
    startTime := time.Now()
    g.Log("access").Info(ctx, "请求开始", g.Map{
        "method": r.Method,
        "path": r.URL.Path,
        "query": r.URL.RawQuery,
        "ip": r.GetClientIp(),
        "userAgent": r.Header.Get("User-Agent"),
        "referer": r.Header.Get("Referer"),
    })
    
    // 4. 执行后续处理
    r.Middleware.Next()
    
    // 5. 记录请求结束
    duration := time.Since(startTime).Milliseconds()
    statusCode := r.Response.Status
    
    logData := g.Map{
        "status": statusCode,
        "duration": duration,
        "size": r.Response.BufferLength(),
    }
    
    // 根据状态码选择日志级别
    if statusCode >= 500 {
        g.Log("error").Error(ctx, "请求失败 (5xx)", logData)
    } else if statusCode >= 400 {
        g.Log("access").Warning(ctx, "请求失败 (4xx)", logData)
    } else {
        g.Log("access").Info(ctx, "请求完成", logData)
    }
}
```

### 2.2 TraceID 提取逻辑

```go
// extractOrGenerateTraceID 提取或生成 TraceID
func extractOrGenerateTraceID(r *ghttp.Request) string {
    // 1. 优先从 Header 读取 (支持分布式追踪)
    traceID := r.Header.Get("X-Trace-Id")
    if traceID != "" {
        return traceID
    }
    
    // 2. 备用 Header
    traceID = r.Header.Get("X-Request-Id")
    if traceID != "" {
        return traceID
    }
    
    // 3. 生成新的 TraceID
    return gtrace.NewID()
}
```

### 2.3 请求日志字段

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| method | string | HTTP 方法 | "GET", "POST" |
| path | string | 请求路径 | "/api/v1/users" |
| query | string | 查询参数 | "page=1&size=10" |
| status | int | HTTP 状态码 | 200, 404, 500 |
| duration | int64 | 请求耗时 (毫秒) | 123 |
| ip | string | 客户端 IP | "192.168.1.100" |
| userAgent | string | User-Agent | "Mozilla/5.0..." |
| referer | string | Referer | "https://example.com" |
| size | int | 响应体大小 (字节) | 1024 |
| traceId | string | 请求追踪 ID | "3f2504e0-..." |

---

## 3. Panic 恢复中间件

### 3.1 接口定义

```go
package middleware

import (
    "runtime/debug"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/os/gtrace"
    "github.com/gogf/gf/v2/frame/g"
)

// PanicRecovery Panic 恢复中间件
func PanicRecovery(r *ghttp.Request) {
    defer func() {
        if err := recover(); err != nil {
            // 1. 记录 Panic 详情到错误日志
            g.Log("error").Error(r.Context(), "Panic 恢复", g.Map{
                "error": err,
                "stack": string(debug.Stack()),
                "method": r.Method,
                "path": r.URL.Path,
                "query": r.URL.RawQuery,
                "ip": r.GetClientIp(),
            })
            
            // 2. 返回统一错误响应
            r.Response.ClearBuffer()
            r.Response.WriteStatus(500)
            r.Response.WriteJson(g.Map{
                "code": 500,
                "message": "服务器内部错误",
                "traceId": gtrace.GetTraceID(r.Context()),
            })
        }
    }()
    
    r.Middleware.Next()
}
```

### 3.2 Panic 日志字段

| 字段 | 类型 | 说明 |
|------|------|------|
| error | interface{} | Panic 错误信息 |
| stack | string | 完整堆栈信息 |
| method | string | HTTP 方法 |
| path | string | 请求路径 |
| query | string | 查询参数 |
| ip | string | 客户端 IP |
| traceId | string | 请求追踪 ID |

---

## 4. 数据库日志 Handler

### 4.1 ORM Handler 注册

```go
package main

import (
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "go-study2/backend/internal/app/http_server/middleware"
)

func initDatabase() {
    db := g.DB()
    
    // 注册全局 Hook
    db.GetCore().AddHook(gdb.Hook{
        Select: middleware.DBLogHandler,
        Insert: middleware.DBLogHandler,
        Update: middleware.DBLogHandler,
        Delete: middleware.DBLogHandler,
    })
}
```

### 4.2 Handler 实现

```go
package middleware

import (
    "context"
    "time"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/genv"
)

// DBLogHandler 数据库日志 Handler
func DBLogHandler(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
    startTime := time.Now()
    
    // 执行 SQL
    result, err = in.Next(ctx)
    
    // 计算耗时
    duration := time.Since(startTime).Milliseconds()
    
    // 构建日志数据
    logData := g.Map{
        "sql": in.SQL,
        "args": in.Args,
        "duration": duration,
        "rows": 0,
    }
    
    if result != nil {
        logData["rows"] = result.Len()
    }
    
    // 根据结果和耗时选择日志级别和实例
    if err != nil {
        // 错误日志
        g.Log("error").Error(ctx, "SQL 执行失败", logData)
    } else if duration > getSlowQueryThreshold() {
        // 慢查询日志 (默认 >1000ms)
        g.Log("slow").Warning(ctx, "慢查询", logData)
    } else {
        // 正常日志 (仅开发环境)
        if genv.Get("ENV") == "dev" {
            g.Log().Debug(ctx, "SQL 执行", logData)
        }
    }
    
    return result, err
}

// getSlowQueryThreshold 获取慢查询阈值 (毫秒)
func getSlowQueryThreshold() int64 {
    threshold := g.Cfg().MustGet(ctx, "database.slow.threshold").Int64()
    if threshold == 0 {
        return 1000  // 默认 1 秒
    }
    return threshold
}
```

### 4.3 慢查询配置

```yaml
# config/config.yaml
database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/db"
    debug: false  # 生产环境关闭 debug
  
  slow:
    threshold: 1000  # 慢查询阈值 (毫秒)
```

### 4.4 SQL 日志字段

| 字段 | 类型 | 说明 |
|------|------|------|
| sql | string | SQL 语句 |
| args | []interface{} | SQL 参数 |
| duration | int64 | 执行耗时 (毫秒) |
| rows | int | 影响/返回行数 |
| error | string | 错误信息 (如有) |
| traceId | string | 请求追踪 ID |

---

## 5. Context 传递机制

### 5.1 Context 流转

```text
[HTTP 请求到达]
    ↓
[PanicRecovery 中间件]
    ↓
[AccessLog 中间件: 注入 TraceID]
    ↓ r.SetCtx(ctx)
[Authentication 中间件: 注入 UserID]
    ↓ r.SetCtx(ctx)
[Handler: 获取 Context]
    ↓ ctx := r.Context()
[Service/DAO: 传递 Context]
    ↓
[数据库 Handler: 使用 Context]
    ↓
[日志自动包含 TraceID/UserID]
```

### 5.2 Context 值设置

```go
// 在 Authentication 中间件中设置 UserID
func Authentication(r *ghttp.Request) {
    // 1. 验证 Token
    userID, err := validateToken(r.Header.Get("Authorization"))
    if err != nil {
        r.Response.WriteStatus(401)
        return
    }
    
    // 2. 注入 UserID 到 Context
    ctx := context.WithValue(r.Context(), "UserId", userID)
    r.SetCtx(ctx)
    
    r.Middleware.Next()
}
```

### 5.3 Context 值获取

```go
// 在 Handler 中获取 TraceID 和 UserID
func UserHandler(r *ghttp.Request) {
    ctx := r.Context()
    
    traceID := gtrace.GetTraceID(ctx)
    userID, _ := ctx.Value("UserId").(string)
    
    logger.LogInfo(ctx, "user", "处理用户请求", g.Map{
        "traceId": traceID,
        "userId": userID,
    })
}
```

---

## 6. 错误处理规范

### 6.1 HTTP 错误响应

所有错误响应必须包含 TraceID:

```go
// 统一错误响应格式
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    TraceId string `json:"traceId"`
}

// 返回错误响应
func respondError(r *ghttp.Request, code int, message string) {
    r.Response.WriteStatus(code)
    r.Response.WriteJson(ErrorResponse{
        Code:    code,
        Message: message,
        TraceId: gtrace.GetTraceID(r.Context()),
    })
}
```

### 6.2 错误日志记录

```go
// 业务错误
if err := service.DoSomething(ctx, params); err != nil {
    logger.LogError(ctx, "service", err, g.Map{
        "params": params,
    })
    respondError(r, 500, "操作失败")
    return
}

// 参数验证错误
if err := params.Validate(); err != nil {
    logger.LogError(ctx, "validation", err, g.Map{
        "params": params,
    })
    respondError(r, 400, "参数无效")
    return
}
```

---

## 7. 日志级别映射

### 7.1 HTTP 状态码与日志级别

| 状态码范围 | 日志级别 | 日志实例 | 说明 |
|-----------|---------|---------|------|
| 200-299 | INFO | access | 成功请求 |
| 300-399 | INFO | access | 重定向 |
| 400-499 | WARNING | access | 客户端错误 |
| 500-599 | ERROR | error | 服务器错误 |

### 7.2 特殊状态码处理

```go
switch statusCode {
case 401:
    g.Log("access").Warning(ctx, "未认证", logData)
case 403:
    g.Log("access").Warning(ctx, "无权限", logData)
case 404:
    g.Log("access").Info(ctx, "资源不存在", logData)
case 429:
    g.Log("access").Warning(ctx, "请求过于频繁", logData)
case 500:
    g.Log("error").Error(ctx, "服务器内部错误", logData)
case 502, 503, 504:
    g.Log("error").Error(ctx, "服务不可用", logData)
}
```

---

## 8. 性能优化

### 8.1 异步日志

```go
// 生产环境启用异步日志
if genv.Get("ENV") == "prod" {
    g.Log("access").SetAsync(true)
    g.Log("error").SetAsync(false)  // 错误日志保持同步
}
```

### 8.2 日志采样

```go
// 高并发场景下采样 (仅记录 10% 的 DEBUG 日志)
if shouldSample(0.1) {
    g.Log().Debug(ctx, "详细调试信息", data)
}

func shouldSample(rate float64) bool {
    return rand.Float64() < rate
}
```

### 8.3 大对象截断

```go
const MaxLogSize = 5 * 1024  // 5KB

// 截断大对象
func truncateLog(data interface{}) string {
    str := gconv.String(data)
    if len(str) > MaxLogSize {
        return str[:MaxLogSize] + "... (truncated)"
    }
    return str
}

// 使用
g.Log().Info(ctx, "大对象日志", g.Map{
    "data": truncateLog(largeObject),
})
```

---

## 9. 测试集成

### 9.1 中间件测试

```go
package middleware_test

import (
    "testing"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/test/gtest"
    "go-study2/internal/middleware"
)

func TestAccessLog(t *testing.T) {
    gtest.C(t, func(t *gtest.T) {
        s := g.Server()
        s.Use(middleware.AccessLog)
        s.BindHandler("/test", func(r *ghttp.Request) {
            r.Response.Write("OK")
        })
        s.Start()
        defer s.Shutdown()
        
        // 发送测试请求
        client := g.Client()
        resp, err := client.Get(ctx, "http://127.0.0.1:8080/test")
        t.AssertNil(err)
        t.Assert(resp.StatusCode, 200)
        
        // 验证日志输出
        // ...
    })
}
```

---

## 10. 最佳实践

### 10.1 DO (推荐)

✅ **使用统一错误响应**
```go
respondError(r, 400, "参数无效")
```

✅ **记录完整请求信息**
```go
g.Log("access").Info(ctx, "请求开始", g.Map{
    "method": r.Method,
    "path": r.URL.Path,
    "ip": r.GetClientIp(),
})
```

✅ **Panic 记录堆栈**
```go
g.Log("error").Error(ctx, "Panic", g.Map{
    "stack": string(debug.Stack()),
})
```

### 10.2 DON'T (禁止)

❌ **不记录敏感 Header**
```go
// 错误: 记录 Authorization Header
g.Log().Info(ctx, "请求", g.Map{
    "authorization": r.Header.Get("Authorization"),
})
```

❌ **不记录完整请求体**
```go
// 错误: 可能包含敏感信息
g.Log().Info(ctx, "请求", g.Map{
    "body": r.GetBody(),
})
```

❌ **不在中间件中阻塞**
```go
// 错误: 同步写入大量日志
for i := 0; i < 1000; i++ {
    g.Log().Info(ctx, "循环日志")
}
```

---

## 11. 验收标准

### 11.1 功能验收

- [ ] 所有 HTTP 请求有完整访问日志 (请求+响应)
- [ ] 所有请求包含 TraceID
- [ ] Panic 被正确捕获并记录堆栈
- [ ] 数据库慢查询被正确记录
- [ ] 错误响应包含 TraceID

### 11.2 性能验收

- [ ] 1000 并发请求下日志开销 <10%
- [ ] 异步日志写入不阻塞请求
- [ ] 大对象日志被正确截断

### 11.3 安全验收

- [ ] 敏感 Header 未被记录
- [ ] 请求体未被完整记录
- [ ] 密码等敏感信息已脱敏

---

**文档版本**: 1.0  
**最后更新**: 2025-12-12  
**审核状态**: 待 Code Review
