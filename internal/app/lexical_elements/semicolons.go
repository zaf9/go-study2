package lexical_elements

import "fmt"

// DisplaySemicolons 展示 Go 语言中分号的使用规则和示例
// Go 的正式语法使用分号来终止语句，但在实际编码中，大多数分号是可选的
func DisplaySemicolons() {
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("【词法元素 - 分号 (Semicolons)】")
	fmt.Println(repeatString("=", 60))

	// 1. 自动分号插入规则
	fmt.Println("\n1. 自动分号插入规则 (Automatic Semicolon Insertion)")
	fmt.Println("   Go 编译器会在以下情况自动插入分号：")
	fmt.Println("   - 当一行的最后一个标记是标识符、基本字面量、或以下标记之一时：")
	fmt.Println("     break, continue, fallthrough, return")
	fmt.Println("     ++, --, ), ], }")

	// 示例：自动插入分号
	fmt.Println("\n   示例代码：")
	x := 10
	y := 20
	fmt.Printf("   x := %d  // 编译器在行尾自动插入分号\n", x)
	fmt.Printf("   y := %d  // 编译器在行尾自动插入分号\n", y)
	fmt.Printf("   结果: x=%d, y=%d\n", x, y)

	// 2. 显式分号的使用场景
	fmt.Println("\n2. 显式分号的使用场景")
	fmt.Println("   虽然大多数情况下不需要显式分号，但在某些场景下必须使用：")

	// 示例：在同一行写多个语句
	fmt.Println("\n   场景 A: 在同一行写多个语句（不推荐，但语法允许）")
	a := 1
	b := 2
	c := a + b
	fmt.Printf("   a := 1; b := 2; c := a + b  // 结果: c=%d\n", c)

	// 示例：for 循环的子句分隔
	fmt.Println("\n   场景 B: for 循环的子句必须用分号分隔")
	fmt.Println("   for i := 0; i < 3; i++ {")
	for i := 0; i < 3; i++ {
		fmt.Printf("       迭代 %d\n", i)
	}
	fmt.Println("   }")

	// 3. 常见的分号相关错误
	fmt.Println("\n3. 常见的分号相关错误")
	fmt.Println("   错误示例 1: 左花括号不能单独成行")
	fmt.Println("   // 错误写法：")
	fmt.Println("   // if x > 0")
	fmt.Println("   // {  // 编译器会在 if x > 0 后自动插入分号，导致语法错误")
	fmt.Println("   // }")
	fmt.Println("   // 正确写法：")
	fmt.Println("   if x > 0 {  // 左花括号必须在同一行")
	fmt.Println("       // 代码块")
	fmt.Println("   }")

	fmt.Println("\n   错误示例 2: return 语句的返回值不能单独成行")
	fmt.Println("   // 错误写法：")
	fmt.Println("   // return")
	fmt.Println("   //     x + y  // 编译器会在 return 后自动插入分号")
	fmt.Println("   // 正确写法：")
	fmt.Println("   // return x + y  // 或者 return (")
	fmt.Println("   //                //     x + y")
	fmt.Println("   //                // )")

	// 4. 实际运行示例
	fmt.Println("\n4. 实际运行示例：函数返回值")
	result := calculateSum(5, 3)
	fmt.Printf("   calculateSum(5, 3) = %d\n", result)

	// 5. 最佳实践
	fmt.Println("\n5. 最佳实践")
	fmt.Println("   ✓ 让编译器自动插入分号，不要手动添加")
	fmt.Println("   ✓ 左花括号 { 始终放在同一行")
	fmt.Println("   ✓ 避免在同一行写多个语句")
	fmt.Println("   ✓ 使用 gofmt 工具自动格式化代码")
	fmt.Println("   ✓ return 语句的返回值要在同一行开始")

	fmt.Println("\n" + repeatString("=", 60))
}

// calculateSum 辅助函数：演示正确的 return 语句写法
func calculateSum(a, b int) int {
	return a + b // 返回值在同一行
}

// repeatString 辅助函数：重复字符串
func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
