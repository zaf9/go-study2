# Tasks: Go-Study2 日志系统重构

**Input**: Design documents from `/specs/012-logging-system/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Per the constitution, features MUST have at least 80% unit test coverage. Test tasks are **MANDATORY**.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `backend/` directory
- **Config**: `backend/configs/`
- **Logger**: `backend/internal/infrastructure/logger/`
- **Middleware**: `backend/internal/app/http_server/middleware/`
- **Tests**: `backend/tests/`

## Constitution Guardrails

- 所有注释与用户文档相关任务必须产出中文内容,且保持清晰一致(Principle V/XV)。
- 需规划达到>=80%测试覆盖,各包包含 *_test.go 与示例,前端核心组件同样达标(Principle III/XXI/XXXVI)。
- 目录/文件/函数保持单一职责与可预测结构,遵循标准 Go 布局(仅根目录 main, go.mod/go.sum 完整)并补齐包 README(Principle IV/VIII/XVIII/XIX)。
- 外部依赖与复杂度最小化,错误处理显式,避免 YAGNI(Principle II/VI/IX)。
- 完成后需包含更新 README 等文档的任务(Principle XI)。

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目初始化和基础结构准备

- [X] T001 创建日志系统目录结构 `backend/internal/infrastructure/logger/`
- [X] T002 创建 HTTP 中间件目录结构 `backend/internal/app/http_server/middleware/`
- [X] T003 [P] 创建日志输出目录 `backend/logs/` 及子目录 (access/, error/, slow/)
- [X] T004 [P] 创建测试目录结构 `backend/tests/unit/logger/` 和 `backend/tests/integration/middleware/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 核心基础设施,必须在所有用户故事之前完成

**⚠️ CRITICAL**: 所有用户故事工作必须等待此阶段完成

- [ ] T005 创建日志配置数据结构 `backend/internal/infrastructure/logger/config.go` (LoggerConfig, InstanceConfig)
- [ ] T006 实现配置文件加载逻辑 `backend/internal/infrastructure/logger/config.go` (LoadConfig 函数)
- [ ] T007 实现配置验证逻辑 `backend/internal/infrastructure/logger/config.go` (Validate 函数)
- [ ] T008 实现目录权限检查 `backend/internal/infrastructure/logger/config.go` (checkDirectoryPermission 函数)
- [ ] T009 [P] 编写配置加载和验证单元测试 `backend/tests/unit/logger/config_test.go`

 **⚠️ CRITICAL**: 所有用户故事工作必须等待此阶段完成

 - [X] T005 创建日志配置数据结构 `backend/internal/infrastructure/logger/config.go` (LoggerConfig, InstanceConfig)
 - [X] T006 实现配置文件加载逻辑 `backend/internal/infrastructure/logger/config.go` (LoadConfig 函数)
 - [X] T007 实现配置验证逻辑 `backend/internal/infrastructure/logger/config.go` (Validate 函数)
 - [X] T008 实现目录权限检查 `backend/internal/infrastructure/logger/config.go` (checkDirectoryPermission 函数)
 - [X] T009 [P] 编写配置加载和验证单元测试 `backend/tests/unit/logger/config_test.go`

**Checkpoint**: 配置基础设施就绪 - 用户故事实现可以并行开始

---

## Phase 3: User Story 1 - 统一日志配置管理 (Priority: P1) 🎯 MVP

**Goal**: 通过 YAML 配置文件统一管理所有日志实例(app/access/error/slow),支持日志级别、格式、路径、分割策略配置

**Independent Test**: 修改配置文件中的日志级别、输出路径等参数,启动应用,验证日志按配置输出到指定位置和格式

