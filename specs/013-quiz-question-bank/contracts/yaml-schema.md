# YAML Schema: 题库文件格式规范

**Branch**: `013-quiz-question-bank` | **Date**: 2025-12-14

## 概述

本文档定义题库YAML文件的详细格式规范，包括字段定义、类型约束、验证规则和示例。此规范用于指导题目生成、人工审核和自动验证。

## 文件结构

### 根对象

```yaml
questions:    # 根级唯一字段，类型：array
  - {...}     # QuizQuestion对象数组
  - {...}
```

### QuizQuestion对象

```yaml
- id: string                    # 必填，题目唯一标识
  type: "single" | "multiple"   # 必填，题型枚举
  difficulty: "easy" | "medium" | "hard"  # 必填，难度枚举
  stem: string                  # 必填，题干
  options: array<string>        # 必填，选项列表
  answer: string                # 必填，正确答案
  explanation: string           # 必填，答案解析
  topic: string                 # 必填，所属主题
  chapter: string               # 必填，所属章节
```

## 字段详细规范

### 1. id (题目ID)

**类型**: `string`  
**必填**: ✅  
**格式**: `{topic_prefix}-{chapter}-{sequence}`  
**长度**: 10-30字符  

**规则**:
- `topic_prefix`: 主题简写（lexical/const/var/type）
- `chapter`: 章节名（英文小写，下划线分隔）
- `sequence`: 三位数字（001-050）

**示例**:
```yaml
id: "lexical-comments-001"      # ✅ 正确
id: "const-boolean-015"         # ✅ 正确
id: "type-interface_basic-030"  # ✅ 正确（下划线在章节名中合法）

id: "comments-001"              # ❌ 错误：缺少主题前缀
id: "lexical-comments-1"        # ❌ 错误：序号不是三位数
id: "LEXICAL-COMMENTS-001"      # ❌ 错误：大写不符合规范
```

**验证规则**:
```go
func ValidateID(id, topic, chapter string) error {
	pattern := fmt.Sprintf(`^%s-%s-\d{3}$`, getTopicPrefix(topic), chapter)
	matched, _ := regexp.MatchString(pattern, id)
	if !matched {
		return fmt.Errorf("ID格式错误: %s, 期望格式: %s", id, pattern)
	}
	return nil
}
```

---

### 2. type (题型)

**类型**: `string`  
**必填**: ✅  
**枚举值**: `"single"` | `"multiple"`  

**含义**:
- `single`: 单选题（只有一个正确答案）
- `multiple`: 多选题（有2个或更多正确答案）

**示例**:
```yaml
type: "single"      # ✅ 正确
type: "multiple"    # ✅ 正确

type: "Single"      # ❌ 错误：大小写敏感
type: "choice"      # ❌ 错误：不在枚举范围
type: ""            # ❌ 错误：不能为空
```

**验证规则**:
```go
func ValidateType(t string) error {
	if t != "single" && t != "multiple" {
		return fmt.Errorf("无效的题型: %s, 必须为 single 或 multiple", t)
	}
	return nil
}
```

---

### 3. difficulty (难度)

**类型**: `string`  
**必填**: ✅  
**枚举值**: `"easy"` | `"medium"` | `"hard"`  

**含义**:
- `easy`: 简单题（直接考察定义、基本规则，不需推理）
- `medium`: 中等题（需要理解规则并应用，辨析相似概念）
- `hard`: 困难题（需要综合运用多个规则，涉及边界情况）

**示例**:
```yaml
difficulty: "easy"      # ✅ 正确
difficulty: "medium"    # ✅ 正确
difficulty: "hard"      # ✅ 正确

difficulty: "Easy"      # ❌ 错误：大小写敏感
difficulty: "normal"    # ❌ 错误：不在枚举范围
difficulty: "1"         # ❌ 错误：不能使用数字
```

