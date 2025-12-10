package constants

import (
	"strings"
	"testing"
)

// TestDisplayInteger 测试整数常量显示函数
func TestDisplayInteger(t *testing.T) {
	content := GetIntegerContent()

	// 验证标题
	if !strings.Contains(content, "Integer Constants") {
		t.Error("内容应包含标题 'Integer Constants'")
	}

	// 验证包含四种进制
	bases := []string{"十进制", "二进制", "八进制", "十六进制"}
	for _, base := range bases {
		if !strings.Contains(content, base) {
			t.Errorf("内容应包含%s说明", base)
		}
	}

	// 验证示例数量 (应该 >= 5)
	exampleCount := strings.Count(content, "【示例")
	if exampleCount < 5 {
		t.Errorf("应至少包含 5 个示例，实际只有 %d 个", exampleCount)
	}

	// 验证包含关键概念
	keyConcepts := []string{
		"任意精度",
		"下划线",
		"无类型",
		"溢出",
	}
	for _, key := range keyConcepts {
		if !strings.Contains(content, key) {
			t.Errorf("内容应包含关键词: %s", key)
		}
	}

	// 验证代码片段特定的 token (检查是否包含了编写的代码)
	codeTokens := []string{
		"1_000_000_000",
		"0b1100100",
		"0o777",
		"0x1A",
		"Huge / 1e10",
		"1 << 100",
	}
	for _, token := range codeTokens {
		if !strings.Contains(content, token) {
			t.Errorf("内容应包含代码Token: %s", token)
		}
	}
}

// TestGetIntegerContent_NotEmpty 测试返回内容非空
func TestGetIntegerContent_NotEmpty(t *testing.T) {
	content := GetIntegerContent()
	if len(content) == 0 {
		t.Error("GetIntegerContent 不应返回空内容")
	}
}
