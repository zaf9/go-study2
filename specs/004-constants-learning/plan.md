# Implementation Plan: Go Constants 学习包

**Branch**: `004-constants-learning` | **Date**: 2025-12-05 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/004-constants-learning/spec.md`

## Summary

本功能为 go-study2 项目添加 Constants 学习模块,支持通过命令行交互式菜单和 HTTP API 两种方式学习 Go 语言规范中的常量相关知识。学习内容涵盖基础常量类型(布尔、符文、整数、浮点、复数、字符串)、常量表达式、类型化/无类型化常量、类型转换、内置函数、iota 特性和编译器实现限制等12个子主题。每个子主题提供详细的中文说明和至少3个可运行的示例代码,确保学习者能够全面理解 Go 常量系统。

技术实现遵循项目现有架构模式:在 `internal/app/constants` 包下创建模块化的学习内容文件,集成到主菜单系统,并通过 HTTP handler 提供 RESTful API 访问。所有代码和文档使用中文,单元测试覆盖率达到 80% 以上。

## Technical Context

**Language/Version**: Go 1.24.5  
**Primary Dependencies**: GoFrame v2.9.5 (HTTP 服务框架)  
**Storage**: 文件系统(学习内容以 .go 文件形式存储,包含注释和示例代码)  
**Testing**: Go 标准测试框架 (`go test`), 目标覆盖率 ≥80%  
**Target Platform**: 跨平台 (Windows/Linux/macOS), 支持命令行和 HTTP 服务两种运行模式  
**Project Type**: 单体应用 (CLI + HTTP Server)  
**Performance Goals**: 
- HTTP 响应时间 <100ms (正常负载 100 并发请求)
- 支持 1000 并发 HTTP 请求无性能下降
- 命令行菜单响应即时 (<50ms)

**Constraints**: 
- 所有用户可见文档和代码注释必须使用中文 (宪章 Principle III, XIII)
- 单元测试覆盖率 ≥80% (宪章 Principle VI)
- 代码简洁清晰,适合 Go 初学者阅读 (宪章 Principle I)
- 避免深层嵌套逻辑 (宪章 Principle IV)
- 遵循 YAGNI 原则,不引入不必要的复杂性 (宪章 Principle V)

**Scale/Scope**: 
- 12 个学习子主题
- 每个子主题至少 3 个示例代码
- 预计总代码量 ~2000-3000 行 (含注释和示例)
- 支持数百名学习者并发访问

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Simplicity):** ✅ 采用与现有 lexical_elements 模块相同的简单架构模式,每个子主题一个独立文件,易于理解和维护
- **Principle II (Comments):** ✅ 每个学习主题文件包含详细的中文注释,分层说明主题概念、语法规则、使用场景和常见错误
- **Principle III (Language):** ✅ 所有代码注释、函数说明、示例输出和用户提示均使用中文
- **Principle IV (Nesting):** ✅ 菜单导航使用 map 映射和 switch/if 单层判断,避免深层嵌套;示例代码保持简洁
- **Principle V (YAGNI):** ✅ 复用现有菜单系统和 HTTP 框架,不引入新的设计模式或抽象层;学习内容直接以函数形式提供
- **Principle VI (Testing):** ✅ 为每个 Display 函数编写单元测试,验证输出内容;为 HTTP handler 编写集成测试;目标覆盖率 ≥80%
- **Principle XVII (Hierarchical Menu):** ✅ 支持主菜单 → Constants 子菜单 → 具体主题的多层导航,提供 'q' 返回上级菜单
- **Principle XVIII (Dual Learning Mode):** ✅ 同时支持命令行交互式学习和 HTTP API 访问两种模式

## Project Structure

### Documentation (this feature)

```text
specs/004-constants-learning/
├── spec.md                    # 功能规范 (已完成)
├── plan.md                    # 本文件 - 实现计划
├── research.md                # Phase 0 输出 - 技术调研
├── data-model.md              # Phase 1 输出 - 数据模型设计
├── contracts/                 # Phase 1 输出 - 接口契约
│   ├── cli-menu.md           # 命令行菜单接口规范
│   └── http-api.md           # HTTP API 接口规范
├── checklists/                # 质量检查清单
│   └── requirements.md       # 需求质量检查 (已完成)
└── tasks.md                   # Phase 2 输出 - 任务分解 (由 /speckit.tasks 生成)
```

### Source Code (repository root)

```text
go-study2/
├── main.go                    # 主入口 (需修改: 添加 Constants 菜单项)
├── main_test.go               # 主程序测试 (需修改: 添加 Constants 菜单测试)
├── go.mod                     # 依赖管理 (无需修改)
├── config.yaml                # 配置文件 (无需修改)
│
├── internal/
│   ├── app/
│   │   ├── constants/         # 【新增】Constants 学习模块
│   │   │   ├── README.md                      # 包文档说明
│   │   │   ├── constants.go                   # 主入口: DisplayMenu 函数
│   │   │   ├── constants_test.go              # 主入口测试
│   │   │   │
│   │   │   ├── boolean.go                     # 布尔常量
│   │   │   ├── boolean_test.go
│   │   │   ├── rune.go                        # 符文常量
│   │   │   ├── rune_test.go
│   │   │   ├── integer.go                     # 整数常量
│   │   │   ├── integer_test.go
│   │   │   ├── floating_point.go              # 浮点常量
│   │   │   ├── floating_point_test.go
│   │   │   ├── complex.go                     # 复数常量
│   │   │   ├── complex_test.go
│   │   │   ├── string.go                      # 字符串常量
│   │   │   ├── string_test.go
│   │   │   │
│   │   │   ├── expressions.go                 # 常量表达式
│   │   │   ├── expressions_test.go
│   │   │   ├── typed_untyped.go               # 类型化/无类型化常量
│   │   │   ├── typed_untyped_test.go
│   │   │   ├── conversions.go                 # 类型转换
│   │   │   ├── conversions_test.go
│   │   │   ├── builtin_functions.go           # 内置函数
│   │   │   ├── builtin_functions_test.go
│   │   │   ├── iota.go                        # iota 特性
│   │   │   ├── iota_test.go
│   │   │   ├── implementation_restrictions.go # 实现限制
│   │   │   └── implementation_restrictions_test.go
│   │   │
│   │   ├── http_server/       # HTTP 服务器 (需修改)
│   │   │   ├── router.go                      # 路由注册 (需添加 Constants 路由)
│   │   │   ├── handler/
│   │   │   │   ├── handler.go                 # Handler 基础结构
│   │   │   │   ├── constants.go               # 【新增】Constants HTTP handler
│   │   │   │   └── constants_test.go          # 【新增】Constants handler 测试
│   │   │   └── middleware/
│   │   │
│   │   └── lexical_elements/  # 现有词法元素模块 (参考实现)
│   │
│   └── config/                # 配置管理 (无需修改)
│
└── tests/                     # 集成测试
    └── integration/
        └── constants_api_test.go  # 【新增】Constants API 集成测试
