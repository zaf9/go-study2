# CLI Menu Contract: Go Constants 学习包

**Feature**: 004-constants-learning  
**Date**: 2025-12-05  
**Phase**: 1 - Design & Contracts

## 接口定义

### DisplayMenu - 主菜单函数

```go
// DisplayMenu 显示 Constants 主菜单,支持子主题选择
//
// 参数:
//   - stdin: 用户输入流
//   - stdout: 正常输出流
//   - stderr: 错误输出流
//
// 行为:
//   - 循环显示菜单直到用户输入 'q' 退出
//   - 输入 0-11 选择对应子主题
//   - 无效输入显示错误提示并重新显示菜单
//
// 示例交互:
//   Constants 学习菜单
//   ---------------------------------
//   请选择要学习的主题:
//   0. Boolean Constants (布尔常量)
//   1. Rune Constants (符文常量)
//   ...
//   q. 返回上级菜单
//
//   请输入您的选择: 0
//   [显示布尔常量内容]
func DisplayMenu(stdin io.Reader, stdout, stderr io.Writer)
```

### Display{Topic} - 主题显示函数

```go
// Display{Topic} 显示特定主题的学习内容
//
// 函数列表:
//   - DisplayBoolean()
//   - DisplayRune()
//   - DisplayInteger()
//   - DisplayFloatingPoint()
//   - DisplayComplex()
//   - DisplayString()
//   - DisplayExpressions()
//   - DisplayTypedUntyped()
//   - DisplayConversions()
//   - DisplayBuiltinFunctions()
//   - DisplayIota()
//   - DisplayImplementationRestrictions()
//
// 参数: 无(直接输出到 stdout)
//
// 输出格式:
//   === {主题标题} ===
//
//   【概念说明】
//   {主题概念的详细说明}
//
//   【语法规则】
//   {语法规则说明}
//
//   【使用场景】
//   {典型使用场景}
//
//   【示例代码】
//   示例 1: {示例标题}
//   {示例代码}
//   // 输出: {预期输出}
//
//   【常见错误】
//   {常见错误和注意事项}
func Display{Topic}()
```

## 菜单编号方案

| 编号 | 主题键 | 主题名称 |
|------|--------|---------|
| 0 | boolean | Boolean Constants (布尔常量) |
| 1 | rune | Rune Constants (符文常量) |
| 2 | integer | Integer Constants (整数常量) |
| 3 | floating_point | Floating-point Constants (浮点常量) |
| 4 | complex | Complex Constants (复数常量) |
| 5 | string | String Constants (字符串常量) |
| 6 | expressions | Constant Expressions (常量表达式) |
| 7 | typed_untyped | Typed and Untyped Constants (类型化/无类型化常量) |
| 8 | conversions | Conversions (类型转换) |
| 9 | builtin_functions | Built-in Functions (内置函数) |
| 10 | iota | Iota (iota 特性) |
| 11 | implementation_restrictions | Implementation Restrictions (实现限制) |
| q | - | 返回上级菜单 |

## 错误处理

### 无效输入

**输入**: 不在 0-11 或 'q' 范围内的字符

**输出**:
```
无效的选择,请重试。
```

**行为**: 重新显示菜单

### 读取错误

**场景**: stdin 读取失败

**输出** (到 stderr):
```
读取输入错误: {错误信息}
```

**行为**: 函数返回

## 集成到主菜单

### main.go 修改

```go
menu: map[string]MenuItem{
    "0": {
        Description: "Lexical elements",
        Action:      lexical_elements.DisplayMenu,
    },
    "1": {  // 新增
        Description: "Constants",
        Action:      constants.DisplayMenu,
    },
    // Add new items here
},
```

## 测试要求

### 单元测试

```go
func TestDisplayMenu(t *testing.T) {
    // 测试菜单显示
    // 测试有效输入处理
    // 测试无效输入处理
    // 测试退出功能
}

func TestDisplay{Topic}(t *testing.T) {
    // 捕获 stdout 输出
    // 验证包含关键内容
}
```

### 集成测试

```go
func TestMainMenuIntegration(t *testing.T) {
    // 测试从主菜单进入 Constants 菜单
    // 测试选择子主题
    // 测试返回主菜单
}
```

## 验收标准

- [ ] 所有 12 个子主题可通过菜单访问
- [ ] 输入 'q' 正确返回上级菜单
- [ ] 无效输入显示友好错误提示
- [ ] 所有提示信息使用中文
- [ ] 菜单编号从 0 开始
- [ ] 与 lexical_elements 菜单风格一致
