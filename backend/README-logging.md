# 日志系统使用说明（节选）

本文档为项目日志系统的快速说明，包含配置示例与常用操作。

配置示例（在 `backend/configs/config.dev.yaml` 中）:

logger:
  timeFormat: "02/Jan/2006:15:04:05 -0700"
  instances:
    access:
      path: "backend/logs/access"
      level: "info"
    error:
      path: "backend/logs/error"
      level: "error"
    slow:
      path: "backend/logs/slow"
      level: "warn"

运行压测示例：

```
go run backend/scripts/stress_client.go -url http://localhost:8080/ -concurrency 100 -requests 1000
```

查询日志示例（代码）：使用 `backend/internal/infrastructure/logger.QueryByTraceID(traceID)`。
