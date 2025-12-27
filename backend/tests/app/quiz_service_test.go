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
		if v.UserID == userID {
			if topic == "" || v.Topic == topic {
				items = append(items, v)
			}
		}
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

func (m *memoryQuizRepo) GetAttemptsBySession(ctx context.Context, sessionID string) ([]quizdom.QuizAttempt, error) {
	var results []quizdom.QuizAttempt
	for _, a := range m.attempts {
		if a.SessionID == sessionID {
			results = append(results, a)
		}
	}
	return results, nil
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

// TestQuizService_GetRecentQuizzes 测试获取最近测验记录功能
func TestQuizService_GetRecentQuizzes(t *testing.T) {
	repo := newMemoryQuizRepo()
	svc := appquiz.NewService(repo)

	// 创建多个已完成的测验会话
	now := time.Now()
	session1 := quizdom.QuizSession{
		ID:             1,
		SessionID:      uuid.NewString(),
		UserID:         1,
		Topic:          "variables",
		Chapter:        "storage",
		TotalQuestions: 5,
		CorrectAnswers: 4,
		Score:          80,
		Passed:         true,
		CompletedAt:    func() *time.Time { t := now.Add(-2 * time.Hour); return &t }(),
	}
	session2 := quizdom.QuizSession{
		ID:             2,
		SessionID:      uuid.NewString(),
		UserID:         1,
		Topic:          "constants",
		Chapter:        "boolean",
		TotalQuestions: 3,
		CorrectAnswers: 2,
		Score:          67,
		Passed:         false,
		CompletedAt:    func() *time.Time { t := now.Add(-1 * time.Hour); return &t }(),
	}
	session3 := quizdom.QuizSession{
		ID:             3,
		SessionID:      uuid.NewString(),
		UserID:         1,
		Topic:          "types",
		Chapter:        "string",
		TotalQuestions: 4,
		CorrectAnswers: 4,
		Score:          100,
		Passed:         true,
		CompletedAt:    func() *time.Time { t := now.Add(-30 * time.Minute); return &t }(),
	}

	// 保存会话到仓库
	repo.CreateSession(ctx(), &session1)
	repo.CreateSession(ctx(), &session2)
	repo.CreateSession(ctx(), &session3)

	// 更新会话结果（标记为已完成）
	repo.UpdateSessionResult(ctx(), session1.SessionID, session1.CorrectAnswers, session1.Score, session1.Passed)
	repo.UpdateSessionResult(ctx(), session2.SessionID, session2.CorrectAnswers, session2.Score, session2.Passed)
	repo.UpdateSessionResult(ctx(), session3.SessionID, session3.CorrectAnswers, session3.Score, session3.Passed)

	// 获取最近测验记录
	recent, err := svc.GetRecentQuizzes(ctx(), 1, 5)
	if err != nil {
		t.Fatalf("获取最近测验记录失败: %v", err)
	}

	if len(recent) == 0 {
		t.Fatal("应返回至少一条测验记录")
	}

	// 验证记录按时间倒序排列（最新的在前）
	if len(recent) >= 2 {
		// 解析时间并比较
		t1, _ := time.Parse(time.RFC3339, recent[0].CompletedAt)
		t2, _ := time.Parse(time.RFC3339, recent[1].CompletedAt)
		if t1.Before(t2) {
			t.Fatalf("记录应按完成时间倒序排列，但 %s 在 %s 之前", recent[0].CompletedAt, recent[1].CompletedAt)
		}
	}

	// 验证记录格式
	for _, r := range recent {
		if r.ID == 0 {
			t.Fatalf("记录 ID 不应为 0: %+v", r)
		}
		if r.TopicName == "" {
			t.Fatalf("主题名称不应为空: %+v", r)
		}
		if r.ChapterName == "" {
			t.Fatalf("章节名称不应为空: %+v", r)
		}
		if r.TotalQuestions == 0 {
			t.Fatalf("总题数不应为 0: %+v", r)
		}
		if r.CompletedAt == "" {
			t.Fatalf("完成时间不应为空: %+v", r)
		}
	}

	// 测试限制数量
	recentLimited, err := svc.GetRecentQuizzes(ctx(), 1, 2)
	if err != nil {
		t.Fatalf("获取限制数量的最近测验记录失败: %v", err)
	}
	if len(recentLimited) > 2 {
		t.Fatalf("应返回最多 2 条记录，但得到 %d 条", len(recentLimited))
	}

	// 测试无记录情况
	recentEmpty, err := svc.GetRecentQuizzes(ctx(), 999, 5)
	if err != nil {
		t.Fatalf("无记录时不应返回错误: %v", err)
	}
	if len(recentEmpty) != 0 {
		t.Fatalf("无记录时应返回空数组，但得到 %d 条", len(recentEmpty))
	}
}
