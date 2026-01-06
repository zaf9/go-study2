package repository

import (
	"context"
	quizdom "go-study2/internal/domain/quiz"
)

// IQuizSessionRepository 定义了测验会话与答题记录相关的仓储操作接口。
// 注意：题目数据由 YAML 文件管理，不再通过数据库获取。
type IQuizSessionRepository interface {
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

	// GetActiveSession 获取24小时内未提交的活跃会话。
	GetActiveSession(ctx context.Context, userID int64, topic, chapter string) (*quizdom.QuizSession, error)

	// MarkAsSubmitted 标记会话为已提交。
	MarkAsSubmitted(ctx context.Context, sessionID string) error
}
