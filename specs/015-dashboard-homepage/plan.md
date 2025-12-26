# Implementation Plan: Dashboard 首页功能

**Branch**: `015-dashboard-homepage` | **Date**: 2025-12-26 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/015-dashboard-homepage/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

创建真正的 Dashboard 首页，作为用户登录后的默认落地页（**路径为根路径 `/`**），提供学习状态快速概览、一键继续学习、主题进度可视化、最近测验记录展示等功能。采用纯前端实现方案，完全复用现有后端 API，无需后端开发（除新增 `/api/progress/last` 接口和 WebSocket 支持）。使用 Next.js 14 + Ant Design 5 + TypeScript 构建响应式 Dashboard 页面，通过 WebSocket 实现实时数据推送。

## Technical Context

**Language/Version**: 
- Frontend: TypeScript 5.x + Next.js 14.2.15 + React 18
- Backend: Go 1.24.5 + GoFrame v2.9.5 (仅用于新增 API 和 WebSocket)

**Primary Dependencies**:
- Frontend: Next.js 14.2.15, Ant Design 5.x, React 18, TypeScript 5.x
- Backend: GoFrame v2.9.5, **WebSocket support (需在实施前确认：优先使用 GoFrame 内置支持，若无则使用 gorilla/websocket)**

**Storage**: 
- 复用现有数据库（学习进度、测验记录、用户数据）
- 无需新增数据表（学习天数和最后学习记录可从现有数据计算）

**Testing**:
- Frontend: Jest + React Testing Library (单元测试)
- Backend: Go testing package (新增 API 的单元测试)
- 目标覆盖率: ≥80%

**Target Platform**: 
- Web 应用（桌面、平板、手机响应式支持）
- 浏览器要求: 现代浏览器（Chrome 90+, Firefox 88+, Safari 14+, Edge 90+）

**Project Type**: Web application (frontend + backend)

**Performance Goals**:
- Dashboard 页面初始加载时间 < 2 秒（正常网络条件）
- 首屏渲染时间 (FCP) < 1.5 秒
- 交互响应时间 < 100ms
- WebSocket 消息延迟 < 500ms
- Lighthouse 性能分数 > 90

**Constraints**:
- 必须与现有 Ant Design 设计风格保持一致
- 不应破坏现有页面功能和导航
- 根路径 `/` 的使用不应与现有路由冲突
- WebSocket 连接必须处理网络不稳定情况
- 支持响应式设计（桌面、平板、手机）

