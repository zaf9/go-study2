package app

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	appquiz "go-study2/internal/app/quiz"
	quizdom "go-study2/internal/domain/quiz"

	"github.com/google/uuid"
)

// memoryQuizRepo 提供内存仓储，便于判分与幂等测试。
type memoryQuizRepo struct {
	questions []quizdom.QuizQuestion
	sessions  map[string]quizdom.QuizSession
	attempts  []quizdom.QuizAttempt
}

func newMemoryQuizRepo() *memoryQuizRepo {
	return &memoryQuizRepo{
		sessions: map[string]quizdom.QuizSession{},
	}
}

func (m *memoryQuizRepo) GetQuestionsByChapter(ctx context.Context, topic, chapter string) ([]quizdom.QuizQuestion, error) {
	return m.questions, nil
}

func (m *memoryQuizRepo) CreateSession(ctx context.Context, session *quizdom.QuizSession) (string, error) {
	if session.SessionID == "" {
		session.SessionID = uuid.NewString()
	}
	m.sessions[session.SessionID] = *session
	return session.SessionID, nil
}

func (m *memoryQuizRepo) SaveAttempts(ctx context.Context, attempts []quizdom.QuizAttempt) error {
	m.attempts = append(m.attempts, attempts...)
	return nil
}

func (m *memoryQuizRepo) GetHistory(ctx context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error) {
	var items []quizdom.QuizSession
	for _, v := range m.sessions {
		items = append(items, v)
	}
	return items, nil
}

func (m *memoryQuizRepo) GetSession(ctx context.Context, sessionID string) (*quizdom.QuizSession, error) {
	if v, ok := m.sessions[sessionID]; ok {
		return &v, nil
	}
	return nil, nil
}

func (m *memoryQuizRepo) UpdateSessionResult(ctx context.Context, sessionID string, correct int, score int, passed bool) error {
	if v, ok := m.sessions[sessionID]; ok {
		v.CorrectAnswers = correct
		v.Score = score
		v.Passed = passed
		now := time.Now()
		v.CompletedAt = &now
		m.sessions[sessionID] = v
	}
	return nil
}

func TestQuizService_FlowAndIdempotency(t *testing.T) {
	repo := newMemoryQuizRepo()
	repo.questions = []quizdom.QuizQuestion{
		{
			ID:             1,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           quizdom.QuestionTypeSingle,
			Difficulty:     quizdom.DifficultyEasy,
			Question:       "单选示例",
			Options:        toJSON([]string{"A1", "A2"}),
			CorrectAnswers: toJSON([]string{"A"}),
			Explanation:    "exp1",
		},
		{
			ID:             2,
			Topic:          "variables",
			Chapter:        "storage",
			Type:           quizdom.QuestionTypeMultiple,
			Difficulty:     quizdom.DifficultyMedium,
			Question:       "多选示例",
			Options:        toJSON([]string{"B1", "B2", "B3"}),
			CorrectAnswers: toJSON([]string{"A", "B"}),
			Explanation:    "exp2",
		},
	}

	svc := appquiz.NewService(repo)
	payload, err := svc.GetQuizQuestions(ctx(), 1, "variables", "storage")
	if err != nil {
		t.Fatalf("获取题目失败: %v", err)
	}
	if payload.SessionID == "" || len(payload.Questions) != 2 {
		t.Fatalf("返回内容不完整: %+v", payload)
	}

	result, err := svc.SubmitQuiz(ctx(), 1, payload.SessionID, "variables", "storage", []appquiz.AnswerSubmission{
		{QuestionID: 1, UserAnswers: []string{"A"}},
		{QuestionID: 2, UserAnswers: []string{"A", "B"}},
	})
	if err != nil {
		t.Fatalf("提交测验失败: %v", err)
	}
	if result.Score <= 0 || result.CorrectAnswers != 2 || !result.Passed {
		t.Fatalf("判分结果异常: %+v", result)
	}

	if _, err := svc.SubmitQuiz(ctx(), 1, payload.SessionID, "variables", "storage", []appquiz.AnswerSubmission{
		{QuestionID: 1, UserAnswers: []string{"A"}},
	}); !errors.Is(err, appquiz.ErrDuplicateSubmit) {
		t.Fatalf("重复提交未被拒绝: %v", err)
	}
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
