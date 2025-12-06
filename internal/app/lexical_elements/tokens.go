package lexical_elements

import (
	"fmt"
	"strings"
)

// GetTokensContent 返回词法单元相关的学习内容
func GetTokensContent() string {
	var sb strings.Builder

	sb.WriteString("\n--- Go 语言的词法单元 (Tokens) ---\n")
	sb.WriteString("词法单元是编译器能够理解的最小语法单位。Go 语言的源代码在编译前会被“词法分析器”分解为一个个的词法单元。\n")
	sb.WriteString("Go 语言主要有四种类型的词法单元：\n")

	// 1. 标识符 (Identifiers)
	sb.WriteString("\n1. 标识符 (Identifiers):\n")
	sb.WriteString("   用于命名程序实体，如变量、函数、类型等。例如 `myVar`, `calculateTotal`。\n")

	// 2. 关键字 (Keywords)
	sb.WriteString("\n2. 关键字 (Keywords):\n")
	sb.WriteString("   语言预定义的、有特殊含义的单词。例如 `if`, `for`, `func`。\n")

	// 3. 操作符和标点 (Operators and Punctuation)
	sb.WriteString("\n3. 操作符和标点 (Operators and Punctuation):\n")
	sb.WriteString("   用于执行操作或分隔代码结构。例如 `+`, `*`, `(`, `)`。\n")

	// 4. 字面量 (Literals)
	sb.WriteString("\n4. 字面量 (Literals):\n")
	sb.WriteString("   表示固定值的文本。例如整数 `100`，字符串 `\"Hello\"`。\n")

	sb.WriteString("\n--- 示例 ---\n")
	sb.WriteString("对于代码行 `price := 100 + 50`，词法分析器会将其分解为以下词法单元：\n")
	sb.WriteString("1. `price` (标识符)\n")
	sb.WriteString("2. `:=` (操作符)\n")
	sb.WriteString("3. `100` (字面量)\n")
	sb.WriteString("4. `+` (操作符)\n")
	sb.WriteString("5. `50` (字面量)\n")
	sb.WriteString("后续的章节将分别详细介绍这些不同类型的词法单元。\n")

	return sb.String()
}

// DisplayTokens 展示并解释 Go 语言中的词法单元 (Tokens)。
// 在 Go 语言中，源代码被分解成一系列的“词法单元”或“符号”。
// 这些是构成语言的基本元素。
func DisplayTokens() {
	fmt.Print(GetTokensContent())
}
