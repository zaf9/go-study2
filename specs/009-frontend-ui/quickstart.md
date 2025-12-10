# Quickstart - 009-frontend-ui

## 开发环境
1) 前提：Go 1.24.5、Node 18+、npm；在仓库根目录 `D:\studyspace\go-study\go-study2`。  
2) 配置：`backend/configs/config.yaml` 补充 `auth.jwtSecret`、`database.path`、`static.path=./frontend/out`。  
3) 后端启动：  
   - 优先执行 `./build.bat`（若存在）；否则 `cd backend && go test ./... && go run main.go -d`。  
4) 前端启动：  
   - `cd frontend`  
   - `npm install`  
   - `npm run dev`（默认 3000，使用 next.config.js 代理到 http://localhost:8080/api）。

## 生产构建
1) 前端静态导出：  
   - `cd frontend`  
   - `npm install`  
   - `npm run build && npm run export`（产物输出至 `frontend/out`）。  
2) 后端构建：  
   - 在仓库根目录执行 `./build.bat`；若无则 `cd backend && go test ./... && go build -o ../bin/gostudy main.go`。  
3) 打包与部署：  
   - 将 `backend/configs`、`backend/data`（或初始化为空）与 `frontend/out` 一并复制到部署目录。  
   - 运行 `./gostudy -d`（或后台方式），确认端口 8080 暴露 `/` 与 `/api/*`。

## 验证要点
- 访问 `http://localhost:8080/` 能加载静态页，受保护路由自动跳转登录。  
- 注册/登录后，主题列表与章节内容可正常拉取，错误提示清晰可见。  
- 学习进度页面显示章节完成率，“继续学习”跳转到最近章节。  
- 登录状态过期时触发刷新，失败则清空凭证并回登录页。  
- `go test ./...` 与前端 `npm test` 均需通过。

