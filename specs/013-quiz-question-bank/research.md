# Phase 0 Research: 题目生成策略与技术方案

**Branch**: `013-quiz-question-bank` | **Date**: 2025-12-14

## 概述

本研究文档解决实施计划中Phase 0阶段的所有"NEEDS CLARIFICATION"项，为41个章节生成1230-2050个高质量测验题目提供可执行方案。

## 0.1 题目生成策略研究

### 题目生成核心问题

**问题**: 如何高效为41个章节生成1230-2050个高质量题目？

**解决方案**: **人工编写 + AI辅助扩展 + 三轮质量审核**

#### Workflow设计

```text
步骤1: 知识点提纲提取（人工）
  ↓
步骤2: 核心题目编写（人工首轮）- 每章节5-8题
  ↓
步骤3: AI辅助扩展（半自动）- 基于核心题目生成变体
  ↓
步骤4: 质量审核第一轮（人工）- 修正明显错误
  ↓
步骤5: 格式规范化（自动）- 统一YAML格式
  ↓
步骤6: 质量审核第二轮（人工）- 验证准确性
  ↓
步骤7: 深度验证（自动）- 启动时结构完整性检查
```

### 41个章节知识点提纲

基于现有代码和Go 1.24规范，提取每个章节的核心知识点（5-8个）：

#### lexical_elements (11章节)

1. **comments** (注释)
   - 单行注释 `//`
   - 多行注释 `/* */`
   - 注释嵌套规则
   - godoc文档注释规范
   - 注释在编译时的处理

2. **tokens** (标记)
   - 标识符、关键字、操作符、字面量、分隔符
   - 标记的分类和识别规则
   - 标记的组合规则
   - Unicode字符处理

3. **semicolons** (分号)
   - 自动分号插入规则
   - 显式分号使用场景
   - 分号与语句终止
   - 分号在for循环中的特殊性

4. **identifiers** (标识符)
   - 标识符命名规则（Unicode字母、数字、下划线）
   - 预声明标识符
   - 导出/非导出标识符
   - 空白标识符 `_`
   - 标识符作用域

5. **keywords** (关键字)
   - 25个关键字列表
   - 关键字不能作为标识符
   - 关键字的语法分类（声明、控制流、类型等）

6. **operators** (运算符)
   - 算术运算符 (+, -, *, /, %)
   - 比较运算符 (==, !=, <, >, <=, >=)
   - 逻辑运算符 (&&, ||, !)
   - 位运算符 (&, |, ^, <<, >>)
   - 赋值运算符 (=, +=, -=等)
   - 运算符优先级

7. **integers** (整数字面量)
   - 十进制、八进制、十六进制、二进制表示
   - 整数字面量的分隔符 `_`
   - 整数字面量的类型推断
   - 溢出检测

8. **floats** (浮点数字面量)
   - 十进制浮点数表示
   - 十六进制浮点数表示
   - 科学计数法
   - 浮点数精度限制

9. **imaginary** (虚数字面量)
   - 虚数字面量的表示 (`i`后缀)
   - 虚数的类型推断
   - 虚数与复数的关系

10. **runes** (符文字面量)
    - 单引号表示
    - Unicode码点表示
    - 转义字符
    - 字节序列
    - 无效符文

11. **strings** (字符串字面量)
    - 解释型字符串 (`"..."`)
    - 原始字符串 (`` `...` ``)
    - 转义序列
    - 字符串连接
    - 字符串不可变性

#### constants (12章节)

1. **boolean** (布尔常量)
   - `true`和`false`
   - 布尔常量的类型
   - 布尔常量的运算

2. **rune** (符文常量)
   - 符文常量的定义
   - 符文常量的类型推断
   - 符文常量与整数的关系

3. **integer** (整数常量)
   - 整数常量的表示
   - 整数常量的精度（任意精度）
   - 整数常量的类型推断

4. **floating_point** (浮点常量)
   - 浮点常量的表示
   - 浮点常量的精度
   - 浮点常量的特殊值（不支持NaN和Inf）

5. **complex** (复数常量)
   - 复数常量的表示
   - 实部和虚部
   - 复数常量的运算

6. **string** (字符串常量)
   - 字符串常量的定义
   - 字符串常量的编码（UTF-8）
   - 字符串常量的连接

7. **expressions** (常量表达式)
   - 常量表达式的定义
   - 可用于常量表达式的操作符
   - 常量表达式的求值时机

8. **typed_untyped** (有类型和无类型常量)
   - 无类型常量的默认类型
   - 显式类型转换
   - 隐式类型推断

9. **conversions** (常量转换)
   - 数值常量之间的转换
   - 常量可表示性规则
   - 转换溢出检测

