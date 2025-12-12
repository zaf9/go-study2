# Implementation Plan: Go-Study2 日志系统重构

**Branch**: `012-logging-system` | **Date**: 2025-12-12 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/012-logging-system/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

本特性旨在重构 Go-Study2 项目的日志系统,充分利用 GoFrame v2.9.5 框架的原生日志能力(glog),建立统一的日志配置、全链路请求追踪、关键操作埋点和日志查询支持。技术方案采用 YAML 配置文件管理多个日志实例(app/access/error/slow),通过 HTTP 中间件自动注入 TraceID 实现请求追踪,利用 ORM Handler 监控慢查询,并提供统一的日志记录接口封装。系统采用 fail-fast 策略处理配置错误和权限问题,确保日志系统的可靠性和可观测性。

## Technical Context

**Language/Version**: Go 1.24.5  
**Primary Dependencies**: GoFrame v2.9.5 (glog, ghttp, gdb, gtrace, gerror)  
**Storage**: 本地文件系统 (./logs/), 支持按日期和大小分割  
**Testing**: Go 标准测试框架 (testing package), 目标覆盖率 ≥80%  
**Target Platform**: Linux/Windows 服务器, 支持 CLI 和 HTTP 双模式  
**Project Type**: Single project (Go 后端应用)  
**Performance Goals**: 
- 日志写入开销 <10% 请求响应时间 (1000 并发)
- TraceID 查询响应时间 <30 秒
- 日志查询响应时间 <5 秒 (1GB 文件)
**Constraints**: 
- 配置文件缺失或无效时拒绝启动 (fail-fast)
- 日志目录权限不足时拒绝启动 (fail-fast)
- 单个日志文件大小 ≤100MB (可配置)
- 日志保留期限默认 30 天 (可配置)
**Scale/Scope**: 
- 支持 4 种日志实例 (app/access/error/slow)
- 支持 1000 并发请求的日志记录
- 单个日志文件最大 1GB (查询性能要求)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (代码质量与可维护性):** ✅ 方案采用 GoFrame 原生能力,避免过度封装,职责清晰(配置/中间件/辅助方法分离),易于测试
- **Principle II (显式错误处理):** ✅ 所有错误场景明确处理:配置无效拒绝启动,权限不足拒绝启动,TraceID 丢失生成新 ID 并警告,日志写入失败降级到 stderr
- **Principle III/XXI/XXXVI (全面测试):** ✅ 规划单元测试覆盖率 ≥80%,包含配置验证、中间件、日志辅助方法、TraceID 传递等测试用例
- **Principle IV (单一职责):** ✅ 配置管理(logger.go)、中间件(middleware/)、辅助方法(helper.go)、实例管理(instance.go)职责分离
- **Principle V/XV (一致文档与中文要求):** ✅ 所有注释和文档使用中文,包含配置示例、使用文档、最佳实践指南
- **Principle VI (YAGNI):** ✅ 不实现应用层访问控制(依赖文件系统权限),不实现敏感信息自动脱敏(由开发者处理),不实现复杂查询语法(仅文本搜索)
- **Principle VII (安全优先):** ✅ 配置文件验证,目录权限检查,敏感信息脱敏规范(文档说明),日志访问依赖文件系统权限
- **Principle VIII/XVIII (可预测结构):** ✅ 遵循标准 Go 布局,日志包位于 `internal/logic/logger/`,配置文件位于 `config/`,日志文件位于 `./logs/`
- **Principle IX (依赖纪律):** ✅ 仅依赖 GoFrame v2.9.5 框架,无额外第三方日志库
- **Principle X (性能优化):** ✅ 异步日志写入,生产环境关闭 File()/Line(),日志采样(高并发场景),大对象截断(5KB 限制)
- **Principle XI (文档同步):** ✅ 完成后更新 README.md 的日志系统说明、配置指南、使用示例
- **Principle XIV (清晰分层注释):** ✅ 配置层、中间件层、业务层、数据层的日志职责均有中文注释说明
- **Principle XVI (浅层逻辑):** ✅ 使用卫语句和早返回,避免深层嵌套,中间件逻辑扁平化
- **Principle XVII (一致开发者体验):** ✅ 提供日志工具函数库、最佳实践文档、配置示例,降低使用门槛
- **Principle XIX (包级 README):** ✅ `internal/logic/logger/` 包含 README.md 说明日志系统架构、配置方法、使用示例
- **Principle XX (代码质量执行):** ✅ 执行 go fmt/go vet/golint/go mod tidy,集成到 CI/CD
- **Principle XXII (分层菜单导航):** N/A (本特性为基础设施,不涉及菜单导航)
- **Principle XXIII (双学习模式):** N/A (本特性为基础设施,不涉及学习内容)
- **Principle XXIV (层次化章节结构):** N/A (本特性为基础设施,不涉及章节结构)
- **Principle XXV (HTTP/CLI 一致性):** ✅ 日志中间件同时支持 HTTP 和 CLI 模式,共享日志配置和实例

