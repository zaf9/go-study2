package constants

import (
	"strings"
	"testing"
)

// TestDisplayConversions 测试常量类型转换显示函数
func TestDisplayConversions(t *testing.T) {
	content := GetConversionsContent()

	// 验证标题
	if !strings.Contains(content, "Conversions") {
		t.Error("内容应包含标题 'Conversions'")
	}

	// 验证包含概念说明
	if !strings.Contains(content, "概念说明") {
		t.Error("内容应包含概念说明部分")
	}

	// 验证包含语法规则
	if !strings.Contains(content, "语法规则") {
		t.Error("内容应包含语法规则部分")
	}

	// 验证至少有 4 个示例
	exampleCount := strings.Count(content, "【示例")
	if exampleCount < 4 {
		t.Errorf("应至少包含 4 个示例,实际只有 %d 个", exampleCount)
	}

	// 验证包含常见错误说明
	if !strings.Contains(content, "常见错误") {
		t.Error("内容应包含常见错误说明")
	}

	// 验证包含关键概念
	keyContent := []string{
		"类型转换",
		"可表示性",
		"representability",
		"int8",
		"float32",
		"float64",
		"complex64",
		"精度损失",
	}
	for _, key := range keyContent {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}
}

// TestGetConversionsContent_NotEmpty 测试返回内容非空
func TestGetConversionsContent_NotEmpty(t *testing.T) {
	content := GetConversionsContent()
	if len(content) == 0 {
		t.Error("GetConversionsContent 不应返回空内容")
	}
}
