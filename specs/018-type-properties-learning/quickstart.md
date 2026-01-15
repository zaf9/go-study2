# Quickstart: Go 类型属性章节学习方案

**Feature**: 018-type-properties-learning  
**Created**: 2026-01-15

## 快速开始

### 1. 验证环境

```bash
cd backend
go version  # 确认 Go 1.24+
go mod tidy
```

### 2. 运行测试

```bash
# 运行类型属性模块的单元测试
cd backend
go test ./src/learning/properties/... -v

# 运行全部测试
go test ./... -v
```

### 3. CLI 模式体验

```bash
cd backend
go run main.go
# 在主菜单中选择 "Properties of types and values (类型属性)"
```

### 4. HTTP 模式体验

```bash
# 启动 HTTP 服务
cd backend
go run main.go --http

# 访问菜单
curl http://localhost:8080/api/v1/topic/properties

# 访问子主题内容
curl http://localhost:8080/api/v1/topic/properties/representation
curl http://localhost:8080/api/v1/topic/properties/underlying_type
curl http://localhost:8080/api/v1/topic/properties/type_identity

# 获取测验题目
curl http://localhost:8080/api/v1/topic/properties/representation/quiz

# 提交测验答案
curl -X POST http://localhost:8080/api/v1/topic/properties/representation/quiz \
  -H "Content-Type: application/json" \
  -d '{"answers":{"q-rep-1":"A","q-rep-2":"B"}}'

# 搜索关键词
curl http://localhost:8080/api/v1/topic/properties/search?keyword=底层类型
```

## 核心文件说明

| 文件 | 说明 |
|------|------|
| `properties.go` | 核心数据结构和注册/查询函数 |
| `content.go` | 内容注册入口（init 调用） |
| `overview.go` | 章节概览和打印提纲 |
| `representation.go` | 值的表示子主题内容 |
| `underlying_type.go` | 底层类型子主题内容 |
| `core_type.go` | 核心类型子主题内容 |
| `type_identity.go` | 类型标识子主题内容 |
| `assignability.go` | 可赋值性子主题内容 |
| `representability.go` | 可表示性子主题内容 |
| `method_set.go` | 方法集子主题内容 |
| `quiz.go` | 综合测验数据 |
| `search.go` | 搜索索引注册 |
| `cli/menu.go` | CLI 菜单交互 |
| `http/handlers.go` | HTTP 处理器 |

## 添加新内容

每个子主题文件遵循相同模式：

```go
package properties

func registerXxx() {
    _ = RegisterContent(TopicXxx, TopicContent{
        Concept: PropertyConcept{
            ID:        "xxx",
            Category:  "basic",
            Title:     "中文标题",
            Summary:   "摘要说明...",
            GoVersion: "1.24",
            Rules:     []string{"规则1", "规则2"},
            Keywords:  []string{"keyword1", "keyword2"},
        },
        Rules: []PropertyRule{...},
        Examples: []ExampleCase{...},
        QuizItems: []QuizItem{...},
        References: []ReferenceIndex{...},
    })
}

func XxxOutline() []string {
    return []string{
        "提纲要点1",
        "提纲要点2",
    }
}
```

## 测试覆盖率

```bash
cd backend
go test ./src/learning/properties/... -cover

# 生成详细覆盖率报告
go test ./src/learning/properties/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

目标：覆盖率 >= 80%
