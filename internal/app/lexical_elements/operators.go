package lexical_elements

import (
	"fmt"
	"strings"
)

// GetOperatorsContent 返回运算符相关的学习内容
func GetOperatorsContent() string {
	var sb strings.Builder

	sb.WriteString("\n" + repeatString("=", 60) + "\n")
	sb.WriteString("【词法元素 - 运算符和标点 (Operators and Punctuation)】\n")
	sb.WriteString(repeatString("=", 60) + "\n")

	// 1. 算术运算符
	sb.WriteString("\n1. 算术运算符 (Arithmetic Operators)\n")
	a, b := 10, 3
	sb.WriteString(fmt.Sprintf("   a = %d, b = %d\n", a, b))
	sb.WriteString(fmt.Sprintf("   a + b = %d  (加法)\n", a+b))
	sb.WriteString(fmt.Sprintf("   a - b = %d  (减法)\n", a-b))
	sb.WriteString(fmt.Sprintf("   a * b = %d  (乘法)\n", a*b))
	sb.WriteString(fmt.Sprintf("   a / b = %d  (除法，整数除法)\n", a/b))
	sb.WriteString(fmt.Sprintf("   a %% b = %d  (取模，余数)\n", a%b))

	// 自增自减
	sb.WriteString("\n   自增自减运算符:\n")
	x := 5
	sb.WriteString(fmt.Sprintf("   x = %d\n", x))
	x++
	sb.WriteString(fmt.Sprintf("   x++ 后: x = %d\n", x))
	x--
	sb.WriteString(fmt.Sprintf("   x-- 后: x = %d\n", x))
	sb.WriteString("   注意: Go 中 ++ 和 -- 是语句，不是表达式\n")

	// 2. 比较运算符
	sb.WriteString("\n2. 比较运算符 (Comparison Operators)\n")
	m, n := 10, 20
	sb.WriteString(fmt.Sprintf("   m = %d, n = %d\n", m, n))
	sb.WriteString(fmt.Sprintf("   m == n: %t  (等于)\n", m == n))
	sb.WriteString(fmt.Sprintf("   m != n: %t  (不等于)\n", m != n))
	sb.WriteString(fmt.Sprintf("   m < n:  %t  (小于)\n", m < n))
	sb.WriteString(fmt.Sprintf("   m <= n: %t  (小于等于)\n", m <= n))
	sb.WriteString(fmt.Sprintf("   m > n:  %t  (大于)\n", m > n))
	sb.WriteString(fmt.Sprintf("   m >= n: %t  (大于等于)\n", m >= n))

	// 3. 逻辑运算符
	sb.WriteString("\n3. 逻辑运算符 (Logical Operators)\n")
	p, q := true, false
	sb.WriteString(fmt.Sprintf("   p = %t, q = %t\n", p, q))
	sb.WriteString(fmt.Sprintf("   p && q: %t  (逻辑与 AND)\n", p && q))
	sb.WriteString(fmt.Sprintf("   p || q: %t  (逻辑或 OR)\n", p || q))
	sb.WriteString(fmt.Sprintf("   !p:     %t  (逻辑非 NOT)\n", !p))

	// 4. 位运算符
	sb.WriteString("\n4. 位运算符 (Bitwise Operators)\n")
	c, d := 12, 25 // 二进制: 1100 和 11001
	sb.WriteString(fmt.Sprintf("   c = %d (二进制: %b)\n", c, c))
	sb.WriteString(fmt.Sprintf("   d = %d (二进制: %b)\n", d, d))
	sb.WriteString(fmt.Sprintf("   c & d  = %d  (按位与 AND)\n", c&d))
	sb.WriteString(fmt.Sprintf("   c | d  = %d  (按位或 OR)\n", c|d))
	sb.WriteString(fmt.Sprintf("   c ^ d  = %d  (按位异或 XOR)\n", c^d))
	sb.WriteString(fmt.Sprintf("   c &^ d = %d  (按位清除 AND NOT)\n", c&^d))
	sb.WriteString(fmt.Sprintf("   c << 1 = %d  (左移)\n", c<<1))
	sb.WriteString(fmt.Sprintf("   c >> 1 = %d  (右移)\n", c>>1))

	// 5. 赋值运算符
	sb.WriteString("\n5. 赋值运算符 (Assignment Operators)\n")
	val := 10
	sb.WriteString(fmt.Sprintf("   val := %d  (简短声明并赋值)\n", val))
	val = 20
	sb.WriteString(fmt.Sprintf("   val = %d   (赋值)\n", val))
	val += 5
	sb.WriteString(fmt.Sprintf("   val += 5   → val = %d  (加法赋值)\n", val))
	val -= 3
	sb.WriteString(fmt.Sprintf("   val -= 3   → val = %d  (减法赋值)\n", val))
	val *= 2
	sb.WriteString(fmt.Sprintf("   val *= 2   → val = %d  (乘法赋值)\n", val))
	val /= 4
	sb.WriteString(fmt.Sprintf("   val /= 4   → val = %d  (除法赋值)\n", val))
	val %= 3
	sb.WriteString(fmt.Sprintf("   val %%= 3   → val = %d  (取模赋值)\n", val))

	// 6. 指针运算符
	sb.WriteString("\n6. 指针运算符 (Pointer Operators)\n")
	num := 42
	ptr := &num // & 取地址
	sb.WriteString(fmt.Sprintf("   num = %d\n", num))
	sb.WriteString(fmt.Sprintf("   ptr = &num  → 地址: %p\n", ptr))
	sb.WriteString(fmt.Sprintf("   *ptr = %d   (解引用，获取指针指向的值)\n", *ptr))
	*ptr = 100 // 通过指针修改值
	sb.WriteString(fmt.Sprintf("   *ptr = 100  → num = %d\n", num))

	// 7. 通道运算符
	sb.WriteString("\n7. 通道运算符 (Channel Operator)\n")
	ch := make(chan int, 1)
	ch <- 99 // <- 发送数据到通道
	sb.WriteString("   ch <- 99  (发送数据到通道)\n")
	received := <-ch // <- 从通道接收数据
	sb.WriteString(fmt.Sprintf("   received := <-ch  → %d\n", received))

	// 8. 其他重要标点符号
	sb.WriteString("\n8. 其他重要标点符号\n")
	sb.WriteString("   ()  - 函数调用、分组表达式、类型转换\n")
	sb.WriteString("   []  - 数组/切片索引、类型声明\n")
	sb.WriteString("   {}  - 代码块、复合字面量\n")
	sb.WriteString("   .   - 选择器（访问字段或方法）\n")
	sb.WriteString("   ,   - 分隔符（参数、元素）\n")
	sb.WriteString("   ;   - 语句终止符（通常自动插入）\n")
	sb.WriteString("   :   - 标签、短变量声明、case 语句\n")
	sb.WriteString("   ...  - 可变参数、数组字面量\n")

	// 9. 运算符优先级
	sb.WriteString("\n9. 运算符优先级（从高到低）\n")
	sb.WriteString("   优先级 5: *  /  %  <<  >>  &  &^\n")
	sb.WriteString("   优先级 4: +  -  |  ^\n")
	sb.WriteString("   优先级 3: ==  !=  <  <=  >  >=\n")
	sb.WriteString("   优先级 2: &&\n")
	sb.WriteString("   优先级 1: ||\n")

	// 优先级示例
	sb.WriteString("\n   优先级示例:\n")
	result := 2 + 3*4
	sb.WriteString(fmt.Sprintf("   2 + 3*4 = %d  (* 优先级高于 +)\n", result))
	result = (2 + 3) * 4
	sb.WriteString(fmt.Sprintf("   (2 + 3)*4 = %d  (括号改变优先级)\n", result))

	// 10. 特殊运算符示例
	sb.WriteString("\n10. 特殊运算符和用法\n")

	// 类型断言
	var i interface{} = "hello"
	s, ok := i.(string)
	sb.WriteString(fmt.Sprintf("   类型断言: i.(string) → %s, ok=%t\n", s, ok))

	// 可变参数
	sum := sumAll(1, 2, 3, 4, 5)
	sb.WriteString(fmt.Sprintf("   可变参数: sumAll(1,2,3,4,5) = %d\n", sum))

	sb.WriteString("\n" + repeatString("=", 60) + "\n")

	return sb.String()
}

// DisplayOperators 展示 Go 语言中的运算符和标点符号
// 包括算术、比较、逻辑、位运算、赋值等运算符
func DisplayOperators() {
	fmt.Print(GetOperatorsContent())
}

// sumAll 演示可变参数 (...)
func sumAll(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
