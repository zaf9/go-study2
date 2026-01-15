package handler

import (
	"go-study2/src/learning/properties/httpprop"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Properties HTTP handlers are delegated to the properties package

var propertiesHandler *httpprop.Handler

func init() {
	propertiesHandler = httpprop.NewHandler()
}

// GetPropertiesMenu 返回 Properties 章节菜单。
func (h *Handler) GetPropertiesMenu(r *ghttp.Request) {
	propertiesHandler.GetPropertiesMenu(r)
}

// GetPropertiesContent 返回 Properties 子主题内容。
func (h *Handler) GetPropertiesContent(r *ghttp.Request) {
	propertiesHandler.GetPropertiesContent(r)
}

// GetPropertiesOutline 返回 Properties 提纲。
func (h *Handler) GetPropertiesOutline(r *ghttp.Request) {
	propertiesHandler.GetPropertiesOutline(r)
}

// SubmitPropertiesQuiz 提交 Properties 测验。
func (h *Handler) SubmitPropertiesQuiz(r *ghttp.Request) {
	propertiesHandler.SubmitPropertiesQuiz(r)
}

// SearchProperties 搜索 Properties 关键词。
func (h *Handler) SearchProperties(r *ghttp.Request) {
	propertiesHandler.SearchProperties(r)
}
