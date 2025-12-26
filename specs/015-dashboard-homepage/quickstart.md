# Quick Start Guide: Dashboard 首页功能

**Feature**: Dashboard 首页  
**Date**: 2025-12-26  
**Related**: [spec.md](./spec.md) | [plan.md](./plan.md)

## Overview

本指南帮助开发者快速搭建和运行 Dashboard 首页功能。

## Prerequisites (前置条件)

### 环境要求

- **Node.js**: >= 18.0.0
- **Go**: >= 1.24.5
- **npm**: >= 9.0.0
- **Git**: 已安装并配置

### 已安装依赖

- Next.js 14.2.15
- Ant Design 5.x
- GoFrame v2.9.5
- 数据库（PostgreSQL/MySQL/SQLite）

## Quick Start (快速开始)

### Step 1: 拉取最新代码

```bash
# 切换到 Dashboard 功能分支
git checkout 015-dashboard-homepage

# 拉取最新代码
git pull origin 015-dashboard-homepage
```

### Step 2: 安装依赖

```bash
# 前端依赖
cd frontend
npm install

# 后端依赖（如有新增）
cd ../backend
go mod tidy
```

### Step 3: 启动后端服务

```bash
cd backend
go run main.go
```

**预期输出**:
```
[INFO] Server started on :8080
[INFO] WebSocket hub started
[INFO] Database connected
```

### Step 4: 启动前端开发服务器

```bash
cd frontend
npm run dev
```

**预期输出**:
```
▲ Next.js 14.2.15
- Local:        http://localhost:3000
- Ready in 2.1s
```

### Step 5: 访问 Dashboard

打开浏览器访问: `http://localhost:3000/dashboard`

## Development Workflow (开发工作流)

### 创建新组件

```bash
# 在 frontend/app/(protected)/dashboard/components/ 目录下创建组件
cd frontend/app/(protected)/dashboard/components
touch WelcomeHeader.tsx
```

**组件模板**:
```typescript
'use client'

import React from 'react'
import { Typography } from 'antd'

const { Title, Text } = Typography

interface WelcomeHeaderProps {
  username: string
  studyDays: number
}

export const WelcomeHeader: React.FC<WelcomeHeaderProps> = ({ username, studyDays }) => {
  return (
    <div className="welcome-header">
      <Title level={2}>欢迎回来，{username}！</Title>
      <Text type="secondary">您已累计学习 {studyDays} 天</Text>
    </div>
  )
}
```

### 添加 API 调用

```bash
# 在 frontend/lib/api.ts 中添加 API 函数
```

**API 函数模板**:
```typescript
export async function getLastLearning(): Promise<LastLearningData | null> {
  try {
    const response = await api.get<ApiResponse<LastLearningData | null>>(
      '/api/v1/progress/last'
    )
    
    if (response.data.code !== 0) {
      throw new Error(response.data.message)
    }
    
    return response.data.data
  } catch (error) {
    console.error('Failed to fetch last learning:', error)
    throw error
  }
}
```

### 运行测试

```bash
# 前端测试
cd frontend
npm test

# 后端测试
cd backend
go test ./...
```

### 代码格式化

```bash
# 前端
cd frontend
npm run lint
npm run format

# 后端
cd backend
go fmt ./...
go vet ./...
```

## Project Structure (项目结构)

```
frontend/
├── app/
│   ├── (protected)/
│   │   └── dashboard/              # Dashboard 页面
│   │       ├── page.tsx            # 主页面
│   │       ├── loading.tsx         # 加载状态
│   │       ├── error.tsx           # 错误边界
│   │       └── components/         # Dashboard 组件
│   │           ├── WelcomeHeader.tsx
│   │           ├── QuickContinue.tsx
│   │           ├── StatsCards.tsx
│   │           ├── TopicProgress.tsx
│   │           └── RecentQuizzes.tsx
│   └── page.tsx                    # 根页面（重定向到 /dashboard）
├── components/
│   └── providers/
│       └── WebSocketProvider.tsx   # WebSocket 上下文
├── lib/
│   ├── api.ts                      # API 调用
│   ├── websocket.ts                # WebSocket 客户端
│   └── utils/
│       ├── time.ts                 # 时间格式化
│       └── progress.ts             # 进度计算
└── types/
    └── dashboard.ts                # 类型定义

backend/
├── internal/
│   ├── controller/
│   │   └── progress_controller.go  # 新增 GetLastLearning 方法
│   ├── service/
│   │   └── progress_service.go     # 学习天数计算、最后学习记录
│   └── websocket/
│       ├── hub.go                  # WebSocket 连接管理
│       ├── client.go               # WebSocket 客户端
│       └── events.go               # 事件定义
└── api/
    └── v1/
        ├── progress.go             # 新增 /api/v1/progress/last 路由
        └── websocket.go            # WebSocket 路由
```

