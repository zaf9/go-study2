package constants

import (
	"strings"
	"testing"
)

// TestDisplayFloatingPoint 测试浮点常量显示函数
func TestDisplayFloatingPoint(t *testing.T) {
	content := GetFloatingPointContent()

	// 验证标题
	if !strings.Contains(content, "Floating-point Constants") {
		t.Error("内容应包含标题 'Floating-point Constants'")
	}

	// 验证示例数量 (应该 >= 4)
	exampleCount := strings.Count(content, "【示例")
	if exampleCount < 4 {
		t.Errorf("应至少包含 4 个示例，实际只有 %d 个", exampleCount)
	}

	// 验证包含关键概念
	keyConcepts := []string{
		"科学计数法",
		"指数",
		"精度",
		"float32",
		"float64",
	}
	for _, key := range keyConcepts {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}

	// 验证代码token
	codeTokens := []string{
		"3.14159",
		"6.02214076e23",
		"1 / Ln2",
	}
	for _, token := range codeTokens {
		if !strings.Contains(content, token) {
			t.Errorf("内容应包含代码Token: %s", token)
		}
	}
}

// TestGetFloatingPointContent_NotEmpty 测试返回内容非空
func TestGetFloatingPointContent_NotEmpty(t *testing.T) {
	content := GetFloatingPointContent()
	if len(content) == 0 {
		t.Error("GetFloatingPointContent 不应返回空内容")
	}
}
