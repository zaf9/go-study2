package variables

import "strings"

func init() {
	registerZeroContent()
	registerZeroQuiz()
}

func registerZeroContent() {
	snippet := strings.Join([]string{
		"var n int          // 零值 0",
		"var s string       // 零值 \"\"",
		"var p *int         // 零值 nil",
		"var arr [3]int     // 元素全为0",
		"var st struct{ X int; P *int } // 字段按类型零值填充",
	}, "\n")
	examples := []Example{
		{
			Title:  "基础类型零值",
			Code:   "var n int\nvar s string\nfmt.Printf(\"%d|%q\", n, s)",
			Output: "0|\"\"",
			Notes: []string{
				"整数零值为0，字符串零值为空串。",
			},
		},
		{
			Title:  "复合元素零值",
			Code:   "var arr [2]*int\nvar st struct{ P *int; V int }\nfmt.Println(arr[0] == nil, st.P == nil, st.V)",
			Output: "true true 0",
			Notes: []string{
				"数组、结构体字段的元素按各自类型填充零值，指针为 nil。",
			},
		},
	}
	RegisterContent(TopicZero, Content{
		Topic:   TopicZero,
		Title:   "零值与取值规则",
		Summary: "所有类型都有零值；复合类型的元素同样自动按零值填充，未赋值即可安全读取零值。",
		Details: []string{
			"未显式赋值的变量会被初始化为零值。",
			"指针、切片、映射、通道、函数、接口的零值为 nil。",
			"复合类型的内部元素（数组、结构体字段）也按零值填充。",
		},
		Examples: examples,
		Snippet:  snippet,
	})
}

func registerZeroQuiz() {
	items := []QuizItem{
		{
			ID:    "q-zero-1",
			Topic: TopicZero,
			Stem:  "以下哪个类型的零值为 nil？",
			Options: []string{
				"A. int",
				"B. string",
				"C. map[string]int",
				"D. bool",
			},
			Answer:      "C",
			Explanation: "引用类型（map/slice/chan/func/pointer/interface）零值均为 nil。",
		},
		{
			ID:    "q-zero-2",
			Topic: TopicZero,
			Stem:  "结构体字段未赋值时的取值规则是？",
			Options: []string{
				"A. 读会 panic",
				"B. 读到零值，按字段类型决定",
				"C. 自动推断为非零默认值",
				"D. 需要手动 new 后才能读",
			},
			Answer:      "B",
			Explanation: "结构体字段按类型零值填充，未赋值即可读取零值。",
		},
	}
	RegisterQuiz(TopicZero, items)
}
