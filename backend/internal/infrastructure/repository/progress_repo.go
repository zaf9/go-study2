package repository

import (
	"context"
	"strings"
	"time"

	"go-study2/internal/domain/progress"

	"github.com/gogf/gf/v2/database/gdb"
)

// ProgressRepository 使用 GoFrame gdb 持久化学习进度。
type ProgressRepository struct {
	db gdb.DB
}

// NewProgressRepository 创建进度仓储实例。
func NewProgressRepository(db gdb.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// Upsert 写入或更新进度记录。
func (r *ProgressRepository) Upsert(ctx context.Context, record *progress.Progress) error {
	now := time.Now()
	record.LastVisit = now
	_, err := r.db.Exec(ctx, `
INSERT INTO learning_progress (user_id, topic, chapter, status, last_visit, last_position, updated_at)
VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, topic, chapter)
DO UPDATE SET status=excluded.status,
              last_visit=excluded.last_visit,
              last_position=excluded.last_position,
              updated_at=CURRENT_TIMESTAMP
`, record.UserID, record.Topic, record.Chapter, record.Status, record.LastVisit, nullableString(record.LastPosition))
	return err
}

// ListByUser 按用户查询全部进度。
func (r *ProgressRepository) ListByUser(ctx context.Context, userID int64) ([]progress.Progress, error) {
	return r.query(ctx, r.db.Model("learning_progress").Where("user_id", userID))
}

// ListByTopic 查询用户在指定主题下的进度。
func (r *ProgressRepository) ListByTopic(ctx context.Context, userID int64, topic string) ([]progress.Progress, error) {
	return r.query(ctx, r.db.Model("learning_progress").Where("user_id", userID).Where("topic", topic))
}

func (r *ProgressRepository) query(ctx context.Context, model *gdb.Model) ([]progress.Progress, error) {
	records, err := model.OrderDesc("last_visit").All(ctx)
	if err != nil {
		return nil, err
	}
	var items []progress.Progress
	if err := records.Structs(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func nullableString(v string) interface{} {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

