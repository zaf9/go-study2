package lexical_elements

import (
	"fmt"
	"unicode"
)

// DisplayRunes 展示 Go 语言中字符字面量（rune）的使用
// rune 是 int32 的别名，表示一个 Unicode 码点
func DisplayRunes() {
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("【词法元素 - 字符字面量 (Rune Literals)】")
	fmt.Println(repeatString("=", 60))

	// 1. rune 的定义
	fmt.Println("\n1. rune 的定义")
	fmt.Println("   rune 是 int32 的别名，用于表示一个 Unicode 码点")
	fmt.Println("   字符字面量用单引号 ' ' 括起来")

	r1 := 'A' // ASCII 字符
	r2 := '中' // 中文字符
	r3 := '😀' // Emoji 表情

	fmt.Printf("   'A' = %c (Unicode: U+%04X, 值: %d)\n", r1, r1, r1)
	fmt.Printf("   '中' = %c (Unicode: U+%04X, 值: %d)\n", r2, r2, r2)
	fmt.Printf("   '😀' = %c (Unicode: U+%04X, 值: %d)\n", r3, r3, r3)

	// 2. 转义序列
	fmt.Println("\n2. 转义序列（Escape Sequences）")
	fmt.Println("   使用反斜杠 \\ 表示特殊字符")

	fmt.Println("\n   常用转义序列:")
	fmt.Printf("   '\\n' = 换行符 (值: %d)\n", '\n')
	fmt.Printf("   '\\t' = 制表符 (值: %d)\n", '\t')
	fmt.Printf("   '\\r' = 回车符 (值: %d)\n", '\r')
	fmt.Printf("   '\\\\' = 反斜杠 (值: %d)\n", '\\')
	fmt.Printf("   '\\'' = 单引号 (值: %d)\n", '\'')
	fmt.Printf("   '\\\"' = 双引号 (值: %d)\n", '"')

	// 演示转义序列的效果
	fmt.Println("\n   转义序列效果演示:")
	fmt.Println("   第一行\\n第二行")    // 换行
	fmt.Println("   列1\\t列2\\t列3") // 制表符

	// 3. Unicode 转义
	fmt.Println("\n3. Unicode 转义表示法")
	fmt.Println("   \\xNN     - 2 位十六进制（8 位字符）")
	fmt.Println("   \\uNNNN   - 4 位十六进制（16 位 Unicode）")
	fmt.Println("   \\UNNNNNNNN - 8 位十六进制（32 位 Unicode）")

	hex8 := '\x41'            // 十六进制 41 = 'A'
	unicode16 := '\u4E2D'     // Unicode 4E2D = '中'
	unicode32 := '\U0001F600' // Unicode 1F600 = '😀'

	fmt.Printf("   '\\x41' = %c (ASCII 'A')\n", hex8)
	fmt.Printf("   '\\u4E2D' = %c (中文 '中')\n", unicode16)
	fmt.Printf("   '\\U0001F600' = %c (Emoji '😀')\n", unicode32)

	// 4. 八进制转义
	fmt.Println("\n4. 八进制转义表示法")
	fmt.Println("   \\NNN - 3 位八进制数字（0-377）")

	octal := '\101' // 八进制 101 = 十进制 65 = 'A'
	fmt.Printf("   '\\101' = %c (八进制 101 = 'A')\n", octal)

	// 5. rune 类型和 byte 类型的区别
	fmt.Println("\n5. rune 类型和 byte 类型的区别")
	fmt.Println("   byte  - uint8 的别名，表示一个字节（ASCII 字符）")
	fmt.Println("   rune  - int32 的别名，表示一个 Unicode 码点")

	var b byte = 'A' // byte 只能表示 ASCII 字符
	var r rune = '中' // rune 可以表示任何 Unicode 字符

	fmt.Printf("   byte: %c (大小: %d 字节)\n", b, 1)
	fmt.Printf("   rune: %c (大小: %d 字节)\n", r, 4)

	// 6. 字符串和 rune 的关系
	fmt.Println("\n6. 字符串和 rune 的关系")
	fmt.Println("   字符串是 rune 的序列，使用 UTF-8 编码")

	str := "Hello,世界"
	fmt.Printf("   字符串: %s\n", str)
	fmt.Printf("   字节长度: %d\n", len(str))

	// 遍历字符串的 rune
	fmt.Println("\n   遍历字符串的每个 rune:")
	for index, runeValue := range str {
		fmt.Printf("   索引 %d: %c (Unicode: U+%04X)\n", index, runeValue, runeValue)
	}

	// 7. rune 数组和切片
	fmt.Println("\n7. rune 数组和切片")

	runes := []rune{'G', 'o', '语', '言'}
	fmt.Printf("   rune 切片: %c\n", runes)
	fmt.Printf("   转换为字符串: %s\n", string(runes))

	// 8. 字符串和 []rune 的转换
	fmt.Println("\n8. 字符串和 []rune 的转换")

	s := "Go编程"
	runeSlice := []rune(s) // 字符串转 rune 切片

	fmt.Printf("   原字符串: %s\n", s)
	fmt.Printf("   rune 切片: %v\n", runeSlice)
	fmt.Printf("   rune 数量: %d\n", len(runeSlice))
	fmt.Printf("   字节数量: %d\n", len(s))

	// 9. 判断字符类型
	fmt.Println("\n9. 判断字符类型（使用 unicode 包）")

	testRunes := []rune{'A', '9', '中', ' ', '!'}
	for _, r := range testRunes {
		fmt.Printf("   '%c': ", r)
		if unicode.IsLetter(r) {
			fmt.Print("字母 ")
		}
		if unicode.IsDigit(r) {
			fmt.Print("数字 ")
		}
		if unicode.IsSpace(r) {
			fmt.Print("空格 ")
		}
		if unicode.IsPunct(r) {
			fmt.Print("标点 ")
		}
		fmt.Println()
	}

	// 10. 最佳实践
	fmt.Println("\n10. 最佳实践")
	fmt.Println("   ✓ 使用 rune 处理 Unicode 字符")
	fmt.Println("   ✓ 使用 byte 处理 ASCII 字符或二进制数据")
	fmt.Println("   ✓ 遍历字符串时使用 range（自动处理 UTF-8）")
	fmt.Println("   ✓ 字符字面量用单引号 ' '，字符串用双引号 \" \"")
	fmt.Println("   ✓ 了解字符串的字节长度和字符数量可能不同")
	fmt.Println("   ✓ 使用 unicode 包进行字符分类和判断")

	fmt.Println("\n" + repeatString("=", 60))
}
