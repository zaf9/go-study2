# Implementation Plan: 学习章节测验题库扩展

**Branch**: `013-quiz-question-bank` | **Date**: 2025-12-14 | **Spec**: [spec.md](./spec.md)

## Summary

为四个Go学习主题(lexical_elements、constants、variables、types)的41个章节直接生成1230-2050个高质量测验题目（YAML格式），每个章节30-50题，单选/多选各50%，难度分级（简单40%/中等40%/困难20%）。同时实现题库加载、深度验证、随机抽题（全局配置）、无状态抽题算法，确保与现有quiz API兼容。

**核心任务**: 
1. **题目内容生成**（P1）- 为41个章节生成YAML题库文件
2. **题库加载与验证**（P1）- Fail-Fast启动验证，结构化日志
3. **随机抽题算法**（P2）- 无状态方案，难度分布控制
4. **API集成**（P2）- 保持现有接口兼容性

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: 
- gopkg.in/yaml.v3 (YAML解析)
- github.com/gogf/gf/v2 (现有HTTP框架)
- math/rand (随机数生成)

**Storage**: YAML外部文件（quiz_data/{topic}/{chapter}.yaml），启动时加载到内存  
**Testing**: Go testing框架，覆盖率≥80%  
**Target Platform**: Linux/Windows server (与现有backend一致)  
**Project Type**: Web后端扩展 (单体应用)  
**Performance Goals**: 
- 题库查询<100ms
- 随机抽题<50ms  
- 支持100+并发用户
- 启动验证<5秒（41个文件）

**Constraints**: 
- 必须向后兼容现有quiz API
- 题库文件更新需重启服务
- 全局抽题配置（不支持章节级定制）
- 无状态设计（不记录用户历史）

**Scale/Scope**: 
- 41个章节
- 1230-2050个题目（30-50题/章节）
- 4个主题（lexical_elements、constants、variables、types）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- ✅ **Principle I (代码质量与可维护性):** 题库加载/验证/抽题职责分离，每个包单一职责
- ✅ **Principle II (显式错误处理):** Fail-Fast验证，所有错误明确返回，无静默失败
- ✅ **Principle III/XXVI (全面测试):** 规划≥80%覆盖率，包含题库验证、抽题算法、并发测试
- ✅ **Principle IV (单一职责):** 题库加载器、验证器、抽题器独立实现
- ✅ **Principle V/XX (中文文档):** 所有代码注释、错误信息、日志使用中文
- ✅ **Principle VI (YAGNI):** 无状态方案，避免过度设计的历史记录系统
- ✅ **Principle VII/XXXIX (安全优先):** 输入验证（YAML结构）、鉴权继承现有系统
- ✅ **Principle VIII/XXIII (可预测结构):** 遵循标准Go布局，题库包独立
- ✅ **Principle IX (依赖纪律):** 仅新增yaml.v3，其他复用现有依赖
- ✅ **Principle X (性能优化):** 启动时一次性加载，内存缓存，并发安全
- ✅ **Principle XI (文档同步):** 计划更新README，新增题库使用说明
- ✅ **Principle XIII (Fail-Fast启动):** 题库验证失败拒绝启动，详细错误日志
- ✅ **Principle XIV (可观测性):** 结构化日志记录加载耗时、抽题请求、验证错误
- ✅ **Principle XV (配置管理):** YAML配置文件管理抽题数量、难度分布、文件路径
- ✅ **Principle XVI (并发处理):** 题库只读，天然线程安全；随机数使用锁保护
- ✅ **Principle XXI (浅层逻辑):** 避免深层嵌套，使用卫语句和函数拆分
- ✅ **Principle XXIV (包级README):** 为quiz包和quiz_data目录添加README
- ✅ **Principle XXX (HTTP/CLI一致性):** 保持现有API路径和响应格式兼容
- ✅ **Principle XXXII (Quiz标准):** 题型多样化、60%及格线、详细解析、支持重测
- ✅ **Principle XXXIII (进度跟踪):** 支持quiz_score和quiz_passed字段存储
- ✅ **Principle XXXIV (功能独立性):** P1题目生成和加载可独立交付MVP

