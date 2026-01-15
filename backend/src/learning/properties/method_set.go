package properties

import "fmt"

// registerMethodSet 注册"方法集"子主题内容。
func registerMethodSet() {
	_ = RegisterContent(TopicMethodSet, TopicContent{
		Concept: PropertyConcept{
			ID:        "method_set",
			Category:  "type_properties",
			Title:     "方法集",
			Summary:   "类型的方法集决定了该类型是否实现了某个接口。值接收者方法同时属于值类型和指针类型的方法集，指针接收者方法只属于指针类型的方法集。",
			GoVersion: "1.24",
			Rules: []string{
				"命名类型 T 的方法集包含所有接收者为 T 的方法",
				"指针类型 *T 的方法集包含所有接收者为 T 和 *T 的方法",
				"结构体类型的方法集包含所有接收者为该结构体类型的方法",
				"结构体指针的方法集包含所有接收者为该结构体及其指针的方法",
				"嵌入字段的方法会被提升到外层类型的方法集",
				"接口类型的方法集就是接口声明的所有方法",
			},
			Keywords: []string{"方法集", "接收者", "接口实现"},
			PrintableOutline: []string{
				"值类型的方法集",
				"指针类型的方法集",
				"嵌入字段的方法提升",
			},
		},
		Rules: []PropertyRule{
			{
				RuleID:      "METHOD-001",
				ConceptID:   "method_set",
				RuleType:    "definition",
				Description: "类型 T 的方法集包含所有接收者为 T 的方法",
			},
			{
				RuleID:      "METHOD-002",
				ConceptID:   "method_set",
				RuleType:    "definition",
				Description: "类型 *T 的方法集包含所有接收者为 T 和 *T 的方法",
			},
			{
				RuleID:      "METHOD-003",
				ConceptID:   "method_set",
				RuleType:    "embedding",
				Description: "嵌入字段的方法会被提升到外层类型的方法集",
			},
		},
		Examples: []ExampleCase{
			{
				ID:        "ex-method-001",
				ConceptID: "method_set",
				Title:     "值类型的方法集",
				Code: `type MyInt int

func (m MyInt) Value() int {
	return int(m)
}

var x MyInt = 42
fmt.Println(x.Value())  // 可以：MyInt 有 Value 方法

var p *MyInt = &x
// p.Value()  // 也可以：*MyInt 的方法集包含 MyInt 的方法
fmt.Println((*p).Value())`,
				ExpectedOutput: "42\n42",
				IsValid:        true,
				RuleRef:        "METHOD-001",
			},
			{
				ID:        "ex-method-002",
				ConceptID: "method_set",
				Title:     "指针接收者方法",
				Code: `type MyInt int

func (m *MyInt) SetValue(x int) {
	*m = MyInt(x)
}

func (m *MyInt) Value() int {
	return int(*m)
}

var x MyInt = 42
// x.Value()  // 编译错误：MyInt 没有 Value 方法
// x.SetValue(100)  // 编译错误：MyInt 没有 SetValue 方法

var p *MyInt = &x
p.Value()     // 可以：*MyInt 有 Value 方法
p.SetValue(100)  // 可以：*MyInt 有 SetValue 方法
fmt.Println(x)`,
				ExpectedOutput: "100",
				IsValid:        true,
				RuleRef:        "METHOD-002",
			},
			{
				ID:        "ex-method-003",
				ConceptID: "method_set",
				Title:     "混合接收者方法",
				Code: `type MyInt int

func (m MyInt) Value() int {
	return int(m)
}

func (m *MyInt) SetValue(x int) {
	*m = MyInt(x)
}

var x MyInt = 42
fmt.Println(x.Value())  // 可以：MyInt 有 Value 方法

x.SetValue(100)  // 可以：编译器自动取地址 (&x)
fmt.Println(x.Value())

var p *MyInt = &x
p.Value()     // 可以：*MyInt 有 Value 方法
p.SetValue(200)  // 可以：*MyInt 有 SetValue 方法
fmt.Println(*p)`,
				ExpectedOutput: "42\n100\n200",
				IsValid:        true,
				RuleRef:        "METHOD-002",
				Notes:          []string{"当值类型可以寻址时，编译器会自动取地址来调用指针接收者方法。"},
			},
			{
				ID:        "ex-method-004",
				ConceptID: "method_set",
				Title:     "接口实现与方法集",
				Code: `type Stringer interface {
	String() string
}

type MyInt int

func (m MyInt) String() string {
	return fmt.Sprintf("MyInt(%d)", m)
}

var x MyInt = 42
var s Stringer
s = x  // 可以：MyInt 实现了 Stringer（有 String 方法）
fmt.Println(s)

var p *MyInt = &x
s = p  // 可以：*MyInt 也实现了 Stringer（继承了 MyInt 的方法）
fmt.Println(s)`,
				ExpectedOutput: "MyInt(42)\nMyInt(42)",
				IsValid:        true,
				RuleRef:        "METHOD-001",
			},
			{
				ID:        "ex-method-005",
				ConceptID: "method_set",
				Title:     "指针方法与接口实现",
				Code: `type Stringer interface {
	String() string
}

type MyInt int

func (m *MyInt) String() string {
	return fmt.Sprintf("MyInt(%d)", *m)
}

var x MyInt = 42
// var s Stringer = x  // 编译错误：MyInt 没有 String 方法
// fmt.Println(s)

var p *MyInt = &x
var s Stringer = p  // 可以：*MyInt 实现了 Stringer
fmt.Println(s)

// 编译器自动取地址
var s2 Stringer = x  // 可以：编译器自动取地址 (&x)
fmt.Println(s2)`,
				ExpectedOutput: "MyInt(42)\nMyInt(42)",
				IsValid:        true,
				RuleRef:        "METHOD-002",
			},
			{
				ID:        "ex-method-006",
				ConceptID: "method_set",
				Title:     "嵌入字段的方法提升",
				Code: `type Inner struct {
	X int
}

func (i Inner) Method1() int {
	return i.X
}

func (i *Inner) Method2(x int) {
	i.X = x
}

type Outer struct {
	Inner  // 嵌入
	Y int
}

var o Outer
o.X = 10
fmt.Println(o.Method1())  // 可以：Inner 的方法提升到 Outer
o.Method2(20)
fmt.Println(o.X)

var p *Outer = &Outer{Inner: Inner{X: 30}}
fmt.Println(p.Method1())  // 可以：*Outer 也有提升的方法
p.Method2(40)
fmt.Println(p.X)`,
				ExpectedOutput: "10\n20\n30\n40",
				IsValid:        true,
				RuleRef:        "METHOD-003",
			},
			{
				ID:        "ex-method-007",
				ConceptID: "method_set",
				Title:     "接口类型的方法集",
				Code: `type Stringer interface {
	String() string
}

type Scanner interface {
	Scan() string
}

type StringScanner interface {
	Stringer
	Scanner
}

// StringScanner 的方法集包含 String() 和 Scan()
// 任何类型同时实现 String() 和 Scan() 方法就实现了 StringScanner

type MyScanner struct {
	data string
}

func (m MyScanner) String() string {
	return m.data
}

func (m *MyScanner) Scan() string {
	m.data = "scanned"
	return m.data
}

var s StringScanner = &MyScanner{data: "hello"}
fmt.Println(s.String())
fmt.Println(s.Scan())`,
				ExpectedOutput: "hello\nscanned",
				IsValid:        true,
				RuleRef:        "METHOD-001",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:        "q-method-1",
				ConceptID: "method_set",
				Stem:      "类型 *T 的方法集包含哪些方法？",
				Options: []string{
					"只有接收者为 *T 的方法",
					"只有接收者为 T 的方法",
					"接收者为 T 和 *T 的所有方法",
					"没有方法",
				},
				Answer:      "C",
				Explanation: "类型 *T 的方法集包含所有接收者为 T 和 *T 的方法。这是 Go 方法集的一个重要特性。",
				RuleRef:     "METHOD-002",
				Difficulty:  "easy",
			},
			{
				ID:        "q-method-2",
				ConceptID: "method_set",
				Stem:      "以下代码是否编译通过？\n\ntype MyInt int\nfunc (m *MyInt) Value() int { return int(*m) }\nvar x MyInt = 42\nfmt.Println(x.Value())",
				Options: []string{
					"通过：编译器自动取地址",
					"不通过：MyInt 没有 Value 方法",
					"不通过：需要显式调用 (&x).Value()",
					"运行时 panic",
				},
				Answer:      "A",
				Explanation: "当值类型可以寻址时，编译器会自动取地址来调用指针接收者方法。x.Value() 会被转换为 (&x).Value()。",
				RuleRef:     "METHOD-002",
				Difficulty:  "medium",
			},
			{
				ID:        "q-method-3",
				ConceptID: "method_set",
				Stem:      "以下代码是否编译通过？\n\ntype Stringer interface { String() string }\ntype MyInt int\nfunc (m *MyInt) String() string { return \"42\" }\nvar s Stringer = MyInt(42)\nfmt.Println(s)",
				Options: []string{
					"通过：MyInt 实现了 Stringer",
					"不通过：MyInt 没有实现 Stringer",
					"不通过：需要使用指针",
					"运行时 panic",
				},
				Answer:      "B",
				Explanation: "MyInt 类型没有 String 方法（只有 *MyInt 有），所以 MyInt 没有实现 Stringer 接口。应该使用 &MyInt(42) 或让编译器自动取地址。",
				RuleRef:     "METHOD-002",
				Difficulty:  "medium",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "方法集",
				ConceptID: "method_set",
				Summary:   "类型的方法集决定了该类型是否实现了某个接口。*T 的方法集包含 T 和 *T 的所有方法。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/method_set",
					"cli":  "properties > method_set",
				},
			},
			{
				Keyword:   "接收者",
				ConceptID: "method_set",
				Summary:   "方法的接收者可以是值类型或指针类型。指针接收者方法只能通过指针类型调用（除非编译器自动取地址）。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/method_set",
					"cli":  "properties > method_set",
				},
			},
			{
				Keyword:   "接口实现",
				ConceptID: "method_set",
				Summary:   "类型实现了接口的所有方法就实现了该接口。方法集决定了类型有哪些方法。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/method_set",
					"cli":  "properties > method_set",
				},
			},
		},
	})
}

// MethodSetOutline 返回"方法集"子主题的提纲。
func MethodSetOutline() []string {
	return []string{
		"类型 T 的方法集包含所有接收者为 T 的方法",
		"类型 *T 的方法集包含所有接收者为 T 和 *T 的方法",
		"嵌入字段的方法会被提升到外层类型的方法集",
		"方法集决定接口实现",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/properties/method_set"),
	}
}
