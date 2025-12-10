package types

import "fmt"

// registerString 注册字符串内容与测验。
func registerString() {
	_ = RegisterContent(TopicString, TopicContent{
		Concept: TypeConcept{
			ID:        "string",
			Category:  "basic",
			Title:     "字符串与不可变性",
			Summary:   "字符串是只读字节序列，零值为空字符串，不可就地修改。",
			GoVersion: "1.24",
			Rules:     []string{"不可变", "切片结果仍为字符串", "len 返回字节数"},
			Keywords:  []string{"string", "immutable"},
			PrintableOutline: []string{
				"零值为空字符串 \"\"",
				"字符串不可变，修改需转换为 []byte/[]rune",
				"切片与索引返回字节，len 为字节长度",
			},
		},
		Rules: []TypeRule{
			{RuleID: "TR-STR-IMMUTABLE", ConceptID: "string", RuleType: "identity", Description: "字符串不可就地修改"},
		},
		Examples: []ExampleCase{
			{
				ID:             "ex-str-slice",
				ConceptID:      "string",
				Title:          "字符串切片",
				Code:           "s := \"go\" \nfmt.Println(s[0:1])",
				ExpectedOutput: "g",
				IsValid:        true,
				RuleRef:        "TR-STR-IMMUTABLE",
			},
			{
				ID:             "ex-str-mutate",
				ConceptID:      "string",
				Title:          "就地修改错误",
				Code:           "s := \"go\"\n// s[0] = 'G' // 编译错误：字符串不可变",
				ExpectedOutput: "",
				IsValid:        false,
				RuleRef:        "TR-STR-IMMUTABLE",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:          "q-str-1",
				ConceptID:   "string",
				Stem:        "字符串可否通过索引直接修改？",
				Options:     []string{"不可修改", "可以修改"},
				Answer:      "A",
				Explanation: "字符串不可变，需转为 []byte 或 []rune 后修改。",
				RuleRef:     "TR-STR-IMMUTABLE",
				Difficulty:  "easy",
			},
			{
				ID:          "q-str-2",
				ConceptID:   "string",
				Stem:        "len(\"你好\") 的结果是？",
				Options:     []string{"2", "4", "6"},
				Answer:      "C",
				Explanation: "UTF-8 下每个汉字 3 字节，len 返回字节数。",
				RuleRef:     "TR-STR-IMMUTABLE",
				Difficulty:  "medium",
			},
		},
	})
}

// StringOutline 返回字符串主题提纲。
func StringOutline() []string {
	return []string{
		"字符串不可变，零值为空字符串",
		"索引/切片按字节，len 返回字节数",
		"修改需转换为可变字节/符文切片",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/types/string"),
	}
}
