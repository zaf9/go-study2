# Implementation Plan: Go 类型属性章节学习方案

**Branch**: `018-type-properties-learning` | **Date**: 2026-01-15 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `D:\studyspace\go-study\go-study2\specs\018-type-properties-learning\spec.md`

## Summary

- 目标：交付 Go Properties of types and values 章节学习内容，覆盖值的表示、底层类型、核心类型、类型标识、可赋值性、可表示性、方法集等7个子主题，提供规则解析、示例代码、测验与检索功能，满足 FR-001~FR-012。
- 技术思路：基于 Go 1.24.5 与现有 GoFrame HTTP/CLI 双模式，参照 Types 章节实现模式，在 `backend/src/learning/` 下创建 `properties` 包，按子主题分文件组织，统一内容源（CLI/HTTP/打印索引），表驱动测验与搜索索引，测试覆盖率维持 80%+。

## Technical Context

**Language/Version**: Go 1.24.5  
**Primary Dependencies**: 标准库（fmt/strings/encoding/json/testing），GoFrame v2.9.5（沿用现有 HTTP 服务器）  
**Storage**: N/A（内容、索引与测验静态内置，必要时以内存结构持有）  
**Testing**: `go test`（单元+契约+集成）；表驱动测验校验；CLI/HTTP 输出一致性断言  
**Target Platform**: 本地 CLI 与 HTTP 服务（ghttp 8080，终端输入输出与 HTML/JSON 双格式）  
**Project Type**: 单体 CLI/HTTP 学习工具  
**Performance Goals**: 内容/测验/检索在本地 <1s 返回；搜索响应满足 SC-003（<=15s）；示例运行输出即刻可得  
**Constraints**: 全中文文档与注释；章节分包与子文件对应子主题；CLI/HTTP 共用内容源与题库；避免深层嵌套与额外三方依赖  
**Scale/Scope**: 单章节类型属性学习（7个子主题）；测验覆盖各子主题规则；索引覆盖关键词与正反例；进度记录沿用轻量数据结构

## Constitution Check

- **Principle I (Simplicity)**: 仅用标准库+既有 GoFrame，示例与索引为静态数据，PASS
- **Principle II (Comments)**: 计划在章节入口、内容生成与测试处添加前置中文分层注释，PASS
- **Principle III (Language)**: 所有文档、注释、测验解析与输出保持中文，PASS
- **Principle IV (Nesting)**: 用卫语句与小函数拆分内容/测验渲染，避免深嵌套，PASS
- **Principle V (YAGNI)**: 不引入新框架/存储，按当前章节最小集实现，PASS
- **Principle VI (Testing)**: 单元+契约+集成计划覆盖内容/题库/索引一致性，目标 >=80%，PASS
- **Principle XVII/XVIII/XX**: 分层菜单、CLI+HTTP 双模式与章节分包/子文件命名（snake_case）均在结构设计中，PASS
- **Principle XXXII (Quiz Standards)**: 本次作为 MVP 实现，采用代码内置少量题目方式，暂不实现全量 YAML 题库，留待后续专项升级 (See Complexity Tracking)，DEVIATION JUSTIFIED

## Project Structure

### Documentation (this feature)

```text
D:\studyspace\go-study\go-study2\specs\018-type-properties-learning\
├── plan.md              # 本文件
├── research.md          # Phase 0 输出
├── data-model.md        # Phase 1 输出
├── quickstart.md        # Phase 1 输出
├── contracts\
│   └── properties-learning.md # Phase 1 输出（CLI/HTTP 契约与API概要）
├── checklists\
│   └── requirements.md
└── tasks.md             # Phase 2 (/speckit.tasks) 生成
```

### Source Code (repository root)

```text
backend/
└── src/
    └── learning/
        └── properties/
            ├── properties.go          # 章节入口与内容聚合（核心数据结构）
            ├── content.go             # 内容注册入口
            ├── overview.go            # 章节概览与提纲聚合
            ├── representation.go      # 值的表示（自包含值与引用值）
            ├── underlying_type.go     # 底层类型推导规则
            ├── core_type.go           # 核心类型判定（含bytestring）
            ├── type_identity.go       # 类型标识/相同性判定
            ├── assignability.go       # 可赋值性规则（9个条件）
            ├── representability.go    # 可表示性规则（常量表示）
            ├── method_set.go          # 方法集规则
            ├── quiz.go                # 综合测验数据
            ├── search.go              # 搜索索引注册
            ├── README.md              # 章节说明
            ├── cli/
            │   └── menu.go            # CLI 菜单与子主题调度
            └── http/
                └── handlers.go        # HTTP 内容/测验/检索输出

tests/
├── unit/
│   └── learning/
│       └── properties/               # 内容、题库、索引的表驱动单测
├── integration/
│   └── learning/
│       └── properties/               # HTTP 路由与格式协商、CLI 菜单一致性
└── contract/
    └── learning/
        └── properties/               # CLI/HTTP 输出契约（JSON/HTML/文本）
```

**Structure Decision**: 采用与 Types 章节相同的布局模式，在 `backend/src/learning/properties/` 按子主题拆分文件，CLI 与 HTTP 共用同一内容/题库数据源；需更新主菜单与路由注册。

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| Principle XXXII (Quiz YAML) | 快速交付核心学习内容 (MVP) | 全量题库和 YAML 解析器开发成本较高，当前优先保证规则覆盖 |
