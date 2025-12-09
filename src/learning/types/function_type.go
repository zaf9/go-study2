package types

import "fmt"

// registerFunction 注册函数类型内容与测验。
func registerFunction() {
	_ = RegisterContent(TopicFunction, TopicContent{
		Concept: TypeConcept{
			ID:        "function",
			Category:  "composite",
			Title:     "函数类型与不可比较性",
			Summary:   "函数类型不可比较（除 nil），可作为变量或参数，零值为 nil。",
			GoVersion: "1.24",
			Rules:     []string{"函数不可比较", "零值为 nil", "可作为一等公民传递"},
			Keywords:  []string{"func type"},
			PrintableOutline: []string{
				"函数是不可比较类型，仅能与 nil 比较",
				"零值为 nil，调用前需判空",
				"可作为参数或返回值，实现回调",
			},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-func-nil",
				ConceptID:      "function",
				Title:          "函数变量零值",
				Code:           "var f func(int) int\nfmt.Println(f == nil)",
				ExpectedOutput: "true",
				IsValid:        true,
				RuleRef:        "TR-FUNC-NIL",
			},
			{
				ID:             "ex-func-compare",
				ConceptID:      "function",
				Title:          "不支持函数比较",
				Code:           "var f1 func()\nvar f2 func()\n// _ = (f1 == f2) // 编译错误",
				ExpectedOutput: "",
				IsValid:        false,
				RuleRef:        "TR-FUNC-NIL",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-FUNC-NIL", ConceptID: "function", RuleType: "comparability", Description: "函数类型仅能与 nil 比较"},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-func-1",
				ConceptID:   "function",
				Stem:        "函数类型能否直接比较是否相等？",
				Options:     []string{"只能与 nil 比较", "可以任意比较"},
				Answer:      "A",
				Explanation: "函数值只允许与 nil 比较。",
				RuleRef:     "TR-FUNC-NIL",
				Difficulty:  "easy",
			},
		},
	})
}

// FunctionOutline 返回函数类型提纲。
func FunctionOutline() []string {
	return []string{
		"函数类型不可比较，仅能与 nil 比较",
		"零值 nil，调用前需判空",
		"支持作为参数/返回值，实现回调",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/function"),
	}
}
