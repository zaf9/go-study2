# Data Model: HTTP学习模式

**Feature**: 003-http-learning-mode  
**Date**: 2025-12-04  
**Purpose**: 定义系统中的核心实体、数据结构和关系

---

## 实体概览

本功能涉及以下核心实体：

1. **Topic (主题)** - 顶层学习主题（如"Lexical Elements"）
2. **Chapter (章节)** - 具体学习章节（如"Comments"、"Tokens"）
3. **Content (内容)** - 章节的学习内容数据
4. **ServerConfig (服务配置)** - HTTP服务器运行配置
5. **Response (响应)** - HTTP接口统一响应结构

---

## 1. Topic (主题)

### 描述
代表顶层学习主题，如"Lexical Elements"、"Types"等。主题包含多个子章节。

### 数据结构

```go
// Topic 学习主题
type Topic struct {
    // ID 主题唯一标识符（如"lexical_elements"）
    ID string `json:"id"`
    
    // Title 主题标题（中文）
    Title string `json:"title"`
    
    // Description 主题描述
    Description string `json:"description"`
    
    // Chapters 子章节列表
    Chapters []ChapterInfo `json:"chapters"`
    
    // Order 显示顺序
    Order int `json:"order"`
}

// ChapterInfo 章节信息（简化版，用于列表展示）
type ChapterInfo struct {
    // ID 章节标识符（如"comments"）
    ID string `json:"id"`
    
    // Title 章节标题
    Title string `json:"title"`
    
    // Path API路径（如"/api/v1/topic/lexical_elements/comments"）
    Path string `json:"path"`
}
```

### 示例数据

```json
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
```

### 验证规则

- `ID`: 必填，只能包含小写字母、数字和下划线
- `Title`: 必填，长度1-100字符
- `Chapters`: 可为空数组，但不能为null
- `Order`: 必须≥0

---

## 2. Chapter (章节)

### 描述
代表具体的学习章节，包含完整的学习内容、代码示例等。

### 数据结构

```go
// Chapter 学习章节
type Chapter struct {
    // ID 章节唯一标识符
    ID string `json:"id"`
    
    // Title 章节标题
    Title string `json:"title"`
    
    // TopicID 所属主题ID
    TopicID string `json:"topic_id"`
    
    // Content 章节内容（Markdown格式）
    Content string `json:"content"`
    
    // Examples 代码示例列表
    Examples []CodeExample `json:"examples,omitempty"`
    
    // RelatedChapters 相关章节ID列表
    RelatedChapters []string `json:"related_chapters,omitempty"`
}

// CodeExample 代码示例
type CodeExample struct {
    // Title 示例标题
    Title string `json:"title"`
    
    // Code 代码内容
    Code string `json:"code"`
    
    // Language 编程语言（默认"go"）
    Language string `json:"language"`
    
    // Explanation 代码说明
    Explanation string `json:"explanation,omitempty"`
}
```

### 示例数据

```json
{
  "id": "comments",
  "title": "注释 (Comments)",
  "topic_id": "lexical_elements",
  "content": "# 注释\n\nGo语言支持两种注释方式...",
  "examples": [
    {
      "title": "行注释示例",
      "code": "// 这是一个行注释\nvar x int // 变量声明后的注释",
      "language": "go",
      "explanation": "行注释以//开头，延续到行尾"
    }
  ],
  "related_chapters": ["tokens", "semicolons"]
}
```

### 验证规则

- `ID`: 必填，格式同Topic.ID
- `Title`: 必填
- `TopicID`: 必填，必须是有效的主题ID
- `Content`: 必填，长度≥10字符
- `Examples`: 可选，但如果提供则每个示例的Code必填

---

## 3. Content (内容)

### 描述
章节内容由`lexical_elements`包中的`Get*Content()`函数生成，无需从文件加载。

### 数据结构

```go
// ContentGenerator 内容生成器接口
type ContentGenerator interface {
    // GetContent 获取章节内容
    GetContent() string
}

// ChapterContentMap 章节ID到内容生成函数的映射
var ChapterContentMap = map[string]func() string{
    "comments":     lexical_elements.GetCommentsContent,
    "tokens":       lexical_elements.GetTokensContent,
    "semicolons":   lexical_elements.GetSemicolonsContent,
    "identifiers":  lexical_elements.GetIdentifiersContent,
    "keywords":     lexical_elements.GetKeywordsContent,
    "operators":    lexical_elements.GetOperatorsContent,
    "integers":     lexical_elements.GetIntegersContent,
    "floats":       lexical_elements.GetFloatsContent,
    "imaginary":    lexical_elements.GetImaginaryContent,
    "runes":        lexical_elements.GetRunesContent,
    "strings":      lexical_elements.GetStringsContent,
}
```