### Tests for User Story 1 (MANDATORY) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

 - [ ] T010 [P] [US1] 编写日志实例初始化单元测试,包括链式调用方式测试 `backend/tests/unit/logger/logger_test.go`
 - [ ] T011 [P] [US1] 编写日志级别过滤测试 `backend/tests/unit/logger/logger_test.go`
 - [ ] T012 [P] [US1] 编写日志格式化测试(JSON/Text) `backend/tests/unit/logger/logger_test.go`
 - [ ] T013 [P] [US1] 编写多实例独立配置测试 `backend/tests/unit/logger/logger_test.go`
 - [ ] T014 [P] [US1] 编写配置文件缺失/无效启动失败测试 `backend/tests/integration/logger/config_validation_test.go`
 
  - [X] T010 [P] [US1] 编写日志实例初始化单元测试,包括链式调用方式测试 `backend/tests/unit/logger/logger_test.go`
  - [X] T011 [P] [US1] 编写日志级别过滤测试 `backend/tests/unit/logger/logger_test.go`
  - [X] T012 [P] [US1] 编写日志格式化测试(JSON/Text) `backend/tests/unit/logger/logger_test.go`
  - [X] T013 [P] [US1] 编写多实例独立配置测试 `backend/tests/unit/logger/logger_test.go`
  - [X] T014 [P] [US1] 编写配置文件缺失/无效启动失败测试 `backend/tests/integration/logger/config_validation_test.go`

### Implementation for User Story 1

- [ ] T015 [US1] 实现日志实例管理器 `backend/internal/infrastructure/logger/logger.go` (Initialize, GetInstance 函数)
- [ ] T016 [US1] 实现多日志实例初始化 `backend/internal/infrastructure/logger/logger.go` (initializeInstances 函数)
 - [ ] T017 [US1] 集成 GoFrame glog 配置,包括异步日志写入配置 `backend/internal/infrastructure/logger/logger.go` (configureGLog 函数)
 - [ ] T018 [US1] 实现日志分割策略配置 `backend/internal/infrastructure/logger/logger.go` (按日期和大小)
 - [ ] T019 [US1] 实现日志保留策略配置 `backend/internal/infrastructure/logger/logger.go` (rotateBackupExpire)
 - [ ] T020 [US1] 在 main.go 中集成日志初始化 `backend/main.go` (调用 logger.Initialize)
 - [ ] T021 [US1] 创建配置文件模板 `backend/configs/config.yaml` (logger 配置段)
 - [ ] T022 [P] [US1] 创建多环境配置文件 `backend/configs/config.dev.yaml`, `config.test.yaml`, `config.prod.yaml`
 - [ ] T023 [US1] 添加配置错误详细诊断信息 `backend/internal/infrastructure/logger/config.go`
 - [ ] T024 [US1] 编写 logger 包 README 文档 `backend/internal/infrastructure/logger/README.md`
 - [X] T015 [US1] 实现日志实例管理器 `backend/internal/infrastructure/logger/logger.go` (Initialize, GetInstance 函数)
 - [X] T016 [US1] 实现多日志实例初始化 `backend/internal/infrastructure/logger/logger.go` (initializeInstances 函数)
  - [X] T017 [US1] 集成 GoFrame glog 配置,包括异步日志写入配置 `backend/internal/infrastructure/logger/logger.go` (configureGLog 函数)
  - [X] T018 [US1] 实现日志分割策略配置 `backend/internal/infrastructure/logger/logger.go` (按日期和大小)
  - [X] T019 [US1] 实现日志保留策略配置 `backend/internal/infrastructure/logger/logger.go` (rotateBackupExpire)
  - [X] T020 [US1] 在 main.go 中集成日志初始化 `backend/main.go` (调用 logger.Initialize)
  - [X] T021 [US1] 创建配置文件模板 `backend/configs/config.yaml` (logger 配置段)
  - [X] T022 [P] [US1] 创建多环境配置文件 `backend/configs/config.dev.yaml`, `config.test.yaml`, `config.prod.yaml`
  - [X] T023 [US1] 添加配置错误详细诊断信息 `backend/internal/infrastructure/logger/config.go`
  - [X] T024 [US1] 编写 logger 包 README 文档 `backend/internal/infrastructure/logger/README.md`

**Checkpoint**: 用户故事1完成 - 日志配置管理功能完全可用且可独立测试

---

## Phase 4: User Story 2 - HTTP 请求全链路追踪 (Priority: P2)

**Goal**: 通过 TraceID 追踪 HTTP 请求的完整生命周期,支持从请求头提取或自动生成 TraceID,在所有日志中自动包含

