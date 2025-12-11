package http

import (
	"net/http"
	"strconv"
	"strings"

	progapp "go-study2/internal/app/progress"
	progressdom "go-study2/internal/domain/progress"

	"github.com/gogf/gf/v2/net/ghttp"
)

// ProgressHandler 提供进度相关 HTTP 接口。
type ProgressHandler struct {
	Service *progapp.Service
}

// RegisterProgressRoutes 注册进度路由。
func RegisterProgressRoutes(s *ghttp.Server, handler *ProgressHandler) {
	group := s.Group("/api/v1")
	group.POST("/progress", handler.PostProgress)
	group.GET("/progress", handler.GetProgress)
	group.GET("/progress/{topic}", handler.GetTopicProgress)
}

// PostProgress 处理进度上报。
func (h *ProgressHandler) PostProgress(r *ghttp.Request) {
	userID, ok := parseUserID(r)
	if !ok {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}
	var req progressBody
	if err := r.GetStruct(&req); err != nil {
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
		return
	}

	result, err := h.Service.CreateOrUpdateProgress(r.GetCtx(), progapp.UpdateProgressRequest{
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
	writeSuccess(r, map[string]interface{}{
		"status":          result.Status,
		"overall":         result.Overall,
		"topic":           result.Topic,
		"read_duration":   result.ReadDuration,
		"scroll_progress": result.ScrollProgress,
		"last_position":   result.LastPosition,
	})
}

// GetProgress 返回整体与主题进度汇总。
func (h *ProgressHandler) GetProgress(r *ghttp.Request) {
	userID, ok := parseUserID(r)
	if !ok {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}
	overall, topics, err := h.Service.GetOverallProgress(r.GetCtx(), userID)
	if err != nil {
		writeError(r, http.StatusBadRequest, 40004, err.Error())
		return
	}
	next, _ := h.Service.GetNextUnfinishedChapter(r.GetCtx(), userID)
	writeSuccess(r, map[string]interface{}{
		"overall": overall,
		"topics":  topics,
		"next":    next,
	})
}

// GetTopicProgress 返回指定主题的章节进度列表。
func (h *ProgressHandler) GetTopicProgress(r *ghttp.Request) {
	userID, ok := parseUserID(r)
	if !ok {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}
	topic := r.Get("topic").String()
	summary, items, err := h.Service.GetTopicProgress(r.GetCtx(), userID, topic)
	if err != nil {
		writeError(r, http.StatusBadRequest, 40004, err.Error())
		return
	}
	writeSuccess(r, map[string]interface{}{
		"topic":    summary,
		"chapters": items,
	})
}

func parseUserID(r *ghttp.Request) (int64, bool) {
	headerVal := r.Header.Get("X-User-ID")
	if strings.TrimSpace(headerVal) == "" {
		return 0, false
	}
	uid, err := strconv.ParseInt(headerVal, 10, 64)
	if err != nil || uid <= 0 {
		return 0, false
	}
	return uid, true
}

type progressBody struct {
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

// ensure interfaces aligned to domain constants，方便 handler 内校验。
var (
	_ = progressdom.StatusCompleted
)
