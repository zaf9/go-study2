package quiz

import "time"

// Option 表示题目选项。
type Option struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// Question 表示单个测验题目（用于业务层）。
type Question struct {
	ID          string   `json:"id"`
	Stem        string   `json:"stem"`
	Options     []Option `json:"options"`
	Multi       bool     `json:"multi"`
	Answer      []string `json:"answer"`
	Explanation string   `json:"explanation,omitempty"`
}

// SubmitAnswer 表示用户提交的答案。
type SubmitAnswer struct {
	ID      string   `json:"id"`
	Choices []string `json:"choices"`
}

// Result 表示测验评分结果。
type Result struct {
	Score       int       `json:"score"`
	Total       int       `json:"total"`
	CorrectIDs  []string  `json:"correctIds"`
	WrongIDs    []string  `json:"wrongIds"`
	SubmittedAt time.Time `json:"submittedAt"`
	DurationMs  int64     `json:"durationMs"`
}

// HistoryItem 表示测验历史记录。
type HistoryItem struct {
	ID         int64     `json:"id"`
	Topic      string    `json:"topic"`
	Chapter    string    `json:"chapter,omitempty"`
	Score      int       `json:"score"`
	Total      int       `json:"total"`
	DurationMs int64     `json:"durationMs"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Record 表示持久化的测验记录。
type Record struct {
	ID         int64     `json:"id" orm:"id"`
	UserID     int64     `json:"userId" orm:"user_id"`
	Topic      string    `json:"topic" orm:"topic"`
	Chapter    string    `json:"chapter" orm:"chapter"`
	Score      int       `json:"score" orm:"score"`
	Total      int       `json:"total" orm:"total"`
	DurationMs int64     `json:"durationMs" orm:"duration_ms"`
	Answers    string    `json:"answers" orm:"answers"`
	CreatedAt  time.Time `json:"createdAt" orm:"created_at"`
}

var supportedTopics = map[string]struct{}{
	"lexical_elements": {},
	"constants":        {},
	"variables":        {},
	"types":            {},
}

// IsSupportedTopic 判断 topic 是否在允许范围内。
func IsSupportedTopic(topic string) bool {
	_, ok := supportedTopics[topic]
	return ok
}