## Project Structure

### Documentation (this feature)

```text
specs/012-logging-system/
├── spec.md              # Feature specification
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output - GoFrame glog 最佳实践研究
├── data-model.md        # Phase 1 output - 日志配置和实体模型
├── quickstart.md        # Phase 1 output - 快速开始指南
├── contracts/           # Phase 1 output
│   ├── cli-integration.md    # CLI 集成规范
│   └── http-middleware.md    # HTTP 中间件规范
├── checklists/
│   └── requirements.md  # Specification quality checklist
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
go-study2/
├── backend/                         # 后端主目录
│   ├── go.mod / go.sum              # Go 模块定义与依赖
│   ├── main.go / main_test.go       # 主入口与测试
│   │
│   ├── configs/                     # 配置目录
│   │   ├── config.yaml              # 基础配置(包含 logger 配置)
│   │   ├── config.dev.yaml          # 开发环境配置
│   │   ├── config.test.yaml         # 测试环境配置
│   │   └── config.prod.yaml         # 生产环境配置
│   │
│   ├── internal/
│   │   ├── app/                     # 应用层
│   │   │   ├── http_server/         # HTTP 服务器(handler、middleware、router)
│   │   │   │   └── middleware/      # HTTP 中间件
│   │   │   │       ├── access_log.go       # 访问日志中间件
│   │   │   │       ├── access_log_test.go  # 访问日志测试
│   │   │   │       ├── panic_recovery.go   # Panic 恢复中间件
│   │   │   │       └── panic_recovery_test.go # Panic 恢复测试
│   │   │   │
│   │   │   ├── lexical_elements/    # 词法元素内容
│   │   │   ├── constants/           # 常量模块内容
│   │   │   └── ...                  # 其他学习主题
│   │   │
│   │   ├── domain/                  # 领域层(user/progress/quiz 实体与服务)
│   │   │
│   │   ├── infrastructure/          # 基础设施层
│   │   │   ├── database/            # 数据库连接与迁移
│   │   │   ├── repository/          # 仓储实现
│   │   │   └── logger/              # 日志系统(新增)
│   │   │       ├── README.md           # 日志系统使用文档
│   │   │       ├── logger.go           # 日志实例管理和初始化
│   │   │       ├── logger_test.go      # 日志实例测试
│   │   │       ├── config.go           # 配置验证和加载
│   │   │       ├── config_test.go      # 配置验证测试
│   │   │       ├── helper.go           # 日志辅助方法封装
│   │   │       ├── helper_test.go      # 辅助方法测试
│   │   │       ├── traceid.go          # TraceID 生成和传递
│   │   │       └── traceid_test.go     # TraceID 测试
│   │   │
│   │   ├── pkg/                     # 共享工具(jwt、password)
│   │   └── config/                  # 配置加载与校验
│   │
│   ├── logs/                        # 日志文件目录(运行时创建)
│   │   ├── app-2025-12-12.log       # 应用日志
│   │   ├── access/                  # 访问日志目录
│   │   ├── error/                   # 错误日志目录
│   │   └── slow/                    # 慢查询日志目录
│   │
│   ├── tests/                       # 测试目录
│   │   ├── unit/                    # 单元测试
│   │   ├── integration/             # 集成测试
│   │   └── contract/                # 契约测试
│   │
│   ├── docs/                        # 后端文档
│   └── scripts/                     # 工具脚本
│
├── frontend/                        # 前端目录(不涉及本特性)
├── specs/                           # 特性规范目录
└── docs/                            # 项目文档
```

