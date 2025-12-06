# Implementation Plan: HTTP学习模式

**Branch**: `003-http-learning-mode` | **Date**: 2025-12-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-http-learning-mode/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

为现有的命令行交互式Go学习工具增加HTTP服务模式。用户可以通过不带参数启动进入命令行模式，或使用`-d`/`--daemon`参数启动HTTP服务。两种模式访问相同的学习内容数据源，确保内容一致性。HTTP服务使用GoFrame框架，所有接口采用POST方法，支持JSON和HTML两种响应格式（通过查询参数`?format`指定）。系统通过YAML配置文件管理HTTP服务参数，实现结构化日志记录，并提供优雅关闭机制。

## Technical Context

**Language/Version**: Go 1.24.5  
**Primary Dependencies**: GoFrame v2.9.5 (HTTP服务框架)  
**Storage**: 无需存储，内容由现有`lexical_elements`包的Display函数生成  
**Testing**: Go标准测试框架 (`go test`)  
**Target Platform**: 跨平台（Linux、Windows、macOS）  
**Project Type**: 单一可执行程序（支持命令行和HTTP两种运行模式）  
**Performance Goals**: HTTP响应时间 <1秒（对于<100KB的章节内容），支持50+并发请求  
**Constraints**: 启动时间 <5秒，优雅关闭时间 <2秒，内存占用 <100MB（正常运行状态）  
**Scale/Scope**: 约11个学习章节（Comments、Tokens、Semicolons、Identifiers、Keywords、Operators、Integers、Floats、Imaginary、Runes、Strings），预期用户并发量 <100

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Simplicity):** ✅ PASS - HTTP模式使用GoFrame框架提供的标准路由和处理器，避免复杂的自定义实现。命令行参数解析使用Go标准库flag包，逻辑清晰简单。
- **Principle II (Comments):** ✅ PASS - 计划中明确要求为HTTP路由层、内容提供者层、配置管理层添加清晰的中文注释，说明各层职责。
- **Principle III (Language):** ✅ PASS - 所有代码注释、日志消息、错误提示、配置文件注释均使用中文。
- **Principle IV (Nesting):** ✅ PASS - 设计采用早返回模式处理错误，避免深层嵌套。HTTP处理器使用中间件模式分离关注点（日志、格式转换、错误处理）。
- **Principle V (YAGNI):** ✅ PASS - 不引入复杂的依赖注入框架或ORM，直接使用文件系统读取内容。配置管理使用GoFrame内置的配置组件，不引入额外的配置库。
- **Principle VI (Testing):** ✅ PASS - 计划包含单元测试策略：内容提供者测试、HTTP处理器测试（使用httptest包）、命令行参数解析测试、配置加载测试。目标覆盖率≥80%。

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
go-study2/
├── main.go                      # 程序入口，解析命令行参数，启动对应模式
├── go.mod                       # Go模块定义
├── go.sum                       # 依赖版本锁定
├── config.yaml                  # HTTP服务配置文件（新增）
│
├── internal/
│   ├── app/
│   │   ├── lexical_elements/    # 现有的词法元素学习内容模块（需重构）
│   │   │   ├── lexical_elements.go  # 主菜单（需重构为支持内容返回）
│   │   │   ├── comments.go      # 注释章节（需重构：DisplayComments() → GetCommentsContent()）
│   │   │   ├── tokens.go        # 标记章节（需重构：DisplayTokens() → GetTokensContent()）
│   │   │   ├── semicolons.go    # 分号章节（需重构）
│   │   │   ├── identifiers.go   # 标识符章节（需重构）
│   │   │   ├── keywords.go      # 关键字章节（需重构）
│   │   │   ├── operators.go     # 运算符章节（需重构）
│   │   │   ├── integers.go      # 整数章节（需重构）
│   │   │   ├── floats.go        # 浮点数章节（需重构）
│   │   │   ├── imaginary.go     # 虚数章节（需重构）
│   │   │   ├── runes.go         # 符文章节（需重构）
│   │   │   ├── strings.go       # 字符串章节（需重构）
│   │   │   └── *_test.go        # 现有测试文件
│   │   │
│   │   └── http_server/         # HTTP服务模块（新增）
│   │       ├── server.go        # HTTP服务器初始化和启动
│   │       ├── router.go        # 路由注册
│   │       ├── handler/         # HTTP处理器
│   │       │   ├── topics.go    # 首页接口处理器
│   │       │   ├── lexical.go   # Lexical Elements菜单接口处理器
│   │       │   └── chapter.go   # 子章节接口处理器
│   │       └── middleware/      # 中间件
│   │           ├── logger.go    # 日志中间件
│   │           └── format.go    # 格式转换中间件
│   │
│   └── config/                  # 配置管理模块（新增）
│       └── config.go            # 配置加载和验证
│
└── tests/                       # 测试文件（新增）
    ├── unit/
    │   ├── lexical_refactor_test.go  # 重构后的lexical_elements测试
    │   ├── handler_test.go      # HTTP处理器单元测试
    │   └── config_test.go       # 配置加载单元测试
    └── integration/
        └── http_test.go         # HTTP服务集成测试
