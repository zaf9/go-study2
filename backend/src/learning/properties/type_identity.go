package properties

import "fmt"

// registerTypeIdentity 注册"类型标识"子主题内容。
func registerTypeIdentity() {
	_ = RegisterContent(TopicTypeIdentity, TopicContent{
		Concept: PropertyConcept{
			ID:        "type_identity",
			Category:  "type_properties",
			Title:     "类型标识",
			Summary:   "两个类型相同当且仅当它们的类型字面量相同，或者它们是相同的类型别名。类型定义创建新类型，与原类型不同。",
			GoVersion: "1.24",
			Rules: []string{
				"两个类型字面量相同，则类型相同",
				"类型别名 type T = base 创建的 T 与 base 类型相同",
				"type T base 创建的新类型 T 与 base 类型不同",
				"指针类型 *T1 和 *T2 相同当且仅当 T1 和 T2 相同",
				"数组类型 [N1]T1 和 [N2]T2 相同当且仅当 N1 == N2 且 T1 == T2",
				"切片类型 []T1 和 []T2 相同当且仅当 T1 == T2",
				"结构体类型相同当且仅当字段名、类型、标签顺序相同",
				"函数类型相同当且仅当参数和返回值类型相同",
				"接口类型相同当且仅当方法集相同",
				"map 类型 map[K1]V1 和 map[K2]V2 相同当且仅当 K1 == K2 且 V1 == V2",
				"channel 类型相同当且仅当元素类型相同、方向相同",
				"泛型实例化类型相同当且仅当类型参数都相同",
			},
			Keywords: []string{"类型相同", "类型标识", "类型定义", "类型别名"},
			PrintableOutline: []string{
				"类型相同的判定规则",
				"类型定义 vs 类型别名的类型标识",
				"复合类型的类型标识",
			},
		},
		Rules: []PropertyRule{
			{
				RuleID:      "TYPE-ID-001",
				ConceptID:   "type_identity",
				RuleType:    "identity",
				Description: "类型定义 type T base 创建的新类型 T 与 base 类型不同",
			},
			{
				RuleID:      "TYPE-ID-002",
				ConceptID:   "type_identity",
				RuleType:    "identity",
				Description: "类型别名 type T = base 创建的 T 与 base 类型相同",
			},
			{
				RuleID:      "TYPE-ID-003",
				ConceptID:   "type_identity",
				RuleType:    "identity",
				Description: "复合类型的类型标识递归依赖于其组成部分的类型标识",
			},
		},
		Examples: []ExampleCase{
			{
				ID:        "ex-id-001",
				ConceptID: "type_identity",
				Title:     "类型定义创建不同类型",
				Code: `type MyInt int
type YourInt int

var mi MyInt = 1
var yi YourInt = 2

// 编译错误：MyInt 和 YourInt 是不同的类型
// mi = yi  // 错误

// 需要类型转换
mi = MyInt(yi)
fmt.Println(mi)`,
				ExpectedOutput: "2",
				IsValid:        true,
				RuleRef:        "TYPE-ID-001",
			},
			{
				ID:        "ex-id-002",
				ConceptID: "type_identity",
				Title:     "类型别名类型相同",
				Code: `type MyInt = int
type YourInt = int

var mi MyInt = 1
var yi YourInt = 2

// 可以直接赋值，类型相同
mi = yi
fmt.Println(mi)`,
				ExpectedOutput: "2",
				IsValid:        true,
				RuleRef:        "TYPE-ID-002",
			},
			{
				ID:        "ex-id-003",
				ConceptID: "type_identity",
				Title:     "指针类型的类型标识",
				Code: `type MyInt int
type YourInt int

var mi MyInt = 1
var yi YourInt = 2

var pmi *MyInt = &mi
var pyi *YourInt = &yi

// 编译错误：*MyInt 和 *YourInt 是不同的类型
// pmi = pyi  // 错误

// 需要类型转换
pmi = (*MyInt)(pyi)
fmt.Println(*pmi)`,
				ExpectedOutput: "2",
				IsValid:        true,
				RuleRef:        "TYPE-ID-003",
			},
			{
				ID:        "ex-id-004",
				ConceptID: "type_identity",
				Title:     "数组类型的类型标识",
				Code: `var a1 [3]int
var a2 [3]int
var a3 [4]int

// a1 和 a2 类型相同（长度和元素类型相同）
a1 = a2

// 编译错误：a1 和 a3 类型不同（长度不同）
// a1 = a3  // 错误

fmt.Println(a1, a2)`,
				ExpectedOutput: "[0 0 0] [0 0 0]",
				IsValid:        true,
				RuleRef:        "TYPE-ID-003",
			},
			{
				ID:        "ex-id-005",
				ConceptID: "type_identity",
				Title:     "结构体类型的类型标识",
				Code: `type Point1 struct { X, Y int }
type Point2 struct { X, Y int }

// Point1 和 Point2 是不同的类型
var p1 Point1 = Point1{X: 1, Y: 2}
var p2 Point2 = Point2{X: 1, Y: 2}

// 编译错误：Point1 和 Point2 类型不同
// p1 = p2  // 错误

type Point3 struct { X, Y int }
var p3 Point3 = Point3{X: 1, Y: 2}

// 编译错误：Point1 和 Point3 类型不同
// p1 = p3  // 错误

fmt.Println(p1, p2)`,
				ExpectedOutput: "{1 2} {1 2}",
				IsValid:        true,
				RuleRef:        "TYPE-ID-001",
			},
			{
				ID:        "ex-id-006",
				ConceptID: "type_identity",
				Title:     "函数类型的类型标识",
				Code: `type Func1 func(int) int
type Func2 func(int) int

var f1 Func1 = func(x int) int { return x }
var f2 Func2 = func(x int) int { return x }

// 编译错误：Func1 和 Func2 类型不同
// f1 = f2  // 错误

type Func3 func(int) string
var f3 Func3 = func(x int) string { return "hello" }

// 编译错误：Func1 和 Func3 类型不同
// f1 = f3  // 错误

fmt.Println(f1(42))`,
				ExpectedOutput: "42",
				IsValid:        true,
				RuleRef:        "TYPE-ID-001",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:        "q-id-1",
				ConceptID: "type_identity",
				Stem:      "类型定义 type T int 创建的类型 T 与 int 的关系是什么？",
				Options: []string{
					"类型相同",
					"类型不同",
					"T 是 int 的别名",
					"编译错误",
				},
				Answer:      "B",
				Explanation: "类型定义 type T int 创建的新类型 T 与 int 类型不同，它们是两个不同的类型，不能直接赋值。",
				RuleRef:     "TYPE-ID-001",
				Difficulty:  "easy",
			},
			{
				ID:        "q-id-2",
				ConceptID: "type_identity",
				Stem:      "类型别名 type T = int 创建的类型 T 与 int 的关系是什么？",
				Options: []string{
					"类型相同",
					"类型不同",
					"T 是 int 的子类型",
					"编译错误",
				},
				Answer:      "A",
				Explanation: "类型别名 type T = int 创建的 T 与 int 类型相同，它们可以互相赋值。",
				RuleRef:     "TYPE-ID-002",
				Difficulty:  "easy",
			},
			{
				ID:        "q-id-3",
				ConceptID: "type_identity",
				Stem:      "以下两个数组类型是否相同？\n\nvar a1 [3]int\nvar a2 [3]int",
				Options: []string{
					"相同",
					"不同",
					"取决于元素值",
					"运行时判断",
				},
				Answer:      "A",
				Explanation: "两个数组类型相同当且仅当长度相同且元素类型相同。[3]int 和 [3]int 类型相同。",
				RuleRef:     "TYPE-ID-003",
				Difficulty:  "easy",
			},
			{
				ID:        "q-id-4",
				ConceptID: "type_identity",
				Stem:      "以下两个结构体类型是否相同？\n\ntype Point1 struct { X, Y int }\ntype Point2 struct { X, Y int }",
				Options: []string{
					"相同",
					"不同",
					"取决于字段值",
					"运行时判断",
				},
				Answer:      "B",
				Explanation: "类型定义创建新类型，即使结构体字段相同，Point1 和 Point2 也是不同的类型。",
				RuleRef:     "TYPE-ID-001",
				Difficulty:  "medium",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "类型相同",
				ConceptID: "type_identity",
				Summary:   "两个类型相同当且仅当它们的类型字面量相同，或者它们是相同的类型别名。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/type_identity",
					"cli":  "properties > type_identity",
				},
			},
			{
				Keyword:   "类型标识",
				ConceptID: "type_identity",
				Summary:   "类型标识判断两个类型是否相同。类型定义创建新类型，类型别名不创建新类型。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/type_identity",
					"cli":  "properties > type_identity",
				},
			},
			{
				Keyword:   "类型定义",
				ConceptID: "type_identity",
				Summary:   "type T base 创建新类型，T 与 base 类型不同。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/type_identity",
					"cli":  "properties > type_identity",
				},
			},
		},
	})
}

// TypeIdentityOutline 返回"类型标识"子主题的提纲。
func TypeIdentityOutline() []string {
	return []string{
		"type T base 创建新类型，T 与 base 类型不同",
		"type T = base 创建类型别名，T 与 base 类型相同",
		"复合类型的类型标识递归依赖于组成部分",
		"数组、指针、切片、结构体、函数、接口的相同性判定",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/properties/type_identity"),
	}
}
