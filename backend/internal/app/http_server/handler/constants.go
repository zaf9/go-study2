package handler

import (
	"fmt"
	"strings"

	"go-study2/internal/app/constants"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Constants 子主题定义
var constantsChapters = []chapterDef{
	{"boolean", "Boolean Constants (布尔常量)", constants.GetBooleanContent},
	{"rune", "Rune Constants (符文常量)", constants.GetRuneContent},
	{"integer", "Integer Constants (整数常量)", constants.GetIntegerContent},
	{"floating_point", "Floating-point Constants (浮点常量)", constants.GetFloatingPointContent},
	{"complex", "Complex Constants (复数常量)", constants.GetComplexContent},
	{"string", "String Constants (字符串常量)", constants.GetStringContent},
	{"expressions", "Constant Expressions (常量表达式)", constants.GetExpressionsContent},
	{"typed_untyped", "Typed and Untyped Constants (类型化/无类型化常量)", constants.GetTypedUntypedContent},
	{"conversions", "Conversions (类型转换)", constants.GetConversionsContent},
	{"builtin_functions", "Built-in Functions (内置函数)", constants.GetBuiltinFunctionsContent},
	{"iota", "Iota (iota 特性)", constants.GetIotaContent},
	{"implementation_restrictions", "Implementation Restrictions (实现限制)", constants.GetImplementationRestrictionsContent},
}

// GetConstantsMenu 获取 Constants 菜单
func (h *Handler) GetConstantsMenu(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()

	items := make([]LexicalMenuItem, len(constantsChapters))
	for i, c := range constantsChapters {
		items[i] = LexicalMenuItem{
			ID:    i, // 使用索引作为简单 ID
			Title: c.Title,
			Name:  c.ID,
		}
	}

	if format == "html" {
		h.sendConstantsMenuHTML(r, items)
	} else {
		h.sendConstantsMenuJSON(r, items)
	}
}

func (h *Handler) sendConstantsMenuJSON(r *ghttp.Request, items []LexicalMenuItem) {
	response := Response{
		Code:    20000,
		Message: "OK",
		Data: LexicalMenuResponse{
			Items: items,
		},
	}
	r.Response.WriteJson(response)
}

func (h *Handler) sendConstantsMenuHTML(r *ghttp.Request, items []LexicalMenuItem) {
	var sb strings.Builder
	sb.WriteString("<h1>Constants Learning</h1>\n<ul>\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf("<li><a href=\"/api/v1/topic/constants/%s?format=html\">%s</a></li>\n", item.Name, item.Title))
	}

	sb.WriteString("</ul>\n")
	sb.WriteString("<a href=\"/api/v1/topics?format=html\" class=\"back-link\">Back to Topics</a>")
	r.Response.Write(getHtmlPage("Constants Learning", sb.String()))
}

// GetConstantsContent 获取具体章节内容
func (h *Handler) GetConstantsContent(r *ghttp.Request) {
	chapterName := r.Get("subtopic").String() // 路由参数 subtopic
	format := r.GetCtxVar("format").String()

	var contentFunc func() string
	var title string

	for _, c := range constantsChapters {
		if c.ID == chapterName {
			contentFunc = c.ContentFunc
			title = c.Title
			break
		}
	}

	if contentFunc == nil {
		r.Response.WriteStatus(404)
		if format == "html" {
			r.Response.Write("Subtopic not found")
		} else {
			r.Response.WriteJson(Response{
				Code:    404,
				Message: "Subtopic not found",
			})
		}
		return
	}

	content := contentFunc()

	if format == "html" {
		var sb strings.Builder
		sb.WriteString("<h1>" + title + "</h1>\n")
		sb.WriteString("<pre>\n" + content + "\n</pre>\n")
		sb.WriteString("<a href=\"/api/v1/topic/constants?format=html\" class=\"back-link\">Back to Menu</a>")
		r.Response.Write(getHtmlPage(title, sb.String()))
	} else {
		// JSON 响应
		r.Response.WriteJson(Response{
			Code:    20000,
			Message: "OK",
			Data: map[string]string{
				"title":   title,
				"content": content,
			},
		})
	}
}
