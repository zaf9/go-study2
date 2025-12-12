// Package constants - 整数常量学习模块
//
// 本文件介绍 Go 语言中的整数常量(Integer Constants)。
// 整数常量可以表示任意精度的整数值，支持多种进制表示。
package constants

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-study2/internal/infrastructure/logger"
)

// GetIntegerContent 返回整数常量相关的学习内容
func GetIntegerContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Integer Constants (整数常量) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("整数常量表示整数值序列。Go 语言的整数常量是无类型的，\n")
	sb.WriteString("并且具有任意精度（仅受限于编译器实现，通常至少 256 位）。\n")
	sb.WriteString("这意味着它们不会溢出，直到赋值给具体的有类型变量。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("支持四种表示法:\n")
	sb.WriteString("1. 十进制: 123 (不能以 0 开头，除非是 0 本身)\n")
	sb.WriteString("2. 二进制: 0b1011 (以 0b 或 0B 开头)\n")
	sb.WriteString("3. 八进制: 0o777 (以 0o 或 0O 开头，旧式写法 0777 也支持但易混淆)\n")
	sb.WriteString("4. 十六进制: 0x1A (以 0x 或 0X 开头)\n")
	sb.WriteString("可以使用下划线 _ 增加可读性: 1_000_000\n\n")

	// 示例 1: 不同进制表示
	sb.WriteString("【示例 1: 不同进制表示】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        d = 100         // 十进制\n")
	sb.WriteString("        b = 0b1100100   // 二进制 (100)\n")
	sb.WriteString("        o = 0o144       // 八进制 (100)\n")
	sb.WriteString("        x = 0x64        // 十六进制 (100)\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"%d %d %d %d\\n\", d, b, o, x)\n")
	sb.WriteString("    // 输出: 100 100 100 100\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 2: 大整数和可读性
	sb.WriteString("【示例 2: 大整数和可读性】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const big = 1_000_000_000\n")
	sb.WriteString("    const flags = 0b1111_0000_1010_0011\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(\"Big:\", big)\n")
	sb.WriteString("    fmt.Printf(\"Flags: %b\\n\", flags)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 3: 任意精度特性
	sb.WriteString("【示例 3: 任意精度特性】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 这是一个超出 int64 范围的巨大数值\n")
	sb.WriteString("    // constant 100000000000000000000 overflows int\n")
	sb.WriteString("    // 但是作为常量定义是合法的，只要不赋值给不够大的变量\n")
	sb.WriteString("    const Huge = 1e20 \n")
	sb.WriteString("    \n")
	sb.WriteString("    // 用于表达式计算不会溢出\n")
	sb.WriteString("    const Result = Huge / 1e10\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(\"Result:\", Result) // 输出: Result: 1e+10\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 4: 所有的整数常量本质上都是无类型的
	sb.WriteString("【示例 4: 无类型常量的灵活性】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const i = 42\n")
	sb.WriteString("    \n")
	sb.WriteString("    var f float64 = i  // 自动转换为 float64\n")
	sb.WriteString("    var u uint8 = i    // 自动转换为 uint8\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"%T(%v), %T(%v)\\n\", f, f, u, u)\n")
	sb.WriteString("    // 输出: float64(42), uint8(42)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 5: 位运算
	sb.WriteString("【示例 5: 位运算】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        Read    = 1 << 0  // 1\n")
	sb.WriteString("        Write   = 1 << 1  // 2\n")
	sb.WriteString("        Execute = 1 << 2  // 4\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    const RW = Read | Write\n")
	sb.WriteString("    fmt.Printf(\"RW mask: %03b\\n\", RW) // 输出: 011\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 八进制混淆:\n")
	sb.WriteString("   const n = 010   // 这是八进制的 8，不是十进制的 10\n")
	sb.WriteString("   建议使用 0o 前缀避免歧义 (如 0o10)。\n")
	sb.WriteString("\n")
	sb.WriteString("2. 赋值溢出:\n")
	sb.WriteString("   const huge = 1 << 100\n")
	sb.WriteString("   var i int = huge // 错误：constant 126... overflows int\n")
	sb.WriteString("   // 常量本身可以很大，但赋值给变量时必须在变量类型范围内\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayInteger 展示并解释 Go 语言中的整数常量。
func DisplayInteger() {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logger.LogBiz(context.Background(), "DisplayInteger", map[string]interface{}{
			"operation": "display_integer_constants",
			"result":    "success",
		}, nil, duration)
	}()

	fmt.Print(GetIntegerContent())
}
