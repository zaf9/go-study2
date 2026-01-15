package properties

import "fmt"

// registerUnderlyingType 注册"底层类型"子主题内容。
func registerUnderlyingType() {
	_ = RegisterContent(TopicUnderlyingType, TopicContent{
		Concept: PropertyConcept{
			ID:        "underlying_type",
			Category:  "type_properties",
			Title:     "底层类型",
			Summary:   "每个类型都有底层类型。类型定义创建新类型，其底层类型是定义中指定的类型；类型别名的底层类型就是别名指定的类型。",
			GoVersion: "1.24",
			Rules: []string{
				"基本类型（int, float64, string 等）的底层类型是自身",
				"类型定义 type T base 创建新类型，底层类型是 base 的底层类型",
				"类型别名 type T = base 的底层类型是 base 的底层类型",
				"指针类型 *T 的底层类型是 *T 的底层类型（递归）",
				"数组类型 [N]T 的底层类型是 [N]T 的底层类型（递归）",
				"结构体类型的底层类型是自身",
				"函数类型的底层类型是自身",
				"接口类型的底层类型是自身",
				"map、channel 类型的底层类型是自身",
				"类型参数的底层类型是其类型约束的底层类型",
			},
			Keywords: []string{"底层类型", "type", "类型定义", "类型别名"},
			PrintableOutline: []string{
				"底层类型推导规则",
				"类型定义 vs 类型别名",
				"复合类型的底层类型推导",
			},
		},
		Rules: []PropertyRule{
			{
				RuleID:      "TYPE-UNDER-001",
				ConceptID:   "underlying_type",
				RuleType:    "derivation",
				Description: "基本类型（预声明类型和类型字面量）的底层类型是自身",
			},
			{
				RuleID:      "TYPE-UNDER-002",
				ConceptID:   "underlying_type",
				RuleType:    "derivation",
				Description: "类型定义 type T base 创建新类型，T 的底层类型是 base 的底层类型",
			},
			{
				RuleID:      "TYPE-UNDER-003",
				ConceptID:   "underlying_type",
				RuleType:    "derivation",
				Description: "类型别名 type T = base 的底层类型是 base 的底层类型",
			},
		},
		Examples: []ExampleCase{
			{
				ID:        "ex-under-001",
				ConceptID: "underlying_type",
				Title:     "基本类型的底层类型",
				Code: `// int 的底层类型是 int
// string 的底层类型是 string
// float64 的底层类型是 float64
var x int
var y string
fmt.Println(x, y)`,
				ExpectedOutput: "0 ",
				IsValid:        true,
				RuleRef:        "TYPE-UNDER-001",
			},
			{
				ID:        "ex-under-002",
				ConceptID: "underlying_type",
				Title:     "类型定义的底层类型",
				Code: `type MyInt int
// MyInt 的底层类型是 int

type Text string
// Text 的底层类型是 string

type MyIntSlice []int
// MyIntSlice 的底层类型是 []int

var mi MyInt = 42
var t Text = "hello"
var mis MyIntSlice = []int{1, 2, 3}
fmt.Println(mi, t, mis)`,
				ExpectedOutput: "42 hello [1 2 3]",
				IsValid:        true,
				RuleRef:        "TYPE-UNDER-002",
			},
			{
				ID:        "ex-under-003",
				ConceptID: "underlying_type",
				Title:     "多层类型定义的底层类型",
				Code: `type MyInt int
type YourInt MyInt
// YourInt 的底层类型是 int（MyInt 的底层类型）

type YourYourInt YourInt
// YourYourInt 的底层类型是 int

var yi YourInt = 1
var yyi YourYourInt = 2
// 不能直接赋值：yi = yyi 编译错误
// 需要类型转换：yi = YourInt(yyi)
fmt.Println(yi, yyi)`,
				ExpectedOutput: "1 2",
				IsValid:        true,
				RuleRef:        "TYPE-UNDER-002",
			},
			{
				ID:        "ex-under-004",
				ConceptID: "underlying_type",
				Title:     "类型别名的底层类型",
				Code: `type MyInt = int
// MyInt 的底层类型是 int
// MyInt 和 int 是相同的类型

type Text = string
// Text 的底层类型是 string
// Text 和 string 是相同的类型

var mi MyInt = 42
var i int = mi  // 可以直接赋值，不需要转换
fmt.Println(mi, i)`,
				ExpectedOutput: "42 42",
				IsValid:        true,
				RuleRef:        "TYPE-UNDER-003",
			},
			{
				ID:        "ex-under-005",
				ConceptID: "underlying_type",
				Title:     "指针类型的底层类型",
				Code: `type MyInt int
type MyIntPtr *MyInt
// MyIntPtr 的底层类型是 *int

var x int = 42
var p MyIntPtr = (*MyInt)(&x)
fmt.Println(*p)`,
				ExpectedOutput: "42",
				IsValid:        true,
				RuleRef:        "TYPE-UNDER-002",
			},
			{
				ID:        "ex-under-006",
				ConceptID: "underlying_type",
				Title:     "数组类型的底层类型",
				Code: `type MyInt int
type MyIntArray [3]MyInt
// MyIntArray 的底层类型是 [3]int

type MyIntArray2 MyIntArray
// MyIntArray2 的底层类型是 [3]int

var arr MyIntArray = [3]MyInt{1, 2, 3}
var arr2 MyIntArray2 = arr
// 编译错误：类型不同，即使底层类型相同
// arr2 = arr  // 错误
// 需要转换：arr2 = MyIntArray2(arr)
fmt.Println(arr, arr2)`,
				ExpectedOutput: "[1 2 3] [1 1 1]",
				IsValid:        false,
				RuleRef:        "TYPE-UNDER-002",
				Notes:          []string{"示例中的 arr2 初始化为零值，不能直接赋值"},
			},
		},
		QuizItems: []QuizItem{
			{
				ID:        "q-under-1",
				ConceptID: "underlying_type",
				Stem:      "类型定义 type T int 创建的新类型 T 的底层类型是什么？",
				Options: []string{
					"T",
					"int",
					"interface{}",
					"any",
				},
				Answer:      "B",
				Explanation: "类型定义 type T int 创建的新类型 T 的底层类型是 int。",
				RuleRef:     "TYPE-UNDER-002",
				Difficulty:  "easy",
			},
			{
				ID:        "q-under-2",
				ConceptID: "underlying_type",
				Stem:      "类型别名 type T = int 创建的 T 的底层类型是什么？",
				Options: []string{
					"T",
					"int",
					"interface{}",
					"any",
				},
				Answer:      "B",
				Explanation: "类型别名 type T = int 的底层类型是 int。实际上 T 和 int 是完全相同的类型。",
				RuleRef:     "TYPE-UNDER-003",
				Difficulty:  "easy",
			},
			{
				ID:        "q-under-3",
				ConceptID: "underlying_type",
				Stem:      "以下代码中 YourInt 的底层类型是什么？\n\ntype MyInt int\ntype YourInt MyInt",
				Options: []string{
					"YourInt",
					"MyInt",
					"int",
					"interface{}",
				},
				Answer:      "C",
				Explanation: "YourInt 的底层类型是 MyInt 的底层类型，即 int。类型定义会递归地查找底层类型。",
				RuleRef:     "TYPE-UNDER-002",
				Difficulty:  "medium",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "底层类型",
				ConceptID: "underlying_type",
				Summary:   "每个类型都有底层类型。类型定义创建新类型，其底层类型是定义中指定的类型的底层类型。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/underlying_type",
					"cli":  "properties > underlying_type",
				},
			},
			{
				Keyword:   "类型定义",
				ConceptID: "underlying_type",
				Summary:   "type T base 创建新类型，T 与 base 是不同的类型，但 T 的底层类型是 base 的底层类型。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/underlying_type",
					"cli":  "properties > underlying_type",
				},
			},
			{
				Keyword:   "类型别名",
				ConceptID: "underlying_type",
				Summary:   "type T = base 创建类型别名，T 和 base 是相同的类型，可以互相赋值。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/underlying_type",
					"cli":  "properties > underlying_type",
				},
			},
		},
	})
}

// UnderlyingTypeOutline 返回"底层类型"子主题的提纲。
func UnderlyingTypeOutline() []string {
	return []string{
		"基本类型的底层类型是自身",
		"type T base: T 的底层类型是 base 的底层类型",
		"type T = base: T 的底层类型是 base 的底层类型",
		"递归推导：多层类型定义的底层类型",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/properties/underlying_type"),
	}
}
