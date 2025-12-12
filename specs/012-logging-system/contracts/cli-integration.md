# Contract: CLI Integration

**Feature**: 012-logging-system  
**Date**: 2025-12-12  
**Purpose**: 定义 CLI 模式下的日志系统集成规范

## 概述

本文档定义了 CLI 应用如何集成日志系统,包括初始化流程、Context 传递、日志记录方法和错误处理规范。所有 CLI 命令必须遵循本规范以确保日志的一致性和可追踪性。

---

## 1. 应用入口集成

### 1.1 main.go 初始化流程

```go
package main

import (
    "fmt"
    "os"
    "github.com/gogf/gf/v2/os/gctx"
    "go-study2/backend/internal/infrastructure/logger"
)

func main() {
    // 1. 初始化日志系统 (必须在所有业务逻辑之前)
    if err := logger.Initialize(); err != nil {
        // 日志初始化失败,输出到 stderr 并退出
        fmt.Fprintf(os.Stderr, "日志系统初始化失败: %v\n", err)
        os.Exit(1)
    }
    
    // 2. 创建带 TraceID 的根 Context
    ctx := gctx.New()
    
    // 3. 记录应用启动日志
    logger.LogInfo(ctx, "app", "应用启动", g.Map{
        "version": "1.0.0",
        "env": os.Getenv("ENV"),
    })
    
    // 4. 运行应用
    if err := app.Run(ctx); err != nil {
        logger.LogError(ctx, "app", err, g.Map{
            "stage": "运行",
        })
        os.Exit(1)
    }
    
    // 5. 记录应用退出日志
    logger.LogInfo(ctx, "app", "应用正常退出", nil)
}
```

**关键点**:
- 日志初始化失败时必须退出,不允许降级运行
- 使用 `gctx.New()` 创建带 TraceID 的 Context
- 记录应用启动和退出日志,便于追踪应用生命周期

---

## 2. CLI 命令集成

### 2.1 命令处理函数签名

所有 CLI 命令处理函数必须遵循以下签名:

```go
type CommandHandler func(ctx context.Context, args []string) error
```

**参数说明**:
- `ctx`: 包含 TraceID 的 Context,必须传递给所有子函数
- `args`: 命令行参数
- 返回值: 错误信息,nil 表示成功

### 2.2 命令处理示例

```go
package commands

import (
    "context"
    "github.com/gogf/gf/v2/frame/g"
    "go-study2/backend/internal/infrastructure/logger"
)

// 用户登录命令
func HandleLogin(ctx context.Context, args []string) error {
    // 1. 记录命令开始
    logger.LogInfo(ctx, "cli", "执行登录命令", g.Map{
        "args": args,
    })
    
    // 2. 参数验证
    if len(args) < 2 {
        err := errors.New("参数不足: 需要用户名和密码")
        logger.LogError(ctx, "cli", err, g.Map{
            "command": "login",
        })
        return err
    }
    
    username := args[0]
    password := args[1]
    
    // 3. 调用业务逻辑 (传递 Context)
    user, err := userService.Login(ctx, username, password)
    if err != nil {
        logger.LogError(ctx, "user", err, g.Map{
            "username": username,
            "action": "login",
        })
        return err
    }
    
    // 4. 记录成功日志
    logger.LogInfo(ctx, "user", "登录成功", g.Map{
        "userId": user.ID,
        "username": username,
    })
    
    return nil
}
```

---

## 3. Context 传递规范

### 3.1 函数调用链

```text
main()
  ├─ ctx := gctx.New()
  └─ app.Run(ctx)
       └─ command.Execute(ctx, args)
            └─ service.DoSomething(ctx, params)
                 └─ dao.Query(ctx, sql)
```

**规则**:
- 所有函数第一个参数必须是 `ctx context.Context`
- Context 必须从上层传递,不允许在中间层重新创建
- 不允许传递 `nil` Context

### 3.2 Context 值设置

```go
// 设置用户 ID (认证后)
ctx = context.WithValue(ctx, "UserId", userID)

// 设置请求 ID (可选)
ctx = context.WithValue(ctx, "RequestId", requestID)

// 获取 TraceID
traceID := gtrace.GetTraceID(ctx)

// 获取用户 ID
userID, ok := ctx.Value("UserId").(string)
if !ok {
    // 未认证
}
```