**分布建议**:
- 每章节：easy 40%, medium 40%, hard 20%
- 30题章节：12简单 + 12中等 + 6困难
- 50题章节：20简单 + 20中等 + 10困难

**验证规则**:
```go
func ValidateDifficulty(d string) error {
	if d != "easy" && d != "medium" && d != "hard" {
		return fmt.Errorf("无效的难度: %s, 必须为 easy, medium 或 hard", d)
	}
	return nil
}
```

---

### 4. stem (题干)

**类型**: `string`  
**必填**: ✅  
**长度**: 10-500字符  
**语言**: 中文或英文（优先中文）  

**规则**:
- 描述清晰，语法正确
- 问题明确，无歧义
- 多选题应标注"（多选）"
- 避免使用"以下哪个是错误的"等双重否定

**示例**:
```yaml
stem: "Go语言中，以下哪种注释方式是正确的？"  # ✅ 正确
stem: "关于Go语言的注释，以下说法正确的是？（多选）"  # ✅ 正确（多选题标注）

stem: "注释"                                  # ❌ 错误：太简短，不成问题
stem: ""                                      # ❌ 错误：不能为空
stem: "以下哪个不是不正确的？"                # ❌ 错误：双重否定，难理解
```

**最佳实践**:
```yaml
# 好的题干
stem: "在Go语言中，iota的值在什么情况下会重置为0？"
stem: "以下代码片段中，哪些变量会发生逃逸分析？（多选）"

# 避免的题干
stem: "iota?"                           # 太简短
stem: "哪个不是错的？"                  # 双重否定
stem: "Go有注释吗？"                    # 问题太宽泛
```

**验证规则**:
```go
func ValidateStem(stem string, qType QuestionType) error {
	if len(stem) < 10 || len(stem) > 500 {
		return fmt.Errorf("题干长度必须在10-500字符之间，当前: %d", len(stem))
	}
	if qType == QuestionTypeMultiple && !strings.Contains(stem, "多选") {
		return errors.New("多选题题干应包含'（多选）'标注")
	}
	return nil
}
```

---

### 5. options (选项列表)

**类型**: `array<string>`  
**必填**: ✅  
**长度**: 单选题2-4个，多选题3-5个  
**格式**: `"X: 选项内容"`（X为A-E）  

**规则**:
- 每个选项必须以 `字母: ` 开头
- 字母必须连续（A, B, C...，不能跳号）
- 字母必须大写
- 选项内容应简洁明确
- 选项之间应互斥（不重叠）

**示例**:
```yaml
# 单选题（2-4个选项）
options:
  - "A: // 这是单行注释"
  - "B: # 这是注释"
  - "C: <!-- 这是注释 -->"
  - "D: -- 这是注释"

# 多选题（3-5个选项）
options:
  - "A: 单行注释以//开头"
  - "B: 多行注释可以嵌套"
  - "C: 注释不影响程序执行"
  - "D: 注释可以用于文档生成"

# 错误示例
options:
  - "A: 选项A"
  - "C: 选项C"              # ❌ 错误：跳过了B

options:
  - "a: 选项a"              # ❌ 错误：小写字母

options:
  - "A选项A"                # ❌ 错误：缺少冒号

options:
  - "A: 选项A"              # ❌ 错误：单选题只有1个选项（至少2个）
```

**验证规则**:
```go
func ValidateOptions(options []string, qType QuestionType) error {
	// 验证数量
	count := len(options)
	if qType == QuestionTypeSingle && (count < 2 || count > 4) {
		return fmt.Errorf("单选题选项数量必须在2-4之间，当前: %d", count)
	}
	if qType == QuestionTypeMultiple && (count < 3 || count > 5) {
		return fmt.Errorf("多选题选项数量必须在3-5之间，当前: %d", count)
	}
	
	// 验证格式和连续性
	expectedLabel := 'A'
	for i, opt := range options {
		if len(opt) < 3 || opt[1] != ':' {
			return fmt.Errorf("选项%d格式错误，必须为'X: 内容'，当前: %s", i, opt)
		}
		if rune(opt[0]) != expectedLabel {
			return fmt.Errorf("选项标签不连续，期望 %c，实际 %c", expectedLabel, opt[0])
		}
		expectedLabel++
	}
	
	return nil
}
```

