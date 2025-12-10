# Tasks: Go-Study2 前端UI界面

**Feature Branch**: `009-frontend-ui`  
**Input**: Design documents from `/specs/009-frontend-ui/`  
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/openapi.yaml ✅

**测试要求**: 根据宪章要求，功能必须达到至少 80% 单元测试覆盖率。测试任务为**强制性**。

**组织方式**: 任务按用户故事分组，以支持每个故事的独立实现与测试。

## 格式说明: `- [ ] [TaskID] [P?] [Story?] Description`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 任务所属用户故事（如 US1, US2, US3）
- 描述中包含精确的文件路径

## 路径约定

本项目采用 Web 应用结构：
- **后端**: `backend/internal/`, `backend/tests/`
- **前端**: `frontend/app/`, `frontend/components/`, `frontend/tests/`
- **配置**: `backend/configs/`, `frontend/`
- **数据**: `backend/data/`

## 宪章检查点

- 所有注释与用户文档必须使用中文，保持清晰一致 (Principle V/XV) ✅
- 测试覆盖率 ≥80%，各包含 *_test.go，前端核心组件同样达标 (Principle III/XXI/XXXVI) ✅
- 目录/文件/函数保持单一职责，遵循标准 Go 布局，补齐包 README (Principle IV/VIII/XVIII/XIX) ✅
- 外部依赖最小化，错误处理显式，遵循 YAGNI (Principle II/VI/IX) ✅
- CLI+HTTP 双模式，内容源共享，路由/响应符合约定 (Principle XXII/XXIII/XXIV/XXV) ✅
- 完成后更新 README 等文档 (Principle XI) ✅

---

## Phase 1: 项目设置（共享基础设施）

**目的**: 项目初始化与基础结构搭建

 - [X] T001 在 `backend/` 目录添加 SQLite、JWT、bcrypt 依赖到 go.mod
 - [X] T002 在 `frontend/` 目录创建 Next.js 14 项目并安装 antd、axios、swr、prismjs 依赖
 - [X] T003 [P] 配置 `frontend/next.config.js` 启用静态导出 (output: 'export')
 - [X] T004 [P] 配置 `frontend/tailwind.config.js` 与 `frontend/tsconfig.json` 路径别名
 - [X] T005 [P] 更新 `backend/configs/config.yaml` 添加 database、jwt、static 配置段

---

## Phase 2: 基础层（阻塞性前置条件）

**目的**: 所有用户故事依赖的核心基础设施，必须在任何用户故事开始前完成

**⚠️ 关键**: 此阶段完成前，不能开始任何用户故事工作

### 后端基础设施

- [ ] T006 实现 `backend/internal/infrastructure/database/sqlite.go` 初始化 SQLite 连接与 WAL 模式
- [ ] T007 实现 `backend/internal/infrastructure/database/migrations.go` 数据库迁移（users/learning_progress/quiz_records/refresh_tokens 表）
- [ ] T008 [P] 实现 `backend/internal/pkg/jwt/jwt.go` JWT 生成与验证工具
- [ ] T009 [P] 实现 `backend/internal/pkg/password/password.go` bcrypt 密码哈希与验证
- [ ] T010 [P] 实现 `backend/internal/app/http_server/middleware/auth.go` JWT 认证中间件
- [ ] T011 [P] 实现 `backend/internal/app/http_server/middleware/cors.go` CORS 中间件（开发环境）
- [ ] T012 更新 `backend/main.go` 初始化数据库连接
- [ ] T013 [P] 为 `backend/internal/pkg/jwt/jwt_test.go` 编写单元测试
- [ ] T014 [P] 为 `backend/internal/pkg/password/password_test.go` 编写单元测试
- [ ] T015 [P] 为 `backend/internal/app/http_server/middleware/auth_test.go` 编写单元测试

### 前端基础设施

