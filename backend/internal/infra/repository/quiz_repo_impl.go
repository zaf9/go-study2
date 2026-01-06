package repository

import (
	"context"
	"errors"
	"time"

	quizdom "go-study2/internal/domain/quiz"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/google/uuid"
)

// quizSessionRepository ? IQuizSessionRepository ? GoFrame gdb ???
type quizSessionRepository struct {
	db gdb.DB
}

// NewQuizRepository ?????????????
func NewQuizRepository(db gdb.DB) IQuizSessionRepository {
	return &quizSessionRepository{db: db}
}

// CreateSession ???????????
func (r *quizSessionRepository) CreateSession(ctx context.Context, session *quizdom.QuizSession) (string, error) {
	if session == nil {
		return "", errors.New("session is nil")
	}
	if session.SessionID == "" {
		session.SessionID = uuid.NewString()
	}
	now := time.Now()
	if session.StartedAt.IsZero() {
		session.StartedAt = now
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}

	_, err := r.db.Model("quiz_sessions").Data(session).FieldsEx("id").Insert()
	if err != nil {
		return "", err
	}
	return session.SessionID, nil
}

// SaveAttempts ???????????????????
func (r *quizSessionRepository) SaveAttempts(ctx context.Context, attempts []quizdom.QuizAttempt) error {
	if len(attempts) == 0 {
		return nil
	}
	return r.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := time.Now()
		for i := range attempts {
			if attempts[i].AttemptedAt.IsZero() {
				attempts[i].AttemptedAt = now
			}
		}
		_, err := tx.Model("quiz_attempts").Data(attempts).FieldsEx("id").Insert()
		return err
	})
}

// GetSession ???? ID ???????
func (r *quizSessionRepository) GetSession(ctx context.Context, sessionID string) (*quizdom.QuizSession, error) {
	var sess *quizdom.QuizSession
	err := r.db.Model("quiz_sessions").Where("session_id", sessionID).Scan(&sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// UpdateSessionResult ??????????
func (r *quizSessionRepository) UpdateSessionResult(ctx context.Context, sessionID string, correct int, score int, passed bool) error {
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

// GetHistory ????????????
func (r *quizSessionRepository) GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error) {
	var sessions []quizdom.QuizSession
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

// GetAttemptsBySession ??????????????
func (r *quizSessionRepository) GetAttemptsBySession(ctx context.Context, sessionID string) ([]quizdom.QuizAttempt, error) {
	var items []quizdom.QuizAttempt
	err := r.db.Model("quiz_attempts").
		Where("session_id", sessionID).
		OrderAsc("id").
		Scan(&items)
	return items, err
}

// GetActiveSession ??24????????????
func (r *quizSessionRepository) GetActiveSession(ctx context.Context, userID int64, topic, chapter string) (*quizdom.QuizSession, error) {
	var sess *quizdom.QuizSession
	// ??24???????????topic?chapter?????
	err := r.db.Model("quiz_sessions").
		Fields("id,session_id,user_id,topic,chapter,total_questions,correct_answers,score,passed,started_at,completed_at,submitted_at,created_at").
		Where("user_id", userID).
		Where("topic", topic).
		Where("chapter", chapter).
		Where("submitted_at IS NULL").
		Where("started_at >= ?", time.Now().Add(-24*time.Hour)).
		OrderDesc("started_at").
		Limit(1).
		Scan(&sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// MarkAsSubmitted ?????????
func (r *quizSessionRepository) MarkAsSubmitted(ctx context.Context, sessionID string) error {
	now := time.Now()
	_, err := r.db.Model("quiz_sessions").
		Where("session_id", sessionID).
		Where("submitted_at IS NULL"). // ??????????????
		Data(gdb.Map{"submitted_at": now}).
		Update()
	return err
}
