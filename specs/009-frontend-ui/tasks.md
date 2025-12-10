---
description: "Tasks for Go-Study2 前端 UI"
---

# Tasks: Go-Study2 前端 UI

**Input**: Design documents from `/specs/009-frontend-ui/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml

**Tests**: 前后端需 ≥80% 覆盖，本文各故事均含强制测试任务。
**Organization**: 按用户故事分组，确保每个故事可独立实现与验证。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行（不同文件且无依赖）
- **[Story]**: 用户故事标签（US1/US2/US3）
- 描述内需给出精确文件路径

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 初始化前端工程与基础配置

- [ ] T001 初始化 Next.js 14 App Router 目录骨架于 `frontend/src/{app,components,lib,types,styles,tests}`
- [ ] T002 [P] 在 `frontend/package.json` 声明 AntD 5、Tailwind CSS, SWR, Axios, Prism 依赖并锁定版本
- [ ] T003 [P] 配置 `frontend/tailwind.config.js` 与 `frontend/src/styles/globals.css`，含断点与基础样式
- [ ] T004 [P] 配置 `frontend/src/app/layout.tsx` 注入 AntD `ConfigProvider` 与 Tailwind 样式
- [ ] T005 配置 `frontend/next.config.js` 为 `output: 'export'` 且代理 `/api` 至 `http://localhost:8080`
- [ ] T006 配置前端代码质量工具：`frontend/.eslintrc.js`、`frontend/.prettierrc`、`frontend/tsconfig.json`、`frontend/jest.config.ts`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成全局基础能力，未完成前禁止进入用户故事

- [ ] T007 创建统一 Axios 客户端与错误拦截、401 刷新逻辑于 `frontend/src/lib/http.ts`
- [ ] T008 [P] 实现内存+localStorage token 管理与刷新队列于 `frontend/src/lib/auth/tokenStore.ts`
- [ ] T009 [P] 定义核心类型（User/Topic/Chapter/Progress/Quiz）于 `frontend/src/types/{auth.ts,content.ts,progress.ts,quiz.ts}`
- [ ] T010 [P] 编写静态路由清单与 `generateStaticParams` 支撑文件于 `frontend/src/lib/staticManifest.ts`
- [ ] T011 [P] 建立通用 SWR Hook 与请求封装于 `frontend/src/lib/hooks/useApi.ts`
- [ ] T012 配置全局错误提示与重试入口 Provider 于 `frontend/src/app/providers.tsx`
- [ ] T013 建立测试基线与 MSW mock 于 `frontend/tests/setupTests.ts`、`frontend/tests/mocks/handlers.ts`

---

## Phase 3: 用户故事 1 - 登录后浏览学习主题 (Priority: P1) 🎯 MVP

**Goal**: 支持注册/登录并浏览主题与章节内容，含代码高亮
**Independent Test**: 仅登录+主题/章节浏览即可交付可演示价值

### Tests for User Story 1 (MANDATORY)

- [ ] T014 [P] [US1] 编写认证与主题列表契约测试于 `frontend/tests/contracts/auth-topics.contract.test.ts`
- [ ] T015 [P] [US1] 编写“登录后浏览主题与章节”集成测试于 `frontend/tests/integration/login-browse.test.tsx`

### Implementation for User Story 1

- [ ] T016 [P] [US1] 实现认证 API 封装（register/login/logout/profile/refresh）于 `frontend/src/lib/api/auth.ts`
- [ ] T017 [P] [US1] 实现主题与章节 API 封装于 `frontend/src/lib/api/topics.ts`
- [ ] T018 [US1] 创建注册/登录表单组件（含错误提示与格式校验）于 `frontend/src/components/AuthForms/index.tsx`
- [ ] T019 [US1] 实现受保护路由布局与登录重定向、退出清理逻辑于 `frontend/src/app/(app)/layout.tsx`
- [ ] T020 [US1] 完成登录/注册页面与路由于 `frontend/src/app/(auth)/login/page.tsx` 与 `frontend/src/app/(auth)/register/page.tsx`
- [ ] T021 [US1] 构建主题列表页（含简介、章节数量、占位态）于 `frontend/src/app/(app)/topics/page.tsx`
- [ ] T022 [US1] 构建章节阅读页并集成 Prism 高亮与锚点导航于 `frontend/src/app/(app)/topics/[topic]/[chapter]/page.tsx`
- [ ] T023 [US1] 实施响应式断点规则与视觉回归校验于 `frontend/src/app/(app)/topics/page.tsx` 与 `frontend/src/app/(app)/topics/[topic]/[chapter]/page.tsx`
- [ ] T024 [US1] 校验退出流程清除内存/localStorage token 与跳转登录于 `frontend/src/app/(app)/layout.tsx`

**Checkpoint**: 完成 US1 后可独立演示登录与内容浏览

---