- [ ] T016 实现 `frontend/lib/api.ts` Axios 实例与请求/响应拦截器
- [ ] T017 实现 `frontend/lib/auth.ts` token 管理工具（内存+localStorage）
- [ ] T018 [P] 实现 `frontend/lib/constants.ts` 常量定义
- [ ] T019 [P] 定义 `frontend/types/api.ts` API 响应类型
- [ ] T020 [P] 定义 `frontend/types/auth.ts` 认证相关类型
- [ ] T021 [P] 定义 `frontend/types/learning.ts` 学习相关类型
- [ ] T022 [P] 定义 `frontend/types/quiz.ts` 测验相关类型
- [ ] T023 实现 `frontend/app/layout.tsx` 根布局与 AntD ConfigProvider
- [ ] T024 [P] 实现 `frontend/components/common/ErrorBoundary.tsx` 全局错误边界
- [ ] T025 [P] 实现 `frontend/components/common/Loading.tsx` 加载状态组件
- [ ] T026 [P] 实现 `frontend/components/common/ErrorMessage.tsx` 错误提示组件
- [ ] T027 配置 `frontend/styles/globals.css` 全局样式与 Tailwind 导入

**检查点**: 基础设施就绪 - 用户故事实现现在可以并行开始

---

## Phase 3: User Story 1 - 登录后浏览学习主题 (Priority: P1) 🎯 MVP

**目标**: 已注册用户登录后，可在浏览器中查看学习主题列表并进入章节阅读，获得代码高亮与分段呈现的内容体验。

**独立测试**: 仅实现登录与内容浏览即可让用户完成基本学习，能独立验证价值。

### 后端 - 认证功能

- [ ] T028 [P] [US1] 实现 `backend/internal/domain/user/entity.go` 用户实体定义
- [ ] T029 [P] [US1] 实现 `backend/internal/domain/user/repository.go` 用户仓储接口
- [ ] T030 [US1] 实现 `backend/internal/domain/user/service.go` 用户服务（注册/登录/登出逻辑）
- [ ] T031 [US1] 实现 `backend/internal/infrastructure/repository/user_repo.go` 用户仓储实现
- [ ] T032 [US1] 实现 `backend/internal/app/http_server/handler/auth.go` 认证 handler（register/login/logout/refresh/profile，支持"记住我"功能）
- [ ] T032a [US1] 在 `backend/internal/app/http_server/handler/auth.go` 中实现根据"记住我"参数设置 refresh token Cookie 过期时间（勾选=7天，未勾选=会话级）
- [ ] T033 [US1] 更新 `backend/internal/app/http_server/router.go` 注册认证路由
- [ ] T034 [P] [US1] 为 `backend/internal/domain/user/service_test.go` 编写单元测试
- [ ] T035 [P] [US1] 为 `backend/internal/infrastructure/repository/user_repo_test.go` 编写单元测试
- [ ] T036 [P] [US1] 为 `backend/internal/app/http_server/handler/auth_test.go` 编写单元测试
- [ ] T037 [US1] 编写 `backend/tests/integration/auth_flow_test.go` 认证流程集成测试
- [ ] T038 [US1] 编写 `backend/tests/contract/auth_api_test.go` 认证 API 契约测试

### 前端 - 认证功能

- [ ] T039 [P] [US1] 实现 `frontend/contexts/AuthContext.tsx` 认证上下文
- [ ] T040 [P] [US1] 实现 `frontend/hooks/useAuth.ts` 认证 Hook
- [ ] T041 [P] [US1] 实现 `frontend/components/auth/LoginForm.tsx` 登录表单组件（包含"记住我"复选框）
- [ ] T041a [US1] 在 `frontend/lib/auth.ts` 中实现"记住我"状态管理（传递给登录 API）
- [ ] T042 [P] [US1] 实现 `frontend/components/auth/RegisterForm.tsx` 注册表单组件
- [ ] T043 [P] [US1] 实现 `frontend/components/auth/AuthGuard.tsx` 认证路由守卫
- [ ] T044 [US1] 实现 `frontend/app/(auth)/login/page.tsx` 登录页
- [ ] T045 [US1] 实现 `frontend/app/(auth)/register/page.tsx` 注册页
- [ ] T046 [US1] 实现 `frontend/app/(protected)/layout.tsx` 受保护路由 Layout（认证检查）
- [ ] T047 [P] [US1] 为 `frontend/tests/components/LoginForm.test.tsx` 编写组件测试
- [ ] T048 [P] [US1] 为 `frontend/tests/components/RegisterForm.test.tsx` 编写组件测试
- [ ] T049 [P] [US1] 为 `frontend/tests/lib/auth.test.ts` 编写单元测试

