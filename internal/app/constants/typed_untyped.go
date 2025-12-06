// Package constants - 类型化/无类型化常量学习模块
//
// 本文件介绍 Go 语言中类型化常量(Typed Constants)和无类型化常量(Untyped Constants)的区别。
// 理解这两种常量的区别对于掌握 Go 的类型系统非常重要。
package constants

import (
	"fmt"
	"strings"
)

// GetTypedUntypedContent 返回类型化/无类型化常量相关的学习内容
func GetTypedUntypedContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Typed and Untyped Constants (类型化/无类型化常量) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("Go 语言中的常量分为两种:\n")
	sb.WriteString("1. 无类型化常量(Untyped Constants): 字面量常量默认是无类型化的，\n")
	sb.WriteString("   具有默认类型，但可以隐式转换为兼容类型，精度不受目标类型限制。\n")
	sb.WriteString("2. 类型化常量(Typed Constants): 通过显式类型声明或类型转换得到的常量，\n")
	sb.WriteString("   类型固定，不能隐式转换，精度受类型限制。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("无类型化常量声明:\n")
	sb.WriteString("  const name = value  // 无类型化，具有默认类型\n")
	sb.WriteString("\n")
	sb.WriteString("类型化常量声明:\n")
	sb.WriteString("  const name type = value  // 显式指定类型\n")
	sb.WriteString("  const name = T(value)      // 通过类型转换\n")
	sb.WriteString("\n")
	sb.WriteString("默认类型映射:\n")
	sb.WriteString("  - 无类型布尔 → bool\n")
	sb.WriteString("  - 无类型符文 → rune (int32)\n")
	sb.WriteString("  - 无类型整数 → int\n")
	sb.WriteString("  - 无类型浮点 → float64\n")
	sb.WriteString("  - 无类型复数 → complex128\n")
	sb.WriteString("  - 无类型字符串 → string\n\n")

	// 示例 1: 无类型常量的灵活性
	sb.WriteString("【示例 1: 无类型常量的灵活性】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 无类型整数常量，可以赋值给多种整数类型\n")
	sb.WriteString("    const untyped = 42\n")
	sb.WriteString("    \n")
	sb.WriteString("    var i8 int8 = untyped\n")
	sb.WriteString("    var i16 int16 = untyped\n")
	sb.WriteString("    var i32 int32 = untyped\n")
	sb.WriteString("    var i64 int64 = untyped\n")
	sb.WriteString("    var f32 float32 = untyped\n")
	sb.WriteString("    var f64 float64 = untyped\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(i8, i16, i32, i64, f32, f64)\n")
	sb.WriteString("    // 输出: 42 42 42 42 42 42\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 无类型常量可以隐式转换为兼容类型，非常灵活。\n\n")

	// 示例 2: 类型化常量的限制
	sb.WriteString("【示例 2: 类型化常量的限制】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 类型化为 int8 的常量\n")
	sb.WriteString("    const typed int8 = 42\n")
	sb.WriteString("    \n")
	sb.WriteString("    var i8 int8 = typed    // OK: 类型匹配\n")
	sb.WriteString("    // var i16 int16 = typed // 编译错误: 类型不匹配\n")
	sb.WriteString("    // var f64 float64 = typed // 编译错误: 类型不匹配\n")
	sb.WriteString("    \n")
	sb.WriteString("    _ = i8\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 类型化常量只能赋值给相同类型的变量，不能隐式转换。\n\n")

	// 示例 3: 默认类型
	sb.WriteString("【示例 3: 默认类型】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        boolConst = true        // 默认类型: bool\n")
	sb.WriteString("        runeConst = 'A'          // 默认类型: rune (int32)\n")
	sb.WriteString("        intConst = 42           // 默认类型: int\n")
	sb.WriteString("        floatConst = 3.14       // 默认类型: float64\n")
	sb.WriteString("        complexConst = 1 + 2i   // 默认类型: complex128\n")
	sb.WriteString("        stringConst = \"hello\"   // 默认类型: string\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 当需要类型时，使用默认类型\n")
	sb.WriteString("    var b bool = boolConst\n")
	sb.WriteString("    var r rune = runeConst\n")
	sb.WriteString("    var i int = intConst\n")
	sb.WriteString("    var f float64 = floatConst\n")
	sb.WriteString("    var c complex128 = complexConst\n")
	sb.WriteString("    var s string = stringConst\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"%T %T %T %T %T %T\\n\", b, r, i, f, c, s)\n")
	sb.WriteString("    // 输出: bool int32 int float64 complex128 string\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 无类型常量有默认类型，在需要类型时会使用默认类型。\n\n")

	// 示例 4: 精度保持
	sb.WriteString("【示例 4: 精度保持】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 无类型浮点常量保持高精度\n")
	sb.WriteString("    const precise = 1.0 / 3.0\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 赋值给 float32 时精度降低\n")
	sb.WriteString("    var f32 float32 = precise\n")
	sb.WriteString("    // 赋值给 float64 时精度更高\n")
	sb.WriteString("    var f64 float64 = precise\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"float32: %.20f\\n\", f32)\n")
	sb.WriteString("    fmt.Printf(\"float64: %.20f\\n\", f64)\n")
	sb.WriteString("    // 输出:\n")
	sb.WriteString("    // float32: 0.33333334326744079590\n")
	sb.WriteString("    // float64: 0.33333333333333331483\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 无类型常量本身精度不受限制\n")
	sb.WriteString("    const huge = 1 << 100  // 远超 int64 范围\n")
	sb.WriteString("    const result = huge >> 100  // 结果: 1\n")
	sb.WriteString("    fmt.Println(result)  // 输出: 1\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 无类型常量保持任意精度，直到赋值给具体类型时才受类型限制。\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 混淆类型化和无类型化常量:\n")
	sb.WriteString("   const a = 42        // 无类型化\n")
	sb.WriteString("   const b int = 42    // 类型化\n")
	sb.WriteString("   // a 可以赋值给多种类型，b 只能赋值给 int\n")
	sb.WriteString("\n")
	sb.WriteString("2. 类型化常量不能混合运算:\n")
	sb.WriteString("   const a int = 10\n")
	sb.WriteString("   const b float64 = 3.14\n")
	sb.WriteString("   // const c = a * b  // 编译错误: 类型不匹配\n")
	sb.WriteString("\n")
	sb.WriteString("3. 无类型常量在需要类型时使用默认类型:\n")
	sb.WriteString("   const x = 42\n")
	sb.WriteString("   // 如果上下文需要 int，x 的类型就是 int\n")
	sb.WriteString("   // 如果上下文需要 float64，x 会隐式转换为 float64\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayTypedUntyped 展示并解释 Go 语言中类型化和无类型化常量的区别。
// 理解这两种常量的区别对于掌握 Go 的类型系统非常重要。
func DisplayTypedUntyped() {
	fmt.Print(GetTypedUntypedContent())
}
