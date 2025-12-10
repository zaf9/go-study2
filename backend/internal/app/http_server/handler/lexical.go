package handler

import (
	"fmt"
	"strings"

	"go-study2/internal/app/lexical_elements"

	"github.com/gogf/gf/v2/net/ghttp"
)

type chapterDef struct {
	ID          string // API 路径参数
	Title       string
	ContentFunc func() string
}

// 词法元素章节列表
var lexicalChapters = []chapterDef{
	{"comments", "Comments (注释)", lexical_elements.GetCommentsContent},
	{"tokens", "Tokens (标记)", lexical_elements.GetTokensContent},
	{"semicolons", "Semicolons (分号)", lexical_elements.GetSemicolonsContent},
	{"identifiers", "Identifiers (标识符)", lexical_elements.GetIdentifiersContent},
	{"keywords", "Keywords (关键字)", lexical_elements.GetKeywordsContent},
	{"operators", "Operators (运算符)", lexical_elements.GetOperatorsContent},
	{"integers", "Integers (整数)", lexical_elements.GetIntegersContent},
	{"floats", "Floats (浮点数)", lexical_elements.GetFloatsContent},
	{"imaginary", "Imaginary (虚数)", lexical_elements.GetImaginaryContent},
	{"runes", "Runes (符文)", lexical_elements.GetRunesContent},
	{"strings", "Strings (字符串)", lexical_elements.GetStringsContent},
}

// GetLexicalMenu 获取词法元素菜单
func (h *Handler) GetLexicalMenu(r *ghttp.Request) {
	format := r.GetCtxVar("format").String()

	items := make([]LexicalMenuItem, len(lexicalChapters))
	for i, c := range lexicalChapters {
		items[i] = LexicalMenuItem{
			ID:    i, // 使用索引作为简单 ID，通过 Name 路由
			Title: c.Title,
			Name:  c.ID,
		}
	}

	if format == "html" {
		h.sendLexicalMenuHTML(r, items)
	} else {
		h.sendLexicalMenuJSON(r, items)
	}
}

func (h *Handler) sendLexicalMenuJSON(r *ghttp.Request, items []LexicalMenuItem) {
	response := Response{
		Code:    20000,
		Message: "OK",
		Data: LexicalMenuResponse{
			Items: items,
		},
	}
	r.Response.WriteJson(response)
}

func (h *Handler) sendLexicalMenuHTML(r *ghttp.Request, items []LexicalMenuItem) {
	var sb strings.Builder
	sb.WriteString("<h1>Lexical Elements</h1>\n<ul>\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf("<li><a href=\"/api/v1/topic/lexical_elements/%s?format=html\">%s</a></li>\n", item.Name, item.Title))
	}

	sb.WriteString("</ul>\n")
	sb.WriteString("<a href=\"/api/v1/topics?format=html\" class=\"back-link\">Back to Topics</a>")
	r.Response.Write(getHtmlPage("Lexical Elements", sb.String()))
}

// GetLexicalContent 获取具体章节内容
func (h *Handler) GetLexicalContent(r *ghttp.Request) {
	chapterName := r.Get("chapter").String()
	format := r.GetCtxVar("format").String()

	var contentFunc func() string
	var title string

	for _, c := range lexicalChapters {
		if c.ID == chapterName {
			contentFunc = c.ContentFunc
			title = c.Title
			break
		}
	}

	if contentFunc == nil {
		r.Response.WriteStatus(404)
		if format == "html" {
			r.Response.Write("Chapter not found")
		} else {
			r.Response.WriteJson(Response{
				Code:    404,
				Message: "Chapter not found",
			})
		}
		return
	}

	content := contentFunc()

	if format == "html" {
		var sb strings.Builder
		sb.WriteString("<h1>" + title + "</h1>\n")
		sb.WriteString("<pre>\n" + content + "\n</pre>\n")
		sb.WriteString("<a href=\"/api/v1/topic/lexical_elements?format=html\" class=\"back-link\">Back to Menu</a>")
		r.Response.Write(getHtmlPage(title, sb.String()))
	} else {
		// JSON 响应：直接返回字符串内容
		// 也可以包装在对象中
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
