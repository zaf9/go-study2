// Package constants - 常量类型转换学习模块
//
// 本文件介绍 Go 语言中常量类型转换(Conversions)的规则和限制。
// 常量转换需要满足可表示性(representability)要求。
package constants

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-study2/internal/infrastructure/logger"
)

// GetConversionsContent 返回常量类型转换相关的学习内容
func GetConversionsContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Conversions (类型转换) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("常量类型转换使用语法 T(x)，其中 T 是目标类型，x 是常量表达式。\n")
	sb.WriteString("转换必须满足可表示性(representability)要求:\n")
	sb.WriteString("1. 常量值必须能被目标类型表示\n")
	sb.WriteString("2. 整数常量 → 整数类型: 值必须在类型范围内\n")
	sb.WriteString("3. 浮点常量 → 浮点类型: 四舍五入到最接近的可表示值\n")
	sb.WriteString("4. 复数常量 → 复数类型: 实部和虚部分别转换\n")
	sb.WriteString("5. 值超出类型范围会导致编译错误\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("转换语法: T(x)\n")
	sb.WriteString("- T: 目标类型 (如 int8, float32, complex64)\n")
	sb.WriteString("- x: 常量表达式\n")
	sb.WriteString("\n")
	sb.WriteString("转换规则:\n")
	sb.WriteString("- 整数常量可以转换为任何整数类型(如果值在范围内)\n")
	sb.WriteString("- 浮点常量可以转换为任何浮点类型(精度可能损失)\n")
	sb.WriteString("- 整数常量可以转换为浮点类型\n")
	sb.WriteString("- 复数常量可以转换为复数类型\n")
	sb.WriteString("- 布尔、符文、字符串常量可以转换为对应类型\n\n")

	// 示例 1: 成功的整数转换
	sb.WriteString("【示例 1: 成功的整数转换】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const value = 100\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 可以转换为各种整数类型(值在范围内)\n")
	sb.WriteString("    const i8 = int8(value)\n")
	sb.WriteString("    const i16 = int16(value)\n")
	sb.WriteString("    const i32 = int32(value)\n")
	sb.WriteString("    const i64 = int64(value)\n")
	sb.WriteString("    const u8 = uint8(value)\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(i8, i16, i32, i64, u8)\n")
	sb.WriteString("    // 输出: 100 100 100 100 100\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 只要常量值在目标类型的范围内，转换就会成功。\n\n")

	// 示例 2: 浮点转换和精度损失
	sb.WriteString("【示例 2: 浮点转换和精度损失】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 高精度浮点常量\n")
	sb.WriteString("    const precise = 1.234567890123456789\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 转换为 float32 时精度降低\n")
	sb.WriteString("    const f32 = float32(precise)\n")
	sb.WriteString("    // 转换为 float64 时精度更高\n")
	sb.WriteString("    const f64 = float64(precise)\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"float32: %.20f\\n\", f32)\n")
	sb.WriteString("    fmt.Printf(\"float64: %.20f\\n\", f64)\n")
	sb.WriteString("    // 输出:\n")
	sb.WriteString("    // float32: 1.23456788063049316406\n")
	sb.WriteString("    // float64: 1.23456789012345678000\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 精度损失是允许的，但值必须在可表示范围内\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 浮点转换允许精度损失，会四舍五入到最接近的可表示值。\n\n")

	// 示例 3: 整数到浮点转换
	sb.WriteString("【示例 3: 整数到浮点转换】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const intValue = 42\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 整数常量可以转换为浮点类型\n")
	sb.WriteString("    const f32 = float32(intValue)\n")
	sb.WriteString("    const f64 = float64(intValue)\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"float32: %.1f, float64: %.1f\\n\", f32, f64)\n")
	sb.WriteString("    // 输出: float32: 42.0, float64: 42.0\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 大整数也可以转换\n")
	sb.WriteString("    const bigInt = 1000000\n")
	sb.WriteString("    const bigFloat = float64(bigInt)\n")
	sb.WriteString("    fmt.Println(bigFloat)  // 输出: 1e+06\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 整数常量可以转换为浮点类型，转换是精确的(整数可以精确表示)。\n\n")

	// 示例 4: 复数转换
	sb.WriteString("【示例 4: 复数转换】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // complex128 常量\n")
	sb.WriteString("    const c128 = 1 + 2i\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 可以转换为 complex64\n")
	sb.WriteString("    const c64 = complex64(c128)\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"complex128: %v\\n\", c128)\n")
	sb.WriteString("    fmt.Printf(\"complex64: %v\\n\", c64)\n")
	sb.WriteString("    // 输出:\n")
	sb.WriteString("    // complex128: (1+2i)\n")
	sb.WriteString("    // complex64: (1+2i)\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 实部和虚部分别转换\n")
	sb.WriteString("    const realPart = real(c128)\n")
	sb.WriteString("    const imagPart = imag(c128)\n")
	sb.WriteString("    fmt.Printf(\"实部: %.1f, 虚部: %.1f\\n\", realPart, imagPart)\n")
	sb.WriteString("    // 输出: 实部: 1.0, 虚部: 2.0\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 复数转换时，实部和虚部分别按照浮点转换规则处理。\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 值超出类型范围:\n")
	sb.WriteString("   const big = 1000\n")
	sb.WriteString("   // const small = int8(big)  // 编译错误: 1000 超出 int8 范围 (-128~127)\n")
	sb.WriteString("\n")
	sb.WriteString("2. 浮点溢出:\n")
	sb.WriteString("   const huge = 1e100\n")
	sb.WriteString("   // const f32 = float32(huge)  // 编译错误: 常量溢出 float32\n")
	sb.WriteString("\n")
	sb.WriteString("3. 类型不兼容:\n")
	sb.WriteString("   const str = \"hello\"\n")
	sb.WriteString("   // const i = int(str)  // 编译错误: 不能将字符串转换为整数\n")
	sb.WriteString("\n")
	sb.WriteString("4. 精度损失警告:\n")
	sb.WriteString("   // 浮点转换允许精度损失，但值必须在可表示范围内\n")
	sb.WriteString("   // 如果值超出范围，会导致编译错误\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayConversions 展示并解释 Go 语言中常量类型转换的规则。
// 常量转换需要满足可表示性要求，值必须在目标类型的范围内。
func DisplayConversions() {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logger.LogBiz(context.Background(), "DisplayConversions", map[string]interface{}{
			"operation": "display_constant_conversions",
			"result":    "success",
		}, nil, duration)
	}()

	fmt.Print(GetConversionsContent())
}
