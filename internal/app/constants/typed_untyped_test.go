package constants

import (
	"strings"
	"testing"
)

// TestDisplayTypedUntyped 测试类型化/无类型化常量显示函数
func TestDisplayTypedUntyped(t *testing.T) {
	content := GetTypedUntypedContent()

	// 验证标题
	if !strings.Contains(content, "Typed and Untyped Constants") {
		t.Error("内容应包含标题 'Typed and Untyped Constants'")
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
		"无类型化常量",
		"类型化常量",
		"默认类型",
		"精度",
		"隐式转换",
		"untyped",
		"typed",
	}
	for _, key := range keyContent {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}
}

// TestGetTypedUntypedContent_NotEmpty 测试返回内容非空
func TestGetTypedUntypedContent_NotEmpty(t *testing.T) {
	content := GetTypedUntypedContent()
	if len(content) == 0 {
		t.Error("GetTypedUntypedContent 不应返回空内容")
	}
}