```

**Structure Decision**: 采用单一Go项目结构。**关键变化**：
1. **重构 `internal/app/lexical_elements`**: 将所有`Display*()`函数改为`Get*Content()`函数，返回字符串而非直接打印
2. **新增 `internal/app/http_server`**: 处理HTTP服务逻辑，调用重构后的`Get*Content()`函数
3. **新增 `internal/config`**: 管理YAML配置文件
4. **移除 `internal/content`**: 不需要单独的内容提供者，直接使用重构后的lexical_elements函数

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项 - 所有Constitution原则检查均通过。

---

## Post-Phase 1 Constitution Re-Check

*Re-evaluation after completing design artifacts (research.md, data-model.md, contracts/, quickstart.md)*

### 设计审查

- **Principle I (Simplicity):** ✅ PASS - 设计保持简单，使用GoFrame标准组件，避免过度设计。数据模型清晰，仅5个核心实体。
- **Principle II (Comments):** ✅ PASS - 所有设计文档和代码示例均包含详细的中文注释，说明各组件职责。
- **Principle III (Language):** ✅ PASS - 所有文档、API响应、错误消息均使用中文。
- **Principle IV (Nesting):** ✅ PASS - 设计采用中间件模式和早返回模式，避免深层嵌套逻辑。
- **Principle V (YAGNI):** ✅ PASS - 未引入不必要的复杂性，如缓存机制标记为"可选优化"，不在初始实现中包含。
- **Principle VI (Testing):** ✅ PASS - 测试策略明确，包含单元测试和集成测试，目标覆盖率≥80%。

### 新增检查项

- **Principle VII (Single Responsibility):** ✅ PASS - 每个模块职责单一：`content`负责内容加载，`http_server`负责HTTP服务，`config`负责配置管理。
- **Principle X (Explicit Error Handling):** ✅ PASS - 所有错误场景均有明确处理，包括404、500、配置缺失、端口占用等。
- **Principle XVII (Hierarchical Menu):** ✅ PASS - HTTP模式通过API路径实现层次化导航（主题→章节→内容），与命令行模式的菜单结构对应。

**结论**: 设计完全符合项目Constitution，无需调整。

---

## Phase 0-1 完成总结

### 已完成的工作

#### Phase 0: Research (研究)
✅ **research.md** - 完成8个关键技术问题的研究：
  - GoFrame HTTP服务器最佳实践
  - POST接口设计模式
  - 多格式响应处理（JSON/HTML）
  - YAML配置管理
  - 优雅关闭机制
  - 并发安全策略
  - 测试策略
  - 重构现有Display函数

#### Phase 1: Design & Contracts (设计与契约)
✅ **data-model.md** - 定义4个核心实体：
  - Topic (主题)
  - Chapter (章节)
  - Content (内容生成) - 通过`lexical_elements`包的`Get*Content()`函数
  - ServerConfig (服务配置)
  - Response (响应)

✅ **contracts/api-spec.md** - 定义3个核心API接口：
  - POST /api/v1/topics - 获取主题列表
  - POST /api/v1/topic/lexical_elements - 获取Lexical Elements菜单
  - POST /api/v1/topic/lexical_elements/{chapter} - 获取章节内容

✅ **quickstart.md** - 提供完整的快速入门指南：
  - 配置文件设置
  - 两种启动方式（命令行/HTTP）
  - 使用示例（浏览器、curl、Go代码）
  - 常见问题解答
  - 开发者指南

✅ **Agent Context Update** - 更新GEMINI.md，添加技术栈信息

### 关键技术决策

| 决策点 | 选择 | 理由 |
|--------|------|------|
| HTTP框架 | GoFrame v2.9.5 | 项目已有依赖，功能完整，适合初学者 |
| 接口方法 | 全部使用POST | 用户明确要求 |
| 响应格式 | JSON/HTML双格式 | 通过`?format`参数控制，满足不同使用场景 |
| 配置格式 | YAML | 符合规范FR-019，易读易写 |
| 配置策略 | 必填项无默认值 | 强制用户明确配置，避免隐式行为 |
| 并发处理 | 依赖GoFrame内置 | 无需额外控制，简化实现 |
| 测试策略 | 分层测试 | 单元测试+集成测试，覆盖率≥80% |
| 内容来源 | 重构现有Display函数 | 保持向后兼容，确保内容一致性 |

### 下一步行动

根据workflow规定，`/speckit.plan`命令在Phase 1完成后停止。

**接下来应该**:
1. 运行`/speckit.tasks`命令生成tasks.md（任务分解）
2. 运行`/speckit.implement`命令执行实现

**或者**:
- 如需进一步澄清需求，运行`/speckit.clarify`
- 如需生成检查清单，运行`/speckit.checklist`

---

## 附录：生成的文件清单

```
specs/003-http-learning-mode/
├── plan.md                          # 本文件
├── spec.md                          # 功能规范（已存在）
├── research.md                      # Phase 0输出 ✅
├── data-model.md                    # Phase 1输出 ✅
├── quickstart.md                    # Phase 1输出 ✅
└── contracts/
    └── api-spec.md                  # Phase 1输出 ✅
```

**计划状态**: Phase 0-1 完成 ✅  
**下一阶段**: Phase 2 - 任务分解（需运行`/speckit.tasks`）