**Structure Decision**: 采用前后端分离的项目结构,后端位于 `backend/` 目录。日志系统作为基础设施,位于 `backend/internal/infrastructure/logger/` 包,HTTP 中间件位于 `backend/internal/app/http_server/middleware/`,配置文件位于 `backend/configs/`,日志输出位于 `backend/logs/`。这种结构符合 DDD 分层架构和 Principle VIII/XVIII 的可预测性要求,便于维护和扩展。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无宪章违规需要说明。所有设计决策均符合项目宪章要求。

---

## Phase 0: Research & Technical Decisions

### Research Topics

以下研究任务将在 `research.md` 中详细记录:

1. **GoFrame glog 配置最佳实践**
   - 研究 GoFrame v2.9.5 的 glog 配置选项和默认行为
   - 确定多实例日志配置方法 (app/access/error/slow)
   - 确定日志分割策略配置 (按日期/大小)
   - 确定日志保留策略实现方法

2. **GoFrame Context 和 TraceID 集成**
   - 研究 `gtrace.NewID()` 和 `gtrace.GetTraceID(ctx)` 的使用
   - 研究 Context 在 HTTP 请求中的传递机制
   - 研究 `ctxKeys` 配置项的使用方法
   - 确定 TraceID 丢失时的检测和恢复策略

3. **GoFrame 中间件开发模式**
   - 研究 `ghttp.MiddlewareHandlerResponse` 的使用
   - 研究中间件执行顺序和 Context 传递
   - 研究 Panic 恢复机制和堆栈记录
   - 确定访问日志中间件的实现模式

4. **GoFrame ORM 日志 Handler**
   - 研究 `gdb.Handler` 的注册和使用
   - 研究 SQL 执行前后的 Hook 机制
   - 研究慢查询检测和记录方法
   - 确定数据库日志的格式和字段

5. **日志配置验证策略**
   - 研究 YAML 配置文件的加载和验证
   - 研究目录权限检查方法 (os.Stat, os.OpenFile)
   - 研究配置错误的详细诊断信息生成
   - 确定启动时配置验证的实现流程

6. **日志性能优化策略**
   - 研究 GoFrame 异步日志的配置和使用
   - 研究日志采样策略 (高并发场景)
   - 研究大对象日志截断方法
   - 研究生产环境性能优化配置

**输出**: `research.md` 文件,包含所有研究结果、技术决策、替代方案对比和最佳实践建议

---

## Phase 1: Design & Contracts

### Data Model (`data-model.md`)

定义日志系统的核心数据结构:

1. **LogConfig** - 日志配置结构
   - 字段: Level, Stdout, TimeFormat, CtxKeys, Instances (map[string]InstanceConfig)
   - 验证规则: Level 必须为有效级别,Path 必须可写,RotateExpire 必须为有效时长

2. **InstanceConfig** - 日志实例配置
   - 字段: Path, File, Format, Level, RotateSize, RotateExpire, RotateBackupLimit
   - 关系: 属于 LogConfig

3. **LogRecord** - 日志记录结构 (JSON 格式)
   - 必须字段: Time, Level, TraceId, File, Content
   - 扩展字段: Module, Action, UserId, Duration, Error, Extra

4. **TraceContext** - 追踪上下文
   - 字段: TraceID, RequestID, UserID
   - 生命周期: HTTP 请求开始到结束

5. **AccessLogEntry** - 访问日志条目
   - 字段: Method, Path, Query, StatusCode, Duration, IP, UserAgent, TraceID

### API Contracts (`contracts/`)

#### 1. CLI Integration (`cli-integration.md`)

定义 CLI 模式下的日志集成规范:

