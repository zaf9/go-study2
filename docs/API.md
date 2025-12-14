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

## 测验题库（需认证）

### 随机抽题

- **接口**: `GET /api/v1/quiz/{topic}/{chapter}`
- **描述**: 从指定主题章节的题库中随机抽取题目并创建测验 session（注意：该接口需认证）
- **参数**:
  - `topic`: 主题名称 (lexical_elements, constants, variables, types)
  - `chapter`: 章节名称 (如 comments, boolean, storage 等)
- **响应**:
  ```json
  {
    "code": 20000,
    "message": "success",
    "data": {
      "topic": "constants",
      "chapter": "boolean",
      "sessionId": "session-123",
      "questions": [
        {
          "id": 101,
          "type": "single",
          "difficulty": "easy",
          "question": "Go语言中，布尔类型的零值是？",
          "options": [
            {"id":"A","label":"true"},
            {"id":"B","label":"false"}
          ],
          "codeSnippet": null
        }
      ]
    }
  }
  ```

**说明**: 每次调用返回不同的题目组合；题目选项在返回给前端前会随机打乱以防止答案位置规律。

### 题库统计

- **接口**: `GET /api/v1/quiz/{topic}/{chapter}/stats`
- **描述**: 获取题库的统计信息
- **参数**: 同上
- **响应**:
  ```json
  {
    "code": 20000,
    "message": "success",
    "data": {
      "total": 35,
      "byType": {
        "single": 18,
        "multiple": 17
      },
      "byDifficulty": {
        "easy": 14,
        "medium": 14,
        "hard": 7
      }
    }
  }
  ```

## 错误码速查

- `40001` 未认证或 token 无效
- `40002` refresh token 过期/无效
- `40004` 参数校验失败
- `50001` 服务器内部错误

