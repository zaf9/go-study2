package properties

import (
	"sort"
	"strings"
)

// PropertyOverview 汇总章节概览与提纲。
type PropertyOverview struct {
	Title        string            `json:"title"`
	Version      string            `json:"version"`
	Concepts     []PropertyConcept `json:"concepts"`
	Printable    []string          `json:"printable"`
	OutlineNotes []string          `json:"outlineNotes"`
}

// GetOverview 汇总所有子主题的概念摘要，便于 CLI/HTTP/打印复用。
func GetOverview() PropertyOverview {
	var concepts []PropertyConcept
	for _, topic := range AllTopics() {
		if c, ok := conceptRegistry[topic]; ok {
			concepts = append(concepts, c)
		}
	}
	sort.Slice(concepts, func(i, j int) bool {
		return concepts[i].ID < concepts[j].ID
	})

	printable := collectPrintableOutline()
	return PropertyOverview{
		Title:     "Go 类型属性章节概览",
		Version:   "Go 1.24.5",
		Concepts:  concepts,
		Printable: printable,
		OutlineNotes: []string{
			"值的表示：自包含值与引用值的区别",
			"底层类型：类型推导规则与类型别名",
			"核心类型：非接口类型与接口类型的核心类型判定",
			"类型标识：两类型相同/不同的判定规则",
			"可赋值性：值可赋值的6个条件和3个额外条件",
			"可表示性：常量在特定类型范围内的可表示性",
			"方法集：值接收者与指针接收者的方法集规则",
		},
	}
}

// RenderPrintableOutline 返回可打印的文本提纲。
func RenderPrintableOutline() string {
	out := collectPrintableOutline()
	return strings.Join(out, "\n")
}

func collectPrintableOutline() []string {
	var outline []string
	outline = append(outline, RepresentationOutline()...)
	outline = append(outline, UnderlyingTypeOutline()...)
	outline = append(outline, CoreTypeOutline()...)
	outline = append(outline, TypeIdentityOutline()...)
	outline = append(outline, AssignabilityOutline()...)
	outline = append(outline, RepresentabilityOutline()...)
	outline = append(outline, MethodSetOutline()...)
	return outline
}
