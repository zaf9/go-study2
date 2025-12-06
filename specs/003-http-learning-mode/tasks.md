# Tasks: HTTP学习模式

**Input**: Design documents from `/specs/003-http-learning-mode/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Per the constitution, features MUST have at least 80% unit test coverage. Test tasks are **MANDATORY**.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

项目采用单一Go项目结构：
- **源代码**: `internal/app/`
- **测试**: `tests/`
- **配置**: 项目根目录的`config.yaml`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目初始化和基础结构搭建

- [ ] T001 创建HTTP服务配置文件 `config.yaml` 在项目根目录
- [ ] T002 创建 `internal/config/` 目录结构
- [ ] T003 创建 `internal/app/http_server/` 目录结构（含handler/和middleware/子目录）
- [ ] T004 创建 `tests/unit/` 和 `tests/integration/` 目录结构
- [ ] T005 [P] 在 `go.mod` 中验证GoFrame v2.9.5依赖已存在

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 核心基础设施，必须在所有用户故事之前完成

**⚠️ CRITICAL**: 所有用户故事工作必须等待此阶段完成

- [ ] T006 实现配置加载和验证 `internal/config/config.go`
- [ ] T007 [P] 为配置加载编写单元测试 `tests/unit/config_test.go`
- [ ] T008 重构 `internal/app/lexical_elements/comments.go` 添加 `GetCommentsContent()` 函数
- [ ] T009 [P] 重构 `internal/app/lexical_elements/tokens.go` 添加 `GetTokensContent()` 函数
- [ ] T010 [P] 重构 `internal/app/lexical_elements/semicolons.go` 添加 `GetSemicolonsContent()` 函数
- [ ] T011 [P] 重构 `internal/app/lexical_elements/identifiers.go` 添加 `GetIdentifiersContent()` 函数
- [ ] T012 [P] 重构 `internal/app/lexical_elements/keywords.go` 添加 `GetKeywordsContent()` 函数
- [ ] T013 [P] 重构 `internal/app/lexical_elements/operators.go` 添加 `GetOperatorsContent()` 函数
- [ ] T014 [P] 重构 `internal/app/lexical_elements/integers.go` 添加 `GetIntegersContent()` 函数
- [ ] T015 [P] 重构 `internal/app/lexical_elements/floats.go` 添加 `GetFloatsContent()` 函数
- [ ] T016 [P] 重构 `internal/app/lexical_elements/imaginary.go` 添加 `GetImaginaryContent()` 函数
- [ ] T017 [P] 重构 `internal/app/lexical_elements/runes.go` 添加 `GetRunesContent()` 函数
- [ ] T018 [P] 重构 `internal/app/lexical_elements/strings.go` 添加 `GetStringsContent()` 函数
- [ ] T019 为所有重构的Get*Content()函数编写单元测试 `tests/unit/lexical_refactor_test.go`
- [ ] T020 实现HTTP服务器初始化 `internal/app/http_server/server.go`
- [ ] T021 [P] 实现日志中间件 `internal/app/http_server/middleware/logger.go`
- [ ] T022 [P] 实现格式转换中间件 `internal/app/http_server/middleware/format.go`

**Checkpoint**: 基础设施就绪 - 用户故事实现现在可以并行开始

---

## Phase 3: User Story 1 - 命令行交互式学习（默认模式） (Priority: P1) 🎯 MVP

**Goal**: 保持现有命令行交互模式正常工作，验证重构后的Get*Content()函数

**Independent Test**: 不带参数运行程序，验证命令行菜单显示和内容展示正常

### Tests for User Story 1 (MANDATORY) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T023 [P] [US1] 为命令行模式编写集成测试 `tests/integration/cli_mode_test.go`
- [ ] T024 [P] [US1] 验证Display*()函数调用Get*Content()的单元测试 `tests/unit/display_wrapper_test.go`

### Implementation for User Story 1

- [ ] T025 [US1] 更新所有Display*()函数调用对应的Get*Content()并打印（11个文件）
- [ ] T026 [US1] 验证 `internal/app/lexical_elements/lexical_elements.go` 中的DisplayMenu()函数仍正常工作
- [ ] T027 [US1] 测试命令行模式：运行程序不带参数，验证所有章节内容正确显示

**Checkpoint**: 命令行模式应完全功能正常，与重构前行为一致

---

## Phase 4: User Story 2 - HTTP服务模式学习 (Priority: P2)

**Goal**: 实现HTTP服务模式，通过POST接口访问学习内容

**Independent Test**: 使用 `-d` 参数启动，通过curl访问HTTP端点验证JSON/HTML响应

### Tests for User Story 2 (MANDATORY) ⚠️

- [ ] T028 [P] [US2] 为 `/api/v1/topics` 接口编写集成测试 `tests/integration/http_topics_test.go`
- [ ] T029 [P] [US2] 为 `/api/v1/topic/lexical_elements` 接口编写集成测试 `tests/integration/http_lexical_test.go`
- [ ] T030 [P] [US2] 为章节接口编写集成测试 `tests/integration/http_chapter_test.go`

### Implementation for User Story 2

- [ ] T031 [P] [US2] 实现Topics处理器 `internal/app/http_server/handler/topics.go`
- [ ] T032 [P] [US2] 实现Lexical Elements菜单处理器 `internal/app/http_server/handler/lexical.go`
- [ ] T033 [US2] 实现章节内容处理器 `internal/app/http_server/handler/chapter.go`（调用Get*Content()函数）
- [ ] T034 [US2] 实现路由注册 `internal/app/http_server/router.go`（所有接口使用POST方法）
- [ ] T035 [US2] 在 `main.go` 中添加命令行参数解析（-d 和 --daemon）
- [ ] T036 [US2] 在 `main.go` 中实现HTTP模式启动逻辑
- [ ] T037 [US2] 实现优雅关闭机制（信号监听）
- [ ] T038 [US2] 测试HTTP模式：启动服务，验证所有接口返回正确的JSON格式内容

**Checkpoint**: HTTP服务模式应完全功能，所有POST接口返回正确内容

---

## Phase 5: User Story 3 - 内容一致性保证 (Priority: P1)

**Goal**: 验证命令行和HTTP两种模式返回相同内容

**Independent Test**: 分别通过两种模式获取相同章节，比较内容一致性

### Tests for User Story 3 (MANDATORY) ⚠️

- [ ] T039 [US3] 编写内容一致性集成测试 `tests/integration/content_consistency_test.go`

### Implementation for User Story 3

- [ ] T040 [US3] 验证所有11个章节在两种模式下内容一致
- [ ] T041 [US3] 添加内容一致性验证到CI/CD流程（如果存在）

**Checkpoint**: 两种模式内容100%一致

---

## Phase 6: User Story 4 - HTTP服务配置灵活性 (Priority: P3)

**Goal**: 支持通过配置文件灵活配置HTTP服务参数

**Independent Test**: 修改config.yaml中的端口和地址，验证服务在新配置下启动

### Tests for User Story 4 (MANDATORY) ⚠️

- [ ] T042 [P] [US4] 为配置验证编写单元测试（缺失必填项） `tests/unit/config_validation_test.go`
- [ ] T043 [P] [US4] 为端口占用场景编写集成测试 `tests/integration/port_conflict_test.go`

### Implementation for User Story 4

- [ ] T044 [US4] 验证配置文件必填项检查（host和port）
- [ ] T045 [US4] 实现端口占用检测和错误提示
- [ ] T046 [US4] 测试不同配置场景（不同端口、不同地址、缺失配置）

**Checkpoint**: 配置灵活性完全实现，错误提示清晰

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: 跨用户故事的改进和完善

- [ ] T047 [P] 验证所有代码注释和用户提示均为中文
- [ ] T048 [P] 验证所有HTTP错误响应格式与请求格式一致（JSON/HTML）
- [ ] T049 [P] 实现HTML响应模板优化 `internal/app/http_server/middleware/format.go`
- [ ] T050 运行所有测试，确保覆盖率≥80%
- [ ] T051 [P] 更新 `README.md` 添加HTTP模式使用说明
- [ ] T052 [P] 创建示例config.yaml文件并添加详细注释
- [ ] T053 性能测试：验证50+并发请求处理能力
- [ ] T054 安全检查：验证输入验证和错误处理
- [ ] T055 按照 `specs/003-http-learning-mode/quickstart.md` 验证所有使用场景

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖 - 可立即开始
- **Foundational (Phase 2)**: 依赖Setup完成 - 阻塞所有用户故事
- **User Stories (Phase 3-6)**: 全部依赖Foundational阶段完成
  - US1 (P1): 可在Foundational后开始 - 无其他故事依赖
  - US2 (P2): 可在Foundational后开始 - 无其他故事依赖
  - US3 (P1): 依赖US1和US2完成（需要两种模式都实现）
  - US4 (P3): 可在Foundational后开始 - 无其他故事依赖
- **Polish (Phase 7)**: 依赖所有期望的用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: Foundational后可开始 - 独立可测试
- **User Story 2 (P2)**: Foundational后可开始 - 独立可测试
- **User Story 3 (P1)**: 需要US1和US2都完成 - 验证两者一致性
- **User Story 4 (P3)**: Foundational后可开始 - 独立可测试

### Within Each User Story

- Tests必须先编写并失败
- 重构的Get*Content()函数 → Display*()函数调用
- HTTP处理器 → 路由注册 → main.go集成
- 核心实现 → 集成测试 → 故事完成

### Parallel Opportunities

- **Phase 1**: 所有标记[P]的任务可并行
- **Phase 2**: T007-T022所有重构任务可并行（不同文件）
- **Phase 3**: US1的测试任务可并行
- **Phase 4**: US2的测试任务和处理器实现可并行
- **Phase 6**: US4的测试任务可并行
- **Phase 7**: 所有标记[P]的任务可并行
- **跨故事**: US1、US2、US4可并行开发（US3需要US1+US2完成）

---

## Parallel Example: Phase 2 (Foundational)

```bash
# 并行重构所有11个Get*Content()函数：
Task: "重构 comments.go 添加 GetCommentsContent()"
Task: "重构 tokens.go 添加 GetTokensContent()"
Task: "重构 semicolons.go 添加 GetSemicolonsContent()"
Task: "重构 identifiers.go 添加 GetIdentifiersContent()"
Task: "重构 keywords.go 添加 GetKeywordsContent()"
Task: "重构 operators.go 添加 GetOperatorsContent()"
Task: "重构 integers.go 添加 GetIntegersContent()"
Task: "重构 floats.go 添加 GetFloatsContent()"
Task: "重构 imaginary.go 添加 GetImaginaryContent()"
Task: "重构 runes.go 添加 GetRunesContent()"
Task: "重构 strings.go 添加 GetStringsContent()"

