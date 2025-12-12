// Package constants - 布尔常量学习模块
//
// 本文件介绍 Go 语言中的布尔常量(Boolean Constants)。
// 布尔常量是最简单的常量类型,只有两个预定义值: true 和 false。
package constants

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-study2/internal/infrastructure/logger"
)

// GetBooleanContent 返回布尔常量相关的学习内容
func GetBooleanContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Boolean Constants (布尔常量) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("布尔常量是预定义的标识符 true 和 false,表示逻辑真值。\n")
	sb.WriteString("它们是无类型化的布尔常量,可以用于任何需要 bool 类型的地方。\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("const name [bool] = true | false\n")
	sb.WriteString("- true: 表示逻辑真\n")
	sb.WriteString("- false: 表示逻辑假\n\n")

	// 示例 1: 基本布尔常量声明
	sb.WriteString("【示例 1: 基本布尔常量声明】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const enabled = true\n")
	sb.WriteString("    const disabled = false\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(\"enabled:\", enabled)   // 输出: enabled: true\n")
	sb.WriteString("    fmt.Println(\"disabled:\", disabled) // 输出: disabled: false\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 布尔常量可以直接使用 true/false 字面量赋值。\n\n")

	// 示例 2: 类型化布尔常量
	sb.WriteString("【示例 2: 类型化布尔常量】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const isReady bool = true\n")
	sb.WriteString("    const isComplete bool = false\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 类型化常量有明确的类型\n")
	sb.WriteString("    fmt.Printf(\"isReady 类型: %T, 值: %v\\n\", isReady, isReady)\n")
	sb.WriteString("    // 输出: isReady 类型: bool, 值: true\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 可以显式指定 bool 类型,这称为类型化布尔常量。\n\n")

	// 示例 3: 布尔表达式常量
	sb.WriteString("【示例 3: 布尔表达式常量】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const a = 10\n")
	sb.WriteString("    const b = 20\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 常量表达式的结果也是常量\n")
	sb.WriteString("    const isGreater = a > b   // false\n")
	sb.WriteString("    const isEqual = a == b    // false\n")
	sb.WriteString("    const isLess = a < b      // true\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(\"a > b:\", isGreater)  // 输出: a > b: false\n")
	sb.WriteString("    fmt.Println(\"a == b:\", isEqual)   // 输出: a == b: false\n")
	sb.WriteString("    fmt.Println(\"a < b:\", isLess)     // 输出: a < b: true\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 比较运算的结果可以作为布尔常量,编译时求值。\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 尝试将非布尔值赋给布尔常量:\n")
	sb.WriteString("   const b = 1      // 错误: 这是整数常量,不是布尔常量\n")
	sb.WriteString("   const b bool = 1 // 编译错误: cannot use 1 as bool\n")
	sb.WriteString("\n")
	sb.WriteString("2. 布尔常量不能进行算术运算:\n")
	sb.WriteString("   const x = true + false // 编译错误: 布尔值不支持加法\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayBoolean 展示并解释 Go 语言中的布尔常量。
// 布尔常量只有 true 和 false 两个值,用于表示逻辑真值。
func DisplayBoolean() {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logger.LogBiz(context.Background(), "DisplayBoolean", map[string]interface{}{
			"operation": "display_boolean_constants",
			"result":    "success",
		}, nil, duration)
	}()

	fmt.Print(GetBooleanContent())
}
