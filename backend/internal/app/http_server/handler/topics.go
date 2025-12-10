package handler

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// GetTopics 获取学习主题列表
func (h *Handler) GetTopics(r *ghttp.Request) {
	// 1. 准备数据
	topics := []Topic{
		{
			ID:          "lexical_elements",
			Title:       "Lexical Elements",
			Description: "Go 语言词法元素学习 (Lexical Elements)",
		},
		{
			ID:          "constants",
			Title:       "Constants",
			Description: "Go 语言常量学习 (Constants)",
		},
		{
			ID:          "variables",
			Title:       "Variables",
			Description: "Go 语言变量学习 (Variables)",
		},
		{
			ID:          "types",
			Title:       "Types",
			Description: "Go 语言类型学习 (Types)",
		},
	}

	// 2. 获取格式 (由中间件设置)
	format := r.GetCtxVar("format").String()

	// 3. 根据格式返回响应
	if format == "html" {
		h.sendTopicsHTML(r, topics)
	} else {
		h.sendTopicsJSON(r, topics)
	}
}

func (h *Handler) sendTopicsJSON(r *ghttp.Request, topics []Topic) {
	response := Response{
		Code:    0,
		Message: "OK",
		Data: TopicListResponse{
			Topics: topics,
		},
	}
	r.Response.WriteJson(response)
}

func (h *Handler) sendTopicsHTML(r *ghttp.Request, topics []Topic) {
	var sb strings.Builder
	sb.WriteString("<h1>Available Learning Topics</h1>\n<ul>\n")

	for _, topic := range topics {
		// 生成链接
		sb.WriteString(fmt.Sprintf("<li><a href=\"/api/v1/topic/%s?format=html\"><strong>%s</strong></a><br><span style=\"color:#888;font-size:0.9em\">%s</span></li>\n",
			topic.ID, topic.Title, topic.Description))
	}

	sb.WriteString("</ul>\n")
	r.Response.Write(getHtmlPage("Go Study Topics", sb.String()))
}