**Scale/Scope**:
- 预计用户规模: 100-1000 并发用户
- Dashboard 组件数量: 5 个主要组件
- API 端点: 复用 3 个现有 + 新增 1 个
- WebSocket 连接: 每用户 1 个持久连接
- 开发工作量: 6-9 小时

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (代码质量与可维护性):** ✅ 方案采用组件化设计，每个组件单一职责（WelcomeHeader、QuickContinue、StatsCards、TopicProgress、RecentQuizzes），易于理解和测试
- **Principle II (显式错误处理):** ✅ 所有 API 调用统一错误处理，WebSocket 连接失败显示友好提示，数据加载失败显示重试按钮，无静默失败
- **Principle III/XXI/XXXVI (全面测试):** ✅ 规划前端组件单元测试覆盖率 ≥80%，后端新增 API 测试覆盖率 ≥80%，包含空状态、错误状态、边界情况测试
- **Principle IV (单一职责):** ✅ 每个组件职责明确：WelcomeHeader 仅负责欢迎信息，QuickContinue 仅负责继续学习，StatsCards 仅负责统计展示，职责清晰可拆分
- **Principle V/XV (一致文档与中文要求):** ✅ 前端组件注释使用中文，用户界面文案全部中文，后端新增 API 注释全部中文
- **Principle VI (YAGNI):** ✅ 避免过度设计，初期不实现复杂的数据分析、学习曲线图、个性化推荐等功能，仅实现核心 Dashboard 功能
- **Principle VII (安全优先):** ✅ 复用现有身份验证机制（JWT/Session），所有 API 端点和 WebSocket 连接受保护，前端不存储敏感信息
- **Principle VIII/XVIII (可预测结构):** ✅ 遵循 Next.js App Router 标准结构，Dashboard 页面位于 `app/(protected)/dashboard/page.tsx`，组件位于 `components/` 子目录
- **Principle IX (依赖纪律):** ✅ 无新增外部依赖，完全复用现有技术栈（Next.js, Ant Design, React），WebSocket 使用标准浏览器 API 或 GoFrame 内置支持
- **Principle X (性能优化):** ✅ 使用 Next.js Server Components 进行 SSR，代码分割独立打包 Dashboard 组件，避免不必要的重渲染，WebSocket 异步推送避免轮询
- **Principle XI (文档同步):** ✅ 完成后更新 `README.md` 的功能列表、使用说明、项目结构和路线图
- **Principle XIV (清晰分层注释):** ✅ 每个组件文件包含功能说明注释，API 调用层有统一错误处理注释，数据格式化函数有参数和返回值注释
- **Principle XVI (浅层逻辑):** ✅ 避免深层嵌套，使用早返回和卫语句处理空状态和错误状态，复杂逻辑拆分为独立函数
- **Principle XVII (一致开发者体验):** ✅ 开发流程与现有前端开发一致（`npm run dev`），组件结构遵循项目现有模式，降低学习成本
- **Principle XIX (包级 README):** ⚠️ Dashboard 功能为前端页面，不涉及 Go 包，无需包级 README（前端组件通过 JSDoc 注释说明）
- **Principle XX (代码质量执行):** ✅ 前端使用 ESLint + Prettier + TypeScript 检查，后端新增代码使用 go fmt/go vet/golint
- **Principle XXII (分层菜单导航):** N/A Dashboard 为 Web 页面，不涉及 CLI 交互菜单
- **Principle XXIII (双学习模式):** N/A Dashboard 为系统功能页面，不属于学习章节内容，仅支持 HTTP 访问
- **Principle XXIV (层次化章节结构):** N/A Dashboard 不属于学习章节
- **Principle XXV (HTTP/CLI 一致性):** N/A Dashboard 不涉及学习内容的 CLI/HTTP 双模式

**Constitution Check Result**: ✅ **PASS** - 所有适用原则均符合要求

## Project Structure

### Documentation (this feature)

```text
specs/015-dashboard-homepage/
├── spec.md              # 功能规格说明（已完成）
├── plan.md              # 本文件 - 实施计划
├── research.md          # Phase 0 输出 - 技术研究
├── data-model.md        # Phase 1 输出 - 数据模型
├── quickstart.md        # Phase 1 输出 - 快速开始指南
├── contracts/           # Phase 1 输出 - API 契约
│   ├── api-progress-last.md      # 新增 API 契约
│   └── websocket-events.md       # WebSocket 事件契约
├── checklists/          # 质量检查清单
│   └── requirements.md  # 需求质量检查（已完成）
└── tasks.md             # Phase 2 输出 - 任务分解（待 /speckit.tasks 生成）
```

### Source Code (repository root)

