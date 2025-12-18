# Tasks: 章节测验体验升级与功能深化 (Quiz UX Enhancement)

**Input**: Design documents from `/specs/014-quiz-ux-enhancement/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/api.yaml ✅, quickstart.md ✅

**Tests**: Per the constitution, features MUST have at least 80% unit test coverage. Test tasks are **MANDATORY**.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `backend/` (Go 1.24 + GoFrame v2.x)
- **Frontend**: `frontend/` (Next.js + React + Ant Design)

## Constitution Guardrails

- 所有注释与用户文档相关任务必须产出中文内容,且保持清晰一致(Principle V/XV)。
- 需规划达到>=80%测试覆盖,各包包含 *_test.go 与示例,前端核心组件同样达标(Principle III/XXI/XXXVI)。
- 目录/文件/函数保持单一职责与可预测结构,遵循标准 Go 布局(仅根目录 main, go.mod/go.sum 完整)并补齐包 README(Principle IV/VIII/XVIII/XIX)。
- 外部依赖与复杂度最小化,错误处理显式,避免 YAGNI(Principle II/VI/IX)。
- 涉及学习章节/菜单/主题时,CLI+HTTP 双模式(历史回顾模式为 Web-only 例外),内容源共享(Principle XXII/XXIII/XXIV/XXV)。
- 完成后需包含更新 README 等文档的任务(Principle XI)。

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目初始化和数据库结构准备

- [x] T001 创建 QuizSession 数据库表迁移文件 in `backend/internal/infra/migrations/`
- [x] T002 创建 QuizAttempt 数据库表迁移文件 in `backend/internal/infra/migrations/`
- [x] T003 [P] 执行数据库迁移并验证表结构

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 所有用户故事都依赖的核心基础设施

**⚠️ CRITICAL**: 必须完成本阶段后才能开始任何用户故事的实现

- [x] T004 创建 QuizSession 实体模型 in `backend/internal/model/entity/quiz_session.go`
- [x] T005 [P] 创建 QuizAttempt 实体模型 in `backend/internal/model/entity/quiz_attempt.go`
- [x] T006 创建 Quiz 仓储接口定义 in `backend/internal/infra/repository/quiz_repo.go`
- [x] T007 实现 Quiz 仓储层（含事务保存会话+答题记录）in `backend/internal/infra/repository/quiz_repo_impl.go`
- [x] T008 [P] 创建 QuizSkeletonLoader 骨架屏组件 in `frontend/components/quiz/QuizSkeletonLoader.tsx`

**Checkpoint**: 基础设施就绪 - 可以开始用户故事的并行实现

---

## Phase 3: User Story 1 - 测验交互优化：防误触与有序标签 (Priority: P1) 🎯 MVP

**Goal**: 实现稳定的 A-D 选项标签和提交二次确认机制，减少用户误操作

**Independent Test**: 前端组件可独立测试标签渲染和确认弹窗拦截逻辑，无需后端依赖

### Tests for User Story 1 (MANDATORY) ⚠️

- [x] T009 [P] [US1] 创建 QuizQuestionCard 组件单元测试 in `frontend/components/quiz/__tests__/QuizQuestionCard.test.tsx`
- [x] T010 [P] [US1] 创建 SubmitConfirmModal 组件单元测试 in `frontend/components/quiz/__tests__/SubmitConfirmModal.test.tsx`

### Implementation for User Story 1

- [x] T011 [P] [US1] 创建 QuizQuestionCard 组件（含有序标签 A-D 渲染）in `frontend/components/quiz/QuizQuestionCard.tsx`
- [x] T012 [P] [US1] 创建 SubmitConfirmModal 组件（显示已答/未答统计）in `frontend/components/quiz/SubmitConfirmModal.tsx`
- [x] T013 [US1] 更新 QuizViewer 组件集成 QuizQuestionCard 和 SubmitConfirmModal in `frontend/app/(protected)/quiz/[topic]/QuizPageClient.tsx`
- [x] T014 [US1] 添加前端提交拦截逻辑和答题状态管理 in `frontend/app/(protected)/quiz/[topic]/QuizPageClient.tsx`

**Checkpoint**: 用户故事 1 完成 - 可独立测试防误触和标签稳定性

---

## Phase 4: User Story 2 - 结果页反馈增强：百分制与题型标识 (Priority: P1)

**Goal**: 实现百分制得分展示、及格状态颜色区分、题型标签显示

**Independent Test**: 可通过 Mock 提交结果数据测试结果页渲染逻辑

### Tests for User Story 2 (MANDATORY) ⚠️

- [x] T015 [P] [US2] 创建 QuizResultPage 组件单元测试 in `frontend/components/quiz/__tests__/QuizResultPage.test.tsx`
- [x] T016 [P] [US2] 创建后端 Submit API 单元测试 in `backend/internal/app/http_server/handler/quiz_submit_test.go`

### Implementation for User Story 2

- [x] T017 [US2] 更新后端 Submit API 返回百分制得分和通过状态 in `backend/internal/app/http_server/handler/quiz.go`
- [x] T018 [US2] 更新后端 Quiz Service 计算百分制得分逻辑 in `backend/internal/app/quiz/scoring_engine.go`
- [x] T019 [P] [US2] 创建 QuizResultPage 组件（百分制得分、颜色区分）in `frontend/components/quiz/QuizResultPage.tsx`
- [x] T020 [P] [US2] 创建 QuestionTypeTag 组件（单选/多选/改错标签）in `frontend/components/quiz/QuestionTypeTag.tsx`
- [x] T021 [US2] 集成题型标签到 QuizQuestionCard 组件 in `frontend/components/quiz/QuizQuestionCard.tsx`
- [x] T022 [US2] 添加多选题操作引导说明 in `frontend/src/components/quiz/QuizQuestionCard.tsx`

**Checkpoint**: 用户故事 2 完成 - 可独立测试结果页展示和题型标识

---

## Phase 5: User Story 3 - 历史回顾模式 (Priority: P2)

**Goal**: 实现测验历史列表和详情回顾功能，支持查看错题、正确答案及解析

**Independent Test**: 前端可通过 Mock 历史数据列表测试详情页渲染

### Tests for User Story 3 (MANDATORY) ⚠️

- [x] T023 [P] [US3] 创建 GET /quiz/history API 契约测试 in `backend/tests/contract/quiz/quiz_history_test.go`
- [x] T024 [P] [US3] 创建 GET /quiz/history/{sessionId} API 契约测试 in `backend/tests/contract/quiz/quiz_review_test.go`
- [x] T025 [P] [US3] 创建 QuizHistoryPage 组件单元测试 in `frontend/components/quiz/__tests__/QuizHistoryPage.test.tsx`
- [x] T026 [P] [US3] 创建 QuizReviewPage 组件单元测试 in `frontend/components/quiz/__tests__/QuizReviewPage.test.tsx`

### Implementation for User Story 3

- [x] T027 [US3] 实现 GET /quiz/history API 端点 in `backend/internal/app/http_server/handler/quiz.go`
- [x] T028 [US3] 实现 GET /quiz/history/{sessionId} API 端点 in `backend/internal/app/http_server/handler/quiz.go`
- [x] T029 [US3] 添加历史查询和详情查询服务方法 in `backend/internal/app/quiz/service.go`
- [x] T030 [US3] 注册新 API 路由 in `backend/internal/app/http_server/router.go`
- [x] T031 [P] [US3] 创建 QuizHistoryPage 页面组件 in `frontend/app/(protected)/quiz/history/page.tsx`
- [x] T032 [P] [US3] 创建 QuizReviewPage 页面组件 in `frontend/app/(protected)/quiz/history/[sessionId]/page.tsx`
- [x] T033 [US3] 更新 QuizViewer 支持 review 模式（禁用选择、显示解析）in `frontend/app/(protected)/quiz/history/[sessionId]/page.tsx`
- [x] T034 [US3] 创建 AnswerIndicator 组件（显示用户答案 vs 正确答案）in `frontend/components/quiz/AnswerIndicator.tsx`

**Checkpoint**: 用户故事 3 完成 - 可独立测试历史列表和回顾功能

---

## Phase 6: User Story 4 - 全链路入口与元数据展示 (Priority: P3)

**Goal**: 在导航栏添加测验入口，显示题库元数据（总题量、预计用时、难度）

**Independent Test**: 检查 UI 布局中是否存在入口链接及元数据绑定

### Tests for User Story 4 (MANDATORY) ⚠️

- [ ] T035 [P] [US4] 创建 QuizMetaInfo 组件单元测试 in `frontend/src/components/quiz/__tests__/QuizMetaInfo.test.tsx`
- [ ] T036 [P] [US4] 创建导航栏测验入口单元测试 in `frontend/src/components/layout/__tests__/Sidebar.test.tsx`

### Implementation for User Story 4

- [ ] T037 [P] [US4] 创建 QuizMetaInfo 组件（总题量、预计用时、难度）in `frontend/src/components/quiz/QuizMetaInfo.tsx`
- [ ] T038 [US4] 更新测验开始页集成 QuizMetaInfo 组件 in `frontend/src/pages/quiz/[topic]/[chapter].tsx`
- [ ] T039 [US4] 添加侧边栏/导航栏测验快速入口 in `frontend/src/components/layout/Sidebar.tsx`
- [ ] T040 [US4] 集成骨架屏到测验加载过程 in `frontend/src/components/quiz/QuizViewer.tsx`

**Checkpoint**: 用户故事 4 完成 - 可独立测试导航入口和元数据展示

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: 跨用户故事的改进和收尾工作

- [ ] T041 [P] 验证所有代码注释和用户文档为中文
- [ ] T042 [P] 更新 backend/README.md 添加新 API 文档
- [ ] T043 [P] 更新 frontend/README.md 添加新组件说明
- [ ] T044 安全加固：验证 Submit API 输入校验（答案数量匹配题目数量）in `backend/internal/controller/quiz/quiz_submit.go`
- [ ] T045 [P] 边缘情况处理：解析字段缺失时显示"暂无解析" in `frontend/src/components/quiz/QuizViewer.tsx`
- [ ] T046 [P] 边缘情况处理：选项超过 26 个时的标签生成逻辑 in `frontend/src/components/quiz/QuizQuestionCard.tsx`
- [ ] T047 网络超时重试机制 in `frontend/src/services/quizApi.ts`
- [ ] T048 运行 quickstart.md 验证所有功能

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖 - 可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成 - **阻塞所有用户故事**
- **User Stories (Phase 3-6)**: 全部依赖 Foundational 阶段完成
  - US1 和 US2 为 P1 优先级，建议优先完成
  - US3 依赖 US2 的提交功能（需要有历史数据）
  - US4 可与其他故事并行
- **Polish (Phase 7)**: 依赖所有期望的用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: Foundational 完成后可开始 - 无其他故事依赖
- **User Story 2 (P1)**: Foundational 完成后可开始 - 与 US1 可并行
- **User Story 3 (P2)**: 依赖 US2 的 Submit API 更新（需要保存历史记录）
- **User Story 4 (P3)**: Foundational 完成后可开始 - 与其他故事独立

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Backend before frontend integration
- Story complete before moving to next priority

### Parallel Opportunities

- **Phase 1**: T001, T002 可并行
- **Phase 2**: T004, T005 可并行; T008 与后端任务可并行
- **Phase 3**: T009, T010 可并行; T011, T012 可并行
- **Phase 4**: T015, T016 可并行; T019, T020 可并行
- **Phase 5**: T023, T024, T025, T026 可并行; T027, T028 顺序执行; T031, T032 可并行
- **Phase 6**: T035, T036 可并行; T037 独立
- **Phase 7**: T041, T042, T043, T045, T046 可并行

---

## Parallel Example: User Story 1

```bash
# 并行启动 US1 的所有测试:
Task: "创建 QuizQuestionCard 组件单元测试 in frontend/src/components/quiz/__tests__/QuizQuestionCard.test.tsx"
Task: "创建 SubmitConfirmModal 组件单元测试 in frontend/src/components/quiz/__tests__/SubmitConfirmModal.test.tsx"

