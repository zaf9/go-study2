package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"go-study2/src/learning/types"

	"github.com/gogf/gf/v2/net/ghttp"
)

// GetTypesMenu 返回 Types 章节菜单。
func (h *Handler) GetTypesMenu(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	items := buildTypesMenuItems()

	if format == "html" {
		h.sendTypesMenuHTML(r, items)
		return
	}

	r.Response.WriteJson(Response{
		Code:    0,
		Message: "OK",
		Data: LexicalMenuResponse{
			Items: items,
		},
	})
}

func (h *Handler) sendTypesMenuHTML(r *ghttp.Request, items []LexicalMenuItem) {
	var sb strings.Builder
	sb.WriteString("<h1>Types Learning</h1>\n<ul>\n")
	for _, item := range items {
		sb.WriteString(fmt.Sprintf("<li><a href=\"/api/v1/topic/types/%s?format=html\">%s</a></li>\n", item.Name, item.Title))
	}
	sb.WriteString("</ul>\n")
	sb.WriteString("<a href=\"/api/v1/topics?format=html\" class=\"back-link\">返回主题列表</a>")
	r.Response.Write(getHtmlPage("Types Learning", sb.String()))
}

// GetTypesContent 返回子主题内容，占位实现提示待上线。
func (h *Handler) GetTypesContent(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	subtopic := r.Get("subtopic").String()
	topic := types.NormalizeTopic(subtopic)
	if !types.IsSupportedTopic(topic) {
		h.writeTypesNotFound(r, format, "未知的 Types 子主题")
		return
	}

	content, err := types.LoadContent(topic)
	if err != nil {
		h.writeTypesNotFound(r, format, err.Error())
		return
	}
	quiz, quizErr := types.LoadQuiz(topic)

	if format == "html" {
		h.sendTypesContentHTML(r, content, quiz, quizErr)
		return
	}

	if quizErr != nil && quizErr != types.ErrQuizUnavailable {
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

// SubmitTypesQuiz 接收综合测验答案并评分。
func (h *Handler) SubmitTypesQuiz(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()

	var payload struct {
		Answers []struct {
			ID     string `json:"id"`
			Choice string `json:"choice"`
		} `json:"answers"`
	}
	body, _ := io.ReadAll(r.Body)
	if len(body) == 0 {
		body = r.GetBody()
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		h.writeErrorJSON(r, 400, "请求体解析失败")
		return
	}
	answerMap := map[string]string{}
	for _, a := range payload.Answers {
		answerMap[a.ID] = a.Choice
	}
	if len(answerMap) == 0 {
		for _, item := range types.LoadComprehensiveQuiz() {
			answerMap[item.ID] = item.Answer
		}
	}

	result, err := types.EvaluateComprehensiveQuiz(answerMap)
	if err != nil {
		h.writeErrorJSON(r, 400, err.Error())
		return
	}

	if format == "html" {
		h.sendTypesQuizHTML(r, "comprehensive", result)
		return
	}

	r.Response.WriteJson(Response{
		Code:    0,
		Message: "OK",
		Data:    result,
	})
}

// SearchTypes 返回检索占位响应。
func (h *Handler) SearchTypes(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	keyword := r.GetQuery("keyword").String()
	results, err := types.SearchReferences(keyword)
	if err != nil {
		h.writeErrorJSON(r, 400, err.Error())
		return
	}
	if len(results) == 0 {
		h.writeTypesNotFound(r, format, "未找到匹配关键词")
		return
	}

	if format == "html" {
		h.sendTypesSearchHTML(r, keyword, results)
		return
	}

	r.Response.WriteJson(Response{
		Code:    0,
		Message: "OK",
		Data: map[string]interface{}{
			"keyword": keyword,
			"results": results,
		},
	})
}

func buildTypesMenuItems() []LexicalMenuItem {
	topics := types.AllTopics()
	items := make([]LexicalMenuItem, 0, len(topics))
	for idx, topic := range topics {
		items = append(items, LexicalMenuItem{
			ID:    idx,
			Title: formatTypesTitle(topic),
			Name:  string(topic),
		})
	}
	return items
}

func formatTypesTitle(topic types.Topic) string {
	switch topic {
	case types.TopicBoolean:
		return "Boolean (布尔类型)"
	case types.TopicNumeric:
		return "Numeric (数值类型)"
	case types.TopicString:
		return "String (字符串)"
	case types.TopicArray:
		return "Array (数组)"
	case types.TopicSlice:
		return "Slice (切片)"
	case types.TopicStruct:
		return "Struct (结构体)"
	case types.TopicPointer:
		return "Pointer (指针)"
	case types.TopicFunction:
		return "Function (函数)"
	case types.TopicInterfaceBasic:
		return "Interface Basic (接口基础)"
	case types.TopicInterfaceEmbedded:
		return "Interface Embedded (接口嵌入)"
	case types.TopicInterfaceGeneral:
		return "Interface General (类型集)"
	case types.TopicInterfaceImpl:
		return "Interface Impl (接口实现)"
	case types.TopicMap:
		return "Map (映射)"
	case types.TopicChannel:
		return "Channel (通道)"
	default:
		return string(topic)
	}
}

func (h *Handler) sendTypesContentHTML(r *ghttp.Request, content types.TopicContent, quiz []types.QuizItem, quizErr error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<h1>%s</h1>\n", content.Concept.Title))
	sb.WriteString(fmt.Sprintf("<p>%s</p>\n", content.Concept.Summary))
	if len(content.Concept.Rules) > 0 {
		sb.WriteString("<h3>规则</h3><ul>\n")
		for _, rule := range content.Concept.Rules {
			sb.WriteString(fmt.Sprintf("<li>%s</li>\n", rule))
		}
		sb.WriteString("</ul>\n")
	}
	if len(content.Examples) > 0 {
		sb.WriteString("<h3>示例</h3>\n")
		for _, ex := range content.Examples {
			sb.WriteString(fmt.Sprintf("<h4>%s</h4><pre>%s</pre>\n", ex.Title, ex.Code))
			if ex.ExpectedOutput != "" {
				sb.WriteString(fmt.Sprintf("<p>输出: %s</p>\n", ex.ExpectedOutput))
			}
		}
	}
	sb.WriteString("<h3>测验</h3>\n")
	if quizErr == types.ErrQuizUnavailable {
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
	sb.WriteString("<a href=\"/api/v1/topic/types?format=html\" class=\"back-link\">返回菜单</a>")
	r.Response.Write(getHtmlPage(content.Concept.Title, sb.String()))
}

func (h *Handler) sendTypesQuizHTML(r *ghttp.Request, topic types.Topic, result types.QuizResult) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<h1>Types 测验结果 - %s</h1>", topic))
	sb.WriteString(fmt.Sprintf("<p>得分: %d / %d</p>", result.Score, result.Total))
	if len(result.Details) > 0 {
		sb.WriteString("<ol>\n")
		for _, d := range result.Details {
			state := "错误"
			if d.Correct {
				state = "正确"
			}
			sb.WriteString(fmt.Sprintf("<li>%s - %s (答案: %s)</li>\n", d.ID, state, d.Answer))
			sb.WriteString(fmt.Sprintf("<p>%s</p>\n", d.Explanation))
		}
		sb.WriteString("</ol>\n")
	}
	sb.WriteString("<a href=\"/api/v1/topic/types?format=html\" class=\"back-link\">返回 Types 菜单</a>")
	r.Response.Write(getHtmlPage("Types Quiz", sb.String()))
}

func (h *Handler) sendTypesSearchHTML(r *ghttp.Request, keyword string, results []types.ReferenceIndex) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<h1>Types 搜索: %s</h1>\n<ul>\n", keyword))
	for _, res := range results {
		sb.WriteString(fmt.Sprintf("<li><strong>%s</strong>: %s</li>\n", res.Keyword, res.Summary))
	}
	sb.WriteString("</ul>\n<a href=\"/api/v1/topic/types?format=html\" class=\"back-link\">返回 Types 菜单</a>")
	r.Response.Write(getHtmlPage("Types Search", sb.String()))
}

// GetTypesOutline 返回类型提纲。
func (h *Handler) GetTypesOutline(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	overview := types.GetOverview()

	if format == "html" {
		var sb strings.Builder
		sb.WriteString("<h1>Types 提纲</h1>\n")
		for _, line := range overview.Printable {
			sb.WriteString(fmt.Sprintf("<p>%s</p>\n", line))
		}
		sb.WriteString("<a href=\"/api/v1/topic/types?format=html\" class=\"back-link\">返回 Types 菜单</a>")
		r.Response.Write(getHtmlPage("Types Outline", sb.String()))
		return
	}

	r.Response.WriteJson(Response{
		Code:    0,
		Message: "OK",
		Data: map[string]interface{}{
			"title":     overview.Title,
			"version":   overview.Version,
			"printable": overview.Printable,
		},
	})
}

func (h *Handler) writeTypesNotFound(r *ghttp.Request, format, msg string) {
	if format == "html" {
		r.Response.WriteStatus(404, getHtmlPage("Not Found", fmt.Sprintf("<p>%s</p><a href=\"/api/v1/topic/types?format=html\" class=\"back-link\">返回 Types 菜单</a>", msg)))
		return
	}
	r.Response.WriteStatusExit(404, Response{
		Code:    404,
		Message: msg,
	})
}
