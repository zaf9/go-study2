# Implementation Plan: Go-Study2 前端 UI

**Branch**: `009-frontend-ui` | **Date**: 2025-12-10 | **Spec**: `D:\studyspace\go-study\go-study2\specs\009-frontend-ui\spec.md`  
**Input**: Feature specification from `D:\studyspace\go-study\go-study2\specs\009-frontend-ui\spec.md`

## Summary

基于 Next.js 14（App Router）+ TypeScript 5 + Ant Design 5 + Tailwind CSS 构建浏览器端 UI，提供注册/登录（含强校验与友好错误）、主题浏览、章节阅读、学习进度、测验记录等体验，并保持现有 CLI/API 的兼容性；前端静态产物使用 Next.js SSG 导出并由 GoFrame 2.9.5 同端口托管，数据落地 SQLite。

## Technical Context

**Language/Version**: 前端 Next.js 14（React 18, TypeScript 5）；后端 Go 1.24.5 + GoFrame v2.9.5  
**Primary Dependencies**: 前端 Ant Design 5、Tailwind CSS、SWR、Axios、Prism.js（按需语言包）；后端 GoFrame ORM (gdb)、JWT、中间件、bcrypt  
**Storage**: SQLite3（文件路径 `backend/data/gostudy.db`，启用 WAL，自动迁移）  
**Testing**: 后端 `go test ./...`；前端 Jest + React Testing Library（SWR 与 Axios mock）；契约测试基于 OpenAPI  
**Target Platform**: GoFrame HTTP Server (8080) 单进程托管 `/api/*` 与静态文件 `/`；现代桌面/移动浏览器  
**Project Type**: 全栈 Web（SSG 前端 + Go 后端）  
**Performance Goals**: 本迭代无强制性能验收指标，重点在功能完整与错误提示一致性；保持 SSG 与按需加载以控制包体。  
**Constraints**: 单端口托管、保持 CLI/HTTP 兼容；仅使用 SQLite，不增引其他存储；JWT 7 天过期，access token 内存+localStorage 恢复，refresh token HttpOnly Cookie + `/api/v1/auth/refresh`；依赖最小化且遵循 AntD/Tailwind 约束  
**Scale/Scope**: 并发 ≤ 50；数据库 < 1GB；单用户记录 < 10MB；核心页面路由 ≤ 7 个

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (代码质量与可维护性):** 是，前后端分层、组件单一职责、SSG 降低运行复杂度。  
- **Principle II (显式错误处理):** 是，Axios 拦截器与 GoFrame 中间件返回统一错误码；需在研究中细化。  
- **Principle III/XXI/XXXVI (全面测试):** 是，后端保持 ≥80% 覆盖，新增模块补齐 *_test.go；前端核心页面/组件 + API 钩子需单测。  
- **Principle IV (单一职责):** 是，Auth/Progress/Quiz/Static 独立模块，组件按功能分区。  
- **Principle V/XV (文档中文):** 是，注释与文档使用中文。  
- **Principle VI (YAGNI):** 是，不引入 SSR/CSR 复杂度，优先 SSG + 客户端数据获取。  
- **Principle VII (安全优先):** 是，采用 HttpOnly refresh + 内存 access token；高亮库按需加载避免 XSS 插桩。  
- **Principle VIII/XVIII (可预测结构):** 是，保持 `backend/` Go 标准布局，前端遵循 Next.js 约定。  
- **Principle IX (依赖纪律):** 是，仅最小新增前端依赖；后端不新增。  
- **Principle X (性能优化):** 是，采用路由级代码分割、按需加载高亮语言；数据库索引与 WAL。  
- **Principle XI (文档同步):** 是，完成后更新根 README 与前后端 README。  
- **Principle XIV (清晰分层注释):** 是，处理器/服务/中间件各自中文注释。  
- **Principle XVI (浅层逻辑):** 是，使用早返回和拆函数，避免深嵌套。  
- **Principle XVII (一致开发者体验):** 是，提供 quickstart，统一脚本（优先 `./build.bat`）。  
- **Principle XIX (包级 README):** 是，前端 `frontend/README.md` 与新增后端子包 README。  
- **Principle XX (代码质量执行):** 是，go fmt/vet、golint、go mod tidy；前端 ESLint/Prettier/TS 校验。  
- **Principle XXII (分层菜单导航):** 是，HTTP/CLI 菜单保持现有约定。  
- **Principle XXIII (双学习模式):** 是，复用现有内容源，HTTP/CLI 同源。  
- **Principle XXIV (层次化章节结构):** 是，新增主题/章节命名遵循 snake_case 与子包组织。  
- **Principle XXV (HTTP/CLI 一致性):** 是，路由 `/api/v1/topic/...` 与菜单同步更新。

## Project Structure

### Documentation (this feature)

```text
specs/009-frontend-ui/
├── plan.md              # 本文件 (/speckit.plan 输出)
├── research.md          # Phase 0 输出 (/speckit.plan)
├── data-model.md        # Phase 1 输出 (/speckit.plan)
├── quickstart.md        # Phase 1 输出 (/speckit.plan)
├── contracts/           # Phase 1 输出 (/speckit.plan)
└── tasks.md             # Phase 2 (/speckit.tasks 生成)
```

### Source Code (repository root)

```text
./
├── backend/                     # 现有 GoFrame 服务
│   ├── configs/                 # 配置与证书
│   ├── internal/                # 业务/路由/中间件/章节内容
│   ├── src/learning/            # 学习内容包
│   ├── tests/                   # contract/integration/unit
│   ├── go.mod / go.sum / main.go
│   └── README.md
├── frontend/                    # 新增 Next.js App Router 前端
│   ├── src/
│   │   ├── app/                 # 页面与路由
│   │   ├── components/          # UI 组件 (auth/layout/learning/quiz)
│   │   ├── lib/                 # api/auth/utils
│   │   ├── types/               # TS 类型
│   │   └── styles/              # 全局样式/Tailwind
│   ├── public/                  # 静态资源
│   ├── tests/                   # 前端单测
│   ├── next.config.js / package.json / tsconfig.json / tailwind.config.js
│   └── README.md
└── specs/                       # 特性文档
    └── 009-frontend-ui/...
```

**Structure Decision**: 采用“backend + frontend”双目录，前端使用 Next.js App Router 静态导出，后端保持 GoFrame 标准布局托管 `/api/*` 与静态文件 `/`。

## Complexity Tracking

无额外复杂度引入，当前无需豁免宪章要求。

## 宪章复核（Phase 1 后）
- 已解决澄清项：代码高亮库选用 Prism；认证刷新策略确定为 HttpOnly refresh + `/api/v1/auth/refresh`。  
- 现有方案符合宪章条款，未引入额外依赖或复杂度。
