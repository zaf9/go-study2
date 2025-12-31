# Contracts: Quiz YAML Schema

## 题库 YAML 文件规范

### 文件命名约定

```
backend/quiz_data/{topic}/{chapter}.yaml
```

**示例**:
- `backend/quiz_data/lexical_elements/comments.yaml`
- `backend/quiz_data/lexical_elements/keywords.yaml`
- `backend/quiz_data/constants/boolean.yaml`

### 完整 YAML Schema

```yaml
# ==================== 元数据部分 ====================
metadata:
  topic: string                    # 必需,主题名称 (例如: lexical_elements, constants, variables, types)
  chapter: string                  # 必需,章节名称 (例如: comments, keywords, boolean)
  total_questions: integer         # 必需,题目总数,必须与 questions 列表长度一致
  last_updated: string             # 必需,最后更新日期,格式: YYYY-MM-DD
  difficulty_distribution:         # 必需,难度分布统计
    easy: integer                  # 简单题数量,建议占比 50%-70%
    medium: integer                # 中等题数量,建议占比 25%-40%
    hard: integer                  # 困难题数量,建议占比 5%-15%
  type_distribution:               # 可选,题型分布统计
    single: integer                # 单选题数量
    multiple: integer              # 多选题数量
    code_output: integer           # 代码输出题数量
    code_fix: integer              # 错误修正题数量

# ==================== 题目列表 ====================
questions:
  - id: string                     # 必需,唯一标识符,格式: {topic}-{chapter}-{3位数字}
                                   # 示例: lexical-comments-001
    type: enum                     # 必需,题型: single | multiple | code_output | code_fix
    difficulty: enum               # 必需,难度: easy | medium | hard
    stem: string                   # 必需,题干(问题描述),长度: 20-500 字符
    options: array                 # 必需,选项列表
      - string                     # 选项格式: "A: 选项内容"
      - string                     # 单选题固定4个选项,多选题4-6个选项
    answer: string                 # 必需,正确答案
                                   # 单选题: 单字母 (例如: "B")
                                   # 多选题: 字母组合 (例如: "AC")
    explanation: string            # 必需,详细解析,≥50 字符
                                   # 必须引用 Go 规范章节
    topic: string                  # 必需,所属主题,必须与 metadata.topic 一致
    chapter: string                # 必需,所属章节,必须与 metadata.chapter 一致
    tags: array                    # 可选,标签列表
      - string                     # 例如: ["keywords", "basic"]
    code: string                   # 可选,代码片段 (仅用于 code_output 和 code_fix 题型)
```

### 详细字段规范

#### 1. metadata 元数据部分

| 字段 | 类型 | 必需 | 说明 | 验证规则 |
|------|------|------|------|----------|
| topic | string | ✅ | 主题名称 | 必须是预定义的主题之一: `lexical_elements`, `constants`, `variables`, `types` |
| chapter | string | ✅ | 章节名称 | 小写字母+下划线,例如 `comments`, `keywords` |
| total_questions | integer | ✅ | 题目总数 | 必须 ≥30 且 ≤50,且等于 questions 数组长度 |
| last_updated | string | ✅ | 最后更新日期 | 格式 `YYYY-MM-DD`,例如 `2026-01-15` |
| difficulty_distribution.easy | integer | ✅ | 简单题数量 | 建议占比 50%-70% |
| difficulty_distribution.medium | integer | ✅ | 中等题数量 | 建议占比 25%-40% |
| difficulty_distribution.hard | integer | ✅ | 困难题数量 | 建议占比 5%-15% |

**难度分布验证公式**:
```
easy_ratio = easy / total_questions
medium_ratio = medium / total_questions
hard_ratio = hard / total_questions

ASSERT: 0.5 <= easy_ratio <= 0.7
ASSERT: 0.25 <= medium_ratio <= 0.4
ASSERT: 0.05 <= hard_ratio <= 0.15
ASSERT: easy + medium + hard == total_questions
```

#### 2. questions 题目列表

