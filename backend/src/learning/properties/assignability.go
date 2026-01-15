package properties

import "fmt"

// registerAssignability 注册"可赋值性"子主题内容。
func registerAssignability() {
	_ = RegisterContent(TopicAssignability, TopicContent{
		Concept: PropertyConcept{
			ID:        "assignability",
			Category:  "type_properties",
			Title:     "可赋值性",
			Summary:   "值 x 可以赋值给类型 T 的变量，当且仅当满足可赋值性条件。可赋值性有 6 个基本条件和 3 个涉及类型参数的额外条件。",
			GoVersion: "1.24",
			Rules: []string{
				"条件 1：x 的类型与 T 相同",
				"条件 2：x 的类型 V 与 T 的底层类型相同，且至少有一个不是命名类型",
				"条件 3：T 是接口类型，x 实现了 T",
				"条件 4：x 是双向 channel 值，T 是 channel 类型，x 的类型 V 和 T 的元素类型相同，且至少有一个不是命名类型",
				"条件 5：x 是预声明标识符 nil，T 是指针、函数、切片、map、channel 或接口类型",
				"条件 6：x 是无类型常量，可表示为 T 类型",
				"额外条件 1：x 的类型 V 和 T 的类型参数相同的实例化",
				"额外条件 2：T 是接口类型，V 有核心类型，且 V 的核心类型实现了 T",
				"额外条件 3：x 的类型 V 和 T 的底层类型相同，忽略结构体标签",
			},
			Keywords: []string{"可赋值性", "赋值", "类型匹配"},
			PrintableOutline: []string{
				"6 个基本可赋值性条件",
				"3 个涉及类型参数的额外条件",
			},
		},
		Rules: []PropertyRule{
			{
				RuleID:      "ASSIGN-001",
				ConceptID:   "assignability",
				RuleType:    "condition",
				Description: "类型相同可以赋值",
			},
			{
				RuleID:      "ASSIGN-002",
				ConceptID:   "assignability",
				RuleType:    "condition",
				Description: "底层类型相同且至少有一个不是命名类型可以赋值",
			},
			{
				RuleID:      "ASSIGN-003",
				ConceptID:   "assignability",
				RuleType:    "condition",
				Description: "值实现了接口类型可以赋值",
			},
			{
				RuleID:      "ASSIGN-004",
				ConceptID:   "assignability",
				RuleType:    "condition",
				Description: "双向 channel 可以赋值给类型兼容的 channel",
			},
			{
				RuleID:      "ASSIGN-005",
				ConceptID:   "assignability",
				RuleType:    "condition",
				Description: "nil 可以赋值给指针、函数、切片、map、channel、接口类型",
			},
			{
				RuleID:      "ASSIGN-006",
				ConceptID:   "assignability",
				RuleType:    "condition",
				Description: "无类型常量可表示为类型 T 时可以赋值",
			},
		},
		Examples: []ExampleCase{
			{
				ID:        "ex-assign-001",
				ConceptID: "assignability",
				Title:     "条件 1：类型相同",
				Code: `var x int = 42
var y int
y = x  // 可以赋值：类型相同
fmt.Println(y)`,
				ExpectedOutput: "42",
				IsValid:        true,
				RuleRef:        "ASSIGN-001",
			},
			{
				ID:        "ex-assign-002",
				ConceptID: "assignability",
				Title:     "条件 2：底层类型相同，至少一个不是命名类型",
				Code: `type MyInt int

var x MyInt = 42
var y int
// y = x  // 编译错误：都是命名类型

type MyInt2 int
var x2 MyInt2 = 42
var y2 int
// y2 = x2  // 编译错误：都是命名类型

var z int = 42
var y3 MyInt
// y3 = z  // 可以：z 不是命名类型（int 是预声明类型）
y3 = MyInt(z)
fmt.Println(y3)`,
				ExpectedOutput: "42",
				IsValid:        true,
				RuleRef:        "ASSIGN-002",
				Notes:          []string{"int 是预声明类型，不是命名类型。MyInt 是命名类型。"},
			},
			{
				ID:        "ex-assign-003",
				ConceptID: "assignability",
				Title:     "条件 2：类型字面量可以赋值",
				Code: `var x []int = []int{1, 2, 3}
var y []int
y = x  // 可以：[]int 是类型字面量，不是命名类型
fmt.Println(y)

type MySlice []int
var z MySlice = []int{4, 5, 6}
// y = z  // 编译错误：都是命名类型
y = []int(z)  // 需要转换
fmt.Println(y)`,
				ExpectedOutput: "[1 2 3]\n[4 5 6]",
				IsValid:        true,
				RuleRef:        "ASSIGN-002",
			},
			{
				ID:        "ex-assign-004",
				ConceptID: "assignability",
				Title:     "条件 3：实现了接口",
				Code: `type Stringer interface {
	String() string
}

type MyInt int

func (m MyInt) String() string {
	return fmt.Sprintf("MyInt(%d)", m)
}

var x MyInt = 42
var s Stringer
s = x  // 可以：MyInt 实现了 Stringer
fmt.Println(s)`,
				ExpectedOutput: "MyInt(42)",
				IsValid:        true,
				RuleRef:        "ASSIGN-003",
			},
			{
				ID:        "ex-assign-005",
				ConceptID: "assignability",
				Title:     "条件 4：双向 channel",
				Code: `var x chan int = make(chan int, 1)
var y chan int
y = x  // 可以：类型相同

var z chan<- int
// z = x  // 可以：双向 channel 可以赋值给单向 channel

var w <-chan int
// w = x  // 可以：双向 channel 可以赋值给单向 channel

x <- 42
fmt.Println(<-y)`,
				ExpectedOutput: "42",
				IsValid:        true,
				RuleRef:        "ASSIGN-004",
			},
			{
				ID:        "ex-assign-006",
				ConceptID: "assignability",
				Title:     "条件 5：nil 可以赋值",
				Code: `var p *int
p = nil  // 可以

var f func()
f = nil  // 可以

var s []int
s = nil  // 可以

var m map[string]int
m = nil  // 可以

var ch chan int
ch = nil  // 可以

var i interface{}
i = nil  // 可以

fmt.Println(p, s, m, ch, i)`,
				ExpectedOutput: "<nil> [] map[] <nil> <nil>",
				IsValid:        true,
				RuleRef:        "ASSIGN-005",
			},
			{
				ID:        "ex-assign-007",
				ConceptID: "assignability",
				Title:     "条件 6：无类型常量",
				Code: `const x = 42  // 无类型常量

var i int
i = x  // 可以：42 可表示为 int

var f float64
f = x  // 可以：42 可表示为 float64

var u uint8
u = x  // 可以：42 可表示为 uint8

var b byte
b = x  // 可以：42 可表示为 byte

fmt.Println(i, f, u, b)`,
				ExpectedOutput: "42 42 42 42",
				IsValid:        true,
				RuleRef:        "ASSIGN-006",
			},
			{
				ID:        "ex-assign-008",
				ConceptID: "assignability",
				Title:     "类型参数的可赋值性",
				Code: `type Numeric interface {
	~int | ~float64
}

func double[T Numeric](x T) T {
	return x + x
}

var i int = 21
var f float64 = 1.5

fmt.Println(double(i))  // 42
fmt.Println(double(f))  // 3`,
				ExpectedOutput: "42\n3",
				IsValid:        true,
				RuleRef:        "ASSIGN-003",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:        "q-assign-1",
				ConceptID: "assignability",
				Stem:      "以下哪个赋值是合法的？\n\ntype MyInt int\nvar x MyInt = 42\nvar y int\ny = x",
				Options: []string{
					"合法",
					"不合法：都是命名类型",
					"不合法：类型不同",
					"编译错误",
				},
				Answer:      "B",
				Explanation: "MyInt 和 int 都是命名类型，虽然底层类型相同，但不能直接赋值。需要类型转换：y = int(x)。",
				RuleRef:     "ASSIGN-002",
				Difficulty:  "medium",
			},
			{
				ID:        "q-assign-2",
				ConceptID: "assignability",
				Stem:      "nil 可以赋值给哪些类型？",
				Options: []string{
					"所有类型",
					"指针、函数、切片、map、channel、接口",
					"仅接口类型",
					"仅指针类型",
				},
				Answer:      "B",
				Explanation: "nil 可以赋值给指针、函数、切片、map、channel 和接口类型的变量。",
				RuleRef:     "ASSIGN-005",
				Difficulty:  "easy",
			},
			{
				ID:        "q-assign-3",
				ConceptID: "assignability",
				Stem:      "以下哪个赋值是合法的？\n\ntype Stringer interface { String() string }\ntype MyInt int\nfunc (m MyInt) String() string { return \"42\" }\nvar x MyInt = 42\nvar s Stringer\ns = x",
				Options: []string{
					"合法：MyInt 实现了 Stringer",
					"不合法：类型不同",
					"不合法：需要类型转换",
					"编译错误",
				},
				Answer:      "A",
				Explanation: "MyInt 实现了 Stringer 接口，所以可以将 MyInt 类型的值赋值给 Stringer 类型的变量。",
				RuleRef:     "ASSIGN-003",
				Difficulty:  "medium",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "可赋值性",
				ConceptID: "assignability",
				Summary:   "值 x 可以赋值给类型 T 的变量，当且仅当满足 6 个基本条件或 3 个额外条件之一。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/assignability",
					"cli":  "properties > assignability",
				},
			},
			{
				Keyword:   "赋值",
				ConceptID: "assignability",
				Summary:   "赋值语句要求右值可赋值给左值的类型。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/assignability",
					"cli":  "properties > assignability",
				},
			},
		},
	})
}

// AssignabilityOutline 返回"可赋值性"子主题的提纲。
func AssignabilityOutline() []string {
	return []string{
		"条件 1：类型相同",
		"条件 2：底层类型相同且至少一个不是命名类型",
		"条件 3：值实现了接口",
		"条件 4：双向 channel 赋值给兼容的 channel",
		"条件 5：nil 赋值给特定类型",
		"条件 6：无类型常量可表示为类型 T",
		"额外条件：类型参数相关",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/properties/assignability"),
	}
}
