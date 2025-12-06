// Package constants - 符文常量学习模块
//
// 本文件介绍 Go 语言中的符文常量(Rune Constants)。
// 符文常量表示 Unicode 代码点,是整数值的一种特殊形式。
package constants

import (
	"fmt"
	"strings"
)

// GetRuneContent 返回符文常量相关的学习内容
func GetRuneContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Rune Constants (符文常量) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("符文常量(Rune Literal)表示一个 Unicode 代码点。\n")
	sb.WriteString("在 Go 中,rune 是 int32 的别名,可以表示任何 Unicode 字符。\n")
	sb.WriteString("符文字面量用单引号括起来,如 'a', '中', '\\n' 等。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("符文字面量的形式:\n")
	sb.WriteString("  'x'     - 单个字符\n")
	sb.WriteString("  '\\t'   - 转义序列\n")
	sb.WriteString("  '\\x41' - 十六进制转义 (2 位)\n")
	sb.WriteString("  '\\u4e2d' - Unicode 转义 (4 位)\n")
	sb.WriteString("  '\\U0001F600' - Unicode 转义 (8 位)\n")
	sb.WriteString("  '\\101' - 八进制转义 (3 位)\n\n")

	// 示例 1: 基本符文常量
	sb.WriteString("【示例 1: 基本符文常量】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const letterA = 'A'\n")
	sb.WriteString("    const chinese = '中'\n")
	sb.WriteString("    const emoji = '😀'\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"'A' = %c, Unicode: U+%04X, 值: %d\\n\", letterA, letterA, letterA)\n")
	sb.WriteString("    fmt.Printf(\"'中' = %c, Unicode: U+%04X, 值: %d\\n\", chinese, chinese, chinese)\n")
	sb.WriteString("    fmt.Printf(\"'😀' = %c, Unicode: U+%04X, 值: %d\\n\", emoji, emoji, emoji)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("输出:\n")
	sb.WriteString("  'A' = A, Unicode: U+0041, 值: 65\n")
	sb.WriteString("  '中' = 中, Unicode: U+4E2D, 值: 20013\n")
	sb.WriteString("  '😀' = 😀, Unicode: U+1F600, 值: 128512\n\n")

	// 示例 2: Unicode 转义
	sb.WriteString("【示例 2: Unicode 转义序列】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    // 使用 \\u 表示 4 位 Unicode\n")
	sb.WriteString("    const zhong = '\\u4e2d'  // 中\n")
	sb.WriteString("    const guo = '\\u56fd'    // 国\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 使用 \\U 表示 8 位 Unicode (用于 emoji 等)\n")
	sb.WriteString("    const smile = '\\U0001F600'  // 😀\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"%c%c\\n\", zhong, guo)  // 输出: 中国\n")
	sb.WriteString("    fmt.Printf(\"%c\\n\", smile)         // 输出: 😀\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: \\u 后跟 4 位十六进制数,\\U 后跟 8 位十六进制数。\n\n")

	// 示例 3: 符文算术运算
	sb.WriteString("【示例 3: 符文常量的算术运算】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const a = 'A'\n")
	sb.WriteString("    const offset = 32  // 大写到小写的偏移量\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 符文常量可以参与算术运算\n")
	sb.WriteString("    const lowerA = a + offset\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Printf(\"'A' + 32 = '%c' (ASCII: %d)\\n\", lowerA, lowerA)\n")
	sb.WriteString("    // 输出: 'A' + 32 = 'a' (ASCII: 97)\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 计算字母表位置\n")
	sb.WriteString("    const letterC = 'C'\n")
	sb.WriteString("    const position = letterC - 'A' + 1\n")
	sb.WriteString("    fmt.Printf(\"'C' 是字母表第 %d 个字母\\n\", position)\n")
	sb.WriteString("    // 输出: 'C' 是字母表第 3 个字母\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 符文本质是整数,可以进行加减等算术运算。\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 使用双引号定义符文:\n")
	sb.WriteString("   const r = \"A\"  // 错误: 这是字符串,不是符文\n")
	sb.WriteString("   const r = 'A'  // 正确: 使用单引号\n")
	sb.WriteString("\n")
	sb.WriteString("2. 单引号内放多个字符:\n")
	sb.WriteString("   const r = 'AB'  // 编译错误: more than one character in rune literal\n")
	sb.WriteString("\n")
	sb.WriteString("3. 混淆符文和字节:\n")
	sb.WriteString("   const r = '中'  // 这是 1 个符文 (rune/int32)\n")
	sb.WriteString("   // 但在 UTF-8 编码中占 3 个字节\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayRune 展示并解释 Go 语言中的符文常量。
// 符文是表示 Unicode 代码点的整数值。
func DisplayRune() {
	fmt.Print(GetRuneContent())
}
