# Go-Study2 - Go语言词法元素学习工具

> 一个交互式命令行工具，帮助Go语言学习者系统掌握词法元素知识，通过菜单驱动的方式提供代码示例和详细解释。

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

- 🎯 **菜单驱动界面** - 清晰的层级菜单，轻松导航各个知识点
- 📖 **全面覆盖** - 涵盖Go语言规范中所有词法元素子主题
- 💻 **可运行示例** - 每个知识点都配有可直接运行的代码示例
- 🇨🇳 **中文注释** - 所有代码注释和说明均为中文，降低学习门槛
- 🧪 **高测试覆盖率** - 80%以上的单元测试覆盖率，保证代码质量
- 🔌 **易于扩展** - 模块化设计，可轻松添加新的学习主题
- 🚀 **零依赖运行** - 编译后的可执行文件无需额外依赖

---

## 🧱 技术栈 Tech Stack

- **语言**: Go 1.24.5
- **框架**: GoFrame v2.9.5（最小化使用）
- **构建工具**: Go Modules
- **测试框架**: Go标准测试库
- **开发环境**: 支持 Windows/Linux/macOS

---

## 🚀 快速开始 Quick Start

**30秒快速体验：**

```bash
# 克隆仓库
git clone https://github.com/yourusername/go-study2.git

# 进入项目目录
cd go-study2

# 运行程序
go run main.go
```

**预期输出：**

```
Go Lexical Elements Learning Tool
---------------------------------
Please select a topic to study:
0. Lexical elements
q. Quit

Enter your choice: 
```

输入 `0` 即可开始学习词法元素！

---

## 📦 安装 Installation

### 方式一：从源码运行（推荐用于学习）

```bash
git clone https://github.com/yourusername/go-study2.git
cd go-study2
go run main.go
```

### 方式二：编译后运行

```bash
git clone https://github.com/yourusername/go-study2.git
cd go-study2
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

1. **启动程序**：运行 `go run main.go` 或编译后的可执行文件
2. **选择主题**：在主菜单中输入 `0` 选择"词法元素"
3. **浏览子主题**：在词法元素菜单中选择具体的子主题（如注释、标识符、关键字等）
4. **查看示例**：程序会显示该主题的代码示例和详细解释
5. **返回或退出**：输入 `b` 返回上级菜单，输入 `q` 退出程序

### 交互示例

```
Go Lexical Elements Learning Tool
---------------------------------
Please select a topic to study:
0. Lexical elements
q. Quit

Enter your choice: 0

Lexical Elements Menu
---------------------
0. Comments (注释)
1. Tokens (标记)
2. Semicolons (分号)
3. Identifiers (标识符)
4. Keywords (关键字)
5. Operators and punctuation (运算符和标点)
6. Integer literals (整数字面量)
7. Floating-point literals (浮点数字面量)
8. Imaginary literals (虚数字面量)
9. Rune literals (字符字面量)
10. String literals (字符串字面量)
b. Back to main menu
q. Quit

Enter your choice: 4
```

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

---

## 🗂 项目结构 Project Structure

```
go-study2/
├── main.go                          # 主入口文件，包含菜单逻辑
├── main_test.go                     # 主程序测试文件
├── go.mod                           # Go模块依赖管理
├── go.sum                           # 依赖校验文件
├── internal/                        # 内部包（不对外暴露）
│   └── app/
│       └── lexical_elements/        # 词法元素学习模块
│           ├── comments.go          # 注释相关示例
│           ├── tokens.go            # 标记相关示例
│           ├── semicolons.go        # 分号规则示例
│           ├── identifiers.go       # 标识符示例
│           ├── keywords.go          # 关键字示例
│           ├── operators.go         # 运算符示例
│           ├── integer_literals.go  # 整数字面量示例
│           ├── float_literals.go    # 浮点数字面量示例
│           ├── imaginary_literals.go # 虚数字面量示例
│           ├── rune_literals.go     # 字符字面量示例
│           ├── string_literals.go   # 字符串字面量示例
│           ├── menu.go              # 词法元素菜单逻辑
│           └── *_test.go            # 各模块测试文件
├── specs/                           # 功能规格说明
│   ├── 001-go-learn-lexical-elements/
│   └── 002-lexical-menu-structure/
├── doc/                             # 文档目录
│   ├── README模板.md
│   ├── Go开发工具生态全景图.md
│   └── The Go-1.24 Programming Language Specification.md
├── .specify/                        # 项目规范和模板
├── .agent/                          # AI辅助开发工作流
└── README.md                        # 本文件
```

**目录说明：**

- `internal/app/lexical_elements/`: 核心学习内容，每个文件对应一个词法元素子主题
- `specs/`: 使用SpecKit方法论管理的功能规格文档
- `doc/`: 参考文档和学习资料
- `main.go`: 应用入口，实现菜单系统和模块调度

---

## ⚙️ 配置 Configuration

**当前版本暂不涉及配置文件。**

本项目采用"零配置"设计理念，开箱即用。所有学习内容都硬编码在代码中，无需外部配置。

**未来可能支持的配置项：**

- 界面语言切换（中文/英文）
- 代码示例输出格式
- 学习进度记录

---

## 📖 API 文档 API Reference

**当前版本暂不涉及公开API。**

本项目是一个独立的CLI工具，不提供对外API。所有功能通过命令行交互完成。

**内部包结构：**

```go
package lexical_elements

// DisplayMenu 显示词法元素子菜单
func DisplayMenu(stdin io.Reader, stdout, stderr io.Writer)

// 各子主题的Display函数
func DisplayComments(stdout io.Writer)
func DisplayTokens(stdout io.Writer)
func DisplaySemicolons(stdout io.Writer)
// ... 其他子主题
```

---

## 🧪 开发与测试 Development & Testing

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/yourusername/go-study2.git
cd go-study2

# 安装依赖（如果有）
go mod download

# 运行程序
go run main.go
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成详细的覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 代码规范

- 遵循Go官方代码规范
- 使用 `gofmt` 格式化代码
- 所有代码注释和文档使用中文
- 提交前确保测试通过且覆盖率 ≥ 80%

### 本地开发工作流

1. 创建功能分支：`git checkout -b feature/your-feature`
2. 编写代码和测试
3. 运行测试：`go test ./...`
4. 格式化代码：`gofmt -w .`
5. 提交代码：`git commit -m "feat: your feature description"`
6. 推送分支：`git push origin feature/your-feature`

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

### 进行中 🚧

- [ ] **v0.4** - 文档完善
  - [x] README.md编写
  - [ ] 贡献指南
  - [ ] 使用教程视频

### 计划中 📋

- [ ] **v1.0** - 正式版本
  - [ ] 完整的词法元素学习内容
  - [ ] 用户学习进度跟踪
  - [ ] 交互式练习题

- [ ] **v1.1** - 扩展主题
  - [ ] 数据类型学习模块
  - [ ] 控制流学习模块
  - [ ] 函数和方法学习模块

- [ ] **v2.0** - 高级功能
  - [ ] Web界面版本
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
