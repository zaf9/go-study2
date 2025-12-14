package handler

import (
	"net/http"

	"go-study2/internal/app/http_server/handler/internal"
	progapp "go-study2/internal/app/progress"

	"github.com/gogf/gf/v2/net/ghttp"
)

type progressRequest struct {
	Topic            string `json:"topic"`
	Chapter          string `json:"chapter"`
	ReadDuration     int64  `json:"read_duration"`
	ScrollProgress   int    `json:"scroll_progress"`
	LastPosition     string `json:"last_position"`
	QuizScore        int    `json:"quiz_score"`
	QuizPassed       bool   `json:"quiz_passed"`
	EstimatedSeconds int64  `json:"estimated_seconds"`
	ForceSync        bool   `json:"force_sync"`
}

// GetAllProgress 返回整体与主题汇总。
func (h *Handler) GetAllProgress(r *ghttp.Request) {
	svc, ok := h.getProgressService(r)
	if !ok {
		return
	}
	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}
	overall, topics, err := svc.GetOverallProgress(r.GetCtx(), userID)
	if err != nil {
		writeError(r, http.StatusBadRequest, 40004, err.Error())
		return
	}
	next, _ := svc.GetNextUnfinishedChapter(r.GetCtx(), userID)
	payload := map[string]interface{}{
		"overall": overall,
		"topics":  topics,
		"next":    next,
	}
	writeSuccess(r, "success", payload)
}

// GetTopicProgress 返回指定主题的章节进度。
func (h *Handler) GetTopicProgress(r *ghttp.Request) {
	svc, ok := h.getProgressService(r)
	if !ok {
		return
	}
	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}
	topic := r.Get("topic").String()
	summary, items, err := svc.GetTopicProgress(r.GetCtx(), userID, topic)
	if err != nil {
		writeError(r, http.StatusBadRequest, 40004, err.Error())
		return
	}
	// 使用服务层富化后的章节表示，包含 percent 并确保 status 与 percent 一致（仅响应层面）
	enriched := svc.EnrichChapters(items)
	writeSuccess(r, "success", map[string]interface{}{
		"topic":    summary,
		"chapters": enriched,
	})
}

// SaveProgress 记录或更新学习进度。
func (h *Handler) SaveProgress(r *ghttp.Request) {
	svc, ok := h.getProgressService(r)
	if !ok {
		return
	}
	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}

	var req progressRequest
	if err := r.Parse(&req); err != nil {
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
		return
	}

	result, err := svc.CreateOrUpdateProgress(r.GetCtx(), progapp.UpdateProgressRequest{
		UserID:         userID,
		Topic:          req.Topic,
		Chapter:        req.Chapter,
		ReadDuration:   req.ReadDuration,
		ScrollProgress: req.ScrollProgress,
		LastPosition:   req.LastPosition,
		QuizScore:      req.QuizScore,
		QuizPassed:     req.QuizPassed,
		EstimatedSec:   req.EstimatedSeconds,
		ForceSync:      req.ForceSync,
	})
	if err != nil {
		writeError(r, http.StatusBadRequest, 40004, err.Error())
		return
	}
	writeSuccess(r, "进度已保存", map[string]interface{}{
		"status":          result.Status,
		"overall":         result.Overall,
		"topic":           result.Topic,
		"read_duration":   result.ReadDuration,
		"scroll_progress": result.ScrollProgress,
		"last_position":   result.LastPosition,
	})
}

func (h *Handler) getProgressService(r *ghttp.Request) (*progapp.Service, bool) {
	if h.progressService != nil {
		return h.progressService, true
	}
	svc, err := internal.BuildProgressService()
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50001, "进度服务不可用")
		return nil, false
	}
	h.progressService = svc
	return svc, true
}