### 前端 - 布局与导航

- [ ] T050 [P] [US1] 实现 `frontend/components/layout/Header.tsx` 页头组件（导航栏+用户菜单）
- [ ] T051 [P] [US1] 实现 `frontend/components/layout/Footer.tsx` 页脚组件
- [ ] T052 [P] [US1] 实现 `frontend/components/layout/Sidebar.tsx` 侧边栏组件（移动端折叠）

### 前端 - 学习内容展示

- [ ] T053 [P] [US1] 实现 `frontend/components/learning/TopicCard.tsx` 主题卡片组件
- [ ] T054 [P] [US1] 实现 `frontend/components/learning/ChapterList.tsx` 章节列表组件
- [ ] T055 [US1] 实现 `frontend/components/learning/ChapterContent.tsx` 章节内容组件（Markdown 渲染 + Prism.js 代码高亮）
- [ ] T056 [US1] 实现 `frontend/app/(protected)/topics/page.tsx` 主题列表页
- [ ] T057 [US1] 实现 `frontend/app/(protected)/topics/[topic]/page.tsx` 主题详情/章节列表页
- [ ] T058 [US1] 实现 `frontend/app/(protected)/topics/[topic]/[chapter]/page.tsx` 章节内容页
- [ ] T059 [P] [US1] 为 `frontend/tests/components/ChapterContent.test.tsx` 编写组件测试
- [ ] T060 [P] [US1] 为 `frontend/tests/components/TopicCard.test.tsx` 编写组件测试

**检查点**: 此时 User Story 1 应完全功能可用且可独立测试

---

## Phase 4: User Story 2 - 进度跟踪与续学 (Priority: P2)

**目标**: 用户在阅读章节时，系统记录学习进度，并在再次登录时提供从上次位置继续学习的入口。

**独立测试**: 仅实现进度记录与续学入口即可独立验证，且不依赖测验功能。

### 后端 - 学习进度功能

- [ ] T061 [P] [US2] 实现 `backend/internal/domain/progress/entity.go` 学习进度实体定义
- [ ] T062 [P] [US2] 实现 `backend/internal/domain/progress/repository.go` 进度仓储接口
- [ ] T063 [US2] 实现 `backend/internal/domain/progress/service.go` 进度服务（记录/查询逻辑）
- [ ] T064 [US2] 实现 `backend/internal/infrastructure/repository/progress_repo.go` 进度仓储实现
- [ ] T065 [US2] 实现 `backend/internal/app/http_server/handler/progress.go` 进度 handler（获取/记录进度）
- [ ] T066 [US2] 更新 `backend/internal/app/http_server/router.go` 注册进度路由
- [ ] T067 [P] [US2] 为 `backend/internal/domain/progress/service_test.go` 编写单元测试
- [ ] T068 [P] [US2] 为 `backend/internal/infrastructure/repository/progress_repo_test.go` 编写单元测试
- [ ] T069 [P] [US2] 为 `backend/internal/app/http_server/handler/progress_test.go` 编写单元测试
- [ ] T070 [US2] 编写 `backend/tests/integration/progress_test.go` 进度记录集成测试

### 前端 - 学习进度功能

