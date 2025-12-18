package entity

import "time"

// QuizAttempt 是单题作答记录的物理模型（映射数据库表 quiz_attempts）。
type QuizAttempt struct {
	Id          int64      `json:"id"          orm:"id"           description:"自增ID"`
	SessionId   string     `json:"sessionId"   orm:"session_id"   description:"关联的会话ID"`
	QuestionId  string     `json:"questionId"  orm:"question_id"  description:"题目标识"`
	UserChoice  string     `json:"userChoice"  orm:"user_choice"  description:"用户选择的答案内容"`
	IsCorrect   bool       `json:"isCorrect"   orm:"is_correct"   description:"是否正确"`
	AttemptedAt *time.Time `json:"attemptedAt" orm:"attempted_at" description:"作答时间"`
}
