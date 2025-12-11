package progress

import "context"

// Repository 定义学习进度的持久化操作。
type Repository interface {
	Upsert(ctx context.Context, record *Progress) error
	ListByUser(ctx context.Context, userID int64) ([]Progress, error)
	ListByTopic(ctx context.Context, userID int64, topic string) ([]Progress, error)
}