**Independent Test**: 发送 HTTP 请求,在日志文件中搜索 TraceID,验证能找到该请求从进入到返回的完整日志链路

### Tests for User Story 2 (MANDATORY) ⚠️

 - [X] T025 [P] [US2] 编写 TraceID 生成单元测试 `backend/tests/unit/logger/traceid_test.go`
 - [X] T026 [P] [US2] 编写 TraceID 提取单元测试 `backend/tests/unit/logger/traceid_test.go`
 - [X] T027 [P] [US2] 编写 TraceID 传递中断恢复测试 `backend/tests/unit/logger/traceid_test.go`
 - [X] T028 [P] [US2] 编写访问日志中间件单元测试 `backend/tests/unit/middleware/access_log_test.go`
 - [X] T029 [P] [US2] 编写 HTTP 请求链路追踪集成测试 `backend/tests/integration/middleware/trace_test.go`

### Implementation for User Story 2

- [ ] T030 [P] [US2] 实现 TraceID 生成函数 `backend/internal/infrastructure/logger/traceid.go` (GenerateTraceID)
- [ ] T031 [P] [US2] 实现 TraceID 提取函数 `backend/internal/infrastructure/logger/traceid.go` (ExtractTraceID)
- [ ] T032 [US2] 实现 TraceID 传递中断检测和恢复 `backend/internal/infrastructure/logger/traceid.go` (EnsureTraceID)
- [ ] T033 [US2] 实现访问日志中间件 `backend/internal/app/http_server/middleware/access_log.go` (AccessLog 函数)
- [ ] T034 [US2] 实现 TraceID 注入到 Context `backend/internal/app/http_server/middleware/access_log.go`
- [ ] T035 [US2] 实现请求开始和结束日志记录 `backend/internal/app/http_server/middleware/access_log.go`
- [ ] T036 [US2] 实现 Panic 恢复中间件 `backend/internal/app/http_server/middleware/panic_recovery.go`
- [ ] T037 [US2] 实现 Panic 堆栈记录到错误日志 `backend/internal/app/http_server/middleware/panic_recovery.go`
- [ ] T038 [US2] 在 main.go 中注册中间件 `backend/main.go` (s.Use(middleware.AccessLog))
- [ ] T039 [US2] 配置 ctxKeys 自动提取 TraceID `backend/configs/config.yaml` (ctxKeys: ["TraceId", "UserId"])

 - [X] T030 [P] [US2] 实现 TraceID 生成函数 `backend/internal/infrastructure/logger/traceid.go` (GenerateTraceID)
 - [X] T031 [P] [US2] 实现 TraceID 提取函数 `backend/internal/infrastructure/logger/traceid.go` (ExtractTraceID)
 - [X] T032 [US2] 实现 TraceID 传递中断检测和恢复 `backend/internal/infrastructure/logger/traceid.go` (EnsureTraceID)
 - [X] T033 [US2] 实现访问日志中间件 `backend/internal/app/http_server/middleware/access_log.go` (AccessLog 函数)
 - [X] T034 [US2] 实现 TraceID 注入到 Context `backend/internal/app/http_server/middleware/access_log.go`
 - [X] T035 [US2] 实现请求开始和结束日志记录 `backend/internal/app/http_server/middleware/access_log.go`
 - [X] T036 [US2] 实现 Panic 恢复中间件 `backend/internal/app/http_server/middleware/panic_recovery.go`
 - [X] T037 [US2] 实现 Panic 堆栈记录到错误日志 `backend/internal/app/http_server/middleware/panic_recovery.go`
 - [X] T038 [US2] 在 main.go 中注册中间件 `backend/main.go` (s.Use(middleware.AccessLog))
 - [X] T039 [US2] 配置 ctxKeys 自动提取 TraceID `backend/configs/config.yaml` (ctxKeys: ["TraceId", "UserId"])

**Checkpoint**: 用户故事2完成 - HTTP 请求全链路追踪功能完全可用且可独立测试

---

## Phase 5: User Story 3 - 关键操作日志埋点 (Priority: P3)

**Goal**: 在关键业务操作点记录结构化日志,包含操作类型、参数、结果、耗时等信息,支持业务监控和性能优化