---

### 6. answer (正确答案)

**类型**: `string`  
**必填**: ✅  
**格式**: 
- 单选题：单个大写字母（`"A"`, `"B"`, `"C"`, `"D"`, `"E"`）
- 多选题：2-4个大写字母组合（`"AB"`, `"ACD"`, `"BDE"`）

**规则**:
- 字母必须大写
- 字母必须在选项范围内
- 多选题答案字母必须升序排列
- 多选题至少包含2个正确答案

**示例**:
```yaml
# 单选题
answer: "A"      # ✅ 正确
answer: "C"      # ✅ 正确

# 多选题
answer: "ACD"    # ✅ 正确（3个答案，升序）
answer: "AB"     # ✅ 正确（2个答案，升序）
answer: "BDE"    # ✅ 正确（3个答案，升序）

# 错误示例
answer: "a"      # ❌ 错误：小写
answer: "AC"     # ❌ 错误：单选题不能有多个答案
answer: "A"      # ❌ 错误：多选题至少2个答案
answer: "CAD"    # ❌ 错误：未升序排列（应为ACD）
answer: "F"      # ❌ 错误：超出选项范围（只有A-E）
answer: ""       # ❌ 错误：不能为空
```

**验证规则**:
```go
func ValidateAnswer(answer string, qType QuestionType, options []string) error {
	if answer == "" {
		return errors.New("答案不能为空")
	}
	
	// 提取有效选项标签
	validLabels := make(map[byte]bool)
	for _, opt := range options {
		if len(opt) > 0 {
			validLabels[opt[0]] = true
		}
	}
	
	// 单选题验证
	if qType == QuestionTypeSingle {
		if len(answer) != 1 {
			return fmt.Errorf("单选题答案必须为单个字母，当前: %s", answer)
		}
		if !validLabels[answer[0]] {
			return fmt.Errorf("答案 %s 不在选项范围内", answer)
		}
	}
	
	// 多选题验证
	if qType == QuestionTypeMultiple {
		if len(answer) < 2 {
			return fmt.Errorf("多选题答案至少包含2个字母，当前: %s", answer)
		}
		
		// 验证升序
		for i := 1; i < len(answer); i++ {
			if answer[i] <= answer[i-1] {
				return fmt.Errorf("多选题答案必须升序排列，当前: %s", answer)
			}
		}
		
		// 验证每个字母
		for _, ch := range answer {
			if !validLabels[byte(ch)] {
				return fmt.Errorf("答案字母 %c 不在选项范围内", ch)
			}
		}
	}
	
	return nil
}
```

---

### 7. explanation (答案解析)

**类型**: `string`  
**必填**: ✅  
**长度**: 20-1000字符  
**语言**: 中文  

**规则**:
- 必须说明为何正确答案是正确的
- 应说明其他选项为何错误（尤其单选题）
- 可补充相关知识点
- 语言清晰，逻辑严密
- 引用规范条文时应标注出处

**示例**:
```yaml
# 单选题解析（说明正确和错误原因）
explanation: "A正确，Go语言支持两种注释方式：单行注释使用//，多行注释使用/* */。选项B是Python风格，C是HTML风格，D是SQL风格，均不适用于Go。"

# 多选题解析（逐项说明）
explanation: "A正确，单行注释使用//。B错误，Go的多行注释/* */不支持嵌套。C正确，注释在编译时被忽略。D正确，godoc工具可以提取注释生成文档。"

# 困难题解析（深度说明 + 知识点补充）
explanation: "A错误，因为Go的多行注释/* */不支持嵌套。当解析器遇到第一个*/时，会关闭外层注释，导致'外层继续'成为非法代码。这是C语言注释风格的延续，目的是避免注释块的意外嵌套。B、C、D均合法。"

# 避免的解析
explanation: "答案是A。"                          # ❌ 太简短，未说明原因
explanation: "看规范。"                           # ❌ 无实质内容
explanation: "Because A is correct."              # ❌ 使用英文
```

