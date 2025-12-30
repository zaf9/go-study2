# Phase 1: 数据模型设计

**Feature**: 学习进度与测验体验优化  
**Created**: 2025-12-30

## 概述

本文档定义学习进度和测验功能的数据模型,包括数据库表结构、领域实体和API数据传输对象(DTO)。

## 数据库表结构

### 1. learning_progress (学习进度表)

**用途**: 记录用户对每个章节的学习状态和进度信息

**表结构**: (已存在,本次优化无需修改schema)

```sql
CREATE TABLE IF NOT EXISTS learning_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'not_started',
    read_duration INTEGER DEFAULT 0,
    scroll_progress INTEGER DEFAULT 0,
    last_position TEXT DEFAULT '',
    quiz_score INTEGER DEFAULT 0,
    quiz_passed BOOLEAN DEFAULT 0,
    first_visit_at DATETIME,
    last_visit_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, topic, chapter)
);

CREATE INDEX IF NOT EXISTS idx_learning_progress_user ON learning_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_progress_topic ON learning_progress(user_id, topic);
```

**字段说明**:
- `status`: 章节状态,取值: `not_started`(未开始), `in_progress`(学习中), `completed`(已完成)
- `quiz_passed`: 是否通过测验(布尔值)
- `completed_at`: 完成时间,仅当`quiz_passed=true`时设置

**状态转换逻辑**:
1. 初始状态: 无记录(隐式`not_started`)
2. 用户访问章节 → 创建记录,`status=in_progress`
3. 完成测验且通过 → `status=completed`, `quiz_passed=true`, 设置`completed_at`

### 2. quiz_sessions (测验会话表)

**用途**: 记录每次测验的会话信息

**表结构**: (已存在,本次添加submitted_at字段用于防重复提交)

```sql
CREATE TABLE IF NOT EXISTS quiz_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    topic TEXT NOT NULL,
    chapter TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at DATETIME,  -- 新增: 提交时间,用于幂等性检查
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_quiz_sessions_user ON quiz_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_session_id ON quiz_sessions(session_id);
```

**关键字段**:
- `session_id`: 唯一会话标识符(MD5哈希)
- `submitted_at`: 提交时间,为NULL表示未提交

### 3. quiz_attempts (答题记录表)

**用途**: 记录用户对每道题的作答详情

**表结构**: (已存在,无需修改)

```sql
CREATE TABLE IF NOT EXISTS quiz_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    question_id TEXT NOT NULL,
    user_choice TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    attempted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES quiz_sessions(session_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_quiz_attempts_session ON quiz_attempts(session_id);
```

---

## 领域实体

### 1. LearningProgress (章节学习进度)

**文件**: `backend/internal/domain/progress/progress.go`

```go
package progress

import "time"

// LearningProgress 描述用户在章节层面的学习进度与测验结果。
type LearningProgress struct {
	ID             int64      `json:"id" orm:"id"`
	UserID         int64      `json:"userId" orm:"user_id"`
	Topic          string     `json:"topic" orm:"topic"`
	Chapter        string     `json:"chapter" orm:"chapter"`
	Status         string     `json:"status" orm:"status"`
	ReadDuration   int64      `json:"readDuration" orm:"read_duration"`
	ScrollProgress int        `json:"scrollProgress" orm:"scroll_progress"`
	LastPosition   string     `json:"lastPosition" orm:"last_position"`
	QuizScore      int        `json:"quizScore" orm:"quiz_score"`
	QuizPassed     bool       `json:"quizPassed" orm:"quiz_passed"`
	FirstVisitAt   time.Time  `json:"firstVisitAt" orm:"first_visit_at"`
	LastVisitAt    time.Time  `json:"lastVisitAt" orm:"last_visit_at"`
	CompletedAt    *time.Time `json:"completedAt" orm:"completed_at"`
	CreatedAt      time.Time  `json:"createdAt" orm:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" orm:"updated_at"`
}

// 学习进度状态枚举
const (
	StatusNotStarted = "not_started"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)
```

**业务方法**:
```go
// IsCompleted 判断章节是否已完成
func (p *LearningProgress) IsCompleted() bool {
	return p.Status == StatusCompleted && p.QuizPassed
}