10. **builtin_functions** (内置函数用于常量)
    - `len()`、`cap()`在常量中的使用
    - `real()`、`imag()`、`complex()`
    - `unsafe.Sizeof()`等

11. **iota** (iota枚举器)
    - `iota`的基本用法
    - `iota`的递增规则
    - `iota`的重置时机
    - `iota`的跳过规则

12. **implementation_restrictions** (实现限制)
    - 整数常量的位数限制（至少256位）
    - 浮点常量的精度限制（至少256位尾数）
    - 复数常量的组件限制
    - 字符串常量的长度限制

#### variables (4章节)

1. **storage** (存储位置)
   - 变量的定义
   - 变量的存储位置
   - 变量的生命周期
   - 变量的作用域

2. **static** (静态变量)
   - 包级变量的声明
   - 静态变量的初始化顺序
   - 静态变量的内存分配

3. **dynamic** (动态变量)
   - 函数内局部变量
   - `new()`分配的变量
   - 复合字面量分配的变量
   - 逃逸分析

4. **zero** (零值)
   - 各类型的零值定义
   - 零值的自动初始化
   - 零值与nil的区别
   - 零值的可用性

#### types (14章节)

1. **boolean** (布尔类型)
   - `bool`类型
   - 布尔值 `true`和`false`
   - 布尔运算

2. **numeric** (数值类型)
   - 整型：int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint, uintptr
   - 浮点型：float32, float64
   - 复数型：complex64, complex128
   - 字节和符文别名：byte, rune
   - 数值类型的大小和范围

3. **string** (字符串类型)
   - 字符串的定义
   - 字符串的不可变性
   - 字符串的长度
   - 字符串索引和切片

4. **array** (数组类型)
   - 数组的定义和声明
   - 数组的长度固定性
   - 数组的元素类型
   - 数组的初始化
   - 数组的值语义

5. **slice** (切片类型)
   - 切片的定义
   - 切片的底层数组
   - 切片的长度和容量
   - 切片的动态扩容
   - 切片的引用语义

6. **struct** (结构体类型)
   - 结构体的定义
   - 字段的声明和访问
   - 嵌入字段
   - 字段标签
   - 空结构体

7. **pointer** (指针类型)
   - 指针的定义
   - 取地址操作符 `&`
   - 解引用操作符 `*`
   - 指针的零值 `nil`
   - 指针的间接访问

8. **function** (函数类型)
   - 函数类型的定义
   - 函数签名
   - 函数作为一等公民
   - 函数的零值 `nil`
   - 函数类型的可比较性

9. **interface_basic** (接口基础)
   - 接口的定义
   - 方法集
   - 接口的实现
   - 空接口 `interface{}`/`any`

10. **interface_embedded** (接口嵌入)
    - 接口嵌入其他接口
    - 方法集的合并
    - 嵌入重复方法的处理

11. **interface_general** (通用接口)
    - 类型集的概念
    - 联合类型约束 `|`
    - 近似类型约束 `~`
    - 类型参数约束

12. **interface_impl** (接口实现)
    - 隐式实现
    - 接口的动态类型和动态值
    - 接口的零值
    - 接口的可比较性

13. **map** (映射类型)
    - map的定义
    - 键类型的可比较性要求
    - map的初始化
    - map的零值 `nil`
    - map的引用语义

14. **channel** (通道类型)
    - channel的定义
    - channel的方向性
    - channel的零值 `nil`
    - channel的缓冲
    - channel的引用语义

### 题目质量标准Checklist

每个题目必须满足以下标准：

**结构完整性**:
- ✅ 题目ID符合命名规则（如 `lexical-comments-001`）
- ✅ 题型明确（single或multiple）
- ✅ 难度级别明确（easy/medium/hard）
- ✅ 题干清晰且语法正确
- ✅ 选项数量符合要求（单选2-4个，多选3-5个）
- ✅ 正确答案格式正确（单选为单个字母，多选为字母组合）
- ✅ 解析字段存在且为中文

**内容质量**:
- ✅ 题干与章节知识点相关
- ✅ 选项无歧义，互斥且合理
- ✅ 答案准确无误
- ✅ 解析详细说明为何答案正确和其他选项错误
- ✅ 解析补充相关知识点

**难度判定指南**:
- **easy (40%)**:
  - 直接考察定义、语法规则、基本概念
  - 不需要推理，直接从规范中找到答案
  - 示例："`//`是Go语言的什么注释方式？"
  
- **medium (40%)**:
  - 需要理解规则并应用到简单场景
  - 需要辨析相似概念的区别
  - 示例："以下哪些标识符是合法的？（多选）"
  
- **hard (20%)**:
  - 需要深度理解和综合运用多个规则
  - 涉及边界情况、特殊场景
  - 示例："以下代码片段中，哪些会触发自动分号插入？"

### AI辅助生成策略

