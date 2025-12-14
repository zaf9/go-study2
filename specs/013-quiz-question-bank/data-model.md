# Data Model: 测验题库数据模型

**Branch**: `013-quiz-question-bank` | **Date**: 2025-12-14

## 概述

本文档定义测验题库系统的核心数据模型，包括Go语言实体定义、YAML文件格式规范、数据验证规则和存储组织结构。

## 核心实体定义

### QuizQuestion - 测验题目实体

```go
package quiz

// QuestionType 题目类型枚举
type QuestionType string

const (
	QuestionTypeSingle   QuestionType = "single"   // 单选题
	QuestionTypeMultiple QuestionType = "multiple" // 多选题
)

// Difficulty 题目难度枚举
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"   // 简单：直接考察定义、语法规则
	DifficultyMedium Difficulty = "medium" // 中等：理解应用、辨析区别
	DifficultyHard   Difficulty = "hard"   // 困难：综合运用、边界情况
)

// QuizQuestion 测验题目完整定义
type QuizQuestion struct {
	// ID 题目唯一标识，格式：{topic_prefix}-{chapter}-{sequence}
	// 示例："lexical-comments-001", "const-boolean-015"
	ID string `yaml:"id" json:"id"`

	// Type 题型：single（单选）或 multiple（多选）
	Type QuestionType `yaml:"type" json:"type"`

	// Difficulty 难度级别：easy/medium/hard
	Difficulty Difficulty `yaml:"difficulty" json:"difficulty"`

	// Stem 题干，清晰描述问题
	Stem string `yaml:"stem" json:"stem"`

	// Options 选项列表，格式："A: 选项内容"
	// 单选题：2-4个选项
	// 多选题：3-5个选项
	Options []string `yaml:"options" json:"options"`

	// Answer 正确答案
	// 单选题：单个字母，如 "A"
	// 多选题：字母组合，如 "ACD"
	Answer string `yaml:"answer" json:"answer"`

	// Explanation 答案解析（中文），说明为何答案正确及其他选项为何错误
	Explanation string `yaml:"explanation" json:"explanation"`

	// Topic 所属主题：lexical_elements/constants/variables/types
	Topic string `yaml:"topic" json:"topic"`

	// Chapter 所属章节，如："comments", "boolean", "storage"
	Chapter string `yaml:"chapter" json:"chapter"`
}

// QuizBank 题库文件结构（单个YAML文件的根对象）
type QuizBank struct {
	Questions []QuizQuestion `yaml:"questions" json:"questions"`
}
```

### QuizConfig - 题库配置

```go
package quiz

// QuizConfig 题库全局配置
type QuizConfig struct {
	// DataPath 题库文件根目录
	DataPath string `yaml:"dataPath"`

	// QuestionCount 每次测验抽取的题目数量
	QuestionCount QuestionCountConfig `yaml:"questionCount"`

	// DifficultyDistribution 难度分布百分比
	DifficultyDistribution DifficultyDistributionConfig `yaml:"difficultyDistribution"`

	// LoadTimeout 题库加载超时时间
	LoadTimeout string `yaml:"loadTimeout"`

	// Validation 验证配置
	Validation ValidationConfig `yaml:"validation"`
}

// QuestionCountConfig 题目数量配置
type QuestionCountConfig struct {
	Single   int `yaml:"single"`   // 单选题数量，默认4
	Multiple int `yaml:"multiple"` // 多选题数量，默认4
}

// DifficultyDistributionConfig 难度分布配置
type DifficultyDistributionConfig struct {
	Easy   int `yaml:"easy"`   // 简单题百分比，默认40
	Medium int `yaml:"medium"` // 中等题百分比，默认40
	Hard   int `yaml:"hard"`   // 困难题百分比，默认20
}

// ValidationConfig 验证配置
type ValidationConfig struct {
	StrictMode bool `yaml:"strictMode"` // 严格验证模式，默认true
	FailFast   bool `yaml:"failFast"`   // 遇到首个错误即停止，默认true
}
```