// MarkAsInProgress 标记为学习中
func (p *LearningProgress) MarkAsInProgress() {
	if p.Status == StatusNotStarted {
		p.Status = StatusInProgress
		now := time.Now()
		p.FirstVisitAt = now
		p.LastVisitAt = now
	}
}

// Complete 标记为已完成
func (p *LearningProgress) Complete(score int, passed bool) {
	p.QuizScore = score
	p.QuizPassed = passed
	if passed {
		p.Status = StatusCompleted
		now := time.Now()
		p.CompletedAt = &now
	}
}
```

### 2. QuizSession (测验会话)

**文件**: `backend/internal/domain/quiz/quiz.go`

```go
package quiz

import "time"

// QuizSession 代表一次测验会话
type QuizSession struct {
	ID          int64        `json:"id" orm:"id"`
	SessionID   string       `json:"sessionId" orm:"session_id"`
	UserID      int64        `json:"userId" orm:"user_id"`
	Topic       string       `json:"topic" orm:"topic"`
	Chapter     string       `json:"chapter" orm:"chapter"`
	Questions   []QuizQuestion `json:"questions" orm:"-"` // 不持久化到数据库
	CreatedAt   time.Time    `json:"createdAt" orm:"created_at"`
	SubmittedAt *time.Time   `json:"submittedAt" orm:"submitted_at"`
}

// QuizQuestion 测验题目
type QuizQuestion struct {
	ID          int      `json:"id"`
	Type        string   `json:"type"`        // "single" 或 "multiple"
	Difficulty  string   `json:"difficulty"`  // "easy", "medium", "hard"
	Question    string   `json:"question"`
	Options     []Option `json:"options"`
	Answer      []string `json:"-"`           // 不返回给前端
	CodeSnippet *string  `json:"codeSnippet,omitempty"`
}

// Option 选项
type Option struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// IsSubmitted 判断是否已提交
func (s *QuizSession) IsSubmitted() bool {
	return s.SubmittedAt != nil
}
```

### 3. QuizAttempt (答题记录)

**文件**: `backend/internal/model/entity/quiz_attempt.go` (已存在)

```go
package entity

import "time"

// QuizAttempt 是单题作答记录的物理模型
type QuizAttempt struct {
	Id          int64      `json:"id"          orm:"id"`
	SessionId   string     `json:"sessionId"   orm:"session_id"`
	QuestionId  string     `json:"questionId"  orm:"question_id"`
	UserChoice  string     `json:"userChoice"  orm:"user_choice"`
	IsCorrect   bool       `json:"isCorrect"   orm:"is_correct"`
	AttemptedAt *time.Time `json:"attemptedAt" orm:"attempted_at"`
}
```

---

## API 数据传输对象 (DTO)

### 1. 进度查询响应

**GET /api/v1/progress/overview**

```typescript
interface ProgressOverviewResponse {
  code: number;
  message: string;
  data: {
    totalChapters: number;       // 所有主题的章节总数
    completedChapters: number;   // 已完成的章节数
    inProgressChapters: number;  // 学习中的章节数
    completionRate: number;      // 完成率 (0-100)
    topics: Array<{
      topic: string;
      totalChapters: number;
      completedChapters: number;
      inProgressChapters: number;
    }>;
    next: {                      // 下一个建议学习的章节
      topic: string;
      chapter: string;
      title: string;
    } | null;
  };
}
```

**GET /api/v1/progress/topic/:topic**

```typescript
interface TopicProgressResponse {
  code: number;
  message: string;
  data: {
    topic: string;
    totalChapters: number;
    chapters: Array<{
      chapter: string;
      status: 'not_started' | 'in_progress' | 'completed';
      quizScore: number;
      quizPassed: boolean;
      lastVisitAt: string | null;  // ISO 8601
      completedAt: string | null;
    }>;
  };
}
```

### 2. 测验会话响应

**GET /api/v1/quiz/:topic/:chapter**

```typescript
interface QuizSessionResponse {
  code: number;
  message: string;
  data: {
    sessionId: string;
    topic: string;
    chapter: string;
    questions: Array<{
      id: number;
      type: 'single' | 'multiple';
      difficulty: 'easy' | 'medium' | 'hard';
      question: string;
      options: Array<{
        id: string;
        label: string;
      }>;
      codeSnippet?: string | null;
    }>;
  };
}
```

### 3. 测验提交请求/响应

**POST /api/v1/quiz/submit**

请求体:
```typescript
interface QuizSubmitRequest {
  sessionId: string;
  topic: string;
  chapter: string;
  durationMs?: number;  // 答题耗时(毫秒)
  answers: Array<{
    questionId: number;
    userAnswers: string[];  // 多选题可能有多个答案
  }>;
}
```

响应:
```typescript
interface QuizSubmitResponse {
  code: number;
  message: string;
  data: {
    sessionId: string;
    score: number;           // 总分
    total: number;           // 总题数
    correctCount: number;    // 答对题数
    passed: boolean;         // 是否通过(score >= 60%)
    details: Array<{
      questionId: number;
      question: string;
      userAnswers: string[];
      correctAnswers: string[];
      isCorrect: boolean;
    }>;
  };
}
```

### 4. 测验历史查询

**GET /api/v1/quiz/history?topic=constants**

```typescript
interface QuizHistoryResponse {
  code: number;
  message: string;
  data: Array<{
    sessionId: string;
    topic: string;
    chapter: string;
    score: number;
    total: number;
    passed: boolean;
    durationMs: number;
    createdAt: string;      // ISO 8601
  }>;
}
```

---

## 前端类型定义

**文件**: `frontend/types/learning.ts`

```typescript
export type ChapterStatus = 'not_started' | 'in_progress' | 'completed';

