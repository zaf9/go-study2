# Variables 章节 Quickstart

## 构建与检查
- `./build.bat` 如存在优先执行；否则运行 `go test ./... && go build ./...`
- 静态检查：`powershell -ExecutionPolicy Bypass -File scripts/check-go.ps1`
- 覆盖率：`go test ./... -cover`

## CLI 模式
- 启动：`go run main.go`，主菜单选择 `2` 进入 Variables
- 子菜单 0-3 覆盖 `storage|static|dynamic|zero`，展示内容和测验（含答案解析）

## HTTP 模式
- 启动：`go run main.go -d`（默认 127.0.0.1:8080）
- API：
  - `GET /api/v1/topic/variables?format=html`（菜单）
  - `GET /api/v1/topic/variables/{storage|static|dynamic|zero}`
  - 内容页面包含测验题目与答案解析，当前未暴露提交接口

## 数据一致性
- CLI 与 HTTP 复用同一内容/测验数据源，确保题目与解析一致。

## 常见问题
- 返回 400：检查 `topic` 是否为 `storage|static|dynamic|zero`，或提交体 JSON 是否正确。
- 返回 501：该主题尚未注册测验数据（当前四个主题均已提供）。***

