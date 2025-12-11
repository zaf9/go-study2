package progress

// topicChapterOrder 定义主题内章节的固定顺序，确保进度计算与“继续学习”提示一致。
var topicChapterOrder = map[string][]string{
	"lexical_elements": {
		"comments",
		"tokens",
		"semicolons",
		"identifiers",
		"keywords",
		"operators",
		"integers",
		"floats",
		"imaginary",
		"runes",
		"strings",
	},
	"constants": {
		"boolean",
		"rune",
		"integer",
		"floating_point",
		"complex",
		"string",
		"expressions",
		"typed_untyped",
		"conversions",
		"builtin_functions",
		"iota",
		"implementation_restrictions",
	},
	"variables": {
		"storage",
		"static",
		"dynamic",
		"zero",
	},
	"types": {
		"boolean",
		"numeric",
		"string",
		"array",
		"slice",
		"struct",
		"pointer",
		"function",
		"interface_basic",
		"interface_embedded",
		"interface_general",
		"interface_impl",
		"map",
		"channel",
	},
}

// topicDisplayNames 定义主题的人类可读名称。
var topicDisplayNames = map[string]string{
	"lexical_elements": "Lexical Elements",
	"constants":        "Constants",
	"variables":        "Variables",
	"types":            "Types",
}

// defaultTopicOrder 提供确定性的主题排序，权重相同场景使用。
var defaultTopicOrder = []string{
	"lexical_elements",
	"constants",
	"variables",
	"types",
}

// defaultChapterTotals 返回主题章节总数，供进度汇总计算。
func defaultChapterTotals() map[string]int {
	result := map[string]int{}
	for topic, chapters := range topicChapterOrder {
		result[topic] = len(chapters)
	}
	return result
}

// topicName 返回主题名称，找不到时回退为主题 ID。
func topicName(topic string) string {
	if name, ok := topicDisplayNames[topic]; ok {
		return name
	}
	return topic
}
