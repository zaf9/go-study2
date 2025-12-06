package unit

import (
	"go-study2/internal/app/lexical_elements"
	"strings"
	"testing"
)

func TestLexicalContentGeneration(t *testing.T) {
	tests := []struct {
		name          string
		generator     func() string
		expectContent []string
	}{
		{
			name:      "Comments",
			generator: lexical_elements.GetCommentsContent,
			expectContent: []string{
				"Go 语言的注释",
				"// 这是一个单行注释",
				"/*",
			},
		},
		{
			name:      "Tokens",
			generator: lexical_elements.GetTokensContent,
			expectContent: []string{
				"Go 语言的词法单元",
				"1. 标识符",
				"2. 关键字",
			},
		},
		{
			name:      "Semicolons",
			generator: lexical_elements.GetSemicolonsContent,
			expectContent: []string{
				"分号 (Semicolons)",
				"自动分号插入规则",
			},
		},
		{
			name:      "Identifiers",
			generator: lexical_elements.GetIdentifiersContent,
			expectContent: []string{
				"标识符 (Identifiers)",
				"userName",
			},
		},
		{
			name:      "Keywords",
			generator: lexical_elements.GetKeywordsContent,
			expectContent: []string{
				"关键字 (Keywords)",
				"if, else",
				"defer",
			},
		},
		{
			name:      "Operators",
			generator: lexical_elements.GetOperatorsContent,
			expectContent: []string{
				"运算符和标点",
				"算术运算符",
				"+",
			},
		},
		{
			name:      "Integers",
			generator: lexical_elements.GetIntegersContent,
			expectContent: []string{
				"整数字面量",
				"1_000_000",
			},
		},
		{
			name:      "Floats",
			generator: lexical_elements.GetFloatsContent,
			expectContent: []string{
				"浮点数字面量",
				"3.14",
				"scientific",
			},
		},
		{
			name:      "Imaginary",
			generator: lexical_elements.GetImaginaryContent,
			expectContent: []string{
				"虚数字面量",
				"1e-3i",
			},
		},
		{
			name:      "Runes",
			generator: lexical_elements.GetRunesContent,
			expectContent: []string{
				"字符字面量",
				"'A'",
				"Unicode",
			},
		},
		{
			name:      "Strings",
			generator: lexical_elements.GetStringsContent,
			expectContent: []string{
				"字符串字面量",
				"Hello, World",
				"转义序列",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := tt.generator()
			if content == "" {
				t.Error("Returned empty content")
			}
			for _, expect := range tt.expectContent {
				// Special handling for "scientific" in Floats because we might not have that exact word if we used Chinese,
				// but "科学计数法" is usually what we wrote.
				// Let's check my implementation. I wrote "科学计数法".
				// I will fix expectation for "scientific" -> "科学计数法" or "e"
				if expect == "scientific" {
					expect = "科学计数法"
				}

				if !strings.Contains(content, expect) {
					t.Errorf("Expected content to contain %q, but it didn't", expect)
				}
			}
		})
	}
}
