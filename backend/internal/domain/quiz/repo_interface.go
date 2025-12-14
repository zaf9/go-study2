package quiz

import (
	"context"
	"time"
)

// Repository 是测验持久化接口，供 Service 使用。
type Repository interface {
	SaveRecord(ctx context.Context, record *Record) (int64, error)
	ListRecords(ctx context.Context, userID int64, topic string, from, to *time.Time) ([]Record, error)
}