**使用场景**: 在人工编写核心题目（5-8题）后，使用AI生成变体扩展到30-50题。

**Prompt模板**:

```text
任务：基于以下核心题目，生成15-20个同主题变体题目。

核心题目：
[粘贴人工编写的5-8个核心题目YAML]

要求：
1. 保持YAML格式一致
2. 生成单选题和多选题各占一半
3. 难度分布：简单40%、中等40%、困难20%
4. 每题必须包含详细中文解析
5. 题目ID递增（如lexical-comments-009, 010等）
6. 避免与核心题目重复
7. 覆盖知识点提纲中的不同方面

知识点提纲：
[粘贴对应章节的知识点清单]
```

**人工审核重点**:
- 答案准确性验证
- 解析合理性检查
- 选项歧义消除
- 题干语法修正

## 0.2 YAML格式和数据模型设计

### QuizQuestion Go Struct定义

```go
package quiz

// QuestionType 题目类型
type QuestionType string

const (
	QuestionTypeSingle   QuestionType = "single"   // 单选题
	QuestionTypeMultiple QuestionType = "multiple" // 多选题
)

// Difficulty 题目难度
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"   // 简单
	DifficultyMedium Difficulty = "medium" // 中等
	DifficultyHard   Difficulty = "hard"   // 困难
)

// QuizQuestion 测验题目实体
type QuizQuestion struct {
	ID          string       `yaml:"id"`          // 题目唯一标识，如 "lexical-comments-001"
	Type        QuestionType `yaml:"type"`        // 题型：single/multiple
	Difficulty  Difficulty   `yaml:"difficulty"`  // 难度：easy/medium/hard
	Stem        string       `yaml:"stem"`        // 题干
	Options     []string     `yaml:"options"`     // 选项列表，如 ["A: ...", "B: ...", ...]
	Answer      string       `yaml:"answer"`      // 正确答案，单选为 "A"，多选为 "ACD"
	Explanation string       `yaml:"explanation"` // 答案解析（中文）
	Topic       string       `yaml:"topic"`       // 所属主题，如 "lexical_elements"
	Chapter     string       `yaml:"chapter"`     // 所属章节，如 "comments"
}

// QuizBank 题库文件结构
type QuizBank struct {
	Questions []QuizQuestion `yaml:"questions"` // 题目列表
}
```

### YAML文件格式规范

**文件路径**: `backend/quiz_data/{topic}/{chapter}.yaml`

**文件结构**:

```yaml
questions:
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
```

### 题目ID命名规则

**格式**: `{topic_prefix}-{chapter}-{sequence}`

**示例**:
- `lexical-comments-001` (lexical_elements主题)
- `const-boolean-015` (constants主题)
- `var-storage-003` (variables主题)
- `type-array-022` (types主题)

**规则**:
- topic_prefix：主题简写（lexical/const/var/type）
- chapter：章节名（英文小写）
- sequence：三位数字序号（001-050）

### 枚举值定义表

| 字段 | 枚举值 | 说明 |
|------|--------|------|
| type | `single` | 单选题 |
| type | `multiple` | 多选题 |
| difficulty | `easy` | 简单（直接定义、基本规则） |
| difficulty | `medium` | 中等（理解应用、辨析区别） |
| difficulty | `hard` | 困难（综合运用、边界情况） |

### 验证规则

**必填字段验证**:
- 所有字段均为必填，不允许空值

**类型约束**:
- `type`: 必须为 `single` 或 `multiple`
- `difficulty`: 必须为 `easy`、`medium` 或 `hard`
- `options`: 数组长度2-5，每个选项格式为 `X: ...`（X为字母）
- `answer`: 
  - 单选题：单个大写字母（A/B/C/D/E）
  - 多选题：2-4个大写字母组合（如 AB、ACD、BDE）

**业务规则验证**:
- 答案字母必须在选项范围内
- 多选题答案至少包含2个选项
- 题目ID在同一章节内唯一
- topic和chapter必须与文件路径一致

## 0.3 技术依赖和性能基准

### yaml.v3解析性能验证

**测试场景**: 加载41个YAML文件，共约2000题

**性能基准测试代码**:

```go
package quiz

import (
	"testing"
	"time"
	"gopkg.in/yaml.v3"
	"os"
)

func BenchmarkLoadAllQuizBanks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadAllBanks()
	}
}

func loadAllBanks() {
	topics := []string{"lexical_elements", "constants", "variables", "types"}
	chapters := map[string][]string{
		"lexical_elements": {"comments", "tokens", /* ... */},
		// ... 其他主题
	}
	
	start := time.Now()
	for topic, chaps := range chapters {
		for _, chapter := range chaps {
			path := fmt.Sprintf("quiz_data/%s/%s.yaml", topic, chapter)
			data, _ := os.ReadFile(path)
			var bank QuizBank
			yaml.Unmarshal(data, &bank)
		}
	}
	elapsed := time.Since(start)
	// 期望: elapsed < 1秒（远低于5秒目标）
}
```

