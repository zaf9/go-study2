package quiz

import (
	"context"
	"time"
)

// Repository 定义测验记录的持久化接口。
type Repository interface {
	SaveRecord(ctx context.Context, record *Record) (int64, error)
	ListRecords(ctx context.Context, userID int64, topic string, from, to *time.Time) ([]Record, error)
}
