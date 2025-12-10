package repository

import (
	"context"
	"time"

	"go-study2/internal/domain/quiz"

	"github.com/gogf/gf/v2/database/gdb"
)

// QuizRepository 持久化测验记录。
type QuizRepository struct {
	db gdb.DB
}

// NewQuizRepository 创建测验仓储。
func NewQuizRepository(db gdb.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

// SaveRecord 保存测验记录。
func (r *QuizRepository) SaveRecord(ctx context.Context, record *quiz.Record) (int64, error) {
	result, err := r.db.Insert(ctx, "quiz_records", map[string]interface{}{
		"user_id":     record.UserID,
		"topic":       record.Topic,
		"chapter":     nullableString(record.Chapter),
		"score":       record.Score,
		"total":       record.Total,
		"duration_ms": record.DurationMs,
		"answers":     record.Answers,
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	record.ID = id
	return id, nil
}

// ListRecords 查询测验历史。
func (r *QuizRepository) ListRecords(ctx context.Context, userID int64, topic string, from, to *time.Time) ([]quiz.Record, error) {
	model := r.db.Model("quiz_records").Where("user_id", userID)
	if topic != "" {
		model = model.Where("topic", topic)
	}
	if from != nil {
		model = model.WhereGTE("created_at", from)
	}
	if to != nil {
		model = model.WhereLTE("created_at", to)
	}
	records, err := model.OrderDesc("created_at").All(ctx)
	if err != nil {
		return nil, err
	}
	var items []quiz.Record
	if err := records.Structs(&items); err != nil {
		return nil, err
	}
	return items, nil
}

