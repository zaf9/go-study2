# Quickstart: Variables章节学习完成

1) 构建
- 优先执行 `./build.bat`；若不存在则运行 `go test ./...` 与 `go build ./...`。

2) 运行 CLI 学习模式
- `go run main.go`，主菜单选择 `2` 进入 Variables
- 在子菜单选择 0-3 查看 `storage|static|dynamic|zero`，展示内容与测验

3) 运行 HTTP 学习模式
- `go run main.go -d` （默认 127.0.0.1:8080）
- `curl http://localhost:8080/api/v1/topic/variables?format=html`
- `curl http://localhost:8080/api/v1/topic/variables/storage`
- （内容页附带测验题目；提交接口暂未在 HTTP 模式暴露）

4) 测试
- `go test ./src/learning/variables/... ./tests/...`
- 覆盖率已达 80%+，确保测验输出与示例保持一致。

5) 文档同步
- `docs/quickstart-variables.md` 记录章节独立启动方式；如有变更同步根级 README 与 `src/learning/variables/README.md`。