### QuizSelection - 抽题结果

```go
package quiz

// QuizSelection 抽题结果（API返回）
type QuizSelection struct {
	Topic     string         `json:"topic"`     // 主题
	Chapter   string         `json:"chapter"`   // 章节
	Questions []QuizQuestion `json:"questions"` // 抽取的题目列表
	Total     int            `json:"total"`     // 题库总题数
	Selected  int            `json:"selected"`  // 抽取题数
}
```

## YAML文件格式规范

### 文件路径规则

```text
backend/quiz_data/
├── lexical_elements/
│   ├── comments.yaml
│   ├── tokens.yaml
│   ├── semicolons.yaml
│   ├── identifiers.yaml
│   ├── keywords.yaml
│   ├── operators.yaml
│   ├── integers.yaml
│   ├── floats.yaml
│   ├── imaginary.yaml
│   ├── runes.yaml
│   └── strings.yaml
├── constants/
│   ├── boolean.yaml
│   ├── rune.yaml
│   ├── integer.yaml
│   ├── floating_point.yaml
│   ├── complex.yaml
│   ├── string.yaml
│   ├── expressions.yaml
│   ├── typed_untyped.yaml
│   ├── conversions.yaml
│   ├── builtin_functions.yaml
│   ├── iota.yaml
│   └── implementation_restrictions.yaml
├── variables/
│   ├── storage.yaml
│   ├── static.yaml
│   ├── dynamic.yaml
│   └── zero.yaml
└── types/
    ├── boolean.yaml
    ├── numeric.yaml
    ├── string.yaml
    ├── array.yaml
    ├── slice.yaml
    ├── struct.yaml
    ├── pointer.yaml
    ├── function.yaml
    ├── interface_basic.yaml
    ├── interface_embedded.yaml
    ├── interface_general.yaml
    ├── interface_impl.yaml
    ├── map.yaml
    └── channel.yaml
```

### 完整YAML示例

```yaml
# backend/quiz_data/lexical_elements/comments.yaml
questions:
  # 单选题示例
  - id: lexical-comments-001
    type: single
    difficulty: easy
    stem: "Go语言中，以下哪种注释方式是正确的？"
    options:
      - "A: // 这是单行注释"
      - "B: # 这是注释"
      - "C: <!-- 这是注释 -->"
      - "D: -- 这是注释"
    answer: "A"
    explanation: "Go语言支持两种注释方式：单行注释使用//，多行注释使用/* */。选项B是Python风格，C是HTML风格，D是SQL风格，均不适用于Go。"
    topic: "lexical_elements"
    chapter: "comments"

  # 多选题示例
  - id: lexical-comments-002
    type: multiple
    difficulty: medium
    stem: "关于Go语言的注释，以下说法正确的是？（多选）"
    options:
      - "A: 单行注释以//开头"
      - "B: 多行注释可以嵌套"
      - "C: 注释不影响程序执行"
      - "D: 注释可以用于文档生成"
    answer: "ACD"
    explanation: "A正确，单行注释使用//。B错误，Go的多行注释/* */不支持嵌套。C正确，注释在编译时被忽略。D正确，godoc工具可以提取注释生成文档。"
    topic: "lexical_elements"
    chapter: "comments"

  # 困难题示例
  - id: lexical-comments-003
    type: single
    difficulty: hard
    stem: "以下代码片段中，哪个注释会导致编译错误？"
    options:
      - "A: /* 外层注释 /* 内层注释 */ 外层继续 */"
      - "B: // 单行注释 // 另一个单行注释"
      - "C: /* 多行注释\n第二行 */"
      - "D: // 单行注释 /* 内嵌多行开始"
    answer: "A"
    explanation: "A错误，因为Go的多行注释/* */不支持嵌套，第一个*/会关闭外层注释，导致'外层继续'成为非法代码。B、C、D均合法。"
    topic: "lexical_elements"
    chapter: "comments"
```

