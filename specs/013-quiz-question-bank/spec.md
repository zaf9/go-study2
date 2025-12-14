# Feature Specification: 学习章节测验题库扩展

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `013-quiz-question-bank`  
**Created**: 2025年12月14日  
**Status**: Draft  
**Input**: User description: "当前已经实现的golang学习的四个主题lexical_elements、constants、variables和types各个章节的测验题目偏少或没有。为了加强学习效果,需要增加每个主题的每个章节候选测试题目库和测验题目。需要为每个主题的每个章节根据难易程度,实现30-50个测验题目,单选题和多选题各占一半。各个章节测验的时候,根据难易程度,从测验题目库中随机获取3-5个单选题,3-5各多选题进行测验。"

## Constitution Guardrails

- 注释与用户文档需清晰且后端全中文(Principle V/XX)。
- 方案需保持可维护性与单一职责,避免过度设计并保持浅层逻辑(Principle I/IV/VI/XXI)。
- 明确错误处理,无静默失败(Principle II)。
- 规划测试覆盖率≥80%,各包具备 *_test.go 与示例; 前端核心组件同样达标(Principle III/XXVI/XLV)。
- 目录/职责可预测且遵循标准 Go 布局,仅根目录 main, go.mod/go.sum 完整,各包需 README 说明(Principle VIII/XXIII/XXIV)。
- 依赖最小且必要(Principle IX)。
- 安全优先: 输入校验、鉴权、HTTPS、敏感信息保护(Principle VII/XXXIX)。
- 如涉及章节/菜单/主题,需同时支持 CLI 与 HTTP,共享内容源,菜单导航与路由/响应格式一致且显式错误(Principle XXVIII/XXIX/XXX)。
- 测验题目需符合Quiz and Assessment Standards: 题型多样化、及格标准60%、详细解析、支持重测(Principle XXXII)。
- 学习进度跟踪需满足Progress Tracking要求: 章节状态管理、进度指标、完成标准、自动保存策略(Principle XXXIII)。
- 功能设计遵循独立可测试原则,每个User Story可独立交付MVP(Principle XXXIV)。
- 完成后需同步更新 README 等文档(Principle XI)。

## Clarifications

### Session 2025-12-14

- Q: 题库数据的存储组织方式？ → A: 外部文件存储，按章节组织（如 quiz_data/constants/boolean.yaml），启动时加载到内存
- Q: 抽题数量的配置方式？ → A: 全局可配置，在配置文件中设置一个全局值，所有章节使用相同数量
- Q: 题目重复率控制方式？ → A: 无状态方案，不记录历史抽题记录，通过随机算法自然避免重复
- Q: 题库数据验证策略？ → A: 深度验证+完全失败，启动时验证所有题目完整性，任一错误则拒绝启动并输出详细错误信息
- Q: 题库文件格式？ → A: YAML格式（每个章节一个YAML文件，如 constants/boolean.yaml），符合项目配置管理标准。**重要补充**：题库需要通过生成机制创建30-50个测试题目（非完全人工编写），人工仅做后续维护和调整
- Q: 题目生成的范围和方式？ → A: **本次需求的核心任务是直接为每个章节生成30-50个测试题目内容（YAML格式），不需要开发生成工具或模板系统，直接产出符合规范的题目数据文件即可**

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

### User Story 1 - 章节题目内容生成 (Priority: P1)

作为题库开发者,我需要为每个学习章节直接生成30-50个高质量测验题目（YAML格式），题目涵盖不同难度级别和题型，以便学习者能够通过多样化的测验巩固知识点。

**Why this priority**: 这是核心价值，没有题目内容就无法实现测验功能，这是所有后续功能的基础。

**Independent Test**: 可以通过检查生成的YAML文件（如 `quiz_data/constants/boolean.yaml`），验证包含30-50个题目，单选题和多选题各占约50%，每题包含完整字段。

**Acceptance Scenarios**:

