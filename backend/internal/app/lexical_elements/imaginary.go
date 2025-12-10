package lexical_elements

import (
	"fmt"
	"math"
	"math/cmplx"
	"strings"
)

// GetImaginaryContent 返回虚数字面量相关的学习内容
func GetImaginaryContent() string {
	var sb strings.Builder

	sb.WriteString("\n" + repeatString("=", 60) + "\n")
	sb.WriteString("【词法元素 - 虚数字面量 (Imaginary Literals)】\n")
	sb.WriteString(repeatString("=", 60) + "\n")

	// 1. 虚数字面量的定义
	sb.WriteString("\n1. 虚数字面量的定义\n")
	sb.WriteString("   虚数字面量是浮点数或整数后跟小写字母 i\n")
	sb.WriteString("   表示数学中的虚数单位 i（i² = -1）\n")

	imag1 := 3i    // 整数虚数
	imag2 := 2.5i  // 浮点数虚数
	imag3 := 1e-3i // 科学计数法虚数

	sb.WriteString(fmt.Sprintf("   3i = %v\n", imag1))
	sb.WriteString(fmt.Sprintf("   2.5i = %v\n", imag2))
	sb.WriteString(fmt.Sprintf("   1e-3i = %v\n", imag3))

	// 2. 复数的构造
	sb.WriteString("\n2. 复数的构造\n")
	sb.WriteString("   复数 = 实部 + 虚部\n")
	sb.WriteString("   Go 提供两种复数类型: complex64 和 complex128\n")

	// 直接构造复数
	c1 := 3 + 4i // 复数字面量
	c2 := 1.5 + 2.5i
	c3 := complex(5, 6) // 使用 complex() 函数

	sb.WriteString(fmt.Sprintf("   3 + 4i = %v\n", c1))
	sb.WriteString(fmt.Sprintf("   1.5 + 2.5i = %v\n", c2))
	sb.WriteString(fmt.Sprintf("   complex(5, 6) = %v\n", c3))

	// 3. 复数类型
	sb.WriteString("\n3. 复数类型\n")
	sb.WriteString("   complex64  - 实部和虚部都是 float32\n")
	sb.WriteString("   complex128 - 实部和虚部都是 float64（默认）\n")

	var c64 complex64 = 1 + 2i
	var c128 complex128 = 3 + 4i

	sb.WriteString(fmt.Sprintf("   complex64:  %v\n", c64))
	sb.WriteString(fmt.Sprintf("   complex128: %v\n", c128))

	// 4. 提取实部和虚部
	sb.WriteString("\n4. 提取实部和虚部\n")
	sb.WriteString("   使用 real() 和 imag() 函数\n")

	z := 6 + 8i
	realPart := real(z)
	imagPart := imag(z)

	sb.WriteString(fmt.Sprintf("   z = %v\n", z))
	sb.WriteString(fmt.Sprintf("   real(z) = %.0f (实部)\n", realPart))
	sb.WriteString(fmt.Sprintf("   imag(z) = %.0f (虚部)\n", imagPart))

	// 5. 复数运算
	sb.WriteString("\n5. 复数运算\n")

	a := 3 + 4i
	b := 1 + 2i

	sb.WriteString(fmt.Sprintf("   a = %v, b = %v\n", a, b))
	sb.WriteString(fmt.Sprintf("   a + b = %v (加法)\n", a+b))
	sb.WriteString(fmt.Sprintf("   a - b = %v (减法)\n", a-b))
	sb.WriteString(fmt.Sprintf("   a * b = %v (乘法)\n", a*b))
	sb.WriteString(fmt.Sprintf("   a / b = %v (除法)\n", a/b))

	// 6. 复数的模（绝对值）
	sb.WriteString("\n6. 复数的模（绝对值）\n")
	sb.WriteString("   使用 math/cmplx 包的 Abs() 函数\n")

	c := 3 + 4i
	magnitude := cmplx.Abs(c)
	sb.WriteString(fmt.Sprintf("   c = %v\n", c))
	sb.WriteString(fmt.Sprintf("   |c| = %.1f (模: √(3² + 4²) = 5)\n", magnitude))

	// 7. 复数的其他操作
	sb.WriteString("\n7. 复数的其他操作（math/cmplx 包）\n")

	z1 := 1 + 1i

	sb.WriteString(fmt.Sprintf("   z = %v\n", z1))
	sb.WriteString(fmt.Sprintf("   共轭: %v (cmplx.Conj)\n", cmplx.Conj(z1)))
	sb.WriteString(fmt.Sprintf("   相位: %.4f 弧度 (cmplx.Phase)\n", cmplx.Phase(z1)))
	sb.WriteString(fmt.Sprintf("   指数: %v (cmplx.Exp)\n", cmplx.Exp(z1)))
	sb.WriteString(fmt.Sprintf("   平方根: %v (cmplx.Sqrt)\n", cmplx.Sqrt(z1)))

	// 8. 实际应用示例
	sb.WriteString("\n8. 实际应用示例：欧拉公式\n")
	sb.WriteString("   e^(iπ) + 1 = 0 (欧拉恒等式)\n")

	// e^(iπ) = cos(π) + i*sin(π) = -1 + 0i
	eulerResult := cmplx.Exp(complex(0, math.Pi))
	sb.WriteString(fmt.Sprintf("   e^(iπ) = %v\n", eulerResult))
	sb.WriteString(fmt.Sprintf("   e^(iπ) + 1 ≈ %v (接近 0)\n", eulerResult+1))

	// 9. 纯虚数和纯实数
	sb.WriteString("\n9. 纯虚数和纯实数\n")

	pureImag := 0 + 5i // 纯虚数
	pureReal := 7 + 0i // 纯实数（复数形式）

	sb.WriteString(fmt.Sprintf("   纯虚数: %v (实部为 0)\n", pureImag))
	sb.WriteString(fmt.Sprintf("   纯实数: %v (虚部为 0)\n", pureReal))

	// 10. 最佳实践
	sb.WriteString("\n10. 最佳实践\n")
	sb.WriteString("   ✓ 虚数字面量必须使用小写 i（不能用大写 I）\n")
	sb.WriteString("   ✓ 默认使用 complex128（更高精度）\n")
	sb.WriteString("   ✓ 使用 complex() 函数从实部和虚部构造复数\n")
	sb.WriteString("   ✓ 使用 real() 和 imag() 提取实部和虚部\n")
	sb.WriteString("   ✓ 导入 math/cmplx 包进行高级复数运算\n")
	sb.WriteString("   ✓ 复数常用于信号处理、量子计算、电路分析等领域\n")

	sb.WriteString("\n" + repeatString("=", 60) + "\n")

	return sb.String()
}

// DisplayImaginary 展示 Go 语言中虚数字面量和复数的使用
// Go 原生支持复数运算，虚数字面量以 i 结尾
func DisplayImaginary() {
	fmt.Print(GetImaginaryContent())
}
