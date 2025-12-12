# Quickstart: 日志系统快速开始指南

**Feature**: 012-logging-system  
**Date**: 2025-12-12  
**Purpose**: 提供日志系统的快速集成和使用指南

## 概述

本指南帮助开发者在 5 分钟内完成日志系统的集成和配置,包括配置文件设置、代码集成、验证测试和常见问题排查。

---

## 1. 配置文件设置 (2 分钟)

### 步骤 1.1: 创建配置文件

```bash
# 进入后端目录
cd backend

# 配置目录已存在,检查配置文件
ls configs/

# 如果需要,复制配置模板
cp configs/config.example.yaml configs/config.yaml
```

### 步骤 1.2: 编辑配置文件

编辑 `backend/configs/config.yaml`:

```yaml
logger:
  level: "all"                    # 开发环境使用 all, 生产环境使用 info
  stdout: true                    # 开发环境输出到终端
  timeFormat: "2006-01-02T15:04:05.000Z07:00"
  ctxKeys: ["TraceId", "UserId"]
  
  default:
    path: "./logs"
    file: "{Y-m-d}.log"
    format: "json"
    rotateSize: "500M"
    rotateBackupExpire: "30d"
    
  access:
    path: "./logs/access"
    file: "access-{Y-m-d}.log"
    
  error:
    path: "./logs/error"
    file: "error-{Y-m-d}.log"
    level: "error"
    
  slow:
    path: "./logs/slow"
    file: "slow-{Y-m-d}.log"
```

### 步骤 1.3: 创建日志目录

```bash
# 在 backend 目录下创建日志目录
mkdir -p logs/access logs/error logs/slow
chmod 755 logs logs/access logs/error logs/slow
```

---

## 2. 代码集成 (2 分钟)

### 步骤 2.1: 在 main.go 中初始化

```go
package main

import (
    "fmt"
    "os"
    "github.com/gogf/gf/v2/os/gctx"
    "github.com/gogf/gf/v2/frame/g"
    "go-study2/backend/internal/infrastructure/logger"
    "go-study2/backend/internal/app/http_server/middleware"
)

func main() {
    // 1. 初始化日志系统
    if err := logger.Initialize(); err != nil {
        fmt.Fprintf(os.Stderr, "日志系统初始化失败: %v\n", err)
        os.Exit(1)
    }
    
    // 2. 创建 Context
    ctx := gctx.New()
    
    // 3. 记录启动日志
    logger.LogInfo(ctx, "app", "应用启动", g.Map{
        "version": "1.0.0",
    })
    
    // 4. 启动 HTTP 服务
    s := g.Server()
    s.Use(middleware.PanicRecovery)
    s.Use(middleware.AccessLog)
    
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("Hello World")
    })
    
    s.Run()
}
```

### 步骤 2.2: 在业务代码中使用

```go
package service

import (
    "context"
    "github.com/gogf/gf/v2/frame/g"
    "go-study2/backend/internal/infrastructure/logger"
)

func CreateUser(ctx context.Context, username string) error {
    // 记录信息日志
    logger.LogInfo(ctx, "user", "创建用户", g.Map{
        "username": username,
    })
    
    // 业务逻辑
    err := dao.Insert(ctx, user)
    if err != nil {
        // 记录错误日志
        logger.LogError(ctx, "user", err, g.Map{
            "username": username,
        })
        return err
    }
    
    return nil
}
```

---

## 3. 验证和测试 (1 分钟)

### 步骤 3.1: 启动应用

```bash
# 在 backend 目录下运行
cd backend
go run main.go
```

### 步骤 3.2: 发送测试请求

```bash
curl http://localhost:8080/
```

### 步骤 3.3: 检查日志文件

```bash
# 查看应用日志
cat logs/$(date +%Y-%m-%d).log | jq

# 查看访问日志
cat logs/access/access-$(date +%Y-%m-%d).log | jq

# 验证 TraceID
cat logs/access/access-$(date +%Y-%m-%d).log | jq '.TraceId'
```

**预期输出**:

```json
{
  "Time": "2025-12-12T14:00:00.123+08:00",
  "Level": "INFO",
  "TraceId": "3f2504e0-4f89-11d3-9a0c-0305e82c3301",
  "File": "access_log.go:45",
  "Content": "请求完成",
  "method": "GET",
  "path": "/",
  "status": 200,
  "duration": 12
}
```

---

## 4. 常见问题排查

### 问题 1: 应用启动失败 - "配置文件不存在"

**错误信息**:
```
日志系统初始化失败: 配置文件不存在: backend/configs/config.yaml
```

**解决方法**:
```bash
# 检查配置文件是否存在
ls -la backend/configs/config.yaml

# 如果不存在,创建配置文件
cd backend
cp configs/config.example.yaml configs/config.yaml
```

---

### 问题 2: 应用启动失败 - "日志目录无写入权限"

**错误信息**:
```
日志系统初始化失败: 实例 default: 日志目录 ./logs 无写入权限
```

**解决方法**:
```bash
# 在 backend 目录下创建日志目录
cd backend
mkdir -p logs/access logs/error logs/slow

# 修改权限
chmod 755 logs logs/access logs/error logs/slow

# 或以 root 运行 (不推荐)
sudo go run main.go
```

