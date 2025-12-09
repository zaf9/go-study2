# Variables 章节说明

本目录提供变量学习章节的共享数据结构、内容装载与 CLI/HTTP 双模式入口骨架。

## 目录结构

- `variables.go`：主题枚举、内容与测验基础结构、公共装载与校验函数。
- `content.go`：内容模板与主题注册占位，后续填充具体示例与文案。
- `cli/menu.go`：CLI 端的章节菜单骨架，提供 list/show/quiz 入口。
- `http/handlers.go`：HTTP 端的路由处理骨架，返回 JSON 结构。
- `README.md`：当前文件，说明职责与使用方式。

## 设计要点

- 统一主题枚举：`storage`、`static`、`dynamic`、`zero`，便于多端复用。
- 双模式复用：CLI 与 HTTP 均调用相同的内容与测验数据源，保持一致性。
- 可扩展骨架：当前内容与测验为占位，将在后续任务中补充实际示例与题目。

## 运行与测试

- 构建/测试：`go test ./... && go build ./...`
- 静态检查：`powershell -ExecutionPolicy Bypass -File scripts/check-go.ps1`
- 覆盖率：`go test ./... -cover`（变量章节子包已达 80%+）
- CLI：`go run ./cmd/cli --chapter variables`
- HTTP：`go run ./cmd/http --chapter variables --addr :8080`