| 字段 | 类型 | 必需 | 说明 | 验证规则 |
|------|------|------|------|----------|
| id | string | ✅ | 唯一标识符 | 正则: `^[a-z_]+-[a-z_]+-\d{3}$` |
| type | enum | ✅ | 题型 | `single` \| `multiple` \| `code_output` \| `code_fix` |
| difficulty | enum | ✅ | 难度 | `easy` \| `medium` \| `hard` |
| stem | string | ✅ | 题干 | 长度 20-500 字符,必须包含问号"?" |
| options | array | ✅ | 选项列表 | 单选4个,多选4-6个,格式 `"A: 选项内容"` |
| answer | string | ✅ | 正确答案 | 单选: `^[A-Z]$`,多选: `^[A-Z]{2,6}$` |
| explanation | string | ✅ | 详细解析 | ≥50 字符,必须包含"Go 规范"或"官方文档" |
| topic | string | ✅ | 所属主题 | 必须与 metadata.topic 一致 |
| chapter | string | ✅ | 所属章节 | 必须与 metadata.chapter 一致 |
| tags | array | ❌ | 标签列表 | 可选,每个标签长度 ≤20 字符 |
| code | string | ❌ | 代码片段 | 仅用于 `code_output` 和 `code_fix` 题型 |

**题型特定验证规则**:

- **single (单选题)**:
  - `options` 长度必须 == 4
  - `answer` 必须是单字母 (A/B/C/D)
  
- **multiple (多选题)**:
  - `options` 长度必须在 4-6 之间
  - `answer` 必须是2-6个不重复字母 (例如: "AC", "ABD")
  - 正确答案数量建议 ≥2 且 ≤选项数-1

- **code_output (代码输出题)**:
  - 必须提供 `code` 字段
  - `stem` 必须包含"输出"或"结果"关键词
  - `code` 必须是有效的 Go 代码,能通过 `go run` 执行

- **code_fix (错误修正题)**:
  - 必须提供 `code` 字段
  - `stem` 必须包含"错误"或"修正"关键词
  - `options` 提供多种修正方案

---

## 示例 YAML 文件

### 示例 1: comments.yaml (基础题库)

```yaml
metadata:
  topic: lexical_elements
  chapter: comments
  total_questions: 35
  last_updated: "2026-01-15"
  difficulty_distribution:
    easy: 21     # 60%
    medium: 11   # 31%
    hard: 3      # 9%
  type_distribution:
    single: 14
    multiple: 10
    code_output: 7
    code_fix: 4

questions:
  # ========== 单选题示例 ==========
  - id: lexical-comments-001
    type: single
    difficulty: easy
    stem: "在Go语言中,哪种注释方式用于单行注释?"
    options:
      - "A: /* 单行注释 */"
      - "B: // 单行注释"
      - "C: # 单行注释"
      - "D: <!-- 单行注释 -->"
    answer: B
    explanation: "Go语言使用 // 开始单行注释,/* ... */ 用于多行注释。参见 Go 规范 Lexical Elements - Comments 章节。"
    topic: lexical_elements
    chapter: comments
    tags: ["basic", "syntax"]

  # ========== 多选题示例 ==========
  - id: lexical-comments-002
    type: multiple
    difficulty: easy
    stem: "关于Go的多行注释,以下说法正确的是哪些?(多选)"
    options:
      - "A: 多行注释使用 /* 和 */ 包裹"
      - "B: 多行注释可以嵌套"
      - "C: 多行注释会被编译器忽略"
      - "D: 多行注释不能跨越多行"
    answer: AC
    explanation: "多行注释用 /* */ 包裹,编译时被忽略。根据 Go 规范,多行注释不支持嵌套,且可以跨越多行。"
    topic: lexical_elements
    chapter: comments
    tags: ["basic", "syntax"]

  # ========== 代码输出题示例 ==========
  - id: lexical-comments-010
    type: code_output
    difficulty: medium
    stem: "以下代码的输出结果是什么?"
    options:
      - "A: Hello"
      - "B: Hello World"
      - "C: // Hello"
      - "D: 编译错误"
    answer: A
    explanation: "注释内容不会被执行,因此只输出 fmt.Println(\"Hello\") 的结果。Go 规范明确说明注释会被编译器忽略。"
    topic: lexical_elements
    chapter: comments
    tags: ["code", "output"]
    code: |
      package main
      import "fmt"
      func main() {
          fmt.Println("Hello") // World
      }

  # ========== 错误修正题示例 ==========
  - id: lexical-comments-020
    type: code_fix
    difficulty: hard
    stem: "以下代码存在注释相关的错误,正确的修正方案是?"
    options:
      - "A: 将 /* 注释 */ 改为 // 注释"
      - "B: 删除嵌套的 /* 符号"
      - "C: 将所有注释移到行尾"
      - "D: 代码没有错误"
    answer: B
    explanation: "Go 不支持嵌套的多行注释,嵌套的 /* 会导致语法错误。根据 Go 规范,多行注释必须正确配对且不能嵌套。"
    topic: lexical_elements
    chapter: comments
    tags: ["code", "fix", "nesting"]
    code: |
      package main
      func main() {
          /* 外层注释
          /* 嵌套注释 */
          */
          println("Hello")
      }
```

