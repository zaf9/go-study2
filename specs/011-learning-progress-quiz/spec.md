# Feature Specification: Go-Study2 学习闭环与测验体系

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `011-learning-progress-quiz`  
**Created**: 2025-12-11  
**Status**: Draft  
**Input**: User description: "注意需求的前缀编号是011。# Go-Study2 产品需求文档 (PRD)"

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

## Clarifications

### Session 2025-12-11

- Q: 阅读进度上报失败的重试策略与卸载时的最终同步如何设计以兼顾数据保留与避免风暴？ → A: 采用指数退避重试（含抖动，最多 5 次），并在页面卸载前强制再同步一次。

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - 统一章节学习与恢复 (Priority: P1)

学习者进入任一主题章节，看到统一结构的内容（概述、要点、详细说明、代码示例、陷阱、实践建议），页面自动恢复上次阅读位置并显示阅读进度。

**Why this priority**: 内容统一和位置恢复是学习闭环的基础，直接影响学习效率与留存。

**Independent Test**: 仅上线单个章节的统一结构与阅读恢复即可独立验证学习体验改进。

**Acceptance Scenarios**:

1. **Given** 用户首次进入章节，**When** 页面加载，**Then** 展示完整标准模块且阅读进度为 0%，记录首次访问时间。
2. **Given** 用户曾阅读到页面中部，**When** 再次进入，**Then** 自动滚动到上次位置并提示已恢复，进度与时长累加。

---

### User Story 2 - 学习进度总览与继续学习 (Priority: P1)

学习者在进度页查看整体与各主题完成度、章节状态，点击“继续学习”快速跳转到第一个未完成章节。

**Why this priority**: 可视化进度与快捷入口驱动持续学习，直接支撑业务目标的学习闭环。

**Independent Test**: 单独发布进度页即可衡量是否正确显示进度与跳转逻辑。

**Acceptance Scenarios**:

1. **Given** 用户已完成部分章节，**When** 打开 `/progress`，**Then** 显示整体进度、各主题进度条、状态图标、完成计数。
2. **Given** 用户有未完成章节，**When** 点击“继续学习”，**Then** 跳转到首个未完成章节详情页并可继续阅读。

---

### User Story 3 - 章节测验与结果反馈 (Priority: P1)

学习者完成章节后发起测验，按题型答题、提交并查看分数与解析，同时更新章节状态。

**Why this priority**: 测验是验证学习效果与闭环完成度的关键环节。

**Independent Test**: 仅上线测验页与提交接口即可验证题目抽取、判分、状态更新是否正确。

**Acceptance Scenarios**:

1. **Given** 用户点击章节底部“开始测验”，**When** 抽题显示多题型并可前后跳题，**Then** 记录作答并显示题目进度与用时。
2. **Given** 用户提交作答，**When** 判分完成，**Then** 展示得分/正确率/用时、标记是否通过，提供解析与重新测验入口，并将章节状态更新为 tested 或 completed（通过时）。

---

### Edge Cases

- 用户快速进入后立即离开：需确保首次访问与最小进度同步不丢失且状态仍为 not_started/in_progress 的正确判定。
- 阅读过程网络波动：防抖/累加上报失败后重试，不应导致进度回退或重复累加。
- 长页未滚动到底但阅读时长已满足：必须同时满足滚动阈值与时长阈值才能标记 completed。
- 测验中途刷新或关闭：应保留未提交的作答缓存，重新进入可继续；已提交后禁止重复提交同一会话结果。
- 多端/多窗口同时学习同章节：以时间戳最新的进度为准，避免状态倒退。

## Assumptions & Dependencies

