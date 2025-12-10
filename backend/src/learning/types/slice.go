package types

import "fmt"

// registerSlice 注册切片内容与测验。
func registerSlice() {
	_ = RegisterContent(TopicSlice, TopicContent{
		Concept: TypeConcept{
			ID:        "slice",
			Category:  "composite",
			Title:     "切片共享底层数组",
			Summary:   "切片是对底层数组的视图，包含指针、长度、容量，零值为 nil。",
			GoVersion: "1.24",
			Rules:     []string{"零值 nil 可安全使用 append", "截取共享底层数组"},
			Keywords:  []string{"slice", "capacity"},
			PrintableOutline: []string{
				"零值 nil，可 append 使用",
				"len/ cap 影响切片运算与扩容",
				"切片共享底层数组，修改会影响视图",
			},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-slice-append",
				ConceptID:      "slice",
				Title:          "零值切片可 append",
				Code:           "var s []int\ns = append(s, 1)\nfmt.Println(len(s), cap(s))",
				ExpectedOutput: "1 1",
				IsValid:        true,
				RuleRef:        "TR-SLICE-APPEND",
			},
			{
				ID:             "ex-slice-share",
				ConceptID:      "slice",
				Title:          "共享底层数组",
				Code:           "arr := [3]int{1,2,3}\ns1 := arr[0:2]\ns2 := arr[1:3]\ns1[1] = 99\nfmt.Println(arr, s2)",
				ExpectedOutput: "[1 99 3] [99 3]",
				IsValid:        true,
				RuleRef:        "TR-SLICE-SHARE",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-SLICE-APPEND", ConceptID: "slice", RuleType: "identity", Description: "零值切片可直接 append"},
			{RuleID: "TR-SLICE-SHARE", ConceptID: "slice", RuleType: "identity", Description: "切片共享底层数组，截取后修改会相互影响"},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-slice-1",
				ConceptID:   "slice",
				Stem:        "零值切片能否直接 append？",
				Options:     []string{"可以", "不可以"},
				Answer:      "A",
				Explanation: "零值切片可直接 append，返回新底层数组。",
				RuleRef:     "TR-SLICE-APPEND",
				Difficulty:  "easy",
			},
			{
				ID:          "q-slice-2",
				ConceptID:   "slice",
				Stem:        "两个切片引用同一底层数组，修改是否互相影响？",
				Options:     []string{"会影响", "不会影响"},
				Answer:      "A",
				Explanation: "共享底层数组，修改会透传到其他切片视图。",
				RuleRef:     "TR-SLICE-SHARE",
				Difficulty:  "easy",
			},
		},
	})
}

// SliceOutline 返回切片主题提纲。
func SliceOutline() []string {
	return []string{
		"零值 nil，可安全 append",
		"len/cap 概念，截取共享底层数组",
		"扩容可能导致底层数组重分配",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/slice"),
	}
}