1. **Given** lexical_elements主题的任一章节, **When** 查询该章节的题库, **Then** 返回30-50个测验题目,包含题目ID、题干、选项、答案、解析和难度级别
2. **Given** constants主题的任一章节, **When** 查询该章节的题库, **Then** 返回30-50个测验题目,其中单选题约占50%,多选题约占50%
3. **Given** variables主题的任一章节, **When** 查询该章节的题库, **Then** 返回30-50个测验题目,难度分为简单、中等、困难三个级别
4. **Given** types主题的任一章节, **When** 查询该章节的题库, **Then** 返回30-50个测验题目,每题包含详细的答案解析

---

### User Story 2 - 智能随机抽题 (Priority: P1)

作为Go语言学习者,我希望每次测验时系统能从题库中随机抽取3-5道单选题和3-5道多选题,且根据难度级别合理分配,以便每次测验都有新鲜感并覆盖不同知识点。

**Why this priority**: 随机抽题是题库扩充后的核心使用场景,直接决定学习体验和知识覆盖面。

**Independent Test**: 可以通过多次调用开始测验接口(如 `/api/v1/quiz/constants/boolean/start`),验证每次返回的题目组合不同,且符合数量和难度分布要求。

**Acceptance Scenarios**:

1. **Given** 用户开始某章节测验, **When** 系统从题库抽题, **Then** 返回3-5道单选题和3-5道多选题,共6-10题
2. **Given** 用户开始某章节测验, **When** 系统抽取题目, **Then** 题目难度分布为:简单题40%、中等题40%、困难题20%
3. **Given** 用户连续两次开始同一章节测验, **When** 对比两次题目, **Then** 题目组合不完全相同(至少50%不同)
4. **Given** 题库中某难度级别题目不足, **When** 系统抽题, **Then** 自动调整其他难度题目占比以满足总题数要求

---

### User Story 3 - 题目质量保证 (Priority: P2)

作为Go语言学习者,我希望所有题目都经过精心设计,题干清晰准确,选项合理,答案正确,解析详细,以便我能通过测验真正理解知识点而非死记硬背。

**Why this priority**: 题目质量直接影响学习效果,但可以在题库建设过程中逐步完善,优先级略低于数量和抽题机制。

**Independent Test**: 可以通过人工审核或自动化测试验证题目格式完整性、答案正确性和解析合理性。

**Acceptance Scenarios**:

1. **Given** 任一测验题目, **When** 检查题目结构, **Then** 必须包含:题干、至少2个选项(单选题2-4个,多选题3-5个)、正确答案、中文解析
2. **Given** 任一单选题, **When** 检查答案, **Then** 答案为A/B/C/D中的一个,且对应的选项能正确解答题干
3. **Given** 任一多选题, **When** 检查答案, **Then** 答案为多个字母组合(如AB、ACD),且所有正确选项都被包含
4. **Given** 任一题目的解析, **When** 阅读解析, **Then** 解析用中文清晰说明为何该答案正确,其他选项为何错误,并补充相关知识点

---

### User Story 4 - 题库数据管理 (Priority: P3)

作为系统维护者,我希望题库数据有清晰的组织结构和易于扩展的设计,以便后续能方便地添加、修改或删除题目。

**Why this priority**: 这是长期可维护性的保障,但不影响用户的直接学习体验,可以后期优化。

**Independent Test**: 可以通过添加新题目、修改现有题目或查询题库统计信息来验证数据管理功能。

**Acceptance Scenarios**:

1. **Given** 需要为某章节新增题目, **When** 调用题库管理接口, **Then** 能成功添加题目并自动分配唯一ID
2. **Given** 发现某题目答案有误, **When** 修正题目数据, **Then** 修改立即生效且不影响正在进行的测验
3. **Given** 需要查看题库统计, **When** 调用统计接口, **Then** 返回各主题、各章节、各难度级别的题目数量分布
4. **Given** 题库数据存储在代码或数据库中, **When** 系统重启, **Then** 所有题目数据保持不变

---

### Edge Cases

- 当某章节题库题目少于6题时,系统如何处理抽题请求?
- 当用户在短时间内多次开始同一章节测验时,如何避免题目重复率过高?
- 当题库中某难度级别题目为0时,抽题算法如何调整?
- 多选题的答案顺序是否影响判分(如AB与BA应视为相同答案)?
- 如何处理题库数据损坏或格式错误的情况?


## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 必须为lexical_elements、constants、variables、types四个主题的每个章节直接生成30-50个测验题目，存储在按章节组织的YAML文件中（如 `quiz_data/constants/boolean.yaml`）
- **FR-002**: 生成的测验题目必须分为单选题和多选题两种类型,且每种类型约占50%
- **FR-003**: 生成的每个测验题目必须分配难度级别(简单、中等、困难),建议分布为简单40%、中等40%、困难20%
- **FR-004**: 生成的每个题目必须包含:唯一ID、题干、选项列表、正确答案、中文解析、难度级别、题型
- **FR-004-A**: 题目内容必须直接生成为符合YAML格式规范的文件，无需开发独立的生成工具或模板系统
- **FR-004-B**: 系统启动时必须深度验证所有YAML题库文件的格式和完整性，包括必填字段、答案有效性、选项数量等，任一验证失败则拒绝启动并输出详细错误信息
- **FR-005**: 系统必须在用户开始章节测验时,从题库中随机抽取题目，抽题数量通过全局配置文件设置（默认8题：4单选+4多选）
- **FR-006**: 系统必须在抽题时根据难度级别进行合理分配,优先选择简单和中等难度题目
- **FR-007**: 系统必须使用随机算法确保题目组合的多样性，无需记录用户历史抽题记录（无状态方案）
- **FR-008**: 系统必须在题库题目不足配置要求的最小题数时,返回明确错误信息而非静默失败
- **FR-009**: 系统必须支持多选题答案的顺序无关判分(AB与BA视为相同答案)
- **FR-010**: 系统必须在返回测验题目时,打乱选项顺序以避免答案位置规律
- **FR-011**: 系统必须保持与现有quiz API的兼容性(`/api/v1/quiz/:topic/:chapter/start`)
- **FR-012**: 系统必须为每个主题的每个章节提供题库查询接口,返回题目总数和难度分布统计

### Key Entities

- **QuizQuestion(测验题目)**: 代表题库中的单个题目,包含题干、选项、答案、解析、难度级别、题型(单选/多选)、所属主题、所属章节
- **QuestionBank(题库)**: 特定章节的所有题目集合,支持按难度、题型筛选和随机抽取
- **QuizSession(测验会话)**: 用户一次测验的题目组合,包含抽取的题目列表、开始时间、用户答案(提交后)
- **DifficultyLevel(难度级别)**: 枚举类型,包含Easy(简单)、Medium(中等)、Hard(困难)三个级别
- **QuestionType(题型)**: 枚举类型,包含SingleChoice(单选)、MultipleChoice(多选)两种类型

### Non-Functional Requirements

- **NFR-001**: 系统必须在1秒内完成题库查询和随机抽题操作(Principle X - Performance Optimization)
- **NFR-002**: 题库数据加载失败时必须拒绝启动服务并输出清晰错误信息(Principle XIII - Fail-Fast Startup)
- **NFR-003**: 所有题库操作必须记录结构化日志,包括抽题请求、题库查询、错误情况(Principle XIV - Observability)
- **NFR-004**: 测验相关配置(题目数量范围、难度分布比例)必须可通过YAML配置文件调整(Principle XV - Configuration Management)
- **NFR-005**: 并发测验请求必须保证线程安全,支持至少100个并发用户同时抽题(Principle XVI - Concurrent Processing)
- **NFR-006**: 代码注释和错误信息必须使用中文(Principle XX - Backend Chinese Documentation)
- **NFR-007**: 题库代码必须保持浅层逻辑,避免深度嵌套(Principle XXI - Shallow Logic)
- **NFR-008**: 题库查询和抽题功能必须达到80%以上单元测试覆盖率(Principle XXVI - Backend Testing)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 每个学习章节的题库包含30-50个高质量测验题目,单选题和多选题各占约50%
- **SC-002**: 用户每次开始测验时,能在1秒内获得6-10道随机抽取的题目
- **SC-003**: 同一章节连续两次测验,题目重复率低于50%
- **SC-004**: 所有题目都包含完整的中文解析,帮助学习者理解知识点
- **SC-005**: 题库覆盖四个主题(lexical_elements、constants、variables、types)的所有现有章节(约41个章节)
- **SC-006**: 测验抽题的难度分布符合预期:简单题约40%、中等题约40%、困难题约20%
- **SC-007**: 多选题答案判分支持顺序无关(AB与BA视为正确)
- **SC-008**: 题库管理达到80%以上的单元测试覆盖率

