package handler

import (
	"fmt"
	"strings"

	"go-study2/src/learning/variables"

	"github.com/gogf/gf/v2/net/ghttp"
)

var variableTopics = []struct {
	ID          string
	Title       string
	Topic       variables.Topic
	Description string
}{
	{"storage", "Storage (存储与取址)", variables.TopicStorage, "变量声明、new 与复合字面量取址"},
	{"static", "Static (静态类型与可赋值性)", variables.TopicStatic, "静态类型决定可赋值性与方法集"},
	{"dynamic", "Dynamic (接口动态类型与 nil)", variables.TopicDynamic, "接口动态类型与带类型 nil 区分"},
	{"zero", "Zero (零值与取值规则)", variables.TopicZero, "各类型零值与复合元素零值规则"},
}

// GetVariablesMenu 获取 Variables 菜单
func (h *Handler) GetVariablesMenu(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()

	items := make([]LexicalMenuItem, len(variableTopics))
	for i, v := range variableTopics {
		items[i] = LexicalMenuItem{
			ID:    i,
			Title: v.Title,
			Name:  v.ID,
		}
	}

	if format == "html" {
		h.sendVariablesMenuHTML(r, items)
		return
	}
	h.sendVariablesMenuJSON(r, items)
}

func (h *Handler) sendVariablesMenuJSON(r *ghttp.Request, items []LexicalMenuItem) {
	r.Response.WriteJson(Response{
		Code:    0,
		Message: "OK",
		Data: LexicalMenuResponse{
			Items: items,
		},
	})
}

func (h *Handler) sendVariablesMenuHTML(r *ghttp.Request, items []LexicalMenuItem) {
	var sb strings.Builder
	sb.WriteString("<h1>Variables Learning</h1>\n<ul>\n")
	for _, item := range items {
		sb.WriteString(fmt.Sprintf("<li><a href=\"/api/v1/topic/variables/%s?format=html\">%s</a></li>\n", item.Name, item.Title))
	}
	sb.WriteString("</ul>\n")
	sb.WriteString("<a href=\"/api/v1/topics?format=html\" class=\"back-link\">Back to Topics</a>")
	r.Response.Write(getHtmlPage("Variables Learning", sb.String()))
}

// GetVariableContent 获取指定变量子主题内容与测验
func (h *Handler) GetVariableContent(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	subtopic := r.Get("subtopic").String()
	topic := variables.NormalizeTopic(subtopic)

	if !variables.IsSupportedTopic(topic) {
		h.writeNotFound(r, format, "Subtopic not found")
		return
	}

	content, err := variables.LoadContent(topic)
	if err != nil {
		h.writeNotFound(r, format, err.Error())
		return
	}
	quiz, quizErr := variables.LoadQuiz(topic)

	if format == "html" {
		h.sendVariableContentHTML(r, content, quiz, quizErr)
		return
	}

	if quizErr != nil && quizErr != variables.ErrQuizUnavailable {
		h.writeErrorJSON(r, 500, quizErr.Error())
		return
	}

	r.Response.WriteJson(Response{
		Code:    0,
		Message: "OK",
		Data: map[string]interface{}{
			"content": content,
			"quiz":    quiz,
		},
	})
}

func (h *Handler) sendVariableContentHTML(r *ghttp.Request, content variables.Content, quiz []variables.QuizItem, quizErr error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<h1>%s (%s)</h1>\n", content.Title, content.Topic))
	sb.WriteString(fmt.Sprintf("<p>%s</p>\n", content.Summary))
	if len(content.Details) > 0 {
		sb.WriteString("<ul>\n")
		for _, d := range content.Details {
			sb.WriteString(fmt.Sprintf("<li>%s</li>\n", d))
		}
		sb.WriteString("</ul>\n")
	}
	if content.Snippet != "" {
		sb.WriteString("<h3>Snippet</h3>\n<pre>\n")
		sb.WriteString(content.Snippet)
		sb.WriteString("\n</pre>\n")
	}
	if len(content.Examples) > 0 {
		sb.WriteString("<h3>Examples</h3>\n")
		for _, ex := range content.Examples {
			sb.WriteString(fmt.Sprintf("<h4>%s</h4>\n<pre>\n%s\n</pre>\n<p>输出: %s</p>\n", ex.Title, ex.Code, ex.Output))
			if len(ex.Notes) > 0 {
				sb.WriteString("<ul>\n")
				for _, n := range ex.Notes {
					sb.WriteString(fmt.Sprintf("<li>%s</li>\n", n))
				}
				sb.WriteString("</ul>\n")
			}
		}
	}
	sb.WriteString("<h3>Quiz</h3>\n")
	if quizErr == variables.ErrQuizUnavailable {
		sb.WriteString("<p>当前主题暂无测验。</p>")
	} else if quizErr != nil {
		sb.WriteString(fmt.Sprintf("<p>测验加载失败: %v</p>", quizErr))
	} else {
		sb.WriteString("<ol>\n")
		for _, item := range quiz {
			sb.WriteString(fmt.Sprintf("<li>%s<ul>", item.Stem))
			for i, opt := range item.Options {
				sb.WriteString(fmt.Sprintf("<li>%c) %s</li>", 'A'+i, opt))
			}
			sb.WriteString(fmt.Sprintf("</ul><strong>答案:</strong> %s<br><em>%s</em></li>\n", item.Answer, item.Explanation))
		}
		sb.WriteString("</ol>\n")
	}
	sb.WriteString("<a href=\"/api/v1/topic/variables?format=html\" class=\"back-link\">Back to Menu</a>")
	r.Response.Write(getHtmlPage(string(content.Topic), sb.String()))
}

func (h *Handler) writeNotFound(r *ghttp.Request, format, msg string) {
	if format == "html" {
		r.Response.WriteStatus(404, msg)
		return
	}
	h.writeErrorJSON(r, 404, msg)
}

func (h *Handler) writeErrorJSON(r *ghttp.Request, code int, msg string) {
	r.Response.WriteJson(Response{
		Code:    code,
		Message: msg,
	})
}
