package repository

import (
	"context"
	"errors"
	"time"

	"go-study2/internal/domain/quiz"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/google/uuid"
)

// QuizRepository 使用 GoFrame gdb 实现测验题库、会话与作答的持久化。
type QuizRepository struct {
	db gdb.DB
}

// NewQuizRepository 创建测验仓储。
func NewQuizRepository(db gdb.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

// GetQuestionsByChapter 返回指定章节的题目列表。
func (r *QuizRepository) GetQuestionsByChapter(ctx context.Context, topic, chapter string) ([]quiz.QuizQuestion, error) {
	records, err := r.db.Model("quiz_questions").
		Where("topic", topic).
		Where("chapter", chapter).
		OrderAsc("id").
		All(ctx)
	if err != nil {
		return nil, err
	}
	var items []quiz.QuizQuestion
	if err := records.Structs(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// CreateSession 创建测验会话并返回 sessionID。
func (r *QuizRepository) CreateSession(ctx context.Context, session *quiz.QuizSession) (string, error) {
	if session == nil {
		return "", errors.New("session is nil")
	}
	if session.SessionID == "" {
		session.SessionID = uuid.NewString()
	}
	if session.StartedAt.IsZero() {
		session.StartedAt = time.Now()
	}

	_, err := r.db.Insert(ctx, "quiz_sessions", gdb.Map{
		"session_id":      session.SessionID,
		"user_id":         session.UserID,
		"topic":           session.Topic,
		"chapter":         session.Chapter,
		"total_questions": session.TotalQuestions,
		"correct_answers": session.CorrectAnswers,
		"score":           session.Score,
		"passed":          session.Passed,
		"started_at":      session.StartedAt,
		"completed_at":    session.CompletedAt,
		"created_at":      time.Now(),
	})
	if err != nil {
		return "", err
	}
	return session.SessionID, nil
}

// SaveAttempts 批量写入测验作答记录。
func (r *QuizRepository) SaveAttempts(ctx context.Context, attempts []quiz.QuizAttempt) error {
	if len(attempts) == 0 {
		return nil
	}
	return r.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, attempt := range attempts {
			if _, err := tx.Insert("quiz_attempts", gdb.Map{
				"session_id":   attempt.SessionID,
				"user_id":      attempt.UserID,
				"topic":        attempt.Topic,
				"chapter":      attempt.Chapter,
				"question_id":  attempt.QuestionID,
				"user_answers": attempt.UserAnswers,
				"is_correct":   attempt.IsCorrect,
				"attempted_at": func() time.Time {
					if attempt.AttemptedAt.IsZero() {
						return time.Now()
					}
					return attempt.AttemptedAt
				}(),
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetSession 根据 sessionID 查询测验会话。
func (r *QuizRepository) GetSession(ctx context.Context, sessionID string) (*quiz.QuizSession, error) {
	record, err := r.db.Model("quiz_sessions").Where("session_id", sessionID).One(ctx)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, nil
	}
	var sess quiz.QuizSession
	if err := record.Struct(&sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

// UpdateSessionResult 写入测验得分与完成时间。
func (r *QuizRepository) UpdateSessionResult(ctx context.Context, sessionID string, correct int, score int, passed bool) error {
	_, err := r.db.Model("quiz_sessions").
		Where("session_id", sessionID).
		Data(gdb.Map{
			"correct_answers": correct,
			"score":           score,
			"passed":          passed,
			"completed_at":    time.Now(),
		}).
		Update()
	return err
}

// GetHistory 返回用户的测验会话列表。
func (r *QuizRepository) GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]quiz.QuizSession, error) {
	model := r.db.Model("quiz_sessions").Where("user_id", userID)
	if topic != "" {
		model = model.Where("topic", topic)
	}
	if limit <= 0 {
		limit = 10
	}

	records, err := model.OrderDesc("completed_at").Limit(limit).All(ctx)
	if err != nil {
		return nil, err
	}
	var sessions []quiz.QuizSession
	if err := records.Structs(&sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}
