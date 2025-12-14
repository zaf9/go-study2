# Tasks: 学习章节测验题库扩展

**Input**: Design documents from `/specs/013-quiz-question-bank/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/yaml-schema.md, quickstart.md

**Tests**: 本功能需达到≥80%单元测试覆盖率。测试任务为**必需项**。

**Organization**: 任务按用户故事分组，每个故事可独立实现和测试。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 所属用户故事（如US1, US2, US3）
- 描述中包含精确文件路径

## Constitution Guardrails

- 所有注释与用户文档必须使用中文(Principle V/XX)。
- 规划测试覆盖率≥80%,各包包含*_test.go与示例(Principle III/XXVI)。
- 目录/文件保持单一职责,遵循标准Go布局,补齐包README(Principle IV/VIII/XXIII/XXIV)。
- 外部依赖最小化,错误处理显式,避免YAGNI(Principle II/VI/IX)。
- 完成后更新README等文档(Principle XI)。

---

## Phase 1: Setup (项目初始化)

**Purpose**: 创建题库目录结构和配置基础

- [ ] T001 创建题库文件根目录结构：`backend/quiz_data/{lexical_elements,constants,variables,types}/`
- [ ] T002 创建题库包目录：`backend/internal/domain/quiz/`及README.md
- [ ] T003 [P] 更新`backend/configs/config.yaml`添加quiz配置项（dataPath, questionCount, difficultyDistribution等）
- [ ] T004 [P] 添加yaml.v3依赖到`backend/go.mod`：`go get gopkg.in/yaml.v3`

---

## Phase 2: Foundational (核心基础设施)

**Purpose**: 实现题库加载、验证、抽题的基础设施，所有用户故事依赖此阶段完成

**⚠️ CRITICAL**: 此阶段完成前，用户故事无法开始实现

- [ ] T005 定义核心实体：在`backend/internal/domain/quiz/entity.go`实现QuizQuestion, QuizBank, QuizConfig结构体（参考data-model.md）
- [ ] T006 [P] 实现题库加载器：`backend/internal/domain/quiz/loader.go`（YAML文件解析,按topic/chapter索引）
- [ ] T007 [P] 实现题库验证器：`backend/internal/domain/quiz/validator.go`（必填字段、枚举值、选项格式、答案格式、路径一致性、ID唯一性验证）
- [ ] T008 [P] 实现抽题选择器：`backend/internal/domain/quiz/selector.go`（难度分布控制、并发安全的随机数生成、Fisher-Yates洗牌算法）
- [ ] T009 实现题库仓储：`backend/internal/domain/quiz/repository.go`（内存索引map[topic][chapter][]QuizQuestion）
- [ ] T010 [P] 单元测试-loader：`backend/internal/domain/quiz/loader_test.go`（测试YAML解析、错误处理）
- [ ] T011 [P] 单元测试-validator：`backend/internal/domain/quiz/validator_test.go`（测试9个验证规则、边界情况）
- [ ] T012 [P] 单元测试-selector：`backend/internal/domain/quiz/selector_test.go`（测试难度分布、并发安全、随机性）

**Checkpoint**: 基础设施就绪 - 用户故事实现可以开始

---

## Phase 3: User Story 1 - 章节题目内容生成 (Priority: P1) 🎯 MVP

**Goal**: 为41个章节直接生成30-50个高质量测验题目（YAML格式），涵盖不同难度和题型

**Independent Test**: 检查生成的YAML文件（如`quiz_data/constants/boolean.yaml`），验证包含30-50题，单选/多选各约50%，每题包含完整字段

### Tests for User Story 1 (MANDATORY) ⚠️

- [ ] T013 [P] [US1] 题库文件格式验证测试：`backend/tests/unit/quiz/yaml_format_test.go`（验证41个YAML文件符合schema规范）
- [ ] T014 [P] [US1] 题目数量分布测试：`backend/tests/integration/quiz/question_distribution_test.go`（验证每章节30-50题，单选/多选各50%±10%）
- [ ] T015 [P] [US1] 难度分布测试：`backend/tests/integration/quiz/difficulty_distribution_test.go`（验证难度比例40/40/20±10%）

### Implementation for User Story 1

**注意**: 题目内容生成主要是内容创作任务，不是代码开发任务

- [ ] T016 [US1] 为lexical_elements主题生成11个章节YAML文件（comments.yaml到strings.yaml），每章节30-50题（参考research.md知识点提纲）
  - T016.1: [P] [US1] 生成`quiz_data/lexical_elements/comments.yaml`（35题）
  - T016.2: [P] [US1] 生成`quiz_data/lexical_elements/tokens.yaml`（40题）
  - T016.3: [P] [US1] 生成`quiz_data/lexical_elements/semicolons.yaml`（30题）
  - T016.4: [P] [US1] 生成`quiz_data/lexical_elements/identifiers.yaml`（40题）
  - T016.5: [P] [US1] 生成`quiz_data/lexical_elements/keywords.yaml`（30题）
  - T016.6: [P] [US1] 生成`quiz_data/lexical_elements/operators.yaml`（45题）
  - T016.7: [P] [US1] 生成`quiz_data/lexical_elements/integers.yaml`（40题）
  - T016.8: [P] [US1] 生成`quiz_data/lexical_elements/floats.yaml`（35题）
  - T016.9: [P] [US1] 生成`quiz_data/lexical_elements/imaginary.yaml`（30题）
  - T016.10: [P] [US1] 生成`quiz_data/lexical_elements/runes.yaml`（40题）
  - T016.11: [P] [US1] 生成`quiz_data/lexical_elements/strings.yaml`（45题）

- [ ] T017 [US1] 为constants主题生成12个章节YAML文件（boolean.yaml到implementation_restrictions.yaml），每章节30-50题
  - T017.1: [P] [US1] 生成`quiz_data/constants/boolean.yaml`（30题）
  - T017.2: [P] [US1] 生成`quiz_data/constants/rune.yaml`（35题）
  - T017.3: [P] [US1] 生成`quiz_data/constants/integer.yaml`（40题）
  - T017.4: [P] [US1] 生成`quiz_data/constants/floating_point.yaml`（40题）
  - T017.5: [P] [US1] 生成`quiz_data/constants/complex.yaml`（35题）
  - T017.6: [P] [US1] 生成`quiz_data/constants/string.yaml`（35题）
  - T017.7: [P] [US1] 生成`quiz_data/constants/expressions.yaml`（45题）
  - T017.8: [P] [US1] 生成`quiz_data/constants/typed_untyped.yaml`（40题）
  - T017.9: [P] [US1] 生成`quiz_data/constants/conversions.yaml`（40题）
  - T017.10: [P] [US1] 生成`quiz_data/constants/builtin_functions.yaml`（35题）
  - T017.11: [P] [US1] 生成`quiz_data/constants/iota.yaml`（45题）
  - T017.12: [P] [US1] 生成`quiz_data/constants/implementation_restrictions.yaml`（30题）

- [ ] T018 [US1] 为variables主题生成4个章节YAML文件（storage.yaml到zero.yaml），每章节30-50题
  - T018.1: [P] [US1] 生成`quiz_data/variables/storage.yaml`（40题）
  - T018.2: [P] [US1] 生成`quiz_data/variables/static.yaml`（35题）
  - T018.3: [P] [US1] 生成`quiz_data/variables/dynamic.yaml`（40题）
  - T018.4: [P] [US1] 生成`quiz_data/variables/zero.yaml`（45题）

- [ ] T019 [US1] 为types主题生成14个章节YAML文件（boolean.yaml到channel.yaml），每章节30-50题
  - T019.1: [P] [US1] 生成`quiz_data/types/boolean.yaml`（30题）
  - T019.2: [P] [US1] 生成`quiz_data/types/numeric.yaml`（50题）
  - T019.3: [P] [US1] 生成`quiz_data/types/string.yaml`（40题）
  - T019.4: [P] [US1] 生成`quiz_data/types/array.yaml`（45题）
  - T019.5: [P] [US1] 生成`quiz_data/types/slice.yaml`（50题）
  - T019.6: [P] [US1] 生成`quiz_data/types/struct.yaml`（50题）
  - T019.7: [P] [US1] 生成`quiz_data/types/pointer.yaml`（40题）
  - T019.8: [P] [US1] 生成`quiz_data/types/function.yaml`（45题）
  - T019.9: [P] [US1] 生成`quiz_data/types/interface_basic.yaml`（45题）
  - T019.10: [P] [US1] 生成`quiz_data/types/interface_embedded.yaml`（40题）
  - T019.11: [P] [US1] 生成`quiz_data/types/interface_general.yaml`（45题）
  - T019.12: [P] [US1] 生成`quiz_data/types/interface_impl.yaml`（40题）
  - T019.13: [P] [US1] 生成`quiz_data/types/map.yaml`（45题）
  - T019.14: [P] [US1] 生成`quiz_data/types/channel.yaml`（40题）

- [ ] T020 [US1] 人工审核所有41个YAML文件，修正明显错误（题干语法、答案准确性、解析完整性）
- [ ] T021 [US1] 创建题库README：`backend/quiz_data/README.md`（说明文件组织、添加题目流程）

**Checkpoint**: 此时User Story 1应完全可用且可独立测试（启动服务会加载所有题库并验证）

---

## Phase 4: User Story 2 - 智能随机抽题 (Priority: P1)

**Goal**: 每次测验从题库随机抽取3-5道单选题和3-5道多选题，难度分布合理

**Independent Test**: 多次调用`/api/v1/quiz/constants/boolean/start`，验证每次返回题目组合不同且符合数量和难度分布要求

### Tests for User Story 2 (MANDATORY) ⚠️

- [ ] T022 [P] [US2] 抽题API契约测试：`backend/tests/contract/quiz/quiz_start_contract_test.go`（验证GET `/api/v1/quiz/:topic/:chapter`返回正确结构）
- [ ] T023 [P] [US2] 随机性集成测试：`backend/tests/integration/quiz/quiz_randomness_test.go`（100次抽题，验证至少50%题目不同）
- [ ] T024 [P] [US2] 难度分布集成测试：`backend/tests/integration/quiz/quiz_difficulty_test.go`（1000次抽题，验证难度比例40/40/20±5%）
- [ ] T025 [P] [US2] 并发抽题测试：`backend/tests/integration/quiz/quiz_concurrency_test.go`（100并发用户同时抽题，无panic，响应<100ms）

### Implementation for User Story 2

- [ ] T026 [US2] 实现抽题服务层：`backend/internal/app/quiz/service.go`（调用selector抽题，返回QuizSelection）
- [ ] T027 [US2] 实现HTTP handler：`backend/internal/app/http_server/handler/quiz.go`（GET `/quiz/:topic/:chapter`路由，兼容现有API）
- [ ] T027-A [US2] 实现选项顺序随机打乱：在`selector.go`或`service.go`中实现Fisher-Yates洗牌算法打乱题目选项顺序（FR-010：避免答案位置规律）
- [ ] T028 [US2] 注册路由：更新`backend/internal/app/http_server/router.go`，添加quiz路由
- [ ] T029 [US2] 集成启动加载：在`backend/main.go`启动时调用quiz.LoadAllBanks()，Fail-Fast验证
- [ ] T030 [US2] 添加结构化日志：在loader、validator、selector中记录加载耗时、验证错误、抽题请求（使用现有logger）
- [ ] T031 [P] [US2] 单元测试-service：`backend/internal/app/quiz/service_test.go`（测试抽题逻辑、错误处理）
- [ ] T032 [P] [US2] 单元测试-handler：`backend/internal/app/http_server/handler/quiz_test.go`（测试HTTP响应格式、错误码）

**Checkpoint**: 此时User Story 1和2应都能独立工作（题库已生成+抽题API可用）

---

## Phase 5: User Story 3 - 题目质量保证 (Priority: P2)

**Goal**: 确保所有题目题干清晰、选项合理、答案正确、解析详细

**Independent Test**: 人工审核或自动化测试验证题目格式完整性、答案正确性和解析合理性

### Tests for User Story 3 (MANDATORY) ⚠️

- [ ] T033 [P] [US3] 题目结构完整性测试：`backend/tests/unit/quiz/question_completeness_test.go`（验证所有题目包含9个必填字段）
- [ ] T034 [P] [US3] 答案格式正确性测试：`backend/tests/unit/quiz/answer_format_test.go`（单选题1字母，多选题2-4字母升序）
- [ ] T034-A [P] [US3] 多选题判分顺序无关测试：`backend/tests/unit/quiz/answer_order_test.go`（验证FR-009：AB与BA、ACD与DCA等视为相同答案）
- [ ] T035 [P] [US3] 解析中文检测测试：`backend/tests/unit/quiz/explanation_chinese_test.go`（验证所有解析包含中文字符）

### Implementation for User Story 3

- [ ] T036 [US3] 实现题目质量审核工具（可选）：`backend/scripts/quiz_quality_check.go`（批量检查题目质量、生成报告）
- [ ] T037 [US3] 第二轮人工审核：复查所有41个YAML文件，重点验证答案准确性和解析详细度
- [ ] T038 [US3] 完善验证器错误信息：在`validator.go`中添加更详细的中文错误提示（包含文件名、题目ID、错误行号）

**Checkpoint**: 所有用户故事应独立可用且题目质量达标

---

## Phase 6: User Story 4 - 题库数据管理 (Priority: P3)

**Goal**: 题库有清晰组织结构，易于扩展，可查询统计信息

**Independent Test**: 添加新题目、修改现有题目、查询题库统计信息

### Tests for User Story 4 (MANDATORY) ⚠️

- [ ] T039 [P] [US4] 题库统计API测试：`backend/tests/contract/quiz/quiz_stats_contract_test.go`（验证GET `/api/v1/quiz/:topic/:chapter/stats`返回正确统计）
- [ ] T040 [P] [US4] 题库重载测试：`backend/tests/integration/quiz/quiz_reload_test.go`（修改YAML文件后重启，验证变更生效）

### Implementation for User Story 4

- [ ] T041 [P] [US4] 实现统计API（可选）：在`handler/quiz.go`添加GET `/quiz/:topic/:chapter/stats`（返回total, byType, byDifficulty）
- [ ] T042 [P] [US4] 实现题库统计服务：在`service.go`添加GetStats方法（从repository聚合统计）
- [ ] T043 [US4] 创建题库维护文档：更新`specs/013-quiz-question-bank/quickstart.md`，补充常见问题FAQ

**Checkpoint**: 所有用户故事功能完整且可维护

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: 跨用户故事的改进和文档更新

- [ ] T044 [P] 验证所有代码注释和日志为中文（检查entity.go, loader.go, validator.go, selector.go, service.go, handler.go）
- [ ] T045 [P] 更新项目README：在`backend/README.md`添加测验题库系统说明
- [ ] T046 [P] 更新API文档：在`docs/API.md`添加quiz相关接口说明
- [ ] T047 代码重构：优化validator.go的验证逻辑，减少重复代码
- [ ] T048 性能基准测试：在`backend/tests/benchmark/quiz_benchmark_test.go`添加加载和抽题性能测试
- [ ] T049 运行quickstart验证：按照`quickstart.md`流程添加一个测试题目，验证完整流程

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖 - 立即开始
- **Foundational (Phase 2)**: 依赖Setup完成 - **阻塞所有用户故事**
- **User Stories (Phase 3-6)**: 都依赖Foundational完成
  - 用户故事可并行进行（如有人力）
  - 或按优先级顺序（P1 → P2 → P3）
- **Polish (Phase 7)**: 依赖所需用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: Foundational完成后可开始 - 无其他故事依赖
- **User Story 2 (P1)**: Foundational完成后可开始 - **依赖User Story 1的题库文件**（T016-T019必须完成）
- **User Story 3 (P2)**: Foundational完成后可开始 - 依赖User Story 1的题库文件
- **User Story 4 (P3)**: Foundational完成后可开始 - 可独立测试但统计API依赖题库存在

### Within Each User Story

- User Story 1: T016-T019可并行生成（41个YAML文件独立），T020人工审核需等待生成完成
- User Story 2: Tests (T022-T025)在实现(T026-T032)之前编写并失败
- 模型→服务→接口的顺序
- 核心实现→集成
- 故事完成后再进入下一优先级

### Parallel Opportunities

- **Setup阶段**: T001-T004全部可并行
- **Foundational阶段**: T006-T008可并行（不同文件），T010-T012可并行（测试文件）
- **User Story 1**: T016的11个子任务、T017的12个子任务、T018的4个子任务、T019的14个子任务全部可并行生成
- **User Story 2**: T022-T025测试可并行，T031-T032测试可并行
- **不同User Story**: 一旦Foundational完成，US1题目生成、US2代码实现、US3质量审核可由不同团队成员并行

---

## Parallel Example: User Story 1 (题目生成)

```bash
# 启动所有lexical_elements章节题目生成（11个文件并行）:
Task: "生成quiz_data/lexical_elements/comments.yaml (35题)"
Task: "生成quiz_data/lexical_elements/tokens.yaml (40题)"
Task: "生成quiz_data/lexical_elements/semicolons.yaml (30题)"
... (共11个)

