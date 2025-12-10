package types

import "fmt"

// registerNumeric 注册数值类型内容与测验。
func registerNumeric() {
	_ = RegisterContent(TopicNumeric, TopicContent{
		Concept: TypeConcept{
			ID:        "numeric",
			Category:  "basic",
			Title:     "数值类型与别名",
			Summary:   "整数、浮点、复数组成数值族，byte/ rune 为整数别名，零值为 0。",
			GoVersion: "1.24",
			Rules:     []string{"整数/浮点/复数均为值类型", "零值为 0", "不同位宽不可隐式赋值"},
			Keywords:  []string{"int", "float", "complex"},
			PrintableOutline: []string{
				"整型/浮点/复数的零值均为 0",
				"不同位宽需要显式转换",
				"运算结果类型遵循操作数类型或类型提升规则",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-NUM-WIDTH", ConceptID: "numeric", RuleType: "comparability", Description: "不同位宽不可直接赋值或比较"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-num-cast",
				ConceptID:      "numeric",
				Title:          "显式转换整型",
				Code:           "var a int32 = 10\nb := int64(a)\nfmt.Println(b)",
				ExpectedOutput: "10",
				IsValid:        true,
				RuleRef:        "TR-NUM-WIDTH",
			},
			{
				ID:             "ex-num-mismatch",
				ConceptID:      "numeric",
				Title:          "位宽不匹配的比较",
				Code:           "var a int32 = 1\nvar b int64 = 1\n// _ = (a == b) // 编译错误，需转换",
				ExpectedOutput: "",
				IsValid:        false,
				RuleRef:        "TR-NUM-WIDTH",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-num-1",
				ConceptID:   "numeric",
				Stem:        "int32 能否直接赋值给 int64？",
				Options:     []string{"不能，需显式转换", "可以，自动扩展"},
				Answer:      "A",
				Explanation: "不同位宽需转换。",
				RuleRef:     "TR-NUM-WIDTH",
				Difficulty:  "easy",
			},
			{
				ID:          "q-num-2",
				ConceptID:   "numeric",
				Stem:        "float64 与 int 相加结果类型？",
				Options:     []string{"编译错误，需转换", "结果为 float64"},
				Answer:      "A",
				Explanation: "不同基础类型需显式转换后运算。",
				RuleRef:     "TR-NUM-WIDTH",
				Difficulty:  "easy",
			},
		},
	})
}

// NumericOutline 返回数值主题提纲。
func NumericOutline() []string {
	return []string{
		"整型/浮点/复数零值为 0",
		"不同位宽与类型间赋值需显式转换",
		"byte=rune=整数别名，保持可比较性",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/numeric"),
	}
}
