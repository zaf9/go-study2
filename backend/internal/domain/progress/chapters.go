package progress

import "strings"

// TopicChapterOrder 导出主题内章节的固定顺序，确保进度计算与"继续学习"提示一致。
var TopicChapterOrder = map[string][]string{
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
	"properties": {
		"representation",
		"underlying_type",
		"core_type",
		"type_identity",
		"assignability",
		"representability",
		"method_set",
	},
}

// ChapterDisplayName 返回章节的显示名称，将章节 ID 格式化为可读名称。
// 例如 "storage" -> "Storage", "floating_point" -> "Floating Point"
func ChapterDisplayName(chapter string) string {
	parts := strings.Split(chapter, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, " ")
}