### 字段约束详解

| 字段 | 类型 | 约束 | 示例 |
|------|------|------|------|
| `id` | string | 必填，格式：`{prefix}-{chapter}-{seq}` | `lexical-comments-001` |
| `type` | string | 必填，枚举：`single`/`multiple` | `single` |
| `difficulty` | string | 必填，枚举：`easy`/`medium`/`hard` | `medium` |
| `stem` | string | 必填，非空，中文或英文 | `"Go语言中..."` |
| `options` | array | 必填，长度2-5，格式：`"X: ..."` | `["A: ...", "B: ..."]` |
| `answer` | string | 必填，单选：1字母，多选：2-4字母 | `"A"` 或 `"ACD"` |
| `explanation` | string | 必填，非空，中文 | `"A正确因为..."` |
| `topic` | string | 必填，匹配文件路径 | `"lexical_elements"` |
| `chapter` | string | 必填，匹配文件名 | `"comments"` |

## 题目ID命名规则

### 格式定义

```text
{topic_prefix}-{chapter}-{sequence}
```

### 主题前缀映射

| 主题 | 前缀 | 示例章节 |
|------|------|----------|
| `lexical_elements` | `lexical` | `lexical-comments-001` |
| `constants` | `const` | `const-boolean-015` |
| `variables` | `var` | `var-storage-003` |
| `types` | `type` | `type-array-022` |

### 序号规则

- 格式：三位数字（001-050）
- 范围：001-050（每章节最多50题）
- 递增：按题目添加顺序递增
- 示例：001, 002, 003, ..., 050

### 完整示例

```text
lexical-comments-001        # lexical_elements/comments 第1题
lexical-comments-002        # lexical_elements/comments 第2题
const-boolean-001           # constants/boolean 第1题
const-iota-015              # constants/iota 第15题
var-storage-001             # variables/storage 第1题
type-interface_basic-030    # types/interface_basic 第30题
```

## 数据验证规则

### 必填字段验证

所有9个字段均为必填：

```go
func (q *QuizQuestion) Validate() error {
	if q.ID == "" {
		return errors.New("题目ID不能为空")
	}
	if q.Type == "" {
		return errors.New("题型不能为空")
	}
	if q.Difficulty == "" {
		return errors.New("难度不能为空")
	}
	if q.Stem == "" {
		return errors.New("题干不能为空")
	}
	if len(q.Options) == 0 {
		return errors.New("选项不能为空")
	}
	if q.Answer == "" {
		return errors.New("答案不能为空")
	}
	if q.Explanation == "" {
		return errors.New("解析不能为空")
	}
	if q.Topic == "" {
		return errors.New("主题不能为空")
	}
	if q.Chapter == "" {
		return errors.New("章节不能为空")
	}
	return nil
}
```

### 枚举值验证

```go
func (q *QuizQuestion) ValidateEnums() error {
	// 验证题型
	if q.Type != QuestionTypeSingle && q.Type != QuestionTypeMultiple {
		return fmt.Errorf("无效的题型: %s, 仅支持 single/multiple", q.Type)
	}
	
	// 验证难度
	if q.Difficulty != DifficultyEasy && 
	   q.Difficulty != DifficultyMedium && 
	   q.Difficulty != DifficultyHard {
		return fmt.Errorf("无效的难度: %s, 仅支持 easy/medium/hard", q.Difficulty)
	}
	
	return nil
}
```

### 选项格式验证

