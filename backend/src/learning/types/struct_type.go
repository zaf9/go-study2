package types

import "fmt"

// registerStruct 注册结构体内容与测验。
func registerStruct() {
	_ = RegisterContent(TopicStruct, TopicContent{
		Concept: TypeConcept{
			ID:        "struct",
			Category:  "composite",
			Title:     "结构体与字段比较",
			Summary:   "结构体可比较当且仅当所有字段可比较，零值字段为对应类型零值。",
			GoVersion: "1.24",
			Rules:     []string{"字段逐个比较", "含不可比较字段则整体不可比较"},
			Keywords:  []string{"struct compare"},
			PrintableOutline: []string{
				"结构体零值为各字段零值",
				"含不可比较字段则整体不可比较",
				"标签与嵌入不改变比较规则",
			},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-struct-compare",
				ConceptID:      "struct",
				Title:          "包含切片不可比较",
				Code:           "type T struct{ S []int }\n// var a, b T; _ = (a == b) // 编译错误",
				ExpectedOutput: "",
				IsValid:        false,
				RuleRef:        "TR-STRUCT-COMP",
			},
			{
				ID:             "ex-struct-ok",
				ConceptID:      "struct",
				Title:          "全部可比较字段",
				Code:           "type U struct{ A int; B string }\nu1 := U{A:1,B:\"x\"}\nu2 := U{A:1,B:\"x\"}\nfmt.Println(u1 == u2)",
				ExpectedOutput: "true",
				IsValid:        true,
				RuleRef:        "TR-STRUCT-COMP",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-STRUCT-COMP", ConceptID: "struct", RuleType: "comparability", Description: "结构体含不可比较字段则整体不可比较"},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-struct-1",
				ConceptID:   "struct",
				Stem:        "结构体比较需要什么条件？",
				Options:     []string{"所有字段可比较", "只要字段数量相同"},
				Answer:      "A",
				Explanation: "任一字段不可比较则整个结构体不可比较。",
				RuleRef:     "TR-STRUCT-COMP",
				Difficulty:  "easy",
			},
		},
	})
}

// StructOutline 返回结构体主题提纲。
func StructOutline() []string {
	return []string{
		"比较按字段逐个比较，需全部可比较",
		"零值字段 = 其类型零值",
		"含切片/映射/函数字段导致不可比较",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/struct"),
	}
}