**Independent Test**: 执行特定业务操作(如访问学习内容),检查日志文件,验证记录了操作的详细信息和性能指标

### Tests for User Story 3 (MANDATORY) ⚠️

- [ ] T040 [P] [US3] 编写日志辅助方法单元测试 `backend/tests/unit/logger/helper_test.go` (LogInfo, LogError, LogSlow, LogBiz)
- [ ] T041 [P] [US3] 编写数据库日志 Handler 单元测试 `backend/tests/unit/middleware/db_log_test.go`
- [ ] T042 [P] [US3] 编写慢查询检测测试 `backend/tests/unit/middleware/db_log_test.go`
- [ ] T043 [P] [US3] 编写业务操作日志埋点集成测试 `backend/tests/integration/logger/business_log_test.go`

 - [X] T040 [P] [US3] 编写日志辅助方法单元测试 `backend/tests/unit/logger/helper_test.go` (LogInfo, LogError, LogSlow, LogBiz)
 - [X] T041 [P] [US3] 编写数据库日志 Handler 单元测试 `backend/tests/unit/middleware/db_log_test.go`
 - [X] T042 [P] [US3] 编写慢查询检测测试 `backend/tests/unit/middleware/db_log_test.go`
 - [X] T043 [P] [US3] 编写业务操作日志埋点集成测试 `backend/tests/integration/logger/business_log_test.go`

### Implementation for User Story 3

- [ ] T044 [P] [US3] 实现 LogInfo 辅助方法 `backend/internal/infrastructure/logger/helper.go`
- [ ] T045 [P] [US3] 实现 LogError 辅助方法 `backend/internal/infrastructure/logger/helper.go` (自动记录堆栈)
- [ ] T046 [P] [US3] 实现 LogSlow 辅助方法 `backend/internal/infrastructure/logger/helper.go`
- [ ] T047 [P] [US3] 实现 LogBiz 辅助方法 `backend/internal/infrastructure/logger/helper.go`
- [ ] T048 [US3] 实现数据库日志 Handler `backend/internal/app/http_server/middleware/db_log.go` (DBLogHandler)
- [ ] T049 [US3] 实现慢查询检测逻辑 `backend/internal/app/http_server/middleware/db_log.go` (threshold 配置)
- [ ] T050 [US3] 实现 SQL 执行日志记录 `backend/internal/app/http_server/middleware/db_log.go`
- [ ] T051 [US3] 在数据库初始化中注册 Handler `backend/internal/infrastructure/database/database.go`
- [ ] T052 [US3] 在现有业务代码中添加日志埋点:
  - `backend/internal/app/lexical_elements/*.go` (学习内容加载)
  - `backend/internal/app/constants/*.go` (常量模块内容加载)
  - `backend/internal/app/http_server/handler/*.go` (菜单导航和请求处理)
  - `backend/internal/app/http_server/middleware/error_handler.go` (错误处理,如存在)
  - 或在相应的 service 层添加业务操作日志
- [ ] T053 [US3] 配置慢查询阈值 `backend/configs/config.yaml` (database.slow.threshold: 1000)

 - [X] T044 [P] [US3] 实现 LogInfo 辅助方法 `backend/internal/infrastructure/logger/helper.go`
 - [X] T045 [P] [US3] 实现 LogError 辅助方法 `backend/internal/infrastructure/logger/helper.go` (自动记录堆栈)
 - [X] T046 [P] [US3] 实现 LogSlow 辅助方法 `backend/internal/infrastructure/logger/helper.go`
 - [X] T047 [P] [US3] 实现 LogBiz 辅助方法 `backend/internal/infrastructure/logger/helper.go`
 - [X] T048 [US3] 实现数据库日志 Handler `backend/internal/app/http_server/middleware/db_log.go` (DBLogHandler)
 - [X] T049 [US3] 实现慢查询检测逻辑 `backend/internal/app/http_server/middleware/db_log.go` (threshold 配置)
 - [X] T050 [US3] 实现 SQL 执行日志记录 `backend/internal/app/http_server/middleware/db_log.go`
 - [X] T051 [US3] 在数据库初始化中注册 Handler `backend/internal/infrastructure/database/database.go`
 - [X] T052 [US3] 在现有业务代码中添加日志埋点:
   - `backend/internal/app/lexical_elements/*.go` (学习内容加载)
   - `backend/internal/app/constants/*.go` (常量模块内容加载)
   - `backend/internal/app/http_server/handler/*.go` (菜单导航和请求处理)
   - `backend/internal/app/http_server/middleware/error_handler.go` (错误处理,如存在)
   - 或在相应的 service 层添加业务操作日志
 - [X] T053 [US3] 配置慢查询阈值 `backend/configs/config.yaml` (database.slow.threshold: 1000)