### 示例 2: keywords.yaml (待填充模板)

```yaml
metadata:
  topic: lexical_elements
  chapter: keywords
  total_questions: 0  # 待填充
  last_updated: null
  difficulty_distribution:
    easy: 0
    medium: 0
    hard: 0
  type_distribution:
    single: 0
    multiple: 0
    code_output: 0
    code_fix: 0

questions: []  # 待添加题目
```

---

## 验证脚本

### 格式验证 (validate_quiz_yaml.go)

```go
package main

import (
    "fmt"
    "os"
    "regexp"
    "gopkg.in/yaml.v3"
)

type ValidationError struct {
    File     string
    Question int
    Field    string
    Message  string
}

func ValidateQuizFile(filePath string) []ValidationError {
    var errors []ValidationError
    
    // 1. 解析 YAML
    data, err := os.ReadFile(filePath)
    if err != nil {
        return []ValidationError{{File: filePath, Message: fmt.Sprintf("文件读取失败: %v", err)}}
    }
    
    var quiz ChapterQuizBank
    if err := yaml.Unmarshal(data, &quiz); err != nil {
        return []ValidationError{{File: filePath, Message: fmt.Sprintf("YAML 解析失败: %v", err)}}
    }
    
    // 2. 验证元数据
    if quiz.Metadata.Topic == "" {
        errors = append(errors, ValidationError{File: filePath, Field: "metadata.topic", Message: "topic 不能为空"})
    }
    if quiz.Metadata.Chapter == "" {
        errors = append(errors, ValidationError{File: filePath, Field: "metadata.chapter", Message: "chapter 不能为空"})
    }
    if quiz.Metadata.TotalQuestions != len(quiz.Questions) {
        errors = append(errors, ValidationError{
            File: filePath, 
            Field: "metadata.total_questions",
            Message: fmt.Sprintf("total_questions (%d) 与实际题目数量 (%d) 不一致", quiz.Metadata.TotalQuestions, len(quiz.Questions)),
        })
    }
    
    // 3. 验证难度分布
    dist := quiz.Metadata.DifficultyDistribution
    total := dist.Easy + dist.Medium + dist.Hard
    if total != quiz.Metadata.TotalQuestions {
        errors = append(errors, ValidationError{
            File: filePath,
            Field: "difficulty_distribution",
            Message: fmt.Sprintf("难度分布总和 (%d) 与 total_questions (%d) 不一致", total, quiz.Metadata.TotalQuestions),
        })
    }
    
    // 4. 验证每个题目
    idRegex := regexp.MustCompile(`^[a-z_]+-[a-z_]+-\d{3}$`)
    answerSingleRegex := regexp.MustCompile(`^[A-Z]$`)
    answerMultipleRegex := regexp.MustCompile(`^[A-Z]{2,6}$`)
    
    for i, q := range quiz.Questions {
        // ID 格式验证
        if !idRegex.MatchString(q.ID) {
            errors = append(errors, ValidationError{
                File: filePath,
                Question: i + 1,
                Field: "id",
                Message: fmt.Sprintf("ID 格式不正确: %s", q.ID),
            })
        }
        
        // 题型验证
        validTypes := map[QuestionType]bool{
            QuestionTypeSingle: true,
            QuestionTypeMultiple: true,
            QuestionTypeCodeOutput: true,
            QuestionTypeCodeFix: true,
        }
        if !validTypes[q.Type] {
            errors = append(errors, ValidationError{
                File: filePath,
                Question: i + 1,
                Field: "type",
                Message: fmt.Sprintf("无效的题型: %s", q.Type),
            })
        }
        
        // 选项数量验证
        if q.Type == QuestionTypeSingle && len(q.Options) != 4 {
            errors = append(errors, ValidationError{
                File: filePath,
                Question: i + 1,
                Field: "options",
                Message: fmt.Sprintf("单选题必须有4个选项,当前有 %d 个", len(q.Options)),
            })
        }
        if q.Type == QuestionTypeMultiple && (len(q.Options) < 4 || len(q.Options) > 6) {
            errors = append(errors, ValidationError{
                File: filePath,
                Question: i + 1,
                Field: "options",
                Message: fmt.Sprintf("多选题必须有4-6个选项,当前有 %d 个", len(q.Options)),
            })
        }
        
        // 答案格式验证
        if q.Type == QuestionTypeSingle && !answerSingleRegex.MatchString(q.Answer) {
            errors = append(errors, ValidationError{
                File: filePath,
                Question: i + 1,
                Field: "answer",
                Message: fmt.Sprintf("单选题答案格式不正确: %s", q.Answer),
            })
        }
        if q.Type == QuestionTypeMultiple && !answerMultipleRegex.MatchString(q.Answer) {
            errors = append(errors, ValidationError{
                File: filePath,
                Question: i + 1,
                Field: "answer",
                Message: fmt.Sprintf("多选题答案格式不正确: %s", q.Answer),
            })
        }
        
        // 解析验证
        if len(q.Explanation) < 50 {
            errors = append(errors, ValidationError{
                File: filePath,
                Question: i + 1,
                Field: "explanation",
                Message: fmt.Sprintf("解析长度不足50字符,当前 %d 字符", len(q.Explanation)),
            })
        }
        if !strings.Contains(q.Explanation, "Go 规范") && !strings.Contains(q.Explanation, "官方文档") {
            errors = append(errors, ValidationError{
                File: filePath,
                Question: i + 1,
                Field: "explanation",
                Message: "解析必须引用 Go 规范或官方文档",
            })
        }
    }
    
    return errors
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("用法: go run validate_quiz_yaml.go <yaml文件路径>")
        os.Exit(1)
    }
    
    errors := ValidateQuizFile(os.Args[1])
    if len(errors) == 0 {
        fmt.Println("✅ 验证通过!")
    } else {
        fmt.Printf("❌ 发现 %d 个错误:\n", len(errors))
        for _, err := range errors {
            if err.Question > 0 {
                fmt.Printf("  [题目 %d] %s: %s\n", err.Question, err.Field, err.Message)
            } else {
                fmt.Printf("  %s: %s\n", err.Field, err.Message)
            }
        }
        os.Exit(1)
    }
}
```

