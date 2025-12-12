# 快速开始验证步骤

此文档说明如何验证 `backend/docs/quickstart.md` 中的快速开始流程已在本地正确运行，主要用于任务 T076 的自动化/手动验证。

前提：已按照仓库根目录 README 中的启动说明在本地启动服务（例如：`go run main.go` 或已生成并运行 `..\\bin\\gostudy.exe`）。

验证要点：

1. 环境与配置
   - 确认使用的配置文件位于 `backend/configs/config.dev.yaml`（或相应环境配置）。
   - 验证 `logger` 配置段存在且 `access`、`error`、`slow` 日志目录已创建。

2. 启动服务并检查日志输出
   - 启动服务后，发送一个简单的 HTTP 请求（例如：GET /）。
   - 检查 `backend/logs/access/` 中是否生成访问日志，并包含 TraceID 与请求路径。
   - 检查 `backend/logs/error/` 在触发错误时是否记录堆栈信息。

3. TraceID 全链路验证
   - 使用 HTTP 客户端在请求头中设置 `TraceId: test-trace-123`，发送请求。
   - 在 `access` 与 `biz` 日志中搜索该 TraceID，验证能找到对应的开始/结束/业务日志条目。

4. 数据库慢查询与 DB 日志验证
   - 在配置中设置慢查询阈值（例如：100ms），执行一个预计耗时超过阈值的查询（或使用测试用例模拟）。
   - 验证 `backend/logs/slow/` 中生成慢查询日志，并包含 SQL 与耗时信息。

5. 日志查询工具验证
   - 使用代码或 `logger` 包提供的查询函数按 TraceID / 时间范围进行检索，确保返回正确的日志条目。

6. 性能与稳定性检查
   - 使用 `backend/scripts/stress_client.go`（详见仓库）进行并发请求压测，观察日志生成与磁盘、CPU 的占用情况。

若任一步骤失败，请参阅 `backend/docs/logging_troubleshooting.md` 进行排查。
