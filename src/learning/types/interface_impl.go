package types

import "fmt"

// registerInterfaceImpl 注册接口实现判定内容。
func registerInterfaceImpl() {
	_ = RegisterContent(TopicInterfaceImpl, TopicContent{
		Concept: TypeConcept{
			ID:        "interface_impl",
			Category:  "interface",
			Title:     "接口实现判定",
			Summary:   "实现要求方法集覆盖，指针接收者方法仅指针类型实现。",
			GoVersion: "1.24",
			Rules:     []string{"方法集覆盖才算实现", "指针接收者方法需指针类型"},
			Keywords:  []string{"interface implement"},
			PrintableOutline: []string{
				"方法集覆盖 = 实现接口",
				"指针方法集包含值方法，值方法集不含指针方法",
				"nil 接口值 vs 带类型 nil 的区别",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-IFACE-IMPL", ConceptID: "interface_impl", RuleType: "identity", Description: "方法集覆盖判定实现"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-iface-impl",
				ConceptID:      "interface_impl",
				Title:          "指针接收者实现",
				Code:           "type S struct{}\nfunc (s *S) Read(p []byte)(int,error){ return 0,nil }\ntype R interface{ Read([]byte)(int,error) }\nvar _ R = &S{}",
				ExpectedOutput: "",
				IsValid:        true,
				RuleRef:        "TR-IFACE-IMPL",
			},
			{
				ID:             "ex-iface-nil",
				ConceptID:      "interface_impl",
				Title:          "接口值为 nil 的陷阱",
				Code:           "var r R\nfmt.Println(r == nil) // true\nvar s *S\nr = s\nfmt.Println(r == nil) // false，动态类型已存在",
				ExpectedOutput: "true\nfalse",
				IsValid:        true,
				RuleRef:        "TR-IFACE-IMPL",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-iface-impl-1",
				ConceptID:   "interface_impl",
				Stem:        "值接收者方法能否被指针类型使用以实现接口？",
				Options:     []string{"可以，指针方法集包含值方法", "不可以"},
				Answer:      "A",
				Explanation: "指针方法集包含值接收者方法，反之不成立。",
				RuleRef:     "TR-IFACE-IMPL",
				Difficulty:  "easy",
			},
			{
				ID:          "q-iface-impl-2",
				ConceptID:   "interface_impl",
				Stem:        "接口变量 r 为 nil 与 r 动态类型为 *S 但值为 nil，是否相等？",
				Options:     []string{"不相等，动态类型已存在", "相等"},
				Answer:      "A",
				Explanation: "带动态类型的接口值非 nil，即便内部指针为 nil。",
				RuleRef:     "TR-IFACE-IMPL",
				Difficulty:  "medium",
			},
		},
	})
}

// InterfaceImplOutline 返回接口实现提纲。
func InterfaceImplOutline() []string {
	return []string{
		"方法集覆盖即实现，隐式判定",
		"指针方法集包含值方法，反之不成立",
		"接口 nil 判定：类型+值同时为 nil 才是 nil",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/interface_impl"),
	}
}
