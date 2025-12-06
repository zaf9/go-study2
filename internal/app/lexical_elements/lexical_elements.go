package lexical_elements

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// DisplayMenu 显示词法元素子菜单，允许用户选择特定的主题进行学习。
// 该函数实现了交互式菜单循环，用户可以：
//   - 选择 0-10 查看对应的主题内容
//   - 输入 'q' 返回主菜单
//
// 参数：
//   - stdin: 用于读取用户输入
//   - stdout: 用于输出菜单和正常信息
//   - stderr: 用于输出错误信息
func DisplayMenu(stdin io.Reader, stdout, stderr io.Writer) {
	reader := bufio.NewReader(stdin)

	// 定义子菜单选项映射（主题编号方案：0-10 顺序对应 11 个词法元素主题）
	subMenu := map[string]string{
		"0":  "Comments (注释)",
		"1":  "Tokens (标记)",
		"2":  "Semicolons (分号)",
		"3":  "Identifiers (标识符)",
		"4":  "Keywords (关键字)",
		"5":  "Operators (运算符)",
		"6":  "Integers (整数)",
		"7":  "Floats (浮点数)",
		"8":  "Imaginary (虚数)",
		"9":  "Runes (符文)",
		"10": "Strings (字符串)",
	}

	// 定义主题执行函数映射（将子菜单选项映射到对应的显示函数）
	topicActions := map[string]func(){
		"0":  DisplayComments,
		"1":  DisplayTokens,
		"2":  DisplaySemicolons,
		"3":  DisplayIdentifiers,
		"4":  DisplayKeywords,
		"5":  DisplayOperators,
		"6":  DisplayIntegers,
		"7":  DisplayFloats,
		"8":  DisplayImaginary,
		"9":  DisplayRunes,
		"10": DisplayStrings,
	}

	for {
		fmt.Fprintln(stdout, "\n词法元素学习菜单")
		fmt.Fprintln(stdout, "---------------------------------")
		fmt.Fprintln(stdout, "请选择要学习的主题:")

		// 显示所有选项
		for i := 0; i <= 10; i++ {
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

		// 验证输入范围（0-10）并执行对应的主题显示函数
		if action, ok := topicActions[choice]; ok {
			// 执行对应的主题显示函数
			action()
			// 执行完毕后，循环会继续显示菜单，允许用户选择其他主题或返回
		} else {
			// 处理无效输入：显示错误消息并重新提示
			fmt.Fprintln(stdout, "无效的选择，请重试。")
		}
	}
}
