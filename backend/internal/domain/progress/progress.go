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
