# Data Model: Dashboard 首页功能

**Feature**: Dashboard 首页  
**Date**: 2025-12-26  
**Related**: [spec.md](./spec.md) | [plan.md](./plan.md)

## Terminology (术语表)

为确保文档一致性，以下术语在整个项目中统一使用：

| 中文术语 | 英文术语 (代码/API) | 说明 |
|---------|-------------------|------|
| 最后学习记录 | last learning record | 用户最近一次访问的章节信息 |
| 主题进度 / 主题完成百分比 | topic progress / topic completion percentage | 某个主题下已完成章节占总章节的比例 |
| 累计学习天数 | total learning days / learning days count | 有学习活动的不同日期数（非连续） |
| Dashboard 首页 | Dashboard homepage | 用户登录后的默认首页，路径为根路径 `/` |

**使用规范**:
- 用户界面文案：使用中文术语
- 代码变量/函数名：使用英文术语（驼峰或蛇形命名）
- API 字段名：使用英文术语（蛇形命名）
- 文档说明：优先使用中文术语，必要时注明英文

## Overview

Dashboard 功能主要复用现有数据模型，无需新增数据表。本文档定义 Dashboard 所需的数据结构、计算逻辑和数据关系。

## Core Entities

### 1. User (用户) - 现有实体

**用途**: 代表学习者

