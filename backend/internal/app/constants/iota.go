// Package constants - iota 特性学习模块
//
// 本文件介绍 Go 语言中的 iota 特性。
// iota 是预声明的标识符，在常量声明中自增，常用于枚举和位掩码。
package constants

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-study2/internal/infrastructure/logger"
)

// GetIotaContent 返回 iota 特性相关的学习内容
func GetIotaContent() string {
	var sb strings.Builder

	sb.WriteString("\n=== Iota (iota 特性) ===\n\n")

	// 概念说明
	sb.WriteString("【概念说明】\n")
	sb.WriteString("iota 是 Go 语言预声明的标识符，在常量声明中使用时表示自增的整数。\n")
	sb.WriteString("iota 的特性:\n")
	sb.WriteString("1. 每个 const 块中，iota 从 0 开始\n")
	sb.WriteString("2. 每个常量声明后，iota 自增 1\n")
	sb.WriteString("3. 同一行多个常量共享同一个 iota 值\n")
	sb.WriteString("4. 新的 const 块会重置 iota 为 0\n")
	sb.WriteString("5. iota 常用于枚举、位掩码、单位转换等场景\n\n")

	// 语法规则
	sb.WriteString("【语法规则】\n")
	sb.WriteString("基本用法:\n")
	sb.WriteString("  const (\n")
	sb.WriteString("    name1 = iota  // 0\n")
	sb.WriteString("    name2         // 1 (隐式重复 = iota)\n")
	sb.WriteString("    name3         // 2\n")
	sb.WriteString("  )\n")
	sb.WriteString("\n")
	sb.WriteString("跳过值:\n")
	sb.WriteString("  const (\n")
	sb.WriteString("    _ = iota  // 跳过 0\n")
	sb.WriteString("    name1     // 1\n")
	sb.WriteString("  )\n")
	sb.WriteString("\n")
	sb.WriteString("表达式复用:\n")
	sb.WriteString("  const (\n")
	sb.WriteString("    name1 = iota * 2  // 0\n")
	sb.WriteString("    name2              // 2\n")
	sb.WriteString("  )\n\n")

	// 示例 1: 基本枚举
	sb.WriteString("【示例 1: 基本枚举】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        Sunday = iota    // 0\n")
	sb.WriteString("        Monday            // 1\n")
	sb.WriteString("        Tuesday           // 2\n")
	sb.WriteString("        Wednesday         // 3\n")
	sb.WriteString("        Thursday          // 4\n")
	sb.WriteString("        Friday            // 5\n")
	sb.WriteString("        Saturday          // 6\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(Sunday, Monday, Saturday)  // 输出: 0 1 6\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 使用枚举\n")
	sb.WriteString("    today := Monday\n")
	sb.WriteString("    fmt.Println(\"今天是:\", today)  // 输出: 今天是: 1\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: iota 从 0 开始自增，常用于定义枚举值。\n\n")

	// 示例 2: 跳过值
	sb.WriteString("【示例 2: 跳过值】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        _ = iota  // 跳过 0\n")
	sb.WriteString("        KB = 1 << (10 * iota)  // 1 << 10 = 1024\n")
	sb.WriteString("        MB                      // 1 << 20 = 1048576\n")
	sb.WriteString("        GB                      // 1 << 30 = 1073741824\n")
	sb.WriteString("        TB                      // 1 << 40\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(KB, MB, GB)  // 输出: 1024 1048576 1073741824\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 使用下划线跳过不需要的值\n")
	sb.WriteString("    fmt.Printf(\"1 KB = %d bytes\\n\", KB)\n")
	sb.WriteString("    // 输出: 1 KB = 1024 bytes\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 使用下划线 _ 可以跳过 iota 的某个值，常用于单位转换。\n\n")

	// 示例 3: 位掩码
	sb.WriteString("【示例 3: 位掩码】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        Read = 1 << iota   // 1 << 0 = 1\n")
	sb.WriteString("        Write               // 1 << 1 = 2\n")
	sb.WriteString("        Execute             // 1 << 2 = 4\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 组合权限\n")
	sb.WriteString("    permission := Read | Write  // 3 (1 | 2)\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(\"权限值:\", permission)  // 输出: 权限值: 3\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 检查权限\n")
	sb.WriteString("    if permission&Read != 0 {\n")
	sb.WriteString("        fmt.Println(\"有读权限\")\n")
	sb.WriteString("    }\n")
	sb.WriteString("    if permission&Write != 0 {\n")
	sb.WriteString("        fmt.Println(\"有写权限\")\n")
	sb.WriteString("    }\n")
	sb.WriteString("    // 输出:\n")
	sb.WriteString("    // 有读权限\n")
	sb.WriteString("    // 有写权限\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: iota 配合位运算可以定义位掩码，用于权限控制等场景。\n\n")

	// 示例 4: 表达式复用
	sb.WriteString("【示例 4: 表达式复用】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        a = iota * 2  // 0 * 2 = 0\n")
	sb.WriteString("        b             // 1 * 2 = 2\n")
	sb.WriteString("        c             // 2 * 2 = 4\n")
	sb.WriteString("        d             // 3 * 2 = 6\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(a, b, c, d)  // 输出: 0 2 4 6\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 表达式会被复用\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        x = iota + 10  // 0 + 10 = 10\n")
	sb.WriteString("        y               // 1 + 10 = 11\n")
	sb.WriteString("        z               // 2 + 10 = 12\n")
	sb.WriteString("    )\n")
	sb.WriteString("    fmt.Println(x, y, z)  // 输出: 10 11 12\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 如果第一个常量使用表达式，后续常量会复用该表达式。\n\n")

	// 示例 5: 多个常量共享 iota
	sb.WriteString("【示例 5: 多个常量共享 iota】\n")
	sb.WriteString("```go\n")
	sb.WriteString("package main\n\n")
	sb.WriteString("import \"fmt\"\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        a, b = iota, iota + 10  // 0, 10\n")
	sb.WriteString("        c, d                   // 1, 11\n")
	sb.WriteString("        e, f                   // 2, 12\n")
	sb.WriteString("    )\n")
	sb.WriteString("    \n")
	sb.WriteString("    fmt.Println(a, b, c, d, e, f)  // 输出: 0 10 1 11 2 12\n")
	sb.WriteString("    \n")
	sb.WriteString("    // 同一行的常量共享同一个 iota 值\n")
	sb.WriteString("    const (\n")
	sb.WriteString("        x, y, z = iota, iota, iota  // 0, 0, 0\n")
	sb.WriteString("        u, v, w                     // 1, 1, 1\n")
	sb.WriteString("    )\n")
	sb.WriteString("    fmt.Println(x, y, z, u, v, w)  // 输出: 0 0 0 1 1 1\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n")
	sb.WriteString("说明: 同一行声明的多个常量共享同一个 iota 值。\n\n")

	// 常见错误
	sb.WriteString("【常见错误】\n")
	sb.WriteString("1. 混淆 iota 在不同 const 块中的值:\n")
	sb.WriteString("   const (\n")
	sb.WriteString("       a = iota  // 0\n")
	sb.WriteString("   )\n")
	sb.WriteString("   const (\n")
	sb.WriteString("       b = iota  // 0 (重新开始)\n")
	sb.WriteString("   )\n")
	sb.WriteString("\n")
	sb.WriteString("2. 不理解表达式复用:\n")
	sb.WriteString("   const (\n")
	sb.WriteString("       a = iota * 2  // 0\n")
	sb.WriteString("       b            // 2 (不是 1)\n")
	sb.WriteString("   )\n")
	sb.WriteString("\n")
	sb.WriteString("3. iota 只能在常量声明中使用:\n")
	sb.WriteString("   // var x = iota  // 编译错误: iota 只能用于常量\n")
	sb.WriteString("\n")
	sb.WriteString("4. 同一行常量共享 iota:\n")
	sb.WriteString("   const (\n")
	sb.WriteString("       a, b = iota, iota  // 都是 0\n")
	sb.WriteString("       c, d                // 都是 1\n")
	sb.WriteString("   )\n")
	sb.WriteString("\n")

	return sb.String()
}

// DisplayIota 展示并解释 Go 语言中的 iota 特性。
// iota 是预声明的标识符，在常量声明中自增，常用于枚举和位掩码。
func DisplayIota() {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logger.LogBiz(context.Background(), "DisplayIota", map[string]interface{}{
			"operation": "display_iota_feature",
			"result":    "success",
		}, nil, duration)
	}()

	fmt.Print(GetIotaContent())
}
