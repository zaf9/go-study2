# Data Model: Go Constants 学习包

**Feature**: 004-constants-learning  
**Date**: 2025-12-05  
**Phase**: 1 - Design & Contracts

## 核心数据结构

### TopicContent - 学习主题内容

表示一个学习主题的完整内容结构。

```go
// TopicContent 表示一个学习主题的完整内容
type TopicContent struct {
    Title        string        `json:"title"`         // 主题标题(中文)
    Description  string        `json:"description"`   // 主题说明
    Syntax       string        `json:"syntax"`        // 语法规则
    Examples     []CodeExample `json:"examples"`      // 示例代码列表
    CommonErrors string        `json:"common_errors"` // 常见错误说明
}
```

**字段说明**:
- `Title`: 主题标题,如 "Boolean Constants (布尔常量)"
- `Description`: 详细的概念说明,解释主题的核心概念
- `Syntax`: 语法规则说明,如 `const name [type] = value`
- `Examples`: 示例代码列表,至少 3 个
- `CommonErrors`: 常见错误和注意事项

### CodeExample - 代码示例

表示一个独立的代码示例。

```go
// CodeExample 表示一个代码示例
type CodeExample struct {
    Title       string `json:"title"`       // 示例标题
    Code        string `json:"code"`        // 示例代码
    Output      string `json:"output"`      // 预期输出
    Explanation string `json:"explanation"` // 说明
}
```

**字段说明**:
- `Title`: 示例标题,如 "基本布尔常量声明"
- `Code`: 完整的示例代码(包含 package 和 main)
- `Output`: 预期输出或注释说明
- `Explanation`: 示例的详细解释

### MenuOption - 菜单选项

表示命令行菜单的一个选项。

```go
// MenuOption 表示菜单选项
type MenuOption struct {
    Key         string // 选项键(如 "0", "1")
    Description string // 选项描述
    Action      func() // 执行函数
}
```

**使用场景**: CLI 菜单导航

### SubtopicInfo - 子主题信息

用于 HTTP API 返回子主题列表。

```go
// SubtopicInfo 子主题信息
type SubtopicInfo struct {
    Key  string `json:"key"`  // 子主题键(如 "boolean")
    Name string `json:"name"` // 子主题名称(如 "Boolean Constants (布尔常量)")
}
```

## 数据流

### CLI 模式数据流

```
用户输入 → DisplayMenu() → 选择主题 → Display{Topic}() → 输出到 stdout
```

### HTTP 模式数据流

```
HTTP 请求 → Handler → 构造 TopicContent → JSON 序列化 → HTTP 响应
```

## 验证规则

1. **Title**: 非空,包含中英文名称
2. **Examples**: 至少 3 个元素
3. **Code**: 可编译的 Go 代码
4. **Output**: 非空,说明预期结果

## 扩展性考虑

- 预留 `RelatedTopics []string` 字段用于主题关联
- 预留 `References []string` 字段用于外部参考链接
- JSON 标签支持未来 API 版本升级