**字段**:
```go
type User struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Dashboard 使用**:
- 显示欢迎信息中的用户名
- WebSocket 连接的用户标识
- 数据查询的过滤条件

---

### 2. LearningProgress (学习进度记录) - 现有实体

**用途**: 记录用户对每个章节的学习状态

**字段**:
```go
type LearningProgress struct {
    ID          uint      `json:"id"`
    UserID      uint      `json:"user_id"`
    TopicID     string    `json:"topic_id"`
    ChapterID   string    `json:"chapter_id"`
    Status      string    `json:"status"`        // not_started, in_progress, completed
    Completed   bool      `json:"completed"`
    LastVisitedAt time.Time `json:"last_visited_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Dashboard 使用**:
- 计算累计学习天数（基于 `created_at` 或 `last_visited_at`）
- 计算总章节完成进度（`completed = true` 的记录数）
- 获取最后学习记录（按 `last_visited_at` 降序排序）
- 计算各主题完成百分比

**计算逻辑**:
```sql
-- 累计学习天数（有学习活动的不同日期数）
SELECT COUNT(DISTINCT DATE(last_visited_at)) 
FROM learning_progress 
WHERE user_id = ? AND last_visited_at IS NOT NULL;

-- 总章节完成进度
SELECT COUNT(*) 
FROM learning_progress 
WHERE user_id = ? AND completed = true;

-- 最后学习记录
SELECT topic_id, chapter_id, last_visited_at
FROM learning_progress
WHERE user_id = ?
ORDER BY last_visited_at DESC
LIMIT 1;

-- 主题完成百分比
SELECT 
    topic_id,
    COUNT(*) as total_chapters,
    SUM(CASE WHEN completed = true THEN 1 ELSE 0 END) as completed_chapters,
    (SUM(CASE WHEN completed = true THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) as percentage
FROM learning_progress
WHERE user_id = ?
GROUP BY topic_id;
```

---

### 3. Topic (主题) - 现有实体

**用途**: 代表一个学习主题

**字段**:
```go
type Topic struct {
    ID          string    `json:"id"`           // 如 "lexical-elements"
    Name        string    `json:"name"`         // 英文名称
    DisplayName string    `json:"display_name"` // 中文显示名称
    Description string    `json:"description"`
    ChapterCount int      `json:"chapter_count"`
    CreatedAt   time.Time `json:"created_at"`
}
```

**Dashboard 使用**:
- 显示所有可学习的主题列表
- 提供主题名称和章节总数
- 跳转到主题详情页的链接

---

### 4. Chapter (章节) - 现有实体

**用途**: 代表主题下的一个学习章节

**字段**:
```go
type Chapter struct {
    ID        string    `json:"id"`
    TopicID   string    `json:"topic_id"`
    Name      string    `json:"name"`
    DisplayName string  `json:"display_name"`
    Content   string    `json:"content"`
    Order     int       `json:"order"`
    CreatedAt time.Time `json:"created_at"`
}
```

**Dashboard 使用**:
- 计算总章节数
- 提供章节名称用于"快速继续"功能
- 跳转到章节详情页的链接

---

### 5. QuizRecord (测验记录) - 现有实体

**用途**: 记录用户的测验历史

**字段**:
```go
type QuizRecord struct {
    ID          uint      `json:"id"`
    UserID      uint      `json:"user_id"`
    TopicID     string    `json:"topic_id"`
    ChapterID   string    `json:"chapter_id"`
    Score       int       `json:"score"`
    TotalQuestions int    `json:"total_questions"`
    Passed      bool      `json:"passed"`
    CompletedAt time.Time `json:"completed_at"`
    CreatedAt   time.Time `json:"created_at"`
}
```

**Dashboard 使用**:
- 显示最近 3-5 条测验记录
- 显示得分、主题/章节名称、完成时间
- 跳转到测验详情页的链接

**查询逻辑**:
```sql
-- 获取最近测验记录
SELECT 
    qr.id,
    t.display_name as topic_name,
    c.display_name as chapter_name,
    qr.score,
    qr.total_questions,
    qr.completed_at
FROM quiz_records qr
JOIN topics t ON qr.topic_id = t.id
JOIN chapters c ON qr.chapter_id = c.id
WHERE qr.user_id = ?
ORDER BY qr.completed_at DESC
LIMIT 5;
```

---

## Derived Data Structures (派生数据结构)

这些数据结构不对应数据库表，而是通过计算或聚合现有数据得出。

### DashboardStats (Dashboard 统计数据)

**用途**: 汇总用户的学习统计信息

```typescript
interface DashboardStats {
  studyDays: number              // 累计学习天数（有学习活动的不同日期数）
  totalChapters: number          // 总章节数
  completedChapters: number      // 已完成章节数
  progressPercentage: number     // 整体完成百分比（保留一位小数）
  weeklyActivity: number         // 本周学习活动次数（本周完成章节数）
}
```

**计算逻辑**:
```go
func CalculateDashboardStats(userID uint) (*DashboardStats, error) {
    stats := &DashboardStats{}
    
    // 1. 累计学习天数
    var studyDays int
    err := db.Model(&LearningProgress{}).
        Where("user_id = ? AND last_visited_at IS NOT NULL", userID).
        Select("COUNT(DISTINCT DATE(last_visited_at))").
        Scan(&studyDays).Error
    if err != nil {
        return nil, err
    }
    stats.StudyDays = studyDays
    
    // 2. 总章节数
    var totalChapters int64
    err = db.Model(&Chapter{}).Count(&totalChapters).Error
    if err != nil {
        return nil, err
    }
    stats.TotalChapters = int(totalChapters)
    
    // 3. 已完成章节数
    var completedChapters int64
    err = db.Model(&LearningProgress{}).
        Where("user_id = ? AND completed = true", userID).
        Count(&completedChapters).Error
    if err != nil {
        return nil, err
    }
    stats.CompletedChapters = int(completedChapters)
    
    // 4. 整体完成百分比
    if stats.TotalChapters > 0 {
        stats.ProgressPercentage = float64(stats.CompletedChapters) / float64(stats.TotalChapters) * 100
        stats.ProgressPercentage = math.Round(stats.ProgressPercentage*10) / 10 // 保留一位小数
    }
    
    // 5. 本周学习活动次数
    weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
    var weeklyActivity int64
    err = db.Model(&LearningProgress{}).
        Where("user_id = ? AND updated_at >= ?", userID, weekStart).
        Count(&weeklyActivity).Error
    if err != nil {
        return nil, err
    }
    stats.WeeklyActivity = int(weeklyActivity)
    
    return stats, nil
}
```

---

### LastLearningRecord (最后学习记录)

**用途**: 用户最后一次学习的主题和章节信息

```typescript
interface LastLearningRecord {
  topicId: string
  topicName: string
  chapterId: string
  chapterName: string
  lastVisitedAt: string  // ISO 8601 格式
}
```

**计算逻辑**:
```go
func GetLastLearningRecord(userID uint) (*LastLearningRecord, error) {
    var progress LearningProgress
    err := db.Model(&LearningProgress{}).
        Where("user_id = ? AND last_visited_at IS NOT NULL", userID).
        Order("last_visited_at DESC").
        First(&progress).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // 无学习记录
        }
        return nil, err
    }
    
    // 获取主题和章节名称
    var topic Topic
    var chapter Chapter
    
    err = db.First(&topic, "id = ?", progress.TopicID).Error
    if err != nil {
        return nil, err
    }
    
    err = db.First(&chapter, "id = ?", progress.ChapterID).Error
    if err != nil {
        return nil, err
    }
    
    return &LastLearningRecord{
        TopicID:       progress.TopicID,
        TopicName:     topic.DisplayName,
        ChapterID:     progress.ChapterID,
        ChapterName:   chapter.DisplayName,
        LastVisitedAt: progress.LastVisitedAt.Format(time.RFC3339),
    }, nil
}
```

---

### TopicProgressSummary (主题进度汇总)

**用途**: 每个主题的学习进度统计

```typescript
interface TopicProgressSummary {
  topicId: string
  topicName: string
  displayName: string
  completedChapters: number
  totalChapters: number
  percentage: number  // 保留一位小数
}
```

**计算逻辑**:
```go
func GetTopicProgressSummary(userID uint) ([]TopicProgressSummary, error) {
    // 1. 获取所有主题
    var topics []Topic
    err := db.Find(&topics).Error
    if err != nil {
        return nil, err
    }
    
    summaries := make([]TopicProgressSummary, 0, len(topics))
    
    for _, topic := range topics {
        // 2. 计算该主题的完成情况
        var totalChapters int64
        err = db.Model(&Chapter{}).
            Where("topic_id = ?", topic.ID).
            Count(&totalChapters).Error
        if err != nil {
            return nil, err
        }
        
        var completedChapters int64
        err = db.Model(&LearningProgress{}).
            Where("user_id = ? AND topic_id = ? AND completed = true", userID, topic.ID).
            Count(&completedChapters).Error
        if err != nil {
            return nil, err
        }
        
        percentage := 0.0
        if totalChapters > 0 {
            percentage = float64(completedChapters) / float64(totalChapters) * 100
            percentage = math.Round(percentage*10) / 10
        }
        
        summaries = append(summaries, TopicProgressSummary{
            TopicID:           topic.ID,
            TopicName:         topic.Name,
            DisplayName:       topic.DisplayName,
            CompletedChapters: int(completedChapters),
            TotalChapters:     int(totalChapters),
            Percentage:        percentage,
        })
    }
    
    return summaries, nil
}
```

---

### RecentQuizSummary (最近测验汇总)

**用途**: 最近测验记录的展示数据

```typescript
interface RecentQuizSummary {
  id: number
  topicName: string
  chapterName: string
  score: number
  totalQuestions: number
  completedAt: string  // ISO 8601 格式
  displayTime: string  // 格式化后的时间（相对或绝对）
}
```

**计算逻辑**:
```go
func GetRecentQuizzes(userID uint, limit int) ([]RecentQuizSummary, error) {
    var records []QuizRecord
    err := db.Model(&QuizRecord{}).
        Where("user_id = ?", userID).
        Order("completed_at DESC").
        Limit(limit).
        Find(&records).Error
    
    if err != nil {
        return nil, err
    }
    
    summaries := make([]RecentQuizSummary, 0, len(records))
    
    for _, record := range records {
        var topic Topic
        var chapter Chapter
        
        db.First(&topic, "id = ?", record.TopicID)
        db.First(&chapter, "id = ?", record.ChapterID)
        
        summaries = append(summaries, RecentQuizSummary{
            ID:             record.ID,
            TopicName:      topic.DisplayName,
            ChapterName:    chapter.DisplayName,
            Score:          record.Score,
            TotalQuestions: record.TotalQuestions,
            CompletedAt:    record.CompletedAt.Format(time.RFC3339),
            // DisplayTime 由前端根据时间差计算
        })
    }
    
    return summaries, nil
}
```

---

## WebSocket Event Data Structures

### ProgressUpdatedEvent (学习进度更新事件)

**触发时机**: 用户完成一个章节的学习

```typescript
interface ProgressUpdatedEvent {
  event: "progress_updated"
  data: {
    user_id: number
    topic_id: string
    chapter_id: string
    completed: boolean
    timestamp: string  // ISO 8601 格式
  }
}
```

**后端发送逻辑**:
```go
func BroadcastProgressUpdate(userID uint, topicID, chapterID string, completed bool) {
    event := map[string]interface{}{
        "event": "progress_updated",
        "data": map[string]interface{}{
            "user_id":    userID,
            "topic_id":   topicID,
            "chapter_id": chapterID,
            "completed":  completed,
            "timestamp":  time.Now().Format(time.RFC3339),
        },
    }
    
    // 发送给该用户的所有 WebSocket 连接
    websocketHub.BroadcastToUser(userID, event)
}
```

---

### QuizCompletedEvent (测验完成事件)

**触发时机**: 用户完成一次测验

```typescript
interface QuizCompletedEvent {
  event: "quiz_completed"
  data: {
    user_id: number
    quiz_id: number
    score: number
    total_questions: number
    timestamp: string  // ISO 8601 格式
  }
}
```

**后端发送逻辑**:
```go
func BroadcastQuizCompletion(userID uint, quizID uint, score, total int) {
    event := map[string]interface{}{
        "event": "quiz_completed",
        "data": map[string]interface{}{
            "user_id":         userID,
            "quiz_id":         quizID,
            "score":           score,
            "total_questions": total,
            "timestamp":       time.Now().Format(time.RFC3339),
        },
    }
    
    websocketHub.BroadcastToUser(userID, event)
}
```

---

## Data Validation Rules

### 学习进度记录
- `user_id`: 必填，必须是有效的用户 ID
- `topic_id`: 必填，必须是有效的主题 ID
- `chapter_id`: 必填，必须是有效的章节 ID
- `completed`: 布尔值，默认 false
- `last_visited_at`: 自动更新为当前时间

### 测验记录
- `user_id`: 必填，必须是有效的用户 ID
- `topic_id`: 必填，必须是有效的主题 ID
- `chapter_id`: 必填，必须是有效的章节 ID
- `score`: 必填，范围 0 到 `total_questions`
- `total_questions`: 必填，大于 0
- `passed`: 根据 `score / total_questions >= 0.6` 自动计算

### WebSocket 事件
- `event`: 必填，枚举值 `progress_updated` 或 `quiz_completed`
- `data.user_id`: 必填，必须是有效的用户 ID
- `data.timestamp`: 必填，ISO 8601 格式

---

## Data Relationships

```
User (1) ──────< (N) LearningProgress
User (1) ──────< (N) QuizRecord

