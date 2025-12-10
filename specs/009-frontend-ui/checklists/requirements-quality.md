# Checklist: Requirements Quality - Go-Study2 前端 UI
Purpose: 确保前端 UI 需求的完整性、清晰度、一致性与可衡量性  
Created: 2025-12-10  
Feature: 009-frontend-ui  
Docs: D:\studyspace\go-study\go-study2\specs\009-frontend-ui\spec.md

## Requirement Completeness
- [ ] CHK001 是否覆盖注册/登录的用户名长度、密码复杂度、重复注册处理等校验要求？[Completeness, Spec §FR-001]
- [ ] CHK002 登录与退出后会话清理（access/refresh/token 缓存、前端状态）是否明确并包含“记住我”场景？[Completeness, Spec §FR-002, Spec §FR-003]
- [ ] CHK003 学习主题列表需展示的字段（标题/简介/章节数/进度摘要）及排序或分页策略是否被完整要求？[Completeness, Spec §FR-005]
- [ ] CHK004 学习进度记录的字段（状态、最近访问时间、位置）及幂等/覆盖规则是否定义？[Completeness, Spec §FR-007]
- [ ] CHK005 测验入口、题型范围、评分算法、历史存储字段（时间、得分、题目结果）是否被要求齐备？[Completeness, Spec §FR-009, Spec §FR-010, Spec §FR-011]
- [ ] CHK006 静态导出路由清单与后端托管目录结构（`frontend/out`、路由 `/api/*` vs `/*`）是否有明确要求？[Completeness, Plan §Project Structure]

## Requirement Clarity
- [ ] CHK007 “友好提示/错误提示”在内容、位置、展示组件上是否有可执行的明确描述？[Clarity, Spec §FR-001, Spec §FR-002, Spec §FR-013]
- [ ] CHK008 “继续上次学习”跳转依据 lastVisit 与 lastPosition 的优先级是否被写清，冲突时的规则是否明确？[Clarity, Spec §FR-008]
- [ ] CHK009 代码高亮的语言范围、主题样式及降级策略（不支持语言时如何展示）是否明示？[Clarity, Spec §FR-006, Research §1]
- [ ] CHK010 响应式断点与布局变化（列数、字号、间距）是否被量化说明？[Clarity, Spec §FR-012]
- [ ] CHK011 “记住我”持续时长与存储位置（Cookie/localStorage）是否有明确要求？[Clarity, Spec §FR-003, Research §2]

## Requirement Consistency
- [ ] CHK012 受保护路由的重定向/返回逻辑是否与会话过期 Edge Case 描述保持一致且无冲突？[Consistency, Spec §FR-004, Spec §Edge Cases]
- [ ] CHK013 测验历史的筛选要求是否与测验记录数据模型字段保持一致（topic、time、score）？[Consistency, Spec §FR-011, Data-Model §QuizRecord]
- [ ] CHK014 API 统一错误码/响应结构是否与既有后端契约保持一致？[Consistency, Contracts §openapi.yaml, Spec §FR-013]

## Acceptance Criteria Quality
- [ ] CHK015 SC-001~SC-004 的测量方法（计时起点、网络/设备环境、样本量）是否被定义以便验证？[Acceptance Criteria, Spec §SC-001..004]
- [ ] CHK016 性能目标（5 秒列表、10 秒测验评分、首屏 JS < 200KB）是否给出度量工具与环境前提？[Measurability, Plan §Performance Goals]

## Scenario Coverage
- [ ] CHK017 未登录、登录过期、刷新失败三类访问受保护页面的处理与提示是否均有要求？[Coverage, Spec §FR-003, Spec §FR-004, Spec §Edge Cases]
- [ ] CHK018 章节内容加载失败或部分加载时的占位、重试与错误信息是否在需求中覆盖？[Coverage, Spec §FR-006, Spec §Edge Cases]
- [ ] CHK019 测验提交网络中断或重复提交的防重与恢复要求是否明确？[Coverage, Spec §FR-010, Spec §Edge Cases]

## Edge Case Coverage
- [ ] CHK020 移动端小屏的分段滚动、代码块溢出与可读性要求是否明确？[Edge Case, Spec §Edge Cases]
- [ ] CHK021 主题/章节为空或未发布时的展示、导航与提示要求是否存在？[Edge Case, Gap]
- [ ] CHK022 SQLite 写入冲突或容量上限接近时的用户提示/限制要求是否记录？[Edge Case, Gap, Plan §Constraints]

## Non-Functional Requirements
- [ ] CHK023 安全要求是否覆盖 token 存储位置、过期策略、refresh 调用频率与失败处理？[Security, Plan §Constraints, Research §2]
- [ ] CHK024 可访问性要求（键盘导航、ARIA、对比度、焦点顺序）是否被写入？[Gap, Frontend Principles §XXIX]
- [ ] CHK025 日志与监控指标（请求耗时、错误率、注册/测验业务指标）是否在需求中提出？[Gap, Spec §Success Criteria]

## Dependencies & Assumptions
- [ ] CHK026 复用既有 API 的兼容性约束（版本、字段、错误码）是否记录且与新前端需求对齐？[Dependency, Spec §Assumptions]
- [ ] CHK027 静态文件与 API 同端口托管的缓存策略、路由优先级、fallback 规则是否被描述？[Dependency, Spec §Assumptions, Plan §Target Platform]

## Ambiguities & Conflicts
- [ ] CHK028 “不破坏现有 CLI/API”与新增前端缓存/重定向逻辑是否存在潜在冲突，是否需要明确回退策略？[Conflict, Spec §FR-015]
- [ ] CHK029 错误提示的文案与展示方式是否要求全局统一以避免页面间不一致？[Consistency, Spec §FR-013]
- [ ] CHK030 学习进度与测验历史的数据保留期限是否明确，以避免与存储约束冲突？[Ambiguity, Plan §Constraints, Spec §FR-011]

