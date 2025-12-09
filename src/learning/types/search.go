package types

import "strings"

var searchIndex = []ReferenceIndex{
	{
		Keyword:   "map key",
		ConceptID: "map",
		Summary:   "map 键必须可比较，如整数、字符串、指针等。",
		Anchors: map[string]string{
			"http": "/api/v1/topic/types/map",
			"cli":  "types > map",
		},
	},
	{
		Keyword:   "~int",
		ConceptID: "interface_general",
		Summary:   "~int 表示底层类型为 int 的命名类型也匹配。",
		Anchors: map[string]string{
			"http": "/api/v1/topic/types/interface_general",
			"cli":  "types > interface_general",
		},
	},
	{
		Keyword:   "interface nil",
		ConceptID: "interface_impl",
		Summary:   "接口值只有类型和值都为 nil 才为 nil，带类型 nil 不等于 nil。",
		Anchors: map[string]string{
			"http": "/api/v1/topic/types/interface_impl",
			"cli":  "types > interface_impl",
		},
	},
	{
		Keyword:   "slice share",
		ConceptID: "slice",
		Summary:   "切片共享底层数组，截取后修改会相互影响。",
		Anchors: map[string]string{
			"http": "/api/v1/topic/types/slice",
			"cli":  "types > slice",
		},
	},
	{
		Keyword:   "array length",
		ConceptID: "array",
		Summary:   "数组长度是类型一部分，长度不同不兼容赋值/比较。",
		Anchors: map[string]string{
			"http": "/api/v1/topic/types/array",
			"cli":  "types > array",
		},
	},
}

// registerSearchIndex 将固定索引写入全局注册表。
func registerSearchIndex() {
	for _, idx := range searchIndex {
		if idx.Keyword == "" {
			continue
		}
		referenceRegistry[strings.ToLower(idx.Keyword)] = idx
	}
}