## Common Tasks (常见任务)

### 添加新的统计卡片

1. 在 `StatsCards.tsx` 中添加新卡片
2. 从 API 获取数据
3. 更新类型定义

```typescript
// types/dashboard.ts
interface DashboardStats {
  studyDays: number
  totalChapters: number
  completedChapters: number
  progressPercentage: number
  weeklyActivity: number
  // 新增字段
  totalQuizzes: number
}
```

### 添加新的 WebSocket 事件

1. 在后端 `events.go` 中定义事件类型
2. 在前端 `WebSocketProvider.tsx` 中处理事件
3. 更新相关组件状态

```typescript
// 前端处理新事件
ws.onmessage = (event) => {
  const message = JSON.parse(event.data)
  
  switch (message.event) {
    case 'progress_updated':
      handleProgressUpdate(message.data)
      break
    case 'quiz_completed':
      handleQuizCompleted(message.data)
      break
    // 新增事件
    case 'achievement_unlocked':
      handleAchievement(message.data)
      break
  }
}
```

### 调试 WebSocket 连接

```bash
# 使用 wscat 工具测试 WebSocket
npm install -g wscat
wscat -c "ws://localhost:8080/api/v1/ws/dashboard?token=YOUR_TOKEN"
```

## Troubleshooting (故障排除)

### 问题 1: Dashboard 页面显示空白

**可能原因**: API 调用失败或数据格式不正确

**解决方案**:
1. 检查浏览器控制台错误
2. 检查后端日志
3. 验证 API 响应格式

```bash
# 测试 API
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/progress/last
```

### 问题 2: WebSocket 连接失败

**可能原因**: Token 无效或后端未启动 WebSocket 服务

**解决方案**:
1. 检查 token 是否有效
2. 确认后端 WebSocket 路由已注册
3. 检查防火墙设置

```bash
# 检查 WebSocket 端点
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Key: test" \
  http://localhost:8080/api/v1/ws/dashboard
```

### 问题 3: 路由冲突（根路径 `/` 无法访问）

**可能原因**: 路由优先级配置错误

**解决方案**:
1. 检查 `app/page.tsx` 是否正确重定向
2. 检查 Next.js 路由配置
3. 清除 `.next` 缓存

```bash
# 清除 Next.js 缓存
cd frontend
rm -rf .next
npm run dev
```

### 问题 4: 数据不实时更新

**可能原因**: WebSocket 事件未正确触发或前端未正确处理

**解决方案**:
1. 检查后端是否在数据变更时触发 WebSocket 事件
2. 检查前端 WebSocket 消息处理逻辑
3. 使用浏览器开发工具的 Network 标签查看 WebSocket 消息

## Testing (测试)

### 运行单元测试

```bash
# 前端
cd frontend
npm test

# 后端
cd backend
go test -v ./internal/controller
go test -v ./internal/service
go test -v ./internal/websocket
```

### 运行集成测试

```bash
# 前端
cd frontend
npm run test:integration

# 后端
cd backend
go test -v -tags=integration ./...
```

### 测试覆盖率

```bash
# 前端
cd frontend
npm run test:coverage

# 后端
cd backend
go test -cover ./...
```

## Deployment (部署)

### 构建生产版本

```bash
# 前端
cd frontend
npm run build

# 后端
cd backend
go build -o bin/server main.go
```

### 运行生产版本

```bash
# 启动后端
cd backend
./bin/server

# 前端静态文件由后端托管
# 访问 http://localhost:8080/
```

## Additional Resources (额外资源)

- [Next.js Documentation](https://nextjs.org/docs)
- [Ant Design Documentation](https://ant.design/docs/react/introduce)
- [GoFrame Documentation](https://goframe.org/pages/viewpage.action?pageId=1114119)
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

## Getting Help (获取帮助)

如果遇到问题，请：

1. 查看本文档的故障排除部分
2. 检查相关文档：[spec.md](./spec.md), [plan.md](./plan.md), [data-model.md](./data-model.md)
3. 查看 API 契约：[contracts/](./contracts/)
4. 提交 Issue 到项目仓库

## Next Steps (下一步)

1. 阅读 [spec.md](./spec.md) 了解功能需求
2. 阅读 [plan.md](./plan.md) 了解技术方案
3. 阅读 [data-model.md](./data-model.md) 了解数据结构
4. 查看 [contracts/](./contracts/) 了解 API 契约
5. 开始实施任务（参考 `tasks.md`，待生成）

---

**Happy Coding!** 🚀
