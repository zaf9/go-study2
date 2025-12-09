# Data Model: Variables章节学习完成

## Entities

- 变量 (Variable)
  - Fields: name/标识, staticType, valueRepresentation, storageKind(declaration|new|composite_addr), addressable(bool)
  - Relationships: 依赖静态类型；对接口变量可关联动态类型
  - Validation Rules: staticType 必填；storageKind 必在定义范围内

- 静态类型 (StaticType)
  - Fields: typeName, origin(declaration|new|composite|element), assignableTo(list)
  - Relationships: 约束变量可赋值性；被变量引用
  - Validation Rules: typeName 必填

- 动态类型 (DynamicType)
  - Fields: concreteType(non-interface), valueState(nil|non-nil)
  - Relationships: 仅接口变量持有；随赋值变化
  - Validation Rules: 当 valueState 为 nil 时 dynamicType 可为空，否则需与存入值一致

- 零值 (ZeroValue)
  - Fields: typeName, valueExample, isCompositeElement(bool)
  - Relationships: 关联变量或结构化元素的类型信息
  - Validation Rules: typeName 必填；valueExample 与类型匹配

- 测验项 (QuizItem)
  - Fields: id, topic(storage|staticType|dynamicType|zeroValue), stem, options[], answer, explanation
  - Relationships: 覆盖功能需求中的三类知识点；可用于CLI/HTTP两端
  - Validation Rules: 至少2个选项且唯一正确答案；explanation 必填