package types

import "fmt"

// registerInterfaceGeneral 注册类型集内容。
func registerInterfaceGeneral() {
	_ = RegisterContent(TopicInterfaceGeneral, TopicContent{
		Concept: TypeConcept{
			ID:        "interface_general",
			Category:  "interface",
			Title:     "类型集与 ~T 约束",
			Summary:   "泛型接口可声明类型集与 ~T 近似约束，用于匹配底层类型。",
			GoVersion: "1.24",
			Rules:     []string{"~T 约束匹配底层类型", "并集使用 |"},
			Keywords:  []string{"type set", "~T"},
			PrintableOutline: []string{
				"类型集用于约束满足的具体类型集合",
				"~T 表示底层类型为 T 的所有命名类型",
				"| 表示并集，可组合多个类型",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-IFACE-TYPESET", ConceptID: "interface_general", RuleType: "type_set", Description: "类型集支持 ~T 与并集"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-iface-typeset",
				ConceptID:      "interface_general",
				Title:          "~int 约束",
				Code:           "type SignedInts interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 }",
				ExpectedOutput: "",
				IsValid:        true,
				RuleRef:        "TR-IFACE-TYPESET",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-iface-typeset-1",
				ConceptID:   "interface_general",
				Stem:        "~int 表示什么含义？",
				Options:     []string{"匹配底层类型为 int 的类型", "仅匹配 int 本身"},
				Answer:      "A",
				Explanation: "~T 约束匹配底层类型为 T 的类型。",
				RuleRef:     "TR-IFACE-TYPESET",
				Difficulty:  "easy",
			},
		},
	})
}

// InterfaceGeneralOutline 返回类型集提纲。
func InterfaceGeneralOutline() []string {
	return []string{
		"类型集可使用并集与 ~T 近似约束",
		"~T 表示底层类型为 T 的命名类型也匹配",
		"用于泛型接口约束可用操作集",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/interface_general"),
	}
}
