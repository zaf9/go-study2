# 题库补全快速开始指南

## 🎯 目标

本指南帮助开发者快速上手题库补全工作,包括:
- 如何生成新章节题库模板
- 如何编写高质量题目内容
- 如何使用验证工具保证质量
- 如何测试 Quiz API

---

## � 术语对照表

为了在用户文档(中文)和技术实现(YAML 枚举值)之间保持一致,请参考以下术语映射:

### 题型术语

| 中文术语 | YAML 枚举值 | 说明 | 示例 |
|---------|------------|------|------|
| 单选题 | `single_choice` | 4个选项,1个正确答案 | 以下哪个是 Go 关键字? |
| 多选题 | `multiple_choice` | 4-6个选项,2-3个正确答案 | Go 语言支持以下哪些数据类型? (多选) |
| 代码输出题 | `code_output` | 预测代码执行结果 | 以下代码的输出是什么? |
| 错误修正题 | `code_fix` | 找出并修复代码错误 | 以下代码无法编译,如何修复? |

### 难度术语

| 中文术语 | YAML 枚举值 | 目标比例 | 判断标准 |
|---------|------------|---------|----------|
| 简单题 | `easy` | 50-70% | 直接概念记忆,无需推理 |
| 中等题 | `medium` | 25-40% | 需要理解规范细节或简单代码分析 |
| 困难题 | `hard` | 5-15% | 需要综合分析、代码执行预测、边界情况处理 |

### 数据实体术语

| 概念层(数据模型) | 存储层(文件格式) | 运行时(内存表示) | 说明 |
|---------------|---------------|----------------|------|
| ChapterQuizBank | `backend/quiz_data/{topic}/{chapter}.yaml` | QuizBank struct | 章节题库:概念实体 vs 存储实现 |
| Question | YAML questions 数组中的一项 | Question struct | 单个题目 |
| QuizAttempt | `data/attempts/{user_id}/{quiz_id}.json` | QuizAttempt struct | 测验记录 |
| ChapterProgress | 内存计算(不持久化) | ChapterProgress struct | 章节进度 |

**使用建议**:
- 编写用户文档时使用**中文术语**
- 编写 YAML 文件时使用 **YAML 枚举值**
- 编写 Go 代码时使用**数据模型名称**
- 文档中首次提到术语时附带枚举值,如"单选题 (single_choice)"

---

## �📋 前置条件

### 开发环境

```bash
# 1. Go 版本 ≥1.24
go version  # 输出应包含 "go1.24"

# 2. 安装依赖
cd backend
go mod download

# 3. 验证工具可用
go run scripts/validate_quiz_yaml.go --help
```

### 必读文档

- **Go 1.24 规范**: https://go.dev/ref/spec (题目内容的权威来源)
- `research.md`: 理解题目生成策略和难度分布
- `data-model.md`: 理解数据实体和验证规则
- `contracts/quiz-yaml-schema.md`: 理解 YAML 格式规范

---

## 🚀 5 分钟快速开始

### 步骤 1: 选择章节

