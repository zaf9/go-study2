package contract

import (
	"context"
	"fmt"
	"testing"

	"go-study2/internal/app/http_server/handler"
	appquiz "go-study2/internal/app/quiz"
	quizdom "go-study2/internal/domain/quiz"
)

func TestQuizStatsContract(t *testing.T) {
	repo := newQuizMemoryRepo([]quizdom.QuizQuestion{
		{
			ID:             201,
			Topic:          "constants",
			Chapter:        "boolean",
			Type:           quizdom.QuestionTypeSingle,
			Difficulty:     quizdom.DifficultyEasy,
			Question:       "布尔常量的零值是？",
			Options:        `["false","true"]`,
			CorrectAnswers: `["A"]`,
			Explanation:    "示例",
		},
		{
			ID:             202,
			Topic:          "constants",
			Chapter:        "boolean",
			Type:           quizdom.QuestionTypeMultiple,
			Difficulty:     quizdom.DifficultyMedium,
			Question:       "关于布尔的描述哪些正确？",
			Options:        `["只能为 true/false","可以与数字互转"]`,
			CorrectAnswers: `["AB"]`,
			Explanation:    "示例",
		},
	})
	svc := appquiz.NewService(repo)
	h := handler.New()
	setQuizService(h, svc)

	// 直接调用 Service 层的 GetStats 以验证统计逻辑
	stats, err := svc.GetStats(context.Background(), "constants", "boolean")
	if err != nil {
		t.Fatalf("GetStats 出错: %v", err)
	}
	if stats.Total != 2 {
		t.Fatalf("Total 期望 2, 但得到 %d", stats.Total)
	}
	if stats.ByType[quizdom.QuestionTypeSingle] != 1 {
		t.Fatalf("ByType single 期望 1")
	}
}

// 以下为简化的内存仓储与注入工具，避免依赖真实数据库
type quizMemoryRepo struct {
	questions []quizdom.QuizQuestion
	sessions  map[string]*quizdom.QuizSession
}

func newQuizMemoryRepo(questions []quizdom.QuizQuestion) *quizMemoryRepo {
	return &quizMemoryRepo{questions: questions, sessions: map[string]*quizdom.QuizSession{}}
}

func (m *quizMemoryRepo) GetQuestionsByChapter(_ context.Context, topic, chapter string) ([]quizdom.QuizQuestion, error) {
	var result []quizdom.QuizQuestion
	for _, q := range m.questions {
		if q.Topic == topic && q.Chapter == chapter {
			result = append(result, q)
		}
	}
	return result, nil
}

func (m *quizMemoryRepo) CreateSession(_ context.Context, session *quizdom.QuizSession) (string, error) {
	id := fmt.Sprintf("session-%d", len(m.sessions)+1)
	copySession := *session
	copySession.SessionID = id
	m.sessions[id] = &copySession
	return id, nil
}

func (m *quizMemoryRepo) SaveAttempts(_ context.Context, attempts []quizdom.QuizAttempt) error {
	return nil
}

func (m *quizMemoryRepo) GetHistory(_ context.Context, userID int64, topic string, limit int) ([]quizdom.QuizSession, error) {
	return nil, nil
}

func (m *quizMemoryRepo) GetSession(_ context.Context, sessionID string) (*quizdom.QuizSession, error) {
	if s, ok := m.sessions[sessionID]; ok {
		copy := *s
		return &copy, nil
	}
	return nil, nil
}

func (m *quizMemoryRepo) UpdateSessionResult(_ context.Context, sessionID string, correct int, score int, passed bool) error {
	return nil
}

func (m *quizMemoryRepo) GetAttemptsBySession(_ context.Context, sessionID string) ([]quizdom.QuizAttempt, error) {
	return nil, nil
}