### 内容生成示例

```go
// 获取注释章节内容
content := lexical_elements.GetCommentsContent()

// 内容示例（部分）
/*
--- Go 语言的注释 ---
注释是代码中非常重要的一部分，用于解释代码的功能、目的和实现方式。
Go 语言支持两种类型的注释：

1. 单行注释 (Single-line Comments):
// 这是一个单行注释。
变量 'variable' 的值是: 10 (这行代码后面就有一个单行注释)。

2. 多行注释 (Multi-line or Block Comments):
/*
 * 这是一个多行注释。
 * 它可以包含很多行的文本。
 */
Go 源码中经常使用多行注释来为包、函数、类型或变量提供文档。
这种文档注释（doc comments）是一种重要的实践，可以使用 go doc 工具来查看。
*/
```

### 章节列表

| 章节ID | 函数名 | 标题 |
|--------|--------|------|
| comments | GetCommentsContent() | 注释 (Comments) |
| tokens | GetTokensContent() | 标记 (Tokens) |
| semicolons | GetSemicolonsContent() | 分号 (Semicolons) |
| identifiers | GetIdentifiersContent() | 标识符 (Identifiers) |
| keywords | GetKeywordsContent() | 关键字 (Keywords) |
| operators | GetOperatorsContent() | 运算符 (Operators) |
| integers | GetIntegersContent() | 整数 (Integers) |
| floats | GetFloatsContent() | 浮点数 (Floats) |
| imaginary | GetImaginaryContent() | 虚数 (Imaginary) |
| runes | GetRunesContent() | 符文 (Runes) |
| strings | GetStringsContent() | 字符串 (Strings) |

---

## 4. ServerConfig (服务配置)

### 描述
HTTP服务器的运行时配置，从config.yaml加载。

### 数据结构

```go
// ServerConfig HTTP服务器配置
type ServerConfig struct {
    // Host 监听地址（必填）
    Host string `json:"host" v:"required|ip"`
    
    // Port 监听端口（必填，范围1-65535）
    Port int `json:"port" v:"required|between:1,65535"`
    
    // ShutdownTimeout 优雅关闭超时时间（秒）
    ShutdownTimeout int `json:"shutdownTimeout" v:"min:1"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
    // Level 日志级别（DEBUG/INFO/WARN/ERROR）
    Level string `json:"level" v:"required|in:DEBUG,INFO,WARN,ERROR"`
    
    // Path 日志文件路径
    Path string `json:"path"`
    
    // Stdout 是否输出到控制台
    Stdout bool `json:"stdout"`
}

// AppConfig 应用配置（根配置）
type AppConfig struct {
    Server ServerConfig `json:"server"`
    Logger LoggerConfig `json:"logger"`
}
```

### 配置文件示例 (config.yaml)

```yaml
server:
  host: "127.0.0.1"
  port: 8080
  shutdownTimeout: 10

logger:
  level: "INFO"
  path: "./logs"
  stdout: true
```

### 验证规则

- `Server.Host`: 必填，必须是有效IP地址或"0.0.0.0"
- `Server.Port`: 必填，范围1-65535
- `Server.ShutdownTimeout`: 可选，默认10秒，最小1秒
- `Logger.Level`: 必填，只能是DEBUG/INFO/WARN/ERROR之一
- `Logger.Path`: 可选，默认"./logs"
- `Logger.Stdout`: 可选，默认true

---

## 5. Response (响应)

### 描述
HTTP接口的统一响应结构，支持JSON和HTML两种格式。

### 数据结构

```go
// APIResponse 统一API响应结构
type APIResponse struct {
    // Code 响应码（0表示成功，非0表示错误）
    Code int `json:"code"`
    
    // Message 响应消息
    Message string `json:"message"`
    
    // Data 响应数据（具体类型根据接口而定）
    Data interface{} `json:"data,omitempty"`
    
    // Timestamp 响应时间戳
    Timestamp int64 `json:"timestamp"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
    Code      int    `json:"code"`
    Message   string `json:"message"`
    Error     string `json:"error"`
    Timestamp int64  `json:"timestamp"`
}
```

### JSON响应示例

#### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "lexical_elements",
    "title": "词法元素 (Lexical Elements)",
    "chapters": [...]
  },
  "timestamp": 1701676800
}
```

#### 错误响应
```json
{
  "code": 404,
  "message": "章节不存在",
  "error": "chapter 'nonexistent' not found",
  "timestamp": 1701676800
}
```

### HTML响应模板

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .content { line-height: 1.6; }
        pre { background: #f4f4f4; padding: 10px; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <div class="content">
        {{.Content}}
    </div>
</body>
</html>
```