```go
// CLI 应用入口集成日志初始化
func main() {
    // 1. 初始化日志系统
    if err := logger.Initialize(); err != nil {
        fmt.Fprintf(os.Stderr, "日志系统初始化失败: %v\n", err)
        os.Exit(1)
    }
    
    // 2. 创建带 TraceID 的 Context
    ctx := gctx.New()
    
    // 3. 记录应用启动日志
    logger.LogInfo(ctx, "app", "应用启动", nil)
    
    // 4. 运行应用
    app.Run(ctx)
}

// CLI 命令中使用日志
func handleCommand(ctx context.Context, args []string) error {
    logger.LogInfo(ctx, "cli", "执行命令", g.Map{"args": args})
    // ... 业务逻辑
    return nil
}
```

#### 2. HTTP Middleware (`http-middleware.md`)

定义 HTTP 中间件的接口规范:

```go
// 访问日志中间件
func AccessLog(r *ghttp.Request) {
    // 1. 请求开始时间
    startTime := time.Now()
    
    // 2. 生成或提取 TraceID
    traceID := extractOrGenerateTraceID(r)
    
    // 3. 注入 TraceID 到 Context
    ctx := context.WithValue(r.Context(), "TraceId", traceID)
    r.SetCtx(ctx)
    
    // 4. 记录请求日志
    logger.AccessLogger().Info(ctx, "请求开始", g.Map{
        "method": r.Method,
        "path": r.URL.Path,
        "ip": r.GetClientIp(),
    })
    
    // 5. 继续处理请求
    r.Middleware.Next()
    
    // 6. 记录响应日志
    duration := time.Since(startTime).Milliseconds()
    logger.AccessLogger().Info(ctx, "请求完成", g.Map{
        "status": r.Response.Status,
        "duration": duration,
    })
}

// Panic 恢复中间件
func PanicRecovery(r *ghttp.Request) {
    defer func() {
        if err := recover(); err != nil {
            // 记录 Panic 堆栈到错误日志
            logger.ErrorLogger().Error(r.Context(), "Panic 恢复", g.Map{
                "error": err,
                "stack": string(debug.Stack()),
            })
            
            // 返回 500 错误
            r.Response.WriteStatus(500, g.Map{
                "code": 500,
                "message": "服务器内部错误",
            })
        }
    }()
    
    r.Middleware.Next()
}
```

### Quickstart Guide (`quickstart.md`)

提供快速开始指南,包含:

1. **配置文件设置**
   - 复制配置模板
   - 修改日志级别和路径
   - 配置多环境差异

2. **代码集成步骤**
   - 在 main.go 中初始化日志系统
   - 注册 HTTP 中间件
   - 注册 ORM Handler
   - 在业务代码中使用日志辅助方法

3. **验证和测试**
   - 启动应用验证配置
   - 发送 HTTP 请求验证 TraceID
   - 检查日志文件格式
   - 验证慢查询日志

4. **常见问题排查**
   - 配置文件格式错误
   - 目录权限不足
   - TraceID 未传递
   - 日志文件未分割

### Agent Context Update

运行 `.specify/scripts/powershell/update-agent-context.ps1 -AgentType gemini` 更新 GEMINI.md,添加:

```markdown
## Active Technologies
- Go 1.24.5 + GoFrame v2.9.5 (glog 日志系统重构)

## Recent Changes
- 012-logging-system: 重构日志系统,采用 GoFrame 原生 glog,支持多实例、TraceID 追踪、慢查询监控
```

---

## Phase 2: Task Generation

**注意**: 任务生成由 `/speckit.tasks` 命令完成,不在本计划范围内。

任务将按以下优先级分解:

### Phase 1: 配置与初始化 (P0)
- T001: 创建日志配置结构和验证逻辑
- T002: 实现配置文件加载和环境差异化
- T003: 实现目录权限检查和启动验证
- T004: 初始化多个日志实例 (app/access/error/slow)
- T005: 编写配置验证单元测试

### Phase 2: 中间件开发 (P0)
- T006: 实现 HTTP 访问日志中间件
- T007: 实现 TraceID 生成和注入逻辑
- T008: 实现 Panic 恢复中间件
- T009: 实现 ORM 慢查询 Handler
- T010: 编写中间件单元测试和集成测试