```text
# Frontend Structure
frontend/
├── app/
│   ├── (protected)/
│   │   └── dashboard/              # 新增 Dashboard 页面（目录名为 dashboard，但路由配置为根路径 /）
│   │       ├── page.tsx            # Dashboard 主页面（Server Component）
│   │       ├── loading.tsx         # 加载状态
│   │       ├── error.tsx           # 错误边界
│   │       └── components/         # Dashboard 专用组件
│   │           ├── WelcomeHeader.tsx      # 欢迎区域组件
│   │           ├── QuickContinue.tsx      # 快速继续卡片
│   │           ├── StatsCards.tsx         # 统计卡片组
│   │           ├── TopicProgress.tsx      # 主题进度列表
│   │           └── RecentQuizzes.tsx      # 最近测验记录
│   ├── page.tsx                    # 修改：重定向到 /dashboard
│   └── layout.tsx                  # 可能需要调整（WebSocket Provider）
├── components/
│   ├── layout/
│   │   └── Sidebar.tsx             # 修改：首页按钮指向 /dashboard
│   └── providers/
│       └── WebSocketProvider.tsx   # 新增：WebSocket 上下文
├── lib/
│   ├── api.ts                      # 可能需要扩展：新增 API 调用
│   ├── websocket.ts                # 新增：WebSocket 客户端
│   └── utils/
│       ├── time.ts                 # 新增：时间格式化工具
│       └── progress.ts             # 新增：进度计算工具
├── types/
│   └── dashboard.ts                # 新增：Dashboard 数据类型定义
└── __tests__/
    └── dashboard/                  # 新增：Dashboard 组件测试
        ├── WelcomeHeader.test.tsx
        ├── QuickContinue.test.tsx
        ├── StatsCards.test.tsx
        ├── TopicProgress.test.tsx
        └── RecentQuizzes.test.tsx

# Backend Structure (仅新增部分)
backend/
├── internal/
│   ├── controller/
│   │   └── progress_controller.go  # 修改：新增 GetLastLearning 方法
│   ├── service/
│   │   └── progress_service.go     # 修改：新增学习天数计算、最后学习记录查询
│   └── websocket/
│       ├── hub.go                  # 新增：WebSocket 连接管理
│       ├── client.go               # 新增：WebSocket 客户端连接
│       └── events.go               # 新增：WebSocket 事件定义
├── api/
│   └── v1/
│       ├── progress.go             # 修改：新增 /api/v1/progress/last 路由
│       └── websocket.go            # 新增：WebSocket 路由
└── tests/
    ├── controller/
    │   └── progress_controller_test.go  # 新增测试
    └── websocket/
        └── hub_test.go             # 新增测试
```

**Structure Decision**: 
- 采用 **Web application** 结构（Option 2）
- 前端使用 Next.js App Router 结构，Dashboard 页面位于 `app/(protected)/dashboard/`
- 后端仅新增 WebSocket 支持和一个 API 端点，最小化后端改动
- 组件按功能模块组织，便于独立开发和测试

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

**无需填充** - Constitution Check 全部通过，无违反项需要说明。

## Phase 0: Research & Technology Selection

### Research Topics

基于用户提供的技术方案和规格澄清，以下研究任务已完成或明确：

#### R1: WebSocket 实现方案
- **Decision**: 使用 GoFrame 内置 WebSocket 支持（如有）或 gorilla/websocket 库
- **Rationale**: GoFrame 可能已内置 WebSocket 支持，优先使用；若无，gorilla/websocket 是 Go 生态最成熟的 WebSocket 库
- **Alternatives Considered**: 
  - 轮询方案：性能差，服务器负载高，已拒绝
  - Server-Sent Events (SSE)：单向推送，不支持客户端发送，功能受限

#### R2: 学习天数计算策略
- **Decision**: 计算有学习活动的不同日期数（非连续天数）
- **Rationale**: 更准确反映用户实际学习投入，不会因中断学习产生误导性数字
- **Implementation**: SQL 查询 `SELECT COUNT(DISTINCT DATE(created_at)) FROM learning_progress WHERE user_id = ?`

#### R3: 最后学习记录获取方式
- **Decision**: 新增 `/api/v1/progress/last` 接口
- **Rationale**: 后端计算更可靠，避免客户端数据不一致
- **Alternatives Considered**:
  - localStorage 客户端记录：数据可能不准确，跨设备不同步，已拒绝
  - 从现有进度列表前端计算：增加前端复杂度，性能较差

#### R4: 身份验证集成方案
- **Decision**: 复用现有身份验证机制（JWT 或 Session）
- **Rationale**: 保持系统一致性，避免重复开发，降低安全风险
- **Implementation**: Dashboard API 和 WebSocket 连接使用现有中间件保护

#### R5: 时间格式化策略
- **Decision**: 混合格式（24 小时内相对时间 + 超过 24 小时绝对时间）
- **Rationale**: 符合主流应用习惯（GitHub、Twitter），提供最佳用户体验
- **Implementation**: 前端工具函数 `formatTime(timestamp)` 实现逻辑判断

