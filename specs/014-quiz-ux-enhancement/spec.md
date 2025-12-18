# Feature Specification: 章节测验体验升级与功能深化 (Chapter Quiz UX & Depth Enhancement)

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `014-quiz-ux-enhancement`  
**Created**: 2025-12-18  
**Status**: Draft  
**Input**: User description: "@[/speckit.specify] 注意spec需求的编号前缀使用：014..."

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

### Session 2025-12-18

- Q: CLI/Web Scope for History Review? → A: **Web-only**. The rich history review and visual tracking will be implemented for the Web UI only. CLI retains its existing "take-quiz-only" behavior without persistent history or review capabilities, as a specific scope exception to Principle XXII for this UX-focused feature.
- Q: Data Persistence Strategy? → A: **Backend Database (SQLite)**. Quiz sessions and granular answer details will be stored in the SQLite database using `gdb`. This ensures robust querying capabilities for the "Review Mode" and supports future analytics.
- Q: Save Timing? → A: **Save on Completion**. Full results are committed to the database only when the user explicitly confirms submission. Intermediate states are managed in the frontend (React state), simplifying the initial implementation while meeting core requirements.
- Q: Analysis Granularity? → A: **Global/General Analysis**. The data model will support a single `analysis` field per question, shown to the user during review regardless of which specific option they selected. Per-option feedback is out of scope for this version.

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

### User Story 1 - 测验交互优化：防误触与有序标签 (Priority: P1)

用户在于测验过程中，需要稳定的视觉锚点（有序标签）以及安全的提交机制，以减少认知负荷和误操作风险。

**Why this priority**: 直接解决用户当前的挫败感痛点（误触提交、乱序干扰），是体验升级的基础。

**Independent Test**: 该功能可独立于后端存储测试。前端组件应能独立渲染固定标签并拦截提交事件。

**Acceptance Scenarios**:

1. **Given** 题目选项内容随机排列, **When** 用户查看题目选项, **Then** 选项编号必须始终显示为有序的 A, B, C, D...。
2. **Given** 用户正在答题, **When** 用户点击“提交测验”按钮, **Then** 系统弹出二次确认对话框，显示“已答题数”和“未答题数”。
3. **Given** 二次确认对话框显示, **When** 用户点击“取消”, **Then** 对话框关闭，停留在答题页面，不做提交处理。
4. **Given** 二次确认对话框显示, **When** 用户点击“确认”, **Then** 系统执行提交逻辑并跳转结果页。

---

### User Story 2 - 结果页反馈增强：百分制与题型标识 (Priority: P1)

用户在提交测验后，需要直观、清晰的评价反馈，包括百分制得分、及格状态颜色区分以及题型的明确标识。

**Why this priority**: 优化反馈机制，帮助用户准确评估学习效果，符合常规心理预期。

**Independent Test**: 可以通过 Mock 提交结果数据来测试结果页面的渲染逻辑。

**Acceptance Scenarios**:

1. **Given** 测验结束, **When** 展示得分, **Then** 显示百分制分数（如 85分）而非仅仅显示“答对/总数”。
2. **Given** 得分高于或等于及格线（如60%）, **When** 展示分数, **Then** 分数显示为绿色（Green）。
3. **Given** 得分低于及格线, **When** 展示分数, **Then** 分数显示为红色（Red），并辅助显示“答对题数/总题数”。
4. **Given** 题目列表展示, **When** 渲染题目, **Then** 题号旁明确标识【单选题】、【多选题】等标签。
5. **Given** 题目为多选题, **When** 用户作答时, **Then** 界面提供额外的操作引导说明（如“请选择所有正确选项”）。

---

### User Story 3 - 历史回顾模式 (Priority: P2)

用户希望查看过往测验的详细记录，包括自己的错题、正确答案及解析，以便查漏补缺。

**Why this priority**: 将测验从“一次性测试”转变为“长期复习资料”，大幅提升学习价值。

**Independent Test**: 依赖后端存储历史记录接口，前端可通过 Mock 历史数据列表测试详情页渲染。

**Acceptance Scenarios**:

1. **Given** 用户在测验历史列表, **When** 点击某条记录的“查看详情/回顾”按钮, **Then** 进入该次测验的回顾视图。
2. **Given** 在回顾视图, **When** 查看某道题, **Then** 显示用户的历史作答选项、正确选项以及专家解析内容。
3. **Given** 题目包含解析字段, **When** 渲染题目详情, **Then** 解析区域应清晰展示，且格式美观。

---

### User Story 4 - 全链路入口与元数据展示 (Priority: P3)

用户希望在学习过程中能更容易发现测验入口，并在开始前了解测验的概况（难度、耗时）。

**Why this priority**: 提升功能触达率，帮助用户做好时间管理和心理准备。

**Independent Test**: 检查 UI 布局中是否存在入口链接及元数据绑定。

**Acceptance Scenarios**:

1. **Given** 侧边栏/页头导航, **When** 用户浏览章节内容时, **Then** 可见常驻的“开始测验”入口。
2. **Given** 测验开始页, **When** 页面加载, **Then** 显示该章节题库的简要画像（总题量、预计用时、平均难度）。
3. **Given** 题目加载过程中, **When** 数据未返回, **Then** 显示骨架屏（Skeleton Screen）占位，减少跳动。

---

### Edge Cases

- 当用户在未回答任何问题的情况下点击提交，确认弹窗应明确提示“未答题数：全部”。
- 当后端返回的题目数据缺失“解析”字段时，回顾模式下应优雅隐藏解析区域或显示“暂无解析”。
- 当题目选项超过 26 个（A-Z范围外），标签生成逻辑应能妥善处理（如 AA, AB 或 使用数字），虽极少见但需防御。
- 网络延迟导致提交超时，应允许用户重试，而不是丢失当前答题状态。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统**必须**在用户点击提交时弹出二次确认框，统计并显示“已答题 X 题，未答题 Y 题”。
- **FR-002**: 前端展示题目选项时，**必须**将选项编号固定为 A, B, C, D... 顺序排列，而选项内容可随机映射。
- **FR-003**: 结果页**必须**优先展示百分制得分，并根据是否及格（>=60%）区分显示绿色或红色。
- **FR-004**: 题目显示区域**必须**明确标注题型（单选/多选/改错），并在多选题提供操作提示。
- **FR-005**: 测验历史列表**必须**提供“回顾”功能入口，点击后可跳转至该次测验的详情快照。
- **FR-006**: 回顾模式下，系统**必须**展示题目题干、用户当时的选择、正确答案及解析（如有）。
- **FR-007**: 测验开始页**必须**展示题库元数据：总题量、预计用时、平均难度。
- **FR-008**: 测验加载过程中**必须**使用骨架屏（Skeleton Screen）进行占位展示。
- **FR-009**: 侧边栏或通用导航栏**必须**包含测验功能的快速入口。

### Key Entities *(include if feature involves data)*

- **QuizSession (测验会话)**: 记录一次测验的元数据（开始时间、结束时间、分数）。
- **QuizAttempt (答题记录)**: 记录单道题目的用户作答情况（QuestionID, UserChoice, IsCorrect）。
- **Question (题目扩展)**: 需确认包含 Analysis (解析), Content, Options, Difficulty, EstimatedTime 等字段。

### Assumptions & Dependencies

- **Data Source**: 题库数据源（如 `questions.json` 或数据库）包含或将更新以包含 `analysis` (解析), `difficulty` (难度), `estimated_time` (预计用时) 等字段。
- **Backend Support**: 后端服务已具备或将实现存储完整测验记录（用户每一题的选项和结果）的能力，而不仅仅是存储最终分数。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 用户误触提交测验的反馈为 0（通过确认弹窗机制保证）。
- **SC-002**: 测验结果页展示百分制得分，且包含及格/不及格的视觉区分。
- **SC-003**: 100% 的测验历史记录均可被重新“回顾”，且信息完整（含解析）。
- **SC-004**: 测验加载阶段无明显的布局抖动，骨架屏覆盖率 100%。
- **SC-005**: 侧边栏/导航栏存在可见的测验入口。