---

## 实体关系图

```
┌─────────────┐
│   Topic     │
│             │
│ - ID        │
│ - Title     │
│ - Chapters[]│───┐
└─────────────┘   │
                  │ 1:N
                  ▼
            ┌─────────────┐
            │  Chapter    │
            │             │
            │ - ID        │
            │ - TopicID   │
            │ - Content   │
            └─────────────┘
                  │
                  │ 生成自
                  ▼
            ┌──────────────────────┐
            │ lexical_elements包   │
            │                      │
            │ Get*Content()函数    │
            └──────────────────────┘
```

---

## 数据流

### 命令行模式数据流

```
用户输入 → 菜单系统 → Display*()函数 
         → Get*Content()函数 → 生成内容字符串 
         → fmt.Print()打印 → 终端显示
```

### HTTP模式数据流

```
HTTP POST请求 → 路由匹配 → Handler处理器 
              → 根据章节ID查找对应的Get*Content()函数
              → 调用函数生成内容字符串
              → 构建Response对象 
              → FormatMiddleware转换格式 → 返回JSON/HTML
```

---

## 状态管理

### 无状态设计原则

- **内容生成函数**: 无状态，每次调用重新生成内容（内容是代码逻辑，无需缓存）
- **HTTP处理器**: 无状态，不保存请求间的数据
- **配置对象**: 启动时加载一次，运行期间只读

### 内容一致性保证

```go
// 命令行模式
func DisplayComments() {
    fmt.Print(GetCommentsContent())  // 调用内容生成函数
}

// HTTP模式
func ChapterHandler(r *ghttp.Request) {
    content := GetCommentsContent()  // 调用相同的内容生成函数
    // ... 包装成JSON/HTML响应
}
```

通过让两种模式调用相同的`Get*Content()`函数，确保内容100%一致。

---

## 数据验证

### 输入验证

所有外部输入必须经过验证：

1. **章节ID验证**: 只允许字母、数字、下划线、斜杠
2. **格式参数验证**: 只允许"json"或"html"
3. **配置验证**: 使用GoFrame的validation标签

### 验证函数示例

```go
// ValidateChapterID 验证章节ID
func ValidateChapterID(id string) error {
    if id == "" {
        return fmt.Errorf("章节ID不能为空")
    }
    
    // 只允许字母、数字、下划线、斜杠
    matched, _ := regexp.MatchString(`^[a-z0-9_/]+$`, id)
    if !matched {
        return fmt.Errorf("章节ID格式无效: %s", id)
    }
    
    return nil
}

// ValidateFormat 验证响应格式
func ValidateFormat(format string) error {
    if format != "json" && format != "html" {
        return fmt.Errorf("不支持的格式: %s (仅支持json或html)", format)
    }
    return nil
}
```

---

## 总结

本数据模型定义了4个核心实体：

1. ✅ **Topic** - 主题结构，包含章节列表
2. ✅ **Chapter** - 章节详细信息，包含内容
3. ✅ **Content** - 通过`lexical_elements`包的`Get*Content()`函数生成
4. ✅ **ServerConfig** - 服务器配置
5. ✅ **Response** - 统一响应结构

所有实体均遵循：
- 明确的字段定义和类型
- 清晰的验证规则
- 中文注释和文档
- 无状态设计原则

**关键设计决策**：
- 内容由现有的`lexical_elements`包函数生成，无需文件存储
- 命令行和HTTP模式调用相同的`Get*Content()`函数，确保内容一致性
- 重构`Display*()`函数为内容生成+打印两步，保持向后兼容

**下一步**: 生成API contracts和quickstart文档。
