package lexical_elements

import "fmt"

// DisplayTokens 展示并解释 Go 语言中的词法单元 (Tokens)。
// 在 Go 语言中，源代码被分解成一系列的“词法单元”或“符号”。
// 这些是构成语言的基本元素。
func DisplayTokens() {
	fmt.Println("\n--- Go 语言的词法单元 (Tokens) ---")
	fmt.Println("词法单元是编译器能够理解的最小语法单位。Go 语言的源代码在编译前会被“词法分析器”分解为一个个的词法单元。")
	fmt.Println("Go 语言主要有四种类型的词法单元：")

	// 1. 标识符 (Identifiers)
	// 例如：main, DisplayTokens, variable
	fmt.Println("\n1. 标识符 (Identifiers):")
	fmt.Println("   用于命名程序实体，如变量、函数、类型等。例如 `myVar`, `calculateTotal`。")

	// 2. 关键字 (Keywords)
	// Go 语言有25个关键字，例如：if, for, func, package, return
	fmt.Println("\n2. 关键字 (Keywords):")
	fmt.Println("   语言预定义的、有特殊含义的单词。例如 `if`, `for`, `func`。")

	// 3. 操作符和标点 (Operators and Punctuation)
	// 例如：+, -, *, /, =, ==, !=, (), {}, []
	fmt.Println("\n3. 操作符和标点 (Operators and Punctuation):")
	fmt.Println("   用于执行操作或分隔代码结构。例如 `+`, `*`, `(`, `)`。")

	// 4. 字面量 (Literals)
	// 表示固定值的符号，例如：100, 3.14, "hello world"
	fmt.Println("\n4. 字面量 (Literals):")
	fmt.Println("   表示固定值的文本。例如整数 `100`，字符串 `\"Hello\"`。")

	fmt.Println("\n--- 示例 ---")
	fmt.Println("对于代码行 `price := 100 + 50`，词法分析器会将其分解为以下词法单元：")
	fmt.Println("1. `price` (标识符)")
	fmt.Println("2. `:=` (操作符)")
	fmt.Println("3. `100` (字面量)")
	fmt.Println("4. `+` (操作符)")
	fmt.Println("5. `50` (字面量)")
	fmt.Println("后续的章节将分别详细介绍这些不同类型的词法单元。")
}