### Phase 3: 辅助方法封装 (P1)
- T011: 封装 LogInfo/LogError/LogSlow/LogBiz 方法
- T012: 实现 TraceID 传递和丢失检测
- T013: 实现日志字段规范化
- T014: 编写辅助方法单元测试

### Phase 4: 代码改造 (P1)
- T015: 在 main.go 中集成日志初始化
- T016: 注册 HTTP 中间件到路由
- T017: 注册 ORM Handler 到数据库连接
- T018: 改造现有错误处理点添加日志
- T019: 添加关键业务流程日志埋点

### Phase 5: 文档和测试 (P2)
- T020: 编写 logger 包 README.md
- T021: 编写日志使用最佳实践文档
- T022: 编写配置示例和环境配置
- T023: 执行性能压测和优化
- T024: 更新项目 README.md

### Phase 6: 日志保留策略 (P3)
- T025: 实现日志文件清理定时任务
- T026: 实现可配置的保留期限
- T027: 编写清理任务测试

---

## Implementation Notes

### 关键技术决策

1. **配置验证策略**: 采用 fail-fast 策略,启动时验证配置文件存在性、格式有效性、目录权限,任何错误立即退出并返回详细诊断信息

2. **TraceID 传递机制**: 
   - HTTP 请求: 从 Header 提取或生成新 TraceID,注入到 Context
   - CLI 模式: 在入口创建带 TraceID 的 Context
   - 传递中断: 检测到 TraceID 丢失时生成新 ID 并记录 WARNING

3. **日志实例管理**: 
   - 使用全局单例模式管理 4 个日志实例
   - 通过 `logger.DefaultLogger()`, `logger.AccessLogger()` 等方法访问
   - 避免重复创建实例

4. **性能优化**: 
   - 生产环境关闭 `File()` 和 `Line()` 调用
   - 启用异步日志写入 (`SetAsync(true)`)
   - 大对象日志截断 (限制 5KB)
   - 高并发场景日志采样

5. **日志保留策略**: 
   - 使用 GoFrame 内置的 `rotateBackupExpire` 配置
   - 默认保留 30 天,可通过配置调整
   - 自动清理过期日志文件

### GoFrame 特有优化

1. **利用框架特性**:
   - `gctx.New()`: 创建带 TraceID 的 Context
   - `gtrace.GetTraceID(ctx)`: 获取链路 ID
   - `gerror.Wrap()`: 包装错误并自动记录堆栈
   - `g.Log()`: 链式调用,代码简洁

2. **避免的坑**:
   - 不混用 `fmt.Println`,统一使用 `g.Log()`
   - 避免在循环中频繁记录日志
   - 敏感信息必须脱敏 (密码、Token、身份证号)
   - 不记录整个 Request/Response,仅记录摘要

### 验收标准

1. **功能验收**:
   - ✅ 所有 HTTP 请求有完整访问日志 (请求+响应)
   - ✅ 所有错误有堆栈信息和 Context
   - ✅ 数据库慢查询可定位 SQL 和参数
   - ✅ 日志可通过 TraceID 串联整个请求链路
   - ✅ 支持多环境配置切换

2. **质量验收**:
   - ✅ 日志格式符合 JSON 规范 100%
   - ✅ 敏感信息已脱敏 100%
   - ✅ 关键业务节点覆盖率 >90%
   - ✅ 错误日志覆盖率 100%
   - ✅ 单元测试覆盖率 ≥80%

3. **性能验收**:
   - ✅ 1000 并发请求下日志开销 <10%
   - ✅ TraceID 查询响应时间 <30 秒
   - ✅ 日志查询响应时间 <5 秒 (1GB 文件)

---

## Next Steps

1. **立即执行**: 运行 `/speckit.plan` 的 Phase 0 研究任务,生成 `research.md`
2. **Phase 1 设计**: 完成 `data-model.md`, `contracts/`, `quickstart.md`
3. **任务分解**: 运行 `/speckit.tasks` 生成详细的实施任务列表
4. **开始实施**: 按优先级执行任务,从配置验证和日志实例初始化开始

**预计工作量**: 5-7 个工作日 (包含编码、测试、文档、Code Review)