---

## 内容质量验证规则

### 1. 自动化内容检查 (80%+)

```yaml
automated_checks:
  # 引用检查
  reference_check:
    rule: "explanation 字段必须包含'Go 规范'或'官方文档'关键词"
    weight: 20%
  
  # 代码验证
  code_validation:
    rule: "code_output 题型的代码必须能通过 go run 执行"
    command: "echo '<code>' > /tmp/test.go && go run /tmp/test.go"
    weight: 15%
  
  # 答案一致性
  answer_consistency:
    rule: "正确答案必须在 options 中存在"
    weight: 10%
  
  # 难度校准
  difficulty_calibration:
    easy:
      explanation_max_length: 100
      stem_max_length: 150
    medium:
      explanation_length: [100, 250]
      stem_length: [150, 300]
    hard:
      explanation_min_length: 200
      stem_min_length: 250
    weight: 15%
  
  # 去重检查
  deduplication:
    rule: "同一章节题目的 stem 不允许完全重复"
    similarity_threshold: 0.9
    weight: 20%
```

### 2. 人工抽查 (20%)

```yaml
manual_spot_check:
  sample_rate: 0.2  # 每章节抽查20%题目
  priority:
    - medium: 40%   # 优先抽查中等难度
    - hard: 40%     # 优先抽查困难题目
    - easy: 20%     # 少量抽查简单题目
  
  checklist:
    - "题干表述是否清晰无歧义?"
    - "选项是否有明显错误或重复?"
    - "正确答案是否唯一且明确?"
    - "解析是否准确引用规范章节?"
    - "代码示例是否能正确运行?"
    - "难度标签是否合理?"
```

---

## 版本控制规范

### Git Commit 消息格式

```
feat(quiz): 新增 {topic}/{chapter} 章节题库 ({n}题)

- 添加 {easy}道简单题, {medium}道中等题, {hard}道困难题
- 题型分布: 单选{single}道, 多选{multiple}道, 代码题{code}道
- 已通过格式验证和内容抽查

Closes #issue-number
```

**示例**:
```
feat(quiz): 新增 lexical_elements/keywords 章节题库 (45题)

- 添加 27道简单题, 14道中等题, 4道困难题
- 题型分布: 单选18道, 多选15道, 代码题12道
- 已通过格式验证和内容抽查

Closes #017-keywords
```

### 文件变更追踪

```yaml
# 每次修改必须更新 metadata.last_updated
# 并在 commit 消息中说明修改原因
fix(quiz): 修正 lexical_elements/comments 第5题答案错误

- 题目ID: lexical-comments-005
- 原答案: A (错误)
- 新答案: B (正确)
- 原因: 根据 Go 1.24 规范勘误,多行注释确实不支持嵌套

Closes #017-fix-comments-005
```