## Out of Scope

本功能明确**不包括**以下内容:

- 题目生成工具或模板系统的开发（直接生成题目内容即可，无需工具抽象层）
- 题目难度的自动化评估算法（难度在生成时人工判定）
- 题库管理的Web界面或管理后台（题库通过YAML文件管理）
- 题目的动态热更新功能（更新需重启服务）
- 基于数据库的题库存储方案（采用YAML外部文件）
- 题目推荐算法或个性化推荐（随机抽取即可）
- 用户历史抽题记录的持久化存储（采用无状态方案）
- 答题时间限制或计时功能（不属于本期范围）
- 题目收藏或错题本功能（可作为后续扩展）
- 题目标签或分类体系（难度级别已足够）
- 题目生成的自动化框架或可复用组件（本次仅生成当前41个章节的题目内容）

## Dependencies and Constraints

### Dependencies

- **现有测验框架**: 依赖于lexical_elements、constants、variables、types四个主题已实现的基础测验API结构
- **章节内容理解**: 依赖对四个主题各章节知识点的深入理解，以生成准确、有价值的测验题目
- **YAML解析库**: 依赖Go YAML解析库（如gopkg.in/yaml.v3）来加载和解析题库文件
- **随机数生成**: 依赖Go标准库`math/rand`进行随机抽题
- **HTTP路由框架**: 依赖现有的GoFrame HTTP服务框架
- **日志系统**: 依赖项目现有的结构化日志基础设施(Principle XIV)
- **配置系统**: 依赖现有的YAML配置管理机制(Principle XV)

### Constraints

- **外部文件约束**: 题库数据存储在YAML文件中（按章节组织，如quiz_data/constants/boolean.yaml），更新题库需修改YAML文件并重启服务
- **全局配置约束**: 抽题数量采用全局配置，所有章节使用相同抽题规则，无法针对单个章节定制
- **无状态约束**: 系统不记录用户历史抽题记录，无法基于历史实现精确的去重控制
- **兼容性约束**: 必须保持与现有quiz API路径和响应格式的向后兼容(`/api/v1/quiz/:topic/:chapter/start`)
- **性能约束**: 题库加载和深度验证发生在服务启动时，大量题目可能影响启动时间（需在启动日志中记录加载耗时）
- **并发安全约束**: 题库数据为只读，天然支持并发访问，但抽题随机数生成需考虑线程安全
- **测试覆盖约束**: 必须达到80%以上单元测试覆盖率(Principle III/XXVI)

## Assumptions

- 假设现有的四个主题(lexical_elements、constants、variables、types)已经实现了基本的测验框架,只需扩充题库数据
- 假设题库数据存储在YAML外部文件中（按章节组织），启动时加载到内存，无需引入数据库
- 假设题目内容将直接生成为符合YAML格式的文件（每个章节30-50题），无需开发独立的生成工具
- 假设题目生成过程可以借助AI辅助、参考现有题目模式等方式，但最终产出是直接可用的YAML题库文件
- 假设题目难度级别在生成时由生成者判定，生成后可由人工审核调整
- 假设测验题目的选项数量为:单选题2-4个选项,多选题3-5个选项
- 假设用户在测验过程中不会修改浏览器本地存储或发送恶意请求
- 假设题目ID在同一章节内唯一即可,无需全局唯一
- 假设抽题算法使用Go标准库的随机数生成器,通过时间种子等机制自然避免重复，无需记录用户历史
- 假设题库YAML文件更新后需要重启服务才能生效
- 假设现有的学习进度跟踪系统(Principle XXXIII)已支持quiz_score和quiz_passed字段存储
- 假设YAML文件路径配置在主配置文件中，支持开发/测试/生产环境使用不同路径
