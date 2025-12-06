package lexical_elements

import "fmt"

// DisplayComments 展示并解释 Go 语言中的注释。
// 注释是程序员为了解释代码而添加的文本，它们会被编译器忽略。
func DisplayComments() {
	fmt.Println("\n--- Go 语言的注释 ---")
	fmt.Println("注释是代码中非常重要的一部分，用于解释代码的功能、目的和实现方式。")
	fmt.Println("Go 语言支持两种类型的注释：")

	// 1. 单行注释
	// 单行注释以 `//` 开始，直到行尾结束。
	// 它们通常用于解释单行代码或简短的代码块。
	fmt.Println("\n1. 单行注释 (Single-line Comments):")
	fmt.Println("// 这是一个单行注释。")
	variable := 10 // 也可以在代码行的末尾添加注释来解释这一行的作用。
	fmt.Printf("变量 'variable' 的值是: %d (这行代码后面就有一个单行注释)。\n", variable)

	// 2. 多行注释
	/*
	   多行注释（也称为块注释）以 slash-star 开始，以 star-slash 结束。
	   它们可以跨越多行，通常用于提供更详细的文档、禁用一大块代码或者在函数/包级别进行说明。
	*/
	fmt.Println("\n2. 多行注释 (Multi-line or Block Comments):")
	fmt.Println("/*")
	fmt.Println(" * 这是一个多行注释。")
	fmt.Println(" * 它可以包含很多行的文本。")
	fmt.Println(" */")
	fmt.Println("Go 源码中经常使用多行注释来为包、函数、类型或变量提供文档。")
	fmt.Println("这种文档注释（doc comments）是一种重要的实践，可以使用 go doc 工具来查看。")
}
