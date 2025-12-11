package progress

import "context"

// ProgressRepository 定义学习进度聚合的完整持久化操作。
type ProgressRepository interface {
	CreateOrUpdate(ctx context.Context, record *LearningProgress) error
	Get(ctx context.Context, userID int64, topic, chapter string) (*LearningProgress, error)
	GetByUser(ctx context.Context, userID int64) ([]LearningProgress, error)
	GetByTopic(ctx context.Context, userID int64, topic string) ([]LearningProgress, error)
}

// Repository 为兼容旧逻辑保留的接口，后续将迁移到 ProgressRepository。
type Repository interface {
	Upsert(ctx context.Context, record *Progress) error
	ListByUser(ctx context.Context, userID int64) ([]Progress, error)
	ListByTopic(ctx context.Context, userID int64, topic string) ([]Progress, error)
}
