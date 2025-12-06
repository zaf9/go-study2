package lexical_elements

import (
	"fmt"
	"strings"
)

// GetIdentifiersContent 返回标识符相关的学习内容
func GetIdentifiersContent() string {
	var sb strings.Builder

	sb.WriteString("\n" + repeatString("=", 60) + "\n")
	sb.WriteString("【词法元素 - 标识符 (Identifiers)】\n")
	sb.WriteString(repeatString("=", 60) + "\n")

	// 1. 标识符的定义
	sb.WriteString("\n1. 标识符的定义\n")
	sb.WriteString("   标识符是由字母和数字组成的序列，用于命名变量、常量、函数、类型等\n")
	sb.WriteString("   规则：\n")
	sb.WriteString("   - 必须以字母（Unicode 字母）或下划线 _ 开头\n")
	sb.WriteString("   - 后续字符可以是字母、数字或下划线\n")
	sb.WriteString("   - 区分大小写\n")

	// 2. 合法的标识符示例
	sb.WriteString("\n2. 合法的标识符示例\n")

	// 英文标识符
	userName := "Alice"
	user_age := 25
	_privateVar := "private"
	sb.WriteString(fmt.Sprintf("   userName = %s\n", userName))
	sb.WriteString(fmt.Sprintf("   user_age = %d\n", user_age))
	sb.WriteString(fmt.Sprintf("   _privateVar = %s\n", _privateVar))

	// Unicode 标识符（支持中文等）
	姓名 := "张三"
	年龄 := 30
	sb.WriteString(fmt.Sprintf("   姓名 = %s  // Unicode 标识符（中文）\n", 姓名))
	sb.WriteString(fmt.Sprintf("   年龄 = %d\n", 年龄))

	// 包含数字的标识符
	var value1, value2, value3 int = 10, 20, 30
	sb.WriteString(fmt.Sprintf("   value1=%d, value2=%d, value3=%d\n", value1, value2, value3))

	// 3. 非法的标识符示例
	sb.WriteString("\n3. 非法的标识符示例（以下是错误示例，不会运行）\n")
	sb.WriteString("   // 1abc      // 错误：不能以数字开头\n")
	sb.WriteString("   // user-name // 错误：不能包含连字符\n")
	sb.WriteString("   // for       // 错误：不能使用关键字\n")
	sb.WriteString("   // user name // 错误：不能包含空格\n")

	// 4. 导出与未导出标识符
	sb.WriteString("\n4. 导出与未导出标识符（可见性规则）\n")
	sb.WriteString("   Go 使用标识符的首字母大小写来控制可见性：\n")
	sb.WriteString("   - 首字母大写：导出的（Public），包外可访问\n")
	sb.WriteString("   - 首字母小写：未导出的（Private），仅包内可访问\n")

	sb.WriteString("\n   示例：\n")
	sb.WriteString("   type Person struct {\n")
	sb.WriteString("       Name string    // 导出字段（首字母大写）\n")
	sb.WriteString("       age  int       // 未导出字段（首字母小写）\n")
	sb.WriteString("   }\n")

	type Person struct {
		Name string // 导出
		age  int    // 未导出
	}
	p := Person{Name: "Bob", age: 28}
	sb.WriteString(fmt.Sprintf("   p.Name = %s (可从包外访问)\n", p.Name))
	sb.WriteString(fmt.Sprintf("   p.age = %d (仅包内可访问)\n", p.age))

	// 5. 空白标识符
	sb.WriteString("\n5. 空白标识符 (_)\n")
	sb.WriteString("   特殊标识符 _ 用于忽略不需要的值\n")

	// 示例：忽略函数返回值
	_, err := divide(10, 2)
	if err != nil {
		sb.WriteString(fmt.Sprintf("   错误: %v\n", err))
	} else {
		sb.WriteString("   使用 _ 忽略了除法结果，只检查错误\n")
	}

	// 示例：在 range 循环中忽略索引
	sb.WriteString("\n   在循环中使用 _:\n")
	numbers := []int{10, 20, 30}
	for _, num := range numbers {
		sb.WriteString(fmt.Sprintf("   数字: %d\n", num))
	}

	// 6. 命名最佳实践
	sb.WriteString("\n6. 命名最佳实践\n")
	sb.WriteString("   ✓ 使用驼峰命名法（camelCase 或 PascalCase）\n")
	sb.WriteString("   ✓ 名称要有意义，避免使用 a, b, c 等无意义名称\n")
	sb.WriteString("   ✓ 缩写词保持一致大小写（如 URL, HTTP, ID）\n")
	sb.WriteString("   ✓ 包名使用小写单词，不使用下划线\n")
	sb.WriteString("   ✓ 接口名通常以 -er 结尾（如 Reader, Writer）\n")
	sb.WriteString("   ✓ 避免使用预定义标识符（如 int, string, true）\n")

	// 7. 预定义标识符
	sb.WriteString("\n7. 预定义标识符（不是关键字，但不建议重新定义）\n")
	sb.WriteString("   类型: bool, byte, complex64, complex128, error, float32, float64\n")
	sb.WriteString("         int, int8, int16, int32, int64, rune, string\n")
	sb.WriteString("         uint, uint8, uint16, uint32, uint64, uintptr\n")
	sb.WriteString("   常量: true, false, iota, nil\n")
	sb.WriteString("   函数: append, cap, close, complex, copy, delete, imag, len\n")
	sb.WriteString("         make, new, panic, print, println, real, recover\n")

	sb.WriteString("\n" + repeatString("=", 60) + "\n")

	return sb.String()
}

// DisplayIdentifiers 展示 Go 语言中标识符的命名规则和示例
// 标识符用于命名程序实体，如变量、类型、函数等
func DisplayIdentifiers() {
	fmt.Print(GetIdentifiersContent())
}

// divide 辅助函数：演示多返回值
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}
