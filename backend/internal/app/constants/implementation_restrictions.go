// Package constants - 实现限制学习模块
//
// 本文件介绍 Go 语言编译器对常量实现的限制和要求。
// 了解这些限制有助于理解常量的精度和范围。
package constants

import (
	"fmt"
	"strings"
)

// GetImplementationRestrictionsContent 返回实现限制相关的学习内容
func GetImplementationRestrictionsContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Implementation Restrictions (实现限制) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("Go 语言规范对编译器实现常量系统提出了最低要求:\n")
	sb.WriteString("1. 整数常量: 至少 256 位精度\n")
	sb.WriteString("2. 浮点常量: 尾数至少 256 位，指数至少 16 位(有符号二进制)\n")
	sb.WriteString("3. 精确表示: 无法精确表示整数常量时报错\n")
	sb.WriteString("4. 溢出处理: 浮点/复数溢出时报错\n")
	sb.WriteString("5. 精度限制: 浮点/复数精度不足时四舍五入到最接近值\n")
	sb.WriteString("这些限制确保了常量计算的精确性和可移植性。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("编译器要求:\n")
	sb.WriteString("- 整数常量必须至少支持 256 位精度\n")
	sb.WriteString("- 浮点常量尾数至少 256 位，指数至少 16 位\n")
	sb.WriteString("- 整数常量必须精确表示，不能有精度损失\n")
	sb.WriteString("- 浮点常量可以四舍五入，但值必须在可表示范围内\n")
	sb.WriteString("- 超出范围的值会导致编译错误\n\n")

	// 示例 1: 大整数常量
	sb.WriteString("【示例 1: 大整数常量】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 编译器必须支持至少 256 位整数常量\n")
	sb.WriteString("    const huge = 1 << 256  // 非常大的数\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 可以在常量表达式中使用\n")
	sb.WriteString("    const result = huge >> 256  // 结果: 1\n")
	sb.WriteString("    fmt.Println(result)  // 输出: 1\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 但不能直接赋值给 int 类型 (会编译错误)\n")
	sb.WriteString("    // var x int = huge  // 错误: constant overflows int\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 可以赋值给更大的类型(如果支持)\n")
	sb.WriteString("    // 但通常需要先进行运算缩小范围\n")
	sb.WriteString("    const small = huge >> 200  // 缩小范围\n")
	sb.WriteString("    fmt.Printf(\"small = %d\\n\", small)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 整数常量支持任意精度(至少 256 位)，不会溢出。\n\n")

	// 示例 2: 高精度浮点常量
	sb.WriteString("【示例 2: 高精度浮点常量】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 编译器必须支持至少 256 位尾数\n")
	sb.WriteString("    const pi = 3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 赋值给 float64 时精度降低\n")
	sb.WriteString("    var f64 float64 = pi\n")
	sb.WriteString("    fmt.Printf(\"float64: %.50f\\n\", f64)\n")
	sb.WriteString("    // 输出精度受 float64 限制，约 15-17 位有效数字\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 常量本身保持高精度\n")
	sb.WriteString("    const halfPi = pi / 2\n")
	sb.WriteString("    fmt.Printf(\"halfPi: %.50f\\n\", halfPi)\n")
	sb.WriteString("    // 常量运算保持高精度\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 浮点常量支持高精度(至少 256 位尾数)，直到赋值给具体类型。\n\n")

	// 示例 3: 溢出错误
	sb.WriteString("【示例 3: 溢出错误】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 整数常量溢出示例\n")
	sb.WriteString("    const maxInt64 = 1<<63 - 1  // int64 最大值\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 以下代码会导致编译错误\n")
	sb.WriteString("    // const overflow = int64(1 << 100)  // 编译错误: constant overflows int64\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 浮点常量溢出示例\n")
	sb.WriteString("    const hugeFloat = 1e100\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 以下代码会导致编译错误\n")
	sb.WriteString("    // const overflowFloat = float32(hugeFloat)  // 编译错误: constant overflows float32\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 正确的做法: 使用 float64\n")
	sb.WriteString("    const okFloat = float64(hugeFloat)\n")
	sb.WriteString("    _ = okFloat\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 当常量值超出目标类型的范围时，编译器会报错。\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 混淆常量精度和类型精度:\n")
	sb.WriteString("   const precise = 1.0 / 3.0  // 常量保持高精度\n")
	sb.WriteString("   var f64 float64 = precise  // 赋值时精度降低\n")
	sb.WriteString("   // 常量精度 ≠ 变量精度\n")
	sb.WriteString("\n")
	sb.WriteString("2. 尝试将超大常量赋值给小类型:\n")
	sb.WriteString("   const huge = 1 << 100\n")
	sb.WriteString("   // var x int = huge  // 编译错误: 溢出\n")
	sb.WriteString("\n")
	sb.WriteString("3. 不理解精度要求:\n")
	sb.WriteString("   // 编译器必须支持至少 256 位整数常量\n")
	sb.WriteString("   // 但赋值给具体类型时受类型限制\n")
	sb.WriteString("\n")
	sb.WriteString("4. 浮点常量精度损失:\n")
	sb.WriteString("   const pi = 3.141592653589793238462643383279\n")
	sb.WriteString("   var f32 float32 = pi  // 精度损失，但允许\n")
	sb.WriteString("   // 如果值超出 float32 范围，会报错\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayImplementationRestrictions 展示并解释 Go 语言编译器对常量实现的限制。
// 了解这些限制有助于理解常量的精度和范围。
func DisplayImplementationRestrictions() {
	fmt.Print(GetImplementationRestrictionsContent())
}
