# Tasks: Variables章节学习完成

**Input**: Design documents from `/specs/006-learn-variables/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/
**Tests**: 必须确保单元/契约/集成测试总覆盖率>=80%。

## Phase 1: Setup (Shared Infrastructure)

- [X] T001 确认工具链可用，优先执行 `./build.bat`，若不存在则运行 `go test ./...` 与 `go build ./...`（仓库根目录）。
- [X] T002 创建章节目录结构 `src/learning/variables/{cli,http}` 与测试目录 `tests/{unit,integration,contract}/learning/variables/`。
- [X] T003 [P] 初始化 gofmt/go vet/golint 常规检查脚本（可放置于仓库根的辅助脚本或文档），确保后续任务统一格式与静态检查。

---

## Phase 2: Foundational (Blocking Prerequisites)

- [X] T004 编写 `src/learning/variables/README.md`，概述子文件职责、CLI/HTTP 双模式说明与目录结构（中文）。
- [X] T005 建立 `src/learning/variables/variables.go` 定义共享的数据结构（内容块、示例、测验项）与公共加载函数。
- [X] T006 [P] 在 `src/learning/variables/content.go` 填充基础常量/模板（主题枚举、通用示例片段占位、校验函数），供各子故事复用。
- [X] T007 [P] 在 `src/learning/variables/cli/menu.go` 搭建章节级菜单框架（支持 list/show/quiz 命令入口，留空实现钩子）。
- [X] T008 [P] 在 `src/learning/variables/http/handlers.go` 搭建HTTP路由/handler骨架（content与quiz占位），返回JSON基础结构与参数校验。

**Checkpoint**: 基础结构就绪，可开始各用户故事。

---

## Phase 3: User Story 1 - 理解变量与存储 (Priority: P1)  MVP

**Goal**: 提供变量概念与存储方式（声明、new、复合字面量取址）的内容与测验，并可通过CLI/HTTP获取。
**Independent Test**: CLI与HTTP均可单独获取storage主题内容与测验，答题反馈准确。

### Tests for User Story 1
- [X] T009 [P] [US1] 在 `tests/unit/learning/variables/storage_test.go` 添加表驱动单元测试，覆盖存储方式内容与示例数据。
- [X] T010 [P] [US1] 在 `tests/contract/learning/variables/storage_http_test.go` 编写HTTP契约测试，校验 `GET /api/variables/content?topic=storage` 与 quiz 接口的返回结构与校验错误。
- [X] T011 [P] [US1] 在 `tests/integration/learning/variables/storage_cli_test.go` 编写CLI集成测试，验证 list/show/quiz 的存储主题交互与错误输入提示。
- [X] T033 [P] [US1] 在 `tests/unit/learning/variables/structured_elements_test.go` 添加数组/切片/结构体元素寻址与零值的单元测试。

### Implementation for User Story 1
- [X] T012 [US1] 在 `src/learning/variables/storage.go` 编写存储主题内容与示例代码，补充测验项（topic=storage）。
- [X] T013 [US1] 在 `src/learning/variables/cli/menu.go` 实现存储主题的 show 与 quiz 逻辑，保证错误输入提示与返回上级菜单。
- [X] T014 [US1] 在 `src/learning/variables/http/handlers.go` 实现 storage 主题的 content/quiz 处理（含参数校验与错误响应）。
- [X] T034 [US1] 在 `src/learning/variables/storage.go` 补充结构化元素寻址与零值示例（含取地址与独立赋值）。
- [X] T035 [P] [US1] 在 `src/learning/variables/storage_example_test.go` 添加 Example 函数覆盖存储与结构化元素示例。

**Checkpoint**: US1 可独立运行（CLI/HTTP），相关测试通过。

---

## Phase 4: User Story 2 - 静态类型与动态类型 (Priority: P2)

**Goal**: 讲解静态类型来源、接口动态类型变化，提供示例与测验，CLI/HTTP 均可获取。
**Independent Test**: 独立获取 static/dynamic 主题内容与测验并通过。

### Tests for User Story 2
- [X] T015 [P] [US2] 在 `tests/unit/learning/variables/types_test.go` 添加单元测试覆盖静态/动态类型内容与示例。
- [X] T016 [P] [US2] 在 `tests/contract/learning/variables/types_http_test.go` 编写HTTP契约测试，覆盖 `topic=static` 与 `topic=dynamic` 的 content/quiz 响应。
- [X] T017 [P] [US2] 在 `tests/integration/learning/variables/types_cli_test.go` 编写CLI集成测试，验证 static/dynamic 主题展示与测验。

### Implementation for User Story 2
- [X] T018 [US2] 在 `src/learning/variables/static_type.go` 实现静态类型与可赋值性示例与测验项（topic=static）。
- [X] T019 [US2] 在 `src/learning/variables/dynamic_type.go` 实现接口动态类型示例与测验项（topic=dynamic），涵盖 nil 与指针案例。
- [X] T020 [US2] 在 `src/learning/variables/cli/menu.go` 扩展 static/dynamic 主题的 show/quiz 流程与菜单。
- [X] T021 [US2] 在 `src/learning/variables/http/handlers.go` 扩展 static/dynamic 主题的 content/quiz 处理与校验。
- [X] T036 [P] [US2] 在 `src/learning/variables/static_type_example_test.go` 添加 Example 展示静态类型与可赋值性。
- [X] T037 [P] [US2] 在 `src/learning/variables/dynamic_type_example_test.go` 添加 Example 展示接口动态类型与 nil/指针案例。

**Checkpoint**: US2 可独立运行并通过测试。

---

## Phase 5: User Story 3 - 零值与取值规则 (Priority: P3)

**Goal**: 讲解零值与未赋值取值规则，覆盖接口、指针、结构化元素；提供测验，CLI/HTTP 均可获取。
**Independent Test**: 独立获取 zero 主题内容与测验并通过。

### Tests for User Story 3
- [X] T022 [P] [US3] 在 `tests/unit/learning/variables/zero_value_test.go` 添加零值与取值规则的单元测试。
- [X] T023 [P] [US3] 在 `tests/contract/learning/variables/zero_http_test.go` 编写HTTP契约测试，覆盖 `topic=zero` 的 content/quiz 响应。
- [X] T024 [P] [US3] 在 `tests/integration/learning/variables/zero_cli_test.go` 编写CLI集成测试，验证 zero 主题展示与测验。

### Implementation for User Story 3
- [X] T025 [US3] 在 `src/learning/variables/zero_value.go` 实现零值与取值规则的示例与测验项（topic=zero）。
- [X] T026 [US3] 在 `src/learning/variables/cli/menu.go` 扩展 zero 主题 show/quiz 流程，保证错误输入提示与返回。
- [X] T027 [US3] 在 `src/learning/variables/http/handlers.go` 扩展 zero 主题 content/quiz 处理与校验。
- [X] T038 [P] [US3] 在 `src/learning/variables/zero_value_example_test.go` 添加 Example 展示零值与取值规则。

**Checkpoint**: US3 可独立运行并通过测试。

---

## Phase 6: Polish & Cross-Cutting

- [X] T028 [P] 全量执行 gofmt/go vet/golint，修复发现的问题。
- [X] T029 [P] 统计覆盖率（`go test ./... -cover`），确保>=80%，补充缺口测试。
- [X] T030 [P] 校验 CLI/HTTP 文案与注释全中文，更新 `specs/006-learn-variables/quickstart.md` 如有变更。
- [X] T031 更新根级 README 与 `src/learning/variables/README.md`，同步新增章节与运行方式（含双模式入口）。
- [ ] T032 依 quickstart 跑通 CLI 与 HTTP 路径，记录验证结果。
- [ ] T039 [P] 运行 `go mod tidy` 并确认无多余依赖，符合质量要求。

---

## Dependencies & Execution Order

- Setup → Foundational → User Stories → Polish。
- Foundational 完成后 US1/US2/US3 可并行；推荐按优先级 US1 → US2 → US3 形成 MVP 再迭代。
- 各故事内部顺序：测试用例先行（确保初始失败）→ 内容/示例实现 → CLI/HTTP 扩展 → 通过测试。
- Polish 依赖所有目标故事完成。

### Parallel Opportunities
- 标记 [P] 任务可在不同文件并行：基础骨架(T006-T008)、各故事测试/实现同主题不同端口、全量静态检查/覆盖率。
- 不同用户故事在 Foundational 后可由不同成员并行推进。

### Implementation Strategy (MVP)
1) 完成 Phase1-2，确保双模式骨架与数据结构可用。
2) 先交付 US1（存储主题）形成 MVP，并验证 CLI/HTTP + 测验闭环。
3) 按优先级补充 US2、US3，保持各故事独立可测。
4) 最后执行 Polish，整理文档、覆盖率与语言规范检查。