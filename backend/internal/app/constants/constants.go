// Package constants 提供 Go 语言常量(Constants)相关知识的学习内容。
//
// 本包涵盖 12 个子主题:
//   - 基础常量类型: boolean, rune, integer, floating_point, complex, string
//   - 常量表达式和类型: expressions, typed_untyped
//   - 转换和内置函数: conversions, builtin_functions
//   - 特殊常量: iota, implementation_restrictions
//
// 使用方式:
//
//	// CLI 模式
//	constants.DisplayMenu(os.Stdin, os.Stdout, os.Stderr)
//
//	// 获取某个主题的内容 (HTTP 模式)
//	content := constants.GetBooleanContent()
package constants

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// DisplayMenu 显示 Constants 子菜单,允许用户选择特定的主题进行学习。
// 该函数实现了交互式菜单循环,用户可以:
//   - 选择 0-11 查看对应的主题内容
//   - 输入 'q' 返回主菜单
//
// 参数:
//   - stdin: 用于读取用户输入
//   - stdout: 用于输出菜单和正常信息
//   - stderr: 用于输出错误信息
func DisplayMenu(stdin io.Reader, stdout, stderr io.Writer) {
	reader := bufio.NewReader(stdin)

	// 定义子菜单选项映射(主题编号方案: 0-11 顺序对应 12 个常量主题)
	subMenu := map[string]string{
		"0":  "Boolean Constants (布尔常量)",
		"1":  "Rune Constants (符文常量)",
		"2":  "Integer Constants (整数常量)",
		"3":  "Floating-point Constants (浮点常量)",
		"4":  "Complex Constants (复数常量)",
		"5":  "String Constants (字符串常量)",
		"6":  "Constant Expressions (常量表达式)",
		"7":  "Typed and Untyped Constants (类型化/无类型化常量)",
		"8":  "Conversions (类型转换)",
		"9":  "Built-in Functions (内置函数)",
		"10": "Iota (iota 特性)",
		"11": "Implementation Restrictions (实现限制)",
	}

	// 定义主题执行函数映射(将子菜单选项映射到对应的显示函数)
	topicActions := map[string]func(){
		"0":  DisplayBoolean,
		"1":  DisplayRune,
		"2":  DisplayInteger,
		"3":  DisplayFloatingPoint,
		"4":  DisplayComplex,
		"5":  DisplayString,
		"6":  DisplayExpressions,
		"7":  DisplayTypedUntyped,
		"8":  DisplayConversions,
		"9":  DisplayBuiltinFunctions,
		"10": DisplayIota,
		"11": DisplayImplementationRestrictions,
	}

	for {
		fmt.Fprintln(stdout, "\nConstants 学习菜单")
		fmt.Fprintln(stdout, "---------------------------------")
		fmt.Fprintln(stdout, "请选择要学习的主题:")

		// 按顺序显示所有选项
		for i := 0; i <= 11; i++ {
			key := fmt.Sprintf("%d", i)
			fmt.Fprintf(stdout, "%s. %s\n", key, subMenu[key])
		}
		fmt.Fprintln(stdout, "q. 返回上级菜单")
		fmt.Fprint(stdout, "\n请输入您的选择: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(stderr, "读取输入错误: %v\n", err)
			return
		}
		choice := strings.TrimSpace(input)

		// 处理返回主菜单
		if choice == "q" {
			return
		}

		// 验证输入范围(0-11)并执行对应的主题显示函数
		if action, ok := topicActions[choice]; ok {
			// 执行对应的主题显示函数
			action()
			// 执行完毕后,循环会继续显示菜单,允许用户选择其他主题或返回
		} else {
			// 处理无效输入: 显示错误消息并重新提示
			fmt.Fprintln(stdout, "无效的选择,请重试。")
		}
	}
}