#### R6: WebSocket 重连策略
- **Decision**: 指数退避（初始 1 秒，最大 30 秒，最多 5 次）
- **Rationale**: 业界标准做法，避免服务器过载，给予足够重连机会
- **Implementation**: 前端 WebSocket 客户端实现重连逻辑，后端无需特殊处理

#### R7: 现有 API 复用确认
- **Assumption**: 以下 API 已存在并可复用
  - `GET /api/v1/progress` - 获取学习进度统计
  - `GET /api/v1/quiz/history` - 获取测验历史
  - `GET /api/v1/topics` - 获取主题列表
- **Action Required**: 实施前需确认这些 API 的实际存在性和数据格式

### Technology Stack Summary

| 层级 | 技术选型 | 版本 | 用途 |
|------|---------|------|------|
| 前端框架 | Next.js | 14.2.15 | 页面渲染、路由、SSR |
| UI 组件库 | Ant Design | 5.x | 统计卡片、进度条、表格等 UI 组件 |
| 编程语言 | TypeScript | 5.x | 类型安全、代码提示 |
| 状态管理 | React Hooks | - | useState, useEffect, useContext |
| WebSocket 客户端 | 浏览器原生 API | - | 实时数据推送 |
| 后端框架 | GoFrame | v2.9.5 | HTTP 服务、WebSocket 支持 |
| WebSocket 服务端 | gorilla/websocket | latest | WebSocket 连接管理（如 GoFrame 无内置支持） |
| 测试框架（前端） | Jest + RTL | latest | 组件单元测试 |
| 测试框架（后端） | Go testing | - | API 单元测试 |

## Phase 1: Design & Architecture

### Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    DashboardPage (SSR)                       │
│  ┌───────────────────────────────────────────────────────┐  │
│  │          WebSocketProvider (Context)                  │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │         WelcomeHeader                           │  │  │
│  │  │  - 欢迎信息 + 用户名                              │  │  │
│  │  │  - 累计学习天数                                  │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │         QuickContinue                           │  │  │
│  │  │  - 最后学习的主题/章节                            │  │  │
│  │  │  - 继续学习按钮 → 跳转                            │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │         StatsCards (Grid 3列)                   │  │  │
│  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐         │  │  │
│  │  │  │总体进度  │ │完成章节  │ │本周活跃  │         │  │  │
│  │  │  │  XX%     │ │  X/Y     │ │  N次     │         │  │  │
│  │  │  └──────────┘ └──────────┘ └──────────┘         │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │         TopicProgress (List)                    │  │  │
│  │  │  - 主题1: ████████░░ 80%                        │  │  │
│  │  │  - 主题2: ████░░░░░░ 40%                        │  │  │
│  │  │  - 主题3: ██████████ 100%                       │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │         RecentQuizzes (Table/List)              │  │  │
│  │  │  - 测验1: 主题A / 章节B | 85/100 | 2小时前       │  │  │
│  │  │  - 测验2: 主题C / 章节D | 90/100 | 昨天 10:30    │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

WebSocket 数据流:
Server → Client: progress_updated / quiz_completed 事件
Client: 更新对应组件状态（StatsCards, TopicProgress, RecentQuizzes）
```

### Data Flow

```
1. 页面加载流程（SSR）:
   ┌──────────┐
   │ 用户访问  │
   │ /dashboard│
   └─────┬────┘
         │
         ▼
   ┌──────────────────┐
   │ Server Component │
   │ 并行获取数据:     │
   │ - GET /api/v1/progress        │
   │ - GET /api/v1/progress/last   │
   │ - GET /api/v1/quiz/history    │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ 渲染 HTML 返回    │
   │ (含初始数据)      │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ 客户端 Hydration  │
   │ 建立 WebSocket    │
   └──────────────────┘

2. WebSocket 实时更新流程:
   ┌──────────────────┐
   │ 用户完成学习/测验 │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ 后端触发事件:     │
   │ progress_updated  │
   │ 或 quiz_completed │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ WebSocket 推送    │
   │ 到所有该用户连接  │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ 前端接收事件      │
   │ 更新组件状态      │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ UI 自动刷新       │
   │ (无需手动刷新)    │
   └──────────────────┘