**验证结果**: ✅ 通过 - 所有原则符合，无阻塞项

## Project Structure

### Documentation (this feature)

```text
specs/013-quiz-question-bank/
├── spec.md              # 功能规范（已完成）
├── plan.md              # 本文件 - 实施计划
├── research.md          # Phase 0 - 题目生成策略研究
├── data-model.md        # Phase 1 - 题库数据模型和YAML格式
├── quickstart.md        # Phase 1 - 题库使用快速指南
├── contracts/           # Phase 1 - API契约定义
│   └── yaml-schema.md   # YAML题库文件格式规范
└── checklists/
    └── requirements.md  # 质量检查清单（已完成）
```

### Source Code (repository root)

```text
backend/
├── internal/
│   └── domain/
│       └── quiz/
│           ├── loader.go          # 【新增】题库加载器
│           ├── loader_test.go     # 【新增】加载器测试
│           ├── validator.go       # 【新增】题库验证器
│           ├── validator_test.go  # 【新增】验证器测试
│           ├── selector.go        # 【新增】抽题选择器
│           ├── selector_test.go   # 【新增】选择器测试
│           ├── entity.go          # 【修改】添加QuizQuestion等实体
│           └── README.md          # 【新增】题库管理包说明
│
├── configs/
│   └── config.yaml            # 【修改】添加quiz配置项
│
└── quiz_data/                 # 【新增】题库YAML文件目录
    ├── README.md              # 【新增】题库文件组织说明
    ├── lexical_elements/      # 【新增】11个章节YAML文件
    │   ├── comments.yaml
    │   ├── tokens.yaml
    │   ├── semicolons.yaml
    │   ├── identifiers.yaml
    │   ├── keywords.yaml
    │   ├── operators.yaml
    │   ├── integers.yaml
    │   ├── floats.yaml
    │   ├── imaginary.yaml
    │   ├── runes.yaml
    │   └── strings.yaml
    ├── constants/             # 【新增】12个章节YAML文件
    │   ├── boolean.yaml
    │   ├── rune.yaml
    │   ├── integer.yaml
    │   ├── floating_point.yaml
    │   ├── complex.yaml
    │   ├── string.yaml
    │   ├── expressions.yaml
    │   ├── typed_untyped.yaml
    │   ├── conversions.yaml
    │   ├── builtin_functions.yaml
    │   ├── iota.yaml
    │   └── implementation_restrictions.yaml
    ├── variables/             # 【新增】4个章节YAML文件
    │   ├── storage.yaml
    │   ├── static.yaml
    │   ├── dynamic.yaml
    │   └── zero.yaml
    └── types/                 # 【新增】14个章节YAML文件
        ├── boolean.yaml
        ├── numeric.yaml
        ├── string.yaml
        ├── array.yaml
        ├── slice.yaml
        ├── struct.yaml
        ├── pointer.yaml
        ├── function.yaml
        ├── interface_basic.yaml
        ├── interface_embedded.yaml
        ├── interface_general.yaml
        ├── interface_impl.yaml
        ├── map.yaml
        └── channel.yaml

tests/
└── integration/
    └── quiz_bank_test.go      # 【新增】题库集成测试
```

**Structure Decision**: 

采用单体应用扩展结构。题库加载/验证/抽题逻辑作为新的domain包`internal/domain/quiz`。题库YAML文件独立存放在`backend/quiz_data/`按主题分目录组织，便于管理和版本控制。

## Phase 0: Research & Outline

**Goal**: 研究题目生成策略、明确数据模型、解决技术未知项

### 0.1 题目生成策略研究

