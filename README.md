# Go-Study2 - Go语言词法元素学习工具

> 一个支持**命令行**和**HTTP服务**双模式的Go语言学习工具，帮助学习者系统掌握词法元素知识，提供交互式菜单和Web API两种访问方式。

[![Go Version](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org)
[![GoFrame](https://img.shields.io/badge/GoFrame-v2.9.5-green.svg)](https://goframe.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## 📝 目录 Table of Contents

- [背景与目标](#-背景与目标-background--motivation)
- [功能特性](#-功能特性-features)
- [技术栈](#-技术栈-tech-stack)
- [快速开始](#-快速开始-quick-start)
- [安装](#-安装-installation)
- [使用方法](#-使用方法-usage)
- [示例](#-示例-examples)
- [项目结构](#-项目结构-project-structure)
- [前端 UI](#-前端-ui)
- [配置](#️-配置-configuration)
- [开发与测试](#-开发与测试-development--testing)
- [Roadmap](#-roadmap)
- [贡献指南](#-贡献指南-contributing)
- [许可证](#-许可证-license)
- [致谢](#-致谢-acknowledgements)

---

## 🎯 背景与目标 Background & Motivation

在学习Go语言的过程中，词法元素（Lexical Elements）是最基础但也是最重要的知识点。然而，现有的学习资源往往缺乏系统性和交互性，学习者需要在大量文档中来回切换，效率低下。

**本项目旨在解决以下痛点：**

- 📚 **知识碎片化**：词法元素知识分散在各处，缺乏系统整理
- 🔍 **缺乏实践**：理论知识多，可运行的代码示例少
- 🌐 **语言障碍**：优质资源多为英文，中文学习者需要额外的理解成本
- 🎯 **学习路径不清晰**：不知道从哪里开始，如何循序渐进

**目标用户：**

- Go语言初学者
- 希望系统复习词法基础的开发者
- 需要中文学习资源的学习者

**项目定位：**

这是一个命令行交互式学习工具，提供结构化的知识体系和可运行的代码示例，帮助学习者快速掌握Go语言词法元素。

---

## ✨ 功能特性 Features

### 核心功能

- 🎯 **双模式运行** - 支持命令行交互模式和HTTP服务模式
- 🖥️ **现代Web界面** 🆕 - Next.js + Ant Design 响应式 UI，桌面/移动端适配
- 🔐 **用户认证** 🆕 - 注册/登录/记住我，JWT 访问令牌 + HttpOnly 刷新令牌，过期自动刷新
- 📍 **学习进度** 🆕 - 记录章节状态与滚动位置，支持“继续上次学习”
- 📝 **主题测验** 🆕 - 单选/多选测验与历史记录，提交即时评分
- 📖 **全面覆盖** - 涵盖Go语言规范中的词法元素和常量系统
- 💻 **可运行示例** - 每个知识点都配有可直接运行的代码示例
- 🇨🇳 **中文注释** - 所有代码注释和说明均为中文，降低学习门槛
- 📚 **多模块支持** - 词法元素模块 + 常量学习模块（12个子主题）

### 命令行模式特性

- 🎯 **菜单驱动界面** - 清晰的层级菜单，轻松导航各个知识点
- 🚀 **零依赖运行** - 编译后的可执行文件无需额外依赖
- ⌨️ **交互式学习** - 即时反馈，边学边练

### HTTP服务模式特性 🆕

- 🌐 **RESTful API** - 标准化的HTTP接口，支持JSON和HTML两种响应格式
- 🪪 **鉴权保护** - 受保护路由统一JWT校验，自动重定向登录
- 🔁 **刷新机制** - 7天刷新令牌，可配置“记住我”延长会话
- 🔌 **灵活访问** - 通过浏览器、curl、Postman或任何HTTP客户端访问
- ⚙️ **YAML配置** - 灵活的配置文件管理服务器参数
- 📊 **结构化日志** - 详细的请求日志和错误追踪
- 🛡️ **优雅关闭** - 支持信号处理和优雅停机
- 🚀 **并发支持** - 可处理多个并发请求

### 前端 UI 特性 🆕

- 📱 **响应式布局** - Mobile <768px / Tablet 768-1024px / Desktop >1024px
- 🧭 **学习导航** - 主题列表、章节锚点、代码高亮与分段呈现
- 🔖 **进度续学** - 展示百分比、最近访问时间、滚动位置恢复
- 🧪 **测验体验** - 题目来源说明、防重复提交、历史筛选
- ⚙️ **一体化部署** - 静态导出到 `frontend/out`，后端同端口托管

### 质量保证

- 🧪 **高测试覆盖率** - 80%以上的单元测试覆盖率，保证代码质量
- ✅ **内容一致性** - CLI和HTTP模式返回完全相同的学习内容
- 🔌 **易于扩展** - 模块化设计，可轻松添加新的学习主题

---

## 🧱 技术栈 Tech Stack

- **语言**: Go 1.24.5；TypeScript 5 + React 18
- **后端**: GoFrame v2.9.5、SQLite3（WAL）、JWT（golang-jwt v5）、bcrypt、GoFrame ORM
- **前端**: Next.js 14（App Router，`output: 'export'`）、Ant Design 5、SWR、Axios、Prism.js（按需语言包）、Tailwind CSS
- **构建工具**: Go Modules、npm；前端静态导出目录 `frontend/out`
- **测试**: Go 标准测试 + 覆盖率工具；前端 Jest + React Testing Library，核心组件/Hook 覆盖率≥80%
- **开发环境**: 支持 Windows/Linux/macOS，前后端同端口一体化部署

---

## 🚀 快速开始 Quick Start

> 提示：如仓库根存在 `./build.bat`，优先执行以完成依赖检查与编译，再按下列方式启动。

### 方式一：命令行模式（传统方式）

**30秒快速体验：**

```bash
# 克隆仓库
git clone https://github.com/yourusername/go-study2.git

# 进入项目目录并切换到后端
cd go-study2/backend

# 运行程序（主菜单含 Lexical / Constants / Variables / Types）
go run main.go
```

**预期输出：**

```
Go Lexical Elements Learning Tool
---------------------------------
Please select a topic to study:
0. Lexical elements
1. Constants
2. Variables
3. Types
q. Quit

Enter your choice: 
```

输入 `0/1/2/3` 进入对应章节学习。Types 子菜单支持：编号查看内容与测验；`o` 打印提纲；`quiz` 综合测验；`search <keyword>` 关键词检索；`q` 返回。

### 方式二：HTTP服务模式 🆕

**60秒启动Web服务：**

```bash
# 1. 克隆仓库（如果尚未克隆）
git clone https://github.com/yourusername/go-study2.git
cd go-study2/backend

# 2. 确保配置文件存在（默认端口 8080）
#   - backend/configs/config.yaml 已预置 server/logger/jwt/database/static 段

# 3. 启动HTTP服务（生产推荐先在根执行 ./build.bat）
go run main.go -d
```

**浏览器访问：**

- 主题列表（HTML）：http://localhost:8080/api/v1/topics?format=html  
- 词法元素章节：http://localhost:8080/api/v1/topic/lexical_elements/comments?format=html  
- Constants 菜单：http://localhost:8080/api/v1/topic/constants?format=html  
- Types 提纲：http://localhost:8080/api/v1/topic/types/outline?format=html  
- 受保护路由示例：`/api/v1/progress`（需登录并携带 Authorization 头）

**API 调用（JSON 示例）：**

```bash
curl http://localhost:8080/api/v1/topics
curl http://localhost:8080/api/v1/topic/constants/boolean
curl http://localhost:8080/api/v1/topic/types/search?keyword=map%20key
```

### 方式三：前端 UI 模式（Web） 🆕

**开发调试：**

```bash
# 后端启动（默认 8080）
cd backend
go run main.go -d   # 若有 ./build.bat 请先在根执行

# 前端启动（默认 3000，已代理到 http://localhost:8080/api）
cd ../frontend
npm install
npm run dev
```

**主要页面与路由：**

- 学习进度页：`/progress`（整体进度、继续学习）
- 主题/章节详情：`/topics/[topic]/[chapter]`（阅读恢复与测验入口）
- 测验：`/quiz/[topic]/[chapter]`，历史记录：`/quiz`

**前后端联调（开发模式）：**

```bash
# 1) 若仓库根存在 build.bat，优先执行（依赖检查与编译）
./build.bat

# 2) 启动后端（默认 8080，提供 /api/v1）
cd backend
go run main.go -d

# 3) 启动前端（显式指向后端 API，避免 404）
cd ../frontend
$env:NEXT_PUBLIC_API_URL="http://localhost:8080/api/v1"
npm install
npm run dev   # http://localhost:3000/
```

说明：
- 前端开发服务器跑在 3000 端口。设置 `NEXT_PUBLIC_API_URL` 后，`/auth/login` 等请求会直接访问 `http://localhost:8080/api/v1`，不会落到前端导致 404。
- 若需同端口访问生产效果，可 `npm run build && npm run export` 生成 `frontend/out`，然后让后端托管静态文件（访问 `http://localhost:8080/`）。

**生产静态导出与托管：**

```bash
cd frontend
npm install
npm run build && npm run export   # 产物输出到 frontend/out
cd ..

# 后端编译（优先 ./build.bat）
./build.bat || (cd backend && go test ./... && go build -o ./bin/gostudy main.go)

# 启动后端托管 / 与 /api/*
./bin/gostudy -d
```

**访问入口：**

- 开发：`http://localhost:3000/`（前端开发服务器）
- 生产：`http://localhost:8080/`（后端托管静态文件与 API，同端口）

**停止服务：** 按 `Ctrl+C` 优雅关闭

---

### 🧭 日志系统快速验证 (Quickstart Validation)

本项目包含日志子系统的可验流程。要在本地一键验证日志生成与性能，请运行仓库中的快速验证脚本（PowerShell）：

```powershell
cd backend\scripts
powershell -ExecutionPolicy Bypass -File quickstart_validate.ps1
```

脚本会启动服务、执行并发请求压测并将日志收集到 `backend/logs/validation-<timestamp>`。详细手动验证步骤见 `backend/docs/quickstart_validation.md`。

文档索引：

- `backend/docs/quickstart_validation.md` — 快速开始验证步骤
- `backend/docs/logging_best_practices.md` — 日志最佳实践（中文）
- `backend/docs/logging_troubleshooting.md` — 日志故障排查指南（中文）
- `backend/README-logging.md` — 日志系统使用说明与示例

## 📦 安装 Installation

### 方式一：从源码运行（推荐用于学习）

```bash
git clone https://github.com/yourusername/go-study2.git
cd go-study2/backend
go run main.go
```

### 方式二：编译后运行

```bash
git clone https://github.com/yourusername/go-study2.git
cd go-study2/backend
go build -o go-study2
./go-study2  # Linux/macOS
# 或
go-study2.exe  # Windows
```

### 方式三：直接安装（需要发布到GitHub）

```bash
go install github.com/yourusername/go-study2@latest
```

**系统要求：**

- Go 1.24.5 或更高版本
- 支持的操作系统：Windows、Linux、macOS

---

## 🛠 使用方法 Usage

### 基本使用流程

1. **启动程序**：进入 `backend/` 运行 `go run main.go` 或编译后的可执行文件
2. **选择主题**：在主菜单中输入 `0` 选择"词法元素"或 `1` 选择"Constants"
3. **浏览子主题**：在子菜单中选择具体的主题（如注释、布尔常量、iota等）
4. **查看示例**：程序会显示该主题的代码示例和详细解释
5. **返回或退出**：输入 `q` 返回上级菜单或退出程序

### 交互示例

```
Go Lexical Elements Learning Tool
---------------------------------
Please select a topic to study:
0. Lexical elements
1. Constants
2. Variables
q. Quit

Enter your choice: 1

Constants 学习菜单
---------------------------------
请选择要学习的主题:
0. Boolean Constants (布尔常量)
1. Rune Constants (符文常量)
2. Integer Constants (整数常量)
3. Floating-point Constants (浮点常量)
4. Complex Constants (复数常量)
5. String Constants (字符串常量)
6. Constant Expressions (常量表达式)
7. Typed and Untyped Constants (类型化/无类型化常量)
8. Conversions (类型转换)
9. Built-in Functions (内置函数)
10. Iota (iota 特性)
11. Implementation Restrictions (实现限制)
q. 返回上级菜单

请输入您的选择: 0
```

### HTTP服务模式使用 🆕

#### 启动HTTP服务

```bash
# 使用 -d 或 --daemon 参数启动
cd backend && go run main.go -d
# 或
cd backend && go run main.go --daemon
```

#### API端点说明

| 端点 | 方法 | 描述 | 示例URL |
|------|------|------|---------|
| `/api/v1/topics` | GET/POST | 获取所有学习主题列表 | `http://localhost:8080/api/v1/topics` |
| `/api/v1/topic/lexical_elements` | GET/POST | 获取词法元素章节菜单 | `http://localhost:8080/api/v1/topic/lexical_elements` |
| `/api/v1/topic/lexical_elements/{chapter}` | GET/POST | 获取词法元素具体章节内容 | `http://localhost:8080/api/v1/topic/lexical_elements/comments` |
| `/api/v1/topic/constants` | GET/POST | 获取常量学习模块菜单 | `http://localhost:8080/api/v1/topic/constants` |
| `/api/v1/topic/constants/{subtopic}` | GET/POST | 获取常量模块具体子主题内容 | `http://localhost:8080/api/v1/topic/constants/boolean` |

#### 响应格式

通过 `format` 查询参数指定响应格式：

**JSON格式（默认，适合API调用）：**

```bash
curl "http://localhost:8080/api/v1/topics?format=json"
```

**HTML格式（适合浏览器访问）：**

```bash
curl "http://localhost:8080/api/v1/topics?format=html"
# 或在浏览器中直接访问
```

#### 可用章节ID

**词法元素模块 (Lexical Elements)**:
- `comments` - 注释
- `tokens` - 标记
- `semicolons` - 分号
- `identifiers` - 标识符
- `keywords` - 关键字
- `operators` - 运算符
- `integers` - 整数
- `floats` - 浮点数
- `imaginary` - 虚数
- `runes` - 符文
- `strings` - 字符串

**常量学习模块 (Constants)** 🆕:
- `boolean` - 布尔常量
- `rune` - 符文常量
- `integer` - 整数常量
- `floating_point` - 浮点常量
- `complex` - 复数常量
- `string` - 字符串常量
- `expressions` - 常量表达式
- `typed_untyped` - 类型化/无类型化常量
- `conversions` - 类型转换
- `builtin_functions` - 内置函数
- `iota` - iota 特性
- `implementation_restrictions` - 实现限制

---

## 📚 示例 Examples

### 示例1：学习Go语言注释

选择"Comments"主题后，你会看到：

```go
// 这是单行注释
// Go语言支持两种注释方式

/*
这是多行注释
可以跨越多行
常用于文档说明
*/
```

### 示例2：学习标识符规则

选择"Identifiers"主题后，程序会展示：

```go
// 合法的标识符示例
var userName string
var _privateVar int
var 中文变量 string  // Go支持Unicode标识符

// 不合法的标识符（会在注释中说明）
// var 123abc string  // 不能以数字开头
// var for string     // 不能使用关键字
```

### 示例3：理解字符串字面量

```go
// 解释型字符串（双引号）
var s1 = "Hello\nWorld"  // 支持转义字符

// 原始字符串（反引号）
var s2 = `Hello
World`  // 保留原始格式，不转义
```

**更多示例**：每个词法元素子主题都包含完整的代码示例和中文解释。

### 示例4：学习常量表达式

选择"Constant Expressions"主题后，你会看到：

```go
package main

import "fmt"

func main() {
    const (
        a = 10
        b = 20
        sum = a + b        // 30
        diff = b - a       // 10
        prod = a * b       // 200
        quot = b / a       // 2
    )
    
    fmt.Println(sum, diff, prod, quot)
    // 输出: 30 10 200 2
}
```

### 示例5：学习 iota 特性

选择"Iota"主题后，程序会展示：

```go
package main

import "fmt"

func main() {
    const (
        Sunday = iota    // 0
        Monday           // 1
        Tuesday          // 2
        Wednesday        // 3
        Thursday         // 4
        Friday           // 5
        Saturday         // 6
    )
    
    fmt.Println(Sunday, Monday, Saturday)  // 输出: 0 1 6
}
```

**更多示例**：每个学习模块的子主题都包含完整的代码示例和中文解释。

---

## 🗂 项目结构 Project Structure

```
go-study2/
├── backend/                         # 后端主目录
│   ├── go.mod / go.sum              # Go 模块定义与依赖
│   ├── main.go / main_test.go       # 主入口与测试
│   ├── configs/                     # 配置（config.yaml、certs/）
│   ├── data/                        # SQLite 数据文件（自动迁移生成）
│   ├── internal/
│   │   ├── app/                     # 应用层：HTTP 服务器、学习内容
│   │   │   ├── http_server/         # handler、middleware、router、server
│   │   │   ├── lexical_elements/    # 词法元素内容
│   │   │   ├── constants/           # 常量模块内容
│   │   │   └── ...                  # 其他学习主题
│   │   ├── domain/                  # 领域层（user/progress/quiz 实体与服务）
│   │   ├── infrastructure/          # 基础设施层（database、repository 实现）
│   │   ├── pkg/                     # 共享工具（jwt、password）
│   │   └── config/                  # 配置加载与校验
│   ├── tests/                       # unit / integration / contract 测试
│   ├── docs/                        # 后端文档 materials
│   └── scripts/                     # 工具脚本（check-go.ps1）
├── frontend/                        # 前端主目录（Next.js 14）
│   ├── app/                         # 路由：auth、topics、quiz、progress、profile
│   ├── components/                  # UI 组件：auth/layout/learning/quiz/common
│   ├── hooks/                       # 自定义 Hooks（useAuth/useProgress/useQuiz 等）
│   ├── lib/                         # Axios 实例、auth 工具、常量
│   ├── types/                       # TypeScript 类型定义
│   ├── styles/                      # 全局样式与 Tailwind
│   ├── tests/                       # 前端单元与集成测试
│   ├── public/                      # 静态资源
│   └── out/                         # 静态导出产物（构建后生成）
├── specs/                           # 功能规格、计划、任务（含 009-frontend-ui）
├── docs/                            # API、部署等文档
├── .specify/                        # 规范与模板
├── .github/                         # GitHub 配置
└── README.md                        # 本文件（根级说明）
```

**目录说明：**

- `backend/internal/app/http_server/`：API 入口与路由、中间件、认证/进度/测验 handler
- `backend/internal/domain/`：用户、进度、测验的实体、仓储接口与服务
- `backend/internal/infrastructure/`：SQLite 连接、迁移与仓储实现
- `backend/internal/pkg/`：JWT、密码工具等复用模块
- `backend/tests/`：单元、集成、契约测试，覆盖认证/进度/测验/学习内容
- `frontend/app/`：登录注册路由 `(auth)`、受保护路由 `(protected)`（topics/progress/quiz/profile）
- `frontend/components/`：AuthGuard、LoginForm、ChapterContent、QuizItem 等核心组件
- `frontend/hooks/`：`useAuth`、`useProgress`、`useQuiz` 管理跨页面状态
- `frontend/lib/`：Axios 实例与 token 管理，统一错误处理
- `frontend/tests/`：Jest + RTL 测试，覆盖核心组件与 API 层

---

## 🌐 前端 UI

- 位置：`frontend/`（Next.js 14 App Router + TypeScript 5 + Ant Design 5，静态导出）
- 功能：登录/注册/记住我，主题列表、章节阅读、学习进度同步、测验作答与历史记录
- 交互：响应式断点 Mobile/Tablet/Desktop，代码高亮（Prism），章节锚点与进度百分比展示
- 开发：`cd frontend && npm install && npm run dev`（默认 3000，API 代理到 http://localhost:8080/api）
- 构建：`npm run build && npm run export`（预生成 topics/quiz 路由，产物位于 `frontend/out/`）
- 部署：后端 `configs/config.yaml` 的 `static.path` 指向 `../frontend/out`，`server.go` 已启用静态托管与 SPA 回退
- 更多：`frontend/README.md`、`docs/DEPLOYMENT.md`

---

## ⚙️ 配置 Configuration

CLI 模式零配置即可运行；HTTP/前端模式需填写 `backend/configs/config.yaml`：

```yaml
# HTTP 配置：本地开发默认开启，便于直接通过 http://127.0.0.1 访问
http:
  # HTTP 监听端口，生产若启用 HTTPS 可关闭 HTTP 监听
  port: 8080

# HTTPS 配置：启用后建议关闭 HTTP 监听并正确配置证书
https:
  # 是否启用 HTTPS；开启需提供证书与私钥
  enabled: false
  # HTTPS 监听端口，通常使用 443 或 8443
  port: 8443
  # 服务端证书路径，支持相对或绝对路径
  certFile: "./configs/certs/server.crt"
  # 服务端私钥路径，需与证书匹配
  keyFile: "./configs/certs/server.key"
  # 是否跳过客户端证书校验，仅限测试/开发使用
  insecureSkipVerify: false
  # 可选 CA 证书路径，自签名证书时用于建立信任链
  caFile: ""

# 服务基础配置
server:
  # 服务监听地址，生产建议绑定 0.0.0.0 或具体内网地址
  host: "127.0.0.1"
  # 优雅停机等待时间（秒），用于处理中的请求
  shutdownTimeout: 10

# 日志配置
logger:
  # 日志级别，支持 DEBUG/INFO/WARN/ERROR
  level: "INFO"
  # 日志文件存储目录
  path: "./logs"
  # 是否输出到 stdout，容器环境可开启
  stdout: true

# 数据库配置（默认使用 SQLite3）
database:
  # 数据库类型，当前支持 sqlite3
  type: "sqlite3"
  # 数据文件路径，确保目录可写
  path: "./data/gostudy.db"
  # 最大打开连接数，SQLite 一般保持较小值
  maxOpenConns: 10
  # 最大空闲连接数，避免频繁创建连接
  maxIdleConns: 5
  # 连接最大生命周期（秒），0 表示无限制
  connMaxLifetime: 3600
  # SQLite PRAGMA 配置列表，可按需调整
  pragmas:
    # 采用 WAL 模式提升并发读性能
    - "journal_mode=WAL"
    # 设置数据库忙等待时间（毫秒）
    - "busy_timeout=5000"
    # 同步策略 NORMAL 在可靠性与性能间平衡
    - "synchronous=NORMAL"
    # 负值表示以 KiB 为单位的缓存大小
    - "cache_size=-64000"
    # 开启外键约束校验
    - "foreign_keys=ON"

# JWT 配置
jwt:
  # 签名密钥，强烈建议通过环境变量注入
  secret: "${JWT_SECRET}"
  # 访问令牌过期时间（秒）
  accessTokenExpiry: 604800
  # 刷新令牌过期时间（秒）
  refreshTokenExpiry: 604800
  # JWT 发行方标识
  issuer: "go-study2"

# 静态资源配置
static:
  # 是否启用静态资源托管
  enabled: true
  # 静态资源目录，默认指向前端构建产物
  path: "../frontend/out"
  # SPA 路由回退到 index.html
  spaFallback: true
```

- HTTP/HTTPS：启用 HTTPS 时建议关闭 HTTP 监听，需配置 cert/key；自签证书可临时配合 `caFile` 与 `insecureSkipVerify`（仅测试）。
- server：生产可改为 `0.0.0.0`；`shutdownTimeout` 用于优雅停机等待在途请求完成。
- logger：`stdout=true` 适合容器化部署，`path` 为文件输出目录。
- database：SQLite WAL 提升并发读；`busy_timeout` 毫秒，`cache_size` 负值为 KiB，`foreign_keys=ON` 开启外键校验。
- jwt：`secret` 必须通过环境变量注入；访问/刷新令牌时间单位为秒。
- static：指向 `frontend/out` 导出目录，`spaFallback=true` 支持 SPA 前端路由。

---

## 📖 API 文档 API Reference

### HTTP API 端点 🆕

本项目现在提供RESTful API接口：

**基础URL**: `http://localhost:8080/api/v1`

**认证与用户**：

| 端点 | 方法 | 描述 |
|------|------|------|
| `/auth/register` | POST | 用户注册（用户名校验、bcrypt 存储） |
| `/auth/login` | POST | 用户登录，返回访问令牌；支持 `rememberMe` |
| `/auth/refresh` | POST | 使用 HttpOnly 刷新令牌换取新访问令牌 |
| `/auth/profile` | GET | 获取当前用户信息（需 `Authorization: Bearer`） |
| `/auth/change-password` | POST | 已登录用户修改密码（默认管理员首登必须完成） |
| `/auth/logout` | POST | 退出并清除刷新令牌（需 `Authorization: Bearer`） |

**学习/进度/测验**：

| 端点 | 方法 | 描述 |
|------|------|------|
| `/topics` | GET/POST | 获取主题列表 |
| `/topic/lexical_elements` | GET/POST | 获取词法元素菜单 |
| `/topic/lexical_elements/{chapter}` | GET/POST | 获取词法元素章节内容 |
| `/topic/constants` | GET/POST | 获取常量学习模块菜单 |
| `/topic/constants/{subtopic}` | GET/POST | 获取常量模块子主题内容 |
| `/topic/variables` | GET/POST | 获取 Variables 菜单 |
| `/topic/variables/{subtopic}` | GET/POST | 获取 Variables 子主题 |
| `/topic/types` | GET/POST | 获取 Types 菜单 |
| `/topic/types/{subtopic}` | GET/POST | 获取 Types 子主题内容 |
| `/topic/types/outline` | GET/POST | 获取 Types 提纲 |
| `/topic/types/search` | GET/POST | Types 搜索 |
| `/topic/types/quiz/submit` | GET/POST | Types 综合测验提交 |
| `/progress` | GET | 获取当前用户全部学习进度（需登录） |
| `/progress/{topic}` | GET | 获取指定主题进度（需登录） |
| `/progress` | POST | 保存/更新章节进度（需登录） |
| `/quiz/{topic}/{chapter}` | GET | 获取测验题目（需登录） |
| `/quiz/submit` | POST | 提交测验并评分（需登录） |
| `/quiz/history` | GET | 查看历史测验记录，可按主题过滤（需登录） |

**响应格式**：`{code, message, data}`；学习内容接口支持 `?format=json|html`。默认管理员首登时 `login`/`register` 响应将返回 `needPasswordChange=true`，此时除改密/资料/退出外的请求会被 403 并写入审计。

**默认管理员**：初始账号 `admin` / `GoStudy@123`，首次登录会被强制改密，改密后旧口令与旧令牌全部失效。

详细 API 文档：`docs/API.md`、`specs/009-frontend-ui/contracts/openapi.yaml`

### 内部包结构

```go
// 词法元素模块
package lexical_elements
func GetCommentsContent() string
func GetTokensContent() string
// ... 其他内容生成函数

// HTTP服务模块
package http_server
func NewServer(cfg *config.Config, names ...string) (*ghttp.Server, error)
func RegisterRoutes(s *ghttp.Server)
```

---

## 🧪 开发与测试 Development & Testing

> 说明：当前分支未实际执行测试，请按下述命令自行验证。

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/yourusername/go-study2.git
cd go-study2/backend

# 安装后端依赖
go mod download

# 启动后端（默认 8080，若存在 ./build.bat 请先在根执行）
go run main.go -d

# 前端（新终端）
cd ../frontend
npm install
npm run dev  # 默认 3000，代理到 http://localhost:8080/api
```

### 运行测试

```bash
# 后端：运行所有测试
cd backend
go test ./...

# 后端：运行测试并显示覆盖率
go test -cover ./...

# 前端：运行单元与集成测试并输出覆盖率
cd ../frontend
npm test -- --coverage

# 生成后端覆盖率报告（可选）
cd ../backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 代码规范

- 遵循 Go 官方代码规范与 ESLint/Prettier 规则
- 使用 `gofmt` 格式化 Go 代码，前端使用 `npm run lint`（如已配置）
- 所有代码注释和文档使用中文
- 提交前确保前后端测试通过且覆盖率 ≥ 80%

### 性能与边界验证计划

- 性能：针对 `/api/v1/progress`、`/api/v1/quiz/*` 使用 `ab`/`wrk` 目标 p95 < 300ms（预估数据量下）。
- 边界：快速进入/离开页面、网络抖动、并发多窗口、测验中途刷新，需人工或脚本验证。
- 质量检查建议：后端 `go fmt && go vet`（必要时 `golint`）、前端 `npm run lint && npm run build`，执行前确保根目录已运行 `./build.bat`。

### 本地开发工作流

1. 创建功能分支：`git checkout -b feature/your-feature`
2. 编写代码和测试
3. 后端测试：`go test ./...`
4. 前端测试：`cd frontend && npm test -- --coverage`
5. 格式化：`gofmt -w .`；前端运行 `npm run lint`（若配置）
6. 提交代码：`git commit -m "feat: your feature description"`
7. 推送分支：`git push origin feature/your-feature`

---

## 🗺 Roadmap

### 已完成 ✅

- [x] **v0.1** - 基础框架搭建
  - [x] 主菜单系统
  - [x] 模块化架构设计
  - [x] 词法元素模块框架

- [x] **v0.2** - 词法元素内容完善
  - [x] 11个词法元素子主题实现
  - [x] 中文代码注释和说明
  - [x] 单元测试覆盖率达到80%

- [x] **v0.3** - 菜单系统优化
  - [x] 层级菜单结构
  - [x] 交互式子菜单
  - [x] 返回和退出功能

- [x] **v0.4** - HTTP学习模式 🆕
  - [x] 双模式支持（CLI + HTTP）
  - [x] RESTful API实现
  - [x] JSON/HTML响应格式
  - [x] YAML配置管理
  - [x] 请求日志和中间件
  - [x] 内容一致性保证
  - [x] 完整测试覆盖

- [x] **v0.5** - Constants 常量学习模块 🆕
  - [x] 12个常量子主题完整实现
  - [x] 基础常量类型（布尔、符文、整数、浮点、复数、字符串）
  - [x] 常量表达式和类型系统
  - [x] 类型转换和内置函数
  - [x] iota 特性和实现限制
  - [x] CLI和HTTP双模式支持
  - [x] 99%测试覆盖率

### 进行中 🚧

- [ ] **v0.6** - 文档完善
  - [x] README.md更新
  - [ ] 贡献指南
  - [ ] 使用教程视频

### 计划中 📋

- [ ] **v1.0** - 正式版本
  - [ ] 完整的词法元素学习内容
  - [ ] 用户学习进度跟踪
  - [ ] 交互式练习题

- [ ] **v1.1** - 扩展主题
  - [x] Constants 常量学习模块 ✅
  - [ ] 数据类型学习模块
  - [ ] 控制流学习模块
  - [ ] 函数和方法学习模块

- [ ] **v2.0** - 高级功能
  - [ ] 增强的Web界面
  - [ ] 学习进度可视化
  - [ ] 社区分享功能
  - [ ] 多语言支持（英文）

---

## 🤝 贡献指南 Contributing

我们欢迎所有形式的贡献！无论是报告bug、提出新功能建议，还是提交代码改进。

### 如何贡献

1. **Fork 本仓库**
2. **克隆到本地**
   ```bash
   git clone https://github.com/your-username/go-study2.git
   ```
3. **创建功能分支**
   ```bash
   git checkout -b feature/amazing-feature
   ```
4. **编写代码**
   - 遵循项目代码规范
   - 添加必要的测试
   - 确保所有注释和文档使用中文
5. **提交更改**
   ```bash
   git commit -m "feat: 添加某某功能"
   ```
   使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：
   - `feat:` 新功能
   - `fix:` 修复bug
   - `docs:` 文档更新
   - `test:` 测试相关
   - `refactor:` 代码重构
6. **推送到分支**
   ```bash
   git push origin feature/amazing-feature
   ```
7. **创建 Pull Request**
   - 清晰描述你的更改
   - 关联相关的 Issue（如果有）
   - 等待代码审查

### 分支模型

- `main`: 主分支，保持稳定可发布状态
- `feature/*`: 功能开发分支
- `bugfix/*`: Bug修复分支
- `docs/*`: 文档更新分支

### 代码审查标准

- ✅ 代码符合Go语言规范
- ✅ 所有测试通过
- ✅ 测试覆盖率不低于80%
- ✅ 代码注释和文档使用中文
- ✅ 提交信息符合Conventional Commits规范

### 报告问题

如果你发现了bug或有功能建议，请[创建Issue](https://github.com/yourusername/go-study2/issues/new)并提供：

- 问题的详细描述
- 复现步骤（如果是bug）
- 期望的行为
- 实际的行为
- 系统环境信息

---

## 📄 许可证 License

本项目采用 **MIT License** 开源协议。

这意味着你可以：

- ✅ 自由使用本项目
- ✅ 修改源代码
- ✅ 用于商业用途
- ✅ 分发和再授权

唯一的要求是在衍生作品中保留原始的版权声明和许可证声明。

详细信息请查看 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢 Acknowledgements

本项目的开发受到以下项目和资源的启发：

- **[The Go Programming Language Specification](https://go.dev/ref/spec)** - Go语言官方规范，本项目的知识来源
- **[GoFrame](https://goframe.org)** - 优秀的Go语言开发框架
- **[SpecKit](https://github.com/speckit/speckit)** - 项目规范管理方法论，用于本项目的需求和设计管理

特别感谢：

- Go语言社区的所有贡献者
- 所有为本项目提供反馈和建议的学习者
- 使用SpecKit方法论帮助我们保持项目质量和一致性

---

## 📞 联系方式

- **项目主页**: [https://github.com/yourusername/go-study2](https://github.com/yourusername/go-study2)
- **问题反馈**: [GitHub Issues](https://github.com/yourusername/go-study2/issues)
- **讨论区**: [GitHub Discussions](https://github.com/yourusername/go-study2/discussions)

---

<div align="center">

**如果这个项目对你有帮助，请给我们一个 ⭐️ Star！**

Made with ❤️ for Go learners

</div>
---

## 🧭 日志系统快速验证 (Quickstart Validation)

本仓库包含可验证的日志子系统，已在 `backend/docs/` 中提供详细指南与故障排查文档。为便于本地一键验证，项目提供了一个 PowerShell 快速验证脚本，完成以下步骤：

- 启动后端服务（若尚未构建，会尝试使用 `go run` 启动）
- 运行并发压测（使用 `backend/scripts/stress_client.go`）
- 收集并归档 `backend/logs/` 下的日志到 `backend/logs/validation-<timestamp>`

运行（Windows PowerShell）：

```powershell
cd backend\\scripts
powershell -ExecutionPolicy Bypass -File quickstart_validate.ps1
```

结果：脚本会在结束时打印压测统计并将日志归档到 `backend/logs/validation-<timestamp>`，详细手动验证步骤见 `backend/docs/quickstart_validation.md`。

文档索引（已添加）:

- `backend/docs/quickstart_validation.md` — 快速开始验证步骤
- `backend/docs/logging_best_practices.md` — 日志最佳实践（中文）
- `backend/docs/logging_troubleshooting.md` — 日志故障排查指南（中文）
- `backend/README-logging.md` — 日志系统使用说明与示例
