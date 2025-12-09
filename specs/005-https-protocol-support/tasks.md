# Tasks: HTTPS 协议支持

**Input**: Design documents from `/specs/005-https-protocol-support/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Tests**: Per the constitution, features MUST have at least 80% unit test coverage. Test tasks are **MANDATORY**.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: Go 项目，使用 `internal/` 目录结构
- Paths based on plan.md structure

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目配置迁移和基础结构准备

- [X] T001 迁移配置文件从根目录到 configs/config.yaml
- [X] T002 [P] 创建证书目录结构 configs/certs/
- [X] T003 [P] 更新 GoFrame 配置加载路径以使用 configs/ 目录

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 配置结构扩展，所有用户故事依赖此阶段

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T004 添加 HttpConfig 结构体到 internal/config/config.go
- [ ] T005 添加 HttpsConfig 结构体到 internal/config/config.go
- [ ] T006 更新 Config 结构体，集成 Http 和 Https 配置
- [X] T004 添加 HttpConfig 结构体到 internal/config/config.go
- [X] T005 添加 HttpsConfig 结构体到 internal/config/config.go
- [X] T006 更新 Config 结构体，集成 Http 和 Https 配置
- [X] T007 [P] 添加配置结构单元测试到 internal/config/config_test.go

**Checkpoint**: 配置结构就绪 - 用户故事实现可以开始

---

## Phase 3: User Story 1 - 启用 HTTPS 安全服务 (Priority: P1) 🎯 MVP

**Goal**: 通过配置启用 HTTPS 服务，使用 TLS 1.2+ 加密通信

**Independent Test**: 设置 `https.enabled = true` 并提供有效证书，服务以 HTTPS 模式启动

### Tests for User Story 1 (MANDATORY) ⚠️

- [X] T008 [P] [US1] HTTPS 服务器启动单元测试 in internal/app/http_server/server_test.go
- [X] T009 [P] [US1] HTTPS 模式集成测试 in tests/integration/https_mode_test.go
- [X] T033 [P] [US1] 自签名证书握手集成测试（含 root CA/跳过校验开关说明）in tests/integration/https_mode_test.go
- [X] T034 [P] [US1] 启用 HTTPS 时禁用 HTTP 端口的集成测试 in tests/integration/https_mode_test.go

### Implementation for User Story 1

- [X] T010 [US1] 创建 TLS 配置函数（MinVersion TLS 1.2）in internal/app/http_server/server.go
- [X] T011 [US1] 实现 HTTPS 启动逻辑，使用 EnableHTTPS() in internal/app/http_server/server.go
- [X] T012 [US1] 更新 NewServer 函数支持 HTTPS 模式切换 in internal/app/http_server/server.go
- [X] T013 [US1] 添加 HTTPS 模式日志输出 in internal/app/http_server/server.go
- [X] T035 [US1] 实现自签名证书加载与可配置信任策略（root CA 或跳过校验开关，仅限测试/开发）in internal/app/http_server/server.go
- [X] T036 [US1] 启用 HTTPS 时显式禁用 HTTP 监听并输出清晰日志 in internal/app/http_server/server.go

**Checkpoint**: HTTPS 服务可独立启动和测试

---

## Phase 4: User Story 2 - 保持 HTTP 模式兼容 (Priority: P2)

**Goal**: 确保 HTTP 模式作为默认行为，向后兼容

**Independent Test**: 设置 `https.enabled = false` 或不配置，服务以 HTTP 模式启动

### Tests for User Story 2 (MANDATORY) ⚠️

- [X] T014 [P] [US2] HTTP 模式单元测试 in internal/app/http_server/server_test.go
- [X] T015 [P] [US2] HTTP 模式向后兼容集成测试 in tests/integration/http_mode_test.go

### Implementation for User Story 2

- [ ] T016 [US2] 确保未配置 https 时默认使用 HTTP 模式 in internal/app/http_server/server.go
- [ ] T017 [US2] 验证现有 HTTP 端点在两种模式下行为一致

**Checkpoint**: HTTP 和 HTTPS 模式均可独立工作

---

## Phase 5: User Story 3 - 证书路径可配置 (Priority: P2)

**Goal**: 支持灵活配置证书文件路径

**Independent Test**: 配置不同的 certFile 和 keyFile 路径，验证加载正确

### Tests for User Story 3 (MANDATORY) ⚠️

- [X] T018 [P] [US3] 证书路径配置单元测试 in internal/config/config_test.go
- [X] T019 [P] [US3] 相对路径和绝对路径解析测试 in internal/config/config_test.go

### Implementation for User Story 3

- [ ] T020 [US3] 实现证书路径解析逻辑（相对路径/绝对路径）in internal/config/config.go
- [ ] T021 [US3] 将证书路径传递给 HTTPS 启动逻辑 in internal/app/http_server/server.go

**Checkpoint**: 证书路径配置功能完整

---

## Phase 6: User Story 4 - 证书文件错误处理 (Priority: P3)

**Goal**: 证书配置错误时提供清晰友好的错误提示

**Independent Test**: 配置不存在的证书路径，验证错误提示包含路径信息

### Tests for User Story 4 (MANDATORY) ⚠️

- [X] T022 [P] [US4] 证书文件不存在错误测试 in internal/config/config_test.go
- [X] T023 [P] [US4] 私钥文件不存在错误测试 in internal/config/config_test.go
- [X] T024 [P] [US4] 证书路径缺失错误测试 in internal/config/config_test.go
- [X] T037 [P] [US4] 证书/私钥权限不足错误测试 in internal/config/config_test.go
- [X] T038 [P] [US4] 证书与私钥不匹配错误测试 in internal/config/config_test.go
- [ ] T039 [P] [US4] 过期证书启动行为测试（不需要，用户确认无需执行）
- [X] T040 [P] [US4] 端口被占用时的错误消息一致性测试 in internal/app/http_server/server_test.go

### Implementation for User Story 4

- [X] T025 [US4] 实现 HTTPS 配置验证逻辑 in internal/config/config.go
- [X] T026 [US4] 添加证书文件存在性检查 in internal/config/config.go
- [X] T027 [US4] 添加友好的中文错误消息 in internal/config/config.go
- [X] T041 [US4] 为自签名/CA 配置添加可选跳过客户端校验开关及风险提示 in internal/config/config.go

**Checkpoint**: 所有错误情况均有友好提示

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: 文档更新和收尾工作

- [X] T028 [P] 验证所有代码注释和文档均为中文
- [X] T029 [P] 更新 README.md 添加 HTTPS 配置说明
- [X] T030 [P] 更新 configs/config.yaml 添加 https 配置示例（注释形式）
- [X] T031 运行 quickstart.md 验证，确保文档准确
- [X] T032 运行完整测试套件，验证覆盖率 ≥80%
- [X] T042 [P] 验证/记录 CLI 学习模式兼容性（如无网络依赖则在文档声明）in quickstart.md & tests
- [X] T043 [P] 执行 go fmt / go vet / golint / go mod tidy 质量门禁并记录结果
- [X] T044 [P] 编写并运行协议切换耗时测量脚本（确保 ≤30 秒）in tests/integration/https_mode_test.go 或脚本目录

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - 立即开始
- **Foundational (Phase 2)**: Depends on Setup - **BLOCKS 所有用户故事**
- **User Stories (Phase 3-6)**: 依赖 Foundational 完成后可并行或按优先级顺序执行
- **Polish (Phase 7)**: 依赖所有用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: Foundational 完成后可开始 - 核心 HTTPS 功能
- **User Story 2 (P2)**: Foundational 完成后可开始 - 可与 US1 并行
- **User Story 3 (P2)**: 依赖 US1 中的证书加载逻辑
- **User Story 4 (P3)**: 依赖 US3 中的证书路径配置

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- 配置结构优先于服务逻辑
- Story complete before moving to next priority

### Parallel Opportunities

**Phase 1 并行**:
- T002 和 T003 可并行执行

**Phase 2 并行**:
- T007 (测试) 可在 T004-T006 完成后立即执行

**User Story 测试并行**:
- 每个用户故事内的测试任务（标记 [P]）可并行执行
- 不同用户故事的测试可并行执行

**Polish 并行**:
- T028, T029, T030 可并行执行

---

## Parallel Example: User Story 1

```text
           T008 (test)
          /           \
T004-T007             T010 → T011 → T012 → T013
(foundation)           \
          \           T009 (integration test - after impl)
           \
            → US1 Complete
```

---

## Implementation Strategy

### MVP Scope

**MVP = User Story 1 (P1)**: 启用 HTTPS 安全服务

完成 Phase 1-3 即可交付可用的 HTTPS 功能。

### Incremental Delivery

1. **Increment 1 (MVP)**: Phase 1 + Phase 2 + Phase 3 (US1)
2. **Increment 2**: Phase 4 (US2) - HTTP 兼容性
3. **Increment 3**: Phase 5 (US3) - 证书路径灵活性
4. **Increment 4**: Phase 6 (US4) - 错误处理完善
5. **Final**: Phase 7 - 收尾和文档

### Task Summary

| Phase | Task Count | Parallel Tasks |
|-------|------------|----------------|
| Setup | 3 | 2 |
| Foundational | 4 | 1 |
| US1 (P1) | 8 | 4 |
| US2 (P2) | 4 | 2 |
| US3 (P2) | 4 | 2 |
| US4 (P3) | 10 | 5 |
| Polish | 8 | 6 |
| **Total** | **41** | **22** |
