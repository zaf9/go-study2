# API Contract: Go 类型属性章节学习

**Feature**: 018-type-properties-learning  
**Created**: 2026-01-15

## CLI Interface

### 菜单入口

```
DisplayMenu(stdin io.Reader, stdout, stderr io.Writer)
```

### 菜单结构

```
Properties of types and values 学习菜单
-----------------------------------------
0. Representation (值的表示)
1. Underlying Type (底层类型)
2. Core Type (核心类型)
3. Type Identity (类型标识)
4. Assignability (可赋值性)
5. Representability (可表示性)
6. Method Set (方法集)
o. 打印提纲
quiz. 综合测验
search <keyword>. 关键词检索
q. 返回上级菜单

请输入您的选择: 
```

### 子主题内容输出格式

```
[concept-id] 标题
摘要说明

规则:
- 规则1
- 规则2

示例:
• 示例标题
示例代码
=> 期望输出

测验: 输入选项字母作答，输入 q 结束测验。
question-id: 题干
A) 选项A
B) 选项B
答案: 

得分: X/Y
- question-id: 正确/错误 (答案: X) - 解析
```

---

## HTTP API

### 菜单接口

**GET** `/api/v1/topic/properties`

**Response (JSON)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "representation",
        "title": "值的表示",
        "path": "/api/v1/topic/properties/representation"
      },
      {
        "id": "underlying_type",
        "title": "底层类型",
        "path": "/api/v1/topic/properties/underlying_type"
      },
      {
        "id": "core_type",
        "title": "核心类型",
        "path": "/api/v1/topic/properties/core_type"
      },
      {
        "id": "type_identity",
        "title": "类型标识",
        "path": "/api/v1/topic/properties/type_identity"
      },
      {
        "id": "assignability",
        "title": "可赋值性",
        "path": "/api/v1/topic/properties/assignability"
      },
      {
        "id": "representability",
        "title": "可表示性",
        "path": "/api/v1/topic/properties/representability"
      },
      {
        "id": "method_set",
        "title": "方法集",
        "path": "/api/v1/topic/properties/method_set"
      }
    ]
  }
}
```

### 内容接口

**GET** `/api/v1/topic/properties/{topic}`

**Parameters**:
- `topic`: 子主题ID（representation, underlying_type, core_type, type_identity, assignability, representability, method_set）
- `format`: 响应格式（json/html，默认 json）

**Response (JSON)**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "concept": {
      "id": "representation",
      "category": "basic",
      "title": "值的表示",
      "summary": "Go 中值的表示方式...",
      "goVersion": "1.24",
      "rules": ["规则1", "规则2"],
      "keywords": ["self-contained", "reference"]
    },
    "rules": [
      {
        "ruleId": "PR-REP-001",
        "conceptId": "representation",
        "ruleType": "representation",
        "description": "预声明类型的值是自包含的..."
      }
    ],
    "examples": [
      {
        "id": "ex-rep-1",
        "conceptId": "representation",
        "title": "自包含值示例",
        "code": "var arr [3]int = [3]int{1, 2, 3}...",
        "expectedOutput": "[1 2 3]",
        "isValid": true,
        "ruleRef": "PR-REP-001"
      }
    ],
    "quizItems": [
      {
        "id": "q-rep-1",
        "conceptId": "representation",
        "stem": "以下哪种类型的值是自包含的？",
        "options": ["数组", "切片", "map", "channel"],
        "difficulty": "easy"
      }
    ]
  }
}
```

**Error Response**:
```json
{
  "code": 404,
  "message": "不支持的类型属性子主题",
  "data": null
}
```

### 测验接口

**GET** `/api/v1/topic/properties/{topic}/quiz`

返回指定子主题的测验题目（不含答案）。

**POST** `/api/v1/topic/properties/{topic}/quiz`

**Request Body**:
```json
{
  "answers": {
    "q-rep-1": "A",
    "q-rep-2": "B"
  }
}
```

**Response**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "score": 1,
    "total": 2,
    "details": [
      {
        "id": "q-rep-1",
        "correct": true,
        "answer": "A",
        "explanation": "数组的值包含完整数据副本...",
        "ruleRef": "PR-REP-001"
      },
      {
        "id": "q-rep-2",
        "correct": false,
        "answer": "B",
        "explanation": "...",
        "ruleRef": "PR-REP-002"
      }
    ]
  }
}
```

### 搜索接口

**GET** `/api/v1/topic/properties/search?keyword={keyword}`

**Response**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "results": [
      {
        "keyword": "底层类型",
        "conceptId": "underlying_type",
        "summary": "底层类型是类型系统的基础...",
        "anchors": {
          "http": "/api/v1/topic/properties/underlying_type",
          "cli": "properties > underlying_type"
        }
      }
    ]
  }
}
```

---

## 路由注册

需要在路由器中添加以下路由：

```go
// backend/internal/app/http_server/router.go

// Properties of types and values 章节
group.GET("/topic/properties", properties_http.MenuHandler)
group.GET("/topic/properties/:topic", properties_http.ContentHandler)
group.GET("/topic/properties/:topic/quiz", properties_http.QuizHandler)
group.POST("/topic/properties/:topic/quiz", properties_http.EvaluateQuizHandler)
group.GET("/topic/properties/search", properties_http.SearchHandler)
```

---

## 主菜单集成

需要在主菜单中添加 Properties of types and values 入口：

```go
// CLI 主菜单
{
    Title: "Properties of types and values (类型属性)",
    Action: properties_cli.DisplayMenu,
}

// HTTP /api/v1/topics 响应
{
    "id": "properties",
    "title": "Properties of types and values",
    "path": "/api/v1/topic/properties"
}
```
