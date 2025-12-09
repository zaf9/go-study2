package types

import "fmt"

// registerMapType 注册 map 内容与测验。
func registerMapType() {
	_ = RegisterContent(TopicMap, TopicContent{
		Concept: TypeConcept{
			ID:        "map",
			Category:  "composite",
			Title:     "map 键可比较性",
			Summary:   "map 键必须可比较，零值为 nil，使用前需 make。",
			GoVersion: "1.24",
			Rules:     []string{"键需可比较", "零值 nil 不能写入"},
			Keywords:  []string{"map key", "comparable"},
			PrintableOutline: []string{
				"键必须可比较：整数、字符串、指针等",
				"零值 nil，不可直接写入，需 make",
				"读取不存在键返回零值与 false",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-MAP-KEY", ConceptID: "map", RuleType: "key_constraint", Description: "map 键类型必须可比较"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-map-key",
				ConceptID:      "map",
				Title:          "不可比较键",
				Code:           "m := map[[]int]int{}\n_ = m",
				ExpectedOutput: "",
				IsValid:        false,
				RuleRef:        "TR-MAP-KEY",
			},
			{
				ID:             "ex-map-ok",
				ConceptID:      "map",
				Title:          "合法键示例",
				Code:           "m := map[string]int{\"go\":1}\nfmt.Println(m[\"go\"])",
				ExpectedOutput: "1",
				IsValid:        true,
				RuleRef:        "TR-MAP-KEY",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-map-1",
				ConceptID:   "map",
				Stem:        "哪种类型可以作为 map 键？",
				Options:     []string{"int", "slice"},
				Answer:      "A",
				Explanation: "切片不可比较，不能作为键。",
				RuleRef:     "TR-MAP-KEY",
				Difficulty:  "easy",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "map key",
				ConceptID: "map",
				Summary:   "map 键必须可比较，如整数、字符串、指针等。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/types/map",
					"cli":  "types > map",
				},
			},
		},
	})
}

// MapOutline 返回 map 提纲。
func MapOutline() []string {
	return []string{
		"键必须可比较，零值 nil 不可写入",
		"读取不存在键返回零值+false",
		"使用 make 初始化 map",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/map"),
	}
}
