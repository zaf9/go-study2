# Feature Specification: Go Constants 学习包

<!--
  NOTE: Per the constitution, all user-facing documentation and code comments
  must be in Chinese.
-->

**Feature Branch**: `004-constants-learning`  
**Created**: 2025-12-05  
**Status**: Draft  
**Input**: User description: "增加一个Constants的package,用于学习Go语言规范中Constants章节的内容,包括详细的示例和说明"

## User Scenarios & Testing *(mandatory)*

<!--
  NOTE: Per the constitution, all features must achieve at least 80% unit test coverage.
  Ensure acceptance scenarios are comprehensive enough to meet this requirement.
-->

### User Story 1 - 基础常量类型学习 (Priority: P1)

作为Go语言学习者,我希望能够通过交互式命令行和HTTP接口学习Go语言中的基础常量类型(布尔常量、符文常量、整数常量、浮点常量、复数常量、字符串常量),以便理解每种常量类型的特性和使用方法。

**Why this priority**: 这是学习常量的基础,必须首先掌握各种常量类型的定义和基本用法,才能进一步学习常量表达式和类型转换等高级特性。

**Independent Test**: 可以通过运行命令行程序选择"Constants"菜单,然后选择任一子主题(如"Boolean Constants"或"Integer Constants"),验证是否能够正确显示该类型的说明和示例代码。也可以通过HTTP GET请求访问对应的学习内容端点,验证返回的内容是否完整准确。

**Acceptance Scenarios**:

1. **Given** 用户启动命令行学习程序, **When** 选择"Constants"主菜单, **Then** 显示所有常量子主题列表(布尔常量、符文常量、整数常量、浮点常量、复数常量、字符串常量)
2. **Given** 用户在Constants菜单中, **When** 选择"Boolean Constants"子主题, **Then** 显示布尔常量的详细说明和至少3个示例代码
3. **Given** 用户在Constants菜单中, **When** 选择"Rune Constants"子主题, **Then** 显示符文常量的详细说明、字面量表示方法和至少3个示例代码
4. **Given** 用户在Constants菜单中, **When** 选择"Integer Constants"子主题, **Then** 显示整数常量的详细说明、不同进制表示方法和至少5个示例代码
5. **Given** 用户在Constants菜单中, **When** 选择"Floating-point Constants"子主题, **Then** 显示浮点常量的详细说明、科学计数法和至少4个示例代码
6. **Given** 用户在Constants菜单中, **When** 选择"Complex Constants"子主题, **Then** 显示复数常量的详细说明和至少3个示例代码
7. **Given** 用户在Constants菜单中, **When** 选择"String Constants"子主题, **Then** 显示字符串常量的详细说明、原始字符串和解释字符串的区别,以及至少4个示例代码
8. **Given** HTTP服务已启动, **When** 发送GET请求到`/learn/constants/boolean`, **Then** 返回JSON格式的布尔常量学习内容
9. **Given** HTTP服务已启动, **When** 发送GET请求到`/learn/constants/integer`, **Then** 返回JSON格式的整数常量学习内容

---

### User Story 2 - 常量表达式和类型学习 (Priority: P2)

作为Go语言学习者,我希望能够学习常量表达式、类型化常量与无类型化常量的区别、常量的默认类型,以便理解Go语言中常量的类型系统和表达式求值规则。

**Why this priority**: 在掌握基础常量类型后,需要理解常量的类型系统和表达式求值,这是编写正确Go代码的重要基础。

**Independent Test**: 可以通过选择"Constant Expressions"或"Typed and Untyped Constants"子主题,验证是否能够正确显示常量表达式的求值规则、类型推断规则和相关示例。

**Acceptance Scenarios**:

1. **Given** 用户在Constants菜单中, **When** 选择"Constant Expressions"子主题, **Then** 显示常量表达式的定义、求值规则和至少5个示例(包括算术、比较、逻辑运算)
2. **Given** 用户在Constants菜单中, **When** 选择"Typed and Untyped Constants"子主题, **Then** 显示类型化常量和无类型化常量的区别、默认类型规则和至少4个对比示例
3. **Given** 用户在Constants菜单中, **When** 选择"Default Types"子主题, **Then** 显示各种无类型化常量的默认类型映射表和至少3个示例
4. **Given** HTTP服务已启动, **When** 发送GET请求到`/learn/constants/expressions`, **Then** 返回JSON格式的常量表达式学习内容

---

### User Story 3 - 常量转换和内置函数学习 (Priority: P3)

作为Go语言学习者,我希望能够学习常量的类型转换、以及可用于常量的内置函数(如min、max、unsafe.Sizeof、cap、len、real、imag、complex),以便理解如何在不同上下文中使用常量。

**Why this priority**: 这是常量使用的高级特性,在掌握基础类型和类型系统后学习,可以帮助学习者更灵活地使用常量。

**Independent Test**: 可以通过选择"Conversions"或"Built-in Functions"子主题,验证是否能够正确显示常量转换规则和内置函数的使用方法。

**Acceptance Scenarios**:

