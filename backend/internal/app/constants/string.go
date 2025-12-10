// Package constants - 字符串常量学习模块
//
// 本文件介绍 Go 语言中的字符串常量(String Constants)。
// 字符串常量是无类型的字符序列，可以使用双引号或反引号表示。
package constants

import (
	"fmt"
	"strings"
)

// GetStringContent 返回字符串常量相关的学习内容
func GetStringContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== String Constants (字符串常量) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("字符串常量 (String Literals) 是通过双引号或反引号创建的字符序列。\n")
	sb.WriteString("Go 语言的字符串是 UTF-8 编码的不可变字节序列。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("1. 解释型字符串 (Interpreted string literals): 使用双引号 \"...\"，\n")
	sb.WriteString("   支持转义字符 (如 \\n, \\t, \\\", \\\\)。\n")
	sb.WriteString("2. 原始字符串 (Raw string literals): 使用反引号 `...`，\n")
	sb.WriteString("   不支持转义，内容原样输出（包括换行符），常用于多行文本或正则表达式。\n\n")

	// 示例 1: 解释型 vs 原始型
	sb.WriteString("【示例 1: 解释型 vs 原始型】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const s1 = \"Hello\\nWorld\"  // 包含换行符转义\n")
	sb.WriteString("    const s2 = `Hello\\nWorld`  // 原样包含反斜杠和n字符\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(\"=== s1 ===\")\n")
	sb.WriteString("    fmt.Println(s1)\n")
	sb.WriteString("    fmt.Println(\"=== s2 ===\")\n")
	sb.WriteString("    fmt.Println(s2)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 2: 多行字符串
	sb.WriteString("【示例 2: 多行字符串】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const menu = `\n")
	sb.WriteString("Select an option:\n")
	sb.WriteString("1. Start\n")
	sb.WriteString("2. Stop\n")
	sb.WriteString("3. Exit\n")
	sb.WriteString("`\n")
	sb.WriteString("    fmt.Println(menu)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 3: 字符串连接 (常量表达式)
	sb.WriteString("【示例 3: 字符串连接】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const greeting = \"Hello\"\n")
	sb.WriteString("    const name = \"Gopher\"\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 字符串连接在编译时完成\n")
	sb.WriteString("    const message = greeting + \", \" + name + \"!\"\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(message) // 输出: Hello, Gopher!\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 示例 4: 内置函数 len()
	sb.WriteString("【示例 4: len() 函数】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const str = \"Hello, 世界\"\n")
	sb.WriteString("    \n")
	sb.WriteString("    // len() 对字符串常量返回字节长度（编译时常量）\n")
	sb.WriteString("    const length = len(str)\n")
	sb.WriteString("    \n")
	sb.WriteString("    // Hello (5) + , (1) + Space (1) + 世界 (3+3) = 13\n")
	sb.WriteString("    fmt.Printf(\"Length: %d\\n\", length)\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 字符串是不可变的:\n")
	sb.WriteString("   const s = \"hello\"\n")
	sb.WriteString("   // s[0] = 'H' // 编译错误: cannot assign to s[0]\n")
	sb.WriteString("\n")
	sb.WriteString("2. 单引号 vs 双引号:\n")
	sb.WriteString("   'a' 是符文 (rune/int32)\n")
	sb.WriteString("   \"a\" 是字符串 (string)\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayString 展示并解释 Go 语言中的字符串常量。
func DisplayString() {
	fmt.Print(GetStringContent())
}
