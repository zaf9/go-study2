# Tasks: Dashboard 首页功能

**Input**: Design documents from `/specs/015-dashboard-homepage/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Per the constitution, features MUST have at least 80% unit test coverage. Test tasks are **MANDATORY**.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `backend/`, `frontend/`
- Frontend paths: `frontend/app/`, `frontend/components/`, `frontend/lib/`
- Backend paths: `backend/internal/`, `backend/api/`

## Constitution Guardrails

- 所有注释与用户文档相关任务必须产出中文内容,且保持清晰一致(Principle V/XV)。
- 需规划达到>=80%测试覆盖,各包包含 *_test.go 与示例,前端核心组件同样达标(Principle III/XXI/XXXVI)。
- 目录/文件/函数保持单一职责与可预测结构,遵循标准 Go 布局(仅根目录 main, go.mod/go.sum 完整)并补齐包 README(Principle IV/VIII/XVIII/XIX)。
- 外部依赖与复杂度最小化,错误处理显式,避免 YAGNI(Principle II/VI/IX)。
- 完成后需包含更新 README 等文档的任务(Principle XI)。

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T000 Confirm GoFrame v2.9.5 WebSocket support - check documentation and test basic WebSocket endpoint; if not available, proceed with T001
- [X] T001 [Conditional] Add gorilla/websocket dependency in backend/go.mod (only if T000 confirms GoFrame lacks WebSocket support)
- [X] T002 [P] Create frontend Dashboard directory structure: frontend/app/(protected)/dashboard/
- [X] T003 [P] Create frontend Dashboard components directory: frontend/app/(protected)/dashboard/components/
- [X] T004 [P] Create frontend types file: frontend/types/dashboard.ts
- [X] T005 [P] Create frontend WebSocket utilities: frontend/lib/websocket.ts
- [X] T006 [P] Create frontend time formatting utilities: frontend/lib/utils/time.ts
- [X] T007 [P] Create frontend progress calculation utilities: frontend/lib/utils/progress.ts
- [X] T008 [P] Create backend WebSocket directory: backend/internal/websocket/
- [X] T009 [P] Create backend test directories: backend/tests/controller/, backend/tests/websocket/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T010 Implement WebSocket Hub in backend/internal/websocket/hub.go (manages user connection pool, supports broadcasting by user_id, handles connection lifecycle)
- [X] T011 Implement WebSocket Client in backend/internal/websocket/client.go (handles individual client connection, message read/write, ping/pong heartbeat)
- [X] T012 Implement WebSocket events definition in backend/internal/websocket/events.go (defines progress_updated and quiz_completed event structures, see contracts/websocket-events.md)
- [X] T013 Add WebSocket route in backend/api/v1/websocket.go (handles WebSocket upgrade, authentication verification, client registration)
- [X] T014 [P] Implement WebSocket Provider in frontend/components/providers/WebSocketProvider.tsx
- [X] T015 [P] Add TypeScript types for Dashboard data in frontend/types/dashboard.ts
- [X] T016 [P] Implement time formatting utility in frontend/lib/utils/time.ts
- [X] T017 [P] Implement progress calculation utility in frontend/lib/utils/progress.ts
- [X] T018 [P] Implement WebSocket client wrapper in frontend/lib/websocket.ts
- [X] T019 Update root page to redirect to /dashboard in frontend/app/page.tsx
- [X] T020 Update Sidebar component to link "首页" to /dashboard in frontend/components/layout/Sidebar.tsx

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - 学习状态快速概览 (Priority: P1) 🎯 MVP

**Goal**: 用户登录后立即看到学习概况（欢迎信息、累计学习天数、总章节完成进度、整体完成百分比）

**Independent Test**: 访问 Dashboard 页面并验证显示的统计数据（学习天数、完成章节数、百分比）与数据库中的实际记录一致

### Tests for User Story 1 (MANDATORY) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T021 [P] [US1] Unit test for WelcomeHeader component in frontend/__tests__/dashboard/WelcomeHeader.test.tsx
- [X] T022 [P] [US1] Unit test for StatsCards component in frontend/__tests__/dashboard/StatsCards.test.tsx
- [X] T023 [P] [US1] Backend unit test for学习天数计算 in backend/tests/service/progress_service_test.go
- [X] T024 [P] [US1] Backend integration test for Dashboard stats API in backend/tests/controller/progress_controller_test.go

### Implementation for User Story 1

- [X] T025 [P] [US1] Implement学习天数计算逻辑 in backend/internal/service/progress_service.go
- [X] T026 [P] [US1] Implement Dashboard stats calculation in backend/internal/service/progress_service.go
- [X] T027 [P] [US1] Create WelcomeHeader component in frontend/app/(protected)/dashboard/components/WelcomeHeader.tsx
- [X] T028 [P] [US1] Create StatsCards component in frontend/app/(protected)/dashboard/components/StatsCards.tsx
- [X] T029 [US1] Create Dashboard main page with SSR data fetching in frontend/app/(protected)/dashboard/page.tsx
- [X] T030 [US1] Add loading state in frontend/app/(protected)/dashboard/loading.tsx
- [X] T031 [US1] Add error boundary in frontend/app/(protected)/dashboard/error.tsx with retry button (implements FR-021)
- [X] T031-Test [US1] Test error handling: simulate API failure, verify error message display, verify retry button triggers re-fetch
- [X] T032 [US1] Integrate WelcomeHeader and StatsCards into Dashboard page
- [X] T033 [US1] Add WebSocket event handling for progress updates in Dashboard page

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - 一键继续学习 (Priority: P1)

**Goal**: 一键继续上次未完成的学习内容（显示最后一次学习的主题和章节，点击后直接跳转到该章节的学习页面）

**Independent Test**: 记录用户最后访问的章节，然后在 Dashboard 上点击"继续学习"按钮，验证是否正确跳转到该章节页面

### Tests for User Story 2 (MANDATORY) ⚠️

- [X] T034 [P] [US2] Unit test for QuickContinue component in frontend/__tests__/dashboard/QuickContinue.test.tsx
- [X] T035 [P] [US2] Backend unit test for /api/v1/progress/last endpoint in backend/tests/controller/progress_controller_test.go
- [X] T036 [P] [US2] Backend integration test for GetLastLearningRecord in backend/tests/service/progress_service_test.go

### Implementation for User Story 2

- [X] T037 [P] [US2] Implement GetLastLearningRecord method in backend/internal/service/progress_service.go
- [X] T038 [P] [US2] Implement GetLastLearning controller method in backend/internal/controller/progress_controller.go
- [X] T039 [US2] Add /api/v1/progress/last route in backend/api/v1/progress.go
- [X] T040 [P] [US2] Add getLastLearning API function in frontend/lib/api.ts
- [X] T041 [P] [US2] Create QuickContinue component in frontend/app/(protected)/dashboard/components/QuickContinue.tsx
- [X] T042 [US2] Integrate QuickContinue component into Dashboard page
- [X] T043 [US2] Handle empty state (no learning record) in QuickContinue component

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - 主题进度可视化 (Priority: P2)

**Goal**: 在首页看到各主题的学习进度（列出所有可学习的主题，每个主题显示其完成百分比，使用进度条可视化展示完成度）

**Independent Test**: 创建多个主题的学习记录（部分完成、全部完成、未开始），然后在 Dashboard 上验证每个主题的进度条和百分比是否正确显示

### Tests for User Story 3 (MANDATORY) ⚠️

- [X] T044 [P] [US3] Unit test for TopicProgress component in frontend/__tests__/dashboard/TopicProgress.test.tsx
- [X] T045 [P] [US3] Backend unit test for GetTopicProgressSummary in backend/tests/service/progress_service_test.go

### Implementation for User Story 3

- [X] T046 [P] [US3] Implement GetTopicProgressSummary method in backend/internal/service/progress_service.go
- [X] T047 [P] [US3] Create TopicProgress component in frontend/app/(protected)/dashboard/components/TopicProgress.tsx
- [X] T048 [US3] Integrate TopicProgress component into Dashboard page
- [X] T049 [US3] Add click handler to navigate to topic detail page
- [X] T050 [US3] Handle WebSocket progress_updated event to update topic progress

**Checkpoint**: User Stories 1, 2, AND 3 should all work independently

---

## Phase 6: User Story 4 - 最近测验记录展示 (Priority: P3)

**Goal**: 看到最近的测验记录（列出最近 3-5 条测验记录，每条记录显示主题/章节名称、得分、完成时间）

**Independent Test**: 完成几次测验，然后在 Dashboard 上验证是否显示最近的测验记录（按时间倒序排列，最多 5 条）

### Tests for User Story 4 (MANDATORY) ⚠️

- [X] T051 [P] [US4] Unit test for RecentQuizzes component in frontend/__tests__/dashboard/RecentQuizzes.test.tsx
- [X] T052 [P] [US4] Backend unit test for GetRecentQuizzes in backend/tests/service/quiz_service_test.go

### Implementation for User Story 4

- [X] T053 [P] [US4] Implement GetRecentQuizzes method in backend/internal/service/quiz_service.go
- [X] T054 [P] [US4] Create RecentQuizzes component in frontend/app/(protected)/dashboard/components/RecentQuizzes.tsx
- [X] T055 [US4] Integrate RecentQuizzes component into Dashboard page
- [X] T056 [US4] Implement time formatting (relative/absolute) in RecentQuizzes
- [X] T057 [US4] Handle empty state (no quiz records) in RecentQuizzes component
- [X] T058 [US4] Add click handler to navigate to quiz detail page (if exists)
- [X] T059 [US4] Handle WebSocket quiz_completed event to update recent quizzes

**Checkpoint**: All primary user stories (1-4) should now be independently functional

---

## Phase 7: User Story 5 - 路由与导航调整 (Priority: P1)

**Goal**: 登录后默认跳转到 Dashboard 首页（而非当前的 `/topics`），并且侧边栏的"首页"按钮链接指向 Dashboard 页面

**Independent Test**: 登录系统并验证是否自动跳转到 Dashboard 页面，以及点击侧边栏"首页"按钮是否跳转到 Dashboard

### Implementation for User Story 5

> **NOTE**: This story has no separate implementation as it's already covered by Foundational phase (T019-T020). Tasks below are verification only.

- [X] T060 [US5] **[Verification]** Verify root page redirect to /dashboard in frontend/app/page.tsx (depends on T019)
- [X] T061 [US5] **[Verification]** Verify Sidebar "首页" link points to /dashboard (depends on T020)
- [X] T062 [US5] **[Integration Test]** Test login flow redirects to /dashboard
- [X] T063 [US5] **[Integration Test]** Test navigation from other pages to Dashboard via Sidebar and verify data refresh

**Checkpoint**: All user stories should now be independently functional with proper navigation

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T064 [P] Add responsive design styles for mobile/tablet in frontend/app/(protected)/dashboard/page.tsx
- [ ] T064-A [P] Implement long text truncation with Ant Design Tooltip (implements FR-023): apply to topic names, chapter names in all Dashboard components
- [ ] T064-A-Test [P] Test tooltip behavior: verify text truncates after N characters, verify tooltip shows on hover (300ms delay), verify full text in tooltip
- [ ] T065 [P] Optimize Dashboard page performance (code splitting, lazy loading)
- [ ] T066 [P] Add error handling for WebSocket connection failures
- [ ] T067 [P] Implement WebSocket reconnection with exponential backoff in frontend/lib/websocket.ts
- [ ] T067-Test [P] Unit test for WebSocket reconnection logic: simulate disconnection, verify exponential backoff intervals (1s, 2s, 4s...), verify max 5 retries, verify error display after failure
- [ ] T068 [P] Add loading skeletons for Dashboard components
- [ ] T069 [P] Verify all code comments and user-facing documentation are in Chinese
- [ ] T070 [P] Add database indexes for performance: learning_progress(user_id, last_visited_at)
- [ ] T071 [P] Add database indexes for performance: quiz_records(user_id, completed_at)
- [ ] T072 [P] Security review: Verify authentication on all Dashboard APIs
- [ ] T073 [P] Security review: Verify WebSocket connection authentication
- [ ] T074 [P] Performance testing: Verify Dashboard loads in < 2 seconds
- [ ] T075 [P] Performance testing: Verify WebSocket message latency < 500ms
- [ ] T076 [P] Accessibility review: Verify WCAG 2.1 AA compliance
- [ ] T077 [P] Cross-browser testing: Chrome, Firefox, Safari, Edge
- [ ] T077-A [P] Edge case testing: test all 6 edge cases from spec.md (incomplete data, no topics/chapters, version changes, API timeout, multi-device sync, long names)
- [ ] T078 [P] Update README.md with Dashboard feature documentation
- [ ] T079 [P] Update project structure documentation
- [ ] T080 Run quickstart.md validation and fix any issues

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 stories first: US1, US2, US5, then P2: US3, then P3: US4)
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - May display data from US1 but independently testable
- **User Story 4 (P3)**: Can start after Foundational (Phase 2) - Completely independent
- **User Story 5 (P1)**: Partially completed in Foundational (Phase 2) - Just needs verification

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Backend service methods before controller methods
- Backend routes after controller methods
- Frontend components before page integration
- Core implementation before WebSocket event handling
- Story complete before moving to next priority

### Parallel Opportunities

- **Phase 1 (Setup)**: All tasks marked [P] can run in parallel (T002-T009)
- **Phase 2 (Foundational)**: Tasks T014-T018 can run in parallel (frontend), T010-T013 can run in parallel (backend)
- **User Story Tests**: All test tasks within a story marked [P] can run in parallel
- **User Story Implementation**: 
  - US1: T025-T028 can run in parallel (different files)
  - US2: T037-T038, T040-T041 can run in parallel
  - US3: T046-T047 can run in parallel
  - US4: T053-T054 can run in parallel
- **Different User Stories**: Can be worked on in parallel by different team members after Foundational phase

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task T021: "Unit test for WelcomeHeader component in frontend/__tests__/dashboard/WelcomeHeader.test.tsx"
Task T022: "Unit test for StatsCards component in frontend/__tests__/dashboard/StatsCards.test.tsx"
Task T023: "Backend unit test for学习天数计算 in backend/tests/service/progress_service_test.go"
Task T024: "Backend integration test for Dashboard stats API in backend/tests/controller/progress_controller_test.go"

# Launch all parallel implementation tasks for User Story 1 together:
Task T025: "Implement学习天数计算逻辑 in backend/internal/service/progress_service.go"
Task T026: "Implement Dashboard stats calculation in backend/internal/service/progress_service.go"
Task T027: "Create WelcomeHeader component in frontend/app/(protected)/dashboard/components/WelcomeHeader.tsx"
Task T028: "Create StatsCards component in frontend/app/(protected)/dashboard/components/StatsCards.tsx"
```

