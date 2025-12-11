# API 文档（简版）

统一响应格式：

```json
{
  "code": 20000,
  "message": "success",
  "data": {}
}
```

- 成功：`code=20000`
- 参数错误：`40004`
- 认证失败：`40001/40002`
- 服务器错误：`50001`

## 认证

- `POST /api/v1/auth/register` — body: `{username, password, remember}` → `data.accessToken`, `expiresIn`
- `POST /api/v1/auth/login` — 同上
- `POST /api/v1/auth/refresh` — 使用 Cookie 中的 refresh token
- `POST /api/v1/auth/logout` — 需 `Authorization: Bearer <access>`
- `GET /api/v1/auth/profile` — 返回 `{id, username}`

## 学习内容

- `GET /api/v1/topics` — 主题列表 `{topics:[{id,title,description}]}`
- `GET /api/v1/topic/{topic}` — 章节菜单 `{items:[{id,title,name}]}`
- `GET /api/v1/topic/{topic}/{chapter}` — 章节内容 `{title, content}`

## 学习进度（需认证）

- `GET /api/v1/progress` — 全量进度列表
- `GET /api/v1/progress/{topic}` — 指定主题进度
- `POST /api/v1/progress` — body: `{topic, chapter, status, position?}`，幂等覆盖

## 测验（需认证）

- `GET /api/v1/quiz/{topic}/{chapter}` — 返回题目列表
- `POST /api/v1/quiz/submit` — body: `{topic, chapter, answers:[{id,choices[]}]}`
- `GET /api/v1/quiz/history` / `quiz/history/{topic}` — 历史记录列表

## 错误码速查

- `40001` 未认证或 token 无效
- `40002` refresh token 过期/无效
- `40004` 参数校验失败
- `50001` 服务器内部错误

