package lexical_elements

import "fmt"

// DisplayIntegers 展示 Go 语言中整数字面量的各种表示形式
// 包括十进制、二进制、八进制、十六进制等
func DisplayIntegers() {
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("【词法元素 - 整数字面量 (Integer Literals)】")
	fmt.Println(repeatString("=", 60))

	// 1. 十进制整数（最常用）
	fmt.Println("\n1. 十进制整数 (Decimal)")
	fmt.Println("   十进制是最常用的整数表示形式，由 0-9 的数字组成")

	decimal1 := 42
	decimal2 := 1000000
	decimal3 := 1_000_000 // 使用下划线提高可读性（Go 1.13+）

	fmt.Printf("   42 = %d\n", decimal1)
	fmt.Printf("   1000000 = %d\n", decimal2)
	fmt.Printf("   1_000_000 = %d (使用下划线分隔，更易读)\n", decimal3)

	// 2. 二进制整数
	fmt.Println("\n2. 二进制整数 (Binary)")
	fmt.Println("   以 0b 或 0B 开头，由 0 和 1 组成")

	binary1 := 0b1010      // 二进制 1010 = 十进制 10
	binary2 := 0b11111111  // 二进制 11111111 = 十进制 255
	binary3 := 0b1111_0000 // 使用下划线分隔

	fmt.Printf("   0b1010 = %d (二进制: %b)\n", binary1, binary1)
	fmt.Printf("   0b11111111 = %d (二进制: %b)\n", binary2, binary2)
	fmt.Printf("   0b1111_0000 = %d (二进制: %b)\n", binary3, binary3)

	// 3. 八进制整数
	fmt.Println("\n3. 八进制整数 (Octal)")
	fmt.Println("   以 0o 或 0O 开头（推荐），或单独的 0 开头（旧式）")
	fmt.Println("   由 0-7 的数字组成")

	octal1 := 0o755 // 八进制 755 = 十进制 493（常用于文件权限）
	octal2 := 0o10  // 八进制 10 = 十进制 8
	octal3 := 0644  // 旧式写法（不推荐，但仍有效）

	fmt.Printf("   0o755 = %d (八进制: %o) - 常用于 Unix 文件权限\n", octal1, octal1)
	fmt.Printf("   0o10 = %d (八进制: %o)\n", octal2, octal2)
	fmt.Printf("   0644 = %d (八进制: %o) - 旧式写法\n", octal3, octal3)

	// 4. 十六进制整数
	fmt.Println("\n4. 十六进制整数 (Hexadecimal)")
	fmt.Println("   以 0x 或 0X 开头，由 0-9 和 A-F（不区分大小写）组成")

	hex1 := 0xFF        // 十六进制 FF = 十进制 255
	hex2 := 0x1A2B      // 十六进制 1A2B = 十进制 6699
	hex3 := 0xDEAD_BEEF // 使用下划线分隔（常见于内存地址）

	fmt.Printf("   0xFF = %d (十六进制: %X)\n", hex1, hex1)
	fmt.Printf("   0x1A2B = %d (十六进制: %X)\n", hex2, hex2)
	fmt.Printf("   0xDEAD_BEEF = %d (十六进制: %X)\n", hex3, hex3)

	// 5. 整数类型
	fmt.Println("\n5. 整数类型")
	fmt.Println("   Go 提供了多种整数类型，根据大小和符号分类：")
	fmt.Println("   有符号整数: int8, int16, int32, int64, int")
	fmt.Println("   无符号整数: uint8, uint16, uint32, uint64, uint, uintptr")
	fmt.Println("   别名: byte (uint8), rune (int32)")

	var i8 int8 = 127                   // 8 位有符号整数，范围 -128 到 127
	var u8 uint8 = 255                  // 8 位无符号整数，范围 0 到 255
	var i32 int32 = 2147483647          // 32 位有符号整数
	var i64 int64 = 9223372036854775807 // 64 位有符号整数

	fmt.Printf("   int8:  %d (范围: -128 到 127)\n", i8)
	fmt.Printf("   uint8: %d (范围: 0 到 255)\n", u8)
	fmt.Printf("   int32: %d\n", i32)
	fmt.Printf("   int64: %d\n", i64)

	// 6. 下划线分隔符（提高可读性）
	fmt.Println("\n6. 使用下划线提高可读性 (Go 1.13+)")
	fmt.Println("   可以在数字之间使用下划线 _ 分隔，提高大数字的可读性")

	population := 1_400_000_000 // 14 亿
	hexColor := 0xFF_AA_00      // RGB 颜色值
	binary := 0b1111_0000_1010_1100

	fmt.Printf("   1_400_000_000 = %d (人口数)\n", population)
	fmt.Printf("   0xFF_AA_00 = #%X (RGB 颜色)\n", hexColor)
	fmt.Printf("   0b1111_0000_1010_1100 = %d (二进制: %b)\n", binary, binary)

	// 7. 整数运算示例
	fmt.Println("\n7. 整数运算示例")
	a := 100
	b := 7
	fmt.Printf("   a = %d, b = %d\n", a, b)
	fmt.Printf("   a + b = %d\n", a+b)
	fmt.Printf("   a - b = %d\n", a-b)
	fmt.Printf("   a * b = %d\n", a*b)
	fmt.Printf("   a / b = %d (整数除法，结果向下取整)\n", a/b)
	fmt.Printf("   a %% b = %d (取余数)\n", a%b)

	// 8. 类型转换
	fmt.Println("\n8. 整数类型转换")
	fmt.Println("   Go 要求显式类型转换，不会自动转换")

	var x int32 = 100
	var y int64 = int64(x) // 显式转换
	fmt.Printf("   int32(%d) → int64(%d)\n", x, y)

	var f float64 = 3.14
	var intFromFloat int = int(f) // 浮点数转整数（截断小数部分）
	fmt.Printf("   float64(%.2f) → int(%d) (小数部分被截断)\n", f, intFromFloat)

	// 9. 最佳实践
	fmt.Println("\n9. 最佳实践")
	fmt.Println("   ✓ 默认使用 int 类型（平台相关，32 或 64 位）")
	fmt.Println("   ✓ 需要特定大小时使用 int8/int16/int32/int64")
	fmt.Println("   ✓ 使用下划线分隔大数字，提高可读性")
	fmt.Println("   ✓ 二进制用 0b，八进制用 0o，十六进制用 0x")
	fmt.Println("   ✓ 避免使用旧式八进制表示法（单独的 0 前缀）")
	fmt.Println("   ✗ 注意整数溢出问题")

	fmt.Println("\n" + repeatString("=", 60))
}
