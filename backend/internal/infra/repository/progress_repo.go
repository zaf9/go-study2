package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"go-study2/internal/domain/progress"

	"github.com/gogf/gf/v2/database/gdb"
)

// ProgressRepository 使用 GoFrame gdb 实现进度表的持久化。
type ProgressRepository struct {
	db gdb.DB
}

// NewProgressRepository 创建仓储实例。
func NewProgressRepository(db gdb.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// CreateOrUpdate 写入或更新学习进度，累加时长并避免状态回退。
func (r *ProgressRepository) CreateOrUpdate(ctx context.Context, record *progress.LearningProgress) error {
	if record == nil {
		return errors.New("record is nil")
	}
	lastPosition := int64(0)
	if record.LastPosition != "" {
		if parsed, err := strconv.ParseInt(record.LastPosition, 10, 64); err == nil {
			lastPosition = parsed
		}
	}

	_, err := r.db.Exec(ctx, `
INSERT INTO learning_progress (user_id, topic, chapter, status, read_duration, scroll_progress, last_position, quiz_score, quiz_passed, first_visit_at, last_visit_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT(user_id, topic, chapter)
DO UPDATE SET
    read_duration = learning_progress.read_duration + excluded.read_duration,
    scroll_progress = excluded.scroll_progress,
    last_position = excluded.last_position,
    quiz_score = COALESCE(excluded.quiz_score, learning_progress.quiz_score),
    quiz_passed = CASE WHEN excluded.quiz_passed = 1 THEN 1 ELSE learning_progress.quiz_passed END,
    status = CASE
                 WHEN learning_progress.status = 'completed' AND excluded.status != 'completed' THEN learning_progress.status
                 WHEN learning_progress.status = 'tested' AND excluded.status = 'in_progress' THEN learning_progress.status
                 ELSE excluded.status
             END,
    last_visit_at = CURRENT_TIMESTAMP,
    completed_at = CASE
                       WHEN excluded.status = 'completed' AND excluded.quiz_passed = 1 THEN COALESCE(learning_progress.completed_at, CURRENT_TIMESTAMP)
                       ELSE learning_progress.completed_at
                   END,
    updated_at = CURRENT_TIMESTAMP
`, record.UserID, record.Topic, record.Chapter, record.Status, record.ReadDuration, record.ScrollProgress, lastPosition, record.QuizScore, record.QuizPassed)
	return err
}

// Get 返回指定章节的进度。
func (r *ProgressRepository) Get(ctx context.Context, userID int64, topic, chapter string) (*progress.LearningProgress, error) {
	one, err := r.db.Model("learning_progress").
		Where("user_id", userID).
		Where("topic", topic).
		Where("chapter", chapter).
		One(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if one.IsEmpty() {
		return nil, nil
	}
	var item progress.LearningProgress
	if err := one.Struct(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByUser 返回用户的全部进度。
func (r *ProgressRepository) GetByUser(ctx context.Context, userID int64) ([]progress.LearningProgress, error) {
	records, err := r.db.Model("learning_progress").
		Where("user_id", userID).
		OrderDesc("last_visit_at").
		All(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []progress.LearningProgress{}, nil
		}
		return nil, err
	}
	var items []progress.LearningProgress
	if err := records.Structs(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetByTopic 返回指定主题下的进度。
func (r *ProgressRepository) GetByTopic(ctx context.Context, userID int64, topic string) ([]progress.LearningProgress, error) {
	records, err := r.db.Model("learning_progress").
		Where("user_id", userID).
		Where("topic", topic).
		OrderDesc("last_visit_at").
		All(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []progress.LearningProgress{}, nil
		}
		return nil, err
	}
	var items []progress.LearningProgress
	if err := records.Structs(&items); err != nil {
		return nil, err
	}
	return items, nil
}
