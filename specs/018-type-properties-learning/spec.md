# Feature Specification: Go 类型属性章节学习方案

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `018-type-properties-learning`  
**Created**: 2026-01-15  
**Status**: Draft  
**Input**: User description: "为golang的Properties of types and values章节增加学习内容，参考types的所有实现的功能，为Properties of types and values章节实现同样类似的功能。"

## Constitution Guardrails

- 注释与用户文档需清晰且后端全中文(Principle V/XV)。
- 方案需保持可维护性与单一职责,避免过度设计并保持浅层逻辑(Principle I/IV/VI/XVI)。
- 明确错误处理,无静默失败(Principle II)。
- 规划测试覆盖率≥80%,各包具备 *_test.go 与示例; 前端核心组件同样达标(Principle III/XXI/XXXVI)。
- 目录/职责可预测且遵循标准 Go 布局,仅根目录 main, go.mod/go.sum 完整,各包需 README 说明(Principle VIII/XVIII/XIX)。
- 依赖最小且必要(Principle IX)。
- 安全优先: 输入校验、鉴权、HTTPS、敏感信息保护(Principle VII)。
- 如涉及章节/菜单/主题,需同时支持 CLI 与 HTTP,共享内容源,菜单导航与路由/响应格式一致且显式错误(Principle XXII/XXIII/XXV)。
- Go 规范章节需按章节->子章节->子包层次组织,文件命名与示例齐备(Principle XXIV)。
- 完成后需同步更新 README 等文档(Principle XI)。

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

### User Story 1 - 快速掌握类型属性全貌 (Priority: P1)

学习者希望在一页内掌握 Go 类型属性体系，包括值的表示、底层类型、核心类型、类型标识、可赋值性、可表示性和方法集的定义与规则。

**Why this priority**: 这是整个章节的入口，决定后续深入学习能否理解上下文。类型属性是理解 Go 类型系统高级特性的基础。

**Independent Test**: 仅通过阅读概览内容并完成 3 个基础判断题即可验证是否掌握核心概念。

**Acceptance Scenarios**:

1. **Given** 学习者打开类型属性章节概览，**When** 浏览各子主题摘要，**Then** 10 分钟内能区分底层类型与核心类型的概念。
2. **Given** 学习者阅读值的表示部分，**When** 查看自包含值与引用值的对比，**Then** 能正确判断哪些类型共享底层数据。
3. **Given** 学习者阅读方法集部分，**When** 理解值接收者与指针接收者的区别，**Then** 能解释为何某类型不实现某接口。

---

### User Story 2 - 通过练习验证对类型属性的理解 (Priority: P2)

学习者需要通过针对类型属性规则的练习（含类型标识判定、可赋值性规则、可表示性验证等）来检验理解。

**Why this priority**: 练习能暴露理解误区，直接影响学习效果。类型属性涉及较多边界情况，需要通过测验强化。

**Independent Test**: 完成至少一组包含 5 题的测验，覆盖类型标识、可赋值性、可表示性的判定规则，并获得即时解析。

**Acceptance Scenarios**:

1. **Given** 学习者完成一套包含 5 题的类型属性测验，**When** 提交答案，**Then** 系统显示得分、正确答案与规则依据，且可重做。
2. **Given** 练习包含类型标识判定题，**When** 学习者选择错误答案，**Then** 解析指出两类型为何不同及对应规则来源。
3. **Given** 练习包含可赋值性判定题，**When** 学习者判断错误，**Then** 解析说明可赋值性条件及不满足原因。

---

### User Story 3 - 快速查找类型属性规则与示例 (Priority: P3)

学习者希望在查找特定规则（如底层类型推导规则、核心类型条件、方法集包含规则等）时能秒级定位到说明与示例。

**Why this priority**: 方便回查与对照，减少在长文档中滚动时间。类型属性规则较为复杂，快速检索功能提升学习效率。

**Independent Test**: 输入关键词（如"底层类型"或"方法集"）即可在 15 秒内看到对应规则摘要与正反例。

**Acceptance Scenarios**:

1. **Given** 学习者搜索"类型标识"，**When** 查看搜索结果，**Then** 结果展示类型相同性判定规则与至少一个相同/不同类型的示例对照。
2. **Given** 学习者搜索"可赋值性"，**When** 查看结果，**Then** 能看到可赋值性条件列表及各条件的示例说明。

