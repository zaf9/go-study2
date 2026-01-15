package properties

import "fmt"

// registerCoreType 注册"核心类型"子主题内容。
func registerCoreType() {
	_ = RegisterContent(TopicCoreType, TopicContent{
		Concept: PropertyConcept{
			ID:        "core_type",
			Category:  "type_properties",
			Title:     "核心类型",
			Summary:   "核心类型用于类型约束和类型推断。非接口类型的核心类型是其底层类型；接口类型的核心类型基于其类型集。",
			GoVersion: "1.24",
			Rules: []string{
				"基本类型（int, float64, string 等）的核心类型是其自身",
				"数组类型的核心类型是其底层类型",
				"指针类型的核心类型是其底层类型",
				"结构体类型的所有字段类型相同且有标签时，核心类型是结构体",
				"函数类型的核心类型是其自身",
				"接口类型的核心类型是其类型集的交集的核心类型",
				"bytestring 核心类型：string 或 []byte",
				"没有核心类型的情况：结构体字段不同、接口类型集为空或无共同类型",
			},
			Keywords: []string{"核心类型", "类型约束", "bytestring", "类型集"},
			PrintableOutline: []string{
				"非接口类型的核心类型",
				"接口类型的核心类型",
				"bytestring 核心类型",
			},
		},
		Rules: []PropertyRule{
			{
				RuleID:      "TYPE-CORE-001",
				ConceptID:   "core_type",
				RuleType:    "derivation",
				Description: "基本类型和大部分类型的核心类型是其底层类型",
			},
			{
				RuleID:      "TYPE-CORE-002",
				ConceptID:   "core_type",
				RuleType:    "derivation",
				Description: "接口类型的核心类型是其类型集中所有类型共享的核心类型",
			},
			{
				RuleID:      "TYPE-CORE-003",
				ConceptID:   "core_type",
				RuleType:    "special",
				Description: "bytestring 是 string 或 []byte 的联合核心类型",
			},
		},
		Examples: []ExampleCase{
			{
				ID:        "ex-core-001",
				ConceptID: "core_type",
				Title:     "基本类型的核心类型",
				Code: `// int 的核心类型是 int
// string 的核心类型是 string
// float64 的核心类型是 float64

type MyInt int
// MyInt 的核心类型是 int

type MyIntSlice []int
// MyIntSlice 的核心类型是 []int

var mi MyInt = 42
fmt.Println(mi)`,
				ExpectedOutput: "42",
				IsValid:        true,
				RuleRef:        "TYPE-CORE-001",
			},
			{
				ID:        "ex-core-002",
				ConceptID: "core_type",
				Title:     "结构体类型的核心类型",
				Code: `type Point struct {
	X int
	Y int
}
// Point 的核心类型是 Point 本身（字段类型相同）

type Data struct {
	a int
	b string
}
// Data 没有核心类型（字段类型不同）

p := Point{X: 1, Y: 2}
fmt.Println(p)`,
				ExpectedOutput: "{1 2}",
				IsValid:        true,
				RuleRef:        "TYPE-CORE-001",
			},
			{
				ID:        "ex-core-003",
				ConceptID: "core_type",
				Title:     "接口类型的核心类型",
				Code: `type Stringer interface {
	String() string
}
// Stringer 没有核心类型（类型集包含所有实现 String() 方法的类型）

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
// Integer 的核心类型是 int（所有类型的共同核心类型）

type Any interface {
	int | float64 | string
}
// Any 没有核心类型（int、float64、string 无共同类型）`,
				IsValid: true,
				RuleRef: "TYPE-CORE-002",
			},
			{
				ID:        "ex-core-004",
				ConceptID: "core_type",
				Title:     "bytestring 核心类型",
				Code: `type Byteseq interface {
	~string | ~[]byte
}
// Byteseq 的核心类型是 bytestring
// bytestring 表示可以是 string 或 []byte

func print[T Byteseq](v T) {
	fmt.Println(v)
}

print("hello")     // string
print([]byte{1, 2}) // []byte`,
				ExpectedOutput: "hello\n[1 2]",
				IsValid:        true,
				RuleRef:        "TYPE-CORE-003",
			},
			{
				ID:        "ex-core-005",
				ConceptID: "core_type",
				Title:     "指针类型的核心类型",
				Code: `type MyInt int
type IntPtr *MyInt
// IntPtr 的核心类型是 *int

type IntPtr2 *int
// IntPtr2 的核心类型是 *int

var x int = 42
var p IntPtr = (*MyInt)(&x)
fmt.Println(*p)`,
				ExpectedOutput: "42",
				IsValid:        true,
				RuleRef:        "TYPE-CORE-001",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:        "q-core-1",
				ConceptID: "core_type",
				Stem:      "基本类型 int 的核心类型是什么？",
				Options: []string{
					"int",
					"interface{}",
					"any",
					"没有核心类型",
				},
				Answer:      "A",
				Explanation: "基本类型的核心类型是其自身。int 的核心类型是 int。",
				RuleRef:     "TYPE-CORE-001",
				Difficulty:  "easy",
			},
			{
				ID:        "q-core-2",
				ConceptID: "core_type",
				Stem:      "以下类型定义中，哪个没有核心类型？",
				Options: []string{
					"type Point struct { X, Y int }",
					"type Data struct { a int; b string }",
					"type IntPtr *int",
					"type MyInt int",
				},
				Answer:      "B",
				Explanation: "结构体类型 Data 的字段类型不同（int 和 string），因此没有核心类型。Point 的字段类型相同，核心类型是 Point 本身。",
				RuleRef:     "TYPE-CORE-001",
				Difficulty:  "medium",
			},
			{
				ID:        "q-core-3",
				ConceptID: "core_type",
				Stem:      "bytestring 核心类型表示什么？",
				Options: []string{
					"只是 string",
					"只是 []byte",
					"string 或 []byte",
					"任意字节序列",
				},
				Answer:      "C",
				Explanation: "bytestring 核心类型表示可以是 string 或 []byte。用于类型约束时，允许类型参数实例化为 string 或 []byte。",
				RuleRef:     "TYPE-CORE-003",
				Difficulty:  "medium",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "核心类型",
				ConceptID: "core_type",
				Summary:   "核心类型用于类型约束和类型推断。基本类型的核心类型是其自身。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/core_type",
					"cli":  "properties > core_type",
				},
			},
			{
				Keyword:   "bytestring",
				ConceptID: "core_type",
				Summary:   "bytestring 是 string 或 []byte 的联合核心类型，用于泛型类型约束。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/core_type",
					"cli":  "properties > core_type",
				},
			},
			{
				Keyword:   "类型集",
				ConceptID: "core_type",
				Summary:   "接口类型的类型集定义了可以满足该接口的所有类型。接口类型的核心类型是其类型集的交集的核心类型。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/core_type",
					"cli":  "properties > core_type",
				},
			},
		},
	})
}

// CoreTypeOutline 返回"核心类型"子主题的提纲。
func CoreTypeOutline() []string {
	return []string{
		"基本类型和大部分类型的核心类型是其底层类型",
		"结构体字段类型相同时，核心类型是结构体本身",
		"接口类型的核心类型是其类型集的交集的核心类型",
		"bytestring: string 或 []byte 的联合核心类型",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/properties/core_type"),
	}
}
