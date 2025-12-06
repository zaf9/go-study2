package constants

import (
	"strings"
	"testing"
)

// TestDisplayComplex 测试复数常量显示函数
func TestDisplayComplex(t *testing.T) {
	content := GetComplexContent()

	// 验证标题
	if !strings.Contains(content, "Complex Constants") {
		t.Error("内容应包含标题 'Complex Constants'")
	}

	// 验证包含关键概念
	keyConcepts := []string{
		"实部",
		"虚部",
		"i",
		"complex128",
	}
	for _, key := range keyConcepts {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}

	// 验证示例数量 (应该 >= 3)
	exampleCount := strings.Count(content, "【示例")
	if exampleCount < 3 {
		t.Errorf("应至少包含 3 个示例，实际只有 %d 个", exampleCount)
	}

	codeTokens := []string{
		"1 + 2i",
		"5i",
		"real(z)",
	}
	for _, token := range codeTokens {
		if !strings.Contains(content, token) {
			t.Errorf("内容应包含代码Token: %s", token)
		}
	}
}

// TestGetComplexContent_NotEmpty 测试返回内容非空
func TestGetComplexContent_NotEmpty(t *testing.T) {
	content := GetComplexContent()
	if len(content) == 0 {
		t.Error("GetComplexContent 不应返回空内容")
	}
}
