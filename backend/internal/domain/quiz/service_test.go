package quiz

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockQuizRepo struct {
	saved   []*Record
	failErr error
}

func (m *mockQuizRepo) SaveRecord(_ context.Context, record *Record) (int64, error) {
	if m.failErr != nil {
		return 0, m.failErr
	}
	record.ID = int64(len(m.saved) + 1)
	record.CreatedAt = time.Now()
	m.saved = append(m.saved, record)
	return record.ID, nil
}

func (m *mockQuizRepo) ListRecords(_ context.Context, userID int64, topic string, _ *time.Time, _ *time.Time) ([]Record, error) {
	if m.failErr != nil {
		return nil, m.failErr
	}
	var res []Record
	for _, rec := range m.saved {
		if rec.UserID == userID && (topic == "" || rec.Topic == topic) {
			res = append(res, *rec)
		}
	}
	return res, nil
}

func TestService_GetQuestionsAndSubmit(t *testing.T) {
	repo := &mockQuizRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	_, err := svc.GetQuestions(ctx, "", "")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("应返回 ErrInvalidInput")
	}

	qs, err := svc.GetQuestions(ctx, "variables", "storage")
	if err != nil {
		t.Fatalf("获取题目失败: %v", err)
	}
	if len(qs) == 0 {
		t.Fatalf("变量主题题目应大于 0")
	}

	answer := SubmitAnswer{
		ID:      qs[0].ID,
		Choices: qs[0].Answer,
	}
	result, err := svc.Submit(ctx, 1, "variables", "storage", []SubmitAnswer{answer}, 1500)
	if err != nil {
		t.Fatalf("提交测验失败: %v", err)
	}
	if result.Total == 0 {
		t.Fatalf("总题数不应为 0")
	}
	if len(repo.saved) != 1 {
		t.Fatalf("记录未被保存")
	}

	history, err := svc.History(ctx, 1, "variables", nil, nil)
	if err != nil {
		t.Fatalf("查询历史失败: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("历史记录数量不正确")
	}
}