# 并行实现中间件：
Task: "实现日志中间件 logger.go"
Task: "实现格式转换中间件 format.go"
```

## Parallel Example: Phase 4 (User Story 2)

```bash
# 并行实现所有HTTP处理器：
Task: "实现Topics处理器 topics.go"
Task: "实现Lexical Elements菜单处理器 lexical.go"

# 并行编写集成测试：
Task: "为 /api/v1/topics 接口编写集成测试"
Task: "为 /api/v1/topic/lexical_elements 接口编写集成测试"
Task: "为章节接口编写集成测试"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - 阻塞所有故事)
3. Complete Phase 3: User Story 1（命令行模式）
4. Complete Phase 4: User Story 2（HTTP模式）
5. **STOP and VALIDATE**: 独立测试US1和US2
6. 部署/演示（如果准备好）

### Incremental Delivery

1. Setup + Foundational → 基础就绪
2. Add User Story 1 → 独立测试 → 部署/演示（命令行MVP）
3. Add User Story 2 → 独立测试 → 部署/演示（HTTP模式）
4. Add User Story 3 → 独立测试 → 验证一致性
5. Add User Story 4 → 独立测试 → 配置灵活性
6. 每个故事增加价值而不破坏之前的故事

### Parallel Team Strategy

多开发者协作：