export interface ChapterProgress {
  chapter: string;
  status: ChapterStatus;
  quizScore: number;
  quizPassed: boolean;
  lastVisitAt: string | null;
  completedAt: string | null;
}

export interface TopicProgressDetail {
  topic: string;
  totalChapters: number;
  chapters: ChapterProgress[];
}

export interface ProgressSnapshot {
  totalChapters: number;
  completedChapters: number;
  inProgressChapters: number;
  completionRate: number;
  topics: Array<{
    topic: string;
    totalChapters: number;
    completedChapters: number;
    inProgressChapters: number;
  }>;
  next: NextChapterHint | null;
}

export interface NextChapterHint {
  topic: string;
  chapter: string;
  title: string;
}
```

**文件**: `frontend/types/quiz.ts`

```typescript
export interface QuizQuestion {
  id: number;
  type: 'single' | 'multiple';
  difficulty: 'easy' | 'medium' | 'hard';
  question: string;
  options: Array<{
    id: string;
    label: string;
  }>;
  codeSnippet?: string | null;
}

export interface QuizSessionPayload {
  sessionId: string;
  topic: string;
  chapter: string;
  questions: QuizQuestion[];
}

export interface QuizSubmitResult {
  sessionId: string;
  score: number;
  total: number;
  correctCount: number;
  passed: boolean;
  details: Array<{
    questionId: number;
    question: string;
    userAnswers: string[];
    correctAnswers: string[];
    isCorrect: boolean;
  }>;
}

export interface QuizHistoryItem {
  sessionId: string;
  topic: string;
  chapter: string;
  score: number;
  total: number;
  passed: boolean;
  durationMs: number;
  createdAt: string;
}
```

---

## 数据关系图

```
users (用户表)
  ├─1:N─> learning_progress (学习进度)
  │         ├─ topic + chapter (唯一约束)
  │         └─ status: not_started | in_progress | completed
  │
  └─1:N─> quiz_sessions (测验会话)
            ├─ session_id (唯一)
            └─1:N─> quiz_attempts (答题记录)
                      ├─ question_id
                      ├─ user_choice
                      └─ is_correct
```

**关键约束**:
1. `learning_progress`: (user_id, topic, chapter) 唯一
2. `quiz_sessions`: session_id 全局唯一
3. `quiz_attempts`: 外键关联 quiz_sessions.session_id,级联删除

---

## 数据迁移

### 迁移脚本

**文件**: `backend/internal/infra/migrations/016_add_submitted_at.sql`

```sql
-- 为quiz_sessions表添加submitted_at字段,用于防重复提交
ALTER TABLE quiz_sessions ADD COLUMN submitted_at DATETIME;

