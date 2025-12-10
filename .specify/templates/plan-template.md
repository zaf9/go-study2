# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: [e.g., Python 3.11, Swift 5.9, Rust 1.75 or NEEDS CLARIFICATION]  
**Primary Dependencies**: [e.g., FastAPI, UIKit, LLVM or NEEDS CLARIFICATION]  
**Storage**: [if applicable, e.g., PostgreSQL, CoreData, files or N/A]  
**Testing**: [e.g., pytest, XCTest, cargo test or NEEDS CLARIFICATION]  
**Target Platform**: [e.g., Linux server, iOS 15+, WASM or NEEDS CLARIFICATION]
**Project Type**: [single/web/mobile - determines source structure]  
**Performance Goals**: [domain-specific, e.g., 1000 req/s, 10k lines/sec, 60 fps or NEEDS CLARIFICATION]  
**Constraints**: [domain-specific, e.g., <200ms p95, <100MB memory, offline-capable or NEEDS CLARIFICATION]  
**Scale/Scope**: [domain-specific, e.g., 10k users, 1M LOC, 50 screens or NEEDS CLARIFICATION]

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (代码质量与可维护性):** 方案是否清晰、可读、单一职责且易测试?
- **Principle II (显式错误处理):** 是否为所有错误设计明确处理,避免静默失败?
- **Principle III/XXI/XXXVI (全面测试):** 是否规划覆盖率≥80%,各包具备 *_test.go/示例,前端核心组件同样达标?
- **Principle IV (单一职责):** 目录/文件/函数职责是否单一且可拆分维护?
- **Principle V/XV (一致文档与中文要求):** 注释与用户文档是否清晰,后端内容是否全部使用中文?
- **Principle VI (YAGNI):** 是否拒绝当前不需要的复杂度或模式?
- **Principle VII (安全优先):** 输入校验、鉴权授权、HTTPS、敏感信息保护是否覆盖到位?
- **Principle VIII/XVIII (可预测结构):** 是否遵循标准 Go 布局,仅根目录定义 main, go.mod/go.sum 完整?
- **Principle IX (依赖纪律):** 外部依赖是否最小且必要?
- **Principle X (性能优化):** 是否考虑关键性能瓶颈(数据库查询/内存/前端渲染与包体积)?
- **Principle XI (文档同步):** 方案是否包含完成后更新根 README 及相关文档?
- **Principle XIV (清晰分层注释):** 是否规划各层职责的中文注释?
- **Principle XVI (浅层逻辑):** 是否避免深层嵌套,采用卫语句与函数拆分?
- **Principle XVII (一致开发者体验):** 初始化与工作流是否对初学者友好且一致?
- **Principle XIX (包级 README):** 是否为每个包规划 README 说明用途与用法?
- **Principle XX (代码质量执行):** 是否计划 go fmt/go vet/golint/go mod tidy 等质量检查?
- **Principle XXII (分层菜单导航):** 如含交互,是否支持多级菜单、统一编号与清晰返回/错误提示?
- **Principle XXIII (双学习模式):** 新章节是否同时支持 CLI 与 HTTP,共享内容源?
- **Principle XXIV (层次化章节结构):** 章节/子章节/子包与文件命名是否按规范组织?
- **Principle XXV (HTTP/CLI 一致性):** 菜单、路由、响应结构与 Topic 注册是否一致且显式错误处理?

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
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
