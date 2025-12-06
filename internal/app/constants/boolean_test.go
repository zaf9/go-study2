package constants

import (
	"strings"
	"testing"
)

// TestDisplayBoolean 测试布尔常量显示函数
func TestDisplayBoolean(t *testing.T) {
	content := GetBooleanContent()

	// 验证标题
	if !strings.Contains(content, "Boolean Constants") {
		t.Error("内容应包含标题 'Boolean Constants'")
	}

	// 验证包含概念说明
	if !strings.Contains(content, "概念说明") {
		t.Error("内容应包含概念说明部分")
	}

	// 验证包含语法规则
	if !strings.Contains(content, "语法规则") {
		t.Error("内容应包含语法规则部分")
	}

	// 验证至少有 3 个示例
	exampleCount := strings.Count(content, "【示例")
	if exampleCount < 3 {
		t.Errorf("应至少包含 3 个示例,实际只有 %d 个", exampleCount)
	}

	// 验证包含常见错误说明
	if !strings.Contains(content, "常见错误") {
		t.Error("内容应包含常见错误说明")
	}

	// 验证包含关键概念
	keyContent := []string{
		"true",
		"false",
		"bool",
		"布尔常量",
	}
	for _, key := range keyContent {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}
}

// TestGetBooleanContent_NotEmpty 测试返回内容非空
func TestGetBooleanContent_NotEmpty(t *testing.T) {
	content := GetBooleanContent()
	if len(content) == 0 {
		t.Error("GetBooleanContent 不应返回空内容")
	}
}
