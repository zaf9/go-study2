package lexical_elements

import (
	"fmt"
	"strconv"
	"strings"
)

// DisplayStrings 展示 Go 语言中字符串字面量的各种形式和操作
// Go 字符串是 UTF-8 编码的字节序列，不可变
func DisplayStrings() {
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("【词法元素 - 字符串字面量 (String Literals)】")
	fmt.Println(repeatString("=", 60))

	// 1. 解释型字符串（双引号）
	fmt.Println("\n1. 解释型字符串（Interpreted Strings）")
	fmt.Println("   使用双引号 \" \" 括起来，支持转义序列")

	str1 := "Hello, World!"
	str2 := "你好，世界！"
	str3 := "第一行\n第二行" // 包含换行符

	fmt.Printf("   \"Hello, World!\" = %s\n", str1)
	fmt.Printf("   \"你好，世界！\" = %s\n", str2)
	fmt.Println("   \"第一行\\n第二行\" 的效果:")
	fmt.Println(str3)

	// 2. 原始字符串（反引号）
	fmt.Println("\n2. 原始字符串（Raw Strings）")
	fmt.Println("   使用反引号 ` ` 括起来，不解释转义序列")
	fmt.Println("   常用于正则表达式、文件路径、多行文本等")

	raw1 := `C:\Users\Documents\file.txt` // Windows 路径
	raw2 := `\n 不会被解释为换行符`
	raw3 := `这是第一行
这是第二行
这是第三行` // 多行字符串

	fmt.Printf("   `C:\\Users\\Documents\\file.txt` = %s\n", raw1)
	fmt.Printf("   `\\n 不会被解释为换行符` = %s\n", raw2)
	fmt.Println("\n   多行原始字符串:")
	fmt.Println(raw3)

	// 3. 转义序列
	fmt.Println("\n3. 转义序列（仅在解释型字符串中有效）")

	fmt.Println("   常用转义序列:")
	fmt.Println("   \\n  - 换行")
	fmt.Println("   \\t  - 制表符")
	fmt.Println("   \\r  - 回车")
	fmt.Println("   \\\\  - 反斜杠")
	fmt.Println("   \\\"  - 双引号")
	fmt.Println("   \\'  - 单引号")

	// 演示
	escaped := "姓名:\t张三\n年龄:\t25\n城市:\t\"北京\""
	fmt.Println("\n   转义序列示例:")
	fmt.Println(escaped)

	// 4. Unicode 转义
	fmt.Println("\n4. Unicode 转义（在解释型字符串中）")

	unicode1 := "\u4E2D\u6587"    // Unicode: 中文
	unicode2 := "\U0001F600"      // Unicode: 😀
	hex := "\x48\x65\x6C\x6C\x6F" // 十六进制: Hello

	fmt.Printf("   \"\\u4E2D\\u6587\" = %s\n", unicode1)
	fmt.Printf("   \"\\U0001F600\" = %s\n", unicode2)
	fmt.Printf("   \"\\x48\\x65\\x6C\\x6C\\x6F\" = %s\n", hex)

	// 5. 字符串的不可变性
	fmt.Println("\n5. 字符串的不可变性")
	fmt.Println("   Go 中的字符串是不可变的，不能修改字符串中的字符")

	s := "Hello"
	fmt.Printf("   原字符串: %s\n", s)
	// s[0] = 'h'  // 错误！不能修改字符串

	// 正确做法：创建新字符串
	newS := "h" + s[1:]
	fmt.Printf("   新字符串: %s (通过拼接创建)\n", newS)

	// 6. 字符串长度和索引
	fmt.Println("\n6. 字符串长度和索引")

	str := "Go语言"
	fmt.Printf("   字符串: %s\n", str)
	fmt.Printf("   字节长度 len(): %d (UTF-8 编码)\n", len(str))
	fmt.Printf("   字符数量: %d (rune 数量)\n", len([]rune(str)))

	// 字节索引
	fmt.Println("\n   按字节索引（可能不是完整字符）:")
	for i := 0; i < len(str); i++ {
		fmt.Printf("   str[%d] = %d (%c)\n", i, str[i], str[i])
	}

	// 7. 遍历字符串
	fmt.Println("\n7. 遍历字符串")

	text := "Go编程"

	// 方法 1: 按字节遍历
	fmt.Println("\n   方法 1: 按字节遍历（不推荐用于 Unicode）")
	for i := 0; i < len(text); i++ {
		fmt.Printf("   字节 %d: %d\n", i, text[i])
	}

	// 方法 2: 按 rune 遍历（推荐）
	fmt.Println("\n   方法 2: 按 rune 遍历（推荐）")
	for index, runeValue := range text {
		fmt.Printf("   索引 %d: %c (Unicode: U+%04X)\n", index, runeValue, runeValue)
	}

	// 8. 字符串拼接
	fmt.Println("\n8. 字符串拼接")

	// 方法 1: + 运算符
	s1 := "Hello"
	s2 := "World"
	result1 := s1 + " " + s2
	fmt.Printf("   方法 1 (+运算符): %s\n", result1)

	// 方法 2: fmt.Sprintf
	result2 := fmt.Sprintf("%s %s", s1, s2)
	fmt.Printf("   方法 2 (fmt.Sprintf): %s\n", result2)

	// 方法 3: strings.Join
	parts := []string{"Go", "is", "awesome"}
	result3 := strings.Join(parts, " ")
	fmt.Printf("   方法 3 (strings.Join): %s\n", result3)

	// 9. 常用字符串操作
	fmt.Println("\n9. 常用字符串操作（strings 包）")

	sample := "  Hello, Go Programming!  "
	fmt.Printf("   原字符串: \"%s\"\n", sample)
	fmt.Printf("   转大写: %s\n", strings.ToUpper(sample))
	fmt.Printf("   转小写: %s\n", strings.ToLower(sample))
	fmt.Printf("   去除空格: \"%s\"\n", strings.TrimSpace(sample))
	fmt.Printf("   是否包含 \"Go\": %t\n", strings.Contains(sample, "Go"))
	fmt.Printf("   是否以 \"  Hello\" 开头: %t\n", strings.HasPrefix(sample, "  Hello"))
	fmt.Printf("   是否以 \"!  \" 结尾: %t\n", strings.HasSuffix(sample, "!  "))
	fmt.Printf("   \"Go\" 的位置: %d\n", strings.Index(sample, "Go"))
	fmt.Printf("   替换 \"Go\" 为 \"Golang\": %s\n", strings.Replace(sample, "Go", "Golang", -1))

	// 分割字符串
	csv := "apple,banana,orange"
	fruits := strings.Split(csv, ",")
	fmt.Printf("   分割 \"%s\": %v\n", csv, fruits)

	// 10. 字符串和其他类型的转换
	fmt.Println("\n10. 字符串和其他类型的转换")

	// 字符串转数字
	numStr := "42"
	num, _ := strconv.Atoi(numStr)
	fmt.Printf("   字符串 \"%s\" → 整数 %d\n", numStr, num)

	floatStr := "3.14"
	floatNum, _ := strconv.ParseFloat(floatStr, 64)
	fmt.Printf("   字符串 \"%s\" → 浮点数 %.2f\n", floatStr, floatNum)

	// 数字转字符串
	n := 100
	nStr := strconv.Itoa(n)
	fmt.Printf("   整数 %d → 字符串 \"%s\"\n", n, nStr)

	// []byte 和 string 的转换
	bytes := []byte("Hello")
	strFromBytes := string(bytes)
	fmt.Printf("   []byte %v → string \"%s\"\n", bytes, strFromBytes)

	// 11. 最佳实践
	fmt.Println("\n11. 最佳实践")
	fmt.Println("   ✓ 使用双引号 \" \" 表示普通字符串")
	fmt.Println("   ✓ 使用反引号 ` ` 表示原始字符串（路径、正则、多行文本）")
	fmt.Println("   ✓ 使用 range 遍历字符串（自动处理 UTF-8）")
	fmt.Println("   ✓ 字符串是不可变的，拼接会创建新字符串")
	fmt.Println("   ✓ 大量拼接使用 strings.Builder 或 bytes.Buffer")
	fmt.Println("   ✓ 使用 strings 包进行字符串操作")
	fmt.Println("   ✓ 使用 strconv 包进行类型转换")
	fmt.Println("   ✓ 注意字节长度和字符数量的区别")

	fmt.Println("\n" + repeatString("=", 60))
}
