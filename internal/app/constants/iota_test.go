package constants

import (
	"strings"
	"testing"
)

// TestDisplayIota 测试 iota 特性显示函数
func TestDisplayIota(t *testing.T) {
	content := GetIotaContent()

	// 验证标题
	if !strings.Contains(content, "Iota") {
		t.Error("内容应包含标题 'Iota'")
	}

	// 验证包含概念说明
	if !strings.Contains(content, "概念说明") {
		t.Error("内容应包含概念说明部分")
	}

	// 验证包含语法规则
	if !strings.Contains(content, "语法规则") {
		t.Error("内容应包含语法规则部分")
	}

	// 验证至少有 5 个示例
	exampleCount := strings.Count(content, "【示例")
	if exampleCount < 5 {
		t.Errorf("应至少包含 5 个示例,实际只有 %d 个", exampleCount)
	}

	// 验证包含常见错误说明
	if !strings.Contains(content, "常见错误") {
		t.Error("内容应包含常见错误说明")
	}

	// 验证包含关键概念
	keyContent := []string{
		"iota",
		"枚举",
		"位掩码",
		"自增",
		"Sunday",
		"KB",
		"Read",
		"表达式复用",
	}
	for _, key := range keyContent {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}
}

// TestGetIotaContent_NotEmpty 测试返回内容非空
func TestGetIotaContent_NotEmpty(t *testing.T) {
	content := GetIotaContent()
	if len(content) == 0 {
		t.Error("GetIotaContent 不应返回空内容")
	}
}
