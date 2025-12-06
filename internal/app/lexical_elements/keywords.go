package lexical_elements

import (
	"fmt"
	"strings"
)

// GetKeywordsContent 返回关键字相关的学习内容
func GetKeywordsContent() string {
	var sb strings.Builder

	sb.WriteString("\n" + repeatString("=", 60) + "\n")
	sb.WriteString("【词法元素 - 关键字 (Keywords)】\n")
	sb.WriteString(repeatString("=", 60) + "\n")

	// 1. Go 语言的 25 个关键字
	sb.WriteString("\n1. Go 语言的 25 个关键字（按类别分组）\n")
	sb.WriteString("\n   声明相关 (5 个):\n")
	sb.WriteString("   var      - 声明变量\n")
	sb.WriteString("   const    - 声明常量\n")
	sb.WriteString("   type     - 声明类型\n")
	sb.WriteString("   func     - 声明函数\n")
	sb.WriteString("   package  - 声明包\n")

	sb.WriteString("\n   控制流 (11 个):\n")
	sb.WriteString("   if, else         - 条件判断\n")
	sb.WriteString("   for              - 循环\n")
	sb.WriteString("   switch, case     - 分支选择\n")
	sb.WriteString("   default          - 默认分支\n")
	sb.WriteString("   break            - 跳出循环或 switch\n")
	sb.WriteString("   continue         - 继续下一次循环\n")
	sb.WriteString("   fallthrough      - 继续执行下一个 case\n")
	sb.WriteString("   goto             - 无条件跳转\n")
	sb.WriteString("   return           - 函数返回\n")

	sb.WriteString("\n   并发相关 (3 个):\n")
	sb.WriteString("   go               - 启动 goroutine\n")
	sb.WriteString("   chan             - 声明通道\n")
	sb.WriteString("   select           - 多路通道选择\n")

	sb.WriteString("\n   其他 (6 个):\n")
	sb.WriteString("   import           - 导入包\n")
	sb.WriteString("   struct           - 定义结构体\n")
	sb.WriteString("   interface        - 定义接口\n")
	sb.WriteString("   map              - 声明映射\n")
	sb.WriteString("   range            - 遍历集合\n")
	sb.WriteString("   defer            - 延迟执行\n")

	// 2. 声明关键字示例
	sb.WriteString("\n2. 声明关键字示例\n")

	// var - 变量声明
	var count int = 10
	sb.WriteString(fmt.Sprintf("   var count int = %d\n", count))

	// const - 常量声明
	const PI = 3.14159
	sb.WriteString(fmt.Sprintf("   const PI = %.5f\n", PI))

	// type - 类型声明
	type UserID int
	var uid UserID = 1001
	sb.WriteString(fmt.Sprintf("   type UserID int; uid = %d\n", uid))

	// 3. 控制流关键字示例
	sb.WriteString("\n3. 控制流关键字示例\n")

	// if-else
	sb.WriteString("\n   if-else 示例:\n")
	score := 85
	if score >= 90 {
		sb.WriteString("   成绩: 优秀\n")
	} else if score >= 60 {
		sb.WriteString("   成绩: 及格\n")
	} else {
		sb.WriteString("   成绩: 不及格\n")
	}

	// for 循环
	sb.WriteString("\n   for 循环示例:\n")
	for i := 1; i <= 3; i++ {
		sb.WriteString(fmt.Sprintf("   循环 %d\n", i))
	}

	// switch-case
	sb.WriteString("\n   switch-case 示例:\n")
	day := 3
	switch day {
	case 1:
		sb.WriteString("   星期一\n")
	case 2:
		sb.WriteString("   星期二\n")
	case 3:
		sb.WriteString("   星期三\n")
	default:
		sb.WriteString("   其他\n")
	}

	// break 和 continue
	sb.WriteString("\n   break 和 continue 示例:\n")
	for i := 1; i <= 5; i++ {
		if i == 3 {
			sb.WriteString("   跳过 3 (continue)\n")
			continue
		}
		if i == 5 {
			sb.WriteString("   遇到 5，退出循环 (break)\n")
			break
		}
		sb.WriteString(fmt.Sprintf("   处理: %d\n", i))
	}

	// fallthrough
	sb.WriteString("\n   fallthrough 示例:\n")
	num := 1
	switch num {
	case 1:
		sb.WriteString("   case 1 执行\n")
		fallthrough // 继续执行下一个 case
	case 2:
		sb.WriteString("   case 2 也执行（因为 fallthrough）\n")
	}

	// 4. 并发关键字示例
	sb.WriteString("\n4. 并发关键字示例\n")

	// chan - 通道
	ch := make(chan int, 1)
	sb.WriteString("   创建通道: ch := make(chan int, 1)\n")

	// go - 启动 goroutine
	go func() {
		ch <- 42 // 发送数据到通道
	}()

	// select - 多路选择
	select {
	case value := <-ch:
		sb.WriteString(fmt.Sprintf("   从通道接收到: %d\n", value))
	default:
		sb.WriteString("   通道为空\n")
	}

	// 5. defer 关键字
	sb.WriteString("\n5. defer 关键字示例\n")
	sb.WriteString("   defer 用于延迟函数调用，常用于资源清理\n")
	deferDemo(&sb)

	// 6. range 关键字
	sb.WriteString("\n6. range 关键字示例\n")
	sb.WriteString("   range 用于遍历数组、切片、映射、通道等\n")
	numbers := []int{10, 20, 30}
	for index, value := range numbers {
		sb.WriteString(fmt.Sprintf("   索引 %d: 值 %d\n", index, value))
	}

	// 7. struct 和 interface
	sb.WriteString("\n7. struct 和 interface 示例\n")

	type Animal struct {
		Name string
	}
	dog := Animal{Name: "旺财"}
	sb.WriteString(fmt.Sprintf("   struct 示例: %s\n", dog.Name))

	// 8. map 关键字
	sb.WriteString("\n8. map 关键字示例\n")
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
	}
	sb.WriteString(fmt.Sprintf("   map 示例: Alice 的年龄是 %d\n", ages["Alice"]))

	// 9. 注意事项
	sb.WriteString("\n9. 关键字使用注意事项\n")
	sb.WriteString("   ✗ 不能将关键字用作变量名、函数名或类型名\n")
	sb.WriteString("   ✗ 关键字区分大小写（如 For 不是关键字，但 for 是）\n")
	sb.WriteString("   ✓ 使用 IDE 的语法高亮可以帮助识别关键字\n")

	sb.WriteString("\n" + repeatString("=", 60) + "\n")

	return sb.String()
}

// DisplayKeywords 展示 Go 语言中的 25 个关键字及其用法
// 关键字是保留字，不能用作标识符
func DisplayKeywords() {
	fmt.Print(GetKeywordsContent())
}

// deferDemo 演示 defer 的使用
func deferDemo(sb *strings.Builder) {
	defer func() { sb.WriteString("   第三步: defer 语句最后执行\n") }()
	sb.WriteString("   第一步: 正常语句\n")
	defer func() { sb.WriteString("   第二步: 多个 defer 按 LIFO 顺序执行\n") }()
}
