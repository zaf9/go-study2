# Tasks: Go 类型章节学习方案

**Input**: 方案文档位于 `D:\studyspace\go-study\go-study2\specs\007-go-types-study\`  
**Prerequisites**: plan.md（必需）、spec.md（必需）、research.md、data-model.md、contracts/、quickstart.md  
**Tests**: 必须确保单元测试覆盖率 ≥80%，各故事需配套契约/集成测试。

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 准备章节目录与基础占位，便于后续并行开发。

 - [X] T001 创建 Types 章节目录与 README 占位：`src/learning/types/README.md`，包含章节结构与版本说明。  
 - [X] T002 创建 CLI/HTTP 目录与空实现文件占位：`src/learning/types/cli/menu.go`、`src/learning/types/http/handlers.go`，保证编译通过并留待后续填充。

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 提供通用数据结构、路由入口和测试基座，解锁各用户故事。

- [X] T003 定义共享域模型与内容注册器：`src/learning/types/types.go`（TypeConcept/TypeRule/ExampleCase/QuizItem/ReferenceIndex/LearningProgress 及加载接口）。  
- [X] T004 注册 Types 入口：更新 `main.go` 菜单编号与 `internal/app/http_server/router.go` 新增 `/api/v1/topic/types` 路由。  
- [X] T005 创建 HTTP handler 骨架：`internal/app/http_server/handler/types_content.go`（菜单/内容/搜索/测验占位返回，中文错误提示）。  
- [X] T006 搭建测试目录与公共测试数据：在 `tests/unit/learning/types/`、`tests/contract/learning/types/`、`tests/integration/learning/types/` 创建基座与 `testdata` 夹，预置示例/反例/题库样例。
- [X] T039 更新主题列表与主菜单：`internal/app/http_server/handler/topics.go` 增加 types 入口，`main.go` 主菜单添加 types 编号，确保与路由标识一致。  
- [X] T040 [P] 主题列表契约测试：`tests/contract/learning/types/types_topics_contract_test.go` 校验 `/api/v1/topics` JSON/HTML 含 types 入口。

---

## Phase 3: User Story 1 - 快速掌握类型全貌 (Priority: P1) 🎯 MVP

**Goal**: 提供类型概览与基础/复合类型内容，含 3 个基础判断题，支持 CLI/HTTP 展示。  
**Independent Test**: 阅读概览后完成 3 个基础判断题得分展示；CLI/HTTP 均可返回内容与示例。

### Tests for User Story 1 (MANDATORY)

- [X] T007 [P] [US1] 编写内容覆盖与结构校验单测：`tests/unit/learning/types/content_test.go`。  
- [X] T008 [P] [US1] 编写菜单/内容 JSON 与 HTML 契约测试：`tests/contract/learning/types/types_menu_contract_test.go`。  
- [X] T009 [P] [US1] 编写 CLI 菜单与展示集成测试：`tests/integration/learning/types/types_cli_menu_test.go`。

### Implementation for User Story 1

 - [X] T010 [US1] 实现概览聚合：`src/learning/types/overview.go`（术语简表、适用版本、打印提纲）。  
 - [X] T011 [P] [US1] 实现基础类型内容：`src/learning/types/boolean.go`、`src/learning/types/numeric.go`、`src/learning/types/string_type.go`（含正反例与规则编号）。  
 - [X] T012 [P] [US1] 实现复合类型内容：`src/learning/types/array.go`、`src/learning/types/slice.go`、`src/learning/types/struct_type.go`、`src/learning/types/pointer.go`、`src/learning/types/function_type.go`、`src/learning/types/map_type.go`、`src/learning/types/channel_type.go`。  
- [X] T013 [US1] 实现 CLI 菜单渲染与 3 题基础测验：`src/learning/types/cli/menu.go`（编号、`q` 返回、得分显示）。  
- [X] T014 [US1] 实现 HTTP 菜单与内容返回（JSON/HTML）：`internal/app/http_server/handler/types_content.go`，复用内容注册器输出示例与基础测验。
 - [X] T030 [P] [US1] 实现接口类型内容：`src/learning/types/interface_basic.go`、`interface_embedded.go`、`interface_general.go`、`interface_impl.go`（规则编号、正反例与适用版本）。  
 - [X] T031 [US1] 扩充 HTTP 接口子主题输出：`internal/app/http_server/handler/types_content.go` 覆盖 interface_* 子主题 JSON/HTML。  
 - [X] T032 [US1] 扩充 CLI 菜单与接口子主题展示：`src/learning/types/cli/menu.go` 增加 interface_* 选项与展示逻辑。  
- [X] T033 [P] [US1] 接口子主题契约/集成测试：`tests/contract/learning/types/types_menu_contract_test.go`、`tests/integration/learning/types/types_cli_menu_test.go` 增加 interface_* 覆盖。  
- [X] T034 [US1] 实现打印/导出提纲：`src/learning/types/overview.go` 与 `internal/app/http_server/handler/types_content.go` 支持可打印文本/HTML 与 CLI 打印命令。  
- [X] T035 [P] [US1] 提纲导出契约/集成测试：`tests/contract/learning/types/types_outline_contract_test.go`、`tests/integration/learning/types/types_cli_outline_test.go`。
- [X] T041 [US1] 完善未知子主题错误处理：`internal/app/http_server/handler/types_content.go` 确保未知 subtopic 返回 JSON/HTML 404/业务错误且结构与链接一致。  
- [X] T042 [P] [US1] 未知子主题错误契约/集成测试：`tests/contract/learning/types/types_menu_contract_test.go` 与 `tests/integration/learning/types/types_cli_menu_test.go` 增加未知 subtopic 404/错误输出校验。

---

## Phase 4: User Story 2 - 通过练习验证理解 (Priority: P2)

**Goal**: 提供身份/可比较性/接口类型集判定测验，含评分、解析与重做；CLI/HTTP 均支持提交。  
**Independent Test**: 提交 5 题测验后返回得分、正确答案与解析；CLI 可重做，HTTP POST 评分返回 200。

### Tests for User Story 2 (MANDATORY)

- [X] T015 [P] [US2] 编写测验评分与重做逻辑单测：`tests/unit/learning/types/quiz_test.go`。  
- [X] T016 [P] [US2] 编写测验提交契约测试：`tests/contract/learning/types/types_quiz_contract_test.go`（POST `/api/v1/topic/types/quiz/submit`）。  
- [X] T017 [P] [US2] 编写 CLI 测验评分与重做集成测试：`tests/integration/learning/types/types_cli_quiz_test.go`。

### Implementation for User Story 2

- [X] T018 [US2] 实现测验数据与评分逻辑：`src/learning/types/quiz.go`（题目/答案/解析/规则引用，重做支持）。  
- [X] T019 [US2] 扩展 CLI 测验流程：`src/learning/types/cli/menu.go`（收集作答、评分、解析、重试提示）。  
- [X] T020 [US2] 实现 HTTP 测验提交与错误处理：`internal/app/http_server/handler/types_content.go`（校验参数、返回得分与解析，中文错误信息）。

---

## Phase 5: User Story 3 - 快速查找规则与反例 (Priority: P3)

**Goal**: 提供关键词检索与边界清单（非法递归、不可比较、接口自包含等），返回规则摘要与正反例。  
**Independent Test**: 输入关键词（如 `map key`、`~int`）15 秒内返回摘要与正反例锚点；CLI/HTTP 均可检索。

### Tests for User Story 3 (MANDATORY)

- [X] T021 [P] [US3] 编写检索索引覆盖与边界校验单测：`tests/unit/learning/types/search_index_test.go`。  
- [X] T022 [P] [US3] 编写检索接口契约测试：`tests/contract/learning/types/types_search_contract_test.go`（GET `/api/v1/topic/types/search`）。  
- [X] T023 [P] [US3] 编写 CLI 搜索集成测试：`tests/integration/learning/types/types_cli_search_test.go`。

### Implementation for User Story 3

- [X] T024 [US3] 构建检索索引与边界清单：`src/learning/types/search.go`（关键词→摘要/正反例/锚点）。  
- [X] T025 [US3] 实现 HTTP 检索响应：`internal/app/http_server/handler/types_content.go`（JSON/HTML 双格式，404/400 友好提示）。  
- [X] T026 [US3] 实现 CLI 搜索命令：`src/learning/types/cli/menu.go`（`search <keyword>`，返回规则摘要与正反例锚点）。
- [X] T036 [P] [US3] 增补高风险边界索引：`src/learning/types/search.go` 覆盖接口自包含、不可比较键类型、递归限制、~T 规则等关键词锚点。  
- [X] T037 [P] [US3] 边界关键词测试覆盖：`tests/unit/learning/types/search_index_test.go`、`tests/contract/learning/types/types_search_contract_test.go` 增加接口自包含、map 键不可比较、递归数组/struct 等用例。

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T027 完善学习进度记录与持久化（内存/文件轻量实现）：`src/learning/types/types.go` 并在 CLI/HTTP 调用处写入进度。  
- [ ] T028 同步文档：更新 `specs/007-go-types-study/quickstart.md`、`README.md`、`docs/quickstart-variables.md` 新增 Types 入口与示例。  
- [ ] T029 代码检查与测试收尾：`gofmt`、`go test ./...`，补充缺失中文注释与错误处理。
- [ ] T038 质量工具链检查：执行 `go vet ./...`、`golint ./...`、`go mod tidy`，并在 quickstart/README 中补充命令说明。

---

## Dependencies & Execution Order

- Phase 1 → Phase 2 → US1(P1) → US2(P2) → US3(P3) → Polish。  
- US2 依赖 US1 的内容与菜单已就绪；US3 依赖 US1 的内容与 US2 的测验数据（用于索引解析）。  
- 标记 [P] 的任务可并行（不同文件且无前置依赖）。

## Parallel Execution Examples

- US1 并行：T007/T008/T009（测试）与 T011/T012（内容文件）可并行，合流到 T013/T014。  
- US2 并行：T015/T016/T017（测试）与 T018（数据）并行，合流到 T019/T020。  
- US3 并行：T021/T022/T023（测试）与 T024（索引数据）并行，合流到 T025/T026。

## Implementation Strategy

- MVP：完成 US1（Phase 3）后可演示概览+基础测验。  
- 增量：US2 增加评分/解析；US3 增加检索与边界清单；最终 Polish 收尾与文档同步。

