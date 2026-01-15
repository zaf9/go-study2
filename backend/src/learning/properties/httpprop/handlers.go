package httpprop

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"go-study2/src/learning/properties"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Handler 提供 Properties 章节的 HTTP 处理。
type Handler struct{}

// NewHandler 创建 Properties Handler。
func NewHandler() *Handler {
	return &Handler{}
}

// GetPropertiesMenu 返回 Properties 章节菜单。
func (h *Handler) GetPropertiesMenu(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	items := buildPropertiesMenuItems()

	if format == "html" {
		h.sendPropertiesMenuHTML(r, items)
		return
	}

	sendMenuJSON(r, items)
}

func (h *Handler) sendPropertiesMenuHTML(r *ghttp.Request, items []MenuItem) {
	var sb strings.Builder
	sb.WriteString("<h1>Properties of Types and Values Learning</h1>\n<ul>\n")
	for _, item := range items {
		sb.WriteString(fmt.Sprintf("<li><a href=\"/api/v1/topic/properties/%s?format=html\">%s</a></li>\n", item.Name, item.Title))
	}
	sb.WriteString("</ul>\n")
	sb.WriteString("<a href=\"/api/v1/topics?format=html\" class=\"back-link\">返回主题列表</a>")
	r.Response.Write(getHtmlPage("Properties of Types and Values Learning", sb.String()))
}

// GetPropertiesContent 返回子主题内容。
func (h *Handler) GetPropertiesContent(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	subtopic := r.Get("subtopic").String()
	topic := properties.NormalizeTopic(subtopic)
	if !properties.IsSupportedTopic(topic) {
		writeNotFound(r, format, "未知的 Properties 子主题")
		return
	}

	content, err := properties.LoadContent(topic)
	if err != nil {
		writeNotFound(r, format, err.Error())
		return
	}
	quiz, quizErr := properties.LoadQuiz(topic)

	if format == "html" {
		h.sendPropertiesContentHTML(r, content, quiz, quizErr)
		return
	}

	if quizErr != nil && quizErr != properties.ErrQuizUnavailable {
		writeErrorJSON(r, 500, quizErr.Error())
		return
	}

	writeSuccess(r, "OK", map[string]interface{}{
		"content": content,
		"quiz":    quiz,
	})
}

// SubmitPropertiesQuiz 接收综合测验答案并评分。
func (h *Handler) SubmitPropertiesQuiz(r *ghttp.Request) {
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
		writeErrorJSON(r, 400, "请求体解析失败")
		return
	}
	answerMap := map[string]string{}
	for _, a := range payload.Answers {
		answerMap[a.ID] = a.Choice
	}
	if len(answerMap) == 0 {
		items, err := properties.GetComprehensiveQuiz()
		if err != nil {
			writeErrorJSON(r, 500, err.Error())
			return
		}
		for _, item := range items {
			answerMap[item.ID] = item.Answer
		}
	}

	subtopic := r.Get("subtopic").String()
	topic := properties.NormalizeTopic(subtopic)
	if !properties.IsSupportedTopic(topic) {
		writeErrorJSON(r, 400, "未知的子主题")
		return
	}

	result, err := properties.EvaluateQuiz(topic, answerMap)
	if err != nil {
		writeErrorJSON(r, 400, err.Error())
		return
	}

	if format == "html" {
		h.sendPropertiesQuizHTML(r, topic, result)
		return
	}

	writeSuccess(r, "OK", result)
}

// SearchProperties 返回检索结果。
func (h *Handler) SearchProperties(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	keyword := r.GetQuery("keyword").String()
	results, err := properties.SearchReferences(keyword)
	if err != nil {
		writeErrorJSON(r, 400, err.Error())
		return
	}
	if len(results) == 0 {
		writeNotFound(r, format, "未找到匹配关键词")
		return
	}

	if format == "html" {
		h.sendPropertiesSearchHTML(r, keyword, results)
		return
	}

	writeSuccess(r, "OK", map[string]interface{}{
		"keyword": keyword,
		"results": results,
	})
}