---

### Edge Cases

- 用户缺乏 Go 基础：需在开头提供术语简表与预备知识提示，避免阅读障碍。
- 章节涉及 Go 1.18+ 泛型与类型参数相关内容：需标明版本适用性，避免与旧版混淆。
- 示例与反例必须覆盖类型参数场景、接口类型集判定、bytestring 核心类型等边界情况。
- 离线或打印需求：需提供可导出/打印的简明提纲版本以便复习。

### Assumptions

- 默认参考 Go 1.20+ 规范，覆盖 1.18 引入的泛型与类型参数相关内容。
- 学习环境可访问示例与测验（在线或本地缓存），并提供中文界面。
- 学习者具备基础编程概念（变量、表达式、函数），以及对 Go 基本类型的初步了解。
- 与 Types 章节共享相同的技术架构和实现模式。
- **MVP 约束**: 本迭代测验题目直接硬编码于 Go 文件中（每子主题 2-3 题），暂不完全遵循 Principle XXXII 的 30-50 题 YAML 题库标准，留待后续统一升级。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 必须提供"值的表示"子主题，涵盖自包含值（数组、结构体等）与引用值（指针、切片、map、channel等）的区别，包含规则、示例和测验。
- **FR-002**: 必须提供"底层类型"子主题，涵盖底层类型的推导规则，特别包括类型别名、类型定义、类型参数的底层类型判定。
- **FR-003**: 必须提供"核心类型"子主题，涵盖非接口类型和接口类型的核心类型判定规则，包括 bytestring 核心类型的概念。
- **FR-004**: 必须提供"类型标识"子主题，涵盖两类型相同/不同的判定规则，包括命名类型、类型字面量、泛型实例化类型的标识判定。
- **FR-005**: 必须提供"可赋值性"子主题，涵盖可赋值性判定的 6 个条件和涉及类型参数的 3 个额外条件。
- **FR-006**: 必须提供"可表示性"子主题，涵盖常量在特定类型范围内的可表示性判定规则。
- **FR-007**: 必须提供"方法集"子主题，涵盖定义类型、指针类型、接口类型的方法集规则，包括嵌入字段对方法集的影响。
- **FR-008**: 必须提供针对各子主题的测验模块，包含评分、解析与重做能力。
- **FR-009**: 必须支持按关键词的快速检索，检索结果需在 15 秒内返回规则摘要与至少一个示例/反例。
- **FR-010**: 必须记录学习进度（已读子主题、测验得分、上次访问位置），遵循 Principle XXXIII 状态转换标准（not_started/in_progress/completed），支持下次继续学习并查看完成度。
- **FR-011**: 必须提供可导出或打印的提纲视图，保留各子主题关键规则与示例的摘要，方便线下复习。
- **FR-012**: 必须同时支持 CLI 和 HTTP 两种访问方式，共享相同的内容源。

### Key Entities *(include if feature involves data)*

- **PropertyConcept**: 描述单个类型属性子主题，包含类别、定义、适用版本、关键规则与注意事项。类似于 Types 模块的 TypeConcept。
- **PropertyRule**: 类型属性相关的规则或约束，包含规则ID、关联概念、规则类型、描述、参考文献、严重程度。
- **ExampleCase**: 关联 PropertyConcept 的示例或反例，包含描述、代码、期望结果、是否合法、违反或满足的规则说明。
- **QuizItem**: 练习题目，包含题干、选项、标准答案、解析、难度标签与关联的 PropertyConcept。
- **LearningProgress**: 学习进度，包含用户标识、已完成的 PropertyConcept、测验得分、最后访问时间与完成百分比。
- **ReferenceIndex**: 检索索引项，包含关键词、关联 PropertyConcept、摘要与跳转位置，用于快速查找。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 80% 的学习者在首次学习会话中于 30 分钟内完成类型属性概览阅读并通过基础测验（得分 >= 80%）。
- **SC-002**: 90% 的测验提交在 15 秒内返回评分与解析，且结果可复用于重做。
- **SC-003**: 90% 的搜索请求在 15 秒内呈现匹配规则摘要和至少一条示例/反例链接。
- **SC-004**: 所有7个子主题均提供完整的概念、规则、示例和测验内容。
- **SC-005**: 单元测试覆盖率达到 80% 以上。
