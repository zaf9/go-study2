// Package constants - 内置函数学习模块
//
// 本文件介绍 Go 语言中可以用于常量的内置函数(Built-in Functions)。
// 这些函数在编译时求值，结果也是常量。
package constants

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-study2/internal/infrastructure/logger"
)

// GetBuiltinFunctionsContent 返回内置函数相关的学习内容
func GetBuiltinFunctionsContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Built-in Functions (内置函数) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("Go 语言提供了一些内置函数，其中部分可以在常量表达式中使用。\n")
	sb.WriteString("这些函数在编译时求值，结果也是常量。\n")
	sb.WriteString("可用于常量的内置函数包括:\n")
	sb.WriteString("1. min/max: 返回多个常量中的最小/最大值 (Go 1.21+)\n")
	sb.WriteString("2. len: 用于字符串常量和数组，返回长度\n")
	sb.WriteString("3. real/imag: 提取复数常量的实部/虚部\n")
	sb.WriteString("4. complex: 从实部和虚部构造复数常量\n")
	sb.WriteString("5. unsafe.Sizeof: 返回类型的大小(字节)\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("内置函数调用语法:\n")
	sb.WriteString("- min(x, y, ...): 返回最小值，所有参数必须是常量\n")
	sb.WriteString("- max(x, y, ...): 返回最大值，所有参数必须是常量\n")
	sb.WriteString("- len(s): s 是字符串常量或数组，返回长度\n")
	sb.WriteString("- real(c): c 是复数常量，返回实部\n")
	sb.WriteString("- imag(c): c 是复数常量，返回虚部\n")
	sb.WriteString("- complex(r, i): r 和 i 是常量，返回复数 r+i\n")
	sb.WriteString("- unsafe.Sizeof(x): x 是类型，返回类型大小\n\n")

	// 示例 1: min 和 max
	sb.WriteString("【示例 1: min 和 max】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        a = 10\n")
	sb.WriteString("        b = 20\n")
	sb.WriteString("        c = 15\n")
	sb.WriteString("        minimum = min(a, b, c)  // 10\n")
	sb.WriteString("        maximum = max(a, b, c)  // 20\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(minimum, maximum)  // 输出: 10 20\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 也可以用于浮点数\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        f1 = 3.14\n")
	sb.WriteString("        f2 = 2.71\n")
	sb.WriteString("        fMin = min(f1, f2)  // 2.71\n")
	sb.WriteString("    )\n")
	sb.WriteString("    fmt.Println(fMin)  // 输出: 2.71\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: min 和 max 函数可以用于整数和浮点常量，返回最小/最大值。\n\n")

	// 示例 2: len (字符串)
	sb.WriteString("【示例 2: len (字符串)】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        str1 = \"Hello\"\n")
	sb.WriteString("        str2 = \"世界\"\n")
	sb.WriteString("        str3 = \"Hello, 世界\"\n")
	sb.WriteString("        len1 = len(str1)  // 5 (ASCII 字符)\n")
	sb.WriteString("        len2 = len(str2)  // 6 (UTF-8 编码，每个中文字符 3 字节)\n")
	sb.WriteString("        len3 = len(str3)  // 13 (5 + 1 + 6 + 1)\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(len1, len2, len3)  // 输出: 5 6 13\n")
	sb.WriteString("    \n")
	sb.WriteString("    // len 返回字节数，不是字符数\n")
	sb.WriteString("    fmt.Printf(\"str2 字节数: %d\\n\", len2)\n")
	sb.WriteString("    // 输出: str2 字节数: 6\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: len 函数用于字符串常量时，返回字符串的字节数(UTF-8 编码)。\n\n")

	// 示例 3: len (数组)
	sb.WriteString("【示例 3: len (数组)】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 数组长度是常量\n")
	sb.WriteString("    const arrayLen = len([5]int{})  // 5\n")
	sb.WriteString("    const sliceLen = len([10]string{})  // 10\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(arrayLen, sliceLen)  // 输出: 5 10\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 注意: len 不能用于切片常量(切片不是常量)\n")
	sb.WriteString("    // const s = []int{1, 2, 3}  // 错误: 切片不是常量\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: len 函数可以用于数组常量，返回数组长度。\n\n")

	// 示例 4: real 和 imag
	sb.WriteString("【示例 4: real 和 imag】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        c = 3 + 4i\n")
	sb.WriteString("        r = real(c)  // 3\n")
	sb.WriteString("        i = imag(c)  // 4\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(r, i)  // 输出: 3 4\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 可以用于计算\n")
	sb.WriteString("    const magnitude = r*r + i*i  // 3*3 + 4*4 = 25\n")
	sb.WriteString("    fmt.Println(magnitude)  // 输出: 25\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: real 和 imag 函数提取复数常量的实部和虚部，结果是常量。\n\n")

	// 示例 5: complex
	sb.WriteString("【示例 5: complex】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        realPart = 5.0\n")
	sb.WriteString("        imagPart = 12.0\n")
	sb.WriteString("        c = complex(realPart, imagPart)  // 5+12i\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(c)  // 输出: (5+12i)\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 可以验证\n")
	sb.WriteString("    fmt.Println(real(c), imag(c))  // 输出: 5 12\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: complex 函数从实部和虚部构造复数常量，参数必须是常量。\n\n")

	// 示例 6: unsafe.Sizeof
	sb.WriteString("【示例 6: unsafe.Sizeof】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("    \"fmt\"\n")
	sb.WriteString("    \"unsafe\"\n")
	sb.WriteString(")\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 类型大小是常量\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        intSize = unsafe.Sizeof(int(0))\n")
	sb.WriteString("        int64Size = unsafe.Sizeof(int64(0))\n")
	sb.WriteString("        float64Size = unsafe.Sizeof(float64(0))\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(intSize, int64Size, float64Size)\n")
	sb.WriteString("    // 输出取决于平台，如 64 位系统: 8 8 8\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 注意: unsafe.Sizeof 返回的是类型大小，不是值大小\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: unsafe.Sizeof 返回类型的大小(字节数)，结果是常量。\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 内置函数参数必须是常量:\n")
	sb.WriteString("   var x = 10\n")
	sb.WriteString("   // const y = min(x, 20)  // 编译错误: x 是变量\n")
	sb.WriteString("\n")
	sb.WriteString("2. len 不能用于切片:\n")
	sb.WriteString("   // const s = []int{1, 2, 3}  // 错误: 切片不是常量\n")
	sb.WriteString("   // const l = len(s)  // 错误: 切片不是常量\n")
	sb.WriteString("\n")
	sb.WriteString("3. min/max 参数类型必须一致:\n")
	sb.WriteString("   const a = 10\n")
	sb.WriteString("   const b = 3.14\n")
	sb.WriteString("   // const c = min(a, b)  // 编译错误: 类型不匹配\n")
	sb.WriteString("\n")
	sb.WriteString("4. unsafe.Sizeof 需要类型，不是值:\n")
	sb.WriteString("   const x = 42\n")
	sb.WriteString("   // const s = unsafe.Sizeof(x)  // 错误: 需要类型\n")
	sb.WriteString("   const s = unsafe.Sizeof(int(0))  // 正确\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayBuiltinFunctions 展示并解释 Go 语言中可以用于常量的内置函数。
// 这些函数在编译时求值，结果也是常量。
func DisplayBuiltinFunctions() {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logger.LogBiz(context.Background(), "DisplayBuiltinFunctions", map[string]interface{}{
			"operation": "display_builtin_functions",
			"result":    "success",
		}, nil, duration)
	}()

	fmt.Print(GetBuiltinFunctionsContent())
}
