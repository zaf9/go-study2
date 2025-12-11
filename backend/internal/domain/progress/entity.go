package progress

import "time"

// Progress 表示用户在某章节的学习状态。
type Progress struct {
	ID           int64     `json:"id" orm:"id"`
	UserID       int64     `json:"userId" orm:"user_id"`
	Topic        string    `json:"topic" orm:"topic"`
	Chapter      string    `json:"chapter" orm:"chapter"`
	Status       string    `json:"status" orm:"status"`
	LastVisit    time.Time `json:"lastVisit" orm:"last_visit"`
	LastPosition string    `json:"lastPosition" orm:"last_position"`
	CreatedAt    time.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" orm:"updated_at"`
}

const (
	StatusNotStarted = "not_started"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

var (
	supportedTopics = map[string]struct{}{
		"lexical_elements": {},
		"constants":        {},
		"variables":        {},
		"types":            {},
	}
	validStatuses = map[string]struct{}{
		StatusNotStarted: {},
		StatusInProgress: {},
		StatusDone:       {},
	}
)

// IsSupportedTopic 判断 topic 是否在允许范围内。
func IsSupportedTopic(topic string) bool {
	_, ok := supportedTopics[topic]
	return ok
}

// IsValidStatus 判断进度状态是否有效。
func IsValidStatus(status string) bool {
	_, ok := validStatuses[status]
	return ok
}
