package handler

import (
	"net/http"
	"time"

	"go-study2/internal/app/http_server/handler/internal"
	"go-study2/internal/domain/quiz"

	"github.com/gogf/gf/v2/net/ghttp"
)

type quizSubmitRequest struct {
	Topic      string              `json:"topic"`
	Chapter    string              `json:"chapter"`
	Answers    []quiz.SubmitAnswer `json:"answers"`
	DurationMs int64               `json:"durationMs"`
}

// GetQuiz 返回指定主题与章节的测验题目。
func (h *Handler) GetQuiz(r *ghttp.Request) {
	svc, ok := h.getQuizService(r)
	if !ok {
		return
	}
	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}
	topic := r.Get("topic").String()
	chapter := r.Get("chapter").String()

	questions, err := svc.GetQuestions(r.GetCtx(), topic, chapter)
	if err != nil {
		h.writeQuizError(r, err)
		return
	}
	writeSuccess(r, "success", questions)
}

// SubmitQuiz 接收测验答案并返回评分结果。
func (h *Handler) SubmitQuiz(r *ghttp.Request) {
	svc, ok := h.getQuizService(r)
	if !ok {
		return
	}
	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}

	var req quizSubmitRequest
	if err := r.Parse(&req); err != nil {
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
		return
	}

	result, err := svc.Submit(r.GetCtx(), userID, req.Topic, req.Chapter, req.Answers, req.DurationMs)
	if err != nil {
		h.writeQuizError(r, err)
		return
	}

	writeSuccess(r, "提交成功", map[string]interface{}{
		"score":       result.Score,
		"total":       result.Total,
		"correctIds":  result.CorrectIDs,
		"wrongIds":    result.WrongIDs,
		"submittedAt": result.SubmittedAt.Format(time.RFC3339),
		"durationMs":  result.DurationMs,
	})
}

// GetQuizHistory 返回当前用户的测验历史。
func (h *Handler) GetQuizHistory(r *ghttp.Request) {
	svc, ok := h.getQuizService(r)
	if !ok {
		return
	}
	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}

	var fromPtr, toPtr *time.Time
	if fromStr := r.GetQuery("from").String(); fromStr != "" {
		if ts, err := time.Parse(time.RFC3339, fromStr); err == nil {
			fromPtr = &ts
		}
	}
	if toStr := r.GetQuery("to").String(); toStr != "" {
		if ts, err := time.Parse(time.RFC3339, toStr); err == nil {
			toPtr = &ts
		}
	}

	topic := r.Get("topic").String()
	items, err := svc.History(r.GetCtx(), userID, topic, fromPtr, toPtr)
	if err != nil {
		h.writeQuizError(r, err)
		return
	}
	writeSuccess(r, "success", items)
}

func (h *Handler) getQuizService(r *ghttp.Request) (*quiz.Service, bool) {
	if h.quizService != nil {
		return h.quizService, true
	}
	svc, err := internal.BuildQuizService()
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50001, "测验服务不可用")
		return nil, false
	}
	h.quizService = svc
	return svc, true
}

func (h *Handler) writeQuizError(r *ghttp.Request, err error) {
	if err == nil {
		return
	}
	switch err {
	case quiz.ErrInvalidInput:
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
	case quiz.ErrQuizUnavailable:
		writeSuccess(r, "当前主题暂无测验", []interface{}{})
	default:
		writeError(r, http.StatusInternalServerError, 50001, "服务器繁忙，请稍后再试")
	}
}
