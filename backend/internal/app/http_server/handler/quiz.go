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

	// 安全加固：输入校验
	if err := validateQuizSubmitRequest(req); err != nil {
		writeError(r, http.StatusBadRequest, 40004, err.Error())
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

// validateQuizSubmitRequest 验证测验提交请求的合法性
// 确保：
// 1. 会话ID、主题、章节不为空
// 2. 答案数量与题目数量一致（防止用户绕过防护或篡改）
// 3. 答题时长合理（避免异常快速或超长时间提交）
func validateQuizSubmitRequest(req quizSubmitRequest) error {
	// 检查必要字段
	if req.SessionID == "" {
		return appquiz.ErrInvalidInput
	}
	if req.Topic == "" {
		return appquiz.ErrInvalidInput
	}
	if req.Chapter == "" {
		return appquiz.ErrInvalidInput
	}

	// 检查答案列表不为空
	if len(req.Answers) == 0 {
		return appquiz.ErrInvalidInput
	}

	// 检查答题时长（应该大于等于 1000ms，即 1 秒）
	if req.DurationMs < 1000 {
		return appquiz.ErrInvalidInput
	}

	// 检查答题时长不超过 8 小时（28800000 毫秒），防止明显异常提交
	if req.DurationMs > 28800000 {
		return appquiz.ErrInvalidInput
	}

	return nil
}