**Unknowns to resolve**:
- 如何高效为41个章节生成1230-2050个高质量题目？
- 题目生成是否可以借助AI辅助？格式和质量如何保证？
- 如何确保题目覆盖章节核心知识点？
- 题目难度如何判定？有无参考标准？

**Research tasks**:
1. 调研Go语言教学最佳实践和现有题库模式
2. 分析现有chapters内容，提取核心知识点清单
3. 设计题目生成workflow（可能包含AI辅助+人工审核）
4. 定义题目质量标准（题干清晰度、选项合理性、解析完整性）

**Output**: `research.md` - 包含：
- 题目生成workflow（步骤、工具、质量控制）
- 41个章节知识点提纲（每章节5-8个核心点）
- 题目质量标准checklist
- 难度判定指南

### 0.2 YAML格式和数据模型设计

**Unknowns to resolve**:
- YAML文件的精确结构（字段、嵌套、枚举值）？
- 如何在YAML中表达单选/多选题的差异？
- 题目ID命名规则？
- 难度和题型的枚举值定义？

**Research tasks**:
1. 设计QuizQuestion数据结构（Go struct和YAML映射）
2. 定义YAML schema规范，包含示例
3. 设计题目ID生成规则（如：lexical-comments-001）
4. 明确枚举值：difficulty (easy/medium/hard), type (single/multiple)

**Output**: `data-model.md` - 包含：
- QuizQuestion Go struct定义
- YAML文件格式规范和完整示例
- 题目ID命名规则
- 枚举值定义表

### 0.3 技术依赖和性能基准

**Unknowns to resolve**:
- yaml.v3解析性能是否满足<5秒启动要求？
- 随机数生成是否需要加密安全级别？
- 并发抽题的锁粒度如何设计？

**Research tasks**:
1. 验证yaml.v3解析41个文件（约2000题）的性能
2. 设计并发安全的随机数生成方案
3. 确定内存占用估算（2000题 × 平均大小）

**Output**: 更新`research.md` - 添加：
- 性能基准测试结果
- 并发方案设计（锁策略）
- 内存占用评估

## Phase 1: Design & Contracts

**Prerequisites**: research.md完成

### 1.1 数据模型最终确认

**Task**: 基于research.md，确定最终数据模型

**Output**: `data-model.md` - 包含：

```yaml
# 示例YAML结构
questions:
  - id: lexical-comments-001
    type: single              # single/multiple
    difficulty: easy          # easy/medium/hard
    stem: "Go语言中，以下哪种注释方式是正确的？"
    options:
      - "A: // 这是单行注释"
      - "B: # 这是注释"
      - "C: <!-- 这是注释 -->"
      - "D: -- 这是注释"
    answer: "A"
    explanation: "Go语言支持两种注释方式：单行注释使用//，多行注释使用/* */。选项B是Python风格，C是HTML风格，D是SQL风格，均不适用于Go。"
    topic: "lexical_elements"
    chapter: "comments"

  - id: lexical-comments-002
    type: multiple
    difficulty: medium
    stem: "关于Go语言的注释，以下说法正确的是？（多选）"
    options:
      - "A: 单行注释以//开头"
      - "B: 多行注释可以嵌套"
      - "C: 注释不影响程序执行"
      - "D: 注释可以用于文档生成"
    answer: "ACD"
    explanation: "A正确，单行注释使用//。B错误，Go的多行注释/* */不支持嵌套。C正确，注释在编译时被忽略。D正确，godoc工具可以提取注释生成文档。"
    topic: "lexical_elements"
    chapter: "comments"
```

### 1.2 YAML Schema规范

**Output**: `contracts/yaml-schema.md` - 包含：
- 必填字段定义
- 字段类型约束
- 枚举值列表
- 验证规则（如：single题answer必须为单个字母，multiple题answer必须为2-4个字母）

### 1.3 题库配置规范

**Output**: 更新`backend/configs/config.yaml`示例：

