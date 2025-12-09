# Quickstart: Variables章节学习完成

1) 构建
- 在仓库根目录执行 `./build.bat`；若不存在则执行 `go test ./...` 与 `go build ./...` 确认可编译。

2) 运行 CLI 学习模式
- `go run ./cmd/cli --chapter variables`
- 使用 `list` 查看子主题，`show static` 查看静态类型示例，`quiz` 进入测验。

3) 运行 HTTP 学习模式
- `go run ./cmd/http --chapter variables --addr :8080`
- `curl http://localhost:8080/api/variables/content?topic=storage`
- `curl -X POST http://localhost:8080/api/variables/quiz/submit -H "Content-Type: application/json" -d '{"answers":[{"id":"q1","choice":"A"}]}'`

4) 测试
- `go test ./src/learning/variables/... ./tests/...`
- 确认单元/契约/集成测试覆盖度>=80%，测验输出与示例保持一致。

5) 文档同步
- 更新 `src/learning/variables/README.md` 与根级 README 中的章节目录，确保双学习模式与菜单导航说明同步。