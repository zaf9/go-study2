package http

import (
	"net/http"

	"go-study2/internal/infra/repository"

	"github.com/gogf/gf/v2/net/ghttp"
)

// QuizHandler 提供测验相关 HTTP 接口。
type QuizHandler struct {
	Repo repository.IQuizSessionRepository
}

// RegisterQuizRoutes 注册测验路由。
func RegisterQuizRoutes(s *ghttp.Server, handler *QuizHandler) {
	group := s.Group("/api/v1/quiz")
	// T112: GET /api/v1/quiz/history - 获取测验历史列表
	group.GET("/history", handler.GetHistory)
	// T113: GET /api/v1/quiz/history/:sessionId - 获取测验详情
	group.GET("/history/:sessionId", handler.GetSessionDetail)
}

// GetHistory 处理获取测验历史列表请求。
// T112: 实现 GET /api/v1/quiz/history
func (h *QuizHandler) GetHistory(r *ghttp.Request) {
	// 1. 验证认证信息
	userID, ok := parseUserID(r)
	if !ok {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}

	// 2. 解析查询参数
	topic := r.Get("topic").String()
	limit := r.Get("limit").Int()
	if limit <= 0 {
		limit = 50 // 默认返回50条记录
	}
	if limit > 100 {
		limit = 100 // 最多返回100条记录
	}

	// 3. 调用仓储层获取测验历史
	sessions, err := h.Repo.GetHistory(r.GetCtx(), userID, topic, limit)
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50000, "获取测验历史失败")
		return
	}

	// 4. 返回成功响应
	writeSuccess(r, map[string]interface{}{
		"list":  sessions,
		"total": len(sessions),
	})
}

// GetSessionDetail 处理获取测验详情请求。
// T113: 实现 GET /api/v1/quiz/history/:sessionId
func (h *QuizHandler) GetSessionDetail(r *ghttp.Request) {
	// 1. 验证认证信息
	userID, ok := parseUserID(r)
	if !ok {
		writeError(r, http.StatusUnauthorized, 40001, "认证信息缺失")
		return
	}

	// 2. 获取路径参数
	sessionID := r.Get("sessionId").String()
	if sessionID == "" {
		writeError(r, http.StatusBadRequest, 40004, "会话ID不能为空")
		return
	}

	// 3. 获取会话信息
	session, err := h.Repo.GetSession(r.GetCtx(), sessionID)
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50000, "获取会话信息失败")
		return
	}
	if session == nil {
		writeError(r, http.StatusNotFound, 40400, "会话不存在")
		return
	}

	// 4. 验证会话所有权
	if session.UserID != userID {
		writeError(r, http.StatusForbidden, 40300, "无权访问此会话")
		return
	}

	// 5. 获取答题详情
	attempts, err := h.Repo.GetAttemptsBySession(r.GetCtx(), sessionID)
	if err != nil {
		writeError(r, http.StatusInternalServerError, 50000, "获取答题详情失败")
		return
	}

	// 6. 返回成功响应
	writeSuccess(r, map[string]interface{}{
		"session":  session,
		"attempts": attempts,
	})
}
