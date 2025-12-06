// Package constants - 常量表达式学习模块
//
// 本文件介绍 Go 语言中的常量表达式(Constant Expressions)。
// 常量表达式是在编译时求值的表达式，可以包含算术、比较、逻辑等运算。
package constants

import (
	"fmt"
	"strings"
)

// GetExpressionsContent 返回常量表达式相关的学习内容
func GetExpressionsContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Constant Expressions (常量表达式) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("常量表达式是在编译时求值的表达式，由常量、运算符和函数调用组成。\n")
	sb.WriteString("常量表达式的结果也是常量，具有任意精度，不会溢出。\n")
	sb.WriteString("常量表达式可以包含算术、比较、逻辑、位运算等操作。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("常量表达式可以包含:\n")
	sb.WriteString("1. 字面量常量: 123, 3.14, true, 'a', \"hello\"\n")
	sb.WriteString("2. 已声明的常量标识符\n")
	sb.WriteString("3. 运算符: +, -, *, /, %, ==, !=, <, <=, >, >=, &&, ||, !, &, |, ^, <<, >>\n")
	sb.WriteString("4. 类型转换: T(x) (结果必须是常量)\n")
	sb.WriteString("5. 内置函数: min, max, len, real, imag, complex, unsafe.Sizeof\n\n")

	// 示例 1: 算术表达式
	sb.WriteString("【示例 1: 算术表达式】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        a = 10\n")
	sb.WriteString("        b = 20\n")
	sb.WriteString("        c = 30\n")
	sb.WriteString("        sum = a + b        // 30\n")
	sb.WriteString("        diff = b - a        // 10\n")
	sb.WriteString("        prod = a * b        // 200\n")
	sb.WriteString("        quot = b / a        // 2\n")
	sb.WriteString("        rem = b % a         // 0\n")
	sb.WriteString("        complexExpr = (a + b) * c / 2  // (10+20)*30/2 = 450\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(sum, diff, prod, quot, rem, complexExpr)\n")
	sb.WriteString("    // 输出: 30 10 200 2 0 450\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 常量表达式支持所有算术运算符，结果在编译时计算。\n\n")

	// 示例 2: 比较表达式
	sb.WriteString("【示例 2: 比较表达式】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        x = 10\n")
	sb.WriteString("        y = 20\n")
	sb.WriteString("        z = 10\n")
	sb.WriteString("        isEqual = x == y      // false\n")
	sb.WriteString("        isNotEqual = x != y   // true\n")
	sb.WriteString("        isLess = x < y        // true\n")
	sb.WriteString("        isLessOrEqual = x <= z // true\n")
	sb.WriteString("        isGreater = y > x      // true\n")
	sb.WriteString("        isGreaterOrEqual = x >= z // true\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(isEqual, isNotEqual, isLess, isLessOrEqual, isGreater, isGreaterOrEqual)\n")
	sb.WriteString("    // 输出: false true true true true true\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 比较运算符的结果是布尔常量，在编译时求值。\n\n")

	// 示例 3: 逻辑表达式
	sb.WriteString("【示例 3: 逻辑表达式】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        a = true\n")
	sb.WriteString("        b = false\n")
	sb.WriteString("        c = true\n")
	sb.WriteString("        and1 = a && b        // false\n")
	sb.WriteString("        and2 = a && c        // true\n")
	sb.WriteString("        or1 = a || b         // true\n")
	sb.WriteString("        or2 = b || b         // false\n")
	sb.WriteString("        not1 = !a            // false\n")
	sb.WriteString("        not2 = !b            // true\n")
	sb.WriteString("        complex = a && b || c  // (true && false) || true = true\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(and1, and2, or1, or2, not1, not2, complex)\n")
	sb.WriteString("    // 输出: false true true false false true true\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 逻辑运算符 &&, ||, ! 可以用于布尔常量，结果也是布尔常量。\n\n")

	// 示例 4: 混合类型表达式
	sb.WriteString("【示例 4: 混合类型表达式】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        intConst = 10\n")
	sb.WriteString("        floatConst = 3.14\n")
	sb.WriteString("        // 无类型常量可以混合运算\n")
	sb.WriteString("        mixed1 = intConst * floatConst  // 31.4 (无类型浮点常量)\n")
	sb.WriteString("        mixed2 = intConst + floatConst  // 13.14\n")
	sb.WriteString("        mixed3 = floatConst / intConst  // 0.314\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(mixed1, mixed2, mixed3)\n")
	sb.WriteString("    // 输出: 31.4 13.14 0.314\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 无类型常量可以混合运算，结果类型由运算决定（整数+浮点=浮点）。\n\n")

	// 示例 5: 嵌套表达式
	sb.WriteString("【示例 5: 嵌套表达式】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        a = 2\n")
	sb.WriteString("        b = 3\n")
	sb.WriteString("        c = 4\n")
	sb.WriteString("        // 嵌套表达式，遵循运算符优先级\n")
	sb.WriteString("        result1 = a + b*c - (a+b)*c  // 2 + 3*4 - (2+3)*4 = 2+12-20 = -6\n")
	sb.WriteString("        result2 = (a + b) * (c - a) // (2+3)*(4-2) = 5*2 = 10\n")
	sb.WriteString("        result3 = a << (b + c)      // 2 << (3+4) = 2 << 7 = 256\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(result1, result2, result3)\n")
	sb.WriteString("    // 输出: -6 10 256\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 常量表达式可以嵌套，遵循标准的运算符优先级规则。\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 除零错误:\n")
	sb.WriteString("   const x = 10 / 0  // 编译错误: division by zero\n")
	sb.WriteString("\n")
	sb.WriteString("2. 类型化常量不能混合运算:\n")
	sb.WriteString("   const a int = 10\n")
	sb.WriteString("   const b float64 = 3.14\n")
	sb.WriteString("   const c = a * b  // 编译错误: 类型不匹配\n")
	sb.WriteString("\n")
	sb.WriteString("3. 常量表达式不能包含变量:\n")
	sb.WriteString("   var x = 10\n")
	sb.WriteString("   const y = x + 5  // 编译错误: x 是变量，不是常量\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayExpressions 展示并解释 Go 语言中的常量表达式。
// 常量表达式在编译时求值，支持算术、比较、逻辑等运算。
func DisplayExpressions() {
	fmt.Print(GetExpressionsContent())
}
