package lexical_elements

import (
	"fmt"
	"math"
	"strings"
)

// GetFloatsContent 返回浮点数字面量相关的学习内容
func GetFloatsContent() string {
	var sb strings.Builder

	sb.WriteString("\n" + repeatString("=", 60) + "\n")
	sb.WriteString("【词法元素 - 浮点数字面量 (Floating-point Literals)】\n")
	sb.WriteString(repeatString("=", 60) + "\n")

	// 1. 基本浮点数表示
	sb.WriteString("\n1. 基本浮点数表示（小数形式）\n")
	sb.WriteString("   浮点数由整数部分、小数点和小数部分组成\n")

	f1 := 3.14
	f2 := 0.5
	f3 := 123.456
	f4 := .25 // 可以省略整数部分的 0
	f5 := 5.  // 可以省略小数部分（等同于 5.0）

	sb.WriteString(fmt.Sprintf("   3.14 = %.2f\n", f1))
	sb.WriteString(fmt.Sprintf("   0.5 = %.1f\n", f2))
	sb.WriteString(fmt.Sprintf("   123.456 = %.3f\n", f3))
	sb.WriteString(fmt.Sprintf("   .25 = %.2f (省略整数部分的 0)\n", f4))
	sb.WriteString(fmt.Sprintf("   5. = %.1f (省略小数部分)\n", f5))

	// 2. 科学计数法（指数形式）
	sb.WriteString("\n2. 科学计数法（指数形式）\n")
	sb.WriteString("   使用 e 或 E 表示 10 的幂次\n")
	sb.WriteString("   格式: 尾数e指数 或 尾数E指数\n")

	sci1 := 1e3       // 1 × 10^3 = 1000
	sci2 := 2.5e2     // 2.5 × 10^2 = 250
	sci3 := 3.14e-2   // 3.14 × 10^-2 = 0.0314
	sci4 := 6.022e23  // 阿伏伽德罗常数
	sci5 := 1.602e-19 // 电子电荷（使用大写 E）

	sb.WriteString(fmt.Sprintf("   1e3 = %.0f (1 × 10³)\n", sci1))
	sb.WriteString(fmt.Sprintf("   2.5e2 = %.0f (2.5 × 10²)\n", sci2))
	sb.WriteString(fmt.Sprintf("   3.14e-2 = %.4f (3.14 × 10⁻²)\n", sci3))
	sb.WriteString(fmt.Sprintf("   6.022e23 = %.3e (阿伏伽德罗常数)\n", sci4))
	sb.WriteString(fmt.Sprintf("   1.602E-19 = %.3e (电子电荷)\n", sci5))

	// 3. 十六进制浮点数（不常用）
	sb.WriteString("\n3. 十六进制浮点数（较少使用）\n")
	sb.WriteString("   以 0x 或 0X 开头，使用 p 或 P 表示 2 的幂次\n")

	hexFloat1 := 0x1.8p1 // 1.5 × 2^1 = 3.0
	hexFloat2 := 0x1p-2  // 1 × 2^-2 = 0.25

	sb.WriteString(fmt.Sprintf("   0x1.8p1 = %.1f (1.5 × 2¹)\n", hexFloat1))
	sb.WriteString(fmt.Sprintf("   0x1p-2 = %.2f (1 × 2⁻²)\n", hexFloat2))

	// 4. 浮点数类型
	sb.WriteString("\n4. 浮点数类型\n")
	sb.WriteString("   Go 提供两种浮点数类型：\n")
	sb.WriteString("   float32 - 32 位浮点数（单精度，约 7 位十进制精度）\n")
	sb.WriteString("   float64 - 64 位浮点数（双精度，约 15 位十进制精度，默认）\n")

	var f32 float32 = 3.14159265
	var f64 float64 = 3.14159265358979323846

	sb.WriteString(fmt.Sprintf("   float32: %.10f (精度有限)\n", f32))
	sb.WriteString(fmt.Sprintf("   float64: %.20f (更高精度)\n", f64))

	// 5. 特殊浮点数值
	sb.WriteString("\n5. 特殊浮点数值\n")

	sb.WriteString(fmt.Sprintf("   正无穷大: %f (math.Inf(1))\n", math.Inf(1)))
	sb.WriteString(fmt.Sprintf("   负无穷大: %f (math.Inf(-1))\n", math.Inf(-1)))
	sb.WriteString(fmt.Sprintf("   非数字 (NaN): %f (math.NaN())\n", math.NaN()))
	sb.WriteString(fmt.Sprintf("   最大 float64: %e (math.MaxFloat64)\n", math.MaxFloat64))
	sb.WriteString(fmt.Sprintf("   最小正 float64: %e (math.SmallestNonzeroFloat64)\n", math.SmallestNonzeroFloat64))

	// 6. 浮点数运算
	sb.WriteString("\n6. 浮点数运算示例\n")
	a := 10.5
	b := 3.2
	sb.WriteString(fmt.Sprintf("   a = %.1f, b = %.1f\n", a, b))
	sb.WriteString(fmt.Sprintf("   a + b = %.1f\n", a+b))
	sb.WriteString(fmt.Sprintf("   a - b = %.1f\n", a-b))
	sb.WriteString(fmt.Sprintf("   a * b = %.2f\n", a*b))
	sb.WriteString(fmt.Sprintf("   a / b = %.4f\n", a/b))

	// 7. 精度问题
	sb.WriteString("\n7. 浮点数精度问题\n")
	sb.WriteString("   浮点数在计算机中以二进制表示，可能存在精度误差\n")

	x := 0.1
	y := 0.2
	sum := x + y
	sb.WriteString(fmt.Sprintf("   0.1 + 0.2 = %.17f (不完全等于 0.3)\n", sum))
	sb.WriteString(fmt.Sprintf("   0.3 == 0.1 + 0.2: %t (精度误差)\n", 0.3 == sum))

	// 正确的比较方法
	epsilon := 1e-9
	sb.WriteString(fmt.Sprintf("   使用 epsilon 比较: %t (|0.3 - sum| < 1e-9)\n",
		math.Abs(0.3-sum) < epsilon))

	// 8. 格式化输出
	sb.WriteString("\n8. 浮点数格式化输出\n")
	val := 123.456789
	sb.WriteString(fmt.Sprintf("   %%f:  %f (默认 6 位小数)\n", val))
	sb.WriteString(fmt.Sprintf("   %%.2f: %.2f (保留 2 位小数)\n", val))
	sb.WriteString(fmt.Sprintf("   %%e:  %e (科学计数法)\n", val))
	sb.WriteString(fmt.Sprintf("   %%E:  %E (科学计数法，大写 E)\n", val))
	sb.WriteString(fmt.Sprintf("   %%g:  %g (自动选择 %%f 或 %%e)\n", val))

	// 9. 类型转换
	sb.WriteString("\n9. 浮点数类型转换\n")

	var intVal int = 42
	var floatVal float64 = float64(intVal)
	sb.WriteString(fmt.Sprintf("   int(%d) → float64(%.1f)\n", intVal, floatVal))

	var f float64 = 3.99
	var i int = int(f) // 截断小数部分
	sb.WriteString(fmt.Sprintf("   float64(%.2f) → int(%d) (截断)\n", f, i))

	// 10. 最佳实践
	sb.WriteString("\n10. 最佳实践\n")
	sb.WriteString("   ✓ 默认使用 float64（更高精度）\n")
	sb.WriteString("   ✓ 使用科学计数法表示极大或极小的数\n")
	sb.WriteString("   ✓ 比较浮点数时使用 epsilon 容差\n")
	sb.WriteString("   ✗ 避免直接用 == 比较浮点数\n")
	sb.WriteString("   ✗ 不要用浮点数表示货币（使用整数分或专门的库）\n")
	sb.WriteString("   ✓ 了解浮点数的精度限制\n")

	sb.WriteString("\n" + repeatString("=", 60) + "\n")

	return sb.String()
}

// DisplayFloats 展示 Go 语言中浮点数字面量的各种表示形式
// 包括小数形式、科学计数法等
func DisplayFloats() {
	fmt.Print(GetFloatsContent())
}
