package http

import (
	"net/http"
)

// Handler 提供 Types 章节的 HTTP 处理占位。
type Handler struct{}

// NewHandler 创建 Types Handler。
func NewHandler() *Handler {
	return &Handler{}
}

// Menu 返回占位响应，后续将替换为真实菜单数据。
func (h *Handler) Menu(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("Types 菜单功能待实现"))
}

// Content 返回占位内容，保证路由可用。
func (h *Handler) Content(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("Types 内容功能待实现"))
}

// Quiz 返回占位测验接口。
func (h *Handler) Quiz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte("Types 测验功能待实现"))
}
