# HTTP API Contract: Go Constants 学习包

**Feature**: 004-constants-learning  
**Date**: 2025-12-05  
**Phase**: 1 - Design & Contracts

## API 端点

### 1. 获取 Constants 子主题列表

**端点**: `GET /api/v1/topic/constants`

**描述**: 返回 Constants 主题的所有子主题列表

**请求参数**: 无

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "title": "Constants (常量)",
    "subtopics": [
      {"key": "boolean", "name": "Boolean Constants (布尔常量)"},
      {"key": "rune", "name": "Rune Constants (符文常量)"},
      {"key": "integer", "name": "Integer Constants (整数常量)"},
      {"key": "floating_point", "name": "Floating-point Constants (浮点常量)"},
      {"key": "complex", "name": "Complex Constants (复数常量)"},
      {"key": "string", "name": "String Constants (字符串常量)"},
      {"key": "expressions", "name": "Constant Expressions (常量表达式)"},
      {"key": "typed_untyped", "name": "Typed and Untyped Constants (类型化/无类型化常量)"},
      {"key": "conversions", "name": "Conversions (类型转换)"},
      {"key": "builtin_functions", "name": "Built-in Functions (内置函数)"},
      {"key": "iota", "name": "Iota (iota 特性)"},
      {"key": "implementation_restrictions", "name": "Implementation Restrictions (实现限制)"}
    ]
  }
}
```

---

### 2. 获取特定子主题内容

**端点**: `GET /api/v1/topic/constants/:subtopic`

**描述**: 返回指定子主题的详细学习内容

**路径参数**:
- `subtopic`: 子主题键,可选值见下表

| subtopic | 说明 |
|----------|------|
| boolean | 布尔常量 |
| rune | 符文常量 |
| integer | 整数常量 |
| floating_point | 浮点常量 |
| complex | 复数常量 |
| string | 字符串常量 |
| expressions | 常量表达式 |
| typed_untyped | 类型化/无类型化常量 |
| conversions | 类型转换 |
| builtin_functions | 内置函数 |
| iota | iota 特性 |
| implementation_restrictions | 实现限制 |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "title": "Boolean Constants (布尔常量)",
    "description": "布尔常量表示真值,只有两个预声明的常量: true 和 false。布尔常量是无类型化的,默认类型为 bool。",
    "syntax": "const name [type] = true | false",
    "examples": [
      {
        "title": "基本布尔常量声明",
        "code": "package main\n\nimport \"fmt\"\n\nconst (\n    enabled  = true\n    disabled = false\n)\n\nfunc main() {\n    fmt.Println(\"enabled:\", enabled)\n    fmt.Println(\"disabled:\", disabled)\n}",
        "output": "enabled: true\ndisabled: false",
        "explanation": "声明两个布尔常量,分别表示启用和禁用状态"
      },
      {
        "title": "类型化布尔常量",
        "code": "...",
        "output": "...",
        "explanation": "..."
      }
    ],
    "common_errors": "常见错误:\n1. 尝试将非布尔值赋给布尔常量\n2. 混淆常量和变量的使用场景"
  }
}
```

**错误响应** (404 Not Found):
```json
{
  "code": 404,
  "message": "subtopic not found: xyz",
  "data": null
}
```

---

## 响应格式规范

### 标准响应结构

```go
type APIResponse struct {
    Code    int         `json:"code"`    // 状态码: 0=成功, 非0=错误
    Message string      `json:"message"` // 消息
    Data    interface{} `json:"data"`    // 数据
}
```

### 状态码定义

| Code | HTTP Status | 说明 |
|------|-------------|------|
| 0 | 200 | 成功 |
| 404 | 404 | 资源不存在 |
| 500 | 500 | 服务器内部错误 |

---

## Handler 实现

### Handler 结构

```go
package handler

import (
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/frame/g"
)

// GetConstantsMenu 获取 Constants 子主题列表
func (h *Handler) GetConstantsMenu(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "code": 0,
        "message": "success",
        "data": g.Map{
            "title": "Constants (常量)",
            "subtopics": []g.Map{
                {"key": "boolean", "name": "Boolean Constants (布尔常量)"},
                // ... 其他子主题
            },
        },
    })
}

// GetConstantsContent 获取特定子主题内容
func (h *Handler) GetConstantsContent(r *ghttp.Request) {
    subtopic := r.Get("subtopic").String()
    
    // 根据 subtopic 构造内容
    content, err := buildConstantsContent(subtopic)
    if err != nil {
        r.Response.WriteJson(g.Map{
            "code": 404,
            "message": "subtopic not found: " + subtopic,
            "data": nil,
        })
        return
    }
    
    r.Response.WriteJson(g.Map{
        "code": 0,
        "message": "success",
        "data": content,
    })
}
```

### 路由注册

```go
// router.go
func RegisterRoutes(s *ghttp.Server) {
    h := handler.New()
    
    s.Group("/api/v1", func(group *ghttp.RouterGroup) {
        group.Middleware(middleware.Format)
        
        // Constants 路由
        group.ALL("/topic/constants", h.GetConstantsMenu)
        group.ALL("/topic/constants/:subtopic", h.GetConstantsContent)
    })
}
```

---

## 性能要求

- **响应时间**: <100ms (正常负载 100 并发请求)
- **并发能力**: 支持 1000 并发请求无错误
- **错误率**: 0%

---

## 测试要求

### 单元测试

```go
func TestGetConstantsMenu(t *testing.T) {
    // 测试菜单端点
    // 验证响应结构
    // 验证子主题列表完整性
}

func TestGetConstantsContent(t *testing.T) {
    // 测试内容端点
    // 验证所有 subtopic 可访问
    // 验证响应包含必需字段
}

func TestGetConstantsContent_NotFound(t *testing.T) {
    // 测试不存在的 subtopic
    // 验证返回 404 错误
}
```

### 集成测试

```go
func TestConstantsAPI_E2E(t *testing.T) {
    // 启动测试服务器
    // 发送真实 HTTP 请求
    // 验证端到端流程
}
```

### 性能测试

```bash
# 使用 wrk 进行压力测试
wrk -t4 -c100 -d30s http://localhost:8080/api/v1/topic/constants/boolean

# 验收标准:
# - p95 响应时间 <100ms
# - 错误率 0%
```

---

## 验收标准

- [ ] 所有端点返回正确的 JSON 格式
- [ ] 12 个子主题都可通过 API 访问
- [ ] 不存在的 subtopic 返回 404 错误
- [ ] 响应时间满足性能要求
- [ ] 并发测试无错误
- [ ] 与 lexical_elements API 风格一致
