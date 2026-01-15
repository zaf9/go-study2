# Properties of Types and Values 学习模块

本模块提供 Go 语言"类型和值的属性"章节的学习内容，涵盖值的表示、底层类型、核心类型、类型标识、可赋值性、可表示性和方法集等 7 个子主题。

## 目录结构

```
properties/
├── properties.go          # 核心数据结构与注册机制
├── content.go             # 内容注册入口
├── overview.go            # 章节概览与提纲
├── representation.go      # 值的表示（自包含值与引用值）
├── underlying_type.go     # 底层类型推导规则
├── core_type.go           # 核心类型判定
├── type_identity.go       # 类型标识/相同性判定
├── assignability.go       # 可赋值性规则
├── representability.go    # 可表示性规则
├── method_set.go          # 方法集规则
├── quiz.go                # 综合测验数据
├── search.go              # 搜索索引注册
├── README.md              # 本文件
├── cli/
│   └── menu.go            # CLI 菜单与子主题调度
└── http/
    └── handlers.go        # HTTP 内容/测验/检索输出
```

## 子主题

### 1. 值的表示 (representation)
- 自包含值：数组、结构体
- 引用值：指针、切片、map、channel、函数
- 值的传递与拷贝语义

### 2. 底层类型 (underlying_type)
- 类型定义的底层类型推导
- 类型别名的底层类型
- 类型参数的底层类型判定

### 3. 核心类型 (core_type)
- 非接口类型的核心类型
- 接口类型的核心类型
- bytestring 核心类型

### 4. 类型标识 (type_identity)
- 两类型相同的判定规则
- 命名类型与类型字面量
- 泛型实例化类型的标识

### 5. 可赋值性 (assignability)
- 可赋值性的 6 个基本条件
- 类型参数的 3 个额外条件
- 赋值规则的完整列表

### 6. 可表示性 (representability)
- 常量在特定类型范围内的表示
- 数值常量的可表示性
- 精度与溢出考虑

### 7. 方法集 (method_set)
- 定义类型的方法集
- 指针类型的方法集
- 接口类型的方法集
- 嵌入字段对方法集的影响

## 使用方式

### CLI 模式

```bash
# 进入主菜单
./gostudy

# 选择 Learn > Properties of types and values
# 然后选择具体子主题
```

### HTTP 模式

```bash
# 启动 HTTP 服务
./gostudy -d

# 访问章节概览
curl http://localhost:8080/api/v1/topic/properties

# 访问子主题内容
curl http://localhost:8080/api/v1/topic/properties/representation

# 获取测验题目
curl http://localhost:8080/api/v1/quiz/properties/representation

# 提交测验答案
curl -X POST http://localhost:8080/api/v1/quiz/properties/representation/submit \
  -H "Content-Type: application/json" \
  -d '{"answers": {"q-rep-1": "A"}}'
```

## API 端点

- `GET /api/v1/topic/properties` - 章节概览
- `GET /api/v1/topic/properties/{topic}` - 子主题内容
- `GET /api/v1/quiz/properties/{topic}` - 测验题目
- `POST /api/v1/quiz/properties/{topic}/submit` - 提交测验答案
- `GET /api/v1/search/properties?keyword=xxx` - 搜索关键词

## 数据结构

### PropertyConcept
子主题的概念元信息，包括标题、摘要、规则列表和关键词。

### PropertyRule
类型属性相关的规则或约束，包含规则ID、类型和描述。

### ExampleCase
演示规则的正例或反例，包含代码、期望结果和规则引用。

### QuizItem
测验题目，包含题干、选项、答案和解析。

### LearningProgress
学习进度记录，跟踪已完成的子主题和测验得分。

## 测试

```bash
# 运行单元测试
go test ./backend/src/learning/properties/...

# 运行集成测试
go test ./backend/tests/integration/learning/properties/...

# 运行契约测试
go test ./backend/tests/contract/learning/properties/...
```

## 版本要求

- Go 1.24.5+
- 支持泛型和类型参数相关特性（Go 1.18+）

## 参考资料

- [Go 语言规范 - Properties of types and values](https://go.dev/ref/spec#Properties_of_types_and_values)
- [Go 语言规范 - Types](https://go.dev/ref/spec#Types)
