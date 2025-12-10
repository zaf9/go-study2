# Constants 学习模块

本包提供 Go 语言常量(Constants)相关知识的学习内容。

## 概述

常量是 Go 语言中重要的概念,本模块涵盖以下 12 个子主题:

### 基础常量类型 (User Story 1 - P1)
1. **Boolean Constants** - 布尔常量 (`true`, `false`)
2. **Rune Constants** - 符文常量 (Unicode 字符)
3. **Integer Constants** - 整数常量 (各种进制表示)
4. **Floating-point Constants** - 浮点常量 (小数和科学计数法)
5. **Complex Constants** - 复数常量 (实部和虚部)
6. **String Constants** - 字符串常量 (解释型和原始型)

### 常量表达式和类型 (User Story 2 - P2)
7. **Constant Expressions** - 常量表达式 (编译时求值)
8. **Typed and Untyped Constants** - 类型化/无类型化常量

### 转换和内置函数 (User Story 3 - P3)
9. **Conversions** - 类型转换
10. **Built-in Functions** - 内置函数 (len, min, max 等)

### 特殊常量 (User Story 4 - P4)
11. **Iota** - iota 特性 (枚举、位掩码)
12. **Implementation Restrictions** - 编译器实现限制

## 使用方式

### 命令行模式

```bash
# 启动 CLI 交互式学习
cd backend && go run main.go

# 在主菜单中选择 "1. Constants"
# 然后选择具体的子主题编号 (0-11)
```

### HTTP 模式

```bash
# 启动 HTTP 服务器
cd backend && go run main.go -d

# 获取 Constants 菜单
curl http://localhost:8080/api/v1/topic/constants

# 获取具体子主题内容
curl http://localhost:8080/api/v1/topic/constants/boolean
curl http://localhost:8080/api/v1/topic/constants/integer
```

## 文件结构

```
constants/
├── README.md                      # 本文件
├── constants.go                   # 主入口: DisplayMenu 函数
├── constants_test.go              # 主入口测试
│
├── boolean.go                     # 布尔常量
├── boolean_test.go
├── rune.go                        # 符文常量
├── rune_test.go
├── integer.go                     # 整数常量
├── integer_test.go
├── floating_point.go              # 浮点常量
├── floating_point_test.go
├── complex.go                     # 复数常量
├── complex_test.go
├── string.go                      # 字符串常量
├── string_test.go
│
├── expressions.go                 # 常量表达式
├── expressions_test.go
├── typed_untyped.go               # 类型化/无类型化常量
├── typed_untyped_test.go
├── conversions.go                 # 类型转换
├── conversions_test.go
├── builtin_functions.go           # 内置函数
├── builtin_functions_test.go
├── iota.go                        # iota 特性
├── iota_test.go
├── implementation_restrictions.go # 实现限制
└── implementation_restrictions_test.go
```

## 开发指南

### 添加新示例

每个子主题文件包含两个主要函数:

1. `Get{Topic}Content() string` - 返回学习内容字符串 (用于 HTTP API)
2. `Display{Topic}()` - 输出学习内容到 stdout (用于 CLI)

示例代码应遵循以下规范:
- 所有注释使用中文
- 每个示例包含完整的可运行代码
- 包含预期输出说明
- 由简到难,覆盖典型用法和边界情况

### 运行测试

```bash
# 运行所有 Constants 测试
cd backend && go test ./internal/app/constants/...

# 查看测试覆盖率
cd backend && go test -cover ./internal/app/constants/...

# 生成覆盖率报告
cd backend && go test -coverprofile=coverage.out ./internal/app/constants/...
go tool cover -html=coverage.out
```

## 相关文档

- [Go 语言规范 - 常量](https://go.dev/ref/spec#Constants)
- [项目规范文档](../../../specs/004-constants-learning/)
