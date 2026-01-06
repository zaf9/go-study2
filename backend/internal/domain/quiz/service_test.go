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

	// 测试输入验证
	_, err := svc.GetQuestions(ctx, "", "")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("应返回 ErrInvalidInput")
	}

	// 测试 GetQuestions 已废弃，应返回 ErrQuizUnavailable
	qs, err := svc.GetQuestions(ctx, "variables", "storage")
	if err == nil {
		t.Fatalf("GetQuestions 已废弃，应返回错误")
	}
	if !errors.Is(err, ErrQuizUnavailable) {
		t.Fatalf("应返回 ErrQuizUnavailable，得到: %v", err)
	}
	if len(qs) != 0 {
		t.Fatalf("题目列表应为空")
	}

	// 测试 Submit 已废弃
	answer := SubmitAnswer{
		ID:      "1",
		Choices: []string{"A"},
	}
	_, err = svc.Submit(ctx, 1, "variables", "storage", []SubmitAnswer{answer}, 1500)
	if err == nil {
		t.Fatalf("Submit 已废弃，应返回错误")
	}

	// 测试 History 仍可正常使用
	record := &Record{
		UserID:     1,
		Topic:      "variables",
		Chapter:    "storage",
		Score:      80,
		Total:      100,
		DurationMs: 60000,
	}
	_, err = repo.SaveRecord(ctx, record)
	if err != nil {
		t.Fatalf("保存记录失败: %v", err)
	}

	history, err := svc.History(ctx, 1, "variables", nil, nil)
	if err != nil {
		t.Fatalf("查询历史失败: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("历史记录数量不正确")
	}
}
