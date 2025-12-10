package handler

import (
	"go-study2/internal/domain/progress"
	"go-study2/internal/domain/quiz"
	"go-study2/internal/domain/user"
)

// Topic 学习主题
type Topic struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Response 标准响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TopicListResponse 主题列表响应数据
type TopicListResponse struct {
	Topics []Topic `json:"topics"`
}

// LexicalMenuItem 词法元素菜单项
type LexicalMenuItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
}

// LexicalMenuResponse 词法元素菜单响应
type LexicalMenuResponse struct {
	Items []LexicalMenuItem `json:"items"`
}

// Handler 处理 HTTP 请求的控制器，包含需要的领域服务。
type Handler struct {
	userService     *user.Service
	progressService *progress.Service
	quizService     *quiz.Service
}

// New 创建默认 Handler。
func New() *Handler {
	return &Handler{}
}

// NewWithUserService 允许注入自定义用户服务，便于测试。
func NewWithUserService(userSvc *user.Service) *Handler {
	return &Handler{
		userService: userSvc,
	}
}