从 [41 个章节列表](research.md#go-124-章节结构) 中选择一个未完成的章节,例如:

```bash
# 查看当前进度
cd backend
./scripts/check_quiz_progress.sh

# 输出示例:
# ✅ lexical_elements/comments: 35 questions COMPLETE
# ⏳ lexical_elements/keywords: PENDING
# ⏳ lexical_elements/identifiers: PENDING
# ...
```

### 步骤 2: 生成模板

使用模板生成器创建章节文件骨架:

```bash
# 语法: generate_quiz_template.go <topic> <chapter> <question_count>
go run scripts/generate_quiz_template.go lexical_elements keywords 40

# 输出: 
# ✅ Created: backend/quiz_data/lexical_elements/keywords.yaml
# 📝 Template contains:
#    - Metadata section
#    - 40 question placeholders (16 single, 12 multiple, 8 code_output, 4 code_fix)
```

生成的 `keywords.yaml` 结构:

```yaml
# Metadata
topic: "lexical_elements"
chapter: "keywords"
total_questions: 40
last_updated: "2025-01-21"
difficulty_distribution:
  easy: 60%
  medium: 30%
  hard: 10%

questions:
  - id: "lexical-keywords-001"
    type: "single_choice"
    difficulty: "easy"
    stem: "TODO: 编写题干"
    options:
      A: "TODO: 选项A"
      B: "TODO: 选项B"
      C: "TODO: 选项C"
      D: "TODO: 选项D"
    answer: "TODO: A/B/C/D"
    explanation: "TODO: 详细解析 (≥50字符)"
  
  # ... 39 more questions
```

### 步骤 3: 编写题目内容

打开生成的 YAML 文件,逐个填充题目内容。

**题目编写指南**:

#### 3.1 题干 (stem) 编写规范

✅ **好的题干示例**:

```yaml
# 基础概念题 (easy)
stem: "在 Go 语言中,以下哪个是关键字?"

# Go 规范引用题 (medium)
stem: "根据 Go 规范,`for` 关键字可以用于以下哪种语句结构?"

# 代码理解题 (hard)
stem: |
  以下代码的输出是什么?
  ```go
  package main
  import "fmt"
  func main() {
      for i := 0; i < 3; i++ {
          defer fmt.Print(i)
      }
  }
  ```
```

❌ **不好的题干示例**:

```yaml
# 过于简单,无学习价值
stem: "Go 有关键字吗?"

# 无 Go 规范引用,无权威性
stem: "有人说 `class` 是 Go 关键字,对吗?"

# 题干模糊,无法判断答案
stem: "关于 `for` 关键字"
```

#### 3.2 选项 (options) 编写规范

✅ **好的选项示例**:

```yaml
# 单选题: 有 3 个干扰项
options:
  A: "for (Go 关键字,正确答案)"
  B: "foreach (不是 Go 关键字)"
  C: "while (不是 Go 关键字)"
  D: "do (不是 Go 关键字)"

# 多选题: 答案 2-4 个,干扰项 2-4 个
options:
  A: "var (变量声明关键字,正确)"
  B: "const (常量声明关键字,正确)"
  C: "let (JavaScript 关键字,错误)"
  D: "def (Python 关键字,错误)"
  E: "type (类型定义关键字,正确)"
  F: "class (Java 关键字,错误)"
```

❌ **不好的选项示例**:

```yaml
# 选项数量不足
options:
  A: "for"
  B: "while"

# 选项无干扰性 (答案太明显)
options:
  A: "for (Go 关键字)"
  B: "这不是关键字"
  C: "这也不是"
  D: "这还不是"

# 多选题只有 1 个正确答案
options:
  A: "for (正确)"
  B: "abc (错误)"
  C: "def (错误)"
  D: "ghi (错误)"
```

#### 3.3 解析 (explanation) 编写规范

✅ **好的解析示例**:

```yaml
explanation: |
  根据 Go 规范 2.1 节"Keywords",Go 语言共有 25 个关键字:`break`, `case`, `chan`, ..., `var`。
  其中 `for` 是循环控制关键字,可以实现传统 for 循环、while 循环和 range 循环三种形式。
  
  示例:
  ```go
  // 传统 for 循环
  for i := 0; i < 10; i++ { ... }
  
  // while 风格循环
  for condition { ... }
  
  // range 循环
  for k, v := range m { ... }
  ```
  
  **选项分析**:
  - A ✅ `for` 是 Go 关键字,正确
  - B ❌ `foreach` 是其他语言的关键字,Go 中不存在
  - C ❌ `while` 不是 Go 关键字,Go 使用 `for condition` 实现 while 逻辑
  - D ❌ `do` 不是 Go 关键字,Go 没有 do-while 循环
  
  **拓展**: Go 的 25 个关键字分为 5 类:声明 (func, var, const, type, import, package)、
  控制流 (if, else, for, switch, case, default, select, return, break, continue, goto, fallthrough)、
  并发 (go, chan)、数据结构 (struct, interface, map) 和 延迟执行 (defer)。
```

❌ **不好的解析示例**:

```yaml
# 太简短,无学习价值
explanation: "答案是 A,因为 `for` 是关键字。"

# 无 Go 规范引用
explanation: "根据我的经验,`for` 应该是关键字。"

# 无选项分析
explanation: "`for` 是 Go 关键字,可以用于循环。"

# 少于 50 字符 (验证工具会报错)
explanation: "`for` 是关键字"
```

#### 3.4 难度 (difficulty) 判断标准

| 难度 | 判断标准 | 示例题型 | 目标比例 |
|------|---------|---------|---------|
| `easy` | 直接概念记忆,Go 规范基础知识,无需推理 | "以下哪个是关键字?" | 50-70% |
| `medium` | 需要理解规范细节,或简单代码分析 | "for 关键字有哪些使用形式?" | 25-40% |
| `hard` | 需要综合分析、代码执行预测、边界情况处理 | "以下代码输出是什么?(涉及 defer 执行顺序)" | 5-15% |

### 步骤 4: 格式验证

完成题目编写后,运行格式验证器:

```bash
# 验证单个文件
go run scripts/validate_quiz_yaml.go backend/quiz_data/lexical_elements/keywords.yaml

# 输出示例:
# ✅ PASS: backend/quiz_data/lexical_elements/keywords.yaml
# 📊 Statistics:
#    - Total questions: 40
#    - Single choice: 16 (40.0%)
#    - Multiple choice: 12 (30.0%)
#    - Code output: 8 (20.0%)
#    - Code fix: 4 (10.0%)
#    - Easy: 24 (60%), Medium: 12 (30%), Hard: 4 (10%)
# ✅ All validations passed!

# 如果有错误,会显示详细信息:
# ❌ FAIL: backend/quiz_data/lexical_elements/keywords.yaml
# 
# Errors:
# - Line 15: Question lexical-keywords-003
#   ❌ ID format invalid: "keywords-003" (expected: "lexical-keywords-003")
# 
# - Line 42: Question lexical-keywords-007
#   ❌ Type invalid: "single" (expected: single_choice/multiple_choice/code_output/code_fix)
# 
# - Line 68: Question lexical-keywords-012
#   ❌ Explanation too short: 28 characters (minimum: 50)
```

**常见验证错误及修复**:

| 错误类型 | 错误示例 | 修复方法 |
|---------|---------|---------|
| ID 格式错误 | `id: "keywords-001"` | 改为 `id: "lexical-keywords-001"` (格式: `{topic}-{chapter}-{序号}`) |
| 题型错误 | `type: "single"` | 改为 `type: "single_choice"` |
| 难度错误 | `difficulty: "简单"` | 改为 `difficulty: "easy"` |
| 选项数量错误 | 单选题只有 2 个选项 | 补充到 4 个选项 (A-D) |
| 答案格式错误 | `answer: "选项A"` | 改为 `answer: "A"` (单选) 或 `answer: "ACD"` (多选) |
| 解析太短 | `explanation: "答案是 A"` | 扩展到 ≥50 字符,包含规范引用和选项分析 |

### 步骤 5: 内容质量检查

运行质量检查工具,确保题目内容符合要求:

```bash
# 检查整个章节
go run scripts/quiz_quality_check.go backend/quiz_data/lexical_elements/

# 输出示例:
# 🔍 Checking: backend/quiz_data/lexical_elements/keywords.yaml
# 
# ✅ Format validation: PASS
# 
# 📊 Content quality checks:
# ✅ Explanation length check: 40/40 (100%) ≥50 chars
# ✅ Go spec reference check: 38/40 (95%) contains "Go 规范" or "官方文档"
#    ⚠️  Missing reference: lexical-keywords-012, lexical-keywords-029
# ✅ Code executability check: 8/8 (100%) code_output questions can compile
# ✅ Duplicate detection: 0 duplicates found
# 
# 🎯 Overall quality score: 92.5% (target: 80%+)
# ✅ PASS: Quality meets requirements
```

**质量检查项说明**:

| 检查项 | 目标 | 自动化 | 说明 |
|-------|------|-------|------|
| 解析长度 | ≥50 字符 | ✅ | 确保解析有足够的学习价值 |
| Go 规范引用 | ≥80% 题目包含"Go 规范"或"官方文档" | ✅ | 确保权威性 |
| 代码可执行性 | 100% `code_output` 题型代码可编译 | ✅ | 防止语法错误 |
| 题目去重 | 无重复题干 | ✅ | 使用相似度算法检测 |
| 人工抽查 | 20% 题目人工审核 | ❌ | 检查语义正确性、选项合理性 |

### 步骤 6: 提交代码

通过验证后,提交到 Git:

```bash
# 1. 切换到功能分支 (如果还没有)
git checkout 017-complete-quiz-bank

# 2. 添加新文件
git add backend/quiz_data/lexical_elements/keywords.yaml

# 3. 提交 (使用标准化 commit message)
git commit -m "feat(quiz): 新增 lexical_elements/keywords 章节题库 (40题)

- 单选题 16 道 (40%)
- 多选题 12 道 (30%)
- 代码输出题 8 道 (20%)
- 错误修正题 4 道 (10%)
- 难度分布: easy 60%, medium 30%, hard 10%
- 通过格式验证和内容质量检查 (质量分数 92.5%)

Refs: #17"

# 4. 推送到远程 (如果需要)
git push origin 017-complete-quiz-bank
```

**标准化 commit message 格式** (参考 `contracts/quiz-yaml-schema.md`):

```
feat(quiz): 新增 {topic}/{chapter} 章节题库 ({题目数量}题)

- 单选题 {数量} 道 ({百分比})
- 多选题 {数量} 道 ({百分比})
- 代码输出题 {数量} 道 ({百分比})
- 错误修正题 {数量} 道 ({百分比})
- 难度分布: easy {百分比}, medium {百分比}, hard {百分比}
- 通过格式验证和内容质量检查 (质量分数 {分数})

Refs: #17
```

---

## 🧪 测试 Quiz API

### 本地启动后端服务

```bash
cd backend
go run main.go

# 输出:
# 🚀 GoFrame QuizBank Server starting...
# 📂 Loading quiz banks from: /quiz_data
# ✅ Loaded 1 chapters: lexical_elements/comments
# 🌐 Listening on: http://localhost:8080
```

### API 测试示例

#### 1. 获取所有章节列表

```bash
curl http://localhost:8080/api/v1/quiz/chapters

# 响应:
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "topics": [
#       {
#         "topic": "lexical_elements",
#         "chapters": [
#           {
#             "chapter": "comments",
#             "total_questions": 35,
#             "difficulty_distribution": {
#               "easy": 60,
#               "medium": 30,
#               "hard": 10
#             },
#             "last_updated": "2025-01-21"
#           },
#           {
#             "chapter": "keywords",
#             "total_questions": 40,
#             "difficulty_distribution": {
#               "easy": 60,
#               "medium": 30,
#               "hard": 10
#             },
#             "last_updated": "2025-01-21"
#           }
#         ]
#       }
#     ]
#   }
# }
```

#### 2. 开始测验 (随机抽题)

```bash
curl http://localhost:8080/api/v1/quiz/lexical_elements/keywords

# 响应:
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "quiz_id": "q-20250121-abc123",
#     "topic": "lexical_elements",
#     "chapter": "keywords",
#     "questions": [
#       {
#         "id": "lexical-keywords-001",
#         "type": "single_choice",
#         "difficulty": "easy",
#         "stem": "在 Go 语言中,以下哪个是关键字?",
#         "options": {
#           "A": "for",
#           "B": "foreach",
#           "C": "while",
#           "D": "do"
#         }
#         // 注意: 不返回 answer 和 explanation
#       },
#       // ... 其他 5-9 道题
#     ],
#     "created_at": "2025-01-21T10:30:00Z"
#   }
# }
```

#### 3. 提交测验答案

```bash
curl -X POST http://localhost:8080/api/v1/quiz/lexical_elements/keywords \
  -H "Content-Type: application/json" \
  -d '{
    "quiz_id": "q-20250121-abc123",
    "user_id": "user-001",
    "answers": {
      "lexical-keywords-001": "A",
      "lexical-keywords-005": "B",
      "lexical-keywords-012": "ACD",
      "lexical-keywords-018": "def main() {...}",
      "lexical-keywords-025": "D"
    }
  }'

# 响应:
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "quiz_id": "q-20250121-abc123",
#     "user_id": "user-001",
#     "score": 80,
#     "total_questions": 6,
#     "correct_count": 5,
#     "wrong_count": 1,
#     "details": [
#       {
#         "question_id": "lexical-keywords-001",
#         "user_answer": "A",
#         "correct_answer": "A",
#         "is_correct": true,
#         "explanation": "根据 Go 规范 2.1 节..."
#       },
#       {
#         "question_id": "lexical-keywords-005",
#         "user_answer": "B",
#         "correct_answer": "C",
#         "is_correct": false,
#         "explanation": "..."
#       },
#       // ... 其他题目详情
#     ],
#     "submitted_at": "2025-01-21T10:35:00Z"
#   }
# }
```

#### 4. 查询历史记录

```bash
curl "http://localhost:8080/api/v1/quiz/history/lexical_elements/keywords?user_id=user-001"

# 响应:
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "attempts": [
#       {
#         "quiz_id": "q-20250121-abc123",
#         "score": 80,
#         "submitted_at": "2025-01-21T10:35:00Z"
#       },
#       {
#         "quiz_id": "q-20250120-def456",
#         "score": 60,
#         "submitted_at": "2025-01-20T15:20:00Z"
#       }
#     ],
#     "best_score": 80,
#     "attempt_count": 2
#   }
# }
```

#### 5. 查询进度概览

```bash
curl "http://localhost:8080/api/v1/quiz/progress?user_id=user-001"

# 响应:
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "total_chapters": 2,
#     "completed_chapters": 1,
#     "in_progress_chapters": 1,
#     "overall_score": 70,
#     "chapters": [
#       {
#         "topic": "lexical_elements",
#         "chapter": "comments",
#         "status": "completed",
#         "best_score": 90,
#         "attempt_count": 3,
#         "last_attempt": "2025-01-19T12:00:00Z"
#       },
#       {
#         "topic": "lexical_elements",
#         "chapter": "keywords",
#         "status": "in_progress",
#         "best_score": 80,
#         "attempt_count": 2,
#         "last_attempt": "2025-01-21T10:35:00Z"
#       }
#     ]
#   }
# }
```

---

## 🛠️ 工具链详细说明

### 1. generate_quiz_template.go

**用途**: 生成章节题库 YAML 模板文件

**用法**:

```bash
go run scripts/generate_quiz_template.go <topic> <chapter> <question_count> [--difficulty-easy=60] [--difficulty-medium=30] [--difficulty-hard=10]

# 示例:
go run scripts/generate_quiz_template.go lexical_elements keywords 40

# 自定义难度分布:
go run scripts/generate_quiz_template.go types slice 50 --difficulty-easy=50 --difficulty-medium=40 --difficulty-hard=10
```

**输出**:
- 创建 `backend/quiz_data/{topic}/{chapter}.yaml` 文件
- 自动生成元数据 (topic, chapter, total_questions, last_updated)
- 自动分配题型 (单选 40%, 多选 30%, code_output 20%, code_fix 10%)
- 自动分配难度 (easy/medium/hard 按指定比例)
- 自动生成题目 ID (`{topic}-{chapter}-001` ~ `{topic}-{chapter}-{N}`)

### 2. validate_quiz_yaml.go

**用途**: 验证 YAML 文件格式和结构

**用法**:

```bash
# 验证单个文件
go run scripts/validate_quiz_yaml.go backend/quiz_data/lexical_elements/keywords.yaml

# 验证整个目录
go run scripts/validate_quiz_yaml.go backend/quiz_data/lexical_elements/

# 严格模式 (停止在第一个错误)
go run scripts/validate_quiz_yaml.go --strict backend/quiz_data/lexical_elements/keywords.yaml
```

**验证规则**:
- ✅ YAML 语法正确性
- ✅ 必填字段完整性 (topic, chapter, questions, etc.)
- ✅ ID 格式 (正则: `^[a-z_]+-[a-z_]+-\d{3}$`)
- ✅ 题型枚举值 (single_choice/multiple_choice/code_output/code_fix)
- ✅ 难度枚举值 (easy/medium/hard)
- ✅ 选项数量 (单选 4 个, 多选 4-6 个)
- ✅ 答案格式 (单选: `^[A-Z]$`, 多选: `^[A-Z]{2,6}$`)
- ✅ 解析长度 (≥50 字符)
- ✅ 难度分布符合规范 (easy 50-70%, medium 25-40%, hard 5-15%)
- ✅ 题型分布符合规范 (single 35-45%, multiple 25-35%, code 25-35%)

### 3. quiz_quality_check.go

**用途**: 检查题目内容质量 (80%+ 自动化)

**用法**:

```bash
# 检查单个文件
go run scripts/quiz_quality_check.go backend/quiz_data/lexical_elements/keywords.yaml

# 检查整个目录
go run scripts/quiz_quality_check.go backend/quiz_data/lexical_elements/

# 生成 HTML 报告
go run scripts/quiz_quality_check.go --output=report.html backend/quiz_data/
```

**检查规则**:
- ✅ 解析长度检查 (≥50 字符,100% 自动化)
- ✅ Go 规范引用检查 (≥80% 题目包含"Go 规范"或"官方文档")
- ✅ 代码可执行性验证 (`code_output` 题型代码可编译)
- ✅ 题目相似度检查 (防止重复题目)
- ⚠️ 选项合理性检查 (检测选项是否有明显错误,需人工确认)
- ⚠️ 干扰项有效性检查 (检测干扰项是否有足够区分度,需人工确认)

**质量分数计算**:

```
Quality Score = (解析长度合格率 × 25%) + 
                (Go 规范引用率 × 25%) + 
                (代码可执行率 × 25%) + 
                (无重复率 × 25%)

通过条件: Quality Score ≥ 80%
```

### 4. check_quiz_progress.sh

**用途**: 统计题库覆盖进度

**用法**:

```bash
cd backend
./scripts/check_quiz_progress.sh

# 输出 JSON 格式:
./scripts/check_quiz_progress.sh --json > progress.json
```

**输出示例**:

```
📊 Quiz Bank Coverage Progress
═══════════════════════════════════════════════════════════

📁 lexical_elements (9 chapters)
  ✅ comments: 35 questions (last_updated: 2025-01-21) COMPLETE
  ✅ keywords: 40 questions (last_updated: 2025-01-21) COMPLETE
  ⏳ identifiers: PENDING
  ⏳ tokens: PENDING
  ⏳ integers: PENDING
  ⏳ floats: PENDING
  ⏳ strings: PENDING
  ⏳ semicolons: PENDING
  ⏳ imaginary: PENDING
  ⏳ runes: PENDING

📁 constants (6 chapters)
  ⏳ boolean: PENDING
  ⏳ numeric: PENDING
  ⏳ string: PENDING
  ⏳ iota: PENDING
  ⏳ typed: PENDING
  ⏳ untyped: PENDING

📁 variables (5 chapters)
  ⏳ declaration: PENDING
  ⏳ initialization: PENDING
  ⏳ scope: PENDING
  ⏳ shadowing: PENDING
  ⏳ zero_values: PENDING

📁 types (21 chapters)
  ⏳ boolean: PENDING
  ⏳ numeric: PENDING
  ... (省略其他 19 章)

═══════════════════════════════════════════════════════════
📊 Summary:
  Total Chapters: 41
  Completed: 2 (4.9%)
  Pending: 39 (95.1%)
  Total Questions: 75

🎯 Target: 1200-2000 questions
📈 Current Progress: 75/1200 (6.3%)
⏱️  Estimated Remaining Work: ~39 chapters × 30-50 questions
```

---

## 📚 常见问题 (FAQ)

### Q1: 如何确定一个章节应该出多少题?

**A**: 参考以下建议:

| 章节复杂度 | 题目数量建议 | 示例章节 |
|-----------|------------|---------|
| 简单 (基础概念为主) | 30-35 题 | comments, keywords, identifiers |
| 中等 (有多个知识点) | 35-45 题 | variables/scope, types/slice, types/map |
| 复杂 (涉及多个子主题) | 45-50 题 | types/interface, types/channel |

**灵活调整原则**: 
- 高频考点章节可以增加 5-10 题
- 边缘知识点章节可以减少 5-10 题
- 总体控制在 1200-2000 题之间即可

### Q2: 题目的难度如何判断?

**A**: 使用 3 个维度判断:

1. **概念难度**: 
   - Easy: 直接概念记忆 (如"哪个是关键字?")
   - Medium: 理解规范细节 (如"`for` 有哪些使用形式?")
   - Hard: 综合分析 (如"预测代码输出,涉及多个知识点")

2. **思维深度**:
   - Easy: 1 步思维 (直接回答)
   - Medium: 2-3 步思维 (理解 → 判断)
   - Hard: 4+ 步思维 (理解 → 分析 → 推理 → 判断)

3. **规范熟悉度**:
   - Easy: Go 基础常识,无需查规范
   - Medium: 需要查阅规范某一节
   - Hard: 需要综合多节规范内容

**难度分布检查**: 使用 `validate_quiz_yaml.go` 自动检查是否符合 easy 50-70%, medium 25-40%, hard 5-15%

### Q3: 如何避免题目重复?

**A**: 使用以下策略:

1. **自动化检测**: 
   ```bash
   go run scripts/quiz_quality_check.go --check-duplicates backend/quiz_data/
   ```
   工具会计算题干的相似度 (Jaccard Similarity),超过 80% 相似度会报警。

2. **ID 管理**: 
   - 严格按照 `{topic}-{chapter}-{序号}` 命名
   - 序号从 001 开始连续递增,不跳号
   - 不同章节可以有相似题目 (如 lexical_elements/keywords 和 types/struct 都可能涉及 `type` 关键字)

3. **人工审查**: 
   - 每个章节完成后,通读一遍所有题干
   - 检查是否有换汤不换药的重复题目

### Q4: 代码题 (code_output/code_fix) 如何编写?

**A**: 遵循以下规范:

#### code_output 题型

```yaml
- id: "lexical-keywords-015"
  type: "code_output"
  difficulty: "medium"
  stem: |
    以下代码的输出是什么?
    ```go
    package main
    import "fmt"
    func main() {
        for i := 0; i < 3; i++ {
            defer fmt.Print(i)
        }
    }
    ```
  options:
    A: "0 1 2"
    B: "2 1 0"
    C: "0 2 1"
    D: "编译错误"
  answer: "B"
  explanation: |
    根据 Go 规范,`defer` 语句会将函数调用推入栈中,在外层函数返回时以 **后进先出 (LIFO)** 的顺序执行。
    
    代码执行过程:
    1. i=0 时,defer fmt.Print(0) 推入栈底
    2. i=1 时,defer fmt.Print(1) 推入栈中
    3. i=2 时,defer fmt.Print(2) 推入栈顶
    4. main() 返回前,依次弹出栈: Print(2) → Print(1) → Print(0)
    
    因此输出为 "2 1 0"。
    
    **验证方法**: 将代码粘贴到 Go Playground (https://go.dev/play/) 运行即可验证。
```

**关键点**:
- ✅ 代码必须可编译可运行 (工具会自动验证)
- ✅ 包含 `package main` 和 `func main()`
- ✅ 代码长度控制在 15 行以内 (太长影响阅读)
- ✅ 涉及的知识点必须在当前章节范围内
- ✅ 解析中说明代码执行过程,不只给答案

#### code_fix 题型

```yaml
- id: "lexical-keywords-020"
  type: "code_fix"
  difficulty: "hard"
  stem: |
    以下代码无法编译,如何修复?
    ```go
    package main
    import "fmt"
    func main() {
        var x int = 10
        if x > 5
            fmt.Println("x is greater than 5")
        }
    }
    ```
  options:
    A: "在 if 后面添加 `{`"
    B: "在 x > 5 后面添加 `{`"
    C: "删除 if 语句中的空格"
    D: "将 var x int = 10 改为 x := 10"
  answer: "B"
  explanation: |
    编译错误原因:根据 Go 规范,`if` 语句的条件表达式后必须紧跟 `{`,不能直接换行。
    
    错误代码:
    ```go
    if x > 5
        fmt.Println(...)  // ❌ 缺少 {
    }
    ```
    
    正确代码:
    ```go
    if x > 5 {  // ✅ 条件后紧跟 {
        fmt.Println(...)
    }
    ```
    
    **选项分析**:
    - A ❌ `if` 后紧跟的是条件表达式,不能直接加 `{`
    - B ✅ 条件表达式 `x > 5` 后面必须加 `{`,正确答案
    - C ❌ 空格不影响语法,Go 编译器会忽略空格
    - D ❌ 变量声明方式不影响 if 语句的语法错误
    
    **拓展**: Go 语言强制要求 `{` 必须在同一行,这是为了避免自动分号插入 (Automatic Semicolon Insertion) 导致的歧义。
```

**关键点**:
- ✅ 代码错误必须是 **语法错误** 或 **编译错误**,不能是逻辑错误
- ✅ 错误点必须明确,不能有多个错误 (除非是组合题)
- ✅ 所有选项都必须是可操作的修复方案 (即使是错误选项)
- ✅ 解析中说明为什么其他选项不对

### Q5: 如何保证题目的权威性?

**A**: 遵循以下原则:

1. **必须引用 Go 1.24 官方规范**:
   - 每道题的解析中至少提及一次"Go 规范"或"官方文档"
   - 引用具体章节号,如"根据 Go 规范 2.1 节 'Keywords'..."
   - Go 规范链接: https://go.dev/ref/spec

2. **可验证性**:
   - `code_output` 题型的代码必须可以在 Go Playground 运行
   - 概念题必须能在规范中找到对应描述
   - 提供参考链接 (如 Go Blog, Effective Go)

3. **避免主观判断**:
   - ❌ "根据我的经验..."
   - ❌ "通常情况下..."
   - ✅ "根据 Go 规范..."
   - ✅ "Go 官方文档明确指出..."

4. **质量检查**:
   ```bash
   go run scripts/quiz_quality_check.go backend/quiz_data/
   # 检查是否 80%+ 题目包含规范引用
   ```

### Q6: 题目编写进度太慢怎么办?

**A**: 使用以下策略提高效率:

1. **分批完成,小步迭代**:
   - 不要一次性完成所有 41 章节
   - 先完成高频章节 (lexical_elements, types 核心章节)
   - 每完成 5-10 章节提交一次,获得成就感

2. **使用模板和参考**:
   - 已完成的 `comments.yaml` (35 题) 是很好的参考模板
   - 复制类似难度/题型的题目结构,修改内容
   - 使用工具生成模板,减少手工编写量

3. **团队协作**:
   - 多人分工,每人负责 10-15 章节
   - 统一规范,互相 review
   - 使用 Git 分支隔离工作

4. **优先级排序** (参考 `plan.md` Sprint 2):
   - P1 (高频): lexical_elements (9章), types 核心 10 章
   - P2 (中频): constants (6章), variables (5章)
   - P3 (低频): types 剩余 11 章

5. **质量 vs 速度平衡**:
   - 不追求完美,80 分即可 (质量分数 ≥80%)
   - 先快速完成,后续迭代优化
   - 允许分阶段上线 (如先上线 lexical_elements 和 types)

### Q7: 如何进行人工抽查 (20%)?

**A**: 使用以下流程:

1. **抽样策略**:
   ```bash
   # 使用质量检查工具生成待抽查列表
   go run scripts/quiz_quality_check.go --sample=20% backend/quiz_data/ > sample_list.txt
   
   # 输出示例:
   # 📋 Random Sample (20% of 40 questions):
   # - lexical-keywords-003 (single_choice, medium)
   # - lexical-keywords-007 (multiple_choice, easy)
   # - lexical-keywords-015 (code_output, hard)
   # - lexical-keywords-022 (code_fix, medium)
   # - lexical-keywords-031 (single_choice, easy)
   # - lexical-keywords-038 (multiple_choice, medium)
   # - lexical-keywords-040 (code_output, easy)
   # - lexical-keywords-012 (single_choice, hard)
   ```

2. **人工审查检查点**:
   - ✅ 题干语义是否清晰无歧义?
   - ✅ 所有选项是否合理? (无明显逻辑错误)
   - ✅ 干扰项是否有足够区分度? (不能太明显)
   - ✅ 正确答案是否无争议?
   - ✅ 解析是否帮助理解知识点? (不只是重复答案)
   - ✅ Go 规范引用是否准确?

3. **审查结果处理**:
   ```bash
   # 如果发现问题,修改 YAML 文件后重新验证
   vim backend/quiz_data/lexical_elements/keywords.yaml
   go run scripts/validate_quiz_yaml.go backend/quiz_data/lexical_elements/keywords.yaml
   go run scripts/quiz_quality_check.go backend/quiz_data/lexical_elements/keywords.yaml
   ```

4. **记录审查结果**:
   ```bash
   # 在章节 README 中记录审查信息
   echo "✅ Manual Review: 2025-01-21, Sampled 8/40 questions, 0 issues found" >> backend/quiz_data/lexical_elements/README.md
   ```

---

## 🎯 最佳实践总结

1. **小步迭代**: 每次完成 1-2 章节就提交,不要积压太多
2. **模板先行**: 先用工具生成模板,再填充内容,减少格式错误
3. **边写边验**: 每完成 5-10 题就运行验证工具,及时发现问题
4. **参考规范**: 编写题目时保持 Go 规范页面打开,随时查阅
5. **代码验证**: `code_output` 题型先在 Go Playground 验证再写入
6. **解析优先**: 先写解析再写选项,确保逻辑清晰
7. **质量优先**: 宁可题目少一些,也要保证每道题有学习价值
8. **团队协作**: 互相 review,分享优秀题目模板

---

## 📞 获取帮助

- **文档**: 查看 `research.md`, `data-model.md`, `contracts/quiz-yaml-schema.md`
- **工具帮助**: 运行 `go run scripts/{tool}.go --help`
- **示例参考**: 查看已完成的 `backend/quiz_data/lexical_elements/comments.yaml`
- **问题反馈**: 在项目中创建 Issue 或联系团队成员

---

**祝题库编写顺利! 🚀**