3. 错误处理流程:
   ┌──────────────────┐
   │ API 调用失败      │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ 显示错误提示      │
   │ + 重试按钮        │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ WebSocket 断开    │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ 指数退避重连      │
   │ (最多5次)         │
   └─────┬────────────┘
         │
         ▼
   ┌──────────────────┐
   │ 重连失败          │
   │ 显示友好提示      │
   └──────────────────┘
```

### API Integration Plan

详见 `contracts/` 目录中的 API 契约文档（Phase 1 输出）。

关键 API 端点：
1. `GET /api/v1/progress` - 复用现有
2. `GET /api/v1/progress/last` - **新增**
3. `GET /api/v1/quiz/history?limit=5` - 复用现有（添加 limit 参数）
4. `GET /api/v1/topics` - 复用现有
5. `WS /api/v1/ws/dashboard` - **新增** WebSocket 连接

### State Management Strategy

```typescript
// 1. 服务端状态（SSR 初始数据）
interface DashboardServerData {
  user: { username: string }
  stats: StatsData
  lastLearning: LastLearningData | null
  topicProgress: TopicProgressData[]
  recentQuizzes: QuizRecordData[]
}

// 2. 客户端状态（WebSocket 实时更新）
const [stats, setStats] = useState<StatsData>(initialStats)
const [topicProgress, setTopicProgress] = useState<TopicProgressData[]>(initialProgress)
const [recentQuizzes, setRecentQuizzes] = useState<QuizRecordData[]>(initialQuizzes)

// 3. WebSocket 上下文
const WebSocketContext = createContext<WebSocketContextValue>(null)

