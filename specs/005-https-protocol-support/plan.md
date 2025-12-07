# Implementation Plan: HTTPS 协议支持

**Branch**: `005-https-protocol-support` | **Date**: 2025-12-08 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-https-protocol-support/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

为现有 HTTP 服务添加可配置的 HTTPS 支持。通过 `https.enabled` 配置项切换协议模式，支持 TLS 1.2+ 自签名证书，提供证书路径配置和友好的错误提示。

## Technical Context

**Language/Version**: Go 1.24.5  
**Primary Dependencies**: github.com/gogf/gf/v2 (GoFrame 框架)  
**Storage**: N/A（仅配置文件）  
**Testing**: go test（标准库）  
**Target Platform**: Linux/Windows 服务器  
**Project Type**: single（单一 Go 项目）  
**Performance Goals**: N/A（与现有 HTTP 模式保持一致）  
**Constraints**: TLS 1.2+ 最低版本要求  
**Scale/Scope**: 单服务器实例，单协议模式运行

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Simplicity):** ✅ 通过 - 方案使用 GoFrame 内置的 TLS 支持，无需引入额外复杂性
- **Principle II (Comments):** ✅ 通过 - 计划在配置结构和服务器启动逻辑中添加清晰的中文注释
- **Principle III (Language):** ✅ 通过 - 所有文档和注释将使用中文
- **Principle IV (Nesting):** ✅ 通过 - 协议选择逻辑使用简单的条件判断，避免深层嵌套
- **Principle V (YAGNI):** ✅ 通过 - 仅实现当前需求（单协议切换），不添加双栈或证书自动生成等额外功能
- **Principle VI (Testing):** ✅ 通过 - 将为配置验证和服务器启动逻辑编写单元测试，目标覆盖率 ≥80%
- **Principle XVIII (Dual Mode):** ✅ 通过 - HTTPS 支持将同时适用于 CLI 和 HTTP 学习模式

## Project Structure

### Documentation (this feature)

```text
specs/005-https-protocol-support/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# 现有项目结构（单一 Go 项目）
configs/
├── config.yaml          # 迁移：从根目录迁移，添加 https 配置节
└── certs/               # 新增：证书目录
    ├── server.crt       # HTTPS 证书文件
    └── server.key       # HTTPS 私钥文件

internal/
├── config/
│   ├── config.go        # 修改：添加 HttpsConfig 结构，更新配置加载路径
│   └── config_test.go   # 修改：添加 HTTPS 配置验证测试
├── app/
│   └── http_server/
│       ├── server.go    # 修改：添加 HTTPS 启动逻辑
│       └── server_test.go # 修改：添加 HTTPS 服务器测试

tests/
├── integration/
│   └── https_mode_test.go  # 新增：HTTPS 模式集成测试
```

**Structure Decision**: 使用现有单一项目结构，配置文件迁移到 `configs/` 目录，证书存放在 `configs/certs/`，在 `internal/config` 中扩展配置，在 `internal/app/http_server` 中添加 HTTPS 支持

## Complexity Tracking

> Constitution Check 已全部通过，无需记录违规项。

## Constitution Re-Check (Post Phase 1 Design)

*设计完成后的重新检查*

- **Principle I (Simplicity):** ✅ 通过 - 设计使用 GoFrame 原生 `EnableHTTPS()` 方法，代码简洁
- **Principle II (Comments):** ✅ 通过 - 数据模型已定义清晰的中文字段说明
- **Principle III (Language):** ✅ 通过 - 所有文档（research.md, data-model.md, quickstart.md）均使用中文
- **Principle IV (Nesting):** ✅ 通过 - 协议选择逻辑为单层 if-else，无嵌套
- **Principle V (YAGNI):** ✅ 通过 - 仅实现配置切换，无额外功能
- **Principle VI (Testing):** ✅ 通过 - 计划包含配置验证和服务器启动的单元测试
- **Principle VII (Single Responsibility):** ✅ 通过 - HttpsConfig 单独负责 HTTPS 配置
- **Principle X (Error Handling):** ✅ 通过 - 定义了完整的错误消息（证书缺失、文件不存在等）
- **Principle XVIII (Dual Mode):** ✅ 通过 - HTTPS 透明支持现有 HTTP 学习模式端点

**结论**: 设计通过所有 Constitution 原则检查，可进入任务分解阶段。