- [ ] T071 [P] [US2] 实现 `frontend/hooks/useProgress.ts` 进度管理 Hook
- [ ] T072 [P] [US2] 实现 `frontend/hooks/useScrollPosition.ts` 滚动位置监听 Hook
- [ ] T073 [P] [US2] 实现 `frontend/components/learning/ProgressBar.tsx` 进度条组件
- [ ] T074 [US2] 在 `frontend/app/(protected)/topics/[topic]/[chapter]/page.tsx` 集成进度记录功能
- [ ] T075 [US2] 在 `frontend/app/(protected)/topics/page.tsx` 添加"继续上次学习"入口
- [ ] T076 [US2] 实现 `frontend/app/(protected)/progress/page.tsx` 学习进度总览页
- [ ] T077 [P] [US2] 为 `frontend/tests/hooks/useProgress.test.ts` 编写单元测试

**检查点**: 此时 User Stories 1 和 2 应都能独立工作

---

## Phase 5: User Story 3 - 主题测验与成绩查看 (Priority: P3)

**目标**: 用户在学习主题后可以参加测验，提交答案后立即获得成绩，并能查看历史测验记录。

**独立测试**: 仅提供测验作答、评分与历史记录展示即可独立完成并验证价值。

### 后端 - 测验功能

- [ ] T078 [P] [US3] 实现 `backend/internal/domain/quiz/entity.go` 测验记录实体定义
- [ ] T079 [P] [US3] 实现 `backend/internal/domain/quiz/repository.go` 测验仓储接口
- [ ] T080 [US3] 实现 `backend/internal/domain/quiz/service.go` 测验服务（获取题目/评分/查询历史）
- [ ] T081 [US3] 实现 `backend/internal/infrastructure/repository/quiz_repo.go` 测验仓储实现
- [ ] T082 [US3] 实现 `backend/internal/app/http_server/handler/quiz.go` 测验 handler（获取题目/提交/历史记录）
- [ ] T083 [US3] 更新 `backend/internal/app/http_server/router.go` 注册测验路由
- [ ] T084 [P] [US3] 为 `backend/internal/domain/quiz/service_test.go` 编写单元测试
- [ ] T085 [P] [US3] 为 `backend/internal/infrastructure/repository/quiz_repo_test.go` 编写单元测试
- [ ] T086 [P] [US3] 为 `backend/internal/app/http_server/handler/quiz_test.go` 编写单元测试
- [ ] T087 [US3] 编写 `backend/tests/integration/quiz_flow_test.go` 测验流程集成测试

### 前端 - 测验功能

- [ ] T088 [P] [US3] 实现 `frontend/hooks/useQuiz.ts` 测验管理 Hook
- [ ] T089 [P] [US3] 实现 `frontend/components/quiz/QuizItem.tsx` 测验题目组件（单选/多选）
- [ ] T090 [P] [US3] 实现 `frontend/components/quiz/QuizResult.tsx` 测验结果组件
- [ ] T091 [P] [US3] 实现 `frontend/components/quiz/QuizHistory.tsx` 测验历史列表组件
- [ ] T092 [US3] 实现 `frontend/app/(protected)/quiz/[topic]/page.tsx` 测验作答页
- [ ] T093 [US3] 实现 `frontend/app/(protected)/quiz/history/page.tsx` 测验历史记录页
- [ ] T094 [P] [US3] 为 `frontend/tests/components/QuizItem.test.tsx` 编写组件测试
- [ ] T095 [P] [US3] 为 `frontend/tests/components/QuizResult.test.tsx` 编写组件测试

**检查点**: 所有用户故事现在应该都能独立功能运行

---

## Phase 6: 部署集成与优化

**目的**: 前后端集成、静态文件服务、性能优化与部署准备

### 后端 - 静态文件服务

