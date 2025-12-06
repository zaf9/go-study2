package lexical_elements

import "fmt"

// DisplayOperators 展示 Go 语言中的运算符和标点符号
// 包括算术、比较、逻辑、位运算、赋值等运算符
func DisplayOperators() {
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("【词法元素 - 运算符和标点 (Operators and Punctuation)】")
	fmt.Println(repeatString("=", 60))

	// 1. 算术运算符
	fmt.Println("\n1. 算术运算符 (Arithmetic Operators)")
	a, b := 10, 3
	fmt.Printf("   a = %d, b = %d\n", a, b)
	fmt.Printf("   a + b = %d  (加法)\n", a+b)
	fmt.Printf("   a - b = %d  (减法)\n", a-b)
	fmt.Printf("   a * b = %d  (乘法)\n", a*b)
	fmt.Printf("   a / b = %d  (除法，整数除法)\n", a/b)
	fmt.Printf("   a %% b = %d  (取模，余数)\n", a%b)

	// 自增自减
	fmt.Println("\n   自增自减运算符:")
	x := 5
	fmt.Printf("   x = %d\n", x)
	x++
	fmt.Printf("   x++ 后: x = %d\n", x)
	x--
	fmt.Printf("   x-- 后: x = %d\n", x)
	fmt.Println("   注意: Go 中 ++ 和 -- 是语句，不是表达式")

	// 2. 比较运算符
	fmt.Println("\n2. 比较运算符 (Comparison Operators)")
	m, n := 10, 20
	fmt.Printf("   m = %d, n = %d\n", m, n)
	fmt.Printf("   m == n: %t  (等于)\n", m == n)
	fmt.Printf("   m != n: %t  (不等于)\n", m != n)
	fmt.Printf("   m < n:  %t  (小于)\n", m < n)
	fmt.Printf("   m <= n: %t  (小于等于)\n", m <= n)
	fmt.Printf("   m > n:  %t  (大于)\n", m > n)
	fmt.Printf("   m >= n: %t  (大于等于)\n", m >= n)

	// 3. 逻辑运算符
	fmt.Println("\n3. 逻辑运算符 (Logical Operators)")
	p, q := true, false
	fmt.Printf("   p = %t, q = %t\n", p, q)
	fmt.Printf("   p && q: %t  (逻辑与 AND)\n", p && q)
	fmt.Printf("   p || q: %t  (逻辑或 OR)\n", p || q)
	fmt.Printf("   !p:     %t  (逻辑非 NOT)\n", !p)

	// 4. 位运算符
	fmt.Println("\n4. 位运算符 (Bitwise Operators)")
	c, d := 12, 25 // 二进制: 1100 和 11001
	fmt.Printf("   c = %d (二进制: %b)\n", c, c)
	fmt.Printf("   d = %d (二进制: %b)\n", d, d)
	fmt.Printf("   c & d  = %d  (按位与 AND)\n", c&d)
	fmt.Printf("   c | d  = %d  (按位或 OR)\n", c|d)
	fmt.Printf("   c ^ d  = %d  (按位异或 XOR)\n", c^d)
	fmt.Printf("   c &^ d = %d  (按位清除 AND NOT)\n", c&^d)
	fmt.Printf("   c << 1 = %d  (左移)\n", c<<1)
	fmt.Printf("   c >> 1 = %d  (右移)\n", c>>1)

	// 5. 赋值运算符
	fmt.Println("\n5. 赋值运算符 (Assignment Operators)")
	val := 10
	fmt.Printf("   val := %d  (简短声明并赋值)\n", val)
	val = 20
	fmt.Printf("   val = %d   (赋值)\n", val)
	val += 5
	fmt.Printf("   val += 5   → val = %d  (加法赋值)\n", val)
	val -= 3
	fmt.Printf("   val -= 3   → val = %d  (减法赋值)\n", val)
	val *= 2
	fmt.Printf("   val *= 2   → val = %d  (乘法赋值)\n", val)
	val /= 4
	fmt.Printf("   val /= 4   → val = %d  (除法赋值)\n", val)
	val %= 3
	fmt.Printf("   val %%= 3   → val = %d  (取模赋值)\n", val)

	// 6. 指针运算符
	fmt.Println("\n6. 指针运算符 (Pointer Operators)")
	num := 42
	ptr := &num // & 取地址
	fmt.Printf("   num = %d\n", num)
	fmt.Printf("   ptr = &num  → 地址: %p\n", ptr)
	fmt.Printf("   *ptr = %d   (解引用，获取指针指向的值)\n", *ptr)
	*ptr = 100 // 通过指针修改值
	fmt.Printf("   *ptr = 100  → num = %d\n", num)

	// 7. 通道运算符
	fmt.Println("\n7. 通道运算符 (Channel Operator)")
	ch := make(chan int, 1)
	ch <- 99 // <- 发送数据到通道
	fmt.Println("   ch <- 99  (发送数据到通道)")
	received := <-ch // <- 从通道接收数据
	fmt.Printf("   received := <-ch  → %d\n", received)

	// 8. 其他重要标点符号
	fmt.Println("\n8. 其他重要标点符号")
	fmt.Println("   ()  - 函数调用、分组表达式、类型转换")
	fmt.Println("   []  - 数组/切片索引、类型声明")
	fmt.Println("   {}  - 代码块、复合字面量")
	fmt.Println("   .   - 选择器（访问字段或方法）")
	fmt.Println("   ,   - 分隔符（参数、元素）")
	fmt.Println("   ;   - 语句终止符（通常自动插入）")
	fmt.Println("   :   - 标签、短变量声明、case 语句")
	fmt.Println("   ...  - 可变参数、数组字面量")

	// 9. 运算符优先级
	fmt.Println("\n9. 运算符优先级（从高到低）")
	fmt.Println("   优先级 5: *  /  %  <<  >>  &  &^")
	fmt.Println("   优先级 4: +  -  |  ^")
	fmt.Println("   优先级 3: ==  !=  <  <=  >  >=")
	fmt.Println("   优先级 2: &&")
	fmt.Println("   优先级 1: ||")

	// 优先级示例
	fmt.Println("\n   优先级示例:")
	result := 2 + 3*4
	fmt.Printf("   2 + 3*4 = %d  (* 优先级高于 +)\n", result)
	result = (2 + 3) * 4
	fmt.Printf("   (2 + 3)*4 = %d  (括号改变优先级)\n", result)

	// 10. 特殊运算符示例
	fmt.Println("\n10. 特殊运算符和用法")

	// 类型断言
	var i interface{} = "hello"
	s, ok := i.(string)
	fmt.Printf("   类型断言: i.(string) → %s, ok=%t\n", s, ok)

	// 可变参数
	sum := sumAll(1, 2, 3, 4, 5)
	fmt.Printf("   可变参数: sumAll(1,2,3,4,5) = %d\n", sum)

	fmt.Println("\n" + repeatString("=", 60))
}

// sumAll 演示可变参数 (...)
func sumAll(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
