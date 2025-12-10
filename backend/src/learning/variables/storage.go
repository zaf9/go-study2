package variables

import "strings"

func init() {
	registerStorageContent()
	registerStorageQuiz()
}

func registerStorageContent() {
	snippet := strings.Join([]string{
		"var a int            // 声明，零值为0",
		"b := new(int)       // new 分配，返回 *int，指向零值",
		"c := &Point{X: 1}   // 复合字面量取址，c 为 *Point",
		"c.Y = 2             // 可通过指针修改字段",
	}, "\n")
	examples := []Example{
		{
			Title:  "声明与零值",
			Code:   "var n int\nvar s string\nfmt.Printf(\"%d|%q\", n, s)",
			Output: "0|\"\"",
			Notes: []string{
				"显式声明未赋值，按类型零值初始化。",
			},
		},
		{
			Title:  "new 返回指针",
			Code:   "p := new(int)\nfmt.Println(*p)\n*p = 7\nfmt.Println(*p)",
			Output: "0\n7",
			Notes: []string{
				"new 为基础类型分配内存并返回指针，指向零值。",
			},
		},
		{
			Title:  "复合字面量取址",
			Code:   "type Point struct{ X, Y int }\nc := &Point{X: 1}\nc.Y = 2\nfmt.Println(*c)",
			Output: "{1 2}",
			Notes: []string{
				"对结构体复合字面量取址，获得可修改的指针值。",
			},
		},
	}
	RegisterContent(TopicStorage, Content{
		Topic:   TopicStorage,
		Title:   "变量存储与取址",
		Summary: "对变量的声明、new 分配与复合字面量取址进行对比，理解零值与可寻址性。",
		Details: []string{
			"声明会生成零值，基础类型零值对应其自然默认值。",
			"new(T) 分配并返回 *T，指向零值，可直接解引用赋值。",
			"对结构体等复合字面量取址可获得可修改指针，常用于初始化并就地修改字段。",
			"数组/切片/结构体字段的元素同样遵循零值规则，未显式赋值即为零值。",
		},
		Examples: examples,
		Snippet:  snippet,
	})
}

func registerStorageQuiz() {
	items := []QuizItem{
		{
			ID:    "q-storage-1",
			Topic: TopicStorage,
			Stem:  "关于 new(T) 返回值，下列说法正确的是？",
			Options: []string{
				"A. 返回 T 类型并已初始化为非零值",
				"B. 返回 *T 指针，指向类型零值",
				"C. 返回 interface{} 需断言后使用",
				"D. 返回 uintptr 表示地址",
			},
			Answer:      "B",
			Explanation: "new 分配零值并返回指向该零值的 *T 指针。",
		},
		{
			ID:    "q-storage-2",
			Topic: TopicStorage,
			Stem:  "对结构体复合字面量取址的主要意义是？",
			Options: []string{
				"A. 必须取址才能访问字段",
				"B. 生成可修改的指针值，便于就地更新字段",
				"C. 只能在全局变量中使用",
				"D. 会绕过零值规则",
			},
			Answer:      "B",
			Explanation: "取址获得指针，可在初始化后直接修改字段；零值规则仍然适用。",
		},
		{
			ID:    "q-storage-3",
			Topic: TopicStorage,
			Stem:  "以下哪项描述了数组/切片/结构体元素的零值？",
			Options: []string{
				"A. 元素不会自动初始化",
				"B. 元素遵循各自类型的零值规则",
				"C. 仅结构体字段有零值",
				"D. 需要手动赋值后才可寻址",
			},
			Answer:      "B",
			Explanation: "复合类型的元素同样在未赋值时为各自类型的零值，且可寻址。",
		},
	}
	RegisterQuiz(TopicStorage, items)
}
