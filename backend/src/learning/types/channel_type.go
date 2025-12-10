package types

import "fmt"

// registerChannel 注册 channel 内容与测验。
func registerChannel() {
	_ = RegisterContent(TopicChannel, TopicContent{
		Concept: TypeConcept{
			ID:        "channel",
			Category:  "composite",
			Title:     "通道方向与缓冲",
			Summary:   "chan 可带方向限制，缓冲区决定发送阻塞行为。",
			GoVersion: "1.24",
			Rules:     []string{"双向 chan 可转换为单向", "容量影响发送阻塞"},
			Keywords:  []string{"chan", "buffer"},
			PrintableOutline: []string{
				"双向 chan 可转换为只读/只写",
				"len/ cap 反映缓冲区使用情况",
				"关闭通道后发送将 panic，接收立即返回零值",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-CHAN-DIR", ConceptID: "channel", RuleType: "direction", Description: "chan 可限制为只读/只写"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-chan-dir",
				ConceptID:      "channel",
				Title:          "只读通道",
				Code:           "ch := make(chan int,1)\nvar recv <-chan int = ch\n_ = recv",
				ExpectedOutput: "",
				IsValid:        true,
				RuleRef:        "TR-CHAN-DIR",
			},
			{
				ID:             "ex-chan-close",
				ConceptID:      "channel",
				Title:          "关闭后接收零值",
				Code:           "ch := make(chan int,1)\nch <- 1\nclose(ch)\nv, ok := <-ch\nfmt.Println(v, ok)",
				ExpectedOutput: "1 false",
				IsValid:        true,
				RuleRef:        "TR-CHAN-DIR",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-chan-1",
				ConceptID:   "channel",
				Stem:        "chan int 能否直接赋值给 <-chan int？",
				Options:     []string{"可以", "不可以"},
				Answer:      "A",
				Explanation: "双向 chan 可赋值给只读通道。",
				RuleRef:     "TR-CHAN-DIR",
				Difficulty:  "easy",
			},
		},
	})
}

// ChannelOutline 返回通道提纲。
func ChannelOutline() []string {
	return []string{
		"通道可限制方向，双向可转单向",
		"缓冲容量决定发送是否阻塞",
		"关闭后发送 panic，接收返回零值和 false",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/channel"),
	}
}