```

**Structure Decision**: 

采用单体应用结构,遵循项目现有的 `internal/app/{module}` 组织模式。Constants 学习模块作为独立包放置在 `internal/app/constants/`,与现有的 `lexical_elements` 模块平级。

每个常量子主题(如 boolean, integer, expressions 等)作为独立的 .go 文件,包含:
1. 详细的中文注释说明(主题概念、语法规则、使用场景)
2. `Display{Topic}()` 函数用于命令行输出
3. 至少 3 个可运行的示例代码(以注释形式嵌入)
4. 对应的 `*_test.go` 测试文件

HTTP 集成通过在 `internal/app/http_server/handler/` 下新增 `constants.go` handler 实现,路由注册在 `router.go` 中添加。

这种结构的优势:
- 模块化清晰,每个子主题独立维护
- 与现有 lexical_elements 模块保持一致性,降低学习成本
- 易于扩展,未来添加新主题只需新增文件
- 测试文件与源文件一一对应,便于维护

## Complexity Tracking

> **无宪章违规** - 本实现完全符合项目宪章所有原则,无需复杂性豁免。

---

## Phase 0: Research & Technical Validation

**目标**: 验证技术可行性,确认实现方案,识别潜在风险

### 0.1 现有架构分析

**任务**: 深入分析 `lexical_elements` 模块的实现模式,作为 Constants 模块的参考蓝图

**输出**: `research.md` 第1节 - 现有架构分析

**关键问题**:
1. `lexical_elements.DisplayMenu()` 的菜单导航实现机制
2. 各子主题 `Display{Topic}()` 函数的标准结构和输出格式
3. HTTP handler 如何调用学习内容并转换为 JSON 响应
4. 测试策略:如何测试 Display 函数的输出内容
5. 主菜单 `main.go` 的集成方式

**验收标准**:
- 文档化 lexical_elements 的完整调用链路(CLI 和 HTTP 两条路径)
- 提取可复用的代码模式和命名约定
- 识别需要修改的集成点(main.go, router.go, handler/)

### 0.2 Go Constants 规范研究

**任务**: 深入研究 Go 1.24 语言规范中 Constants 章节,提取学习要点

**输出**: `research.md` 第2节 - Constants 规范要点

**关键内容**:
1. 6 种常量类型的定义和字面量表示方法
2. 常量表达式的求值规则和运算符优先级
3. 类型化常量 vs 无类型化常量的区别和默认类型映射
4. 常量转换规则和可表示性(representability)要求
5. 可用于常量的内置函数列表及使用限制
6. iota 的自增规则和常见使用模式
7. 编译器实现限制(精度、范围)

**验收标准**:
- 为每个子主题列出至少 5 个示例场景
- 识别常见错误和易混淆点
- 确认所有内容符合 Go 1.24 规范

### 0.3 示例代码设计

**任务**: 为 12 个子主题设计高质量的示例代码

**输出**: `research.md` 第3节 - 示例代码清单

**设计原则**:
- 每个示例独立可运行,包含完整的 package 和 main 函数
- 示例由简到难,覆盖典型用法和边界情况
- 包含预期输出的注释,便于学习者验证理解
- 代码简洁清晰,符合 Go 编码规范和项目宪章

**验收标准**:
- 每个子主题至少 3 个示例(部分复杂主题如 expressions, iota 需 5+ 个)
- 所有示例代码通过 `go fmt`, `go vet` 检查
- 示例覆盖 spec.md 中定义的所有验收场景

### 0.4 性能和并发考虑

**任务**: 评估 HTTP 模式下的性能需求,设计并发安全方案

**输出**: `research.md` 第4节 - 性能和并发设计

**关键问题**:
1. 学习内容是否需要缓存?(内容静态,可考虑启动时加载)
2. HTTP handler 的并发安全性(GoFrame 框架的并发处理机制)
3. 如何满足 100ms 响应时间和 1000 并发请求的性能目标

**验收标准**:
- 确定内容加载策略(动态生成 vs 预加载)
- 识别潜在的性能瓶颈和优化方案
- 设计性能测试方案

### 0.5 风险识别

**任务**: 识别实现过程中的技术风险和缓解措施

**输出**: `research.md` 第5节 - 风险和缓解措施

**已识别风险**:
1. **示例代码维护成本高**: 12 个主题 × 3+ 示例 = 36+ 代码片段需保持正确性
   - 缓解: 编写自动化测试验证示例代码可编译性
2. **中文注释质量不一致**: 多个文件的注释风格可能不统一
   - 缓解: 制定注释模板,代码审查时检查
3. **HTTP 响应格式设计**: JSON 结构需平衡可读性和扩展性
   - 缓解: 在 Phase 1 设计明确的 API 契约
4. **测试覆盖率达标难度**: Display 函数输出字符串,测试断言复杂
   - 缓解: 使用 golden file 测试或关键内容片段匹配

**验收标准**:
- 每个风险有明确的缓解措施
- 高风险项有备选方案

---

## Phase 1: Design & Contracts

**目标**: 完成详细设计,定义清晰的接口契约

### 1.1 数据模型设计

**任务**: 定义学习内容的数据结构

**输出**: `data-model.md`

**核心结构**:

```go
// TopicContent 表示一个学习主题的完整内容
type TopicContent struct {
    Title       string        // 主题标题(中文)
    Description string        // 主题说明
    Syntax      string        // 语法规则
    Examples    []CodeExample // 示例代码列表
    CommonErrors string       // 常见错误说明
}

