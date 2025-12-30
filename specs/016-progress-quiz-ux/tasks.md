# Tasks: 学习进度与测验体验优化

**Feature Branch**: `016-progress-quiz-ux`  
**Input**: Design documents from `/specs/016-progress-quiz-ux/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: 遵循Constitution原则,所有功能必须达到≥80%单元测试覆盖率。本功能采用TDD方法,先编写测试,再实现功能。

**Organization**: 任务按用户故事(User Story)分组,每个故事可独立实现和测试,支持增量交付。

## Format: `- [ ] [ID] [P?] [Story?] Description`

- `- [ ]`: Markdown任务复选框(必需)
- **[ID]**: 任务序号(T001, T002...)
- **[P]**: 可并行执行标记(不同文件,无依赖)
- **[Story]**: 用户故事标签(US1, US2, US3...),仅用于用户故事阶段
- 任务描述必须包含精确的文件路径

## Project Structure Context

- **Backend**: `backend/internal/` (DDD分层: domain/, interfaces/, infra/)
- **Frontend**: `frontend/` (Next.js App Router: app/, components/, hooks/, services/)
- **Tests**: `backend/tests/`, `frontend/__tests__/`
- **Database**: SQLite (tables: learning_progress, quiz_sessions, quiz_attempts)

## Constitution Compliance

- ✅ 所有代码注释和用户文档必须使用中文(Principle V/XV)
- ✅ 测试覆盖率≥80%,包含单元测试、集成测试和契约测试(Principle III/XXI/XXXVI)
- ✅ 单一职责,可预测的目录结构,遵循Go标准布局和Next.js规范(Principle IV/VIII/XVIII)
- ✅ 显式错误处理,无静默失败(Principle II)
- ✅ 最小化依赖,避免过度设计(Principle VI/IX)
- ✅ 完成后更新README和文档(Principle XI)

---

## Implementation Strategy

**MVP优先**: User Story 1(准确章节数量)是最小可交付功能,其他故事按优先级递增实现。

**独立测试**: 每个用户故事都包含独立的验收标准和测试用例,可单独演示和验证。

**并行机会**: 标记[P]的任务可并行执行,提高开发效率。

**TDD流程**: 先编写失败的测试,再实现功能,最后重构优化。

---

## Phase 1: Setup (项目初始化)

**目的**: 准备开发环境和基础设施

- [X] T001 检查开发环境依赖(Go 1.24+, Node.js 18+, SQLite3)
- [X] T002 创建feature分支 `016-progress-quiz-ux`
- [X] T003 [P] 安装后端测试依赖(testify, mockery等)到backend/go.mod
- [X] T004 [P] 安装前端测试依赖(@testing-library/react, jest等)到frontend/package.json
- [X] T005 配置后端测试数据库(SQLite in-memory)在backend/tests/config_test.go

---

## Phase 2: Foundational (基础设施 - 必须完成才能开始用户故事)

**目的**: 所有用户故事共享的基础设施,必须先完成

**⚠️ 关键检查点**: 此阶段完成后才能开始用户故事实现

### 数据库迁移

- [X] T006 创建数据库迁移文件backend/internal/infra/migrations/016_add_submitted_at.sql
- [X] T007 在迁移文件中添加ALTER TABLE quiz_sessions ADD COLUMN submitted_at DATETIME
- [X] T008 在迁移文件中添加UPDATE语句回填已有session的submitted_at值
- [X] T009 执行数据库迁移并验证schema变更成功

### 后端核心工具函数

- [X] T010 [P] 在backend/internal/domain/progress/progress.go中定义状态常量(StatusNotStarted, StatusInProgress, StatusCompleted)
- [X] T011 [P] 在backend/internal/domain/progress/progress.go中实现IsCompleted()方法
- [X] T012 [P] 在backend/internal/domain/progress/progress.go中实现MarkAsInProgress()方法
- [X] T013 [P] 在backend/internal/domain/progress/progress.go中实现Complete()方法

### 前端共享工具

- [X] T014 [P] 创建frontend/lib/chapter-count.ts,实现getChapterCount()函数从static-routes.ts获取章节数
- [X] T015 [P] 创建frontend/lib/chapter-status.ts,实现getChapterStatus()函数计算章节状态
- [X] T016 [P] 验证frontend/types/learning.ts中的进度类型定义(ChapterProgress, TopicProgress, ProgressOverview已存在)
- [X] T017 [P] 验证frontend/types/quiz.ts中的测验类型定义(QuizSession, QuizQuestion, QuizResult已存在)

### 基础测试工具

- [X] T018 [P] 创建backend/tests/testutil/db.go,实现SetupTestDB()函数用于测试数据库初始化
- [X] T019 [P] 创建backend/tests/testutil/fixtures.go,实现创建测试数据的辅助函数
- [X] T020 [P] 创建frontend/__tests__/test-utils.tsx,配置React Testing Library的自定义render函数

**检查点**: 基础设施就绪,可开始用户故事并行开发

---

## Phase 3: User Story 1 - 准确显示章节数量 (Priority: P1) 🎯 MVP

**目标**: 修复主题列表页面章节数量不准确问题,确保显示真实的章节总数

**独立测试**: 访问主题列表页面,验证每个主题显示的章节数与frontend/lib/static-routes.ts中定义的数量完全一致

**价值**: 这是用户对系统信任度的基础,错误的数据直接误导学习规划,必须优先修复

### 测试 - User Story 1 (TDD: 先编写测试)

- [X] T021 [P] [US1] 创建frontend/__tests__/lib/chapter-count.test.ts,测试getChapterCount()返回正确的章节数
- [X] T022 [P] [US1] 创建frontend/__tests__/components/TopicCard.test.tsx,测试TopicCard组件显示准确的章节数量
- [X] T023 [P] [US1] 创建frontend/__tests__/app/topics/page.test.tsx,测试主题列表页面渲染所有主题及章节数

### 实现 - User Story 1

- [X] T024 [US1] 修改frontend/components/learning/TopicCard.tsx,使用getChapterCount()获取章节总数
- [X] T025 [US1] 修改frontend/app/(protected)/topics/page.tsx,显示每个主题的准确章节数
- [X] T026 [P] [US1] 在TopicCard中添加零章节的友好提示(“暂无章节”而非“0章节”)
- [X] T027 [US1] 验证主题列表页面的章节数显示,运行npm test确保测试通过

### 集成测试 - User Story 1

- [X] T028 [US1] 创建frontend/__tests__/integration/topic-list.integration.test.tsx,端到端测试主题列表加载和章节数显示

**User Story 1 验收**: 
- ✅ 主题列表页面显示的章节数量100%准确
- ✅ 数据源唯一(static-routes.ts),无API调用
- ✅ 零章节主题显示友好提示

---

## Phase 4: User Story 2 - 完整的章节学习状态 (Priority: P1)

**目标**: 实现未开始/学习中/已完成三种章节状态,清晰展示学习进度

**独立测试**: 访问章节列表,验证每个章节都有明确的状态标识(图标+颜色),状态转换逻辑正确

**价值**: 状态清晰是高效学习的前提,让用户快速识别新内容和待完成任务

### 测试 - User Story 2 (TDD: 先编写测试)

#### 后端测试

- [X] T029 [P] [US2] 创建backend/tests/repository/progress_repo_test.go,测试CreateOrUpdate()方法
- [X] T030 [P] [US2] 在progress_repo_test.go中测试GetByUserAndTopic()返回主题的所有章节记录
- [X] T031 [P] [US2] 在progress_repo_test.go中测试GetOverview()聚合统计功能
- [X] T032 [P] [US2] 创建backend/tests/domain/progress_service_test.go,测试UpdateChapterStatus()创建学习记录
- [X] T033 [P] [US2] 在progress_service_test.go中测试CompleteChapter()更新状态为completed

#### 前端测试

- [X] T034 [P] [US2] 创建frontend/__tests__/lib/chapter-status.test.ts,测试getChapterStatus()状态计算逻辑
- [X] T035 [P] [US2] 创建frontend/__tests__/components/ChapterStatusBadge.test.tsx,测试状态徽章组件渲染
- [X] T036 [P] [US2] 创建frontend/__tests__/components/ChapterList.test.tsx,测试章节列表显示三种状态

### 实现 - User Story 2

#### 后端实现

- [X] T037 [P] [US2] 在backend/internal/infra/repository/progress_repo.go中实现CreateOrUpdate()方法
- [X] T038 [P] [US2] 在backend/internal/infra/repository/progress_repo.go中实现GetByUserAndTopic()方法
- [X] T039 [P] [US2] 在backend/internal/infra/repository/progress_repo.go中实现GetOverview()方法
- [X] T040 [US2] 在backend/internal/domain/progress/service.go中实现UpdateChapterStatus()业务逻辑
- [X] T041 [US2] 在backend/internal/domain/progress/service.go中实现CompleteChapter()业务逻辑
- [X] T042 [US2] 在backend/internal/domain/progress/service.go中实现GetTopicProgressWithStatus()填充所有章节状态
- [X] T043 [US2] 在backend/internal/interfaces/progress_handler.go中添加GET /api/v1/progress/topic/:topic端点
- [X] T044 [US2] 运行go test ./internal/domain/progress/...确保后端测试通过

#### 前端实现

- [X] T045 [P] [US2] 创建frontend/components/learning/ChapterStatusBadge.tsx,实现状态徽章(未开始🔘/学习中📖/已完成✅)
- [X] T046 [US2] 修改frontend/components/learning/ChapterList.tsx,为每个章节显示ChapterStatusBadge
- [X] T047 [US2] 修改frontend/app/(protected)/topics/[topic]/page.tsx,调用/api/v1/progress/topic/:topic获取章节状态
- [X] T048 [P] [US2] 在frontend/services/progressService.ts中添加getTopicProgress()方法
- [X] T049 [US2] 创建frontend/hooks/useChapterStatus.ts,封装章节状态计算逻辑
- [X] T050 [US2] 运行npm test确保前端测试通过

### 集成测试 - User Story 2

- [X] T051 [US2] 创建backend/tests/interfaces/progress_handler_test.go,测试GET /progress/topic/:topic端点完整流程
- [X] T052 [US2] 创建frontend/__tests__/integration/chapter-status.integration.test.tsx,测试章节状态显示和更新

**User Story 2 验收**:
- ✅ 章节列表清晰显示三种状态
- ✅ 首次访问章节自动标记为"学习中"
- ✅ 完成测验后自动更新为"已完成"
- ✅ 状态有明显的视觉区分(颜色+图标)

---

## Phase 5: User Story 3 - 可用的章节测验功能 (Priority: P1)

**目标**: 修复测验加载失败和提交错误,实现完整的测验流程

**独立测试**: 完整走通"打开章节 → 点击测验 → 答题 → 提交 → 查看结果"流程,无错误

**价值**: 测验是学习闭环的关键,当前完全不可用严重影响学习质量

### 测试 - User Story 3 (TDD: 先编写测试)

#### 后端契约测试

- [X] T053 [P] [US3] 创建backend/tests/contract/quiz_api_test.go,测试GET /api/v1/quiz/:topic/:chapter返回正确的题目结构 ✅
- [X] T054 [P] [US3] 在quiz_api_test.go中测试POST /api/v1/quiz/submit返回正确的评分结果 ✅
- [X] T055 [P] [US3] 在quiz_api_test.go中测试重复提交被拒绝(幂等性) ✅ (修复: 添加durationMs字段)

#### 后端单元测试

- [X] T056 [P] [US3] 创建backend/tests/repository/quiz_repo_test.go,测试CreateSession()方法
- [X] T057 [P] [US3] 在quiz_repo_test.go中测试GetActiveSession()查询24小时内的会话
- [X] T058 [P] [US3] 在quiz_repo_test.go中测试MarkAsSubmitted()设置submitted_at字段
- [X] T059 [P] [US3] 创建backend/tests/domain/quiz_service_test.go,测试GetOrCreateSession()业务逻辑
- [X] T060 [P] [US3] 在quiz_service_test.go中测试loadQuestions()从YAML文件加载题目
- [X] T061 [P] [US3] 在quiz_service_test.go中测试SubmitAnswers()计算得分和更新进度

#### 前端单元测试

- [X] T062 [P] [US3] 创建frontend/__tests__/hooks/useQuiz.test.ts,测试useQuiz Hook的加载和提交逻辑
- [X] T063 [P] [US3] 创建frontend/__tests__/components/QuizSession.test.tsx,测试测验组件渲染和交互
- [X] T064 [P] [US3] 在QuizSession.test.tsx中测试防重复提交逻辑

### 实现 - User Story 3

#### 后端实现

- [X] T065 [P] [US3] 在backend/internal/infra/repository/quiz_repo.go中实现CreateSession()方法
- [X] T066 [P] [US3] 在backend/internal/infra/repository/quiz_repo.go中实现GetActiveSession()方法(查询24小时内)
- [X] T067 [P] [US3] 在backend/internal/infra/repository/quiz_repo.go中实现MarkAsSubmitted()方法
- [X] T068 [P] [US3] 在backend/internal/infra/repository/quiz_repo.go中实现SaveAttempts()批量保存答题记录
- [X] T069 [US3] 在backend/internal/domain/quiz/service.go中实现generateSessionID()生成唯一会话ID
- [X] T070 [US3] 在backend/internal/domain/quiz/service.go中实现loadQuestions()从quiz_data目录加载YAML题目
- [X] T071 [US3] 在backend/internal/domain/quiz/service.go中实现GetOrCreateSession()业务逻辑
- [X] T072 [US3] 在backend/internal/domain/quiz/service.go中实现SubmitAnswers()评分和保存记录
- [X] T073 [US3] 在backend/internal/interfaces/quiz_handler.go中实现GET /api/v1/quiz/:topic/:chapter端点
- [X] T074 [US3] 在backend/internal/interfaces/quiz_handler.go中实现POST /api/v1/quiz/submit端点
- [X] T075 [US3] 在quiz_handler.go的submit端点中添加重复提交检查(submitted_at不为空则拒绝)
- [X] T076 [US3] 运行go test ./internal/domain/quiz/...确保后端测试通过

#### 前端实现

- [X] T077 [P] [US3] 在frontend/services/quizService.ts中实现getQuizSession()方法调用GET /quiz/:topic/:chapter
- [X] T078 [P] [US3] 在frontend/services/quizService.ts中实现submitQuiz()方法调用POST /quiz/submit
- [X] T079 [US3] 创建frontend/hooks/useQuiz.ts,封装测验加载、答题选择、提交和结果展示逻辑
- [X] T080 [US3] 在useQuiz.ts中添加防重复提交状态(isSubmitting)
- [X] T081 [P] [US3] 创建frontend/components/quiz/QuizSession.tsx,实现测验答题界面
- [X] T082 [P] [US3] 创建frontend/components/quiz/QuizResult.tsx,实现测验结果展示组件
- [X] T083 [US3] 修改frontend/app/(protected)/topics/[topic]/[chapter]/page.tsx,集成测验入口按钮
- [X] T084 [US3] 在章节页面添加"开始测验"按钮,点击后加载QuizSession组件
- [X] T085 [US3] 在QuizSession中添加空题目提示("该章节暂无测验")
- [X] T086 [US3] 运行npm test确保前端测试通过

### 集成测试 - User Story 3

- [X] T087 [US3] 创建backend/tests/integration/quiz_flow_test.go,测试完整测验流程(创建session → 提交 → 查询结果)
- [X] T088 [US3] 创建frontend/__tests__/integration/quiz-flow.integration.test.tsx,端到端测试用户答题和提交

**User Story 3 验收**:
- ✅ 测验题目正确加载
- ✅ 用户可以选择答案并提交
- ✅ 提交后显示得分、正确率和正确答案
- ✅ 防重复提交机制生效
- ✅ 无题目章节显示友好提示
- ✅ 所有前端单元测试通过(160个测试)
- ✅ 前端构建成功
- ✅ 所有后端单元测试通过(包括3个契约测试)
- ✅ 后端构建成功 (bin/server.exe, ~33MB)

---

## Phase 6: User Story 4 - 统一的进度数据展示 (Priority: P2)

**目标**: 统一所有页面的进度数据源,确保数据一致性

**独立测试**: 在多个页面(主题列表、进度页、章节详情)对比相同指标,数值完全一致

**价值**: 数据一致性增强用户信任,避免混淆

### 测试 - User Story 4 (TDD: 先编写测试)

#### 后端测试

- [ ] T089 [P] [US4] 在backend/tests/interfaces/progress_handler_test.go中测试GET /api/v1/progress/overview返回完整统计
- [ ] T090 [P] [US4] 测试overview接口计算的completionRate = (completedChapters / totalChapters) * 100

#### 前端测试

- [ ] T091 [P] [US4] 创建frontend/__tests__/hooks/useProgress.test.ts,测试useProgress Hook缓存和刷新逻辑
- [ ] T092 [P] [US4] 测试useProgress在不同组件中返回相同的数据(SWR缓存生效)
- [ ] T093 [P] [US4] 创建frontend/__tests__/components/ProgressDashboard.test.tsx,测试进度页面数据展示

### 实现 - User Story 4

#### 后端实现

- [ ] T094 [US4] 在backend/internal/interfaces/progress_handler.go中实现GET /api/v1/progress/overview端点
- [ ] T095 [US4] 在progress_handler.go中调用progressService.GetOverview()获取全局统计
- [ ] T096 [US4] 在overview响应中包含totalChapters(从static-routes.ts同步计算),completedChapters和completionRate
- [ ] T097 [US4] 运行go test ./internal/interfaces/...确保测试通过

#### 前端实现

- [ ] T098 [P] [US4] 在frontend/services/progressService.ts中实现getProgressOverview()方法
- [ ] T099 [US4] 修改frontend/hooks/useProgress.ts,使用SWR统一管理进度数据缓存
- [ ] T100 [US4] 在useProgress中配置SWR的revalidateOnFocus和dedupingInterval确保数据新鲜度
- [ ] T101 [US4] 修改frontend/app/(protected)/progress/page.tsx,使用useProgress Hook获取数据
- [ ] T102 [US4] 修改frontend/components/learning/TopicCard.tsx,同样使用useProgress Hook
- [ ] T103 [US4] 在frontend/contexts/AuthContext.tsx中添加mutateProgress()方法,用于手动刷新进度
- [ ] T104 [US4] 在QuizResult组件提交成功后调用mutateProgress()刷新所有页面进度
- [ ] T105 [US4] 运行npm test确保测试通过

### 集成测试 - User Story 4

- [ ] T106 [US4] 创建frontend/__tests__/integration/progress-consistency.integration.test.tsx,测试跨页面数据一致性

**User Story 4 验收**:
- ✅ 所有页面显示相同的进度统计
- ✅ 完成测验后所有页面立即更新
- ✅ 使用SWR缓存减少重复请求
- ✅ 完成率计算准确

---

## Phase 7: User Story 5 - 便捷的测验导航体验 (Priority: P2)

**目标**: 新增测验中心页面,优化测验访问路径

**独立测试**: 从主页到测验中心≤2步,可查看所有测验记录

**价值**: 提升用户体验,便于集中复习

### 测试 - User Story 5 (TDD: 先编写测试)

#### 后端测试

- [ ] T107 [P] [US5] 创建backend/tests/interfaces/quiz_history_handler_test.go,测试GET /api/v1/quiz/history返回测验记录列表
- [ ] T108 [P] [US5] 测试GET /api/v1/quiz/history/:sessionId返回测验详情

#### 前端测试

- [ ] T109 [P] [US5] 创建frontend/__tests__/components/QuizCenter.test.tsx,测试测验中心组件渲染
- [ ] T110 [P] [US5] 创建frontend/__tests__/components/QuizHistoryCard.test.tsx,测试历史记录卡片显示
- [ ] T111 [P] [US5] 创建frontend/__tests__/app/quiz-center/page.test.tsx,测试测验中心页面

### 实现 - User Story 5

#### 后端实现

- [ ] T112 [P] [US5] 在backend/internal/infra/repository/quiz_repo.go中实现GetUserHistory()方法查询用户所有已提交的测验
- [ ] T113 [P] [US5] 在backend/internal/infra/repository/quiz_repo.go中实现GetSessionDetail()方法查询测验详情
- [ ] T114 [US5] 在backend/internal/domain/quiz/service.go中实现GetQuizHistory()业务逻辑
- [ ] T115 [US5] 在backend/internal/interfaces/quiz_handler.go中实现GET /api/v1/quiz/history端点
- [ ] T116 [US5] 在backend/internal/interfaces/quiz_handler.go中实现GET /api/v1/quiz/history/:sessionId端点
- [ ] T117 [US5] 运行go test确保后端测试通过

#### 前端实现

- [ ] T118 [P] [US5] 在frontend/services/quizService.ts中实现getQuizHistory()方法
- [ ] T119 [P] [US5] 在frontend/services/quizService.ts中实现getQuizDetail()方法
- [ ] T120 [P] [US5] 创建frontend/components/quiz/QuizCenter.tsx,测验中心容器组件
- [ ] T121 [P] [US5] 创建frontend/components/quiz/QuizHistoryCard.tsx,历史记录卡片组件
- [ ] T122 [US5] 创建frontend/app/(protected)/quiz-center/page.tsx,测验中心页面
- [ ] T123 [US5] 在QuizCenter中实现主题筛选功能
- [ ] T124 [US5] 在QuizHistoryCard中显示章节名称、得分、正确率、时间
- [ ] T125 [US5] 修改frontend/components/layout/MainNav.tsx,添加"测验中心"导航链接
- [ ] T126 [US5] 运行npm test确保前端测试通过

### 集成测试 - User Story 5

- [ ] T127 [US5] 创建frontend/__tests__/integration/quiz-center-navigation.integration.test.tsx,测试导航路径

**User Story 5 验收**:
- ✅ 主导航包含"测验中心"入口
- ✅ 从主页到测验中心≤2步
- ✅ 测验记录列表完整显示
- ✅ 可查看测验详情

---

## Phase 8: Polish & Cross-Cutting Concerns (优化与收尾)

**目的**: 错误处理、边缘情况、文档更新

### 错误处理与边缘情况

- [ ] T128 [P] 在backend/internal/interfaces/error_handler.go中统一API错误响应格式
- [ ] T129 [P] 在frontend/lib/error-handler.ts中实现统一的错误提示组件
- [ ] T130 [P] 处理空数据场景: 新用户进度页显示"开始学习"引导在frontend/app/(protected)/progress/page.tsx
- [ ] T131 [P] 处理并发提交: 前端在QuizSession中添加isSubmitting状态禁用按钮
- [ ] T132 [P] 处理测验数据丢失: 在useQuiz中添加localStorage暂存答案,支持"继续上次测验"
- [ ] T133 [P] 处理章节删除: 后端在GetTopicProgressWithStatus中过滤不存在的章节
- [ ] T134 [P] 处理慢网络: 在QuizSession提交时显示Loading状态

### 性能优化

- [ ] T135 [P] 在backend/internal/infra/repository/progress_repo.go中添加数据库索引(user_id, topic)
- [ ] T136 [P] 在frontend中配置SWR的staleTime和cacheTime优化缓存策略
- [ ] T137 [P] 在TopicCard组件中使用React.memo避免不必要的重渲染

### 文档更新

- [ ] T138 在根目录README.md中更新功能清单,标记feature 016为"已完成"
- [ ] T139 在docs/API.md中添加新增的进度和测验API文档
- [ ] T140 更新frontend/components/quiz/README.md,说明测验组件使用方法
- [ ] T141 在specs/016-progress-quiz-ux/quickstart.md中添加验收测试结果

### 代码质量检查

- [ ] T142 运行后端代码格式化: cd backend && go fmt ./...
- [ ] T143 运行后端代码检查: cd backend && go vet ./...
- [ ] T144 运行前端代码格式化: cd frontend && npm run format
- [ ] T145 运行前端ESLint检查: cd frontend && npm run lint

### 最终验证

- [ ] T146 运行所有后端测试并生成覆盖率报告: cd backend && go test ./... -coverprofile=coverage.out
- [ ] T147 验证后端测试覆盖率≥80%: cd backend && go tool cover -func=coverage.out
- [ ] T148 运行所有前端测试并生成覆盖率报告: cd frontend && npm run test:coverage
- [ ] T149 验证前端测试覆盖率≥80%: 检查coverage/lcov-report/index.html
- [ ] T150 执行手动验收测试,完成quickstart.md中的所有场景
- [ ] T151 提交代码并创建Pull Request: git add . && git commit -m "feat: 学习进度与测验体验优化" && git push origin 016-progress-quiz-ux

---

## Dependency Graph (用户故事完成顺序)

```
Setup (Phase 1) → Foundational (Phase 2)
                        ↓
        ┌───────────────┼───────────────┐
        ↓               ↓               ↓
    US1 (P1)        US2 (P1)        US3 (P1)  ← MVP核心功能,可并行开发
        ↓               ↓               ↓
        └───────────────┼───────────────┘
                        ↓
                    US4 (P2)  ← 依赖US1-US3的数据
                        ↓
                    US5 (P2)  ← 依赖US3的测验功能
                        ↓
                Polish (Phase 8)
