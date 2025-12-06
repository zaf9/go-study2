package lexical_elements

import (
	"fmt"
	"strconv"
	"strings"
)

// GetStringsContent 返回字符串字面量相关的学习内容
func GetStringsContent() string {
	var sb strings.Builder

	sb.WriteString("\n" + repeatString("=", 60) + "\n")
	sb.WriteString("【词法元素 - 字符串字面量 (String Literals)】\n")
	sb.WriteString(repeatString("=", 60) + "\n")

	// 1. 解释型字符串（双引号）
	sb.WriteString("\n1. 解释型字符串（Interpreted Strings）\n")
	sb.WriteString("   使用双引号 \" \" 括起来，支持转义序列\n")

	str1 := "Hello, World!"
	str2 := "你好，世界！"
	str3 := "第一行\n第二行" // 包含换行符

	sb.WriteString(fmt.Sprintf("   \"Hello, World!\" = %s\n", str1))
	sb.WriteString(fmt.Sprintf("   \"你好，世界！\" = %s\n", str2))
	sb.WriteString("   \"第一行\\n第二行\" 的效果:\n")
	sb.WriteString(str3 + "\n")

	// 2. 原始字符串（反引号）
	sb.WriteString("\n2. 原始字符串（Raw Strings）\n")
	sb.WriteString("   使用反引号 ` ` 括起来，不解释转义序列\n")
	sb.WriteString("   常用于正则表达式、文件路径、多行文本等\n")

	raw1 := `C:\Users\Documents\file.txt` // Windows 路径
	raw2 := `\n 不会被解释为换行符`
	raw3 := `这是第一行
这是第二行
这是第三行` // 多行字符串

	sb.WriteString(fmt.Sprintf("   `C:\\Users\\Documents\\file.txt` = %s\n", raw1))
	sb.WriteString(fmt.Sprintf("   `\\n 不会被解释为换行符` = %s\n", raw2))
	sb.WriteString("\n   多行原始字符串:\n")
	sb.WriteString(raw3 + "\n")

	// 3. 转义序列
	sb.WriteString("\n3. 转义序列（仅在解释型字符串中有效）\n")

	sb.WriteString("   常用转义序列:\n")
	sb.WriteString("   \\n  - 换行\n")
	sb.WriteString("   \\t  - 制表符\n")
	sb.WriteString("   \\r  - 回车\n")
	sb.WriteString("   \\\\  - 反斜杠\n")
	sb.WriteString("   \\\"  - 双引号\n")
	sb.WriteString("   \\'  - 单引号\n")

	// 演示
	escaped := "姓名:\t张三\n年龄:\t25\n城市:\t\"北京\""
	sb.WriteString("\n   转义序列示例:\n")
	sb.WriteString(escaped + "\n")

	// 4. Unicode 转义
	sb.WriteString("\n4. Unicode 转义（在解释型字符串中）\n")

	unicode1 := "\u4E2D\u6587"    // Unicode: 中文
	unicode2 := "\U0001F600"      // Unicode: 😀
	hex := "\x48\x65\x6C\x6C\x6F" // 十六进制: Hello

	sb.WriteString(fmt.Sprintf("   \"\\u4E2D\\u6587\" = %s\n", unicode1))
	sb.WriteString(fmt.Sprintf("   \"\\U0001F600\" = %s\n", unicode2))
	sb.WriteString(fmt.Sprintf("   \"\\x48\\x65\\x6C\\x6C\\x6F\" = %s\n", hex))

	// 5. 字符串的不可变性
	sb.WriteString("\n5. 字符串的不可变性\n")
	sb.WriteString("   Go 中的字符串是不可变的，不能修改字符串中的字符\n")

	s := "Hello"
	sb.WriteString(fmt.Sprintf("   原字符串: %s\n", s))
	// s[0] = 'h'  // 错误！不能修改字符串

	// 正确做法：创建新字符串
	newS := "h" + s[1:]
	sb.WriteString(fmt.Sprintf("   新字符串: %s (通过拼接创建)\n", newS))

	// 6. 字符串长度和索引
	sb.WriteString("\n6. 字符串长度和索引\n")

	str := "Go语言"
	sb.WriteString(fmt.Sprintf("   字符串: %s\n", str))
	sb.WriteString(fmt.Sprintf("   字节长度 len(): %d (UTF-8 编码)\n", len(str)))
	sb.WriteString(fmt.Sprintf("   字符数量: %d (rune 数量)\n", len([]rune(str))))

	// 字节索引
	sb.WriteString("\n   按字节索引（可能不是完整字符）:\n")
	for i := 0; i < len(str); i++ {
		sb.WriteString(fmt.Sprintf("   str[%d] = %d (%c)\n", i, str[i], str[i]))
	}

	// 7. 遍历字符串
	sb.WriteString("\n7. 遍历字符串\n")

	text := "Go编程"

	// 方法 1: 按字节遍历
	sb.WriteString("\n   方法 1: 按字节遍历（不推荐用于 Unicode）\n")
	for i := 0; i < len(text); i++ {
		sb.WriteString(fmt.Sprintf("   字节 %d: %d\n", i, text[i]))
	}

	// 方法 2: 按 rune 遍历（推荐）
	sb.WriteString("\n   方法 2: 按 rune 遍历（推荐）\n")
	for index, runeValue := range text {
		sb.WriteString(fmt.Sprintf("   索引 %d: %c (Unicode: U+%04X)\n", index, runeValue, runeValue))
	}

	// 8. 字符串拼接
	sb.WriteString("\n8. 字符串拼接\n")

	// 方法 1: + 运算符
	s1 := "Hello"
	s2 := "World"
	result1 := s1 + " " + s2
	sb.WriteString(fmt.Sprintf("   方法 1 (+运算符): %s\n", result1))

	// 方法 2: fmt.Sprintf
	result2 := fmt.Sprintf("%s %s", s1, s2)
	sb.WriteString(fmt.Sprintf("   方法 2 (fmt.Sprintf): %s\n", result2))

	// 方法 3: strings.Join
	parts := []string{"Go", "is", "awesome"}
	result3 := strings.Join(parts, " ")
	sb.WriteString(fmt.Sprintf("   方法 3 (strings.Join): %s\n", result3))

	// 9. 常用字符串操作
	sb.WriteString("\n9. 常用字符串操作（strings 包）\n")

	sample := "  Hello, Go Programming!  "
	sb.WriteString(fmt.Sprintf("   原字符串: \"%s\"\n", sample))
	sb.WriteString(fmt.Sprintf("   转大写: %s\n", strings.ToUpper(sample)))
	sb.WriteString(fmt.Sprintf("   转小写: %s\n", strings.ToLower(sample)))
	sb.WriteString(fmt.Sprintf("   去除空格: \"%s\"\n", strings.TrimSpace(sample)))
	sb.WriteString(fmt.Sprintf("   是否包含 \"Go\": %t\n", strings.Contains(sample, "Go")))
	sb.WriteString(fmt.Sprintf("   是否以 \"  Hello\" 开头: %t\n", strings.HasPrefix(sample, "  Hello")))
	sb.WriteString(fmt.Sprintf("   是否以 \"!  \" 结尾: %t\n", strings.HasSuffix(sample, "!  ")))
	sb.WriteString(fmt.Sprintf("   \"Go\" 的位置: %d\n", strings.Index(sample, "Go")))
	sb.WriteString(fmt.Sprintf("   替换 \"Go\" 为 \"Golang\": %s\n", strings.Replace(sample, "Go", "Golang", -1)))

	// 分割字符串
	csv := "apple,banana,orange"
	fruits := strings.Split(csv, ",")
	sb.WriteString(fmt.Sprintf("   分割 \"%s\": %v\n", csv, fruits))

	// 10. 字符串和其他类型的转换
	sb.WriteString("\n10. 字符串和其他类型的转换\n")

	// 字符串转数字
	numStr := "42"
	num, _ := strconv.Atoi(numStr)
	sb.WriteString(fmt.Sprintf("   字符串 \"%s\" → 整数 %d\n", numStr, num))

	floatStr := "3.14"
	floatNum, _ := strconv.ParseFloat(floatStr, 64)
	sb.WriteString(fmt.Sprintf("   字符串 \"%s\" → 浮点数 %.2f\n", floatStr, floatNum))

	// 数字转字符串
	n := 100
	nStr := strconv.Itoa(n)
	sb.WriteString(fmt.Sprintf("   整数 %d → 字符串 \"%s\"\n", n, nStr))

	// []byte 和 string 的转换
	bytes := []byte("Hello")
	strFromBytes := string(bytes)
	sb.WriteString(fmt.Sprintf("   []byte %v → string \"%s\"\n", bytes, strFromBytes))

	// 11. 最佳实践
	sb.WriteString("\n11. 最佳实践\n")
	sb.WriteString("   ✓ 使用双引号 \" \" 表示普通字符串\n")
	sb.WriteString("   ✓ 使用反引号 ` ` 表示原始字符串（路径、正则、多行文本）\n")
	sb.WriteString("   ✓ 使用 range 遍历字符串（自动处理 UTF-8）\n")
	sb.WriteString("   ✓ 字符串是不可变的，拼接会创建新字符串\n")
	sb.WriteString("   ✓ 大量拼接使用 strings.Builder 或 bytes.Buffer\n")
	sb.WriteString("   ✓ 使用 strings 包进行字符串操作\n")
	sb.WriteString("   ✓ 使用 strconv 包进行类型转换\n")
	sb.WriteString("   ✓ 注意字节长度和字符数量的区别\n")

	sb.WriteString("\n" + repeatString("=", 60) + "\n")

	return sb.String()
}

// DisplayStrings 展示 Go 语言中字符串字面量的各种形式和操作
// Go 字符串是 UTF-8 编码的字节序列，不可变
func DisplayStrings() {
	fmt.Print(GetStringsContent())
}
