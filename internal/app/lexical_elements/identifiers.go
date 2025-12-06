package lexical_elements

import "fmt"

// DisplayIdentifiers 展示 Go 语言中标识符的命名规则和示例
// 标识符用于命名程序实体，如变量、类型、函数等
func DisplayIdentifiers() {
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("【词法元素 - 标识符 (Identifiers)】")
	fmt.Println(repeatString("=", 60))

	// 1. 标识符的定义
	fmt.Println("\n1. 标识符的定义")
	fmt.Println("   标识符是由字母和数字组成的序列，用于命名变量、常量、函数、类型等")
	fmt.Println("   规则：")
	fmt.Println("   - 必须以字母（Unicode 字母）或下划线 _ 开头")
	fmt.Println("   - 后续字符可以是字母、数字或下划线")
	fmt.Println("   - 区分大小写")

	// 2. 合法的标识符示例
	fmt.Println("\n2. 合法的标识符示例")

	// 英文标识符
	userName := "Alice"
	user_age := 25
	_privateVar := "private"
	fmt.Printf("   userName = %s\n", userName)
	fmt.Printf("   user_age = %d\n", user_age)
	fmt.Printf("   _privateVar = %s\n", _privateVar)

	// Unicode 标识符（支持中文等）
	姓名 := "张三"
	年龄 := 30
	fmt.Printf("   姓名 = %s  // Unicode 标识符（中文）\n", 姓名)
	fmt.Printf("   年龄 = %d\n", 年龄)

	// 包含数字的标识符
	var value1, value2, value3 int = 10, 20, 30
	fmt.Printf("   value1=%d, value2=%d, value3=%d\n", value1, value2, value3)

	// 3. 非法的标识符示例
	fmt.Println("\n3. 非法的标识符示例（以下是错误示例，不会运行）")
	fmt.Println("   // 1abc      // 错误：不能以数字开头")
	fmt.Println("   // user-name // 错误：不能包含连字符")
	fmt.Println("   // for       // 错误：不能使用关键字")
	fmt.Println("   // user name // 错误：不能包含空格")

	// 4. 导出与未导出标识符
	fmt.Println("\n4. 导出与未导出标识符（可见性规则）")
	fmt.Println("   Go 使用标识符的首字母大小写来控制可见性：")
	fmt.Println("   - 首字母大写：导出的（Public），包外可访问")
	fmt.Println("   - 首字母小写：未导出的（Private），仅包内可访问")

	fmt.Println("\n   示例：")
	fmt.Println("   type Person struct {")
	fmt.Println("       Name string    // 导出字段（首字母大写）")
	fmt.Println("       age  int       // 未导出字段（首字母小写）")
	fmt.Println("   }")

	type Person struct {
		Name string // 导出
		age  int    // 未导出
	}
	p := Person{Name: "Bob", age: 28}
	fmt.Printf("   p.Name = %s (可从包外访问)\n", p.Name)
	fmt.Printf("   p.age = %d (仅包内可访问)\n", p.age)

	// 5. 空白标识符
	fmt.Println("\n5. 空白标识符 (_)")
	fmt.Println("   特殊标识符 _ 用于忽略不需要的值")

	// 示例：忽略函数返回值
	_, err := divide(10, 2)
	if err != nil {
		fmt.Printf("   错误: %v\n", err)
	} else {
		fmt.Println("   使用 _ 忽略了除法结果，只检查错误")
	}

	// 示例：在 range 循环中忽略索引
	fmt.Println("\n   在循环中使用 _:")
	numbers := []int{10, 20, 30}
	for _, num := range numbers {
		fmt.Printf("   数字: %d\n", num)
	}

	// 6. 命名最佳实践
	fmt.Println("\n6. 命名最佳实践")
	fmt.Println("   ✓ 使用驼峰命名法（camelCase 或 PascalCase）")
	fmt.Println("   ✓ 名称要有意义，避免使用 a, b, c 等无意义名称")
	fmt.Println("   ✓ 缩写词保持一致大小写（如 URL, HTTP, ID）")
	fmt.Println("   ✓ 包名使用小写单词，不使用下划线")
	fmt.Println("   ✓ 接口名通常以 -er 结尾（如 Reader, Writer）")
	fmt.Println("   ✓ 避免使用预定义标识符（如 int, string, true）")

	// 7. 预定义标识符
	fmt.Println("\n7. 预定义标识符（不是关键字，但不建议重新定义）")
	fmt.Println("   类型: bool, byte, complex64, complex128, error, float32, float64")
	fmt.Println("         int, int8, int16, int32, int64, rune, string")
	fmt.Println("         uint, uint8, uint16, uint32, uint64, uintptr")
	fmt.Println("   常量: true, false, iota, nil")
	fmt.Println("   函数: append, cap, close, complex, copy, delete, imag, len")
	fmt.Println("         make, new, panic, print, println, real, recover")

	fmt.Println("\n" + repeatString("=", 60))
}

// divide 辅助函数：演示多返回值
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}
