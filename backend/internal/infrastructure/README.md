# 基础设施层说明

提供数据库连接与仓储实现，支撑领域服务落地。

## 目录

- `database/`：SQLite 初始化与迁移（WAL、busy_timeout、索引）。
- `repository/`：用户、进度、测验仓储实现，基于 GoFrame ORM。

## 配置要点

- 数据库路径与 PRAGMA 在 `configs/config.yaml` 中配置。
- 启动时调用 `database.Init` 完成目录创建、PRAGMA 设置与建表。
- 仓储实现依赖领域接口，保持单一职责。

## 测试

- 集成测试使用临时数据库文件隔离数据。
- 推荐在 CI 中执行 `go test ./...` 确认仓储与迁移可用。*** End Patch

