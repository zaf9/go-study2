package entity

import "time"

// QuizSession 是测验会话的物理模型（映射数据库表 quiz_sessions）。
type QuizSession struct {
	SessionID      string     `json:"sessionId"      orm:"session_id"      description:"会话ID，UUID"`
	UserID         int64      `json:"userId"         orm:"user_id"         description:"用户ID"`
	Topic          string     `json:"topic"          orm:"topic"           description:"主题标识"`
	Chapter        string     `json:"chapter"        orm:"chapter"         description:"章节标识"`
	TotalQuestions int        `json:"totalQuestions" orm:"total_questions" description:"总题数"`
	CorrectAnswers int        `json:"correctAnswers" orm:"correct_answers" description:"正确题数"`
	Score          int        `json:"score"          orm:"score"           description:"百分制得分 (0-100)"`
	Passed         bool       `json:"passed"         orm:"passed"          description:"是否通过 (score >= 60)"`
	StartedAt      *time.Time `json:"startedAt"      orm:"started_at"      description:"开始时间"`
	CompletedAt    *time.Time `json:"completedAt"    orm:"completed_at"    description:"结束时间"`
	CreatedAt      *time.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`
}
