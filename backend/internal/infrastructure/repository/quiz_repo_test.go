package repository

import (
	"testing"
	"time"

	"go-study2/internal/domain/quiz"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestQuizRepository_SaveAndList(t *testing.T) {
	db := setupProgressRepoDB(t)
	repo := NewQuizRepository(db)
	ctx := gctx.New()

	record := &quiz.Record{
		UserID:     2,
		Topic:      "variables",
		Chapter:    "storage",
		Score:      2,
		Total:      3,
		DurationMs: 1200,
		Answers:    `[{"id":"q1","choices":["A"]}]`,
	}
	id, err := repo.SaveRecord(ctx, record)
	if err != nil {
		t.Fatalf("保存测验记录失败: %v", err)
	}
	if id == 0 {
		t.Fatalf("返回的记录 ID 无效")
	}

	history, err := repo.ListRecords(ctx, 2, "variables", nil, nil)
	if err != nil {
		t.Fatalf("查询测验历史失败: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("应返回 1 条历史记录，得到 %d", len(history))
	}

	future := time.Now().Add(24 * time.Hour)
	filtered, err := repo.ListRecords(ctx, 2, "variables", &future, nil)
	if err != nil {
		t.Fatalf("按时间筛选失败: %v", err)
	}
	if len(filtered) != 0 {
		t.Fatalf("未来时间应返回空结果")
	}
}