```go
func (q *QuizQuestion) ValidateOptions() error {
	// 验证选项数量
	optCount := len(q.Options)
	if optCount < 2 || optCount > 5 {
		return fmt.Errorf("选项数量必须在2-5之间，当前: %d", optCount)
	}
	
	// 单选题：2-4个选项
	if q.Type == QuestionTypeSingle && (optCount < 2 || optCount > 4) {
		return fmt.Errorf("单选题选项数量必须在2-4之间，当前: %d", optCount)
	}
	
	// 多选题：3-5个选项
	if q.Type == QuestionTypeMultiple && (optCount < 3 || optCount > 5) {
		return fmt.Errorf("多选题选项数量必须在3-5之间，当前: %d", optCount)
	}
	
	// 验证选项格式："A: ...", "B: ...", ...
	optionLabels := make(map[string]bool)
	for _, opt := range q.Options {
		if len(opt) < 3 || opt[1] != ':' {
			return fmt.Errorf("选项格式错误，必须为'X: 内容'，当前: %s", opt)
		}
		label := string(opt[0])
		if label < "A" || label > "E" {
			return fmt.Errorf("选项标签必须为A-E，当前: %s", label)
		}
		if optionLabels[label] {
			return fmt.Errorf("选项标签重复: %s", label)
		}
		optionLabels[label] = true
	}
	
	return nil
}
```

### 答案格式验证

```go
func (q *QuizQuestion) ValidateAnswer() error {
	if q.Answer == "" {
		return errors.New("答案不能为空")
	}
	
	// 提取选项标签
	validLabels := make(map[string]bool)
	for _, opt := range q.Options {
		if len(opt) > 0 {
			validLabels[string(opt[0])] = true
		}
	}
	
	// 单选题验证
	if q.Type == QuestionTypeSingle {
		if len(q.Answer) != 1 {
			return fmt.Errorf("单选题答案必须为单个字母，当前: %s", q.Answer)
		}
		if !validLabels[q.Answer] {
			return fmt.Errorf("答案 %s 不在选项范围内", q.Answer)
		}
	}
	
	// 多选题验证
	if q.Type == QuestionTypeMultiple {
		if len(q.Answer) < 2 {
			return fmt.Errorf("多选题答案必须至少包含2个字母，当前: %s", q.Answer)
		}
		if len(q.Answer) > 4 {
			return fmt.Errorf("多选题答案不能超过4个字母，当前: %s", q.Answer)
		}
		
		// 验证每个字母都在选项中
		seen := make(map[rune]bool)
		for _, ch := range q.Answer {
			label := string(ch)
			if !validLabels[label] {
				return fmt.Errorf("答案字母 %s 不在选项范围内", label)
			}
			if seen[ch] {
				return fmt.Errorf("答案字母重复: %s", label)
			}
			seen[ch] = true
		}
	}
	
	return nil
}
```

### 路径一致性验证

```go
func (q *QuizQuestion) ValidatePath(fileTopic, fileChapter string) error {
	if q.Topic != fileTopic {
		return fmt.Errorf("题目主题 %s 与文件路径 %s 不一致", q.Topic, fileTopic)
	}
	if q.Chapter != fileChapter {
		return fmt.Errorf("题目章节 %s 与文件名 %s 不一致", q.Chapter, fileChapter)
	}
	return nil
}
```

### ID唯一性验证

```go
func ValidateUniqueIDs(questions []QuizQuestion) error {
	idMap := make(map[string]int)
	for i, q := range questions {
		if existingIndex, exists := idMap[q.ID]; exists {
			return fmt.Errorf("题目ID重复: %s (索引 %d 和 %d)", q.ID, existingIndex, i)
		}
		idMap[q.ID] = i
	}
	return nil
}
```

## 存储组织结构

### 目录结构

```text
backend/
├── quiz_data/                    # 题库数据根目录
│   ├── README.md                 # 题库文件组织说明
│   ├── lexical_elements/         # 词法元素主题
│   │   ├── comments.yaml         # 11个章节YAML文件
│   │   ├── tokens.yaml
│   │   └── ...
│   ├── constants/                # 常量主题
│   │   ├── boolean.yaml          # 12个章节YAML文件
│   │   ├── rune.yaml
│   │   └── ...
│   ├── variables/                # 变量主题
│   │   ├── storage.yaml          # 4个章节YAML文件
│   │   └── ...
│   └── types/                    # 类型主题
│       ├── boolean.yaml          # 14个章节YAML文件
│       └── ...
└── internal/
    └── domain/
        └── quiz/                 # 题库管理包
            ├── entity.go         # 实体定义（QuizQuestion等）
            ├── loader.go         # 题库加载器
            ├── validator.go      # 题库验证器
            ├── selector.go       # 抽题选择器
            └── README.md         # 包说明文档
```

