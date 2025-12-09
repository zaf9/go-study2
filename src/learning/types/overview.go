package types

import (
	"sort"
	"strings"
)

// TypeOverview 汇总章节概览与提纲。
type TypeOverview struct {
	Title        string        `json:"title"`
	Version      string        `json:"version"`
	Concepts     []TypeConcept `json:"concepts"`
	Printable    []string      `json:"printable"`
	OutlineNotes []string      `json:"outlineNotes"`
}

// GetOverview 汇总所有子主题的概念摘要，便于 CLI/HTTP/打印复用。
func GetOverview() TypeOverview {
	var concepts []TypeConcept
	for _, topic := range AllTopics() {
		if c, ok := conceptRegistry[topic]; ok {
			concepts = append(concepts, c)
		}
	}
	sort.Slice(concepts, func(i, j int) bool {
		return concepts[i].ID < concepts[j].ID
	})

	printable := collectPrintableOutline()
	return TypeOverview{
		Title:     "Go 类型章节概览",
		Version:   "Go 1.24.5",
		Concepts:  concepts,
		Printable: printable,
		OutlineNotes: []string{
			"基础类型：布尔、数值、字符串",
			"复合类型：数组、切片、结构体、指针、函数、map、chan",
			"接口与类型集：基础/嵌入/类型集/实现判定",
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
	outline = append(outline, BooleanOutline()...)
	outline = append(outline, NumericOutline()...)
	outline = append(outline, StringOutline()...)
	outline = append(outline, ArrayOutline()...)
	outline = append(outline, SliceOutline()...)
	outline = append(outline, StructOutline()...)
	outline = append(outline, PointerOutline()...)
	outline = append(outline, FunctionOutline()...)
	outline = append(outline, InterfaceBasicOutline()...)
	outline = append(outline, InterfaceEmbeddedOutline()...)
	outline = append(outline, InterfaceGeneralOutline()...)
	outline = append(outline, InterfaceImplOutline()...)
	outline = append(outline, MapOutline()...)
	outline = append(outline, ChannelOutline()...)
	return outline
}