---

## Implementation Strategy

### MVP First (User Stories 1, 2, 5 Only - All P1)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (学习状态快速概览)
4. Complete Phase 4: User Story 2 (一键继续学习)
5. Complete Phase 7: User Story 5 (路由与导航调整)
6. **STOP and VALIDATE**: Test all P1 stories independently
7. Deploy/demo MVP

**MVP Scope**: 80 tasks (T001-T063, T064-T080 optional)

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready (T001-T020)
2. Add User Story 1 → Test independently → Deploy/Demo (MVP Core!) (T021-T033)
3. Add User Story 2 → Test independently → Deploy/Demo (T034-T043)
4. Add User Story 5 → Test independently → Deploy/Demo (T060-T063)
5. Add User Story 3 → Test independently → Deploy/Demo (T044-T050)
6. Add User Story 4 → Test independently → Deploy/Demo (T051-T059)
7. Polish & Optimize → Final release (T064-T080)

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (T001-T020)
2. Once Foundational is done:
   - **Developer A**: User Story 1 (T021-T033)
   - **Developer B**: User Story 2 (T034-T043)
   - **Developer C**: User Story 3 (T044-T050)
   - **Developer D**: User Story 4 (T051-T059)
3. Stories complete and integrate independently
4. Team completes Polish together (T064-T080)

---

## Task Summary

- **Total Tasks**: 80
- **Setup Tasks**: 9 (T001-T009)
- **Foundational Tasks**: 11 (T010-T020)
- **User Story 1 Tasks**: 13 (T021-T033) - P1 🎯 MVP
- **User Story 2 Tasks**: 10 (T034-T043) - P1 🎯 MVP
- **User Story 3 Tasks**: 7 (T044-T050) - P2
- **User Story 4 Tasks**: 9 (T051-T059) - P3
- **User Story 5 Tasks**: 4 (T060-T063) - P1 🎯 MVP
- **Polish Tasks**: 17 (T064-T080)

**MVP Scope**: 47 tasks (Setup + Foundational + US1 + US2 + US5)  
**Full Feature**: 80 tasks

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- All Chinese comments and documentation per constitution
- Maintain ≥80% test coverage per constitution
