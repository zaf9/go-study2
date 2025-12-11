# API Contract: 学习进度

Base: `/api/v1`

## POST /progress
- Purpose: 创建/更新章节进度（阅读/滚动/测验结果）
- Auth: 必须登录
- Body:
```json
{
  "topic": "lexical_elements",
  "chapter": "comments",
  "read_duration": 120,        // 可选，累加
  "scroll_progress": 75,       // 可选，覆盖 0-100
  "last_position": 1200,       // 可选，覆盖
  "quiz_score": 85,            // 可选，测验提交
  "quiz_passed": true          // 可选，测验提交
}
```
- Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "in_progress",
    "overall_progress": 45,
    "topic_progress": 73,
    "read_duration": 180,
    "scroll_progress": 80,
    "last_position": 1400
  }
}
```
- Errors: `400` 参数缺失/非法；`401` 未认证；`409` 状态回退冲突（返回最新状态）；`500` 服务器错误。

## GET /progress
- Purpose: 获取整体进度与主题汇总
- Query: none
- Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "overall": {
      "progress": 45,
      "completedChapters": 12,
      "totalChapters": 48,
      "studyDays": 15,
      "totalStudyTime": 3600
    },
    "topics": [
      {
        "id": "lexical_elements",
        "name": "Lexical Elements",
        "weight": 20,
        "progress": 73,
        "completedChapters": 8,
        "totalChapters": 11,
        "lastVisitAt": "2025-12-10T14:30:00Z"
      }
    ]
  }
}
```

## GET /progress/{topic}
- Purpose: 获取指定主题的章节进度列表
- Path: `topic` = lexical_elements|constants|variables|types
- Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "lexical_elements",
    "name": "Lexical Elements",
    "weight": 20,
    "progress": 73,
    "completedChapters": 8,
    "totalChapters": 11,
    "chapters": [
      {
        "id": "comments",
        "name": "注释",
        "status": "completed",
        "progress": 100,
        "quizScore": 85,
        "quizPassed": true,
        "lastVisitAt": "2025-12-10T14:30:00Z"
      }
    ]
  }
}
```

## Status Codes
- `not_started` | `in_progress` | `completed` | `tested` (测验但未通过)  
- 完成条件：read_duration ≥ 80% 预估、scroll_progress ≥ 90%、quiz_passed==true；通过测验自动置 completed。

