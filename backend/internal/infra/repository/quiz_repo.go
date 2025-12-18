package repository

import (
	"context"
	quizdom "go-study2/internal/domain/quiz"
)

// IQuizRepository 定义了测验相关的仓储操作接口。
type IQuizRepository interface {
	// GetQuestionsByChapter 获取指定章节的题目。
	GetQuestionsByChapter(ctx context.Context, topic, chapter string) ([]quizdom.QuizQuestion, error)

	// CreateSession 创建一个新的测验会话。
	CreateSession(ctx context.Context, session *quizdom.QuizSession) (string, error)

	// SaveAttempts 批量保存用户的答题记录（需支持事务）。
	SaveAttempts(ctx context.Context, attempts []quizdom.QuizAttempt) error

	// GetSession 根据会话 ID 获取会话信息。
	GetSession(ctx context.Context, sessionID string) (*quizdom.QuizSession, error)

	// UpdateSessionResult 更新会话的测试结果。
	UpdateSessionResult(ctx context.Context, sessionID string, correct int, score int, passed bool) error

	// GetHistory 获取用户的测验历史列表。
	GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error)

	// GetAttemptsBySession 获取指定会话的所有答题详情。
	GetAttemptsBySession(ctx context.Context, sessionID string) ([]quizdom.QuizAttempt, error)
}
