package types

import "fmt"

// registerInterfaceBasic 注册接口基础内容。
func registerInterfaceBasic() {
	_ = RegisterContent(TopicInterfaceBasic, TopicContent{
		Concept: TypeConcept{
			ID:        "interface_basic",
			Category:  "interface",
			Title:     "接口基础定义",
			Summary:   "接口是方法集合，零值为 nil，满足鸭子类型。",
			GoVersion: "1.24",
			Rules:     []string{"方法集合决定实现", "零值 nil"},
			Keywords:  []string{"interface", "method set"},
			PrintableOutline: []string{
				"接口由方法集定义，隐式实现",
				"零值 nil，可与 nil 比较",
				"类型实现需覆盖全部方法集",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-IFACE-METHOD", ConceptID: "interface_basic", RuleType: "identity", Description: "接口由方法集合定义"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-iface-basic",
				ConceptID:      "interface_basic",
				Title:          "隐式实现",
				Code:           "type Reader interface{ Read(p []byte)(int,error) }\ntype MyR struct{}\nfunc (MyR) Read(p []byte)(int,error){ return 0,nil }\nvar r Reader = MyR{}",
				ExpectedOutput: "",
				IsValid:        true,
				RuleRef:        "TR-IFACE-METHOD",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-iface-1",
				ConceptID:   "interface_basic",
				Stem:        "Go 接口需要显式关键字 implements 吗？",
				Options:     []string{"不需要，隐式满足", "需要 implements 声明"},
				Answer:      "A",
				Explanation: "Go 采用结构性类型系统，隐式实现。",
				RuleRef:     "TR-IFACE-METHOD",
				Difficulty:  "easy",
			},
		},
	})
}

// InterfaceBasicOutline 返回接口基础提纲。
func InterfaceBasicOutline() []string {
	return []string{
		"接口由方法集定义，隐式实现",
		"零值 nil，可与 nil 比较",
		"实现要求方法集覆盖",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/interface_basic"),
	}
}
