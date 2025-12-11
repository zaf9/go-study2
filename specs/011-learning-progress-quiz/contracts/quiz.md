# API Contract: 章节测验

Base: `/api/v1`

## GET /quiz/{topic}/{chapter}
- Purpose: 抽取并返回章节题目（不含答案）
- Auth: 必须登录
- Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "topic": "lexical_elements",
    "chapter": "operators",
    "sessionId": "uuid-string",
    "questions": [
      {
        "id": 1,
        "type": "single",
        "difficulty": "easy",
        "question": "以下哪个标识符是合法的？",
        "options": ["123abc","_privateVar","for","user-name"],
        "codeSnippet": null
      }
    ]
  }
}
```
- Errors: `404` 章节或题库不存在；`401` 未认证；`500` 服务器错误。

## POST /quiz/submit
- Purpose: 提交章节测验答案并判分
- Body:
```json
{
  "sessionId": "uuid-string",
  "topic": "lexical_elements",
  "chapter": "operators",
  "answers": [
    { "questionId": 1, "userAnswers": [1] },
    { "questionId": 2, "userAnswers": [0,2] }
  ]
}
```
- Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "score": 85,
    "total_questions": 10,
    "correct_answers": 9,
    "passed": true,
    "details": [
      {
        "question_id": 1,
        "is_correct": true,
        "correct_answers": [1],
        "explanation": " `_privateVar` 合法，其他非法原因同解析。"
      }
    ]
  }
}
```
- Errors: `400` 参数无效；`401` 未认证；`409` 重复提交同一 session；`500` 服务器错误。

## GET /quiz/history
- Purpose: 拉取测验历史（可分页/过滤）
- Query: `topic`(可选), `limit`(可选,默认10)
- Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "sessions": [
      {
        "id": 1,
        "topic": "lexical_elements",
        "chapter": "operators",
        "score": 85,
        "passed": true,
        "duration": 323,
        "completed_at": "2025-12-10T14:30:00Z"
      }
    ]
  }
}
```

## 判分约束
- 单选/判断：全对得分；多选/改错：全对满分，部分正确=正确数/应选数；含错误选项则 0 分；代码输出题同单选。  
- 通过标准：score ≥ 60；通过后进度状态设为 completed，否则 tested。

