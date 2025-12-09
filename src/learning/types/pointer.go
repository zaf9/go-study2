package types

import "fmt"

// registerPointer 注册指针内容与测验。
func registerPointer() {
	_ = RegisterContent(TopicPointer, TopicContent{
		Concept: TypeConcept{
			ID:        "pointer",
			Category:  "composite",
			Title:     "指针与零值 nil",
			Summary:   "指针零值为 nil，可比较；解引用前需非 nil。",
			GoVersion: "1.24",
			Rules:     []string{"零值 nil 可比较", "解引用需非 nil"},
			Keywords:  []string{"pointer", "nil"},
			PrintableOutline: []string{
				"零值 nil，可与 nil 比较",
				"解引用前需判空",
				"指针方法集与值方法集关系",
			},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-ptr-nil",
				ConceptID:      "pointer",
				Title:          "指针零值比较",
				Code:           "var p *int\nfmt.Println(p == nil)",
				ExpectedOutput: "true",
				IsValid:        true,
				RuleRef:        "TR-PTR-NIL",
			},
			{
				ID:             "ex-ptr-deref",
				ConceptID:      "pointer",
				Title:          "解引用前判空",
				Code:           "var p *int\n// fmt.Println(*p) // 运行时 panic",
				ExpectedOutput: "",
				IsValid:        false,
				RuleRef:        "TR-PTR-NIL",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-PTR-NIL", ConceptID: "pointer", RuleType: "identity", Description: "指针零值 nil，可比较"},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-ptr-1",
				ConceptID:   "pointer",
				Stem:        "指针零值是否可与 nil 比较？",
				Options:     []string{"可以", "不可以"},
				Answer:      "A",
				Explanation: "指针是可比较类型，零值为 nil。",
				RuleRef:     "TR-PTR-NIL",
				Difficulty:  "easy",
			},
		},
	})
}

// PointerOutline 返回指针提纲。
func PointerOutline() []string {
	return []string{
		"指针零值为 nil，可比较",
		"解引用前需判空，避免 panic",
		"指针方法集包含值接收者方法",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/pointer"),
	}
}
