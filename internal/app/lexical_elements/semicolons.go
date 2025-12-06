package lexical_elements

import (
	"fmt"
	"strings"
)

// GetSemicolonsContent 返回分号相关的学习内容
func GetSemicolonsContent() string {
	var sb strings.Builder

	sb.WriteString("\n" + repeatString("=", 60) + "\n")
	sb.WriteString("【词法元素 - 分号 (Semicolons)】\n")
	sb.WriteString(repeatString("=", 60) + "\n")

	// 1. 自动分号插入规则
	sb.WriteString("\n1. 自动分号插入规则 (Automatic Semicolon Insertion)\n")
	sb.WriteString("   Go 编译器会在以下情况自动插入分号：\n")
	sb.WriteString("   - 当一行的最后一个标记是标识符、基本字面量、或以下标记之一时：\n")
	sb.WriteString("     break, continue, fallthrough, return\n")
	sb.WriteString("     ++, --, ), ], }\n")

	// 示例：自动插入分号
	sb.WriteString("\n   示例代码：\n")
	x := 10
	y := 20
	sb.WriteString(fmt.Sprintf("   x := %d  // 编译器在行尾自动插入分号\n", x))
	sb.WriteString(fmt.Sprintf("   y := %d  // 编译器在行尾自动插入分号\n", y))
	sb.WriteString(fmt.Sprintf("   结果: x=%d, y=%d\n", x, y))

	// 2. 显式分号的使用场景
	sb.WriteString("\n2. 显式分号的使用场景\n")
	sb.WriteString("   虽然大多数情况下不需要显式分号，但在某些场景下必须使用：\n")

	// 示例：在同一行写多个语句
	sb.WriteString("\n   场景 A: 在同一行写多个语句（不推荐，但语法允许）\n")
	a := 1
	b := 2
	c := a + b
	sb.WriteString(fmt.Sprintf("   a := 1; b := 2; c := a + b  // 结果: c=%d\n", c))

	// 示例：for 循环的子句分隔
	sb.WriteString("\n   场景 B: for 循环的子句必须用分号分隔\n")
	sb.WriteString("   for i := 0; i < 3; i++ {\n")
	for i := 0; i < 3; i++ {
		sb.WriteString(fmt.Sprintf("       迭代 %d\n", i))
	}
	sb.WriteString("   }\n")

	// 3. 常见的分号相关错误
	sb.WriteString("\n3. 常见的分号相关错误\n")
	sb.WriteString("   错误示例 1: 左花括号不能单独成行\n")
	sb.WriteString("   // 错误写法：\n")
	sb.WriteString("   // if x > 0\n")
	sb.WriteString("   // {  // 编译器会在 if x > 0 后自动插入分号，导致语法错误\n")
	sb.WriteString("   // }\n")
	sb.WriteString("   // 正确写法：\n")
	sb.WriteString("   if x > 0 {  // 左花括号必须在同一行\n")
	sb.WriteString("       // 代码块\n")
	sb.WriteString("   }\n")

	sb.WriteString("\n   错误示例 2: return 语句的返回值不能单独成行\n")
	sb.WriteString("   // 错误写法：\n")
	sb.WriteString("   // return\n")
	sb.WriteString("   //     x + y  // 编译器会在 return 后自动插入分号\n")
	sb.WriteString("   // 正确写法：\n")
	sb.WriteString("   // return x + y  // 或者 return (\n")
	sb.WriteString("   //                //     x + y\n")
	sb.WriteString("   //                // )\n")

	// 4. 实际运行示例
	sb.WriteString("\n4. 实际运行示例：函数返回值\n")
	result := calculateSum(5, 3)
	sb.WriteString(fmt.Sprintf("   calculateSum(5, 3) = %d\n", result))

	// 5. 最佳实践
	sb.WriteString("\n5. 最佳实践\n")
	sb.WriteString("   ✓ 让编译器自动插入分号，不要手动添加\n")
	sb.WriteString("   ✓ 左花括号 { 始终放在同一行\n")
	sb.WriteString("   ✓ 避免在同一行写多个语句\n")
	sb.WriteString("   ✓ 使用 gofmt 工具自动格式化代码\n")
	sb.WriteString("   ✓ return 语句的返回值要在同一行开始\n")

	sb.WriteString("\n" + repeatString("=", 60) + "\n")

	return sb.String()
}

// DisplaySemicolons 展示 Go 语言中分号的使用规则和示例
// Go 的正式语法使用分号来终止语句，但在实际编码中，大多数分号是可选的
func DisplaySemicolons() {
	fmt.Print(GetSemicolonsContent())
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
