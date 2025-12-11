package quiz

import (
	"context"
	"time"
)

// QuizRepository 定义测验题目、会话与作答的完整持久化接口。
type QuizRepository interface {
	GetQuestionsByChapter(ctx context.Context, topic, chapter string) ([]QuizQuestion, error)
	CreateSession(ctx context.Context, session *QuizSession) (string, error)
	SaveAttempts(ctx context.Context, attempts []QuizAttempt) error
	GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]QuizSession, error)
}

// Repository 定义测验记录的持久化接口。
type Repository interface {
	SaveRecord(ctx context.Context, record *Record) (int64, error)
	ListRecords(ctx context.Context, userID int64, topic string, from, to *time.Time) ([]Record, error)
}
