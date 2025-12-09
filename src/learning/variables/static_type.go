package variables

import "strings"

func init() {
	registerStaticContent()
	registerStaticQuiz()
}

func registerStaticContent() {
	snippet := strings.Join([]string{
		"var x int = 1        // 显式静态类型 int",
		"y := 2               // 推断为 int",
		"var r io.Reader      // 接口静态类型 io.Reader",
		"r = strings.NewReader(\"go\")",
	}, "\n")
	examples := []Example{
		{
			Title:  "静态类型决定可赋值性",
			Code:   "var x int\nvar y int32\n// x = y // 编译错误：int32 不能赋给 int\nx = int(y)",
			Output: "",
			Notes: []string{
				"静态类型在编译期决定可赋值性，不匹配需显式转换。",
			},
		},
		{
			Title:  "推断的静态类型",
			Code:   "s := \"go\"\nvar r io.Reader\nr = strings.NewReader(s)\n_ = r",
			Output: "",
			Notes: []string{
				"短变量声明会推断静态类型；接口静态类型决定可接受的具体类型集合。",
			},
		},
	}
	RegisterContent(TopicStatic, Content{
		Topic:   TopicStatic,
		Title:   "静态类型与可赋值性",
		Summary: "静态类型在编译期固定，决定可赋值性与方法集，推断与显式声明都形成确定的静态类型。",
		Details: []string{
			"显式声明或推断都会得到确定的静态类型，编译期检查可赋值性。",
			"接口的静态类型约束可接受的具体类型，方法集需满足接口。",
			"需要跨类型赋值时使用显式转换，避免隐式丢失或不兼容。",
		},
		Examples: examples,
		Snippet:  snippet,
	})
}

func registerStaticQuiz() {
	items := []QuizItem{
		{
			ID:    "q-static-1",
			Topic: TopicStatic,
			Stem:  "静态类型的主要作用是什么？",
			Options: []string{
				"A. 决定运行时分配的内存大小",
				"B. 决定编译期的可赋值性与方法集",
				"C. 仅影响调试输出",
				"D. 仅在接口类型中生效",
			},
			Answer:      "B",
			Explanation: "静态类型在编译期检查可赋值性和可用方法。",
		},
		{
			ID:    "q-static-2",
			Topic: TopicStatic,
			Stem:  "以下哪种情况需要显式类型转换？",
			Options: []string{
				"A. int32 赋给 int",
				"B. 同一类型的变量赋值",
				"C. 将字符串赋给 interface{}",
				"D. 将 nil 赋给接口变量",
			},
			Answer:      "A",
			Explanation: "不同的具体数字类型需要显式转换，避免隐式不兼容。",
		},
	}
	RegisterQuiz(TopicStatic, items)
}
