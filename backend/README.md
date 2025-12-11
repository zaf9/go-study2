# 后端架构与运行说明

## 目录结构

- `go.mod` / `go.sum`：模块与依赖管理  
- `main.go` / `main_test.go`：入口与顶层测试  
- `configs/`：`config.yaml`、证书等配置  
- `internal/`：配置加载、HTTP 服务、学习内容（lexical_elements、constants 等）  
- `src/`：CLI/HTTP 复用的章节内容与测验逻辑  
- `tests/`：unit / integration / contract 测试  
- `scripts/`：`check-go.ps1`（gofmt、go vet、golint 统一入口）  
- `doc/`、`docs/`：文档材料与章节 quickstart  

## 构建与运行

```bash
# 在仓库根执行
cd backend

# 构建/运行
go build ./...         # 或 go run main.go
go run main.go -d      # 启动 HTTP 模式（默认 127.0.0.1:8080）

# 测试
go test ./...
```

## 配置

- 默认配置位于 `configs/config.yaml`，HTTPS 证书放置于 `configs/certs/`。  
- 配置加载通过 `internal/config`，工作目录以 `backend/` 为基准。  

## 默认管理员

- 初始账号：`admin` / `GoStudy@123`。  
- 首次登录会被强制修改密码；改密后旧口令与旧令牌全部失效。

## API 速览

- 主题列表：`GET /api/v1/topics?format=json|html`
- 词法元素菜单：`GET /api/v1/topic/lexical_elements`
- 词法元素子主题：`GET /api/v1/topic/lexical_elements/{chapter}`
- 常量菜单：`GET /api/v1/topic/constants`
- 常量子主题：`GET /api/v1/topic/constants/{subtopic}`

## 开发辅助

- 代码检查：`powershell -ExecutionPolicy Bypass -File backend/scripts/check-go.ps1`
- 覆盖率示例：`go test -cover ./...`