# 并行启动 US1 的组件实现:
Task: "创建 QuizQuestionCard 组件 in frontend/src/components/quiz/QuizQuestionCard.tsx"
Task: "创建 SubmitConfirmModal 组件 in frontend/src/components/quiz/SubmitConfirmModal.tsx"
```

---

## Parallel Example: User Story 3

```bash
# 并行启动 US3 的所有契约测试:
Task: "创建 GET /quiz/history API 契约测试 in backend/internal/controller/quiz/quiz_history_test.go"
Task: "创建 GET /quiz/history/{sessionId} API 契约测试 in backend/internal/controller/quiz/quiz_review_test.go"
Task: "创建 QuizHistoryPage 组件单元测试 in frontend/src/components/quiz/__tests__/QuizHistoryPage.test.tsx"
Task: "创建 QuizReviewPage 组件单元测试 in frontend/src/components/quiz/__tests__/QuizReviewPage.test.tsx"

# 并行启动 US3 的前端页面:
Task: "创建 QuizHistoryPage 页面组件 in frontend/src/pages/quiz/history/index.tsx"
Task: "创建 QuizReviewPage 页面组件 in frontend/src/pages/quiz/history/[sessionId].tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 + 2)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (防误触与有序标签)
4. Complete Phase 4: User Story 2 (百分制与题型标识)
5. **STOP and VALIDATE**: 测试 US1 + US2 独立功能
6. Deploy/demo if ready - 这是可交付的 MVP

### Incremental Delivery

1. Setup + Foundational → 基础设施就绪
2. Add User Story 1 + 2 → Test → Deploy/Demo (**MVP!**)
3. Add User Story 3 → Test → Deploy/Demo (历史回顾)
4. Add User Story 4 → Test → Deploy/Demo (入口与元数据)
5. 每个故事独立增加价值，不破坏之前的功能

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (前端交互优化)
   - Developer B: User Story 2 (结果页增强)
3. After US2 complete:
   - Developer A: User Story 3 (历史回顾 - 后端)
   - Developer B: User Story 3 (历史回顾 - 前端)
   - Developer C: User Story 4 (入口与元数据)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- **Web-only Exception**: 历史回顾模式 (US3) 仅限 Web 端，CLI 保持现有行为
- **Persistence**: 使用 SQLite + gdb ORM，事务保证 QuizSession + QuizAttempt 原子性保存

