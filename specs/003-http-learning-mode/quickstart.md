# Quick Start: HTTP学习模式

**Feature**: 003-http-learning-mode  
**Date**: 2025-12-04  
**Audience**: 开发者和用户

---

## 目录

1. [功能概述](#功能概述)
2. [前置条件](#前置条件)
3. [配置文件设置](#配置文件设置)
4. [启动方式](#启动方式)
5. [使用示例](#使用示例)
6. [常见问题](#常见问题)
7. [开发者指南](#开发者指南)

---

## 功能概述

HTTP学习模式为Go学习工具提供了基于Web的访问方式。用户可以通过两种模式使用本工具：

- **命令行模式**（默认）：传统的终端交互式学习
- **HTTP模式**：通过浏览器或API客户端访问学习内容

两种模式访问相同的学习内容，确保学习体验一致。

### 核心特性

✅ 支持命令行和HTTP两种运行模式  
✅ 所有接口使用POST方法  
✅ 支持JSON和HTML两种响应格式  
✅ YAML配置文件管理  
✅ 结构化日志记录  
✅ 优雅关闭机制  
✅ 并发请求支持

---

## 前置条件

### 系统要求

- **操作系统**: Windows、Linux或macOS
- **Go版本**: 1.24.5或更高
- **内存**: 至少100MB可用内存
- **磁盘空间**: 至少50MB可用空间

### 依赖安装

确保已安装Go环境：

```bash
# 检查Go版本
go version
# 输出应显示: go version go1.24.5 或更高版本
```

### 获取项目

```bash
# 克隆项目（如果尚未克隆）
git clone <repository-url>
cd go-study2

# 下载依赖
go mod download
```

---

## 配置文件设置

### 创建配置文件

在项目根目录创建`config.yaml`文件：

```bash
# Windows PowerShell
New-Item -Path "config.yaml" -ItemType File

# Linux/macOS
touch config.yaml
```

### 配置文件内容

编辑`config.yaml`，添加以下内容：

```yaml
# HTTP服务器配置
server:
  # 监听地址（必填）
  # 使用 127.0.0.1 仅允许本地访问
  # 使用 0.0.0.0 允许外部访问
  host: "127.0.0.1"
  
  # 监听端口（必填，范围1-65535）
  port: 8080
  
  # 优雅关闭超时时间（秒，可选，默认10）
  shutdownTimeout: 10

# 日志配置
logger:
  # 日志级别: DEBUG, INFO, WARN, ERROR
  level: "INFO"
  
  # 日志文件路径（可选）
  path: "./logs"
  
  # 是否输出到控制台（可选，默认true）
  stdout: true
```

### 配置说明

| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `server.host` | string | ✅ | 无 | HTTP服务监听地址 |
| `server.port` | int | ✅ | 无 | HTTP服务监听端口 |
| `server.shutdownTimeout` | int | ❌ | 10 | 优雅关闭超时时间（秒） |
| `logger.level` | string | ✅ | INFO | 日志级别 |
| `logger.path` | string | ❌ | ./logs | 日志文件路径 |
| `logger.stdout` | bool | ❌ | true | 是否输出到控制台 |

---

## 启动方式

### 方式1: 命令行模式（默认）

不带任何参数启动，进入传统的命令行交互模式：

```bash
# Windows
.\go-study2.exe

# Linux/macOS
./go-study2
```

**预期输出**:

```
=== Go语言学习工具 ===

请选择学习主题:
0. 退出
1. 词法元素 (Lexical Elements)

请输入选项 (0-1):
```

### 方式2: HTTP服务模式

使用`-d`或`--daemon`参数启动HTTP服务：

```bash
# 使用 -d 参数
.\go-study2.exe -d

# 或使用 --daemon 参数（效果相同）
.\go-study2.exe --daemon
```

**预期输出**:

```
2025-12-04 18:00:00 [INFO] 正在加载配置文件...
2025-12-04 18:00:00 [INFO] 配置加载成功
2025-12-04 18:00:00 [INFO] HTTP服务已启动
2025-12-04 18:00:00 [INFO] 监听地址: http://127.0.0.1:8080
2025-12-04 18:00:00 [INFO] 按 Ctrl+C 停止服务
```

### 停止HTTP服务

在HTTP模式下，按`Ctrl+C`优雅关闭服务：

```
^C
2025-12-04 18:30:00 [INFO] 正在关闭HTTP服务...
2025-12-04 18:30:01 [INFO] HTTP服务已安全关闭
```

---

## 使用示例

### 示例1: 浏览器访问

#### 步骤1: 启动HTTP服务

```bash
.\go-study2.exe -d
```

#### 步骤2: 打开浏览器

访问以下URL：

- **主题列表**: http://localhost:8080/api/v1/topics?format=html
- **Lexical Elements菜单**: http://localhost:8080/api/v1/topic/lexical_elements?format=html
- **注释章节**: http://localhost:8080/api/v1/topic/lexical_elements/comments?format=html

### 示例2: 使用curl调用API

#### 获取主题列表（JSON格式）

```bash
curl -X POST "http://localhost:8080/api/v1/topics?format=json" \
  -H "Content-Type: application/json"
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "topics": [
      {
        "id": "lexical_elements",
        "title": "词法元素 (Lexical Elements)",
        "description": "Go语言的基本词法元素",
        "chapters": [
          {
            "id": "comments",
            "title": "注释 (Comments)",
            "path": "/api/v1/topic/lexical_elements/comments"
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

#### 获取章节内容

```bash
curl -X POST "http://localhost:8080/api/v1/topic/lexical_elements/comments?format=json" \
  -H "Content-Type: application/json"
```

### 示例3: 使用Postman

1. **创建新请求**
   - Method: `POST`
   - URL: `http://localhost:8080/api/v1/topics`
   - Params: `format=json`

2. **设置Headers**
   - `Content-Type`: `application/json`

3. **发送请求**
   - 点击"Send"按钮
   - 查看响应数据

### 示例4: 使用Go代码调用

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
    // 创建POST请求
    url := "http://localhost:8080/api/v1/topic/lexical_elements/comments?format=json"
    resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("{}")))
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("读取响应失败: %v\n", err)
        return
    }
    
    // 解析JSON
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        fmt.Printf("解析JSON失败: %v\n", err)
        return
    }
    
    // 打印结果
    fmt.Printf("响应: %+v\n", result)
}
```

---

## 常见问题

### Q1: 启动时提示"配置项 server.host 为必填项"

**原因**: 配置文件缺失或未正确设置`server.host`。

**解决方案**:

1. 确认`config.yaml`文件存在于项目根目录
2. 检查配置文件中是否包含`server.host`配置
3. 确保配置格式正确（YAML语法）

```yaml
server:
  host: "127.0.0.1"  # 确保此行存在
  port: 8080
```

### Q2: 启动时提示"端口被占用"

**原因**: 指定的端口已被其他程序使用。

**解决方案**:

**方案1**: 修改配置文件中的端口号

```yaml
server:
  port: 9090  # 改为其他未占用的端口
```

**方案2**: 查找并关闭占用端口的程序

```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <进程ID> /F

# Linux/macOS
lsof -i :8080
kill -9 <进程ID>
```

### Q3: 请求返回404错误

**原因**: 请求的章节不存在。

**解决方案**:

1. 检查章节ID是否正确（区分大小写）
2. 确认章节文件存在于`resource/lexical_elements/`目录
3. 先调用`/api/v1/topics`接口查看可用章节列表

### Q4: 响应格式不是预期的HTML/JSON

**原因**: 未正确设置`format`查询参数。

**解决方案**:

确保URL包含正确的格式参数：

- JSON格式: `?format=json`
- HTML格式: `?format=html`

示例: `http://localhost:8080/api/v1/topics?format=html`

### Q5: 如何允许外部访问HTTP服务？

**原因**: 默认配置仅允许本地访问（`127.0.0.1`）。

**解决方案**:

修改配置文件中的`host`为`0.0.0.0`：

```yaml
server:
  host: "0.0.0.0"  # 允许所有网络接口访问
  port: 8080
```

**安全提示**: 仅在受信任的网络环境中使用`0.0.0.0`。

---

## 开发者指南

### 编译项目

```bash
# 编译当前平台可执行文件
go build -o go-study2

# 编译Windows可执行文件
GOOS=windows GOARCH=amd64 go build -o go-study2.exe

# 编译Linux可执行文件
GOOS=linux GOARCH=amd64 go build -o go-study2
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/content

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 添加新章节

#### 步骤1: 创建内容文件

在`resource/lexical_elements/`目录下创建新的Markdown文件：

```bash
# 例如: 添加"运算符"章节
New-Item -Path "resource/lexical_elements/operators.md" -ItemType File
```

#### 步骤2: 编写章节内容

编辑`operators.md`文件：

```markdown
# 运算符 (Operators)

## 概述

Go语言支持多种运算符...

## 算术运算符

...
```

#### 步骤3: 更新主题配置

在内容提供者中注册新章节（具体实现见开发文档）。

#### 步骤4: 测试

```bash
# 启动HTTP服务
.\go-study2.exe -d

# 访问新章节
curl -X POST "http://localhost:8080/api/v1/topic/lexical_elements/operators?format=json"
```

### 日志查看

日志文件位于`logs/`目录（可通过配置文件修改）：

```bash
# 查看最新日志
tail -f logs/app.log

# Windows PowerShell
Get-Content logs/app.log -Tail 50 -Wait
```

### 性能测试

使用Apache Bench进行并发测试：

```bash
# 安装ab工具（如果未安装）
# Ubuntu: sudo apt-get install apache2-utils
# macOS: 已预装

# 执行并发测试（100个请求，并发50）
ab -n 100 -c 50 -p empty.json -T application/json \
  "http://localhost:8080/api/v1/topics?format=json"
```

创建`empty.json`文件（空JSON对象）：

```json
{}
```

---

## 下一步

- 📖 阅读[API契约文档](./contracts/api-spec.md)了解详细接口定义
- 🏗️ 查看[数据模型文档](./data-model.md)了解系统实体设计
- 🔬 参考[研究文档](./research.md)了解技术决策
- 📋 查看[实现计划](./plan.md)了解整体架构

---

## 支持

如有问题，请：

1. 查看本文档的[常见问题](#常见问题)部分
2. 检查日志文件获取详细错误信息
3. 提交Issue到项目仓库

---

**版本**: 1.0.0  
**最后更新**: 2025-12-04
