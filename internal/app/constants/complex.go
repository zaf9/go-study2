// Package constants - 复数常量学习模块
//
// 本文件介绍 Go 语言中的复数常量(Complex Constants)。
// 复数常量由实部和虚部组成，形式上通常表示为 a + bi。
package constants

import (
	"fmt"
	"strings"
)

// GetComplexContent 返回复数常量相关的学习内容
func GetComplexContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Complex Constants (复数常量) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("复数常量包含实部和虚部。在 Go 中，通过在浮点数或整数后加 `i` 来表示虚部。\n")
	sb.WriteString("形式如: 1.2 + 3.4i。实部和虚部都是浮点数。\n")
	sb.WriteString("如果省略实部，则实部默认为 0。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("1. 基本形式: 1 + 2i\n")
	sb.WriteString("2. 纯虚数: 2i (等同于 0 + 2i)\n")
	sb.WriteString("3. 实部虚部可以是整数或浮点数表示: 1.5 + 2.5i\n")
	sb.WriteString("4. 虚数单位 `i` 紧跟在数字后面，中间不能有空格\n\n")

	// 示例 1: 基本复数声明
	sb.WriteString("【示例 1: 基本复数声明】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const a = 1 + 2i\n")
	sb.WriteString("    const b = 5i  // 纯虚数\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"a: %v (Type: %T)\\n\", a, a)\n")
	sb.WriteString("    fmt.Printf(\"b: %v (Type: %T)\\n\", b, b)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 2: 复数运算
	sb.WriteString("【示例 2: 复数运算】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const c1 = 1 + 2i\n")
	sb.WriteString("    const c2 = 3 - 4i\n")
	sb.WriteString("    \n")
	sb.WriteString("    const sum = c1 + c2\n")
	sb.WriteString("    const product = c1 * c2\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(\"Sum:\", sum)      // 输出: Sum: (4-2i)\n")
	sb.WriteString("    fmt.Println(\"Product:\", product) // 输出: Product: (11+2i)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 3: 实部虚部提取
	sb.WriteString("【示例 3: 实部虚部提取】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 注意：real() 和 imag() 函数只能用于变量或类型化常量\n")
	sb.WriteString("    // 对于无类型复数常量，需要先转换或赋值\n")
	sb.WriteString("    const c = 3 + 4i\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 或者直接使用内置函数 complex 构建\n")
	sb.WriteString("    // const built = complex(5, 6)\n")
	sb.WriteString("    \n")
	sb.WriteString("    var z complex128 = c\n")
	sb.WriteString("    fmt.Printf(\"Real: %f, Imag: %f\\n\", real(z), imag(z))\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 虚数单位写法:\n")
	sb.WriteString("   const c = 1 + 2*i // 错误: i 被视为变量名\n")
	sb.WriteString("   const c = 1 + 2i  // 正确\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayComplex 展示并解释 Go 语言中的复数常量。
func DisplayComplex() {
	fmt.Print(GetComplexContent())
}