**最佳实践结构**:
```text
1. 正确答案说明（为何正确）
2. 错误选项说明（为何错误）
3. 知识点补充（可选，增强理解）
```

**验证规则**:
```go
func ValidateExplanation(explanation string) error {
	if len(explanation) < 20 || len(explanation) > 1000 {
		return fmt.Errorf("解析长度必须在20-1000字符之间，当前: %d", len(explanation))
	}
	
	// 检查是否包含中文
	hasChinese := false
	for _, r := range explanation {
		if unicode.Is(unicode.Han, r) {
			hasChinese = true
			break
		}
	}
	if !hasChinese {
		return errors.New("解析必须包含中文说明")
	}
	
	return nil
}
```

---

### 8. topic (所属主题)

**类型**: `string`  
**必填**: ✅  
**枚举值**: `"lexical_elements"` | `"constants"` | `"variables"` | `"types"`  

**规则**:
- 必须与文件路径中的主题目录名一致
- 小写，下划线分隔

**示例**:
```yaml
# 文件路径: quiz_data/lexical_elements/comments.yaml
topic: "lexical_elements"    # ✅ 正确

topic: "constants"           # ❌ 错误：与文件路径不符
topic: "LexicalElements"     # ❌ 错误：大小写不符
topic: "lexical-elements"    # ❌ 错误：使用连字符而非下划线
```

**验证规则**:
```go
func ValidateTopic(topic, filePath string) error {
	// 从文件路径提取主题：quiz_data/lexical_elements/comments.yaml -> lexical_elements
	expectedTopic := extractTopicFromPath(filePath)
	if topic != expectedTopic {
		return fmt.Errorf("主题 %s 与文件路径 %s 不一致", topic, expectedTopic)
	}
	
	// 验证主题在允许范围内
	validTopics := []string{"lexical_elements", "constants", "variables", "types"}
	if !contains(validTopics, topic) {
		return fmt.Errorf("无效的主题: %s", topic)
	}
	
	return nil
}
```

---

### 9. chapter (所属章节)

**类型**: `string`  
**必填**: ✅  
**格式**: 小写英文，下划线分隔  
**长度**: 3-30字符  

**规则**:
- 必须与文件名（去掉.yaml后缀）一致
- 必须在该主题的有效章节列表中

**有效章节列表**:

```yaml
# lexical_elements (11章节)
- comments, tokens, semicolons, identifiers, keywords, operators
- integers, floats, imaginary, runes, strings

# constants (12章节)
- boolean, rune, integer, floating_point, complex, string
- expressions, typed_untyped, conversions, builtin_functions
- iota, implementation_restrictions

# variables (4章节)
- storage, static, dynamic, zero

# types (14章节)
- boolean, numeric, string, array, slice, struct, pointer
- function, interface_basic, interface_embedded
- interface_general, interface_impl, map, channel
```

**示例**:
```yaml
# 文件路径: quiz_data/constants/boolean.yaml
chapter: "boolean"           # ✅ 正确

chapter: "Boolean"           # ❌ 错误：大小写不符
chapter: "bool"              # ❌ 错误：简写不符合规范
chapter: "comments"          # ❌ 错误：不属于constants主题
```

**验证规则**:
```go
func ValidateChapter(chapter, topic, filePath string) error {
	// 从文件路径提取章节：quiz_data/constants/boolean.yaml -> boolean
	expectedChapter := extractChapterFromPath(filePath)
	if chapter != expectedChapter {
		return fmt.Errorf("章节 %s 与文件名 %s 不一致", chapter, expectedChapter)
	}
	
	// 验证章节在该主题的有效列表中
	validChapters := getValidChapters(topic)
	if !contains(validChapters, chapter) {
		return fmt.Errorf("章节 %s 不属于主题 %s", chapter, topic)
	}
	
	return nil
}
```

