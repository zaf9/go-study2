package repository

import (
	"context"
	"errors"
	"time"

	"go-study2/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/google/uuid"
)

// quizRepository 是 IQuizRepository 的 GoFrame gdb 实现。
type quizRepository struct {
	db gdb.DB
}

// NewQuizRepository 创建一个测验仓储实例。
func NewQuizRepository(db gdb.DB) IQuizRepository {
	return &quizRepository{db: db}
}

// GetQuestionsByChapter 获取指定章节的题目。
func (r *quizRepository) GetQuestionsByChapter(ctx context.Context, topic, chapter string) ([]entity.QuizQuestion, error) {
	var items []entity.QuizQuestion
	err := r.db.Model("quiz_questions").
		Where("topic", topic).
		Where("chapter", chapter).
		OrderAsc("id").
		Scan(&items)
	return items, err
}

// CreateSession 创建一个新的测验会话。
func (r *quizRepository) CreateSession(ctx context.Context, session *entity.QuizSession) (string, error) {
	if session == nil {
		return "", errors.New("session is nil")
	}
	if session.SessionID == "" {
		session.SessionID = uuid.NewString()
	}
	now := time.Now()
	if session.StartedAt == nil {
		session.StartedAt = &now
	}
	if session.CreatedAt == nil {
		session.CreatedAt = &now
	}

	_, err := r.db.Model("quiz_sessions").Data(session).Insert()
	if err != nil {
		return "", err
	}
	return session.SessionID, nil
}

// SaveAttempts 批量保存用户的答题记录（需支持事务）。
func (r *quizRepository) SaveAttempts(ctx context.Context, attempts []entity.QuizAttempt) error {
	if len(attempts) == 0 {
		return nil
	}
	return r.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := time.Now()
		for i := range attempts {
			if attempts[i].AttemptedAt == nil {
				attempts[i].AttemptedAt = &now
			}
		}
		_, err := tx.Model("quiz_attempts").Data(attempts).Insert()
		return err
	})
}

// GetSession 根据会话 ID 获取会话信息。
func (r *quizRepository) GetSession(ctx context.Context, sessionID string) (*entity.QuizSession, error) {
	var sess *entity.QuizSession
	err := r.db.Model("quiz_sessions").Where("session_id", sessionID).Scan(&sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// UpdateSessionResult 更新会话的测试结果。
func (r *quizRepository) UpdateSessionResult(ctx context.Context, sessionID string, correct int, score int, passed bool) error {
	now := time.Now()
	_, err := r.db.Model("quiz_sessions").
		Where("session_id", sessionID).
		Data(gdb.Map{
			"correct_answers": correct,
			"score":           score,
			"passed":          passed,
			"completed_at":    now,
		}).
		Update()
	return err
}

// GetHistory 获取用户的测验历史列表。
func (r *quizRepository) GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]entity.QuizSession, error) {
	var sessions []entity.QuizSession
	m := r.db.Model("quiz_sessions").Where("user_id", userID)
	if topic != "" {
		m = m.Where("topic", topic)
	}
	if limit <= 0 {
		limit = 20
	}
	err := m.OrderDesc("completed_at").Limit(limit).Scan(&sessions)
	return sessions, err
}

// GetAttemptsBySession 获取指定会话的所有答题详情。
func (r *quizRepository) GetAttemptsBySession(ctx context.Context, sessionID string) ([]entity.QuizAttempt, error) {
	var items []entity.QuizAttempt
	err := r.db.Model("quiz_attempts").
		Where("session_id", sessionID).
		OrderAsc("id").
		Scan(&items)
	return items, err
}
