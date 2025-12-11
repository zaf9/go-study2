# Quickstart: Go-Study2 学习闭环与测验体系

**Branch**: 011-learning-progress-quiz  
**Spec**: specs/011-learning-progress-quiz/spec.md

## 前置
- Node 18+，Go 1.24.5，SQLite3。  
- 先执行项目根 `./build.bat`（Windows 下可直接双击或在 PowerShell 运行）。  
- 确保 backend/.env（或 config.yaml）中数据库路径可写。

## Backend
```powershell
cd backend
go mod tidy
go fmt ./...
go vet ./...
go test ./...
go run ./cmd/main.go
```

## Frontend
```powershell
cd frontend
npm install
npm run lint
npm test
npm run dev   # 开发
# npm run build && npm run start # 生产预览
```

## 接口与页面
- API 主要路由：`/api/v1/progress`、`/api/v1/quiz/{topic}/{chapter}`、`/api/v1/quiz/history`。  
- 页面：`/progress` 学习进度、`/topics/[topic]` 主题列表与继续学习、`/topics/[topic]/[chapter]` 章节详情与测验入口、`/quiz/[topic]/[chapter]` 测验、`/quiz` 测验历史。

## 测试目标
- 后端与前端核心功能覆盖率 ≥80%。  
- 关键用例：进度同步（防抖+指数退避+卸载前同步）、阅读恢复、章节完成判定、测验抽题与判分、测验历史与继续学习跳转。