# 启动所有constants章节题目生成（12个文件并行）:
Task: "生成quiz_data/constants/boolean.yaml (30题)"
Task: "生成quiz_data/constants/rune.yaml (35题)"
... (共12个)

# 启动所有variables章节题目生成（4个文件并行）:
Task: "生成quiz_data/variables/storage.yaml (40题)"
Task: "生成quiz_data/variables/static.yaml (35题)"
Task: "生成quiz_data/variables/dynamic.yaml (40题)"
Task: "生成quiz_data/variables/zero.yaml (45题)"

# 启动所有types章节题目生成（14个文件并行）:
Task: "生成quiz_data/types/boolean.yaml (30题)"
Task: "生成quiz_data/types/numeric.yaml (50题)"
... (共14个)

# 总计: 41个YAML文件可完全并行生成
```

---

## Implementation Strategy

### MVP First (仅User Story 1 + 2)

1. 完成Phase 1: Setup
2. 完成Phase 2: Foundational（关键-阻塞所有故事）
3. 完成Phase 3: User Story 1（生成题库文件）
4. 完成Phase 4: User Story 2（实现抽题API）
5. **STOP并验证**: 测试题库加载和抽题功能
6. 如就绪则部署/演示

### Incremental Delivery

1. 完成Setup + Foundational → 基础就绪
2. 添加User Story 1 → 独立测试 → 部署/演示（题库文件可用！）
3. 添加User Story 2 → 独立测试 → 部署/演示（MVP完整！）
4. 添加User Story 3 → 独立测试 → 部署/演示（质量提升）
5. 添加User Story 4 → 独立测试 → 部署/演示（可维护性增强）
6. 每个故事独立增加价值，不破坏已有故事

### Parallel Team Strategy

多开发者场景：

1. 团队共同完成Setup + Foundational
2. Foundational完成后：
   - **内容团队A**: User Story 1（题目生成，可2-3人分主题并行）
   - **开发者B**: User Story 2（抽题API实现）
   - **QA C**: User Story 3（质量审核工具）
3. 故事独立完成并集成

---

## Notes

- [P]任务 = 不同文件，无依赖关系
- [Story]标签映射任务到特定用户故事便于追踪
- 每个用户故事应可独立完成和测试
- 实现前先验证测试失败
- 每个任务或逻辑组完成后提交
- 在任何检查点停止以独立验证故事
- 避免: 模糊任务、同文件冲突、破坏独立性的跨故事依赖
- **题目生成是内容创作任务**: T016-T019的子任务可使用AI辅助+人工审核（参考research.md的Prompt模板）

---

**Tasks Status**: ✅ Complete  
**Total Tasks**: 49 (主任务) + 41 (题目生成子任务) = 90个任务  
**MVP Scope**: Phase 1-2 + User Story 1-2 (T001-T032) = 32个任务  
**Parallel Opportunities**: 41个题库文件生成可完全并行，基础设施测试可并行，多用户故事可并行开发