---

## 完整示例文件

### 单选题为主的章节示例

```yaml
# backend/quiz_data/constants/boolean.yaml
questions:
  - id: const-boolean-001
    type: single
    difficulty: easy
    stem: "Go语言中，布尔常量有哪几个？"
    options:
      - "A: true和false"
      - "B: 0和1"
      - "C: yes和no"
      - "D: TRUE和FALSE"
    answer: "A"
    explanation: "A正确，Go语言的布尔常量只有true和false两个，均为小写。B错误，0和1是整数。C错误，Go没有yes/no。D错误，大写的TRUE/FALSE不是预声明标识符。"
    topic: "constants"
    chapter: "boolean"

  - id: const-boolean-002
    type: single
    difficulty: medium
    stem: "以下哪个表达式的结果是无类型布尔常量？"
    options:
      - "A: true && false"
      - "B: bool(1)"
      - "C: 1 == 1"
      - "D: var b = true"
    answer: "C"
    explanation: "C正确，比较表达式1==1产生无类型布尔常量。A错误，true和false是有类型布尔常量。B错误，bool(1)是类型转换，产生有类型值。D错误，var声明产生布尔变量。"
    topic: "constants"
    chapter: "boolean"
```

### 多选题为主的章节示例

```yaml
# backend/quiz_data/variables/zero.yaml
questions:
  - id: var-zero-001
    type: multiple
    difficulty: easy
    stem: "关于Go语言的零值，以下说法正确的是？（多选）"
    options:
      - "A: 布尔类型的零值是false"
      - "B: 数值类型的零值是0"
      - "C: 字符串的零值是空字符串"
      - "D: 指针的零值是0"
    answer: "ABC"
    explanation: "A正确，bool的零值是false。B正确，所有数值类型（整型、浮点、复数）的零值都是0。C正确，string的零值是空字符串''。D错误，指针的零值是nil，不是0。"
    topic: "variables"
    chapter: "zero"

  - id: var-zero-002
    type: multiple
    difficulty: medium
    stem: "以下哪些类型的零值可以直接使用而不会panic？（多选）"
    options:
      - "A: slice"
      - "B: map"
      - "C: channel"
      - "D: interface"
    answer: "ACD"
    explanation: "A正确，nil slice可以用于len/cap/range，不会panic。B错误，nil map读取返回零值，但写入会panic。C正确，nil channel在select中可用。D正确，nil interface可以进行类型断言，但调用方法会panic。"
    topic: "variables"
    chapter: "zero"
```

---

## 验证清单

题目提交前，请确保通过以下所有检查：

### 结构完整性
- [ ] 所有9个字段均存在
- [ ] 无额外未定义字段
- [ ] YAML语法正确（缩进、引号、换行）

### 字段格式
- [ ] ID符合命名规则（prefix-chapter-seq）
- [ ] type为single或multiple
- [ ] difficulty为easy/medium/hard
- [ ] stem长度10-500字符
- [ ] options数量符合题型要求（单选2-4，多选3-5）
- [ ] options格式为"X: ..."
- [ ] answer格式符合题型要求
- [ ] explanation长度20-1000字符，包含中文
- [ ] topic与文件路径一致
- [ ] chapter与文件名一致

### 内容质量
- [ ] 题干清晰无歧义
- [ ] 选项互斥且合理
- [ ] 答案准确无误
- [ ] 解析详细说明正确和错误原因
- [ ] 多选题题干包含"（多选）"标注

### 业务规则
- [ ] 答案字母在选项范围内
- [ ] 多选题至少2个正确答案
- [ ] 同一文件内ID唯一
- [ ] 难度分布符合40/40/20比例

---

**YAML Schema Status**: ✅ Complete  
**Next**: 生成 quickstart.md