// CodeExample 表示一个代码示例
type CodeExample struct {
    Title       string // 示例标题
    Code        string // 示例代码
    Output      string // 预期输出
    Explanation string // 说明
}

// MenuOption 表示菜单选项
type MenuOption struct {
    Key         string // 选项键(如 "0", "1")
    Description string // 选项描述
    Action      func() // 执行函数
}
```

**验收标准**:
- 数据结构支持所有 spec.md 中定义的内容要素
- 结构设计简洁,避免过度抽象
- 包含 JSON 序列化标签,支持 HTTP 响应

### 1.2 CLI 菜单接口契约

**任务**: 定义命令行菜单的交互规范

**输出**: `contracts/cli-menu.md`

**接口定义**:

```go
// DisplayMenu 显示 Constants 主菜单,支持子主题选择
// 参数:
//   - stdin: 用户输入流
//   - stdout: 正常输出流
//   - stderr: 错误输出流
// 行为:
//   - 循环显示菜单直到用户输入 'q' 退出
//   - 输入 0-11 选择对应子主题
//   - 无效输入显示错误提示并重新显示菜单
func DisplayMenu(stdin io.Reader, stdout, stderr io.Writer)

// Display{Topic} 显示特定主题的学习内容
// 例如: DisplayBoolean(), DisplayInteger(), DisplayExpressions()
// 参数: 无(直接输出到 stdout)
// 行为:
//   - 输出主题标题、说明、语法规则
//   - 输出所有示例代码和预期输出
//   - 输出常见错误说明
func Display{Topic}()
```

**菜单编号方案**:
```
0. Boolean Constants (布尔常量)
1. Rune Constants (符文常量)
2. Integer Constants (整数常量)
3. Floating-point Constants (浮点常量)
4. Complex Constants (复数常量)
5. String Constants (字符串常量)
6. Constant Expressions (常量表达式)
7. Typed and Untyped Constants (类型化/无类型化常量)
8. Conversions (类型转换)
9. Built-in Functions (内置函数)
10. Iota (iota 特性)
11. Implementation Restrictions (实现限制)
q. 返回上级菜单
```

**验收标准**:
- 菜单编号从 0 开始,符合项目约定
- 所有提示信息使用中文
- 错误处理清晰,用户体验友好

### 1.3 HTTP API 接口契约

**任务**: 定义 HTTP API 的请求/响应规范

**输出**: `contracts/http-api.md`

**端点定义**:

```
GET /api/v1/topic/constants
描述: 获取 Constants 主题的子菜单列表
响应:
{
  "code": 0,
  "message": "success",
  "data": {
    "title": "Constants (常量)",
    "subtopics": [
      {"key": "boolean", "name": "Boolean Constants (布尔常量)"},
      {"key": "rune", "name": "Rune Constants (符文常量)"},
      {"key": "integer", "name": "Integer Constants (整数常量)"},
      {"key": "floating_point", "name": "Floating-point Constants (浮点常量)"},
      {"key": "complex", "name": "Complex Constants (复数常量)"},
      {"key": "string", "name": "String Constants (字符串常量)"},
      {"key": "expressions", "name": "Constant Expressions (常量表达式)"},
      {"key": "typed_untyped", "name": "Typed and Untyped Constants (类型化/无类型化常量)"},
      {"key": "conversions", "name": "Conversions (类型转换)"},
      {"key": "builtin_functions", "name": "Built-in Functions (内置函数)"},
      {"key": "iota", "name": "Iota (iota 特性)"},
      {"key": "implementation_restrictions", "name": "Implementation Restrictions (实现限制)"}
    ]
  }
}