---

## 4. 日志记录方法

### 4.1 辅助方法列表

```go
// 信息日志
logger.LogInfo(ctx context.Context, module string, msg string, extra g.Map)

// 错误日志 (自动记录堆栈)
logger.LogError(ctx context.Context, module string, err error, extra g.Map)

// 慢操作日志
logger.LogSlow(ctx context.Context, operation string, duration int64, extra g.Map)

// 业务日志
logger.LogBiz(ctx context.Context, action string, result string, extra g.Map)
```

### 4.2 使用示例

#### 信息日志

```go
logger.LogInfo(ctx, "user", "用户注册", g.Map{
    "username": "zhangsan",
    "email": "zhangsan@example.com",
})
```

#### 错误日志

```go
err := userService.CreateUser(ctx, user)
if err != nil {
    logger.LogError(ctx, "user", err, g.Map{
        "username": user.Username,
        "action": "create",
    })
    return err
}
```

#### 慢操作日志

```go
startTime := time.Now()
result := heavyOperation(ctx)
duration := time.Since(startTime).Milliseconds()

if duration > 1000 {
    logger.LogSlow(ctx, "heavyOperation", duration, g.Map{
        "params": params,
        "result": result,
    })
}
```

#### 业务日志

```go
logger.LogBiz(ctx, "订单创建", "成功", g.Map{
    "orderId": order.ID,
    "amount": order.Amount,
    "userId": order.UserID,
})
```

---

## 5. 错误处理规范

### 5.1 错误记录和返回

```go
func DoSomething(ctx context.Context, params Params) error {
    // 1. 验证参数
    if err := params.Validate(); err != nil {
        logger.LogError(ctx, "validation", err, g.Map{
            "params": params,
        })
        return gerror.Wrap(err, "参数验证失败")
    }
    
    // 2. 执行业务逻辑
    result, err := businessLogic(ctx, params)
    if err != nil {
        logger.LogError(ctx, "business", err, g.Map{
            "params": params,
        })
        return gerror.Wrap(err, "业务逻辑执行失败")
    }
    
    // 3. 记录成功日志
    logger.LogInfo(ctx, "business", "操作成功", g.Map{
        "result": result,
    })
    
    return nil
}
```

**关键点**:
- 使用 `gerror.Wrap()` 包装错误,保留堆栈信息
- 错误日志必须包含 Context (TraceID)
- 错误日志必须包含足够的上下文信息 (参数、状态等)

### 5.2 Panic 处理

```go
func SafeExecute(ctx context.Context, fn func() error) (err error) {
    defer func() {
        if r := recover(); r != nil {
            // 记录 Panic 堆栈
            logger.LogError(ctx, "panic", fmt.Errorf("panic: %v", r), g.Map{
                "stack": string(debug.Stack()),
            })
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    
    return fn()
}
```

---

## 6. 模块标识规范

### 6.1 模块命名

模块名称使用小写英文,多个单词用下划线分隔:

| 模块 | 说明 | 示例 |
|------|------|------|
| app | 应用级别 | 启动、退出、配置加载 |
| cli | CLI 命令 | 命令解析、参数验证 |
| user | 用户模块 | 登录、注册、权限 |
| order | 订单模块 | 创建、支付、取消 |
| payment | 支付模块 | 支付、退款 |
| db | 数据库 | 连接、查询、事务 |
| cache | 缓存 | 读取、写入、失效 |
| external_api | 外部 API | 第三方接口调用 |

### 6.2 操作标识

操作名称使用动词+名词形式:

| 操作 | 说明 |
|------|------|
| create_user | 创建用户 |
| update_profile | 更新资料 |
| delete_order | 删除订单 |
| query_list | 查询列表 |
| send_email | 发送邮件 |

---

## 7. 日志级别使用规范

### 7.1 级别选择