### 内存索引结构

```go
package quiz

// QuizRepository 题库仓储（内存缓存）
type QuizRepository struct {
	// banks 按主题和章节索引的题库
	// map[topic][chapter][]QuizQuestion
	banks map[string]map[string][]QuizQuestion
	
	// loaded 标记题库是否已加载
	loaded bool
}

// 使用示例
func (r *QuizRepository) GetQuestions(topic, chapter string) ([]QuizQuestion, error) {
	if !r.loaded {
		return nil, errors.New("题库未加载")
	}
	
	topicBank, exists := r.banks[topic]
	if !exists {
		return nil, fmt.Errorf("主题不存在: %s", topic)
	}
	
	questions, exists := topicBank[chapter]
	if !exists {
		return nil, fmt.Errorf("章节不存在: %s/%s", topic, chapter)
	}
	
	return questions, nil
}
```

## 配置示例

### config.yaml 配置项

```yaml
# backend/configs/config.yaml
quiz:
  # 题库文件根目录（相对于backend目录）
  dataPath: "quiz_data"
  
  # 每次测验抽取的题目数量
  questionCount:
    single: 4      # 单选题数量
    multiple: 4    # 多选题数量
  
  # 难度分布百分比（总和应为100）
  difficultyDistribution:
    easy: 40       # 简单题占比40%
    medium: 40     # 中等题占比40%
    hard: 20       # 困难题占比20%
  
  # 题库加载超时时间
  loadTimeout: 5s
  
  # 验证配置
  validation:
    strictMode: true   # 严格验证模式
    failFast: true     # 遇到首个错误即停止
```

## 数据流程图

```text
应用启动
    ↓
加载配置 (config.yaml)
    ↓
遍历题库目录 (quiz_data/*/*.yaml)
    ↓
解析YAML文件 → QuizBank
    ↓
验证题目完整性 → Validator
    ↓
构建内存索引 → Repository
    ↓
【如有错误，Fail-Fast停止启动】
    ↓
题库就绪，等待抽题请求
    ↓
用户发起测验 (topic, chapter)
    ↓
从Repository查询题目列表
    ↓
Selector按配置随机抽题
    ↓
返回QuizSelection给前端
```

## 实体关系图

```text
QuizConfig (配置)
    |
    | 控制
    ↓
Selector (选择器)
    |
    | 读取
    ↓
QuizRepository (仓储)
    |
    | 管理
    ↓
QuizBank (题库文件)
    |
    | 包含
    ↓
QuizQuestion (题目实体)
```

## 扩展性考虑

### 未来可能扩展

1. **新增题型**: 如填空题、判断题、代码题
   - 扩展 `QuestionType` 枚举
   - 调整 `QuizQuestion` 结构（可能需要新字段）

2. **个性化难度**: 根据用户水平调整难度分布
   - 在 `QuizConfig` 中增加用户级配置
   - Selector支持动态难度调整

3. **题目标签**: 增加知识点标签，支持按标签抽题
   - 在 `QuizQuestion` 中增加 `Tags []string`
   - 扩展Selector支持标签过滤

4. **题目版本**: 支持题目更新和版本管理
   - 增加 `Version` 字段
   - 题库文件支持多版本共存

### 设计原则保障

- **单一职责**: 每个实体职责明确，QuizQuestion只表示题目，Selector只负责抽题
- **开闭原则**: 通过接口扩展，无需修改现有代码
- **YAGNI**: 当前不实现未确定需求，保持简单

---

**Data Model Status**: ✅ Complete  
**Next**: 生成 contracts/yaml-schema.md