## Phase 4: 用户故事 2 - 进度跟踪与续学 (Priority: P2)

**Goal**: 记录学习进度、显示完成率，并提供“继续上次学习”入口
**Independent Test**: 单独实现进度记录与续学跳转即可交付价值

### Tests for User Story 2 (MANDATORY)

- [ ] T025 [P] [US2] 编写进度读取/写入契约测试于 `frontend/tests/contracts/progress.contract.test.ts`
- [ ] T026 [P] [US2] 编写“记录进度并续学”集成测试于 `frontend/tests/integration/progress-continue.test.tsx`

### Implementation for User Story 2

- [ ] T027 [P] [US2] 实现进度 API 封装于 `frontend/src/lib/api/progress.ts`
- [ ] T028 [US2] 编写进度 SWR Hook 与状态计算（完成率、lastVisit）于 `frontend/src/lib/hooks/useProgress.ts`
- [ ] T029 [US2] 在主题列表组件展示进度百分比与最近访问于 `frontend/src/components/TopicList/index.tsx`
- [ ] T030 [US2] 在章节页上报进度与滚动位置并幂等更新于 `frontend/src/app/(app)/topics/[topic]/[chapter]/page.tsx`
- [ ] T031 [US2] 在主题列表页提供“继续上次学习”入口与跳转逻辑于 `frontend/src/app/(app)/topics/page.tsx`

**Checkpoint**: 完成 US2 后，进度记录与续学路径可独立验证

---

## Phase 5: 用户故事 3 - 主题测验与成绩查看 (Priority: P3)

**Goal**: 支持测验作答、评分结果展示与历史记录查看
**Independent Test**: 独立的测验流程与历史列表即可交付价值

### Tests for User Story 3 (MANDATORY)

- [ ] T032 [P] [US3] 编写测验获取/提交/历史契约测试于 `frontend/tests/contracts/quiz.contract.test.ts`
- [ ] T033 [P] [US3] 编写“作答测验并查看成绩历史（含网络中断与防重复计分）”集成测试于 `frontend/tests/integration/quiz-flow.test.tsx`

### Implementation for User Story 3

- [ ] T034 [P] [US3] 实现测验 API 封装于 `frontend/src/lib/api/quiz.ts`
- [ ] T035 [US3] 实现测验状态管理与评分逻辑 Hook（含幂等提交与重试提示）于 `frontend/src/lib/hooks/useQuiz.ts`
- [ ] T036 [US3] 构建测验作答与防重复提交组件于 `frontend/src/components/QuizRunner/index.tsx`
- [ ] T037 [US3] 构建测验历史页（筛选/排序）于 `frontend/src/app/(app)/quiz/history/page.tsx`
- [ ] T038 [US3] 在主题/章节页面集成测验入口与结果提示于 `frontend/src/app/(app)/topics/[topic]/page.tsx`

**Checkpoint**: 完成 US3 后，测验全流程可独立演示

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: 多故事共用的完善与收尾

- [ ] T039 [P] 完成统一错误与空态体验核查，补充 `frontend/src/components/StateHints/index.tsx`
- [ ] T040 [P] 实现 404/回退页面并在静态导出后验证路由优先级 `/api/*` 于 `frontend/src/app/not-found.tsx` 与 `frontend/src/app/fallback/page.tsx`
- [ ] T041 [P] 校验 CLI/HTTP 兼容与既有路由/响应契约回归于 `backend/tests` 或 `frontend/tests/contracts`
- [ ] T042 [P] 安全检查：HTTPS 配置、token 存储、bcrypt/敏感信息校验清单于 `specs/009-frontend-ui/quickstart.md`
- [ ] T043 运行 `./build.bat`（无则按 quickstart）并验证 `frontend/out` 静态导出与 `npm test`

---

## Dependencies & Execution Order

- Phase 1 → Phase 2 → 各用户故事（US1 → US2 → US3）→ Polish
- 所有用户故事依赖 Phase 2 完成方可开始；故事间按优先级执行，但可在基础完成后并行开发
- 每个故事内：先测试任务，再模型/Hook，再页面与交互

## Parallel Execution Examples

- US1 可并行：契约测试 (T014) + 集成测试 (T015) + API 封装 (T016, T017)
- US2 可并行：契约测试 (T025) + 进度 API 封装 (T027) + 进度 UI (T029)
- US3 可并行：契约测试 (T032) + 测验 API (T034) + 历史页 (T037)
- 跨故事并行：US1 与 US2、US3 可在 Phase 2 完成后由不同成员各自推进

## Implementation Strategy

- MVP：完成 Phase 1-2 后先交付 US1，验证登录+内容浏览闭环
- 增量：在保证 US1 稳定后追加 US2（进度），再追加 US3（测验），每次完成均可独立演示与测试
- 质量：所有故事先写测试确保失败，再实现；完成后运行 `npm test` 与 `./build.bat` 验证静态导出

