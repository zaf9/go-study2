package lexical_elements

import (
	"fmt"
	"strings"
	"unicode"
)

// GetRunesContent 返回字符字面量相关的学习内容
func GetRunesContent() string {
	var sb strings.Builder

	sb.WriteString("\n" + repeatString("=", 60) + "\n")
	sb.WriteString("【词法元素 - 字符字面量 (Rune Literals)】\n")
	sb.WriteString(repeatString("=", 60) + "\n")

	// 1. rune 的定义
	sb.WriteString("\n1. rune 的定义\n")
	sb.WriteString("   rune 是 int32 的别名，用于表示一个 Unicode 码点\n")
	sb.WriteString("   字符字面量用单引号 ' ' 括起来\n")

	r1 := 'A' // ASCII 字符
	r2 := '中' // 中文字符
	r3 := '😀' // Emoji 表情

	sb.WriteString(fmt.Sprintf("   'A' = %c (Unicode: U+%04X, 值: %d)\n", r1, r1, r1))
	sb.WriteString(fmt.Sprintf("   '中' = %c (Unicode: U+%04X, 值: %d)\n", r2, r2, r2))
	sb.WriteString(fmt.Sprintf("   '😀' = %c (Unicode: U+%04X, 值: %d)\n", r3, r3, r3))

	// 2. 转义序列
	sb.WriteString("\n2. 转义序列（Escape Sequences）\n")
	sb.WriteString("   使用反斜杠 \\ 表示特殊字符\n")

	sb.WriteString("\n   常用转义序列:\n")
	sb.WriteString(fmt.Sprintf("   '\\n' = 换行符 (值: %d)\n", '\n'))
	sb.WriteString(fmt.Sprintf("   '\\t' = 制表符 (值: %d)\n", '\t'))
	sb.WriteString(fmt.Sprintf("   '\\r' = 回车符 (值: %d)\n", '\r'))
	sb.WriteString(fmt.Sprintf("   '\\\\' = 反斜杠 (值: %d)\n", '\\'))
	sb.WriteString(fmt.Sprintf("   '\\'' = 单引号 (值: %d)\n", '\''))
	sb.WriteString(fmt.Sprintf("   '\\\"' = 双引号 (值: %d)\n", '"'))

	// 演示转义序列的效果
	sb.WriteString("\n   转义序列效果演示:\n")
	sb.WriteString("   第一行\\n第二行\n")    // 换行
	sb.WriteString("   列1\\t列2\\t列3\n") // 制表符

	// 3. Unicode 转义
	sb.WriteString("\n3. Unicode 转义表示法\n")
	sb.WriteString("   \\xNN     - 2 位十六进制（8 位字符）\n")
	sb.WriteString("   \\uNNNN   - 4 位十六进制（16 位 Unicode）\n")
	sb.WriteString("   \\UNNNNNNNN - 8 位十六进制（32 位 Unicode）\n")

	hex8 := '\x41'            // 十六进制 41 = 'A'
	unicode16 := '\u4E2D'     // Unicode 4E2D = '中'
	unicode32 := '\U0001F600' // Unicode 1F600 = '😀'

	sb.WriteString(fmt.Sprintf("   '\\x41' = %c (ASCII 'A')\n", hex8))
	sb.WriteString(fmt.Sprintf("   '\\u4E2D' = %c (中文 '中')\n", unicode16))
	sb.WriteString(fmt.Sprintf("   '\\U0001F600' = %c (Emoji '😀')\n", unicode32))

	// 4. 八进制转义
	sb.WriteString("\n4. 八进制转义表示法\n")
	sb.WriteString("   \\NNN - 3 位八进制数字（0-377）\n")

	octal := '\101' // 八进制 101 = 十进制 65 = 'A'
	sb.WriteString(fmt.Sprintf("   '\\101' = %c (八进制 101 = 'A')\n", octal))

	// 5. rune 类型和 byte 类型的区别
	sb.WriteString("\n5. rune 类型和 byte 类型的区别\n")
	sb.WriteString("   byte  - uint8 的别名，表示一个字节（ASCII 字符）\n")
	sb.WriteString("   rune  - int32 的别名，表示一个 Unicode 码点\n")

	var b byte = 'A' // byte 只能表示 ASCII 字符
	var r rune = '中' // rune 可以表示任何 Unicode 字符

	sb.WriteString(fmt.Sprintf("   byte: %c (大小: %d 字节)\n", b, 1))
	sb.WriteString(fmt.Sprintf("   rune: %c (大小: %d 字节)\n", r, 4))

	// 6. 字符串和 rune 的关系
	sb.WriteString("\n6. 字符串和 rune 的关系\n")
	sb.WriteString("   字符串是 rune 的序列，使用 UTF-8 编码\n")

	str := "Hello,世界"
	sb.WriteString(fmt.Sprintf("   字符串: %s\n", str))
	sb.WriteString(fmt.Sprintf("   字节长度: %d\n", len(str)))

	// 遍历字符串的 rune
	sb.WriteString("\n   遍历字符串的每个 rune:\n")
	for index, runeValue := range str {
		sb.WriteString(fmt.Sprintf("   索引 %d: %c (Unicode: U+%04X)\n", index, runeValue, runeValue))
	}

	// 7. rune 数组和切片
	sb.WriteString("\n7. rune 数组和切片\n")

	runes := []rune{'G', 'o', '语', '言'}
	sb.WriteString(fmt.Sprintf("   rune 切片: %c\n", runes))
	sb.WriteString(fmt.Sprintf("   转换为字符串: %s\n", string(runes)))

	// 8. 字符串和 []rune 的转换
	sb.WriteString("\n8. 字符串和 []rune 的转换\n")

	s := "Go编程"
	runeSlice := []rune(s) // 字符串转 rune 切片

	sb.WriteString(fmt.Sprintf("   原字符串: %s\n", s))
	sb.WriteString(fmt.Sprintf("   rune 切片: %v\n", runeSlice))
	sb.WriteString(fmt.Sprintf("   rune 数量: %d\n", len(runeSlice)))
	sb.WriteString(fmt.Sprintf("   字节数量: %d\n", len(s)))

	// 9. 判断字符类型
	sb.WriteString("\n9. 判断字符类型（使用 unicode 包）\n")

	testRunes := []rune{'A', '9', '中', ' ', '!'}
	for _, tr := range testRunes {
		sb.WriteString(fmt.Sprintf("   '%c': ", tr))
		if unicode.IsLetter(tr) {
			sb.WriteString("字母 ")
		}
		if unicode.IsDigit(tr) {
			sb.WriteString("数字 ")
		}
		if unicode.IsSpace(tr) {
			sb.WriteString("空格 ")
		}
		if unicode.IsPunct(tr) {
			sb.WriteString("标点 ")
		}
		sb.WriteString("\n")
	}

	// 10. 最佳实践
	sb.WriteString("\n10. 最佳实践\n")
	sb.WriteString("   ✓ 使用 rune 处理 Unicode 字符\n")
	sb.WriteString("   ✓ 使用 byte 处理 ASCII 字符或二进制数据\n")
	sb.WriteString("   ✓ 遍历字符串时使用 range（自动处理 UTF-8）\n")
	sb.WriteString("   ✓ 字符字面量用单引号 ' '，字符串用双引号 \" \"\n")
	sb.WriteString("   ✓ 了解字符串的字节长度和字符数量可能不同\n")
	sb.WriteString("   ✓ 使用 unicode 包进行字符分类和判断\n")

	sb.WriteString("\n" + repeatString("=", 60) + "\n")

	return sb.String()
}

// DisplayRunes 展示 Go 语言中字符字面量（rune）的使用
// rune 是 int32 的别名，表示一个 Unicode 码点
func DisplayRunes() {
	fmt.Print(GetRunesContent())
}
