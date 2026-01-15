package properties

import "fmt"

// registerRepresentation 注册"值的表示"子主题内容。
func registerRepresentation() {
	_ = RegisterContent(TopicRepresentation, TopicContent{
		Concept: PropertyConcept{
			ID:        "representation",
			Category:  "value_properties",
			Title:     "值的表示",
			Summary:   "Go 中的值分为自包含值和引用值。自包含值直接包含数据，引用值包含指向底层数据的指针。",
			GoVersion: "1.24",
			Rules: []string{
				"自包含值：数组、结构体、基本类型（bool, 数值, string）",
				"引用值：指针、切片、map、channel、函数",
				"拷贝自包含值会复制整个值",
				"拷贝引用值只复制引用，不复制底层数据",
			},
			Keywords: []string{"值", "引用", "拷贝", "自包含值", "引用值"},
			PrintableOutline: []string{
				"自包含值：数组、结构体、基本类型",
				"引用值：指针、切片、map、channel、函数",
				"值拷贝语义：自包含值完全拷贝，引用值拷贝引用",
			},
		},
		Rules: []PropertyRule{
			{
				RuleID:      "VAL-REP-001",
				ConceptID:   "representation",
				RuleType:    "classification",
				Description: "自包含值直接存储数据，拷贝时复制整个值",
			},
			{
				RuleID:      "VAL-REP-002",
				ConceptID:   "representation",
				RuleType:    "classification",
				Description: "引用值存储指向底层数据的指针，拷贝时只复制指针",
			},
		},
		Examples: []ExampleCase{
			{
				ID:        "ex-rep-001",
				ConceptID: "representation",
				Title:     "自包含值拷贝 - 数组",
				Code: `a := [3]int{1, 2, 3}
b := a          // 完全拷贝数组
b[0] = 99
fmt.Println(a)  // [1 2 3] - a 不受影响
fmt.Println(b)  // [99 2 3]`,
				ExpectedOutput: "[1 2 3]\n[99 2 3]",
				IsValid:        true,
				RuleRef:        "VAL-REP-001",
				Notes:          []string{"数组是自包含值，拷贝时完全复制"},
			},
			{
				ID:        "ex-rep-002",
				ConceptID: "representation",
				Title:     "引用值拷贝 - 切片",
				Code: `a := []int{1, 2, 3}
b := a          // 拷贝切片结构（指针、长度、容量）
b[0] = 99
fmt.Println(a)  // [99 2 3] - a 被修改
fmt.Println(b)  // [99 2 3]`,
				ExpectedOutput: "[99 2 3]\n[99 2 3]",
				IsValid:        true,
				RuleRef:        "VAL-REP-002",
				Notes:          []string{"切片是引用值，拷贝后共享底层数组"},
			},
			{
				ID:        "ex-rep-003",
				ConceptID: "representation",
				Title:     "自包含值拷贝 - 结构体",
				Code: `type Point struct { X, Y int }
p1 := Point{X: 1, Y: 2}
p2 := p1         // 完全拷贝结构体
p2.X = 99
fmt.Println(p1)  // {1 2} - p1 不受影响
fmt.Println(p2)  // {99 2}`,
				ExpectedOutput: "{1 2}\n{99 2}",
				IsValid:        true,
				RuleRef:        "VAL-REP-001",
				Notes:          []string{"结构体是自包含值，拷贝时完全复制"},
			},
			{
				ID:        "ex-rep-004",
				ConceptID: "representation",
				Title:     "引用值拷贝 - map",
				Code: `m1 := map[string]int{"a": 1, "b": 2}
m2 := m1         // 拷贝 map 指针
m2["a"] = 99
fmt.Println(m1)  // map[a:99 b:2] - m1 被修改
fmt.Println(m2)  // map[a:99 b:2]`,
				ExpectedOutput: "map[a:99 b:2]\nmap[a:99 b:2]",
				IsValid:        true,
				RuleRef:        "VAL-REP-002",
				Notes:          []string{"map 是引用值，拷贝后共享底层数据"},
			},
			{
				ID:        "ex-rep-005",
				ConceptID: "representation",
				Title:     "引用值拷贝 - 指针",
				Code: `x := 42
p1 := &x
p2 := p1         // 拷贝指针
*p2 = 99
fmt.Println(x)   // 99 - x 被修改
fmt.Println(*p1) // 99
fmt.Println(*p2) // 99`,
				ExpectedOutput: "99\n99\n99",
				IsValid:        true,
				RuleRef:        "VAL-REP-002",
				Notes:          []string{"指针是引用值，拷贝后指向同一地址"},
			},
		},
		QuizItems: []QuizItem{
			{
				ID:        "q-rep-1",
				ConceptID: "representation",
				Stem:      "以下哪些类型是自包含值？",
				Options: []string{
					"数组、结构体、基本类型",
					"切片、map、channel",
					"指针、函数",
					"以上都是",
				},
				Answer:      "A",
				Explanation: "自包含值包括数组、结构体和基本类型（bool、数值、string）。它们直接包含数据，拷贝时完全复制。",
				RuleRef:     "VAL-REP-001",
				Difficulty:  "easy",
			},
			{
				ID:        "q-rep-2",
				ConceptID: "representation",
				Stem:      "拷贝一个切片时，会发生什么？",
				Options: []string{
					"完全复制底层数组",
					"只复制切片结构（指针、长度、容量），共享底层数组",
					"创建新的底层数组",
					"编译错误",
				},
				Answer:      "B",
				Explanation: "切片是引用值，拷贝时只复制切片结构（包含指向底层数组的指针、长度和容量），不复制底层数组。因此两个切片共享同一底层数组。",
				RuleRef:     "VAL-REP-002",
				Difficulty:  "easy",
			},
			{
				ID:        "q-rep-3",
				ConceptID: "representation",
				Stem:      "执行以下代码后，a 的值是什么？\n\na := [3]int{1, 2, 3}\nb := a\nb[0] = 99",
				Options: []string{
					"[99 2 3]",
					"[1 2 3]",
					"编译错误",
					"运行时 panic",
				},
				Answer:      "B",
				Explanation: "数组是自包含值，b := a 会完全复制数组。修改 b 不会影响 a。因此 a 保持为 [1 2 3]。",
				RuleRef:     "VAL-REP-001",
				Difficulty:  "medium",
			},
		},
		References: []ReferenceIndex{
			{
				Keyword:   "自包含值",
				ConceptID: "representation",
				Summary:   "自包含值直接包含数据，拷贝时完全复制。包括数组、结构体、基本类型。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/representation",
					"cli":  "properties > representation",
				},
			},
			{
				Keyword:   "引用值",
				ConceptID: "representation",
				Summary:   "引用值包含指向底层数据的指针，拷贝时只复制指针。包括指针、切片、map、channel、函数。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/representation",
					"cli":  "properties > representation",
				},
			},
			{
				Keyword:   "值拷贝",
				ConceptID: "representation",
				Summary:   "自包含值拷贝时完全复制数据，引用值拷贝时只复制引用。",
				Anchors: map[string]string{
					"http": "/api/v1/topic/properties/representation",
					"cli":  "properties > representation",
				},
			},
		},
	})
}

// RepresentationOutline 返回"值的表示"子主题的提纲。
func RepresentationOutline() []string {
	return []string{
		"自包含值：数组、结构体、基本类型（bool, 数值, string）",
		"引用值：指针、切片、map、channel、函数",
		"拷贝语义：自包含值完全拷贝，引用值拷贝引用",
		fmt.Sprintf("示例与测验：%s", "/api/v1/topic/properties/representation"),
	}
}