Topic (1) ──────< (N) Chapter
Topic (1) ──────< (N) LearningProgress
Topic (1) ──────< (N) QuizRecord

Chapter (1) ────< (N) LearningProgress
Chapter (1) ────< (N) QuizRecord
```

---

## Performance Considerations

### 索引建议

```sql
-- 学习进度表
CREATE INDEX idx_learning_progress_user_id ON learning_progress(user_id);
CREATE INDEX idx_learning_progress_user_topic ON learning_progress(user_id, topic_id);
CREATE INDEX idx_learning_progress_last_visited ON learning_progress(user_id, last_visited_at DESC);

-- 测验记录表
CREATE INDEX idx_quiz_records_user_id ON quiz_records(user_id);
CREATE INDEX idx_quiz_records_completed_at ON quiz_records(user_id, completed_at DESC);
```

### 查询优化

1. **学习天数计算**: 使用 `COUNT(DISTINCT DATE(...))` 可能较慢，考虑缓存结果
2. **主题进度汇总**: 批量查询所有主题，避免 N+1 查询问题
3. **最近测验记录**: 使用 `LIMIT 5` 限制结果集大小

---

## Summary

Dashboard 功能完全复用现有数据模型，无需新增数据表。主要通过以下方式获取数据：

1. **聚合查询**: 计算学习天数、完成进度、主题进度
2. **排序查询**: 获取最后学习记录、最近测验记录
3. **关联查询**: 获取主题和章节的显示名称
4. **实时推送**: 通过 WebSocket 推送学习进度和测验完成事件

所有数据计算逻辑在后端完成，前端仅负责展示和格式化。
