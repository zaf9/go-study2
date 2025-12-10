package constants

import (
	"strings"
	"testing"
)

// TestDisplayRune 测试符文常量显示函数
func TestDisplayRune(t *testing.T) {
	content := GetRuneContent()

	// 验证标题
	if !strings.Contains(content, "Rune Constants") {
		t.Error("内容应包含标题 'Rune Constants'")
	}

	// 验证两个主要概念
	keyConcepts := []string{
		"Unicode",
		"int32",
		"单引号",
	}
	for _, key := range keyConcepts {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含对应概念: %s", key)
		}
	}

	// 验证示例代码 presence
	if !strings.Contains(content, "【示例 1: 基本符文常量】") {
		t.Error("应包含基本符文常量示例")
	}

	// 验证转义序列说明
	escapes := []string{"\\u", "\\U", "\\x"}
	for _, esc := range escapes {
		if !strings.Contains(content, esc) {
			t.Errorf("内容应该解释转义序列: %s", esc)
		}
	}

	// 验证常见错误
	if !strings.Contains(content, "常见错误") {
		t.Error("内容应包含常见错误部分")
	}
}

// TestGetRuneContent_NotEmpty 测试返回内容非空
func TestGetRuneContent_NotEmpty(t *testing.T) {
	content := GetRuneContent()
	if len(content) == 0 {
		t.Error("GetRuneContent 不应返回空内容")
	}
}