**Checkpoint**: 用户故事3完成 - 关键操作日志埋点功能完全可用且可独立测试

---

## Phase 6: User Story 4 - 日志查询与分析支持 (Priority: P4)

**Goal**: 提供基础的日志查询能力,支持按时间范围、日志级别、TraceID、关键字过滤日志

**Independent Test**: 生成测试日志,使用日志查询功能按不同条件(时间、级别、TraceID)查询,验证返回结果准确性

**Note**: 本期仅提供基础文本搜索能力,不实现复杂查询语法

### Tests for User Story 4 (MANDATORY) ⚠️

- [ ] T054 [P] [US4] 编写日志文件读取测试 `backend/tests/unit/logger/query_test.go`
- [ ] T055 [P] [US4] 编写 TraceID 查询测试 `backend/tests/unit/logger/query_test.go`
- [ ] T056 [P] [US4] 编写时间范围查询测试 `backend/tests/unit/logger/query_test.go`
- [ ] T057 [P] [US4] 编写日志级别过滤测试 `backend/tests/unit/logger/query_test.go`

 - [X] T054 [P] [US4] 编写日志文件读取测试 `backend/tests/unit/logger/query_test.go`
 - [X] T055 [P] [US4] 编写 TraceID 查询测试 `backend/tests/unit/logger/query_test.go`
 - [X] T056 [P] [US4] 编写时间范围查询测试 `backend/tests/unit/logger/query_test.go`
 - [X] T057 [P] [US4] 编写日志级别过滤测试 `backend/tests/unit/logger/query_test.go`

### Implementation for User Story 4

- [ ] T058 [P] [US4] 实现日志文件读取函数 `backend/internal/infrastructure/logger/query.go` (ReadLogFile)
- [ ] T059 [P] [US4] 实现 TraceID 查询函数 `backend/internal/infrastructure/logger/query.go` (QueryByTraceID)
- [ ] T060 [P] [US4] 实现时间范围查询函数 `backend/internal/infrastructure/logger/query.go` (QueryByTimeRange)
- [ ] T061 [P] [US4] 实现日志级别过滤函数 `backend/internal/infrastructure/logger/query.go` (QueryByLevel)
- [ ] T062 [US4] 实现关键字搜索函数 `backend/internal/infrastructure/logger/query.go` (QueryByKeyword)
- [ ] T063 [US4] 优化大文件查询性能 `backend/internal/infrastructure/logger/query.go` (流式读取)
- [ ] T064 [US4] 编写日志查询使用文档 `backend/internal/infrastructure/logger/README.md` (查询部分)

 - [X] T058 [P] [US4] 实现日志文件读取函数 `backend/internal/infrastructure/logger/query.go` (ReadLogFile)
 - [X] T059 [P] [US4] 实现 TraceID 查询函数 `backend/internal/infrastructure/logger/query.go` (QueryByTraceID)
 - [X] T060 [P] [US4] 实现时间范围查询函数 `backend/internal/infrastructure/logger/query.go` (QueryByTimeRange)
 - [X] T061 [P] [US4] 实现日志级别过滤函数 `backend/internal/infrastructure/logger/query.go` (QueryByLevel)
 - [X] T062 [US4] 实现关键字搜索函数 `backend/internal/infrastructure/logger/query.go` (QueryByKeyword)
 - [X] T063 [US4] 优化大文件查询性能 `backend/internal/infrastructure/logger/query.go` (流式读取)
 - [X] T064 [US4] 编写日志查询使用文档 `backend/internal/infrastructure/logger/README.md` (查询部分)

