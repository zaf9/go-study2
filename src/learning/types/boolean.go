package types

import "fmt"

// registerBoolean 注册布尔类型内容与测验。
func registerBoolean() {
	_ = RegisterContent(TopicBoolean, TopicContent{
		Concept: TypeConcept{
			ID:        "boolean",
			Category:  "basic",
			Title:     "布尔类型与条件判断",
			Summary:   "布尔类型只有 true/false，零值为 false，常用于条件与逻辑运算。",
			GoVersion: "1.24",
			Rules:     []string{"只能取 true/false", "比较运算结果为布尔值", "零值为 false"},
			Keywords:  []string{"bool", "condition"},
			PrintableOutline: []string{
				"bool 零值为 false，不会是未定义值",
				"比较与逻辑运算结果为 bool",
				"禁止将非布尔值作为条件",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-BOOL-ONLY", ConceptID: "boolean", RuleType: "identity", Description: "布尔类型只接受 true/false"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-bool-true",
				ConceptID:      "boolean",
				Title:          "合法布尔比较",
				Code:           "var ok bool = (3 > 1)\nfmt.Println(ok)",
				ExpectedOutput: "true",
				IsValid:        true,
				RuleRef:        "TR-BOOL-ONLY",
			},
			{
				ID:             "ex-bool-condition",
				ConceptID:      "boolean",
				Title:          "非布尔条件错误",
				Code:           "var n int = 1\n// if n { } // 编译错误，需显式比较",
				ExpectedOutput: "",
				IsValid:        false,
				RuleRef:        "TR-BOOL-ONLY",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-bool-1",
				ConceptID:   "boolean",
				Stem:        "布尔类型的零值是 false 吗？",
				Options:     []string{"是", "否"},
				Answer:      "A",
				Explanation: "bool 零值固定为 false。",
				RuleRef:     "TR-BOOL-ONLY",
				Difficulty:  "easy",
			},
			{
				ID:        "q-bool-2",
				ConceptID: "boolean",
				Stem:      "下列语句哪一行会编译错误？",
				Options: []string{
					"if 3 > 1 {}",
					"var n int = 1; if n {}",
				},
				Answer:      "B",
				Explanation: "条件必须是 bool，int 需比较后使用。",
				RuleRef:     "TR-BOOL-ONLY",
				Difficulty:  "easy",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "bool zero",
				ConceptID: "boolean",
				Summary:   "布尔类型零值为 false，用于条件判断初始值。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/types/boolean",
					"cli":  "types > boolean",
				},
			},
		},
	})
}

// BooleanOutline 返回布尔主题的提纲。
func BooleanOutline() []string {
	return []string{
		"bool 仅 true/false，零值 false",
		"比较和逻辑运算结果类型为 bool",
		"条件判断必须使用 bool 表达式",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/boolean"),
	}
}
