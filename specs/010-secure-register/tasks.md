# Tasks: 安全注册与默认管理员（登录后强制改密）

**Input**: Design documents from `/specs/010-secure-register/`  
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, quickstart.md  
**Tests**: 覆盖率需≥80%，后端与前端核心路径均需测试任务。  
**Organization**: 任务按用户故事分组以便独立实现与验证。

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 本阶段无额外任务，可直接进入基础与用户故事实现。

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 共享前后端口令校验与测试基线，完成后方可进入各用户故事

- [X] T002 [P] 更新口令策略校验（≥8 位且包含大小写/数字/特殊字符）并补充单元测试 `backend/internal/pkg/password/password.go` `backend/internal/pkg/password/password_test.go`
- [X] T003 [P] 调整域服务校验逻辑复用新策略并完善测试 `backend/internal/domain/user/service.go` `backend/internal/domain/user/service_test.go`
- [X] T004 [P] 前端表单校验升级为新策略（含特殊字符）`frontend/components/auth/RegisterForm.tsx` `frontend/components/auth/LoginForm.tsx`
- [X] T005 [P] 更新前端表单测试以覆盖新策略与错误提示 `frontend/tests/components/RegisterForm.test.tsx` `frontend/tests/components/LoginForm.test.tsx`

**Checkpoint**: 新口令策略在前后端与测试中生效。

---

## Phase 3: User Story 1 - 已登录管理员创建新用户 (Priority: P1) 🎯 MVP

**Goal**: 仅已登录且具管理员权限的用户可调用注册接口；无令牌/无权限请求被拒并审计。

**Independent Test**: 使用管理员 JWT 调用注册成功；缺少或非管理员 JWT 被拒绝且无账户创建。

### Tests for User Story 1 (MANDATORY) ⚠️

- [X] T006 [P] [US1] 合同测试：注册接口需携带有效管理员 JWT，否则 401/403 `backend/tests/contract/auth_api_test.go`
- [X] T007 [P] [US1] 集成测试：管理员成功注册与非管理员/无令牌失败场景 `backend/tests/integration/auth_flow_test.go`

### Implementation for User Story 1

- [X] T008 [US1] 后端为注册路由增加鉴权与管理员校验，拒绝未授权并返回明确错误 `backend/internal/app/http_server/handler/auth.go` `backend/internal/app/http_server/router.go`
- [X] T009 [US1] 服务层支持管理员校验与审计记录（注册成功/拒绝）`backend/internal/domain/user/service.go` `backend/internal/infrastructure/repository/user_repo.go`
- [X] T010 [US1] 前端注册入口改为受保护的管理员页面/路由并使用现有会话 Token 调用注册接口 `frontend/components/auth/RegisterForm.tsx` `frontend/app/(protected)/...`
- [X] T011 [P] [US1] 前端 API 调用补充 Authorization 头与错误提示映射 `frontend/lib/api.ts` `frontend/contexts/AuthContext.tsx`

**Checkpoint**: 管理员可完成注册，未授权路径全部被拒绝且有审计。

---

## Phase 4: User Story 2 - 启动时保障默认管理员存在 (Priority: P1)

**Goal**: 启动自动创建默认管理员（admin/gostudy@123）且幂等，不覆盖已有账号；标记需改密并审计。

**Independent Test**: 空用户库启动后存在 admin 且 must_change_password=true；已有 admin 时不会重置密码或重复创建。

### Tests for User Story 2 (MANDATORY) ⚠️

- [X] T012 [P] [US2] 集成测试：空库启动创建默认 admin；已有 admin 时不被覆盖 `backend/tests/integration/auth_flow_test.go`（或新增 `default_admin_test.go`）

### Implementation for User Story 2

- [X] T013 [US2] 启动钩子检查/创建默认管理员（含 must_change_password 标记，密码符合策略）`backend/internal/app/http_server/server.go` `backend/internal/domain/user/service.go`
- [X] T014 [US2] 默认管理员创建写入审计，避免并发重复创建（唯一约束/事务）`backend/internal/domain/user/service.go` `backend/internal/infrastructure/repository/user_repo.go`
**Checkpoint**: 启动后必有可用管理员且不被重复/覆盖。

