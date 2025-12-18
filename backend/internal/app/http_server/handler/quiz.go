package handler

import (
	"net/http"
	"time"

	"go-study2/internal/app/http_server/handler/internal"
	appquiz "go-study2/internal/app/quiz"

	"github.com/gogf/gf/v2/net/ghttp"
)

type quizSubmitRequest struct {
	SessionID  string                     `json:"sessionId"`
	Topic      string                     `json:"topic"`
	Chapter    string                     `json:"chapter"`
	Answers    []appquiz.AnswerSubmission `json:"answers"`
	DurationMs int64                      `json:"durationMs"`
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

	payload, err := svc.GetQuizQuestions(r.GetCtx(), userID, topic, chapter)
	if err != nil {
		h.writeQuizError(r, err)
		return
	}
	writeSuccess(r, "success", map[string]interface{}{
		"topic":     payload.Topic,
		"chapter":   payload.Chapter,
		"sessionId": payload.SessionID,
		"questions": payload.Questions,
	})
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

	result, err := svc.SubmitQuiz(r.GetCtx(), userID, req.SessionID, req.Topic, req.Chapter, req.Answers)
	if err != nil {
		h.writeQuizError(r, err)
		return
	}

	writeSuccess(r, "提交成功", map[string]interface{}{
		"score":           result.Score,
		"total_questions": result.TotalQuestions,
		"correct_answers": result.CorrectAnswers,
		"passed":          result.Passed,
		"details":         result.Details,
		"submittedAt":     time.Now().Format(time.RFC3339),
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

	topic := r.Get("topic").String()
	items, err := svc.GetQuizHistory(r.GetCtx(), userID, topic, 20)
	if err != nil {
		h.writeQuizError(r, err)
		return
	}
	writeSuccess(r, "success", items)
}

// GetQuizReview 返回指定会话的回顾详情。
func (h *Handler) GetQuizReview(r *ghttp.Request) {
	svc, ok := h.getQuizService(r)
	if !ok {
		return
	}
	userID := r.GetCtxVar("user_id").Int64()
	if userID <= 0 {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}
	sessionID := r.Get("sessionId").String()
	if sessionID == "" {
		writeError(r, http.StatusBadRequest, 40004, "会话ID不能为空")
		return
	}

	detail, err := svc.GetQuizReview(r.GetCtx(), userID, sessionID)
	if err != nil {
		h.writeQuizError(r, err)
		return
	}

	writeSuccess(r, "success", detail)
}

// GetQuizStats 返回章节题库统计信息: total, byType, byDifficulty
func (h *Handler) GetQuizStats(r *ghttp.Request) {
	svc, ok := h.getQuizService(r)
	if !ok {
		return
	}
	topic := r.Get("topic").String()
	chapter := r.Get("chapter").String()
	stats, err := svc.GetStats(r.GetCtx(), topic, chapter)
	if err != nil {
		h.writeQuizError(r, err)
		return
	}
	writeSuccess(r, "success", stats)
}

func (h *Handler) getQuizService(r *ghttp.Request) (*appquiz.Service, bool) {
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
	case appquiz.ErrInvalidInput:
		writeError(r, http.StatusBadRequest, 40004, "请求参数无效")
	case appquiz.ErrQuizUnavailable:
		writeSuccess(r, "当前主题暂无测验", []interface{}{})
	case appquiz.ErrDuplicateSubmit:
		writeError(r, http.StatusConflict, 40009, "重复提交会话")
	default:
		writeError(r, http.StatusInternalServerError, 50001, "服务器繁忙，请稍后再试")
	}
}
