# Research: Go 类型属性章节学习方案

**Feature**: 018-type-properties-learning  
**Created**: 2026-01-15

## Research Questions

### RQ-001: Properties of types and values 章节结构

**Question**: Go 语言规范中 Properties of types and values 章节包含哪些子章节？

**Decision**: 根据 Go 语言规范，该章节包含7个子章节：
1. Representation of values（值的表示）
2. Underlying types（底层类型）
3. Core types（核心类型）
4. Type identity（类型标识）
5. Assignability（可赋值性）
6. Representability（可表示性）
7. Method sets（方法集）

**Rationale**: 这些子章节完整覆盖了 Go 类型系统的属性与规则，是理解 Go 类型系统高级特性的基础。

**Alternatives Considered**: 无，直接采用规范原文结构。

---

### RQ-002: 现有 Types 章节实现模式

**Question**: 现有 Types 章节的实现模式是什么？

**Decision**: 采用以下实现模式：
- 核心数据结构定义在主文件（`types.go`）
- 每个子主题对应独立的 `.go` 文件
- 使用 `registerXxx()` 函数在 `init()` 中注册内容
- CLI 菜单在 `cli/menu.go`，HTTP handler 在 `http/handlers.go`
- 提供 `AllTopics()`, `LoadContent()`, `LoadQuiz()`, `EvaluateQuiz()`, `SearchReferences()` 等核心函数

**Rationale**: 沿用现有成熟模式可确保一致性，减少学习成本和维护复杂度。

**Alternatives Considered**: 无需改变架构。

---

### RQ-003: 值的表示规则

**Question**: 值的表示（Representation of values）的核心规则是什么？

**Decision**: 核心规则包括：
1. **自包含值**（self-contained）：预声明类型、数组、结构体的值包含完整数据副本
2. **引用值**：指针、函数、切片、map、channel 包含对底层数据的引用
3. **共享底层数据**：多个值共享底层数据时，修改一个可能影响另一个
4. **nil 值**：指针、函数、切片、map、channel、接口的零值为 nil

**Rationale**: 这是理解 Go 值语义的基础。

---

### RQ-004: 底层类型推导规则

**Question**: 底层类型（Underlying types）的推导规则是什么？

**Decision**: 核心规则：
1. 预声明类型（bool, int, string 等）的底层类型是它自身
2. 类型字面量的底层类型是它自身
3. 类型别名（`type A = B`）的底层类型是 B 的底层类型
4. 类型定义（`type A B`）的底层类型是 B 的底层类型
5. 类型参数的底层类型是其类型约束的底层类型（总是接口）

**Rationale**: 底层类型用于判断类型标识和可赋值性。

---

### RQ-005: 核心类型规则

**Question**: 核心类型（Core types）的判定规则是什么？

**Decision**: 核心规则：
1. 非接口类型的核心类型等于其底层类型
2. 接口有核心类型的条件：
   - 类型集中所有类型有相同的底层类型 U，核心类型为 U
   - 类型集仅包含相同元素类型的 channel，核心类型为 chan E（或带方向）
3. 核心类型不是定义类型、类型参数或接口类型
4. **bytestring**：当接口类型集仅包含 `[]byte` 和 `string` 作为底层类型时

**Rationale**: 核心类型用于泛型和某些操作的类型推断。

---

### RQ-006: 类型标识规则

**Question**: 类型标识（Type identity）的判定规则是什么？

**Decision**: 核心规则：
1. 命名类型始终与其他类型不同
2. 两类型相同当底层类型字面量结构等价
3. 具体的结构等价规则：
   - 数组：相同元素类型和长度
   - 切片：相同元素类型
   - 结构体：相同字段序列（名称、类型、标签、嵌入）
   - 指针：相同基础类型
   - 函数：相同参数/返回值类型和可变参数性
   - 接口：定义相同的类型集
   - Map：相同键和值类型
   - Channel：相同元素类型和方向
   - 实例化类型：相同定义类型和类型参数

**Rationale**: 类型标识是赋值、比较和接口实现判定的基础。

---

### RQ-007: 可赋值性规则

**Question**: 可赋值性（Assignability）的判定条件是什么？

**Decision**: 标准条件（6个）：
1. V 和 T 类型相同
2. V 和 T 有相同底层类型，且至少一个不是命名类型（非类型参数）
3. V 和 T 是 channel 类型，V 是双向 channel，至少一个不是命名类型
4. T 是接口类型（非类型参数），x 实现 T
5. x 是 nil，T 是指针/函数/切片/map/channel/接口类型（非类型参数）
6. x 是无类型常量，可表示为 T 类型的值

涉及类型参数的额外条件（3个）：
1. x 是 nil，T 是类型参数，x 可赋值给 T 类型集中的每个类型
2. V 不是命名类型，T 是类型参数，x 可赋值给 T 类型集中的每个类型
3. V 是类型参数，T 不是命名类型，V 类型集中每个类型的值可赋值给 T

**Rationale**: 可赋值性是变量赋值和函数调用参数传递的核心规则。

---

### RQ-008: 可表示性规则

**Question**: 可表示性（Representability）的判定规则是什么？

**Decision**: 常量 x 可表示为类型 T 的值当：
1. x 在 T 决定的值集合中
2. T 是浮点类型，x 可舍入到 T 的精度且不溢出（IEEE 754 round-to-even）
3. T 是复数类型，x 的实部和虚部可表示为 T 的组件类型

类型参数：x 可表示为 T 类型集中每个类型的值

**Rationale**: 主要用于常量赋值的合法性判定。

---

### RQ-009: 方法集规则

**Question**: 方法集（Method sets）的规则是什么？

**Decision**: 核心规则：
1. 定义类型 T 的方法集 = 接收者为 T 的所有方法
2. 指针类型 *T 的方法集 = 接收者为 *T 或 T 的所有方法
3. 接口类型的方法集 = 类型集中所有类型方法集的交集
4. 嵌入字段规则适用于结构体和其指针
5. 方法集中每个方法必须有唯一的非空白方法名

**Rationale**: 方法集决定了类型是否实现某接口。

## Summary

研究完成，所有关键概念已明确。实现将：
- 为每个子章节创建独立的 `.go` 文件
- 遵循 Types 章节的实现模式
- 提供详细的规则、示例和测验
- 支持 CLI 和 HTTP 双模式访问
