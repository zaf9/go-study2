# Implementation Plan: Go-Study2 学习闭环与测验体系

**Branch**: `011-learning-progress-quiz` | **Date**: 2025-12-11 | **Spec**: `specs/011-learning-progress-quiz/spec.md`
**Input**: Feature specification from `/specs/011-learning-progress-quiz/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

实现学习闭环：标准化四主题章节内容（概述/要点/详细说明/示例/陷阱/实践建议）、追踪与展示学习进度（权重加权、状态与阅读恢复、防抖+指数退避同步）、提供章节测验（多题型抽题、判分与历史）。后端基于 Go 1.24.5 + GF + SQLite 复用现有架构；前端基于 Next.js 14 + SWR + Ant Design，按统一 API 合约交互。测试覆盖率前后端核心功能 ≥80%。

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.24.5（backend），TypeScript + Next.js 14（frontend）  
**Primary Dependencies**: backend: gogf/gf v2.9.5、go-sqlite3；frontend: React 18、AntD 5、SWR、Axios、React Markdown/Prism  
**Storage**: SQLite3（现有），使用 learning_progress / quiz_* 表  
**Testing**: backend `go test ./...` + table-driven +示例；frontend Jest + React Testing Library + Next lint；覆盖率核心 ≥80%  
**Target Platform**: Backend: Windows/Linux（当前 dev），Frontend: Next.js SPA/SSG 输出  
**Project Type**: web（前后端分离，REST 接口）  
**Performance Goals**: 进度/测验 API p95 < 300ms（单用户典型查询 < 10 条记录）；前端交互首屏 < 2s（本地静态资源）  
**Constraints**: 进度同步防抖 10s、指数退避最多 5 次并卸载前强制同步；数据库单机 SQLite，需事务确保一致性  
**Scale/Scope**: 4 主题约 48 章节，题库 400-500 题，单机并发中低（<100 并发），重点在一致性与体验

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (代码质量与可维护性):** 采用分层（handler/service/repo），函数职责单一，可测试。  
- **Principle II (显式错误处理):** 所有 API 返回统一 Response{code,message,data}，校验失败/DB 错误显式返回。  
- **Principle III/XXI/XXXVI (全面测试):** 计划 backend 单元+示例覆盖核心逻辑≥80%，frontend Jest/RTL 覆盖核心组件≥80%。  
- **Principle IV (单一职责):** 进度/测验/内容模块解耦，文件按领域拆分。  
- **Principle V/XV (中文要求):** 后端注释/文档全中文，前端文案一致。  
- **Principle VI (YAGNI):** 不引入新消息队列或分布式组件，保持单机 SQLite。  
- **Principle VII (安全优先):** 依赖现有鉴权中间件，新增接口做参数校验与权限检查，HTTPS 由网关配置。  
- **Principle VIII/XVIII (可预测结构):** 复用现有 backend 标准布局；不新增额外 main。  
- **Principle IX (依赖纪律):** 复用现有 GF/SWR/AntD，无新增重量依赖。  
- **Principle X (性能优化):** 进度上报防抖+指数退避，批量更新；前端按需渲染/分页。  
- **Principle XI (文档同步):** 完成后更新根 README 的功能/路由/运行说明。  
- **Principle XIV (清晰分层注释):** 为 service/repo/handler 增补中文职责注释。  
- **Principle XVI (浅层逻辑):** 使用早返回拆分校验与业务。  
- **Principle XVII (一致开发者体验):** 提供 quickstart（build.bat 优先），脚本化命令。  
- **Principle XIX (包级 README):** 新增包补充 README（progress、quiz）。  
- **Principle XX (代码质量执行):** 计划 gofmt/go vet/golint/go test/go mod tidy；frontend eslint/tsc/jest。  
- **Principle XXII-XXV (菜单/双模式/结构/一致性):** CLI/HTTP 共用内容与测验数据源，路由与菜单注册同步，编号/响应结构遵循现有规范。

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
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
specs/011-learning-progress-quiz/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
└── tasks.md

backend/
├── api/...
├── internal/
│   ├── app/            # 内容/进度/测验服务实现
│   ├── domain/         # 领域实体与仓储接口
│   ├── infra/          # 基础设施、数据库
│   └── interfaces/     # HTTP/CLI handler
├── configs/
└── tests/              # 单元/集成测试

frontend/
├── app/ or pages/      # Next.js 路由
├── components/         # 进度条、题目组件、答案解析
├── services/           # SWR + Axios 封装
└── tests/              # Jest + RTL
```

**Structure Decision**: 采用现有前后端分离布局，新增进度/测验模块遵循 backend 分层与 frontend 组件化结构。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| (None) | - | - |

## Phase 0: Research (complete)
- 产物：`research.md`（指数退避策略、统一响应格式、SQLite+时间戳防回退、性能与容量目标）。  
- 未决项：无。

## Phase 1: Design & Contracts (complete)
- 数据模型：`data-model.md` 定义 LearningProgress / QuizQuestion / QuizSession / QuizAttempt 及约束、索引、状态规则。  
- 合约：`contracts/progress.md`、`contracts/quiz.md`，覆盖进度与测验核心 API。  
- Quickstart：`quickstart.md`，含 build.bat 先行要求、前后端命令、测试目标。  
- Agent 上下文：本计划已纳入宪章与现有依赖，无需新增技术。

## Phase 2: Implementation Prep
- 目标：在 `tasks.md` 细化开发任务（/speckit.tasks）。  
- 重点落地：进度服务（防抖+指数退避+状态机）、测验抽题与判分、前端进度/测验页面与组件、测试与 README 更新。
