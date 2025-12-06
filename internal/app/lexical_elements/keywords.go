package lexical_elements

import "fmt"

// DisplayKeywords 展示 Go 语言中的 25 个关键字及其用法
// 关键字是保留字，不能用作标识符
func DisplayKeywords() {
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("【词法元素 - 关键字 (Keywords)】")
	fmt.Println(repeatString("=", 60))

	// 1. Go 语言的 25 个关键字
	fmt.Println("\n1. Go 语言的 25 个关键字（按类别分组）")
	fmt.Println("\n   声明相关 (5 个):")
	fmt.Println("   var      - 声明变量")
	fmt.Println("   const    - 声明常量")
	fmt.Println("   type     - 声明类型")
	fmt.Println("   func     - 声明函数")
	fmt.Println("   package  - 声明包")

	fmt.Println("\n   控制流 (11 个):")
	fmt.Println("   if, else         - 条件判断")
	fmt.Println("   for              - 循环")
	fmt.Println("   switch, case     - 分支选择")
	fmt.Println("   default          - 默认分支")
	fmt.Println("   break            - 跳出循环或 switch")
	fmt.Println("   continue         - 继续下一次循环")
	fmt.Println("   fallthrough      - 继续执行下一个 case")
	fmt.Println("   goto             - 无条件跳转")
	fmt.Println("   return           - 函数返回")

	fmt.Println("\n   并发相关 (3 个):")
	fmt.Println("   go               - 启动 goroutine")
	fmt.Println("   chan             - 声明通道")
	fmt.Println("   select           - 多路通道选择")

	fmt.Println("\n   其他 (6 个):")
	fmt.Println("   import           - 导入包")
	fmt.Println("   struct           - 定义结构体")
	fmt.Println("   interface        - 定义接口")
	fmt.Println("   map              - 声明映射")
	fmt.Println("   range            - 遍历集合")
	fmt.Println("   defer            - 延迟执行")

	// 2. 声明关键字示例
	fmt.Println("\n2. 声明关键字示例")

	// var - 变量声明
	var count int = 10
	fmt.Printf("   var count int = %d\n", count)

	// const - 常量声明
	const PI = 3.14159
	fmt.Printf("   const PI = %.5f\n", PI)

	// type - 类型声明
	type UserID int
	var uid UserID = 1001
	fmt.Printf("   type UserID int; uid = %d\n", uid)

	// 3. 控制流关键字示例
	fmt.Println("\n3. 控制流关键字示例")

	// if-else
	fmt.Println("\n   if-else 示例:")
	score := 85
	if score >= 90 {
		fmt.Println("   成绩: 优秀")
	} else if score >= 60 {
		fmt.Println("   成绩: 及格")
	} else {
		fmt.Println("   成绩: 不及格")
	}

	// for 循环
	fmt.Println("\n   for 循环示例:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("   循环 %d\n", i)
	}

	// switch-case
	fmt.Println("\n   switch-case 示例:")
	day := 3
	switch day {
	case 1:
		fmt.Println("   星期一")
	case 2:
		fmt.Println("   星期二")
	case 3:
		fmt.Println("   星期三")
	default:
		fmt.Println("   其他")
	}

	// break 和 continue
	fmt.Println("\n   break 和 continue 示例:")
	for i := 1; i <= 5; i++ {
		if i == 3 {
			fmt.Println("   跳过 3 (continue)")
			continue
		}
		if i == 5 {
			fmt.Println("   遇到 5，退出循环 (break)")
			break
		}
		fmt.Printf("   处理: %d\n", i)
	}

	// fallthrough
	fmt.Println("\n   fallthrough 示例:")
	num := 1
	switch num {
	case 1:
		fmt.Println("   case 1 执行")
		fallthrough // 继续执行下一个 case
	case 2:
		fmt.Println("   case 2 也执行（因为 fallthrough）")
	}

	// 4. 并发关键字示例
	fmt.Println("\n4. 并发关键字示例")

	// chan - 通道
	ch := make(chan int, 1)
	fmt.Println("   创建通道: ch := make(chan int, 1)")

	// go - 启动 goroutine
	go func() {
		ch <- 42 // 发送数据到通道
	}()

	// select - 多路选择
	select {
	case value := <-ch:
		fmt.Printf("   从通道接收到: %d\n", value)
	default:
		fmt.Println("   通道为空")
	}

	// 5. defer 关键字
	fmt.Println("\n5. defer 关键字示例")
	fmt.Println("   defer 用于延迟函数调用，常用于资源清理")
	deferDemo()

	// 6. range 关键字
	fmt.Println("\n6. range 关键字示例")
	fmt.Println("   range 用于遍历数组、切片、映射、通道等")
	numbers := []int{10, 20, 30}
	for index, value := range numbers {
		fmt.Printf("   索引 %d: 值 %d\n", index, value)
	}

	// 7. struct 和 interface
	fmt.Println("\n7. struct 和 interface 示例")

	type Animal struct {
		Name string
	}
	dog := Animal{Name: "旺财"}
	fmt.Printf("   struct 示例: %s\n", dog.Name)

	// 8. map 关键字
	fmt.Println("\n8. map 关键字示例")
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
	}
	fmt.Printf("   map 示例: Alice 的年龄是 %d\n", ages["Alice"])

	// 9. 注意事项
	fmt.Println("\n9. 关键字使用注意事项")
	fmt.Println("   ✗ 不能将关键字用作变量名、函数名或类型名")
	fmt.Println("   ✗ 关键字区分大小写（如 For 不是关键字，但 for 是）")
	fmt.Println("   ✓ 使用 IDE 的语法高亮可以帮助识别关键字")

	fmt.Println("\n" + repeatString("=", 60))
}

// deferDemo 演示 defer 的使用
func deferDemo() {
	defer fmt.Println("   第三步: defer 语句最后执行")
	fmt.Println("   第一步: 正常语句")
	defer fmt.Println("   第二步: 多个 defer 按 LIFO 顺序执行")
}