GET /api/v1/topic/constants/:subtopic
描述: 获取特定子主题的学习内容
路径参数: subtopic (如 boolean, integer, expressions)
响应:
{
  "code": 0,
  "message": "success",
  "data": {
    "title": "Boolean Constants (布尔常量)",
    "description": "布尔常量表示真值...",
    "syntax": "const name [type] = value",
    "examples": [
      {
        "title": "基本布尔常量声明",
        "code": "const enabled = true\nconst disabled = false",
        "output": "// 编译时确定,无运行时输出",
        "explanation": "布尔常量只有 true 和 false 两个值"
      }
    ],
    "common_errors": "常见错误: 1. 尝试将非布尔值赋给布尔常量..."
  }
}

错误响应 (404):
{
  "code": 404,
  "message": "subtopic not found: xyz",
  "data": null
}
```

**性能要求**:
- 正常负载(100 并发)响应时间 <100ms
- 支持 1000 并发请求无错误

**验收标准**:
- API 遵循 RESTful 规范
- 响应格式与现有 lexical_elements API 保持一致
- 包含完整的错误处理规范

### 1.4 测试策略设计

**任务**: 设计测试方案,确保 80% 覆盖率

**输出**: `contracts/testing-strategy.md`

**测试层次**:

1. **单元测试** (覆盖目标: 85%)
   - 每个 `Display{Topic}()` 函数测试输出内容包含关键字
   - 菜单导航逻辑测试(输入验证、选项映射)
   - HTTP handler 函数测试(请求解析、响应构造)

2. **集成测试** (覆盖目标: 主要流程)
   - CLI 菜单完整交互流程
   - HTTP API 端到端测试(启动服务器、发送请求、验证响应)

3. **示例代码验证测试**
   - 自动提取示例代码并编译验证
   - 确保所有示例可运行

**测试工具**:
- Go 标准测试框架 (`testing` 包)
- `httptest` 包用于 HTTP handler 测试
- 自定义测试辅助函数(如 `assertContains` 检查输出内容)

**验收标准**:
- 测试策略覆盖所有验收场景
- 明确覆盖率计算方法和目标
- 包含性能测试方案

---

## Phase 2: Implementation (由 /speckit.tasks 生成详细任务)

Phase 2 的详细任务分解将在执行 `/speckit.tasks` 命令时生成到 `tasks.md` 文件中。

**预期任务类别**:

1. **基础设施任务** (优先级: P1)
   - 创建 `internal/app/constants/` 包结构
   - 编写包文档 `README.md`
   - 创建主入口 `constants.go` 和 `DisplayMenu()` 函数

2. **学习内容实现任务** (优先级: P1-P4,对应 spec.md 的用户故事优先级)
   - P1: 实现 6 个基础常量类型文件(boolean, rune, integer, floating_point, complex, string)
   - P2: 实现常量表达式和类型相关文件(expressions, typed_untyped)
   - P3: 实现转换和内置函数文件(conversions, builtin_functions)
   - P4: 实现特殊常量文件(iota, implementation_restrictions)

3. **CLI 集成任务** (优先级: P1)
   - 修改 `main.go` 添加 Constants 菜单项
   - 修改 `main_test.go` 添加菜单测试

4. **HTTP 集成任务** (优先级: P1)
   - 实现 `handler/constants.go` handler
   - 修改 `router.go` 注册 Constants 路由
   - 编写 handler 单元测试

5. **测试任务** (优先级: P1,与实现任务并行)
   - 为每个学习内容文件编写单元测试
   - 编写集成测试
   - 验证测试覆盖率 ≥80%

6. **文档和收尾任务** (优先级: P2)
   - 更新 `README.md` 添加 Constants 学习模块说明
   - 代码审查和优化
   - 性能测试和调优

---

## Phase 3: Validation & Polish

**目标**: 确保实现质量,验证所有验收标准

### 3.1 功能验收测试

**验证项**:
- [ ] 所有 12 个子主题的 CLI 菜单可正常访问
- [ ] 每个子主题至少包含 3 个示例代码
- [ ] 所有示例代码可编译和运行
- [ ] HTTP API 所有端点正常响应
- [ ] 错误处理符合规范(无效输入、404 等)

### 3.2 质量指标验证

**验证项**:
- [ ] 单元测试覆盖率 ≥80% (运行 `go test -cover`)
- [ ] 所有代码通过 `go fmt`, `go vet`, `golint` 检查
- [ ] 所有注释和文档使用中文
- [ ] HTTP 响应时间 <100ms (100 并发负载测试)
- [ ] 支持 1000 并发请求无错误(压力测试)

### 3.3 宪章合规性审查

**验证项**:
- [ ] Principle I: 代码简洁清晰,适合初学者
- [ ] Principle II: 注释分层清晰,说明充分
- [ ] Principle III: 所有文档和注释使用中文
- [ ] Principle IV: 无深层嵌套逻辑
- [ ] Principle V: 无不必要的复杂性
- [ ] Principle VI: 测试覆盖率 ≥80%
- [ ] Principle XVII: 菜单导航层次清晰
- [ ] Principle XVIII: 支持 CLI 和 HTTP 双模式

### 3.4 用户体验测试

**验证项**:
- [ ] 菜单提示清晰易懂
- [ ] 错误消息友好且有指导性
- [ ] 学习内容组织合理,由浅入深
- [ ] 示例代码易于理解和实践

### 3.5 性能优化

**优化方向**:
- 如性能测试发现瓶颈,考虑内容预加载或缓存
- 优化 JSON 序列化性能
- 减少不必要的内存分配

---

## Success Criteria Mapping

将 spec.md 中的成功标准映射到实现计划:

| Success Criteria | 实现方案 | 验证方法 |
|------------------|---------|---------|
| SC-001: 5分钟内访问任意子主题 | 简洁的菜单导航,快速的 HTTP 响应 | 用户体验测试,性能测试 |
| SC-002: 每个子主题 ≥3 个示例 | Phase 0.3 设计示例清单,Phase 2 实现 | 代码审查,自动化计数 |
| SC-003: 所有示例可编译运行 | 示例代码验证测试 | 自动化编译测试 |
| SC-004: 测试覆盖率 ≥80% | 完善的测试策略(Phase 1.4) | `go test -cover` |
| SC-005: HTTP 响应 <100ms | 简洁的 handler 实现,可选缓存 | 性能测试(100 并发) |
| SC-006: 内容 100% 准确 | Phase 0.2 深入研究规范,代码审查 | 专家审查,规范对照 |
| SC-007: 90% 首次成功率 | 清晰的菜单和提示,友好的错误处理 | 用户体验测试 |
| SC-008: 1000 并发无错误 | GoFrame 框架并发支持,无状态设计 | 压力测试 |

---

## Dependencies & Integration Points

### 依赖项

**外部依赖**:
- GoFrame v2.9.5 (已在 go.mod 中)
- Go 标准库 (io, bufio, fmt, strings, testing 等)

**内部依赖**:
- `internal/config`: 配置管理(HTTP 端口等)
- `internal/app/http_server`: HTTP 服务器框架
- `main.go`: 主菜单集成

### 集成点

**需要修改的现有文件**:
1. `main.go`: 在 `NewApp()` 的 menu map 中添加 Constants 菜单项
2. `main_test.go`: 添加 Constants 菜单选项测试
3. `internal/app/http_server/router.go`: 注册 Constants 路由
4. `internal/app/http_server/handler/handler.go`: 可能需要添加辅助方法(如果需要)

**新增文件**:
- `internal/app/constants/` 下所有文件(约 26 个文件: 12 个主题 × 2 + constants.go + README.md)
- `internal/app/http_server/handler/constants.go` 和测试文件
- `tests/integration/constants_api_test.go`

---

## Risk Mitigation

| 风险 | 影响 | 概率 | 缓解措施 | 责任人 |
|------|------|------|---------|--------|
| 示例代码维护成本高 | 中 | 高 | 编写自动化编译测试,代码审查时重点检查 | 开发者 |
| 中文注释质量不一致 | 低 | 中 | 制定注释模板,使用统一的术语表 | 开发者 |
| 测试覆盖率难达标 | 高 | 中 | 使用 golden file 测试,关键内容片段匹配 | 开发者 |
| HTTP 性能不达标 | 中 | 低 | 早期性能测试,必要时引入缓存 | 开发者 |
| 与现有模块集成问题 | 高 | 低 | Phase 0.1 深入分析现有架构,遵循既定模式 | 开发者 |

---

## Timeline Estimate

基于任务复杂度的粗略估算(实际时间线将在 tasks.md 中细化):

- **Phase 0 (Research)**: 1-2 天
  - 分析现有架构: 0.5 天
  - 研究 Constants 规范: 0.5 天
  - 设计示例代码: 0.5-1 天
  - 性能和风险分析: 0.5 天

- **Phase 1 (Design)**: 1-1.5 天
  - 数据模型设计: 0.25 天
  - CLI 契约设计: 0.25 天
  - HTTP API 契约设计: 0.5 天
  - 测试策略设计: 0.5 天

- **Phase 2 (Implementation)**: 4-6 天
  - 基础设施: 0.5 天
  - P1 学习内容(6 个基础类型): 2 天
  - P2 学习内容(表达式和类型): 1 天
  - P3 学习内容(转换和函数): 1 天
  - P4 学习内容(iota 和限制): 0.5 天
  - CLI 集成: 0.5 天
  - HTTP 集成: 0.5 天
  - 测试编写: 1-2 天(与实现并行)

- **Phase 3 (Validation)**: 1-1.5 天
  - 功能验收: 0.5 天
  - 质量指标验证: 0.25 天
  - 宪章合规审查: 0.25 天
  - 性能测试和优化: 0.5 天

**总计**: 7-11 天 (取决于示例代码复杂度和测试编写效率)

---

## Next Steps

1. **立即执行**: 运行 `/speckit.tasks` 生成详细的任务分解到 `tasks.md`
2. **开始 Phase 0**: 按照本计划的 Phase 0 章节进行技术调研,输出 `research.md`
3. **完成 Phase 1**: 完成详细设计,输出 `data-model.md` 和 `contracts/` 下的契约文档
4. **执行 Phase 2**: 按照 `tasks.md` 中的任务顺序进行开发
5. **Phase 3 验证**: 完成所有验收测试,确保质量达标
6. **合并代码**: 通过代码审查后合并到主分支
7. **更新文档**: 按照宪章 Principle XIX 更新 `README.md`

---

**计划状态**: ✅ 完成  
**下一步**: 执行 `/speckit.tasks` 生成任务分解
