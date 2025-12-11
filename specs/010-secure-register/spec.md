# Feature Specification: 安全注册与默认管理员

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `010-secure-register`  
**Created**: 2025-12-11  
**Status**: Draft  
**Input**: User description: "为了安全起见，注册接口只有在前端登录之后，到jwt认证之后才能调用。同时，每次启动时检查默认用户是否已经存在，不存在则创建。首次登录使用默认用户名和密码：admin/gostudy@123   。首次登录之后强制修改密码。"

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

### User Story 1 - 已登录管理员创建新用户 (Priority: P1)

管理员已通过前端登录并持有有效 JWT，希望调用注册接口为团队创建新账号，注册请求只有在鉴权通过后才被接受。

**Why this priority**: 直接保护注册入口，避免未授权批量注册，是防止账号滥用的首要安全环节。

**Independent Test**: 仅凭拥有有效登录会话与管理员权限即可完成注册操作；缺少令牌或权限时注册被拒绝且返回明确错误。

**Acceptance Scenarios**:

1. **Given** 管理员已登录且持有未过期 JWT，**When** 提交包含合法用户名/密码的新用户注册请求，**Then** 新用户被创建且返回成功信息。
2. **Given** 请求缺少 JWT 或 JWT 已过期，**When** 调用注册接口，**Then** 请求被拒绝并返回鉴权错误，不会创建账户。

---

### User Story 2 - 启动时保障默认管理员存在 (Priority: P1)

运维人员启动服务时，系统自动检查是否存在默认管理员账号；若缺失则创建 admin/gostudy@123 并标记为需改密，保证首次可访问入口。

**Why this priority**: 没有默认管理员将导致环境不可登录；自动创建是确保可运维性的基线。

**Independent Test**: 在空用户库启动服务即可验证默认账号被创建；重复启动不会生成重复账号。

**Acceptance Scenarios**:

1. **Given** 用户表为空，**When** 服务启动完成，**Then** 存在用户名为 admin 的账号，密码为 gostudy@123，且被标记为“首次登录需改密”。
2. **Given** admin 账号已存在（可能已改密），**When** 服务再次启动，**Then** 不会重置其密码或重复创建账户。

---

### User Story 3 - 默认管理员首次登录强制改密 (Priority: P2)

默认管理员首次登录后必须完成密码修改后才能继续任何操作，避免默认口令长期有效。

**Why this priority**: 默认口令是高风险项，强制改密降低被撞库风险。

**Independent Test**: 仅凭默认口令登录会被引导至改密流程；改密完成后才能访问其他功能。

**Acceptance Scenarios**:

1. **Given** admin 使用默认口令首次登录成功，**When** 进入系统，**Then** 所有功能被阻断并提示必须修改密码，直至新密码保存成功。
2. **Given** admin 已完成改密并重新登录，**When** 访问注册接口或其他功能，**Then** 不再出现强制改密提示，且历史默认口令无法再次登录。

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

- 请求携带格式正确但签名无效的 JWT 时，注册接口应返回鉴权失败且无任何副作用。
- 非管理员用户持有有效 JWT 调用注册接口时，应返回权限不足并记录审计。
- 多实例并发启动时默认管理员创建需保证幂等，仅创建一次且不覆盖已有密码。
- 默认管理员已被锁定或禁用时，启动时不得解锁或覆盖其状态，需输出告警。
- 强制改密流程被中断（网络/浏览器关闭）后再次登录，需继续提示改密且不允许其他操作。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 注册接口仅在携带有效且未过期的前端登录 JWT 时可调用；缺少或无效令牌时返回鉴权失败且不创建任何账号。
- **FR-002**: 仅具备管理员权限的已认证用户可注册新用户；权限不足时返回明确错误并记录审计。
- **FR-003**: 每次服务启动完成后，系统检查是否存在用户名为 admin 的账号；若不存在则创建，密码为 gostudy@123，并标记为“首次登录需改密”。
- **FR-004**: 若 admin 账号已存在，启动流程不得覆盖其密码、权限或锁定状态，仅确保“需改密”标记保持或更新为已有值。
- **FR-005**: admin 使用默认口令首次登录后，必须完成密码修改才能继续访问任意业务接口；在改密完成前，所有业务请求应被阻断并提示改密。
- **FR-006**: 密码修改与注册新用户均需满足口令策略（至少8位，且同时包含大小写字母、数字与特殊字符），不符合策略时返回校验错误；现有实现需补充特殊字符校验。
- **FR-007**: 成功创建默认管理员、拒绝未授权注册、完成首次改密等关键动作需生成审计记录，包含时间、操作者身份与结果。

### Key Entities *(include if feature involves data)*

- **User**: 唯一用户名、密码哈希、角色/权限集合、首次登录需改密标记、创建时间、状态（正常/锁定/禁用）。
- **AuthSession/JWT**: 绑定用户身份与权限的令牌，包含签发时间、过期时间、签名；用于校验注册接口访问。
- **AuditEvent**: 记录关键安全事件（默认账号创建、注册被拒、首次改密完成），包含事件类型、主体、时间戳、结果。

## Assumptions

- 注册权限仅授予管理员角色，普通登录用户无注册权限。
- 密码策略按 FR-006 要求执行，适用于注册与改密。
- 部署可能为多实例，默认管理员创建流程需幂等且可在共享存储上检测已有账号。
- 默认管理员被删除或禁用时需要人工介入恢复，本功能仅在缺失时创建，不自动解锁或重置密码。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% 的未携带或无效令牌的注册请求被拒绝并产生审计记录，无任何账号被创建。
- **SC-002**: 在空用户库启动后 5 秒内可通过 admin/gostudy@123 完成首次登录，证明默认管理员已自动创建。
- **SC-003**: 默认管理员首次登录后，在完成改密之前所有业务接口返回“需改密”提示且无一成功调用；改密后再次登录可正常访问。
- **SC-004**: 默认口令被修改后，使用旧口令或旧令牌的登录/注册尝试成功率为 0，相关失败被记录以备审计。
- **SC-005**: 审计记录可在事件发生后 1 分钟内被检索到，覆盖默认账号创建、注册拒绝、首次改密完成三个场景。