- 依赖现有用户体系与鉴权，学习进度与测验记录均以已登录用户为前提。
- 章节清单与预估阅读时长已在内容侧准备并可提供标准化标识（topic、chapter）供前后端共享。
- variables 与 types 章节最终数量按内容团队定稿范围（约 8-10 / 15-20），进度计算需随实际章节数配置。
- 前端可获取页面内容总高度/滚动位置并稳定上报，后端支持幂等更新与时间戳对比。

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: 4 个主题（lexical_elements、constants、variables、types）全部章节按统一结构呈现：概述、知识要点、详细说明、≥2 个可运行代码示例（含标题/行内中文注释/预期输出）、常见陷阱（可选）、实践建议。
- **FR-002**: 每章节代码示例可单独运行且输出与文档描述一致，示例长度 15-40 行，遵循官方 Go 代码风格。
- **FR-003**: 章节阅读完成判定需同时满足：阅读时长≥预估的 80%、滚动进度≥90%、章节测验通过（正确率≥60%）；未满足则为 in_progress。
- **FR-004**: 章节状态标识：not_started（未访问）、in_progress（已访问未完成）、completed（满足完成条件）、tested（已测验，无论通过与否）；通过测验时自动置为 completed 并记录通过分数。
- **FR-005**: 学习进度计算：整体进度 = Σ(主题进度 × 权重)，权重为 lexical_elements 20%、constants 20%、variables 25%、types 35%；主题进度 = 已完成章节 / 该主题总章节 × 100%。
- **FR-006**: 进度数据跟踪字段包含 read_duration（秒，累加）、scroll_progress（0-100，覆盖）、last_position（像素，覆盖）、quiz_score、quiz_passed、状态时间戳（first_visit_at/last_visit_at/completed_at），并保持 user_id+topic+chapter 唯一。
- **FR-007**: 前端进度上报策略：进入章节创建/更新记录；阅读中每 10 秒防抖上报累加时长并覆盖滚动/位置；失败时指数退避+抖动重试（最多 5 次），并在页面卸载前强制再同步一次；测验提交时同步分数与通过标记。
- **FR-008**: `/progress` 页面展示整体进度条、完成章节计数、学习天数与总学习时长；按主题折叠卡显示权重、进度百分比、章节列表及状态图标，支持按状态/主题筛选与章节顺序或最近访问排序；提供“继续学习”跳转至首个未完成章节。
- **FR-009**: `/topics/[topic]` 页面顶部显示当前主题进度条与已完成/总数；章节卡显示状态/分数/章节序号及“查看”“重新测验”；提供“继续学习 [章节名]”快捷入口跳转未完成章节。
- **FR-010**: `/topics/[topic]/[chapter]` 页面显示阅读进度和预计剩余时间提示，支持顶部进度指示与底部导航（上一章/返回列表/下一章/开始测验）；加载时自动滚动到 last_position 并展示恢复提示。
- **FR-011**: 每章节测验题量与题型分配符合章节复杂度要求（基础 5-8 题，中等 8-12 题，复杂 12-15 题），覆盖单选/多选/判断/代码输出/改错，并给出通过标准 60%。
- **FR-012**: 测验流程：开始时抽取题目并记录 started_at；答题可前后跳题与跳过；提交后判分，返回分数、正确率、通过与题目级结果/解析；历史列表提供最近会话记录、筛选与排序。
- **FR-013**: 测验历史页 `/quiz` 列表展示主题/章节、分数、通过状态、用时、完成时间，并提供查看详情与重新测验入口；详情页可筛选仅错题。
- **FR-014**: 数据完整性要求：进度与测验表具备必要索引；删除用户时相关记录级联删除；并行访问时以最新更新时间为准避免状态回退。
- **FR-015**: 易用性与容错：网络异常时前端提示并重试，不得丢失已累计的阅读时长与作答；重复提交需被阻止或幂等处理。

### Key Entities *(include if feature involves data)*

- **LearningProgress**: 记录用户在特定主题与章节的状态、阅读时长、滚动进度、最后位置、测验成绩、时间戳；唯一键 user_id+topic+chapter。
- **QuizQuestion**: 章节题目，含题型、难度、题干、选项、正确答案、解析、可选代码片段；按 topic+chapter 查询与抽题。
- **QuizAttempt**: 用户对单题的作答记录，存储用户答案与是否正确，关联用户与题目。
- **QuizSession**: 一次测验会话的总体成绩、题量、正确数、得分、通过标识、开始/完成时间，关联用户与章节。

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: 4 个主题的全部章节 100% 按标准结构落地且每章≥2 个代码示例通过运行验证。
- **SC-002**: 在测试样本中，整体进度与真实完成度差异≤2 个百分点，章节状态无错判（准确率≥98%）。
- **SC-003**: 阅读位置自动恢复成功率≥95%，恢复提示在 2 秒内展示；离开页面后进度数据丢失率为 0。
- **SC-004**: 章节测验得分计算与正确答案一致率 100%；提交后 3 秒内返回结果并更新章节状态。
- **SC-005**: 进度页与主题页“继续学习”跳转准确率≥98%，用户完成单章测验+完成的平均时间 ≤ 10 分钟（以可用性测试为准）。
