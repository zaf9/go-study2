package types

import "fmt"

// registerInterfaceEmbedded 注册接口嵌入内容。
func registerInterfaceEmbedded() {
	_ = RegisterContent(TopicInterfaceEmbedded, TopicContent{
		Concept: TypeConcept{
			ID:        "interface_embedded",
			Category:  "interface",
			Title:     "接口嵌入与方法集合并",
			Summary:   "嵌入接口会合并方法集，重复方法需签名一致。",
			GoVersion: "1.24",
			Rules:     []string{"嵌入接口方法集合并", "重复方法签名需一致"},
			Keywords:  []string{"embedded interface"},
			PrintableOutline: []string{
				"接口可嵌入其他接口，方法集合并",
				"重复方法签名必须一致",
				"嵌入不改变实现的隐式规则",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-IFACE-EMBED", ConceptID: "interface_embedded", RuleType: "identity", Description: "嵌入接口合并方法集"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-iface-embed",
				ConceptID:      "interface_embedded",
				Title:          "ReaderAt 嵌入 Reader",
				Code:           "type R interface{ Read(p []byte)(int,error) }\ntype RA interface{ R; ReadAt(p []byte, off int64)(int,error) }",
				ExpectedOutput: "",
				IsValid:        true,
				RuleRef:        "TR-IFACE-EMBED",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-iface-embed-1",
				ConceptID:   "interface_embedded",
				Stem:        "接口嵌入会如何影响方法集？",
				Options:     []string{"合并方法集", "清空方法集"},
				Answer:      "A",
				Explanation: "嵌入将被嵌接口的方法并入当前接口。",
				RuleRef:     "TR-IFACE-EMBED",
				Difficulty:  "easy",
			},
		},
	})
}

// InterfaceEmbeddedOutline 返回接口嵌入提纲。
func InterfaceEmbeddedOutline() []string {
	return []string{
		"接口可嵌入接口，方法集合并",
		"重复方法需签名一致，否则编译错误",
		"实现方需满足合并后的完整方法集",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/interface_embedded"),
	}
}