**Checkpoint**: 用户故事4完成 - 日志查询与分析功能完全可用且可独立测试

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: 跨用户故事的改进和完善

- [ ] T065 [P] 验证所有代码注释和用户文档使用中文 (Principle V/XV)
- [ ] T066 [P] 运行 `go fmt` 格式化所有代码
- [ ] T067 [P] 运行 `go vet` 静态分析
- [ ] T068 [P] 运行 `golint` 代码风格检查
- [ ] T069 [P] 运行 `go mod tidy` 清理依赖
- [ ] T070 验证单元测试覆盖率 ≥80% `go test -cover ./backend/internal/infrastructure/logger/...`
- [ ] T071 验证集成测试覆盖率 ≥80% `go test -cover ./backend/tests/integration/...`
- [X] T065 [P] 验证所有代码注释和用户文档使用中文 (Principle V/XV)
- [X] T066 [P] 运行 `go fmt` 格式化所有代码
- [X] T067 [P] 运行 `go vet` 静态分析
- [X] T068 [P] 运行 `golint` 代码风格检查
- [X] T069 [P] 运行 `go mod tidy` 清理依赖
- [X] T070 验证单元测试覆盖率 ≥80% `go test -cover ./backend/internal/infrastructure/logger/...`
- [X] T071 验证集成测试覆盖率 ≥80% `go test -cover ./backend/tests/integration/...`
- [ ] T072 [P] 性能压测 (1000 并发请求,验证日志开销 <10%)
- [ ] T073 [P] 验证 TraceID 查询性能 (<30 秒)
- [ ] T074 [P] 验证日志查询性能 (<5 秒,1GB 文件)
- [ ] T075 [P] 敏感信息脱敏检查 (密码、Token 等)
- [ ] T076 运行 quickstart.md 验证流程 `backend/docs/quickstart_validation.md`
- [ ] T077 更新项目 README.md (添加日志系统说明、配置指南、使用示例)
- [ ] T078 [P] 创建日志最佳实践文档 `backend/docs/logging_best_practices.md`
- [ ] T079 [P] 创建故障排查文档 `backend/docs/logging_troubleshooting.md`
- [ ] T080 生成最终 Git commit 消息 (遵循 Conventional Commits 规范)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖 - 可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成 - **阻塞所有用户故事**
- **User Stories (Phase 3-6)**: 所有依赖 Foundational 阶段完成
  - 用户故事可以并行进行 (如果有足够人力)
  - 或按优先级顺序执行 (P1 → P2 → P3 → P4)
- **Polish (Phase 7)**: 依赖所有期望的用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: Foundational 完成后可开始 - 无其他故事依赖
- **User Story 2 (P2)**: Foundational 完成后可开始 - 依赖 US1 的日志实例管理
- **User Story 3 (P3)**: Foundational 完成后可开始 - 依赖 US1 的日志实例管理和 US2 的 TraceID
- **User Story 4 (P4)**: Foundational 完成后可开始 - 依赖 US1 的日志文件生成

### Within Each User Story

- 测试必须先编写并失败,然后再实现
- 配置/数据结构 → 核心逻辑 → 集成 → 文档
- 故事完成后再进入下一优先级

### Parallel Opportunities

- Phase 1 所有 [P] 任务可并行
- Phase 2 所有 [P] 任务可并行
- Foundational 完成后,所有用户故事可并行开始 (如果团队容量允许)
- 每个用户故事内的 [P] 任务可并行
- 不同用户故事可由不同团队成员并行工作

---

## Parallel Example: User Story 1

```bash
# 并行启动用户故事1的所有测试:
Task: "编写日志实例初始化单元测试 backend/tests/unit/logger/logger_test.go"
Task: "编写日志级别过滤测试 backend/tests/unit/logger/logger_test.go"
Task: "编写日志格式化测试 backend/tests/unit/logger/logger_test.go"
Task: "编写多实例独立配置测试 backend/tests/unit/logger/logger_test.go"
Task: "编写配置文件缺失/无效启动失败测试 backend/tests/integration/logger/config_validation_test.go"

# 并行启动用户故事1的配置文件创建:
Task: "创建多环境配置文件 backend/configs/config.dev.yaml"
Task: "创建多环境配置文件 backend/configs/config.test.yaml"
Task: "创建多环境配置文件 backend/configs/config.prod.yaml"
```

