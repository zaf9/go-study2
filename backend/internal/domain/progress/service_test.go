package progress

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockRepo struct {
	items []Progress
	err   error
}

func (m *mockRepo) Upsert(_ context.Context, record *Progress) error {
	if m.err != nil {
		return m.err
	}
	record.ID = int64(len(m.items) + 1)
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	m.items = append(m.items, *record)
	return nil
}

func (m *mockRepo) ListByUser(_ context.Context, userID int64) ([]Progress, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []Progress
	for _, item := range m.items {
		if item.UserID == userID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *mockRepo) ListByTopic(_ context.Context, userID int64, topic string) ([]Progress, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []Progress
	for _, item := range m.items {
		if item.UserID == userID && item.Topic == topic {
			result = append(result, item)
		}
	}
	return result, nil
}

func TestService_SaveAndList(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	_, err := svc.Save(ctx, 0, "variables", "storage", StatusInProgress, "")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("应返回 ErrInvalidInput")
	}

	record, err := svc.Save(ctx, 1, "variables", "storage", StatusInProgress, `{"scroll":120}`)
	if err != nil {
		t.Fatalf("保存进度失败: %v", err)
	}
	if record.Topic != "variables" || record.Status != StatusInProgress {
		t.Fatalf("保存进度返回内容不正确")
	}

	list, err := svc.ListAll(ctx, 1)
	if err != nil {
		t.Fatalf("查询全部进度失败: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("期望 1 条进度，得到 %d", len(list))
	}

	listTopic, err := svc.ListByTopic(ctx, 1, "variables")
	if err != nil {
		t.Fatalf("按主题查询失败: %v", err)
	}
	if len(listTopic) != 1 {
		t.Fatalf("按主题查询进度数量不正确")
	}
}

