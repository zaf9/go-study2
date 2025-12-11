# Data Model - 010-secure-register

## Entities

- **User**
  - Fields: id, username (unique, 3-50, `[A-Za-z0-9_]+`), password_hash, roles (包含 admin 标记), must_change_password (bool), status (active/locked/disabled), created_at, updated_at.
  - Rules: username 唯一；密码哈希使用 bcrypt；must_change_password 初始为 true 对默认管理员；删除/禁用不自动重建。

- **AuthSession/JWT**
  - Fields: access_token, expires_at; refresh_token (cookie), refresh_expires_at, token_hash（存库）。
  - Rules: 登录/注册/改密后刷新令牌对；改密成功后旧令牌作废；refresh_token 存储哈希。

- **AuditEvent**
  - Fields: id, event_type (default_admin_created, register_denied, password_changed), user (nullable for denied), timestamp, result, metadata (IP/UA 摘要).
  - Rules: 记录关键安全事件；不存敏感原文密码。

## Relationships
- User 1 - n AuthSession/RefreshToken。
- AuditEvent 关联 User（若有），无则记录匿名上下文。

## State/Transitions
- 默认管理员创建：缺失时插入 (admin, must_change_password=true)。
- 首次登录：若 must_change_password=true，则登录成功后强制改密；改密成功后 must_change_password=false，旧令牌失效。
- 注册：仅管理员 JWT 可创建新 User；失败不写 User。

## Validation
- Password policy：≥8 位，且包含大小写字母、数字、特殊字符；前后端一致校验；不满足即拒绝。
- JWT：需有效签名与未过期；无效/过期/无权限时注册接口拒绝并审计。

