# Data Model: Go 类型属性章节学习方案

**Feature**: 018-type-properties-learning  
**Created**: 2026-01-15

## Core Entities

### Topic (子主题枚举)

```go
type Topic string

const (
    TopicRepresentation  Topic = "representation"   // 值的表示
    TopicUnderlyingType  Topic = "underlying_type"  // 底层类型
    TopicCoreType        Topic = "core_type"        // 核心类型
    TopicTypeIdentity    Topic = "type_identity"    // 类型标识
    TopicAssignability   Topic = "assignability"    // 可赋值性
    TopicRepresentability Topic = "representability" // 可表示性
    TopicMethodSet       Topic = "method_set"       // 方法集
)
```

### PropertyConcept (类型属性概念)

描述单个类型属性子主题的元信息。

| 字段 | 类型 | 描述 |
|------|------|------|
| ID | string | 唯一标识符（与 Topic 对应） |
| Category | string | 类别（basic/advanced） |
| Title | string | 中文标题 |
| Summary | string | 摘要说明 |
| GoVersion | string | 适用的 Go 版本 |
| Rules | []string | 关键规则列表 |
| Keywords | []string | 关键词（用于搜索） |
| PrintableOutline | []string | 可打印提纲 |

### PropertyRule (类型属性规则)

描述具体的规则或约束。

| 字段 | 类型 | 描述 |
|------|------|------|
| RuleID | string | 规则唯一标识（如 "PR-ASSIGN-001"） |
| ConceptID | string | 关联的概念ID |
| RuleType | string | 规则类型（identity/assignability/representability 等） |
| Description | string | 规则描述（中文） |
| References | []string | 参考文献/规范章节 |
| Severity | string | 严重程度（info/warning/error） |

### ExampleCase (示例用例)

用于演示规则的正例或反例。

| 字段 | 类型 | 描述 |
|------|------|------|
| ID | string | 示例唯一标识 |
| ConceptID | string | 关联的概念ID |
| Title | string | 示例标题（中文） |
| Code | string | 示例代码 |
| ExpectedOutput | string | 期望输出 |
| IsValid | bool | 是否为合法示例 |
| RuleRef | string | 引用的规则ID |
| Notes | []string | 注释说明 |

### QuizItem (测验题目)

| 字段 | 类型 | 描述 |
|------|------|------|
| ID | string | 题目唯一标识 |
| ConceptID | string | 关联的概念ID |
| Stem | string | 题干（中文） |
| Options | []string | 选项列表 |
| Answer | string | 正确答案（A/B/C/D） |
| Explanation | string | 答案解析（中文） |
| RuleRef | string | 引用的规则ID |
| Difficulty | string | 难度等级（easy/medium/hard） |

### QuizAnswerFeedback (测验答题反馈)

| 字段 | 类型 | 描述 |
|------|------|------|
| ID | string | 题目ID |
| Correct | bool | 是否正确 |
| Answer | string | 正确答案 |
| Explanation | string | 解析 |
| RuleRef | string | 规则引用 |

### QuizResult (测验结果)

| 字段 | 类型 | 描述 |
|------|------|------|
| Score | int | 得分 |
| Total | int | 总题数 |
| Details | []QuizAnswerFeedback | 各题详情 |

### ReferenceIndex (搜索索引)

| 字段 | 类型 | 描述 |
|------|------|------|
| Keyword | string | 关键词 |
| ConceptID | string | 关联概念ID |
| Summary | string | 摘要 |
| PositiveExampleID | string | 正例ID |
| NegativeExampleID | string | 反例ID |
| Anchors | map[string]string | 跳转锚点（http/cli） |

### TopicContent (聚合内容)

聚合单个子主题的全部素材，用于内容加载。

| 字段 | 类型 | 描述 |
|------|------|------|
| Concept | PropertyConcept | 概念信息 |
| Rules | []PropertyRule | 规则列表 |
| Examples | []ExampleCase | 示例列表 |
| QuizItems | []QuizItem | 测验题目 |
| References | []ReferenceIndex | 搜索索引 |

### LearningProgress (学习进度)

| 字段 | 类型 | 描述 |
|------|------|------|
| UserID | string | 用户标识 |
| CompletedConcepts | []string | 已完成的概念ID列表 |
| LastVisited | string | 最后访问的概念ID |
| QuizScores | map[string]QuizResult | 各主题的测验得分 |

## Entity Relationships

```
PropertyConcept 1---* PropertyRule
PropertyConcept 1---* ExampleCase
PropertyConcept 1---* QuizItem
PropertyConcept 1---* ReferenceIndex
LearningProgress *---* PropertyConcept (通过 CompletedConcepts)
LearningProgress *---* QuizResult (通过 QuizScores)
```

## Registry Structure

使用 map 结构在内存中管理所有内容：

```go
var (
    conceptRegistry   = map[Topic]PropertyConcept{}
    ruleRegistry      = map[Topic][]PropertyRule{}
    exampleRegistry   = map[Topic][]ExampleCase{}
    quizRegistry      = map[Topic][]QuizItem{}
    referenceRegistry = map[string]ReferenceIndex{}  // key: 小写关键词
)
```

## Core Functions

| 函数 | 描述 |
|------|------|
| `AllTopics() []Topic` | 返回所有子主题（按教学顺序） |
| `NormalizeTopic(raw string) Topic` | 归一化子主题名称 |
| `IsSupportedTopic(topic Topic) bool` | 判断是否支持该子主题 |
| `RegisterContent(topic Topic, content TopicContent) error` | 注册子主题内容 |
| `LoadContent(topic Topic) (TopicContent, error)` | 加载子主题内容 |
| `LoadQuiz(topic Topic) ([]QuizItem, error)` | 加载子主题测验 |
| `EvaluateQuiz(topic Topic, answers map[string]string) (QuizResult, error)` | 评估测验答案 |
| `SearchReferences(keyword string) ([]ReferenceIndex, error)` | 搜索关键词索引 |
| `GetOverview() PropertyOverview` | 获取章节概览 |
| `RenderPrintableOutline() string` | 渲染可打印提纲 |
