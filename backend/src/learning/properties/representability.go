package properties

import "fmt"

// registerRepresentability 注册"可表示性"子主题内容。
func registerRepresentability() {
	_ = RegisterContent(TopicRepresentability, TopicContent{
		Concept: PropertyConcept{
			ID:        "representability",
			Category:  "type_properties",
			Title:     "可表示性",
			Summary:   "常量 x 可以表示为类型 T，当 x 的值在 T 类型的值范围内。可表示性决定了无类型常量能否赋值给特定类型的变量。",
			GoVersion: "1.24",
			Rules: []string{
				"整数常量可表示为整数类型 T，当其在 T 的范围内",
				"浮点常量可表示为浮点类型 T，当其在 T 的范围内",
				"复数常量可表示为复数类型 T，当其实部和虚部都可表示为 T 的元素类型",
				"整数常量可表示为浮点类型 T，当其在 T 的精确整数范围内",
				"浮点常量不能直接表示为整数类型（需要转换）",
				"常量溢出会导致编译错误",
			},
			Keywords: []string{"可表示性", "常量", "溢出"},
			PrintableOutline: []string{
				"整数常量的可表示性",
				"浮点常量的可表示性",
				"复数常量的可表示性",
			},
		},
		Rules: []PropertyRule{
			{
				RuleID:      "REPR-001",
				ConceptID:   "representability",
				RuleType:    "condition",
				Description: "整数常量可表示为整数类型 T，当其值在 T 的范围内",
			},
			{
				RuleID:      "REPR-002",
				ConceptID:   "representability",
				RuleType:    "condition",
				Description: "浮点常量可表示为浮点类型 T，当其值在 T 的范围内",
			},
			{
				RuleID:      "REPR-003",
				ConceptID:   "representability",
				RuleType:    "condition",
				Description: "整数常量可表示为浮点类型 T，当其在 T 的精确整数范围内",
			},
		},
		Examples: []ExampleCase{
			{
				ID:        "ex-repr-001",
				ConceptID: "representability",
				Title:     "整数常量表示为整数类型",
				Code: `const x = 42

var i int
i = x  // 可以：42 在 int 范围内

var i8 int8
i8 = x  // 可以：42 在 int8 范围内 (-128 到 127)

var i16 int16
i16 = x  // 可以：42 在 int16 范围内

var u8 uint8
u8 = x  // 可以：42 在 uint8 范围内 (0 到 255)

fmt.Println(i, i8, i16, u8)`,
				ExpectedOutput: "42 42 42 42",
				IsValid:        true,
				RuleRef:        "REPR-001",
			},
			{
				ID:        "ex-repr-002",
				ConceptID: "representability",
				Title:     "整数常量溢出",
				Code: `const x = 300

var i8 int8
// i8 = x  // 编译错误：300 超出 int8 范围 (-128 到 127)

var i16 int16
i16 = x  // 可以：300 在 int16 范围内

fmt.Println(i16)

const y = -100
var u8 uint8
// u8 = y  // 编译错误：-100 超出 uint8 范围 (0 到 255)

var i int
i = y  // 可以：-100 在 int 范围内
fmt.Println(i)`,
				ExpectedOutput: "300\n-100",
				IsValid:        true,
				RuleRef:        "REPR-001",
			},
			{
				ID:        "ex-repr-003",
				ConceptID: "representability",
				Title:     "浮点常量表示为浮点类型",
				Code: `const x = 3.14

var f32 float32
f32 = x  // 可以：3.14 可表示为 float32

var f64 float64
f64 = x  // 可以：3.14 可表示为 float64

fmt.Println(f32, f64)

const y = 1.7976931348623157e+308  // 接近 float64 最大值
var f64_2 float64
f64_2 = y  // 可以：在 float64 范围内

const z = 1.8e+309  // 超出 float64 范围
// var f64_3 float64 = z  // 编译错误：常量溢出`,
				ExpectedOutput: "3.14 3.14",
				IsValid:        true,
				RuleRef:        "REPR-002",
			},
			{
				ID:        "ex-repr-004",
				ConceptID: "representability",
				Title:     "整数常量表示为浮点类型",
				Code: `const x = 42

var f32 float32
f32 = x  // 可以：42 在 float32 精确整数范围内

var f64 float64
f64 = x  // 可以：42 在 float64 精确整数范围内

fmt.Println(f32, f64)

const y = 9007199254740993  // 超过 float64 精确整数范围 (2^53)
var f64_2 float64
f64_2 = y  // 可以赋值，但可能损失精度
fmt.Println(f64_2)  // 9007199254740992 (精度损失)`,
				ExpectedOutput: "42 42\n9007199254740992",
				IsValid:        true,
				RuleRef:        "REPR-003",
			},
			{
				ID:        "ex-repr-005",
				ConceptID: "representability",
				Title:     "复数常量的可表示性",
				Code: `const x = 3 + 4i

var c64 complex64
c64 = x  // 可以：实部 3 和虚部 4 都在 float32 范围内

var c128 complex128
c128 = x  // 可以：实部 3 和虚部 4 都在 float64 范围内

fmt.Println(c64, c128)

const y = 1.8e+309 + 2i
// var c128_2 complex128 = y  // 编译错误：实部超出范围`,
				ExpectedOutput: "(3+4i) (3+4i)",
				IsValid:        true,
				RuleRef:        "REPR-002",
			},
			{
				ID:        "ex-repr-006",
				ConceptID: "representability",
				Title:     "类型范围的边界",
				Code: `const maxInt8 = 127
const minInt8 = -128

var i8 int8
i8 = maxInt8  // 可以：127 在 int8 范围内
i8 = minInt8  // 可以：-128 在 int8 范围内

const overflowInt8 = 128
// i8 = overflowInt8  // 编译错误：128 超出 int8 范围

const maxUint8 = 255
const minUint8 = 0

var u8 uint8
u8 = maxUint8  // 可以：255 在 uint8 范围内
u8 = minUint8  // 可以：0 在 uint8 范围内

const overflowUint8 = 256
// u8 = overflowUint8  // 编译错误：256 超出 uint8 范围

fmt.Println(i8, u8)`,
				ExpectedOutput: "-128 255",
				IsValid:        true,
				RuleRef:        "REPR-001",
			},
		},
		QuizItems: []QuizItem{
			{
				ID:        "q-repr-1",
				ConceptID: "representability",
				Stem:      "常量 300 可以表示为以下哪些类型？",
				Options: []string{
					"int8",
					"uint8",
					"int16",
					"以上都可以",
				},
				Answer:      "C",
				Explanation: "300 超出 int8 范围 (-128 到 127) 和 uint8 范围 (0 到 255)，但在 int16 范围内 (-32768 到 32767)。",
				RuleRef:     "REPR-001",
				Difficulty:  "medium",
			},
			{
				ID:        "q-repr-2",
				ConceptID: "representability",
				Stem:      "整数常量 42 可以表示为 float32 类型吗？",
				Options: []string{
					"可以：42 在 float32 精确整数范围内",
					"不可以：整数不能表示为浮点类型",
					"可以：但会损失精度",
					"编译错误",
				},
				Answer:      "A",
				Explanation: "整数常量 42 在 float32 的精确整数范围内，可以完美表示为 float32。",
				RuleRef:     "REPR-003",
				Difficulty:  "easy",
			},
			{
				ID:        "q-repr-3",
				ConceptID: "representability",
				Stem:      "以下哪个常量赋值会导致编译错误？\n\nconst x = 128\nvar i8 int8 = x\nconst y = -1\nvar u8 uint8 = y",
				Options: []string{
					"只有第一个",
					"只有第二个",
					"两个都会",
					"两个都不会",
				},
				Answer:      "C",
				Explanation: "128 超出 int8 范围 (最大 127)，-1 超出 uint8 范围 (最小 0)。两个都会导致编译错误。",
				RuleRef:     "REPR-001",
				Difficulty:  "medium",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "可表示性",
				ConceptID: "representability",
				Summary:   "常量 x 可以表示为类型 T，当 x 的值在 T 类型的值范围内。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/representability",
					"cli":  "properties > representability",
				},
			},
			{
				Keyword:   "常量溢出",
				ConceptID: "representability",
				Summary:   "常量值超出类型范围会导致编译错误。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/representability",
					"cli":  "properties > representability",
				},
			},
		},
	})
}

// RepresentabilityOutline 返回"可表示性"子主题的提纲。
func RepresentabilityOutline() []string {
	return []string{
		"整数常量可表示为整数类型 T，当其在 T 的范围内",
		"浮点常量可表示为浮点类型 T，当其在 T 的范围内",
		"整数常量可表示为浮点类型 T，当其在 T 的精确整数范围内",
		"复数常量的实部和虚部都必须可表示",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/properties/representability"),
	}
}
