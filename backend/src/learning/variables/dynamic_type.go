package variables

import "strings"

func init() {
	registerDynamicContent()
	registerDynamicQuiz()
}

func registerDynamicContent() {
	snippet := strings.Join([]string{
		"var r io.Reader          // 静态类型 io.Reader",
		"r = &bytes.Buffer{}      // 动态类型 *bytes.Buffer",
		"r = (*bytes.Buffer)(nil) // 动态类型为 *bytes.Buffer，值为 nil",
		"var e error = nil        // 动态类型与值均为 nil，接口为零值",
	}, "\n")
	examples := []Example{
		{
			Title:  "接口动态类型变化",
			Code:   "var r io.Reader\nr = strings.NewReader(\"go\")\nfmt.Printf(\"%T\", r)\nr = &bytes.Buffer{}\nfmt.Printf(\" %T\", r)",
			Output: "*strings.Reader *bytes.Buffer",
			Notes: []string{
				"接口的静态类型固定为 io.Reader，动态类型随赋值变化。",
			},
		},
		{
			Title:  "nil 接口与带类型 nil",
			Code:   "var r io.Reader = (*bytes.Buffer)(nil)\nif r == nil { fmt.Println(\"nil\") } else { fmt.Println(\"non-nil\") }\nvar e error = nil\nfmt.Println(e == nil)",
			Output: "non-nil\ntrue",
			Notes: []string{
				"接口持有带静态类型的 nil 值时，接口本身非 nil；接口零值需动态类型与值都为 nil。",
			},
		},
	}
	RegisterContent(TopicDynamic, Content{
		Topic:   TopicDynamic,
		Title:   "接口动态类型与 nil",
		Summary: "接口的静态类型固定，动态类型随赋值变化；带类型的 nil 会使接口非零值，需要区分接口零值与带类型 nil。",
		Details: []string{
			"接口变量的静态类型由声明确定，动态类型由赋值的具体类型决定。",
			"当接口保存某具体类型的 nil 值时，接口本身不为 nil。",
			"接口零值需要动态类型与内部值均为 nil。",
		},
		Examples: examples,
		Snippet:  snippet,
	})
}

func registerDynamicQuiz() {
	items := []QuizItem{
		{
			ID:    "q-dynamic-1",
			Topic: TopicDynamic,
			Stem:  "接口变量的动态类型何时确定？",
			Options: []string{
				"A. 在编译期固定",
				"B. 每次赋值时更新为具体类型",
				"C. 仅在使用 type switch 时确定",
				"D. 由垃圾回收器决定",
			},
			Answer:      "B",
			Explanation: "接口动态类型在每次赋值时随被赋值的具体类型更新。",
		},
		{
			ID:    "q-dynamic-2",
			Topic: TopicDynamic,
			Stem:  "为何 `var r io.Reader = (*bytes.Buffer)(nil)` 不为 nil？",
			Options: []string{
				"A. 因为 bytes.Buffer 不支持 nil",
				"B. 接口包含动态类型 *bytes.Buffer，接口非零值",
				"C. Go 不允许接口为 nil",
				"D. 需要手动置为 nil 才生效",
			},
			Answer:      "B",
			Explanation: "接口存有动态类型信息和 nil 值，对比全零值接口。",
		},
	}
	RegisterQuiz(TopicDynamic, items)
}
