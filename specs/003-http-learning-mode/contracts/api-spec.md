# API Contracts: HTTP学习模式

**Feature**: 003-http-learning-mode  
**Date**: 2025-12-04  
**Version**: 1.0.0

---

## 概述

本文档定义HTTP学习模式的所有API接口契约。所有接口均使用POST方法，支持JSON和HTML两种响应格式。

### 基础信息

- **Base URL**: `http://{host}:{port}`
- **API Version**: v1
- **Content-Type**: `application/json` (请求) / `application/json` 或 `text/html` (响应)
- **Character Encoding**: UTF-8

### 通用响应格式

#### 成功响应 (JSON)

```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "timestamp": 1701676800
}
```

#### 错误响应 (JSON)

```json
{
  "code": 404,
  "message": "章节不存在",
  "error": "chapter 'nonexistent' not found",
  "timestamp": 1701676800
}
```

#### 格式参数

所有接口支持通过查询参数`format`指定响应格式：

- `?format=json` - 返回JSON格式（默认）
- `?format=html` - 返回HTML格式

---

## 接口列表

### 1. 获取主题列表

获取所有可用的学习主题。

**Endpoint**: `POST /api/v1/topics`

#### 请求

**URL**: `/api/v1/topics?format=json`

**Method**: `POST`

**Headers**:
```
Content-Type: application/json
```

**Body**: 无需请求体（或空JSON对象`{}`）

#### 响应

**Status Code**: `200 OK`

**JSON格式**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "topics": [
      {
        "id": "lexical_elements",
        "title": "词法元素 (Lexical Elements)",
        "description": "Go语言的基本词法元素，包括注释、标记、分号等",
        "chapters": [
          {
            "id": "comments",
            "title": "注释 (Comments)",
            "path": "/api/v1/topic/lexical_elements/comments"
          },
          {
            "id": "tokens",
            "title": "标记 (Tokens)",
            "path": "/api/v1/topic/lexical_elements/tokens"
          }
        ],
        "order": 1
      }
    ],
    "total": 1
  },
  "timestamp": 1701676800
}
```

**HTML格式** (`?format=html`):

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>Go学习主题列表</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .topic { margin-bottom: 30px; padding: 15px; border: 1px solid #ddd; }
        .topic h2 { color: #0066cc; }
        .chapters { margin-left: 20px; }
        .chapter-link { display: block; margin: 5px 0; color: #0066cc; }
    </style>
</head>
<body>
    <h1>Go学习主题列表</h1>
    <div class="topic">
        <h2>词法元素 (Lexical Elements)</h2>
        <p>Go语言的基本词法元素，包括注释、标记、分号等</p>
        <div class="chapters">
            <a class="chapter-link" href="/api/v1/topic/lexical_elements/comments">注释 (Comments)</a>
            <a class="chapter-link" href="/api/v1/topic/lexical_elements/tokens">标记 (Tokens)</a>
        </div>
    </div>
</body>
</html>
```

#### 错误响应

无特定错误（此接口总是返回可用主题列表，即使为空）

---

### 2. 获取Lexical Elements菜单

获取"Lexical Elements"主题的章节菜单。

**Endpoint**: `POST /api/v1/topic/lexical_elements`

#### 请求

**URL**: `/api/v1/topic/lexical_elements?format=json`

**Method**: `POST`

**Headers**:
```
Content-Type: application/json
```

**Body**: 无需请求体

#### 响应

**Status Code**: `200 OK`

