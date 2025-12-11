# Tasks: Go-Study2 学习闭环与测验体系

**Input**: Design documents from `/specs/011-learning-progress-quiz/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Per the constitution, features MUST have at least 80% unit test coverage. Test tasks are **MANDATORY**.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `backend/internal/`, `backend/tests/`
- **Frontend**: `frontend/app/`, `frontend/components/`, `frontend/services/`, `frontend/tests/`
- Paths shown below follow plan.md structure

## Constitution Guardrails

- 所有注释与用户文档相关任务必须产出中文内容,且保持清晰一致(Principle V/XV)。
- 需规划达到>=80%测试覆盖,各包包含 *_test.go 与示例,前端核心组件同样达标(Principle III/XXI/XXXVI)。
- 目录/文件/函数保持单一职责与可预测结构,遵循标准 Go 布局(仅根目录 main, go.mod/go.sum 完整)并补齐包 README(Principle IV/VIII/XVIII/XIX)。
- 外部依赖与复杂度最小化,错误处理显式,避免 YAGNI(Principle II/VI/IX)。
- 涉及学习章节/菜单/主题时,需交付 CLI+HTTP 双模式且内容源共享,菜单导航/路由/响应符合约定,结构按章节->子章节->子包组织,Topic 注册一致且显式错误处理(Principle XXII/XXIII/XXIV/XXV)。
- 完成后需包含更新 README 等文档的任务(Principle XI)。

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

 - [X] T001 创建数据库迁移脚本 backend/internal/infra/migrations/011_learning_progress_quiz.sql（LearningProgress/QuizQuestion/QuizSession/QuizAttempt 表与索引）
 - [X] T002 [P] 在 backend/internal/domain/progress/ 创建 progress.go 定义 LearningProgress 实体与状态常量
 - [X] T003 [P] 在 backend/internal/domain/quiz/ 创建 quiz.go 定义 QuizQuestion/QuizSession/QuizAttempt 实体与题型常量
 - [X] T004 [P] 在 backend/internal/domain/progress/ 创建 repository.go 定义 ProgressRepository 接口
 - [X] T005 [P] 在 backend/internal/domain/quiz/ 创建 repository.go 定义 QuizRepository 接口

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T006 在 backend/internal/infra/repository/progress_repo.go 实现 ProgressRepository（CreateOrUpdate/Get/GetByUser/GetByTopic）
- [X] T007 在 backend/internal/infra/repository/quiz_repo.go 实现 QuizRepository（GetQuestionsByChapter/CreateSession/SaveAttempts/GetHistory）
- [X] T008 [P] 在 backend/tests/repository/progress_repo_test.go 编写 ProgressRepository 单元测试（覆盖 CRUD 与唯一约束）
- [X] T009 [P] 在 backend/tests/repository/quiz_repo_test.go 编写 QuizRepository 单元测试（覆盖题目抽取与会话创建）
- [X] T010 [P] 在 backend/internal/domain/progress/ 创建 README.md 说明进度模块职责与接口
- [X] T011 [P] 在 backend/internal/domain/quiz/ 创建 README.md 说明测验模块职责与接口

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - 统一章节学习与恢复 (Priority: P1) 🎯 MVP

**Goal**: 学习者进入任一主题章节，看到统一结构的内容（概述、要点、详细说明、代码示例、陷阱、实践建议），页面自动恢复上次阅读位置并显示阅读进度。

**Independent Test**: 仅上线单个章节的统一结构与阅读恢复即可独立验证学习体验改进。

### Backend Tests for User Story 1 (MANDATORY) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

 - [X] T012 [P] [US1] 在 backend/tests/app/progress_service_test.go 编写进度服务单元测试（CreateOrUpdate/Get/CalculateStatus）
 - [X] T013 [P] [US1] 在 backend/tests/interfaces/progress_handler_test.go 编写进度 API 集成测试（POST/GET /progress）

### Backend Implementation for User Story 1

- [X] T014 [US1] 在 backend/internal/app/progress/service.go 实现 ProgressService（CreateOrUpdateProgress/GetProgress/CalculateChapterStatus/CalculateOverallProgress）
- [X] T015 [US1] 在 backend/internal/app/progress/calculator.go 实现进度计算引擎（主题权重/整体进度算法/状态判断逻辑）
- [X] T016 [US1] 在 backend/internal/interfaces/http/progress_handler.go 实现 POST /api/v1/progress 接口（参数校验/防抖幂等/指数退避响应）
- [X] T017 [US1] 在 backend/internal/interfaces/http/progress_handler.go 实现 GET /api/v1/progress 接口（整体进度与主题汇总）
- [X] T018 [US1] 在 backend/internal/interfaces/http/progress_handler.go 实现 GET /api/v1/progress/{topic} 接口（章节列表与状态）
- [X] T019 [US1] 在 backend/internal/app/progress/service.go 增加卸载前强制同步逻辑（最后更新时间戳防回退）

### Frontend Tests for User Story 1 (MANDATORY) ⚠️

- [X] T020 [P] [US1] 在 frontend/tests/services/progressService.test.ts 编写进度服务单元测试（updateProgress/getProgress/指数退避重试）
- [X] T021 [P] [US1] 在 frontend/tests/components/ChapterProgress.test.tsx 编写章节进度组件测试（进度条/恢复提示）

### Frontend Implementation for User Story 1

- [X] T022 [P] [US1] 在 frontend/services/progressService.ts 实现进度服务（updateProgress/getProgress/getTopicProgress，含 SWR 缓存与 Axios 拦截器）
- [X] T023 [P] [US1] 在 frontend/services/progressService.ts 实现指数退避+抖动重试逻辑（最多 5 次）与卸载前强制同步（beforeunload 事件）
- [X] T024 [P] [US1] 在 frontend/components/progress/ProgressBar.tsx 创建进度条组件（支持分段显示与百分比）
- [X] T025 [P] [US1] 在 frontend/components/progress/ChapterProgress.tsx 创建章节进度指示器（阅读进度/预计剩余时间）
- [X] T026 [US1] 在 frontend/app/topics/[topic]/[chapter]/page.tsx 集成阅读位置恢复（获取 last_position 并自动滚动，显示恢复提示）
- [X] T027 [US1] 在 frontend/app/topics/[topic]/[chapter]/page.tsx 集成滚动监听与防抖上报（每 10 秒累加 read_duration 并覆盖 scroll_progress/last_position）
- [X] T028 [US1] 在 frontend/app/topics/[topic]/[chapter]/page.tsx 添加底部导航栏（上一章/返回列表/下一章/开始测验）

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - 学习进度总览与继续学习 (Priority: P1)

**Goal**: 学习者在进度页查看整体与各主题完成度、章节状态，点击"继续学习"快速跳转到第一个未完成章节。

**Independent Test**: 单独发布进度页即可衡量是否正确显示进度与跳转逻辑。

### Backend Tests for User Story 2 (MANDATORY) ⚠️

- [X] T029 [P] [US2] 在 backend/tests/app/progress_service_test.go 补充整体进度计算测试（权重加权/主题进度/学习天数）
- [X] T030 [P] [US2] 在 backend/tests/interfaces/progress_handler_test.go 补充 GET /progress 响应格式测试

### Backend Implementation for User Story 2

- [X] T031 [US2] 在 backend/internal/app/progress/service.go 补充 GetOverallProgress 方法（计算整体进度/完成章节计数/学习天数/总学习时长）
- [X] T032 [US2] 在 backend/internal/app/progress/service.go 补充 GetNextUnfinishedChapter 方法（按主题优先级与章节序号查找首个未完成章节）

### Frontend Tests for User Story 2 (MANDATORY) ⚠️

- [X] T033 [P] [US2] 在 frontend/tests/components/ProgressOverview.test.tsx 编写进度总览组件测试（整体进度条/主题卡片/筛选排序）
- [X] T034 [P] [US2] 在 frontend/tests/pages/progress.test.tsx 编写进度页面测试（继续学习跳转/状态图标）

### Frontend Implementation for User Story 2

- [X] T035 [P] [US2] 在 frontend/components/progress/ProgressOverview.tsx 创建整体进度卡片（进度条/完成计数/学习天数/总时长）
- [X] T036 [P] [US2] 在 frontend/components/progress/TopicProgressCard.tsx 创建主题进度卡片（权重/进度百分比/章节列表/状态图标/可折叠）
- [X] T037 [P] [US2] 在 frontend/components/progress/ChapterStatusIcon.tsx 创建章节状态图标组件（not_started/in_progress/completed/tested）
- [X] T038 [US2] 在 frontend/app/progress/page.tsx 实现学习进度页面（整体进度/主题卡片/筛选与排序/继续学习按钮）
- [X] T039 [US2] 在 frontend/app/topics/[topic]/page.tsx 添加主题进度条与"继续学习 [章节名]"快捷入口

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - 章节测验与结果反馈 (Priority: P1)

**Goal**: 学习者完成章节后发起测验，按题型答题、提交并查看分数与解析，同时更新章节状态。

**Independent Test**: 仅上线测验页与提交接口即可验证题目抽取、判分、状态更新是否正确。

### Backend Tests for User Story 3 (MANDATORY) ⚠️

- [X] T040 [P] [US3] 在 backend/tests/app/quiz_service_test.go 编写测验服务单元测试（GetQuestions/SubmitQuiz/EvaluateAnswers/判分算法）
- [X] T041 [P] [US3] 在 backend/tests/interfaces/quiz_handler_test.go 编写测验 API 集成测试（GET /quiz/{topic}/{chapter} 与 POST /quiz/submit）

### Backend Implementation for User Story 3

- [X] T042 [US3] 在 backend/internal/app/quiz/service.go 实现 QuizService（GetQuizQuestions/SubmitQuiz/EvaluateAnswers/GetQuizHistory）
- [X] T043 [US3] 在 backend/internal/app/quiz/question_manager.go 实现题目管理器（按难度分层抽取/题型均衡/随机打乱）
- [X] T044 [US3] 在 backend/internal/app/quiz/scoring_engine.go 实现判分引擎（单选/多选/判断/代码输出/改错题判分逻辑）
- [X] T045 [US3] 在 backend/internal/interfaces/http/quiz_handler.go 实现 GET /api/v1/quiz/{topic}/{chapter} 接口（抽题并创建 session）
- [X] T046 [US3] 在 backend/internal/interfaces/http/quiz_handler.go 实现 POST /api/v1/quiz/submit 接口（判分/保存 attempts/更新 session/同步进度状态）
- [X] T047 [US3] 在 backend/internal/interfaces/http/quiz_handler.go 实现 GET /api/v1/quiz/history 接口（测验历史列表/分页/过滤）
- [X] T048 [US3] 在 backend/internal/app/quiz/service.go 实现防重复提交逻辑（session_id 幂等校验）

### Frontend Tests for User Story 3 (MANDATORY) ⚠️

- [X] T049 [P] [US3] 在 frontend/tests/services/quizService.test.ts 编写测验服务单元测试（getQuestions/submitQuiz/getHistory）
- [X] T050 [P] [US3] 在 frontend/tests/components/QuizQuestion.test.tsx 编写题目组件测试（多题型渲染/答案选择）
- [X] T051 [P] [US3] 在 frontend/tests/components/QuizResult.test.tsx 编写结果组件测试（得分/解析/重新测验）

### Frontend Implementation for User Story 3

- [X] T052 [P] [US3] 在 frontend/services/quizService.ts 实现测验服务（getQuestions/submitQuiz/getHistory，含 SWR 缓存）
- [X] T053 [P] [US3] 在 frontend/components/quiz/QuizQuestion.tsx 创建题目组件（支持 5 种题型/代码高亮/选项选择）
- [X] T054 [P] [US3] 在 frontend/components/quiz/QuizNavigation.tsx 创建测验导航组件（题目进度/用时/上一题/下一题/跳过/提交）
- [X] T055 [P] [US3] 在 frontend/components/quiz/QuizResult.tsx 创建结果组件（得分/正确率/通过标识/查看解析/重新测验）
- [X] T056 [P] [US3] 在 frontend/components/quiz/AnswerExplanation.tsx 创建答案解析组件（逐题展示/对错图标/筛选错题）
- [X] T057 [US3] 在 frontend/app/quiz/[topic]/[chapter]/page.tsx 实现章节测验页面（题目展示/导航/本地缓存答案）
- [X] T058 [US3] 在 frontend/app/quiz/[topic]/[chapter]/page.tsx 实现测验提交与结果展示（防重复提交/结果页/解析页）
- [X] T059 [US3] 在 frontend/app/quiz/page.tsx 实现测验历史页面（会话列表/筛选/排序/查看详情）
- [X] T060 [US3] 在 frontend/app/topics/[topic]/[chapter]/page.tsx 添加"开始测验"按钮（跳转到测验页）

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T061 [P] 在 backend/internal/app/progress/ 补充中文注释（service/calculator 职责说明）
- [X] T062 [P] 在 backend/internal/app/quiz/ 补充中文注释（service/question_manager/scoring_engine 职责说明）
- [X] T063 [P] 在 backend/tests/ 补充示例函数（ExampleProgressService/ExampleQuizService）
- [X] T064 [P] 在 frontend/components/ 补充 JSDoc 中文注释（进度/测验组件参数与返回值）
- [ ] T065 [P] 执行 backend 质量检查（go fmt/go vet/golint/go mod tidy）
- [ ] T066 [P] 执行 frontend 质量检查（eslint/prettier/tsc）
- [ ] T067 [P] 运行 backend 测试并验证覆盖率≥80%（go test -cover ./...）
- [ ] T068 [P] 运行 frontend 测试并验证覆盖率≥80%（npm test -- --coverage）
- [X] T069 在根 README.md 更新功能列表（学习进度追踪/章节测验）与路由说明（/progress、/quiz）
- [X] T070 在根 README.md 更新运行说明（build.bat 先行要求/前后端启动命令）
- [ ] T071 运行 quickstart.md 验证流程（build.bat → backend 启动 → frontend 启动 → 测试核心场景）
- [X] T072 [P] 性能验证：进度/测验 API p95 < 300ms（使用 ab/wrk 工具）
- [X] T073 [P] 边界测试：快速进入离开/网络波动/多窗口并发/测验中途刷新（手动或自动化）

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (all P1, so order flexible)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Reuses US1 progress API but independently testable
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - Integrates with US1 progress update but independently testable

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Backend: Repository → Service → Handler
- Frontend: Service → Components → Pages
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Models/components within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "backend/tests/app/progress_service_test.go 编写进度服务单元测试"
Task: "backend/tests/interfaces/progress_handler_test.go 编写进度 API 集成测试"
Task: "frontend/tests/services/progressService.test.ts 编写进度服务单元测试"
Task: "frontend/tests/components/ChapterProgress.test.tsx 编写章节进度组件测试"

# Launch all parallel frontend components for User Story 1 together:
Task: "frontend/services/progressService.ts 实现进度服务"
Task: "frontend/services/progressService.ts 实现指数退避+抖动重试逻辑"
Task: "frontend/components/progress/ProgressBar.tsx 创建进度条组件"
Task: "frontend/components/progress/ChapterProgress.tsx 创建章节进度指示器"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP: 章节学习与恢复)
3. Add User Story 2 → Test independently → Deploy/Demo (进度总览)
4. Add User Story 3 → Test independently → Deploy/Demo (章节测验)
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Backend + Frontend)
   - Developer B: User Story 2 (Backend + Frontend)
   - Developer C: User Story 3 (Backend + Frontend)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- 所有后端注释与文档必须使用中文
- 前端文案与提示保持中文一致性
- 测试覆盖率目标：backend ≥80%，frontend 核心组件 ≥80%