// 4. 状态更新逻辑
useEffect(() => {
  const ws = new WebSocket('/api/v1/ws/dashboard')
  
  ws.onmessage = (event) => {
    const { event: eventType, data } = JSON.parse(event.data)
    
    if (eventType === 'progress_updated') {
      // 更新 stats 和 topicProgress
      setStats(prev => ({ ...prev, completedChapters: data.completed }))
      setTopicProgress(prev => updateTopicProgress(prev, data))
    }
    
    if (eventType === 'quiz_completed') {
      // 更新 recentQuizzes
      setRecentQuizzes(prev => [data, ...prev.slice(0, 4)])
    }
  }
  
  return () => ws.close()
}, [])
```

### Performance Optimization Strategy

1. **服务端渲染 (SSR)**:
   - 使用 Next.js Server Components 预渲染 Dashboard 数据
   - 减少客户端 JavaScript 执行时间
   - 提升首屏加载速度

2. **代码分割**:
   - Dashboard 页面独立打包
   - 组件按需加载（如果组件较大）
   - 减少初始 bundle 大小

3. **数据缓存**:
   - 利用 Next.js 数据缓存机制
   - 设置合理的 revalidate 时间
   - 减少不必要的 API 调用

4. **WebSocket 优化**:
   - 仅推送必要的数据字段
   - 客户端去重处理（避免重复更新）
   - 连接池管理（后端）

5. **渲染优化**:
   - 使用 React.memo 避免不必要的重渲染
   - 合理使用 useMemo 和 useCallback
   - 虚拟滚动（如果列表很长）

## Phase 2: Implementation Roadmap

详细任务分解将由 `/speckit.tasks` 命令生成到 `tasks.md` 文件。

### High-Level Phases

#### Phase 1: 基础搭建 (1-2 小时)
- 创建 `/dashboard` 页面和目录结构
- 实现基本布局框架
- 配置路由调整（根路径重定向、侧边栏链接）

#### Phase 2: 组件实现 (2-3 小时)
- 实现 WelcomeHeader 组件
- 实现 StatsCards 组件
- 实现 QuickContinue 组件
- 实现 TopicProgress 组件
- 实现 RecentQuizzes 组件

#### Phase 3: 数据集成 (1-2 小时)
- API 接口调用和数据处理
- 错误处理和空状态
- 数据格式化工具函数
- WebSocket 客户端实现

#### Phase 4: 后端开发 (1-2 小时)
- 实现 `/api/v1/progress/last` 接口
- 实现 WebSocket 服务端
- 实现事件推送逻辑
- 单元测试

#### Phase 5: 样式优化 (1 小时)
- 响应式布局调整
- 细节样式优化
- 动画和交互效果

#### Phase 6: 测试和调试 (1 小时)
- 单元测试编写
- 手动测试各场景
- 修复 bug 和优化

**总计**: 6-9 小时

## Risk Assessment

### Technical Risks

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|---------|
| 现有 API 数据格式不符合预期 | 高 | 中 | 实施前确认 API 契约，必要时调整前端数据处理逻辑 |
| WebSocket 连接不稳定 | 中 | 中 | 实现指数退避重连，降级到页面刷新模式 |
| 根路径 `/` 路由冲突 | 高 | 低 | 仔细检查现有路由配置，确保优先级正确 |
| 性能不达标（加载时间 > 2秒） | 中 | 低 | 使用 SSR、代码分割、数据缓存优化 |
| 跨浏览器兼容性问题 | 低 | 低 | 使用标准 Web API，测试主流浏览器 |

### Implementation Risks

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|---------|
| 开发时间超出预期 | 中 | 中 | 按优先级实施（P1 > P2 > P3），P3 功能可延后 |
| 测试覆盖率不足 | 中 | 低 | 提前规划测试用例，边开发边测试 |
| 设计与现有风格不一致 | 低 | 低 | 严格遵循 Ant Design 规范，复用现有组件 |
| WebSocket 后端实现复杂度高 | 中 | 中 | 使用成熟库（gorilla/websocket），参考最佳实践 |

## Pre-Implementation Checklist

以下问题需要在实施前确认（**责任人：后端开发负责人，截止时间：Phase 1 开始前**）：

1. **现有 API 确认**:
   - [ ] ⏳ `/api/v1/progress` 接口是否存在？返回数据格式是什么？
   - [ ] ⏳ `/api/v1/quiz/history` 接口是否存在？是否支持 `limit` 参数？
   - [ ] ⏳ `/api/v1/topics` 接口是否存在？返回数据格式是什么？

2. **数据库字段确认**:
   - [ ] ⏳ 学习进度表是否有 `created_at` 字段用于计算学习天数？
   - [ ] ⏳ 学习进度表是否有 `last_visited_at` 字段用于获取最后学习记录？

3. **WebSocket 实现确认**:
   - [ ] ⏳ GoFrame v2.9.5 是否内置 WebSocket 支持？如果有，API 是什么？
   - [ ] ⏳ 如 GoFrame 无内置支持，是否批准使用 gorilla/websocket 库？

4. **身份验证确认**:
   - [ ] ⏳ 现有身份验证机制是 JWT 还是 Session？
   - [ ] ⏳ 身份验证中间件如何使用？

5. **部署确认**:
   - [ ] ⏳ 前端静态文件是否由后端托管？
   - [ ] ⏳ WebSocket 连接是否需要特殊的代理配置（如 Nginx）？

**符号说明**:
- ✅ 已确认
- ⏳ 待确认
- ❌ 已确认不可用/需要调整

## Next Steps

1. **执行 Phase 0 研究** - 生成 `research.md`（已在本文档中完成）
2. **执行 Phase 1 设计** - 生成以下文档：
   - `data-model.md` - 数据模型定义
   - `contracts/api-progress-last.md` - 新增 API 契约
   - `contracts/websocket-events.md` - WebSocket 事件契约
   - `quickstart.md` - 快速开始指南
3. **更新 Agent 上下文** - 运行 `update-agent-context.ps1`
4. **执行 `/speckit.tasks`** - 生成详细任务分解到 `tasks.md`
5. **开始实施** - 按照任务优先级逐步实现功能

---

**Plan Status**: ✅ Phase 0 Complete | ⏳ Phase 1 Pending | ⏳ Phase 2 Pending
