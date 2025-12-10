// Package constants - 浮点常量学习模块
//
// 本文件介绍 Go 语言中的浮点常量(Floating-point Constants)。
// 浮点常量表示小数，支持普通小数写法和科学计数法。
package constants

import (
	"fmt"
	"strings"
)

// GetFloatingPointContent 返回浮点常量相关的学习内容
func GetFloatingPointContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Floating-point Constants (浮点常量) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("浮点常量用于表示包含小数部分的数值。Go 的浮点常量也是无类型的，\n")
	sb.WriteString("并且具有非常高的精度（通常至少 256 位）。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("1. 小数表示法: 1.23, .45, 1.\n")
	sb.WriteString("2. 科学计数法: 1.23e6 (表示 1.23 * 10^6), 1E-3\n")
	sb.WriteString("   - 指数部分 (e 或 E) 是必须的，如果使用科学计数法\n")
	sb.WriteString("   - 必须包含小数点或指数部分，以区别于整数常量\n\n")

	// 示例 1: 基本浮点表示
	sb.WriteString("【示例 1: 基本浮点表示】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        Pi    = 3.14159\n")
	sb.WriteString("        Zero  = 0.0\n")
	sb.WriteString("        Half  = .5\n")
	sb.WriteString("        One   = 1.\n")
	sb.WriteString("    )\n")
	sb.WriteString("    fmt.Printf(\"%v, %v, %v, %v\\n\", Pi, Zero, Half, One)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 2: 科学计数法
	sb.WriteString("【示例 2: 科学计数法】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const Avogadro = 6.02214076e23\n")
	sb.WriteString("    const Planck = 6.62607015e-34\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"Avogadro: %e\\n\", Avogadro)\n")
	sb.WriteString("    fmt.Printf(\"Planck: %E\\n\", Planck)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 3: 高精度特性
	sb.WriteString("【示例 3: 高精度计算】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 浮点常量具有极高的精度，计算过程中不会丢失精度\n")
	sb.WriteString("    const Ln2 = 0.693147180559945309417232121458\n")
	sb.WriteString("    const Log2E = 1 / Ln2\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 只有赋值给 float32 或 float64 时才会截断\n")
	sb.WriteString("    var f32 float32 = Ln2\n")
	sb.WriteString("    var f64 float64 = Ln2\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"Full: %.30f\\n\", Ln2)\n")
	sb.WriteString("    fmt.Printf(\"f64:  %.30f\\n\", f64)\n")
	sb.WriteString("    fmt.Printf(\"f32:  %.30f\\n\", f32)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 4: 浮点数表达式
	sb.WriteString("【示例 4: 浮点数表达式】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        SecondsPerMinute = 60.0\n")
	sb.WriteString("        MinutesPerHour   = 60.0\n")
	sb.WriteString("        SecondsPerHour   = SecondsPerMinute * MinutesPerHour\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 结果也是常量\n")
	sb.WriteString("    fmt.Println(\"Seconds/Hour:\", SecondsPerHour)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 整数除法意外截断:\n")
	sb.WriteString("   const half = 1 / 2      // 结果是 0 (整数除法)\n")
	sb.WriteString("   const realHalf = 1.0 / 2.0 // 结果是 0.5 (浮点除法)\n")
	sb.WriteString("\n")
	sb.WriteString("2. 精度丢失:\n")
	sb.WriteString("   浮点常量本身精度很高，但赋值给 float32/float64 变量时可能丢失精度。\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayFloatingPoint 展示并解释 Go 语言中的浮点常量。
func DisplayFloatingPoint() {
	fmt.Print(GetFloatingPointContent())
}