1. **Given** 用户在Constants菜单中, **When** 选择"Conversions"子主题, **Then** 显示常量类型转换的规则、可表示性要求和至少4个示例(包括成功和失败的转换)
2. **Given** 用户在Constants菜单中, **When** 选择"Built-in Functions"子主题, **Then** 显示可用于常量的内置函数列表、每个函数的说明和至少6个示例(覆盖min、max、unsafe.Sizeof、len、real、imag、complex)
3. **Given** 用户在Constants菜单中, **When** 选择"Representability"子主题, **Then** 显示常量可表示性的定义、溢出处理规则和至少3个示例
4. **Given** HTTP服务已启动, **When** 发送GET请求到`/learn/constants/conversions`, **Then** 返回JSON格式的常量转换学习内容

---

### User Story 4 - 特殊常量和实现限制学习 (Priority: P4)

作为Go语言学习者,我希望能够学习特殊常量(如iota、true、false)的使用方法,以及Go编译器对常量的实现限制,以便全面理解Go语言的常量系统。

**Why this priority**: 这是常量学习的补充内容,包括实用的iota特性和需要了解的实现限制,优先级相对较低但对完整学习很重要。

**Independent Test**: 可以通过选择"Iota"或"Implementation Restrictions"子主题,验证是否能够正确显示iota的使用模式和编译器限制说明。

**Acceptance Scenarios**:

1. **Given** 用户在Constants菜单中, **When** 选择"Iota"子主题, **Then** 显示iota的定义、自增规则和至少5个实用示例(包括枚举、位掩码等常见模式)
2. **Given** 用户在Constants菜单中, **When** 选择"Predeclared Constants"子主题, **Then** 显示预声明常量(true、false)的说明和示例
3. **Given** 用户在Constants菜单中, **When** 选择"Implementation Restrictions"子主题, **Then** 显示编译器对常量的实现限制(整数至少256位、浮点数尾数至少256位、指数至少16位等)和相关说明
4. **Given** HTTP服务已启动, **When** 发送GET请求到`/learn/constants/iota`, **Then** 返回JSON格式的iota学习内容

---

### Edge Cases

- 当用户在Constants菜单中输入无效选项时,系统应提示错误并重新显示菜单
- 当HTTP请求访问不存在的常量子主题时,应返回404错误和友好的错误消息
- 当常量示例代码包含特殊字符(如反引号、引号)时,应正确转义和显示
- 当学习内容文件不存在或无法读取时,应返回明确的错误提示而不是崩溃
- 当用户快速连续发送多个HTTP请求时,系统应能够正确处理并发请求

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须在`internal/learn/constants`包下创建Constants学习模块的目录结构
- **FR-002**: 系统必须为每个一级子主题(布尔常量、符文常量、整数常量、浮点常量、复数常量、字符串常量、常量表达式、类型化/无类型化常量、转换、内置函数、iota、实现限制)创建独立的.go文件
- **FR-003**: 每个.go文件必须包含该主题的详细中文说明(作为注释)和至少3个可运行的示例代码
- **FR-004**: 系统必须提供`Display()`函数,用于在命令行模式下显示所有常量学习内容
- **FR-005**: 系统必须支持通过HTTP GET请求访问每个子主题的学习内容,返回JSON格式数据
- **FR-006**: 所有示例代码必须是可编译和可运行的,并包含预期输出的注释
- **FR-007**: 系统必须在主菜单中添加"Constants"选项,并提供子菜单导航
- **FR-008**: 每个子主题的学习内容必须包含:主题说明、语法规则、使用场景、示例代码、常见错误
- **FR-009**: 系统必须为所有学习内容提供至少80%的单元测试覆盖率
- **FR-010**: 所有用户可见的文档和代码注释必须使用中文(符合项目宪章要求)
- **FR-011**: HTTP接口路径必须遵循RESTful风格,格式为`/learn/constants/{subtopic}`
- **FR-012**: 系统必须处理文件读取错误、无效输入等异常情况,并返回友好的错误消息

### Key Entities *(include if feature involves data)*

- **ConstantTopic**: 表示一个常量学习主题,包含主题名称、说明、示例代码列表、相关链接等属性
- **CodeExample**: 表示一个代码示例,包含示例标题、代码内容、预期输出、说明等属性
- **LearningContent**: 表示完整的学习内容,包含主题信息、详细说明、示例列表、相关主题链接等属性

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 学习者能够在5分钟内通过命令行或HTTP接口访问任意常量子主题的学习内容
- **SC-002**: 每个常量子主题至少包含3个可运行的示例代码,覆盖该主题的主要使用场景
- **SC-003**: 所有示例代码能够成功编译和运行,并产生预期的输出结果
- **SC-004**: 单元测试覆盖率达到80%以上,确保代码质量和正确性
- **SC-005**: HTTP接口响应时间在正常负载下(100并发请求)保持在100ms以内
- **SC-006**: 学习内容的准确性达到100%,所有说明和示例符合Go 1.24语言规范
- **SC-007**: 90%的学习者能够在首次使用时成功导航到目标学习主题并理解示例代码
- **SC-008**: 系统能够处理至少1000个并发HTTP请求而不出现错误或性能下降