```yaml
quiz:
  dataPath: "backend/quiz_data"           # 题库文件根目录
  questionCount:
    single: 4                              # 单选题数量
    multiple: 4                            # 多选题数量
  difficultyDistribution:
    easy: 40                               # 简单题占比%
    medium: 40                             # 中等题占比%
    hard: 20                               # 困难题占比%
  loadTimeout: 5s                          # 题库加载超时
  validation:
    strictMode: true                       # 严格验证模式
    failFast: true                         # 首个错误即停止
```

### 1.4 API契约确认

**Output**: `contracts/api-spec.md` - 确认现有API保持兼容：

```http
# 现有API（保持不变）
GET /api/v1/quiz/:topic/:chapter/start
Response: {
  "code": 0,
  "message": "success",
  "data": {
    "questions": [...]  # 随机抽取的题目列表
  }
}

# 新增API（可选）
GET /api/v1/quiz/:topic/:chapter/stats
Response: {
  "code": 0,
  "message": "success",
  "data": {
    "total": 35,
    "byType": {"single": 18, "multiple": 17},
    "byDifficulty": {"easy": 14, "medium": 14, "hard": 7}
  }
}
```

### 1.5 快速开始指南

**Output**: `quickstart.md` - 包含：
- 题库文件组织说明
- 如何添加/修改题目
- 如何触发重新加载（重启服务）
- 验证错误排查指南

### 1.6 更新Agent上下文

**Task**: 运行`.specify/scripts/powershell/update-agent-context.ps1 -AgentType copilot`

**目的**: 将新增的技术栈（yaml.v3）和概念（题库、抽题算法）添加到agent上下文，便于后续开发。

## Phase 2: Task Breakdown

**Note**: 此阶段由`/speckit.tasks`命令生成，不在`/speckit.plan`范围内。

**预期任务分解方向**:
1. **Task Group 1**: 题目内容生成（P1）
   - 为41个章节生成YAML题库文件
   - 人工审核和质量调整
   
2. **Task Group 2**: 题库加载与验证（P1）
   - 实现loader.go（YAML文件加载）
   - 实现validator.go（深度验证）
   - 单元测试（coverage≥80%）

3. **Task Group 3**: 随机抽题算法（P2）
   - 实现selector.go（按难度抽题）
   - 并发安全保证
   - 单元测试

4. **Task Group 4**: API集成（P2）
   - 集成到现有quiz API
   - 配置文件支持
   - 集成测试

## Risk Analysis

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| 题目生成质量不达标 | High | Medium | 建立质量checklist，多轮人工审核，抽样测试 |
| 题库文件过大影响启动 | Medium | Low | 性能基准测试，必要时引入懒加载 |
| YAML解析错误难以定位 | Medium | Medium | Fail-Fast详细错误日志，包含文件名、行号 |
| 41个章节题目生成耗时长 | High | High | 可借助AI辅助生成初稿，人工精修 |
| 抽题算法难度分布不准 | Low | Low | 单元测试覆盖各种边界情况 |

## Success Metrics

- ✅ 41个章节全部拥有30-50个高质量题目
- ✅ 单选/多选比例符合50%±5%
- ✅ 难度分布符合40/40/20±10%
- ✅ 题库验证错误率<1%（启动时检测）
- ✅ 单元测试覆盖率≥80%
- ✅ 抽题响应时间<50ms (p95)
- ✅ 启动加载时间<5秒
- ✅ 100并发用户无性能退化

## Next Steps

1. ✅ **Completed**: Phase 0 Outline (本plan.md)
2. ⏭ **Next**: 执行`/speckit.phase0`生成research.md（题目生成策略）
3. ⏭ **Then**: 执行`/speckit.phase1`生成data-model.md + contracts/
4. ⏭ **Finally**: 执行`/speckit.tasks`分解具体任务

---

**Plan Status**: ✅ Complete  
**Ready for**: Phase 0 Research