| 级别 | CLI 使用场景 | 示例 |
|------|-------------|------|
| DEBUG | 详细调试信息 (开发环境) | 变量值、函数入参出参 |
| INFO | 关键业务流程节点 | 命令开始、命令完成、用户操作 |
| NOTICE | 需要注意的普通事件 | 配置加载、资源初始化 |
| WARNING | 可恢复的异常 | 参数校验失败、缓存未命中 |
| ERROR | 错误异常 | 命令执行失败、数据库错误 |
| CRITICAL | 严重错误 | 应用无法继续运行 |

### 7.2 使用示例

```go
// DEBUG: 开发环境调试
g.Log().Debug(ctx, "函数调用", g.Map{
    "function": "DoSomething",
    "params": params,
})

// INFO: 关键流程节点
logger.LogInfo(ctx, "user", "用户登录成功", g.Map{
    "userId": user.ID,
})

// WARNING: 可恢复异常
g.Log().Warning(ctx, "缓存未命中", g.Map{
    "key": cacheKey,
})

// ERROR: 错误异常
logger.LogError(ctx, "db", err, g.Map{
    "sql": sql,
})

// CRITICAL: 严重错误
g.Log().Critical(ctx, "数据库连接失败", g.Map{
    "error": err,
})
```

---

## 8. 测试集成

### 8.1 测试中的日志初始化

```go
package commands_test

import (
    "testing"
    "github.com/gogf/gf/v2/os/gctx"
    "go-study2/internal/logic/logger"
)

func TestMain(m *testing.M) {
    // 测试环境日志初始化
    if err := logger.InitializeForTest(); err != nil {
        panic(err)
    }
    
    os.Exit(m.Run())
}

func TestHandleLogin(t *testing.T) {
    ctx := gctx.New()
    
    err := HandleLogin(ctx, []string{"testuser", "testpass"})
    if err != nil {
        t.Errorf("HandleLogin failed: %v", err)
    }
    
    // 验证日志输出 (可选)
    // ...
}
```

### 8.2 测试日志配置

```yaml
# config/config.test.yaml
logger:
  level: "debug"
  stdout: true  # 测试环境输出到终端
  
  default:
    path: "./logs/test"
    file: "test-{Y-m-d}.log"
```

---

## 9. 最佳实践

### 9.1 DO (推荐)

✅ **传递 Context**
```go
func DoSomething(ctx context.Context, params Params) error {
    logger.LogInfo(ctx, "module", "操作开始", g.Map{"params": params})
    // ...
}
```

✅ **使用辅助方法**
```go
logger.LogError(ctx, "module", err, g.Map{"context": "additional info"})
```

✅ **记录关键节点**
```go
logger.LogInfo(ctx, "user", "用户登录成功", g.Map{"userId": user.ID})
```

✅ **敏感信息脱敏**
```go
logger.LogInfo(ctx, "user", "用户注册", g.Map{
    "username": user.Username,
    "password": "***",  // 脱敏
})
```

### 9.2 DON'T (禁止)

❌ **不传递 Context**
```go
func DoSomething(params Params) error {
    g.Log().Info(nil, "操作开始")  // 错误: 缺少 TraceID
}
```

❌ **混用 fmt.Println**
```go
fmt.Println("用户登录成功")  // 错误: 不统一
```

❌ **记录敏感信息**
```go
logger.LogInfo(ctx, "user", "用户登录", g.Map{
    "password": password,  // 错误: 泄露密码
})
```

❌ **在循环中频繁记录**
```go
for _, item := range items {
    logger.LogInfo(ctx, "process", "处理", g.Map{"item": item})  // 错误: 日志过多
}
```

---

## 10. 验收标准

### 10.1 功能验收

- [ ] 所有 CLI 命令在 main.go 中正确初始化日志系统
- [ ] 所有命令处理函数签名包含 `ctx context.Context`
- [ ] 所有日志调用传递 Context
- [ ] 所有错误使用 `logger.LogError()` 记录
- [ ] 所有关键业务节点有日志记录

### 10.2 质量验收

- [ ] 所有日志包含 TraceID
- [ ] 敏感信息已脱敏
- [ ] 模块标识符合规范
- [ ] 日志级别使用正确
- [ ] 测试用例包含日志验证

---

**文档版本**: 1.0  
**最后更新**: 2025-12-12  
**审核状态**: 待 Code Review
