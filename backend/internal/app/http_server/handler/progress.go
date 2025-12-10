package handler

import (
	"net/http"
	"time"

	"go-study2/internal/app/http_server/handler/internal"
	"go-study2/internal/domain/progress"

	"github.com/gogf/gf/v2/net/ghttp"
)

type progressRequest struct {
	Topic    string `json:"topic"`
	Chapter  string `json:"chapter"`
	Status   string `json:"status"`
	Position string `json:"position"`
}

type progressResponse struct {
	Topic        string `json:"topic"`
	Chapter      string `json:"chapter"`
	Status       string `json:"status"`
	LastVisit    string `json:"lastVisit"`
	LastPosition string `json:"lastPosition,omitempty"`
}

// GetAllProgress 返回当前用户的全部学习进度。
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

	items, err := svc.ListAll(r.GetCtx(), userID)
	if err != nil {
		h.writeProgressError(r, err)
		return
	}
	writeSuccess(r, "success", toProgressResponses(items))
}

// GetTopicProgress 返回指定主题下的进度。
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
	items, err := svc.ListByTopic(r.GetCtx(), userID, topic)
	if err != nil {
		h.writeProgressError(r, err)
		return
	}
	writeSuccess(r, "success", toProgressResponses(items))
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

	record, err := svc.Save(r.GetCtx(), userID, req.Topic, req.Chapter, req.Status, req.Position)
	if err != nil {
		h.writeProgressError(r, err)
		return
	}
	writeSuccess(r, "进度已保存", progressResponse{
		Topic:        record.Topic,
		Chapter:      record.Chapter,
		Status:       record.Status,
		LastVisit:    record.LastVisit.Format(time.RFC3339),
		LastPosition: record.LastPosition,
	})
}

func (h *Handler) getProgressService(r *ghttp.Request) (*progress.Service, bool) {
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

func (h *Handler) writeProgressError(r *ghttp.Request, err error) {
	if err == nil {
		return
	}
	switch err {
	case progress.ErrInvalidInput:
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
	default:
		writeError(r, http.StatusInternalServerError, 50001, "服务器繁忙，请稍后再试")
	}
}

func toProgressResponses(items []progress.Progress) []progressResponse {
	list := make([]progressResponse, 0, len(items))
	for _, item := range items {
		resp := progressResponse{
			Topic:        item.Topic,
			Chapter:      item.Chapter,
			Status:       item.Status,
			LastVisit:    item.LastVisit.Format(time.RFC3339),
			LastPosition: item.LastPosition,
		}
		list = append(list, resp)
	}
	return list
}

