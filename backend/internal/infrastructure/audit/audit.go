package audit

import (
	"context"
	"time"

	"go-study2/internal/infrastructure/database"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// Event 表示审计事件模型。
type Event struct {
	ID        int64     `json:"id"`
	EventType string    `json:"eventType"`
	UserID    int64     `json:"userId"`
	Result    string    `json:"result"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"createdAt"`
}

// Record 写入审计事件，失败时记录日志但不阻断业务。
func Record(ctx context.Context, eventType string, userID int64, result string, metadata string) {
	db := database.Default()
	if db == nil {
		return
	}
	_, err := db.Insert(ctx, "audit_events", g.Map{
		"event_type": eventType,
		"user_id":    userID,
		"result":     result,
		"metadata":   metadata,
	})
	if err != nil {
		g.Log().Warningf(gctx.New(), "写入审计事件失败: %v", err)
	}
}