---

### 问题 3: 日志文件未生成

**可能原因**:
1. 日志级别过高,过滤了所有日志
2. 日志目录权限不足
3. 磁盘空间不足

**排查步骤**:
```bash
# 1. 检查日志级别配置
grep "level:" backend/configs/config.yaml

# 2. 检查目录权限
ls -la backend/logs/

# 3. 检查磁盘空间
df -h

# 4. 检查应用日志输出 (stdout)
cd backend && go run main.go 2>&1 | grep -i log
```

---

### 问题 4: TraceID 未传递

**症状**: 日志中 TraceId 字段为空

**解决方法**:
```go
// 确保所有函数传递 Context
func DoSomething(ctx context.Context, params Params) error {
    // 正确: 传递 ctx
    logger.LogInfo(ctx, "module", "message", nil)
    
    // 错误: 不传递 ctx
    // logger.LogInfo(nil, "module", "message", nil)
}

// 确保 HTTP 中间件注册顺序正确
s.Use(middleware.AccessLog)  // 必须在业务 Handler 之前
```

---

### 问题 5: 日志文件未分割

**症状**: 日志文件持续增长,未按日期分割

**解决方法**:
```yaml
# 检查配置文件中的 file 字段
default:
  file: "{Y-m-d}.log"  # 正确: 使用日期占位符
  # file: "app.log"    # 错误: 固定文件名
```

---

### 问题 6: 慢查询日志未记录

**症状**: 明显慢查询未出现在 slow 日志中

**解决方法**:
```yaml
# 检查慢查询阈值配置
database:
  slow:
    threshold: 1000  # 单位: 毫秒

# 或降低阈值用于测试
database:
  slow:
    threshold: 100  # 100ms
```

---

## 5. 配置示例

### 开发环境配置 (`config.dev.yaml`)

```yaml
logger:
  level: "all"
  stdout: true  # 输出到终端
  
  default:
    path: "./logs"
    file: "{Y-m-d}.log"
    format: "json"
```

### 测试环境配置 (`config.test.yaml`)

```yaml
logger:
  level: "debug"
  stdout: false
  
  default:
    path: "./logs/test"
    file: "test-{Y-m-d}.log"
```

### 生产环境配置 (`config.prod.yaml`)

```yaml
logger:
  level: "info"  # 只记录 INFO 及以上
  stdout: false
  
  default:
    path: "/var/log/go-study2"
    file: "{Y-m-d}.log"
    format: "json"
    rotateSize: "500M"
    rotateBackupExpire: "30d"
    rotateBackupCompress: 9  # 压缩备份文件
```

---

## 6. 使用示例

### 示例 1: 记录用户登录

```go
func Login(ctx context.Context, username, password string) error {
    logger.LogInfo(ctx, "user", "用户登录尝试", g.Map{
        "username": username,
    })
    
    user, err := dao.GetUserByUsername(ctx, username)
    if err != nil {
        logger.LogError(ctx, "user", err, g.Map{
            "username": username,
            "action": "查询用户",
        })
        return err
    }
    
    if !verifyPassword(password, user.Password) {
        logger.LogInfo(ctx, "user", "密码错误", g.Map{
            "username": username,
        })
        return errors.New("密码错误")
    }
    
    logger.LogInfo(ctx, "user", "登录成功", g.Map{
        "userId": user.ID,
        "username": username,
    })
    
    return nil
}
```

### 示例 2: 记录慢操作

```go
func HeavyOperation(ctx context.Context, params Params) error {
    startTime := time.Now()
    
    // 执行耗时操作
    result := doHeavyWork(params)
    
    duration := time.Since(startTime).Milliseconds()
    if duration > 1000 {
        logger.LogSlow(ctx, "heavyOperation", duration, g.Map{
            "params": params,
            "result": result,
        })
    }
    
    return nil
}
```

### 示例 3: 记录业务流程

```go
func CreateOrder(ctx context.Context, order Order) error {
    logger.LogBiz(ctx, "订单创建", "开始", g.Map{
        "userId": order.UserID,
        "amount": order.Amount,
    })
    
    // 业务逻辑
    if err := dao.Insert(ctx, order); err != nil {
        logger.LogBiz(ctx, "订单创建", "失败", g.Map{
            "userId": order.UserID,
            "error": err.Error(),
        })
        return err
    }
    
    logger.LogBiz(ctx, "订单创建", "成功", g.Map{
        "orderId": order.ID,
        "userId": order.UserID,
        "amount": order.Amount,
    })
    
    return nil
}
```

---

## 7. 下一步

- ✅ **完成**: 日志系统已集成
- ⏭️ **建议**: 阅读 [logger 包 README](../../../backend/internal/infrastructure/logger/README.md)
- ⏭️ **建议**: 阅读 [日志最佳实践文档](./best-practices.md)
- ⏭️ **建议**: 配置日志采集系统 (ELK/Loki)

---

**文档版本**: 1.0  
**最后更新**: 2025-12-12  
**预计完成时间**: 5 分钟
