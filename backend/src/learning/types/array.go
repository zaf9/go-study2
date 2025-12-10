package types

import "fmt"

// registerArray 注册数组内容与测验。
func registerArray() {
	_ = RegisterContent(TopicArray, TopicContent{
		Concept: TypeConcept{
			ID:        "array",
			Category:  "composite",
			Title:     "数组与长度是类型一部分",
			Summary:   "数组长度属于类型，比较需要可比较元素且长度一致。",
			GoVersion: "1.24",
			Rules:     []string{"长度是类型一部分", "元素类型可比较时数组可比较"},
			Keywords:  []string{"array length", "comparable"},
			PrintableOutline: []string{
				"数组类型 = 元素类型 + 长度",
				"长度不同不可赋值/比较",
				"元素可比较时数组可比较",
			},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-array-len",
				ConceptID:      "array",
				Title:          "数组长度类型约束",
				Code:           "var a [2]int\nvar b [3]int\n// a = b // 编译错误",
				ExpectedOutput: "",
				IsValid:        false,
				RuleRef:        "TR-ARRAY-LEN",
			},
			{
				ID:             "ex-array-compare",
				ConceptID:      "array",
				Title:          "可比较数组",
				Code:           "a := [2]int{1,2}\nb := [2]int{1,2}\nfmt.Println(a == b)",
				ExpectedOutput: "true",
				IsValid:        true,
				RuleRef:        "TR-ARRAY-LEN",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-ARRAY-LEN", ConceptID: "array", RuleType: "identity", Description: "数组长度属于类型，必须相同才能赋值或比较"},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-array-1",
				ConceptID:   "array",
				Stem:        "数组长度是否影响赋值兼容性？",
				Options:     []string{"是，长度不同不兼容", "否，只看元素类型"},
				Answer:      "A",
				Explanation: "长度属于类型定义。",
				RuleRef:     "TR-ARRAY-LEN",
				Difficulty:  "easy",
			},
		},
	})
}

// ArrayOutline 返回数组主题提纲。
func ArrayOutline() []string {
	return []string{
		"数组长度是类型一部分",
		"元素类型决定可比较性，长度必须相同",
		"适用场景：固定长度数据、栈上分配",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/array"),
	}
}