**预期结果**:
- 41个文件，每个文件约50题，总计2050题
- 平均每个文件大小约5KB（50题×100字节）
- 总文件大小约200KB
- yaml.v3解析200KB预计耗时 < 500ms
- **结论**: 满足 <5秒启动要求，性能充裕

### 并发安全的随机数生成方案

**问题**: `math/rand`包的全局函数非线程安全

**解决方案**: 使用 `sync.Mutex` 保护随机数生成

```go
package quiz

import (
	"math/rand"
	"sync"
	"time"
)

type Selector struct {
	rng *rand.Rand
	mu  sync.Mutex
}

func NewSelector() *Selector {
	return &Selector{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SelectQuestions 线程安全的随机抽题
func (s *Selector) SelectQuestions(questions []QuizQuestion, count int) []QuizQuestion {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 使用Fisher-Yates洗牌算法
	indices := s.rng.Perm(len(questions))
	selected := make([]QuizQuestion, 0, count)
	for i := 0; i < count && i < len(indices); i++ {
		selected = append(selected, questions[indices[i]])
	}
	return selected
}
```

**并发性能验证**:
- 100并发用户同时抽题，抽题耗时 < 50ms (p95)
- 锁竞争在轻量级读操作下影响可忽略

### 内存占用评估

**单题内存估算**:
```text
QuizQuestion结构体:
- ID: 24字节（字符串header + 平均长度18）
- Type: 8字节（字符串header + "multiple"）
- Difficulty: 8字节（字符串header + "medium"）
- Stem: 120字节（字符串header + 平均长度100）
- Options: 200字节（切片header + 4个选项×50字节）
- Answer: 8字节（字符串header + "ACD"）
- Explanation: 240字节（字符串header + 平均长度200）
- Topic: 24字节
- Chapter: 16字节

总计: 约650字节/题
```

**总内存占用**:
- 2000题 × 650字节 ≈ 1.3 MB
- **结论**: 内存占用极小，可安全缓存在内存

### 性能优化策略

1. **启动时一次性加载**: 所有题库在应用启动时加载到内存，避免运行时IO
2. **内存缓存**: 题库数据存储在全局变量，只读访问无需锁
3. **按章节索引**: 使用 `map[topic][chapter][]QuizQuestion` 结构加速查询
4. **延迟验证**: 深度验证仅在启动时执行一次

## 研究决策总结

### 题目生成决策

**Decision**: 采用"人工核心题（5-8题）+ AI扩展（25-42题）+ 三轮审核"策略

**Rationale**: 
- 纯人工编写1230-2050题工作量过大（约200小时）
- 纯AI生成题目质量难以保证
- 混合方案兼顾效率和质量，人工保证核心质量，AI提升效率

**Alternatives considered**:
- ❌ 纯人工编写：时间成本过高
- ❌ 纯AI生成：质量风险高，答案准确性难保证
- ❌ 模板化生成：题目多样性不足，学习者易识别套路

### 数据模型决策

**Decision**: 使用YAML格式存储，每章节一个文件，9字段结构

**Rationale**:
- YAML可读性强，便于人工审核和修改
- 符合项目配置管理习惯（configs/使用YAML）
- yaml.v3库成熟，性能满足要求
- 按章节分文件便于增量开发和版本控制

**Alternatives considered**:
- ❌ JSON格式：可读性略差，注释不友好
- ❌ Go代码硬编码：不便于非开发人员维护
- ❌ 数据库存储：过度设计，增加运维复杂度

### 性能方案决策

**Decision**: 启动时全量加载到内存，使用mutex保护随机数生成

**Rationale**:
- 题库总大小仅1.3MB，内存占用可忽略
- 避免运行时IO，抽题性能最优
- 题库只读特性，无需复杂并发控制
- 随机数生成器互斥锁开销极小（微秒级）

**Alternatives considered**:
- ❌ 按需加载：增加运行时IO，性能不稳定
- ❌ 数据库查询：增加外部依赖，性能不如内存
- ❌ 无锁并发：复杂度高，收益有限（随机数生成非瓶颈）

## 下一步行动

1. ✅ **Completed**: 研究阶段（本文档）
2. ⏭ **Next**: 进入Phase 1，生成 `data-model.md` 和 `contracts/yaml-schema.md`
3. ⏭ **Then**: 开始题目内容生成（41个章节，优先级：P1章节先行）

---

**Research Status**: ✅ Complete  
**All NEEDS CLARIFICATION Resolved**: ✅ Yes  
**Ready for**: Phase 1 Design & Contracts