-- 为现有已提交的会话设置submitted_at(根据quiz_attempts推断)
UPDATE quiz_sessions
SET submitted_at = (
    SELECT MIN(attempted_at)
    FROM quiz_attempts
    WHERE quiz_attempts.session_id = quiz_sessions.session_id
)
WHERE id IN (
    SELECT DISTINCT qs.id
    FROM quiz_sessions qs
    INNER JOIN quiz_attempts qa ON qs.session_id = qa.session_id
);
```

### 回滚脚本

```sql
-- 如需回滚,删除submitted_at列
-- SQLite不支持直接DROP COLUMN,需重建表
-- 生产环境谨慎操作,建议备份后执行
```

---

## 数据验证规则

### 后端验证

1. **Topic 和 Chapter 枚举检查**:
   ```go
   var supportedTopics = []string{"lexical_elements", "constants", "variables", "types"}
   
   func validateTopic(topic string) error {
       if !contains(supportedTopics, topic) {
           return fmt.Errorf("不支持的主题: %s", topic)
       }
       return nil
   }
   ```

2. **Status 状态检查**:
   ```go
   var validStatuses = []string{"not_started", "in_progress", "completed"}
   
   func validateStatus(status string) error {
       if !contains(validStatuses, status) {
           return fmt.Errorf("无效的状态: %s", status)
       }
       return nil
   }
   ```

3. **Quiz 答案格式验证**:
   ```go
   func validateAnswers(answers []UserAnswer, questions []QuizQuestion) error {
       if len(answers) == 0 {
           return errors.New("答案不能为空")
       }
       
       for _, ans := range answers {
           if ans.QuestionID <= 0 {
               return fmt.Errorf("无效的题目ID: %d", ans.QuestionID)
           }
           if len(ans.UserAnswers) == 0 {
               return fmt.Errorf("题目 %d 的答案不能为空", ans.QuestionID)
           }
       }
       
       return nil
   }
   ```

### 前端验证

1. **提交前检查**:
   ```typescript
   function validateBeforeSubmit(answers: Record<number, string[]>, totalQuestions: number): string | null {
     const answeredCount = Object.keys(answers).length;
     
     if (answeredCount === 0) {
       return '请至少回答一道题';
     }
     
     if (answeredCount < totalQuestions) {
       return `您还有 ${totalQuestions - answeredCount} 道题未作答,确定提交吗?`;
     }
     
     return null;
   }
   ```

---

## 性能优化

### 数据库索引

已有索引:
- `idx_learning_progress_user`: (user_id)
- `idx_learning_progress_topic`: (user_id, topic)
- `idx_quiz_sessions_user`: (user_id)
- `idx_quiz_sessions_session_id`: (session_id)
- `idx_quiz_attempts_session`: (session_id)

### 查询优化

1. **批量查询章节进度**:
   ```go
   // 一次查询获取某主题的所有章节进度,避免N+1问题
   func (r *Repository) GetTopicProgress(ctx context.Context, userID int64, topic string) ([]LearningProgress, error) {
       var list []LearningProgress
       err := r.db.Model("learning_progress").
           Where("user_id", userID).
           Where("topic", topic).
           Scan(&list)
       return list, err
   }
   ```

2. **使用事务提交测验**:
   ```go
   // 确保quiz_sessions更新和quiz_attempts插入的原子性
   func (r *Repository) SubmitQuizWithTransaction(ctx context.Context, sessionID string, attempts []QuizAttempt) error {
       return r.db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
           // 更新session
           // 插入attempts
           // 更新learning_progress
           return nil
       })
   }
   ```

---

## 总结

本数据模型设计遵循以下原则:

1. **最小化变更**: 复用现有表结构,仅添加submitted_at字段
2. **明确约束**: 使用唯一索引和外键保证数据一致性
3. **状态明确**: 三状态(not_started/in_progress/completed)清晰定义
4. **防重复提交**: 通过submitted_at字段实现幂等性
5. **性能优先**: 合理使用索引,批量查询避免N+1

所有实体和DTO定义完整,为Phase 2实现提供明确规范。
