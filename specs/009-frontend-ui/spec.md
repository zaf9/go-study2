# Feature Specification: Go-Study2 前端UI界面

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `009-frontend-ui`  
**Created**: 2025-12-10  
**Status**: Draft  
**Input**: User description: "为 Go-Study2 提供现代化 Web 界面，支持用户注册/登录、学习内容展示与进度跟踪、测验记录、响应式布局与一体化部署。"

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

### User Story 1 - 登录后浏览学习主题 (Priority: P1)

已注册用户登录后，可在浏览器中查看学习主题列表并进入章节阅读，获得代码高亮与分段呈现的内容体验。

**Why this priority**: 直接提供核心学习体验，是后续进度跟踪与测验的基础路径。

**Independent Test**: 仅实现登录与内容浏览即可让用户完成基本学习，能独立验证价值。

**Acceptance Scenarios**:

1. **Given** 用户已成功注册并登录，**When** 访问学习主页，**Then** 可看到按主题分组的学习列表及每个主题的简介与状态。
2. **Given** 用户处于登录状态，**When** 点击某主题进入章节页面，**Then** 可看到格式化的章节内容与代码高亮且无空白占位。

---

### User Story 2 - 进度跟踪与续学 (Priority: P2)

用户在阅读章节时，系统记录学习进度，并在再次登录时提供从上次位置继续学习的入口。

**Why this priority**: 帮助用户保持学习连贯性，避免重复查找，提升留存。

**Independent Test**: 仅实现进度记录与续学入口即可独立验证，且不依赖测验功能。

**Acceptance Scenarios**:

1. **Given** 用户在章节页阅读并滚动，**When** 返回主题列表，**Then** 该章节在列表中被标记为已学习并显示最近访问时间。
2. **Given** 用户退出后再次登录，**When** 进入学习区，**Then** 系统提供“继续上次学习”入口并直接跳转到最近阅读的章节位置。

---

### User Story 3 - 主题测验与成绩查看 (Priority: P3)

用户在学习主题后可以参加测验，提交答案后立即获得成绩，并能查看历史测验记录。

**Why this priority**: 加强学习巩固与反馈，提升互动性。

**Independent Test**: 仅提供测验作答、评分与历史记录展示即可独立完成并验证价值。

**Acceptance Scenarios**:

1. **Given** 用户已登录并选择某主题，**When** 开始测验并提交全部题目，**Then** 系统在同一会话内展示得分与正确/错误摘要。
2. **Given** 用户已有多次测验记录，**When** 打开“历史测验记录”，**Then** 能看到按时间排序的分数、用时与主题信息。

---

### Edge Cases

- 登录会话过期时，访问受保护页面需提示重新登录并在成功后返回原目标。
- 网络或 API 请求失败时，页面应显示错误提示与重试入口，不应留下空白状态。
- 测验提交过程中发生网络中断时，需避免重复计分，并提示用户重新提交或恢复草稿答案（若可用）。
- 移动端小屏设备上，列表与代码区需切换为分段滚动，确保内容可读且不溢出。

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: 系统必须提供用户注册界面，校验用户名与密码的必填性和格式，并对失败原因进行友好提示。
- **FR-002**: 系统必须提供登录与退出功能，登录失败时提示原因，退出后需清除会话状态。
- **FR-003**: 已登录用户会话需保持直至过期或主动退出，支持“记住我”场景并在过期时要求重新登录。
- **FR-004**: 未登录用户访问受保护学习或测验页面时，需被引导至登录页并在完成登录后返回原页面。
- **FR-005**: 登录后必须展示学习主题列表，含主题标题、简介、章节数量及用户当前进度摘要。
- **FR-006**: 用户必须能够进入任一主题查看章节内容，章节需分段呈现并提供代码高亮与基础导航。
- **FR-007**: 系统必须记录用户每个章节的学习状态（未开始/进行中/已完成）及最近访问时间。
- **FR-008**: 系统必须提供“继续上次学习”入口，能跳转到用户最近一次阅读的章节位置。
- **FR-009**: 每个主题必须提供测验入口；用户可开始测验、答题并提交。
- **FR-010**: 系统必须在测验提交后计算得分、展示结果摘要，并保存答题记录（时间、得分、题目结果）。
- **FR-011**: 用户必须能够查看学习进度总览及历史测验记录，支持按时间或主题筛选。
- **FR-012**: 页面必须具备响应式布局，确保桌面与移动设备上登录、学习、测验、进度查看等核心功能均可正常使用。
- **FR-013**: 所有 API/网络请求失败时必须显示可理解的错误提示，并提供重试或返回安全页面的操作。
- **FR-014**: 用户账号信息、学习进度与测验成绩必须进行安全存储与传输，密码不得以明文形式保存。
- **FR-015**: 新增前端界面不得破坏现有命令行模式和既有 HTTP API 的可用性与兼容性。

### Key Entities *(include if feature involves data)*

- **用户（User）**：代表注册账户，包含用户名、认证凭据（安全存储）、会话状态、偏好设置。
- **学习主题（Topic）**：代表某类学习内容，包含标题、简介、章节数量，与章节集合关联。
- **章节（Lesson）**：主题下的具体章节，包含标题、正文内容、代码片段与序号。
- **学习进度（Progress）**：关联用户与章节的学习状态、最近访问时间、最近位置。
- **测验（Quiz）**：关联主题的测验集合，包含题目列表与标准答案。
- **测验记录（QuizAttempt）**：用户对测验的作答记录，包含得分、用时、答题详情、完成时间。

## Assumptions & Dependencies

- 复用现有 Go-Study2 后端 HTTP API 与测验题库，接口语义保持兼容，不影响 CLI 模式。
- 前后端仍由单个服务实例托管，静态资源与 API 同端口提供以简化部署。
- 继续沿用现有持久化方案与配置，用户凭据需使用安全方式保存与验证。
- 目标运行环境为现代桌面与移动浏览器，无需额外客户端安装。

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: 80% 新用户可在 1 分钟内完成注册并登录至学习首页。
- **SC-002**: 90% 已登录用户在 2 秒内看到学习主题列表并可进入任一章节。
- **SC-003**: 95% 用户可在再次登录后通过“一键续学”回到最近学习章节且进度标记准确。
- **SC-004**: 90% 测验在提交后 5 秒内返回得分并写入历史记录可查询。