---

## Parallel Example: User Story 2

```bash
# 并行启动用户故事2的所有测试:
Task: "编写 TraceID 生成单元测试 backend/tests/unit/logger/traceid_test.go"
Task: "编写 TraceID 提取单元测试 backend/tests/unit/logger/traceid_test.go"
Task: "编写 TraceID 传递中断恢复测试 backend/tests/unit/logger/traceid_test.go"
Task: "编写访问日志中间件单元测试 backend/tests/unit/middleware/access_log_test.go"
Task: "编写 HTTP 请求链路追踪集成测试 backend/tests/integration/middleware/trace_test.go"

# 并行启动用户故事2的核心实现:
Task: "实现 TraceID 生成函数 backend/internal/infrastructure/logger/traceid.go"
Task: "实现 TraceID 提取函数 backend/internal/infrastructure/logger/traceid.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 1: Setup
2. 完成 Phase 2: Foundational (关键 - 阻塞所有故事)
3. 完成 Phase 3: User Story 1
4. **停止并验证**: 独立测试用户故事1
5. 如果就绪,部署/演示

### Incremental Delivery

1. 完成 Setup + Foundational → 基础就绪
2. 添加 User Story 1 → 独立测试 → 部署/演示 (MVP!)
3. 添加 User Story 2 → 独立测试 → 部署/演示
4. 添加 User Story 3 → 独立测试 → 部署/演示
5. 添加 User Story 4 → 独立测试 → 部署/演示
6. 每个故事增加价值而不破坏之前的故事

### Parallel Team Strategy

多个开发者:

1. 团队一起完成 Setup + Foundational
2. Foundational 完成后:
   - 开发者 A: User Story 1 (配置管理)
   - 开发者 B: User Story 2 (TraceID 追踪)
   - 开发者 C: User Story 3 (日志埋点)
   - 开发者 D: User Story 4 (日志查询)
3. 故事独立完成和集成

---

## Notes

- [P] 任务 = 不同文件,无依赖
- [Story] 标签将任务映射到特定用户故事以便追踪
- 每个用户故事应该可以独立完成和测试
- 实现前验证测试失败
- 每个任务或逻辑组后提交
- 在任何检查点停止以独立验证故事
- 避免: 模糊任务、同文件冲突、破坏独立性的跨故事依赖

---

## Task Summary

- **Total Tasks**: 80
- **Setup Phase**: 4 tasks
- **Foundational Phase**: 5 tasks
- **User Story 1 (P1)**: 15 tasks (5 tests + 10 implementation)
- **User Story 2 (P2)**: 15 tasks (5 tests + 10 implementation)
- **User Story 3 (P3)**: 14 tasks (4 tests + 10 implementation)
- **User Story 4 (P4)**: 11 tasks (4 tests + 7 implementation)
- **Polish Phase**: 16 tasks

**Parallel Opportunities**: 45 tasks marked [P] can run in parallel within their phase

**Independent Test Criteria**:
- US1: 修改配置文件,启动应用,验证日志按配置输出
- US2: 发送 HTTP 请求,搜索 TraceID,验证完整日志链路
- US3: 执行业务操作,检查日志,验证操作详情和性能指标
- US4: 生成测试日志,查询过滤,验证结果准确性

**Suggested MVP Scope**: User Story 1 only (统一日志配置管理)

**Estimated Effort**: 5-7 工作日 (包含编码、测试、文档、Code Review)

---

## Phase 8: Import Cycle Resolution

**Purpose**: 解决后端包之间的导入循环问题，确保所有测试能够正常运行

- [ ] T073 解决 database ↔ middleware 导入循环问题
  - 创建独立的 `backend/internal/infrastructure/db_logging` 包
  - 将数据库日志处理逻辑从 middleware 移动到 db_logging 包
  - 更新 database/sqlite.go 使用新的 db_logging 包
  - 修复配置文件中的重复 logger 配置
  - 验证所有后端测试通过
