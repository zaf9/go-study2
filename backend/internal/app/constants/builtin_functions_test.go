package constants

import (
	"strings"
	"testing"
)

// TestDisplayBuiltinFunctions 测试内置函数显示函数
func TestDisplayBuiltinFunctions(t *testing.T) {
	content := GetBuiltinFunctionsContent()

	// 验证标题
	if !strings.Contains(content, "Built-in Functions") {
		t.Error("内容应包含标题 'Built-in Functions'")
	}

	// 验证包含概念说明
	if !strings.Contains(content, "概念说明") {
		t.Error("内容应包含概念说明部分")
	}

	// 验证包含语法规则
	if !strings.Contains(content, "语法规则") {
		t.Error("内容应包含语法规则部分")
	}

	// 验证至少有 6 个示例
	exampleCount := strings.Count(content, "【示例")
	if exampleCount < 6 {
		t.Errorf("应至少包含 6 个示例,实际只有 %d 个", exampleCount)
	}

	// 验证包含常见错误说明
	if !strings.Contains(content, "常见错误") {
		t.Error("内容应包含常见错误说明")
	}

	// 验证包含关键概念
	keyContent := []string{
		"内置函数",
		"min",
		"max",
		"len",
		"real",
		"imag",
		"complex",
		"unsafe.Sizeof",
	}
	for _, key := range keyContent {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}
}

// TestGetBuiltinFunctionsContent_NotEmpty 测试返回内容非空
func TestGetBuiltinFunctionsContent_NotEmpty(t *testing.T) {
	content := GetBuiltinFunctionsContent()
	if len(content) == 0 {
		t.Error("GetBuiltinFunctionsContent 不应返回空内容")
	}
}