**JSON格式**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "lexical_elements",
    "title": "词法元素 (Lexical Elements)",
    "description": "Go语言的基本词法元素",
    "chapters": [
      {
        "id": "comments",
        "title": "注释 (Comments)",
        "path": "/api/v1/topic/lexical_elements/comments"
      },
      {
        "id": "tokens",
        "title": "标记 (Tokens)",
        "path": "/api/v1/topic/lexical_elements/tokens"
      },
      {
        "id": "semicolons",
        "title": "分号 (Semicolons)",
        "path": "/api/v1/topic/lexical_elements/semicolons"
      },
      {
        "id": "identifiers",
        "title": "标识符 (Identifiers)",
        "path": "/api/v1/topic/lexical_elements/identifiers"
      }
    ]
  },
  "timestamp": 1701676800
}
```

**HTML格式** (`?format=html`):

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>词法元素 (Lexical Elements)</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #0066cc; }
        .chapter-list { list-style: none; padding: 0; }
        .chapter-item { margin: 10px 0; padding: 10px; background: #f9f9f9; }
        .chapter-link { color: #0066cc; text-decoration: none; font-size: 18px; }
    </style>
</head>
<body>
    <h1>词法元素 (Lexical Elements)</h1>
    <p>Go语言的基本词法元素</p>
    <ul class="chapter-list">
        <li class="chapter-item">
            <a class="chapter-link" href="/api/v1/topic/lexical_elements/comments">注释 (Comments)</a>
        </li>
        <li class="chapter-item">
            <a class="chapter-link" href="/api/v1/topic/lexical_elements/tokens">标记 (Tokens)</a>
        </li>
        <!-- 更多章节... -->
    </ul>
</body>
</html>
```

#### 错误响应

**404 Not Found** - 主题不存在

```json
{
  "code": 404,
  "message": "主题不存在",
  "error": "topic 'lexical_elements' not found",
  "timestamp": 1701676800
}
```

---

### 3. 获取章节内容

获取指定章节的详细学习内容。

**Endpoint**: `POST /api/v1/topic/lexical_elements/{chapter}`

#### 请求

**URL**: `/api/v1/topic/lexical_elements/comments?format=json`

**Method**: `POST`

**Headers**:
```
Content-Type: application/json
```

**Path Parameters**:
- `chapter` (string, required): 章节ID，如 `comments`, `tokens`

**Body**: 无需请求体

#### 响应

**Status Code**: `200 OK`

**JSON格式**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "comments",
    "title": "注释 (Comments)",
    "topic_id": "lexical_elements",
    "content": "# 注释\n\n## 概述\n\nGo语言支持两种注释方式：行注释和块注释。\n\n## 行注释\n\n行注释以 `//` 开头...",
    "examples": [
      {
        "title": "行注释示例",
        "code": "// 这是一个行注释\nvar x int // 变量声明后的注释",
        "language": "go",
        "explanation": "行注释以//开头，延续到行尾"
      },
      {
        "title": "块注释示例",
        "code": "/*\n这是一个\n多行块注释\n*/",
        "language": "go",
        "explanation": "块注释可以跨越多行"
      }
    ],
    "related_chapters": ["tokens", "semicolons"]
  },
  "timestamp": 1701676800
}
```

**HTML格式** (`?format=html`):

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>注释 (Comments) - Go学习</title>
    <style>
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 { color: #0066cc; border-bottom: 2px solid #0066cc; padding-bottom: 10px; }
        h2 { color: #333; margin-top: 30px; }
        pre { 
            background: #f4f4f4;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
            border-left: 4px solid #0066cc;
        }
        code { 
            font-family: 'Courier New', monospace;
            background: #f4f4f4;
            padding: 2px 6px;
            border-radius: 3px;
        }
        .example {
            margin: 20px 0;
            padding: 15px;
            background: #f9f9f9;
            border-radius: 5px;
        }
        .example-title {
            font-weight: bold;
            color: #0066cc;
            margin-bottom: 10px;
        }
        .related {
            margin-top: 30px;
            padding: 15px;
            background: #e8f4f8;
            border-radius: 5px;
        }
        .related a {
            color: #0066cc;
            text-decoration: none;
            margin-right: 15px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>注释 (Comments)</h1>
        
        <h2>概述</h2>
        <p>Go语言支持两种注释方式：行注释和块注释。</p>
        
        <h2>行注释</h2>
        <p>行注释以 <code>//</code> 开头，延续到行尾。</p>
        
        <div class="example">
            <div class="example-title">行注释示例</div>
            <pre><code class="language-go">// 这是一个行注释
var x int // 变量声明后的注释</code></pre>
            <p>行注释以//开头，延续到行尾</p>
        </div>
        
        <div class="example">
            <div class="example-title">块注释示例</div>
            <pre><code class="language-go">/*
这是一个
多行块注释
*/</code></pre>
            <p>块注释可以跨越多行</p>
        </div>
        
        <div class="related">
            <strong>相关章节：</strong>
            <a href="/api/v1/topic/lexical_elements/tokens">标记 (Tokens)</a>
            <a href="/api/v1/topic/lexical_elements/semicolons">分号 (Semicolons)</a>
        </div>
    </div>
</body>
</html>
```

