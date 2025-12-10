package lexical_elements

import (
	"fmt"
	"strings"
)

// GetCommentsContent 返回注释相关的学习内容
func GetCommentsContent() string {
	var sb strings.Builder

	sb.WriteString("\n--- Go 语言的注释 ---\n")
	sb.WriteString("注释是代码中非常重要的一部分，用于解释代码的功能、目的和实现方式。\n")
	sb.WriteString("Go 语言支持两种类型的注释：\n")

	// 1. 单行注释
	sb.WriteString("\n1. 单行注释 (Single-line Comments):\n")
	sb.WriteString("// 这是一个单行注释。\n")
	variable := 10
	sb.WriteString(fmt.Sprintf("变量 'variable' 的值是: %d (这行代码后面就有一个单行注释)。\n", variable))

	// 2. 多行注释
	sb.WriteString("\n2. 多行注释 (Multi-line or Block Comments):\n")
	sb.WriteString("/*\n")
	sb.WriteString(" * 这是一个多行注释。\n")
	sb.WriteString(" * 它可以包含很多行的文本。\n")
	sb.WriteString(" */\n")
	sb.WriteString("Go 源码中经常使用多行注释来为包、函数、类型或变量提供文档。\n")
	sb.WriteString("这种文档注释（doc comments）是一种重要的实践，可以使用 go doc 工具来查看。\n")

	return sb.String()
}

// DisplayComments 展示并解释 Go 语言中的注释。
// 注释是程序员为了解释代码而添加的文本，它们会被编译器忽略。
func DisplayComments() {
	fmt.Print(GetCommentsContent())
}
