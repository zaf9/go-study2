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

- **FR-001**: 注册界面需校验用户名与密码格式，用户名长度 3-50 且仅限字母数字下划线，密码最少 8 位且包含大小写字母与数字；重复用户名需提示“用户名已存在”且不创建账户；错误提示内容需明确字段原因，展示位置在表单项下方。
- **FR-002**: 登录与退出需提供错误提示（账号不存在/密码错误/锁定），退出必须清除前端会话状态（内存、localStorage 中的 access token 状态位）。
- **FR-003**: 会话策略：access token 有效期 7 天，存储于内存并以 localStorage 仅做页面刷新恢复；refresh token 仅存 HttpOnly Cookie，续期接口 `/api/v1/auth/refresh`；“记住我”勾选时 refresh Cookie 7 天，否则会话级；过期或刷新失败需要求重新登录。
- **FR-004**: 未登录或 access token 过期访问受保护页面时，必须重定向至登录页并在登录成功后返回原目标路由；受保护路由与 Edge Case 描述保持一致，不得出现空白页面。
- **FR-005**: 学习主题列表必须包含标题、简介、章节数量、用户进度百分比（已完成章节数/总章节数），默认按预设顺序排列；若列表为空需给出占位与下一步引导。
- **FR-006**: 章节内容需分段呈现并支持代码高亮；高亮语言范围至少包含 Go、TypeScript、JavaScript、JSON、bash、markdown；不支持的语言需回退为纯文本；需提供章节内导航或分段锚点。
- **FR-007**: 学习进度记录必须包含状态（not_started/in_progress/done）、最近访问时间 lastVisit、最近位置 lastPosition（滚动或锚点）；(user, topic, chapter) 组合应唯一，更新为幂等覆盖。
- **FR-008**: “继续上次学习”入口需优先使用 lastPosition，其次 lastVisit；若两者缺失则跳转到该主题首个未完成章节；需在列表页可见且指向明确。
- **FR-009**: 每个主题需提供测验入口，题型支持单选/多选；需说明题目来源与数量范围（由后端返回）；提交需防重复。
- **FR-010**: 评分算法需定义：单选/多选均按全对计分，score ≤ total；需保存答题明细（题号、选择、正确答案、得分、用时）。
- **FR-011**: 测验历史需支持按主题与时间筛选，展示时间、得分、题目结果摘要；列表为空时需给出引导或提示。
- **FR-012**: 响应式要求需明确断点：Mobile <768px，Tablet 768-1024px，Desktop >1024px；列表列数、字号、间距随断点变化的规则需给出（至少：桌面 3-4 列，移动单列）。
- **FR-013**: 所有 API/网络请求失败需显示统一格式错误提示（含原因与重试入口），提示位置需固定（页面顶部或组件内显著区域），禁止留空白状态。
- **FR-014**: 安全存储：密码使用 bcrypt 哈希；access token 仅前端内存 + localStorage 恢复，refresh token HttpOnly Cookie；敏感信息传输需使用 HTTPS，避免明文存储。
- **FR-015**: 新增前端界面不得破坏现有命令行模式和既有 HTTP API 的可用性与兼容性，需复用既有路由/响应格式。
- **FR-016**: 静态导出目录为 `frontend/out`，GoFrame 静态文件服务需托管 `/`，API 走 `/api/*`，路由冲突时 `/api/*` 优先；需提供 404/回退页。

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

- **SC-001**: 注册失败提示覆盖 100% 常见校验错误（格式、重复用户名），无空白或模糊错误信息。
- **SC-002**: 登录/访问受保护路由的过期与刷新失败场景，100% 出现明确提示并能返回登录页。
- **SC-003**: 章节学习状态与“继续上次学习”跳转的一致性达到 100%（最近位置与进度标记一致，无跳转偏差）。
- **SC-004**: 测验提交后，评分结果与历史记录字段（时间、得分、题目结果摘要）完整写入，缺失率为 0。
