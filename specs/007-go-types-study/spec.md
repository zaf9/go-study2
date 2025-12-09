# Feature Specification: Go 类型章节学习方案

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: 007-go-types-study  
**Created**: 2025-12-09  
**Status**: Draft  
**Input**: User description: "实现golanng的types的章节学习需求，遵循constitution中的规则要求。"

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

### User Story 1 - 快速掌握类型全貌 (Priority: P1)

学习者希望在一页内掌握 Go 类型体系（基础、复合、接口、泛型类型参数等）的定义、取值范围和常见用法。

**Why this priority**: 这是整个章节的入口，决定后续练习能否理解上下文。

**Independent Test**: 仅通过阅读概览内容并完成 3 个基础判断题即可验证是否掌握核心概念。

**Acceptance Scenarios**:

1. **Given** 学习者打开类型章节概览，**When** 浏览基础类型和命名/别名说明，**Then** 10 分钟内能回答布尔、数值、字符串的定义与区别。
2. **Given** 学习者阅读复合类型部分，**When** 查看数组、切片、结构体示例，**Then** 能指出长度与容量的差异且不混淆多维组合的限制。

---

### User Story 2 - 通过练习验证理解 (Priority: P2)

学习者需要基于类型规则的练习（含接口类型集、方法集、非法递归类型判定等）来检验理解。

**Why this priority**: 练习能暴露理解误区，直接影响学习效果。

**Independent Test**: 完成至少一组针对类型身份、可比较性、接口实现判定的测验，并获得即时解析。

**Acceptance Scenarios**:

1. **Given** 学习者完成一套包含 5 题的类型测验，**When** 提交答案，**Then** 系统显示得分、正确答案与规则依据，且可重做。
2. **Given** 练习包含接口类型集判定题，**When** 学习者选择错误类型，**Then** 解析指出该类型不在接口类型集中以及对应规则来源。

---

### User Story 3 - 快速查找规则与反例 (Priority: P3)

学习者希望在查找特定规则（如禁止的递归组合、可作 map 键的类型、通道方向性等）时能秒级定位到说明与示例。

**Why this priority**: 方便回查与对照，减少在长文档中滚动时间。

**Independent Test**: 输入关键词（如map 键类型或~int 接口）即可在 15 秒内看到对应规则摘要与一个正反例。

**Acceptance Scenarios**:

1. **Given** 学习者搜索递归数组类型限制，**When** 查看搜索结果，**Then** 结果展示限制条款与至少一个非法/合法示例的对照。
2. **Given** 学习者搜索接口类型集 union，**When** 查看结果，**Then** 能看到并理解交并集规则以及不可用的组合限制。

---

### Edge Cases

- 用户缺乏 Go 基础：需在开头提供术语简表与预备知识提示，避免阅读障碍。
- 章节涉及 Go 1.18+ 泛型与接口类型集：需标明版本适用性，避免与旧版混淆。
- 示例与反例必须覆盖不可比较类型、非法递归类型、通道方向性等易错点，避免遗漏关键边界。
- 离线或打印需求：需提供可导出/打印的简明提纲版本以便复习。

### Assumptions

- 默认参考 Go 1.20+ 规范，至少覆盖 1.18 引入的泛型与接口类型集相关内容。
- 学习环境可访问示例与测验（在线或本地缓存），并提供中文界面。
- 学习者具备基础编程概念（变量、表达式），无需已有 Go 实战经验。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 必须提供覆盖基础类型、数值范围、别名与命名类型区别的概览，学习者单独阅读即可获取定义与取值范围。
- **FR-002**: 必须为每类复合类型（数组、切片、结构体、指针、函数、接口、map、channel）提供规则摘要、典型用例、常见陷阱与正反示例。
- **FR-003**: 必须提供针对类型身份、可比较性、方法集与接口类型集判定的测验模块，包含评分、解析与重做能力。
- **FR-004**: 必须支持按关键词或类型名称的快速检索，检索结果需在 15 秒内返回规则摘要与至少一个示例/反例。
- **FR-005**: 必须记录学习进度（已读章节、测验得分、上次访问位置），支持下次继续学习并查看完成度。
- **FR-006**: 必须提供可导出或打印的提纲视图，保留各类型关键规则与示例的摘要，方便线下复习。
- **FR-007**: 必须明确列出不合规的递归定义、不可作为 map 键的类型、接口自包含等边界清单，便于对照校验。

### Key Entities *(include if feature involves data)*

- **TypeConcept**: 描述单个类型主题，包含类别（基础/复合/接口）、定义、适用版本、关键规则与注意事项。
- **ExampleCase**: 关联 TypeConcept 的示例或反例，包含描述、期望结果、违反或满足的规则说明。
- **QuizItem**: 练习题目，包含题干、选项或判断、标准答案、解析、难度标签与关联的 TypeConcept。
- **LearningProgress**: 学习进度，包含用户标识、已完成的 TypeConcept、测验得分、最后访问时间与完成百分比。
- **ReferenceIndex**: 检索索引项，包含关键词、关联 TypeConcept、摘要与跳转位置，用于快速查找。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 80% 的学习者在首次学习会话中于 30 分钟内完成类型概览阅读并通过基础测验（得分 >= 80%）。
- **SC-002**: 90% 的测验提交在 15 秒内返回评分与解析，且结果可复用于重做。
- **SC-003**: 90% 的搜索请求在 15 秒内呈现匹配规则摘要和至少一条示例/反例链接。
- **SC-004**: 课程满意度调查中，关于类型章节清晰度/实用性的平均评分 >= 4.3/5，且回访率提升 >= 20%。
