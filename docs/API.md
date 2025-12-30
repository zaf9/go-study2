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

### 获取进度概览

- `GET /api/v1/progress/overview` — 获取全局学习进度概览
  - **响应**:
    ```json
    {
      "code": 20000,
      "message": "success",
      "data": {
        "totalChapters": 100,
        "completedChapters": 25,
        "completionRate": 25.0,
        "topics": [
          {
            "topic": "lexical_elements",
            "totalChapters": 11,
            "completedChapters": 5
          }
        ],
        "next": {
          "topic": "constants",
          "chapter": "boolean",
          "title": "布尔常量"
        }
      }
    }
    ```

### 获取主题进度

- `GET /api/v1/progress/topic/{topic}` — 获取指定主题的章节进度详情
  - **响应**:
    ```json
    {
      "code": 20000,
      "message": "success",
      "data": {
        "topic": "lexical_elements",
        "totalChapters": 11,
        "chapters": [
          {
            "chapter": "comments",
            "status": "completed",
            "quizScore": 90,
            "quizPassed": true,
            "lastVisitAt": "2025-12-30T10:00:00Z",
            "completedAt": "2025-12-30T10:15:00Z"
          },
          {
            "chapter": "tokens",
            "status": "in_progress",
            "lastVisitAt": "2025-12-30T14:00:00Z"
          },
          {
            "chapter": "semicolons",
            "status": "not_started"
          }
        ]
      }
    }
    ```

### 获取全部进度（旧版，兼容保留）

- `GET /api/v1/progress` — 全量进度列表
- `GET /api/v1/progress/{topic}` — 指定主题进度

### 更新进度

- `POST /api/v1/progress` — body: `{topic, chapter, status, position?}`，幂等覆盖

## 测验（需认证）

### 获取测验题目

- `GET /api/v1/quiz/{topic}/{chapter}` — 获取或创建测验会话
  - **描述**: 从指定主题章节的题库中随机抽取题目并创建测验 session
  - **参数**: 
    - `topic`: 主题ID（如 lexical_elements, constants, variables, types）
    - `chapter`: 章节ID（如 comments, boolean, declarations）
  - **响应**:
    ```json
    {
      "code": 20000,
      "message": "success",
      "data": {
        "sessionId": "uuid-session-id",
        "topic": "lexical_elements",
        "chapter": "comments",
        "questions": [
          {
            "id": 1,
            "question": "Go语言支持哪些注释方式？",
            "options": ["A. 单行注释", "B. 多行注释", "C. 文档注释"],
            "type": "multiple",
            "difficulty": "easy"
          }
        ]
      }
    }
    ```

### 提交测验

- `POST /api/v1/quiz/submit` — 提交测验答案并评分
  - **请求体**:
    ```json
    {
      "sessionId": "uuid-session-id",
      "topic": "lexical_elements",
      "chapter": "comments",
      "durationMs": 120000,
      "answers": [
        {
          "questionId": 1,
          "userAnswers": ["A", "B"]
        }
      ]
    }
    ```
  - **响应**:
    ```json
    {
      "code": 20000,
      "message": "success",
      "data": {
        "score": 80,
        "totalQuestions": 10,
        "correctAnswers": 8,
        "details": [
          {
            "question_id": 1,
            "is_correct": true,
            "user_answers": ["A", "B"],
            "correct_answers": ["A", "B"]
          }
        ]
      }
    }
    ```
  - **幂等性**: 重复提交同一 sessionId 会返回 409 Conflict

### 获取测验历史

- `GET /api/v1/quiz/history` — 获取所有测验记录
  - **查询参数**:
    - `topic` (可选): 按主题过滤
  - **响应**:
    ```json
    {
      "code": 20000,
      "message": "success",
      "data": [
        {
          "sessionId": "uuid-session-id",
          "topic": "lexical_elements",
          "chapter": "comments",
          "score": 80,
          "totalQuestions": 10,
          "correctAnswers": 8,
          "createdAt": "2025-12-30T10:00:00Z",
          "submittedAt": "2025-12-30T10:15:00Z"
        }
      ]
    }
    ```

- `GET /api/v1/quiz/history/{sessionId}` — 获取特定测验详情
  - **响应**: 包含完整的题目、用户答案和正确答案对比

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

