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

// 学习进度状态枚举。
const (
	StatusNotStarted = "not_started"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusTested     = "tested"
	StatusDone       = StatusCompleted
)

var (
	// SupportedTopics 定义允许的主题集合。
	SupportedTopics = map[string]struct{}{
		"lexical_elements": {},
		"constants":        {},
		"variables":        {},
		"types":            {},
		"properties":       {},
	}
	// ValidStatuses 定义允许的进度状态集合。
	ValidStatuses = map[string]struct{}{
		StatusNotStarted: {},
		StatusInProgress: {},
		StatusCompleted:  {},
		StatusTested:     {},
	}
)

// IsSupportedTopic 判断主题是否在允许范围内。
func IsSupportedTopic(topic string) bool {
	_, ok := SupportedTopics[topic]
	return ok
}

// IsValidStatus 判断状态是否有效。
func IsValidStatus(status string) bool {
	_, ok := ValidStatuses[status]
	return ok
}

// IsCompleted 判断章节是否已完成学习
func (lp *LearningProgress) IsCompleted() bool {
	return lp.Status == StatusCompleted && lp.CompletedAt != nil
}

// MarkAsInProgress 标记章节为学习中状态
func (lp *LearningProgress) MarkAsInProgress() {
	if lp.Status == StatusNotStarted {
		lp.Status = StatusInProgress
		now := time.Now()
		if lp.FirstVisitAt.IsZero() {
			lp.FirstVisitAt = now
		}
		lp.LastVisitAt = now
	}
}

// Complete 标记章节为已完成状态
func (lp *LearningProgress) Complete() {
	lp.Status = StatusCompleted
	now := time.Now()
	lp.CompletedAt = &now
	lp.LastVisitAt = now
}

// ChapterStatusInfo 章节状态信息
type ChapterStatusInfo struct {
	Chapter     string     `json:"chapter"`
	Status      string     `json:"status"`
	QuizScore   int        `json:"quizScore"`
	QuizPassed  bool       `json:"quizPassed"`
	LastVisitAt *time.Time `json:"lastVisitAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

// TopicProgressSummary 主题进度汇总
type TopicProgressSummary struct {
	Topic              string  `json:"topic"`
	TotalChapters      int     `json:"totalChapters"`
	CompletedChapters  int     `json:"completedChapters"`
	InProgressChapters int     `json:"inProgressChapters"`
	Weight             int     `json:"weight"`   // 主题权重
	Progress           float64 `json:"progress"` // 主题进度百分比
}

// ProgressOverview 进度概览
type ProgressOverview struct {
	TotalChapters      int                    `json:"totalChapters"`
	CompletedChapters  int                    `json:"completedChapters"`
	InProgressChapters int                    `json:"inProgressChapters"`
	CompletionRate     float64                `json:"completionRate"`
	StudyDays          int                    `json:"studyDays"`      // 学习天数
	TotalStudyTime     int64                  `json:"totalStudyTime"` // 学习总时长（秒）
	Topics             []TopicProgressSummary `json:"topics"`
	NextChapter        *NextChapterHint       `json:"next,omitempty"`
}

// NextChapterHint 下一个建议学习的章节
type NextChapterHint struct {
	Topic   string `json:"topic"`
	Chapter string `json:"chapter"`
	Title   string `json:"title"`
}

// TopicProgressDetail 主题进度详情
type TopicProgressDetail struct {
	Topic         string              `json:"topic"`
	TotalChapters int                 `json:"totalChapters"`
	Chapters      []ChapterStatusInfo `json:"chapters"`
}
