# Data Model: Quiz UX Enhancement

## Entities

### QuizSession
*Extended from existing concept*
| Field | Type | Description |
|-------|------|-------------|
| `session_id` | UUID | Primary Key, Unique Identifier |
| `user_id` | Int64 | User ID |
| `topic` | String | Topic Identifier (e.g. "lexical") |
| `chapter` | String | Chapter Identifier (e.g. "comment") |
| `total_questions` | Int | Total number of questions in this session |
| `correct_answers` | Int | Number of correct answers |
| `score` | Int | Percentage score (0-100) |
| `passed` | Bool | Whether score >= 60 |
| `started_at` | DateTime | Session start time |
| `completed_at` | DateTime | Session submission time |

### QuizAttempt
*Granular record of each question*
| Field | Type | Description |
|-------|------|-------------|
| `id` | Int64 | Auto-increment PK |
| `session_id` | UUID | Foreign Key to QuizSession |
| `question_id` | String | ID of the specific question (from JSON/DB) |
| `user_choice` | String | User's selected option key (e.g. "A", "B", "AB") |
| `is_correct` | Bool | Whether the specific answer was correct |
| `attempted_at` | DateTime | Time of answering |

## API Contracts

### 1. Submit Quiz (Update)
**POST** `/api/v1/quiz/submit`
- **Request**: Same as existing (Topic, Chapter, Answers map).
- **Response**:
```json
{
  "code": 0,
  "data": {
    "sessionId": "uuid-string",
    "score": 85,
    "passed": true,
    "correctCount": 8,
    "totalCount": 10,
    "details": [ ... ] // Optional immediate feedback
  }
}
```

### 2. Get Quiz History Context
**GET** `/api/v1/quiz/history`
- **Query Params**: `topic` (optional), `chapter` (optional)
- **Response**:
```json
{
  "code": 0,
  "data": {
    "list": [
      {
        "sessionId": "uuid-1",
        "topic": "lexical",
        "chapter": "comment",
        "score": 100,
        "completedAt": "2025-12-18T10:00:00Z"
      }
    ]
  }
}
```

### 3. Get Session Details (Review Mode)
**GET** `/api/v1/quiz/history/{sessionId}`
- **Response**:
```json
{
  "code": 0,
  "data": {
    "meta": {
      "score": 85,
      "passed": true,
      "completedAt": "..."
    },
    "items": [
      {
        "questionId": "q1",
        "stem": "Questions text...",
        "options": ["Option 1", "Option 2"],
        "userChoice": "A", // Index based or Value based? -> Index based for simplicity in frontend mapping? NO, SPEC says stable A-D.
                           // CLARIFICATION IMPLEMENTATION: Backend stores "Value" or "Key"?
                           // DECISION: Backend stores the OPTION TEXT HASH or KEY. 
                           // For MVP: Store the index '0', '1', '2' converted to 'A', 'B', 'C'.
        "correctChoice": "B",
        "isCorrect": false,
        "explanation": "Expert analysis..."
      }
    ]
  }
}
```
---

## Final Decision: userChoice 存储与回顾策略

> **状态**: ✅ RESOLVED (2025-12-18)

### 问题背景

由于选项在前端显示时会被随机打乱（Label A-D 固定，但内容随机），需要明确：
1. 前端提交时发送什么数据？
2. 后端存储什么数据？
3. 回顾模式如何正确显示用户的历史选择？

### 最终方案：**内容值存储 (Content-Based Storage)**

#### 1. 提交流程

```
用户看到:
  A. Go 是编译型语言  ← 用户选择这个
  B. Go 是解释型语言
  C. Go 是脚本语言
  D. Go 是标记语言

前端提交:
{
  "topic": "lexical",
  "chapter": "comment",
  "answers": {
    "q1": "Go 是编译型语言"  // 发送选项内容，不是 "A"
  }
}
```

#### 2. 后端存储

| 字段 | 存储值 | 说明 |
|------|--------|------|
| `user_choice` | `"Go 是编译型语言"` | 存储选项的完整文本内容 |
| `is_correct` | `true` | 后端比对 `user_choice == correct_answer` |

#### 3. 回顾模式响应

```json
{
  "questionId": "q1",
  "stem": "Go 语言的类型是？",
  "options": ["Go 是编译型语言", "Go 是解释型语言", "Go 是脚本语言", "Go 是标记语言"],
  "userChoice": "Go 是编译型语言",
  "correctChoice": "Go 是编译型语言",
  "isCorrect": true,
  "explanation": "Go 是静态类型、编译型语言..."
}
```

#### 4. 前端回顾渲染

```tsx
// 前端根据内容匹配显示标签
options.map((opt, index) => {
  const label = String.fromCharCode(65 + index); // A, B, C, D
  const isUserChoice = opt === userChoice;
  const isCorrect = opt === correctChoice;
  
  return (
    <Option key={index}>
      {label}. {opt}
      {isUserChoice && <Tag color="blue">你的答案</Tag>}
      {isCorrect && <Tag color="green">正确答案</Tag>}
    </Option>
  );
});
```

### 方案优势

| 优势 | 说明 |
|------|------|
| **无状态** | 不需要存储 shuffle seed 或创建 Draft Session |
| **简单可靠** | 内容值比对，无歧义 |
| **回顾友好** | 无论选项顺序如何变化，内容匹配始终准确 |
| **防作弊** | 后端不信任前端发送的 "A/B/C/D"，而是验证实际内容 |

### 边缘情况处理

| 场景 | 处理方式 |
|------|----------|
| 选项内容重复 | 题库设计应避免；如发生，以第一个匹配为准 |
| 选项内容过长 | 存储完整内容，数据库字段使用 TEXT 类型 |
| 多选题 | `user_choice` 存储为逗号分隔的内容列表，如 `"选项A内容,选项C内容"` |

### 实现检查清单

- [ ] 前端 `QuizViewer` 提交时发送选项内容而非标签
- [ ] 后端 `Submit API` 接收并存储选项内容
- [ ] 后端 `Review API` 返回 `userChoice` 和 `correctChoice` 的内容值
- [ ] 前端 `QuizReviewPage` 根据内容匹配渲染标签和状态
