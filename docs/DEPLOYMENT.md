# 部署指南

## 环境要求

- Go 1.24.5+
- Node.js 18+
- SQLite3（内置，无需额外服务）
- 必需环境变量：`JWT_SECRET`（≥32 字符）、`DB_PATH`（可选，默认 `backend/data/gostudy.db`）

## 构建步骤

1. 前端构建
   ```bash
   cd frontend
   npm install
   npm run build        # 已启用静态导出，产物在 frontend/out/
   ```
2. 后端编译
   ```bash
   cd backend
   go test ./...
   go build -o ../bin/gostudy main.go
   ```
3. 产物收集
   - 复制 `backend/configs/`（含 config.yaml）
   - 复制 `frontend/out/`（静态文件）
   - 保留 `bin/gostudy` 可执行文件

## 配置要点

- `configs/config.yaml`
  - `http.port` / `https.*` 按需调整
  - `database.path` 指向实际数据目录（如 `./data/gostudy.db`）
  - `static.enabled=true`，`static.path=../frontend/out`，`spaFallback=true`
- 生产环境务必设置 `JWT_SECRET`；建议通过系统环境变量或启动参数传入。

## 启动

```bash
./bin/gostudy -d
# 或设置工作目录后直接运行
JWT_SECRET="your-secret-32chars" DB_PATH="./data/gostudy.db" ./bin/gostudy -d
```

默认监听 `:8080`，静态页与 API 共用端口。

## 验证

- 访问 `http://localhost:8080/` 应能加载前端静态页。
- `curl http://localhost:8080/api/v1/topics` 应返回 `code=20000`。
- 登录后验证进度/测验接口，确保 401/400xx 错误码符合预期。

## 常见问题

- **构建失败**：检查 Node 版本是否 ≥18，`npm install` 是否成功。
- **静态路由 404**：确认 `static.path` 指向 `frontend/out`，`spaFallback=true`。
- **JWT 校验失败**：确认 `JWT_SECRET` 一致且长度足够。***