1. 团队一起完成Setup + Foundational
2. Foundational完成后：
   - Developer A: User Story 1（命令行模式）
   - Developer B: User Story 2（HTTP模式）
   - Developer C: User Story 4（配置管理）
3. Developer A + B完成后：
   - Developer A or B: User Story 3（一致性验证）
4. 故事独立完成和集成

---

## Task Summary

- **Total Tasks**: 55
- **Setup Phase**: 5 tasks
- **Foundational Phase**: 17 tasks (11个重构 + 配置 + 中间件 + 测试)
- **User Story 1**: 5 tasks (2 tests + 3 implementation)
- **User Story 2**: 11 tasks (3 tests + 8 implementation)
- **User Story 3**: 3 tasks (1 test + 2 implementation)
- **User Story 4**: 5 tasks (2 tests + 3 implementation)
- **Polish Phase**: 9 tasks

### Parallel Opportunities

- **Phase 2**: 15个任务可并行（所有重构+中间件）
- **Phase 3**: 2个测试任务可并行
- **Phase 4**: 5个任务可并行（3个测试+2个处理器）
- **Phase 6**: 2个测试任务可并行
- **Phase 7**: 5个任务可并行
- **跨故事**: US1、US2、US4可并行（3个故事）

### Test Coverage Target

- **目标**: ≥80% 单元测试覆盖率
- **测试任务**: 13个（包含单元测试和集成测试）
- **关键测试点**:
  - 配置加载和验证
  - 所有Get*Content()函数
  - HTTP处理器
  - 中间件
  - 内容一致性
  - 端口冲突处理

---

## Notes

- [P] 任务 = 不同文件，无依赖，可并行
- [Story] 标签将任务映射到特定用户故事，便于追踪
- 每个用户故事应独立可完成和可测试
- 在实现前验证测试失败
- 每个任务或逻辑组后提交
- 在任何检查点停止以独立验证故事
- 避免：模糊任务、相同文件冲突、破坏独立性的跨故事依赖