- [ ] T096 更新 `backend/internal/app/http_server/server.go` 配置静态文件托管 (frontend/out/)
- [ ] T097 实现 SPA 回退逻辑（非 /api/* 路径返回 index.html）
- [ ] T098 测试 API 与静态资源路由优先级
- [ ] T098a [P] 实现 `frontend/app/not-found.tsx` 404 页面（显示友好错误提示与返回首页链接）

### 前端 - 构建与优化

- [ ] T099 [P] 配置 `frontend/next.config.js` 的 generateStaticParams 预生成路由
- [ ] T100 [P] 实现代码分割（动态导入 ChapterContent 等重组件）
- [ ] T101 [P] 配置 SWR 缓存策略（revalidateOnFocus、dedupingInterval）
- [ ] T102 [P] 优化 Prism.js 按需导入语言包（Go/TypeScript/JavaScript/JSON/bash/markdown）
- [ ] T103 执行前端构建测试 (`npm run build`)
- [ ] T104 验证构建产物 `frontend/out/` 目录结构

### 集成测试

- [ ] T105 编写 `backend/tests/integration/learning_flow_test.go` 学习流程端到端测试
- [ ] T106 编写 `backend/tests/contract/api_contract_test.go` API 契约测试（验证响应格式）
- [ ] T107 [P] 编写 `frontend/tests/integration/auth.test.tsx` 前端认证流程集成测试
- [ ] T108 [P] 编写 `frontend/tests/lib/api.test.ts` API 层集成测试

### 文档与配置

- [ ] T109 [P] 创建 `backend/internal/pkg/jwt/README.md` JWT 工具使用文档
- [ ] T110 [P] 创建 `backend/internal/domain/README.md` 领域层架构说明
- [ ] T111 [P] 创建 `backend/internal/infrastructure/README.md` 基础设施层说明
- [ ] T112 [P] 创建 `frontend/README.md` 前端开发指南（安装/运行/构建）
- [ ] T113 更新根目录 `README.md` 添加新功能章节（用户认证/学习进度跟踪/测验功能）
- [ ] T114 [P] 创建 `docs/API.md` API 文档（基于 openapi.yaml）
- [ ] T115 [P] 创建 `docs/DEPLOYMENT.md` 部署指南（环境变量/构建流程/启动步骤）

---

## Phase 7: 最终检查与交付

**目的**: 影响多个用户故事的改进与最终验证

- [ ] T116 [P] 验证所有代码注释与用户文档为中文
- [ ] T117 [P] 执行后端代码质量检查（go fmt, go vet, golint, go mod tidy）
- [ ] T118 [P] 执行前端代码质量检查（ESLint, Prettier）
- [ ] T119 运行所有后端测试并验证覆盖率 ≥80% (`go test -cover ./...`)
- [ ] T120 运行所有前端测试并验证覆盖率 ≥80% (`npm test -- --coverage`)
- [ ] T121 执行 `quickstart.md` 验证（开发环境启动测试）
- [ ] T122 执行生产构建与部署验证
- [ ] T123 验证响应式布局在 Mobile/Tablet/Desktop 断点正确显示
- [ ] T124 验证所有 Edge Cases（会话过期/网络错误/测验重复提交/移动端布局）
- [ ] T125 执行安全检查（JWT secret 环境变量/密码哈希/HttpOnly Cookie/HTTPS 配置）

---

## 依赖关系与执行顺序

### 阶段依赖

- **Phase 1 (Setup)**: 无依赖 - 可立即开始
- **Phase 2 (Foundational)**: 依赖 Phase 1 完成 - **阻塞所有用户故事**
- **Phase 3-5 (User Stories)**: 所有依赖 Phase 2 完成
  - 用户故事可并行进行（如有人力）
  - 或按优先级顺序（P1 → P2 → P3）
- **Phase 6 (Deployment)**: 依赖所需的用户故事完成
- **Phase 7 (Polish)**: 依赖所有用户故事完成

### 用户故事依赖

- **User Story 1 (P1)**: Phase 2 完成后可开始 - 无其他故事依赖
- **User Story 2 (P2)**: Phase 2 完成后可开始 - 可能集成 US1 但应独立可测
- **User Story 3 (P3)**: Phase 2 完成后可开始 - 可能集成 US1/US2 但应独立可测

### 每个用户故事内部

- 后端：实体 → 仓储接口 → 服务 → 仓储实现 → handler → 路由注册 → 测试
- 前端：类型定义 → Hooks → 组件 → 页面 → 测试
- 测试应在实现前编写并确保失败（TDD 可选）

### 并行机会

- Phase 1 中所有任务可并行
- Phase 2 中标记 [P] 的任务可并行
- Phase 2 完成后，所有用户故事可并行开始（如团队容量允许）
- 每个用户故事内标记 [P] 的任务可并行
- 不同用户故事可由不同团队成员并行工作

---

## 并行示例：User Story 1

```bash
# Phase 2 完成后，同时启动 User Story 1 的所有并行任务：

# 后端实体定义（可并行）:
T028: "实现 backend/internal/domain/user/entity.go"
T029: "实现 backend/internal/domain/user/repository.go"

# 前端组件（可并行）:
T041: "实现 frontend/components/auth/LoginForm.tsx"
T042: "实现 frontend/components/auth/RegisterForm.tsx"
T043: "实现 frontend/components/auth/AuthGuard.tsx"

# 测试（可并行）:
T034: "编写 backend/internal/domain/user/service_test.go"
T035: "编写 backend/internal/infrastructure/repository/user_repo_test.go"
T036: "编写 backend/internal/app/http_server/handler/auth_test.go"
```

---

## 实施策略

### MVP 优先（仅 User Story 1）

1. 完成 Phase 1: Setup
2. 完成 Phase 2: Foundational（关键 - 阻塞所有故事）
3. 完成 Phase 3: User Story 1
4. **停止并验证**: 独立测试 User Story 1
5. 如果就绪则部署/演示

### 增量交付

1. 完成 Setup + Foundational → 基础就绪
2. 添加 User Story 1 → 独立测试 → 部署/演示（MVP！）
3. 添加 User Story 2 → 独立测试 → 部署/演示
4. 添加 User Story 3 → 独立测试 → 部署/演示
5. 每个故事添加价值而不破坏之前的故事

### 并行团队策略

多个开发者情况下：

1. 团队一起完成 Setup + Foundational
2. Foundational 完成后：
   - 开发者 A: User Story 1（认证+内容浏览）
   - 开发者 B: User Story 2（进度跟踪）
   - 开发者 C: User Story 3（测验功能）
3. 故事独立完成并集成

---

## 注意事项

- **[P] 任务** = 不同文件，无依赖，可并行
- **[Story] 标签** 将任务映射到特定用户故事以便追溯
- 每个用户故事应该可独立完成和测试
- 实施前验证测试失败（TDD 方法）
- 每个任务或逻辑组后提交
- 在任何检查点停止以独立验证故事
- 避免：模糊任务、同文件冲突、破坏独立性的跨故事依赖

---

## 任务统计

**总任务数**: 129 个任务

**按用户故事分布**:
- Setup (Phase 1): 5 个任务
- Foundational (Phase 2): 22 个任务
- User Story 1 (Phase 3): 36 个任务（新增"记住我"功能）
- User Story 2 (Phase 4): 17 个任务
- User Story 3 (Phase 5): 18 个任务
- Deployment (Phase 6): 21 个任务（新增 404 页面）
- Polish (Phase 7): 10 个任务

**并行机会**: 约 60% 的任务标记为 [P]，可在各自阶段内并行执行

**独立测试标准**:
- User Story 1: 用户可注册/登录并浏览学习内容
- User Story 2: 用户可记录进度并从上次位置继续
- User Story 3: 用户可参加测验并查看历史记录

**建议 MVP 范围**: Phase 1 + Phase 2 + Phase 3 (User Story 1 only) = 63 个任务

---

## 格式验证 ✅

所有任务遵循严格的检查清单格式：
- ✅ 每个任务以 `- [ ]` 开头（Markdown 复选框）
- ✅ 任务 ID 按执行顺序编号（T001-T125）
- ✅ 可并行任务标记 [P]
- ✅ 用户故事任务标记 [Story]（US1/US2/US3）
- ✅ 描述包含精确文件路径
- ✅ 无模糊任务，每个任务可独立执行