```

**关键路径**: Setup → Foundational → US1/US2/US3(并行) → US4 → US5 → Polish

**MVP范围**: US1 + US2 + US3(章节数量准确 + 状态追踪 + 测验可用)

---

## Parallel Execution Examples (并行执行示例)

### Setup阶段可并行任务

- T003(后端测试依赖) + T004(前端测试依赖) + T005(测试数据库配置)

### Foundational阶段可并行任务

- T010-T013(后端状态方法) || T014-T017(前端工具函数) || T018-T020(测试工具)

### User Story 1可并行任务

- T021-T023(所有前端测试) → T024-T027(所有实现任务,依赖测试先完成)

### User Story 2可并行任务

**测试阶段**: T029-T033(后端测试) || T034-T036(前端测试)  
**实现阶段**: T037-T039(repository方法) || T045(ChapterStatusBadge组件)

### User Story 3可并行任务

**测试阶段**: T053-T061(后端测试) || T062-T064(前端测试)  
**实现阶段**: T065-T068(repository方法) || T077-T078(前端service方法) || T081-T082(前端组件)

### User Story 4可并行任务

**测试阶段**: T089-T090(后端) || T091-T093(前端)  
**实现阶段**: T098(service方法) || T101-T102(页面修改)

### User Story 5可并行任务

**测试阶段**: T107-T108(后端) || T109-T111(前端)  
**实现阶段**: T112-T113(repository) || T118-T119(service) || T120-T121(组件)

### Polish阶段可并行任务

- T128-T137(所有错误处理和性能优化任务)

---

## Summary

**总任务数**: 151个任务

**任务分布**:
- Setup: 5个任务
- Foundational: 15个任务(必须完成才能开始用户故事)
- User Story 1 (P1): 8个任务(MVP核心)
- User Story 2 (P1): 24个任务(MVP核心)
- User Story 3 (P1): 36个任务(MVP核心)
- User Story 4 (P2): 18个任务
- User Story 5 (P2): 20个任务
- Polish: 24个任务

**并行机会**: 约60%的任务标记[P],可显著缩短开发周期

**测试覆盖**: 每个用户故事包含单元测试、集成测试和契约测试,确保≥80%覆盖率

**MVP建议**: 完成Setup + Foundational + US1 + US2 + US3即可交付核心功能(约68个任务)

**增量交付**: 每个用户故事都是独立可测试的功能增量,支持逐步发布

---

## Format Validation ✅

所有任务严格遵循格式要求:
- ✅ 每个任务以 `- [ ]` 开头(Markdown复选框)
- ✅ 包含任务ID(T001-T151)
- ✅ 可并行任务标记[P]
- ✅ 用户故事任务标记[Story]
- ✅ 包含精确的文件路径
- ✅ 按用户故事分组,支持独立实现和测试
