# API Contract: GET /api/v1/progress/last

**Feature**: Dashboard 首页 - 获取最后学习记录  
**Date**: 2025-12-26  
**Version**: v1  
**Related**: [spec.md](../spec.md) | [data-model.md](../data-model.md)

## Overview

此 API 返回用户最后一次学习的主题和章节信息，用于 Dashboard 的"快速继续"功能。

## Endpoint

```
GET /api/v1/progress/last
```

## Authentication

**Required**: Yes  
**Method**: JWT Token 或 Session Cookie  
**Header**: `Authorization: Bearer <token>` 或通过 Cookie 自动携带

## Request

### Headers

```http
GET /api/v1/progress/last HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

### Query Parameters

无

### Request Body

无

## Response

### Success Response (200 OK)

**有学习记录时**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "topic_id": "lexical-elements",
    "topic_name": "Lexical Elements",
    "topic_display_name": "词法元素",
    "chapter_id": "identifiers",
    "chapter_name": "Identifiers",
    "chapter_display_name": "标识符",
    "last_visited_at": "2025-12-26T10:30:00+08:00"
  }
}
```

**无学习记录时**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### Error Responses

#### 401 Unauthorized (未认证)

```json
{
  "code": 40001,
  "message": "未登录或登录已过期",
  "data": null
}
```

#### 500 Internal Server Error (服务器错误)

```json
{
  "code": 50000,
  "message": "服务器内部错误",
  "data": null
}
```

## Response Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `code` | integer | Yes | 响应码，0 表示成功 |
| `message` | string | Yes | 响应消息 |
| `data` | object\|null | Yes | 数据对象，无记录时为 null |
| `data.topic_id` | string | Yes | 主题 ID |
| `data.topic_name` | string | Yes | 主题英文名称 |
| `data.topic_display_name` | string | Yes | 主题中文显示名称 |
| `data.chapter_id` | string | Yes | 章节 ID |
| `data.chapter_name` | string | Yes | 章节英文名称 |
| `data.chapter_display_name` | string | Yes | 章节中文显示名称 |
| `data.last_visited_at` | string | Yes | 最后访问时间（ISO 8601 格式） |

## Business Rules

1. **用户隔离**: 只返回当前登录用户的学习记录
2. **最新记录**: 按 `last_visited_at` 降序排序，取第一条
3. **数据完整性**: 如果主题或章节已被删除，返回 `data: null`
4. **时区**: 时间戳使用服务器时区（UTC+8）

## Implementation Notes

### Backend (Go)

```go
// Controller
func (c *ProgressController) GetLastLearning(ctx *ghttp.Request) {
    userID := ctx.GetCtxVar("user_id").Uint()
    
    lastLearning, err := c.progressService.GetLastLearningRecord(userID)
    if err != nil {
        ctx.Response.WriteJsonExit(g.Map{
            "code":    50000,
            "message": "服务器内部错误",
            "data":    nil,
        })
        return
    }
    
    ctx.Response.WriteJsonExit(g.Map{
        "code":    0,
        "message": "success",
        "data":    lastLearning,
    })
}

// Service
func (s *ProgressService) GetLastLearningRecord(userID uint) (*LastLearningRecord, error) {
    var progress LearningProgress
    err := s.db.Model(&LearningProgress{}).
        Where("user_id = ? AND last_visited_at IS NOT NULL", userID).
        Order("last_visited_at DESC").
        First(&progress).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // 无学习记录，返回 nil
        }
        return nil, err
    }
    
    // 获取主题和章节信息
    var topic Topic
    var chapter Chapter
    
    err = s.db.First(&topic, "id = ?", progress.TopicID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // 主题已删除
        }
        return nil, err
    }
    
    err = s.db.First(&chapter, "id = ?", progress.ChapterID).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // 章节已删除
        }
        return nil, err
    }
    
    return &LastLearningRecord{
        TopicID:           progress.TopicID,
        TopicName:         topic.Name,
        TopicDisplayName:  topic.DisplayName,
        ChapterID:         progress.ChapterID,
        ChapterName:       chapter.Name,
        ChapterDisplayName: chapter.DisplayName,
        LastVisitedAt:     progress.LastVisitedAt.Format(time.RFC3339),
    }, nil
}
```

### Frontend (TypeScript)

```typescript
// API 调用
export async function getLastLearning(): Promise<LastLearningData | null> {
  const response = await api.get<ApiResponse<LastLearningData | null>>(
    '/api/v1/progress/last'
  )
  
  if (response.data.code !== 0) {
    throw new Error(response.data.message)
  }
  
  return response.data.data
}

// 类型定义
interface LastLearningData {
  topic_id: string
  topic_name: string
  topic_display_name: string
  chapter_id: string
  chapter_name: string
  chapter_display_name: string
  last_visited_at: string
}

interface ApiResponse<T> {
  code: number
  message: string
  data: T
}
```

## Testing

### Unit Test (Backend)

```go
func TestProgressController_GetLastLearning(t *testing.T) {
    // Test case 1: 有学习记录
    t.Run("WithLearningRecord", func(t *testing.T) {
        // Setup mock data
        // Call API
        // Assert response
    })
    
    // Test case 2: 无学习记录
    t.Run("NoLearningRecord", func(t *testing.T) {
        // Setup empty user
        // Call API
        // Assert data is null
    })
    
    // Test case 3: 未认证
    t.Run("Unauthorized", func(t *testing.T) {
        // Call API without token
        // Assert 401 error
    })
    
    // Test case 4: 主题已删除
    t.Run("TopicDeleted", func(t *testing.T) {
        // Setup progress with deleted topic
        // Call API
        // Assert data is null
    })
}
```

### Integration Test (Frontend)

```typescript
describe('getLastLearning', () => {
  it('should return last learning record', async () => {
    const mockData = {
      topic_id: 'lexical-elements',
      topic_display_name: '词法元素',
      chapter_id: 'identifiers',
      chapter_display_name: '标识符',
      last_visited_at: '2025-12-26T10:30:00+08:00'
    }
    
    mockApi.onGet('/api/v1/progress/last').reply(200, {
      code: 0,
      message: 'success',
      data: mockData
    })
    
    const result = await getLastLearning()
    expect(result).toEqual(mockData)
  })
  
  it('should return null when no record', async () => {
    mockApi.onGet('/api/v1/progress/last').reply(200, {
      code: 0,
      message: 'success',
      data: null
    })
    
    const result = await getLastLearning()
    expect(result).toBeNull()
  })
})
```

## Performance

- **Expected Response Time**: < 100ms
- **Database Queries**: 1-3 queries (progress + topic + chapter)
- **Caching**: 可考虑缓存 5-10 秒（用户学习记录变化频率低）

## Security

- **Authentication**: 必须验证用户身份
- **Authorization**: 只能访问自己的学习记录
- **SQL Injection**: 使用参数化查询防止注入
- **Data Leakage**: 不返回其他用户的数据

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| v1 | 2025-12-26 | 初始版本 |