#### 错误响应

**404 Not Found** - 章节不存在

```json
{
  "code": 404,
  "message": "章节不存在",
  "error": "chapter 'nonexistent' not found in topic 'lexical_elements'",
  "timestamp": 1701676800
}
```

**HTML格式错误响应**:

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>错误 - 404</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 50px; text-align: center; }
        .error-code { font-size: 72px; color: #cc0000; }
        .error-message { font-size: 24px; color: #333; margin: 20px 0; }
        .error-detail { color: #666; }
    </style>
</head>
<body>
    <div class="error-code">404</div>
    <div class="error-message">章节不存在</div>
    <div class="error-detail">chapter 'nonexistent' not found in topic 'lexical_elements'</div>
    <p><a href="/api/v1/topics">返回主题列表</a></p>
</body>
</html>
```

---

## 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|-------|-----------|------|
| 0 | 200 | 成功 |
| 400 | 400 | 请求参数错误（如format参数无效） |
| 404 | 404 | 资源不存在（主题或章节） |
| 500 | 500 | 服务器内部错误 |

---

## 请求示例

### 使用curl

#### 获取主题列表（JSON格式）

```bash
curl -X POST "http://localhost:8080/api/v1/topics?format=json" \
  -H "Content-Type: application/json"
```

#### 获取章节内容（HTML格式）

```bash
curl -X POST "http://localhost:8080/api/v1/topic/lexical_elements/comments?format=html" \
  -H "Content-Type: application/json"
```

### 使用JavaScript (fetch)

```javascript
// 获取主题列表
fetch('http://localhost:8080/api/v1/topics?format=json', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(response => response.json())
.then(data => console.log(data));

// 获取章节内容
fetch('http://localhost:8080/api/v1/topic/lexical_elements/comments?format=json', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(response => response.json())
.then(data => console.log(data));
```

### 使用Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

func main() {
    // 获取主题列表
    resp, err := http.Post(
        "http://localhost:8080/api/v1/topics?format=json",
        "application/json",
        bytes.NewBuffer([]byte("{}")),
    )
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    var result map[string]interface{}
    json.Unmarshal(body, &result)
    
    fmt.Printf("响应: %+v\n", result)
}
```

---

## 版本历史

| 版本 | 日期 | 变更说明 |
|------|------|---------|
| 1.0.0 | 2025-12-04 | 初始版本，定义3个核心接口 |

---

## 附录：完整接口清单

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 获取主题列表 | POST | `/api/v1/topics` | 返回所有学习主题 |
| 获取Lexical Elements菜单 | POST | `/api/v1/topic/lexical_elements` | 返回词法元素章节列表 |
| 获取Comments章节 | POST | `/api/v1/topic/lexical_elements/comments` | 返回注释章节内容 |
| 获取Tokens章节 | POST | `/api/v1/topic/lexical_elements/tokens` | 返回标记章节内容 |
| 获取其他章节 | POST | `/api/v1/topic/lexical_elements/{chapter}` | 返回指定章节内容 |

所有接口均支持`?format=json`或`?format=html`参数。
