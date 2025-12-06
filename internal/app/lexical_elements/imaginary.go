package lexical_elements

import (
	"fmt"
	"math"
	"math/cmplx"
)

// DisplayImaginary 展示 Go 语言中虚数字面量和复数的使用
// Go 原生支持复数运算，虚数字面量以 i 结尾
func DisplayImaginary() {
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("【词法元素 - 虚数字面量 (Imaginary Literals)】")
	fmt.Println(repeatString("=", 60))

	// 1. 虚数字面量的定义
	fmt.Println("\n1. 虚数字面量的定义")
	fmt.Println("   虚数字面量是浮点数或整数后跟小写字母 i")
	fmt.Println("   表示数学中的虚数单位 i（i² = -1）")

	imag1 := 3i    // 整数虚数
	imag2 := 2.5i  // 浮点数虚数
	imag3 := 1e-3i // 科学计数法虚数

	fmt.Printf("   3i = %v\n", imag1)
	fmt.Printf("   2.5i = %v\n", imag2)
	fmt.Printf("   1e-3i = %v\n", imag3)

	// 2. 复数的构造
	fmt.Println("\n2. 复数的构造")
	fmt.Println("   复数 = 实部 + 虚部")
	fmt.Println("   Go 提供两种复数类型: complex64 和 complex128")

	// 直接构造复数
	c1 := 3 + 4i // 复数字面量
	c2 := 1.5 + 2.5i
	c3 := complex(5, 6) // 使用 complex() 函数

	fmt.Printf("   3 + 4i = %v\n", c1)
	fmt.Printf("   1.5 + 2.5i = %v\n", c2)
	fmt.Printf("   complex(5, 6) = %v\n", c3)

	// 3. 复数类型
	fmt.Println("\n3. 复数类型")
	fmt.Println("   complex64  - 实部和虚部都是 float32")
	fmt.Println("   complex128 - 实部和虚部都是 float64（默认）")

	var c64 complex64 = 1 + 2i
	var c128 complex128 = 3 + 4i

	fmt.Printf("   complex64:  %v\n", c64)
	fmt.Printf("   complex128: %v\n", c128)

	// 4. 提取实部和虚部
	fmt.Println("\n4. 提取实部和虚部")
	fmt.Println("   使用 real() 和 imag() 函数")

	z := 6 + 8i
	realPart := real(z)
	imagPart := imag(z)

	fmt.Printf("   z = %v\n", z)
	fmt.Printf("   real(z) = %.0f (实部)\n", realPart)
	fmt.Printf("   imag(z) = %.0f (虚部)\n", imagPart)

	// 5. 复数运算
	fmt.Println("\n5. 复数运算")

	a := 3 + 4i
	b := 1 + 2i

	fmt.Printf("   a = %v, b = %v\n", a, b)
	fmt.Printf("   a + b = %v (加法)\n", a+b)
	fmt.Printf("   a - b = %v (减法)\n", a-b)
	fmt.Printf("   a * b = %v (乘法)\n", a*b)
	fmt.Printf("   a / b = %v (除法)\n", a/b)

	// 6. 复数的模（绝对值）
	fmt.Println("\n6. 复数的模（绝对值）")
	fmt.Println("   使用 math/cmplx 包的 Abs() 函数")

	c := 3 + 4i
	magnitude := cmplx.Abs(c)
	fmt.Printf("   c = %v\n", c)
	fmt.Printf("   |c| = %.1f (模: √(3² + 4²) = 5)\n", magnitude)

	// 7. 复数的其他操作
	fmt.Println("\n7. 复数的其他操作（math/cmplx 包）")

	z1 := 1 + 1i

	fmt.Printf("   z = %v\n", z1)
	fmt.Printf("   共轭: %v (cmplx.Conj)\n", cmplx.Conj(z1))
	fmt.Printf("   相位: %.4f 弧度 (cmplx.Phase)\n", cmplx.Phase(z1))
	fmt.Printf("   指数: %v (cmplx.Exp)\n", cmplx.Exp(z1))
	fmt.Printf("   平方根: %v (cmplx.Sqrt)\n", cmplx.Sqrt(z1))

	// 8. 实际应用示例
	fmt.Println("\n8. 实际应用示例：欧拉公式")
	fmt.Println("   e^(iπ) + 1 = 0 (欧拉恒等式)")

	// e^(iπ) = cos(π) + i*sin(π) = -1 + 0i
	eulerResult := cmplx.Exp(complex(0, math.Pi))
	fmt.Printf("   e^(iπ) = %v\n", eulerResult)
	fmt.Printf("   e^(iπ) + 1 ≈ %v (接近 0)\n", eulerResult+1)

	// 9. 纯虚数和纯实数
	fmt.Println("\n9. 纯虚数和纯实数")

	pureImag := 0 + 5i // 纯虚数
	pureReal := 7 + 0i // 纯实数（复数形式）

	fmt.Printf("   纯虚数: %v (实部为 0)\n", pureImag)
	fmt.Printf("   纯实数: %v (虚部为 0)\n", pureReal)

	// 10. 最佳实践
	fmt.Println("\n10. 最佳实践")
	fmt.Println("   ✓ 虚数字面量必须使用小写 i（不能用大写 I）")
	fmt.Println("   ✓ 默认使用 complex128（更高精度）")
	fmt.Println("   ✓ 使用 complex() 函数从实部和虚部构造复数")
	fmt.Println("   ✓ 使用 real() 和 imag() 提取实部和虚部")
	fmt.Println("   ✓ 导入 math/cmplx 包进行高级复数运算")
	fmt.Println("   ✓ 复数常用于信号处理、量子计算、电路分析等领域")

	fmt.Println("\n" + repeatString("=", 60))
}
