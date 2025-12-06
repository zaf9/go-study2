package constants

import (
	"strings"
	"testing"
)

// TestDisplayString 测试字符串常量显示函数
func TestDisplayString(t *testing.T) {
	content := GetStringContent()

	// 验证标题
	if !strings.Contains(content, "String Constants") {
		t.Error("内容应包含标题 'String Constants'")
	}

	// 验证包含关键概念
	keyConcepts := []string{
		"解释型",
		"原始字符串",
		"反引号",
		"UTF-8",
		"不可变",
	}
	for _, key := range keyConcepts {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}

	// 验证示例数量 (应该 >= 4)
	exampleCount := strings.Count(content, "【示例")
	if exampleCount < 4 {
		t.Errorf("应至少包含 4 个示例，实际只有 %d 个", exampleCount)
	}

	codeTokens := []string{
		"Hello\\nWorld",
		"`Hello\\nWorld`",
		"greeting + \", \" + name",
		"len(str)",
	}
	for _, token := range codeTokens {
		if !strings.Contains(content, token) {
			t.Errorf("内容应包含代码Token: %s", token)
		}
	}
}

// TestGetStringContent_NotEmpty 测试返回内容非空
func TestGetStringContent_NotEmpty(t *testing.T) {
	content := GetStringContent()
	if len(content) == 0 {
		t.Error("GetStringContent 不应返回空内容")
	}
}
