# Implementation Plan: Go 类型章节学习方案

**Branch**: `007-go-types-study` | **Date**: 2025-12-09 | **Spec**: `D:\studyspace\go-study\go-study2\specs\007-go-types-study\spec.md`  
**Input**: Feature specification from `D:\studyspace\go-study\go-study2\specs\007-go-types-study\spec.md`

## Summary

- 目标：交付 Go 类型章节学习内容，覆盖基础/复合/接口/类型集规则与非法递归、可比较性、通道方向等边界，并提供测验、检索与打印友好提纲，满足 FR-001~FR-007。
- 技术思路：基于 Go 1.24.5 与现有 GoFrame HTTP/CLI 双模式，按章节分包到 `src/learning/types`，统一内容源（CLI/HTTP/打印索引），表驱动测验与搜索索引，测试覆盖率维持 80%+。

## Technical Context

**Language/Version**: Go 1.24.5  
**Primary Dependencies**: 标准库（fmt/strings/encoding/json/testing），GoFrame v2.9.5（沿用现有 HTTP 服务器）  
**Storage**: N/A（内容、索引与测验静态内置，必要时以内存结构持有）  
**Testing**: `go test`（单元+契约+集成）；表驱动测验校验；CLI/HTTP 输出一致性断言  
**Target Platform**: 本地 CLI 与 HTTP 服务（ghttp 8080，终端输入输出与 HTML/JSON 双格式）  
**Project Type**: 单体 CLI/HTTP 学习工具  
**Performance Goals**: 内容/测验/检索在本地 <1s 返回；搜索响应满足 SC-003（<=15s）；示例运行输出即刻可得  
**Constraints**: 全中文文档与注释；章节分包与子文件对应子主题；CLI/HTTP 共用内容源与题库；避免深层嵌套与额外三方依赖  
**Scale/Scope**: 单章节类型学习（基础+7类复合+接口+类型集）；≥1 组测验覆盖身份/可比较性/接口实现；索引覆盖关键词与正反例；进度记录沿用轻量数据结构

## Constitution Check

- **Principle I (Simplicity)**: 仅用标准库+既有 GoFrame，示例与索引为静态数据，PASS
- **Principle II (Comments)**: 计划在章节入口、内容生成与测试处添加前置中文分层注释，PASS
- **Principle III (Language)**: 所有文档、注释、测验解析与输出保持中文，PASS
- **Principle IV (Nesting)**: 用卫语句与小函数拆分内容/测验渲染，避免深嵌套，PASS
- **Principle V (YAGNI)**: 不引入新框架/存储，按当前章节最小集实现，PASS
- **Principle VI (Testing)**: 单元+契约+集成计划覆盖内容/题库/索引一致性，目标 >=80%，PASS
- **Principle XVII/XVIII/XX**: 分层菜单、CLI+HTTP 双模式与章节分包/子文件命名（snake_case）均在结构设计中，PASS

## Project Structure

### Documentation (this feature)

```text
D:\studyspace\go-study\go-study2\specs\007-go-types-study\
├── plan.md              # 本文件
├── research.md          # Phase 0 输出
├── data-model.md        # Phase 1 输出
├── quickstart.md        # Phase 1 输出
├── contracts\
│   └── types-learning.md # Phase 1 输出（CLI/HTTP 契约与API概要）
├── checklists\
│   └── requirements.md
└── tasks.md             # Phase 2 (/speckit.tasks) 生成
```

### Source Code (repository root)

```text
src/
└── learning/
    └── types/
        ├── types.go              # 章节入口与内容聚合
        ├── overview.go           # 类型定义/身份/可比较性总览
        ├── boolean.go            # 布尔类型
        ├── numeric.go            # 整数/浮点/复数与别名 byte/rune
        ├── string_type.go        # 字符串类型与索引规则
        ├── array.go              # 数组与非法递归约束
        ├── slice.go              # 切片与容量/共享底层数组
        ├── struct_type.go        # 结构体与嵌入/递归限制/标签
        ├── pointer.go            # 指针类型
        ├── function_type.go      # 函数签名与可变参数
        ├── interface_basic.go    # 基础接口定义与方法集
        ├── interface_embedded.go # 嵌入接口与方法集合并
        ├── interface_general.go  # 类型集/union/~T 规则与限制
        ├── interface_impl.go     # 接口实现判定与类型集子集关系
        ├── map_type.go           # map 键限制与删除/clear 规则
        ├── channel_type.go       # 通道方向与缓冲
        ├── README.md             # 章节说明
        ├── cli/
        │   └── menu.go           # CLI 菜单与子主题调度
        └── http/
            └── handlers.go       # HTTP 内容/测验/检索输出

tests/
├── unit/
│   └── learning/
│       └── types/                # 内容、题库、索引的表驱动单测
├── integration/
│   ├── types_api_test.go         # HTTP 路由与格式协商
│   └── types_cli_test.go         # CLI 菜单与题库一致性
└── contract/
    └── learning/
        └── types/                # CLI/HTTP 输出契约（JSON/HTML/文本）
```

**Structure Decision**: 采用单体布局并在 `src/learning/types/` 按子主题拆分文件，CLI 与 HTTP 共用同一内容/题库数据源；更新 `main.go` 菜单、`internal/app/http_server/router.go` 与 handler，注册 `/api/v1/topic/types` 路由，测试按单元/集成/契约分层保证一致性与 80% 覆盖。

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| None | N/A | N/A |