// GetPropertiesOutline 返回 Properties 提纲。
func (h *Handler) GetPropertiesOutline(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()
	overview := properties.GetOverview()

	if format == "html" {
		var sb strings.Builder
		sb.WriteString("<h1>Properties of Types and Values 提纲</h1>\n")
		for _, line := range overview.Printable {
			sb.WriteString(fmt.Sprintf("<p>%s</p>\n", line))
		}
		sb.WriteString("<a href=\"/api/v1/topic/properties?format=html\" class=\"back-link\">返回 Properties 菜单</a>")
		r.Response.Write(getHtmlPage("Properties Outline", sb.String()))
		return
	}

	writeSuccess(r, "OK", map[string]interface{}{
		"title":     overview.Title,
		"version":   overview.Version,
		"printable": overview.Printable,
	})
}

func buildPropertiesMenuItems() []MenuItem {
	topics := properties.AllTopics()
	items := make([]MenuItem, 0, len(topics))
	for idx, topic := range topics {
		items = append(items, MenuItem{
			ID:    idx,
			Title: formatPropertiesTitle(topic),
			Name:  string(topic),
		})
	}
	return items
}

func formatPropertiesTitle(topic properties.Topic) string {
	switch topic {
	case properties.TopicRepresentation:
		return "Representation (值的表示)"
	case properties.TopicUnderlyingType:
		return "Underlying Type (底层类型)"
	case properties.TopicCoreType:
		return "Core Type (核心类型)"
	case properties.TopicTypeIdentity:
		return "Type Identity (类型标识)"
	case properties.TopicAssignability:
		return "Assignability (可赋值性)"
	case properties.TopicRepresentability:
		return "Representability (可表示性)"
	case properties.TopicMethodSet:
		return "Method Set (方法集)"
	default:
		return string(topic)
	}
}

func (h *Handler) sendPropertiesContentHTML(r *ghttp.Request, content properties.TopicContent, quiz []properties.QuizItem, quizErr error) {
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
	if quizErr == properties.ErrQuizUnavailable {
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
	sb.WriteString("<a href=\"/api/v1/topic/properties?format=html\" class=\"back-link\">返回菜单</a>")
	r.Response.Write(getHtmlPage(content.Concept.Title, sb.String()))
}

func (h *Handler) sendPropertiesQuizHTML(r *ghttp.Request, topic properties.Topic, result properties.QuizResult) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<h1>Properties 测验结果 - %s</h1>", topic))
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
	sb.WriteString("<a href=\"/api/v1/topic/properties?format=html\" class=\"back-link\">返回 Properties 菜单</a>")
	r.Response.Write(getHtmlPage("Properties Quiz", sb.String()))
}

func (h *Handler) sendPropertiesSearchHTML(r *ghttp.Request, keyword string, results []properties.ReferenceIndex) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<h1>Properties 搜索: %s</h1>\n<ul>\n", keyword))
	for _, res := range results {
		sb.WriteString(fmt.Sprintf("<li><strong>%s</strong>: %s</li>\n", res.Keyword, res.Summary))
	}
	sb.WriteString("</ul>\n<a href=\"/api/v1/topic/properties?format=html\" class=\"back-link\">返回 Properties 菜单</a>")
	r.Response.Write(getHtmlPage("Properties Search", sb.String()))
}

// MenuItem 菜单项结构
type MenuItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
}

func sendMenuJSON(r *ghttp.Request, items []MenuItem) {
	r.Response.WriteJson(map[string]interface{}{
		"code":    20000,
		"message": "success",
		"data": map[string]interface{}{
			"items": items,
		},
	})
}

func writeSuccess(r *ghttp.Request, message string, data interface{}) {
	r.Response.WriteJson(map[string]interface{}{
		"code":    20000,
		"message": message,
		"data":    data,
	})
}

func writeErrorJSON(r *ghttp.Request, code int, message string) {
	r.Response.WriteJson(map[string]interface{}{
		"code":    code,
		"message": message,
	})
}

func writeNotFound(r *ghttp.Request, format string, message string) {
	if format == "html" {
		r.Response.Write(getHtmlPage("Not Found", fmt.Sprintf("<p>%s</p>", message)))
		return
	}
	writeErrorJSON(r, 404, message)
}

func getHtmlPage(title, body string) []byte {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #333; }
        pre { background: #f4f4f4; padding: 10px; border-radius: 5px; overflow-x: auto; }
        code { background: #f4f4f4; padding: 2px 5px; border-radius: 3px; }
        .back-link { display: inline-block; margin-top: 20px; color: #0066cc; }
    </style>
</head>
<body>
%s
</body>
</html>`, title, body)
	return []byte(html)
}