---

## Phase 5: User Story 3 - 默认管理员首次登录强制改密 (Priority: P2)

**Goal**: 首次登录的默认管理员必须完成改密后才能访问其他功能；改密后旧口令与旧令牌全部失效。

**Independent Test**: 默认口令登录返回 needPasswordChange；未改密访问业务接口被拦截；改密成功后需重新登录且旧口令无法使用。

### Tests for User Story 3 (MANDATORY) ⚠️

- [X] T016 [P] [US3] 集成测试：登录返回 needPasswordChange，改密前业务请求被阻断；改密后旧令牌失效 `backend/tests/integration/auth_flow_test.go`
- [X] T017 [P] [US3] 前端集成/组件测试：needPasswordChange 时重定向改密页，改密成功后跳转登录且旧口令失败 `frontend/tests/integration/auth.test.tsx`

### Implementation for User Story 3

- [X] T018 [US3] 登录响应返回 needPasswordChange，增加中间件阻断需改密用户访问除改密外的接口 `backend/internal/app/http_server/handler/auth.go` `backend/internal/app/http_server/middleware/*.go`
- [X] T019 [US3] 实现改密接口：校验旧密码与新策略、清除旧刷新令牌、更新 must_change_password=false 并审计 `backend/internal/app/http_server/handler/auth.go` `backend/internal/domain/user/service.go`
- [X] T020 [US3] 前端改密页面与流程：表单校验新策略，成功后清理本地令牌并跳转登录 `frontend/app/(auth)/change-password/page.tsx` `frontend/components/auth/ChangePasswordForm.tsx`
- [X] T021 [P] [US3] 前端 AuthContext/路由守卫支持 needPasswordChange 状态与改密 API 调用 `frontend/contexts/AuthContext.tsx` `frontend/lib/auth.ts`

### Audit Coverage for FR-007 (Tests)

- [X] T025 [P] [US1] 审计测试：管理员注册成功与未授权注册被拒均写入审计事件 `backend/tests/integration/auth_flow_test.go`
- [X] T026 [P] [US2] 审计测试：空库启动创建默认管理员事件被记录且无重复写入 `backend/tests/integration/default_admin_test.go`
- [X] T027 [P] [US3] 审计测试：首次改密成功记录事件，旧令牌失效被拒绝操作也有审计 `backend/tests/integration/auth_flow_test.go`

**Checkpoint**: 改密流程闭环，未改密无法访问业务，改密后必须重新登录。

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: 收尾、文档与质量保障

- [X] T022 [P] 更新文档：密码策略与强制改密说明（`README.md`、`specs/010-secure-register/quickstart.md`）
- [X] T023 [P] 审计/安全巡检：核查关键动作日志与告警配置 `backend/internal/app/http_server/handler/auth.go` `backend/logs/`
- [X] T024 整体质量门禁：后端 `go test ./...`、前端 `npm test`/`npm run lint`、必要时 `npm run build`

---

## Dependencies & Execution Order

- Phase 2 → User Stories → Polish（Phase 1 无阻塞任务）。User Stories 可在 Phase 2 完成后并行，但建议按优先级 US1 (P1) → US2 (P1) → US3 (P2)。
- 任务依赖：US1 依赖口令策略（T002-T005）；US2 依赖口令策略；US3 依赖口令策略与默认管理员存在；审计测试 T025-T027 依赖对应故事实现。

## Parallel Examples

- 基础并行：T002/T003/T004/T005 可并行；US1 的 T006/T007 可并行，T010/T011 可与后端实现并行。
- 故事并行：Phase 2 完成后，US1 与 US2 可并行开发，US3 建议在 US2 完成后进行（需默认管理员流程）。

## Implementation Strategy

- MVP：完成 Phase 1-2 + US1（管理员带 JWT 注册）即可形成可演示版本。  
- 增量：US2 确保启动可登录，US3 完成安全闭环。  
- 每个用户故事完成后执行对应测试与 quickstart 步骤以验证独立可用性。

