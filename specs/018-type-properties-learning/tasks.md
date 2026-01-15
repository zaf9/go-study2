# Tasks: Go 类型属性章节学习方案

**Feature**: 018-type-properties-learning  
**Created**: 2026-01-15  
**Status**: Ready for Implementation

## Overview

实现 Go Properties of types and values 章节的学习内容，包含7个子主题，支持 CLI 和 HTTP 双模式。

## Phase 1: 基础设施搭建

### T001: 创建 properties 包基础结构
**Priority**: P1  
**Dependencies**: None  
**Effort**: S

创建 `backend/src/learning/properties/` 目录及核心文件：
- `properties.go`: 核心数据结构和 API
- `content.go`: 内容注册入口
- `README.md`: 包文档

**Acceptance Criteria**:
- [ ] 目录结构已创建
- [ ] 核心数据类型已定义（Topic, PropertyConcept, PropertyRule, ExampleCase, QuizItem 等）
- [ ] 核心函数已实现（AllTopics, RegisterContent, LoadContent, LoadQuiz, EvaluateQuiz, SearchReferences）
- [ ] 核心数据类型支持 Principle XXXIII 状态转换字段（not_started/in_progress/completed）
- [ ] 包可编译通过

---

### T002: 创建 overview.go
**Priority**: P1  
**Dependencies**: T001  
**Effort**: S

实现章节概览功能：
- `GetOverview()` 函数
- `RenderPrintableOutline()` 函数

**Acceptance Criteria**:
- [ ] 概览函数返回所有子主题摘要
- [ ] 打印提纲函数返回格式化文本

---

## Phase 2: 子主题内容实现

### T003: 实现 representation.go (值的表示)
**Priority**: P1  
**Dependencies**: T001  
**Effort**: M

实现"值的表示"子主题：
- 自包含值 vs 引用值的规则
- 示例代码（数组、切片、map、指针等）
- 测验题目（2-3道）
- 搜索索引

**Acceptance Criteria**:
- [ ] 内容注册成功
- [ ] 包含至少2个代码示例
- [ ] 包含至少2道测验题
- [ ] 关键词可搜索

---

### T004: 实现 underlying_type.go (底层类型)
**Priority**: P1  
**Dependencies**: T001  
**Effort**: M

实现"底层类型"子主题：
- 底层类型推导规则
- 类型别名、类型定义、类型参数的示例
- 测验题目

**Acceptance Criteria**:
- [ ] 覆盖所有底层类型推导规则
- [ ] 包含正反例对比
- [ ] 测验覆盖边界情况

---

### T005: 实现 core_type.go (核心类型)
**Priority**: P1  
**Dependencies**: T001  
**Effort**: M

实现"核心类型"子主题：
- 非接口类型的核心类型
- 接口类型的核心类型判定条件
- bytestring 核心类型
- 示例和测验

**Acceptance Criteria**:
- [ ] 覆盖所有核心类型判定规则
- [ ] 包含有核心类型和无核心类型的接口示例
- [ ] 包含 bytestring 示例

---

### T006: 实现 type_identity.go (类型标识)
**Priority**: P1  
**Dependencies**: T001  
**Effort**: L

实现"类型标识"子主题：
- 命名类型 vs 类型字面量
- 各种类型的结构等价规则
- 泛型实例化类型的标识
- 相同/不同类型的示例对比
- 测验题目

**Acceptance Criteria**:
- [ ] 覆盖所有类型标识判定规则
- [ ] 包含多种类型的对比示例
- [ ] 测验包含边界情况

---

### T007: 实现 assignability.go (可赋值性)
**Priority**: P1  
**Dependencies**: T001  
**Effort**: L

实现"可赋值性"子主题：
- 6个标准可赋值条件
- 3个类型参数额外条件
- 各条件的示例代码
- 测验题目

**Acceptance Criteria**:
- [ ] 覆盖全部9个可赋值性条件
- [ ] 每个条件有对应示例
- [ ] 测验覆盖常见误区

---

### T008: 实现 representability.go (可表示性)
**Priority**: P1  
**Dependencies**: T001  
**Effort**: M

实现"可表示性"子主题：
- 常量可表示性规则
- 各类型的可表示性判定
- 边界情况（溢出、精度）
- 测验题目

**Acceptance Criteria**:
- [ ] 覆盖所有可表示性规则
- [ ] 包含边界情况示例
- [ ] 测验包含常见陷阱

---

### T009: 实现 method_set.go (方法集)
**Priority**: P1  
**Dependencies**: T001  
**Effort**: M

实现"方法集"子主题：
- 定义类型方法集
- 指针类型方法集
- 接口类型方法集
- 嵌入字段规则
- 测验题目

**Acceptance Criteria**:
- [ ] 覆盖所有方法集规则
- [ ] 包含值接收者/指针接收者对比
- [ ] 包含接口实现判定示例

---

### T010: 实现 quiz.go (综合测验)
**Priority**: P2  
**Dependencies**: T003-T009  
**Effort**: S

实现综合测验功能：
- 从各子主题抽取代表性题目
- `LoadComprehensiveQuiz()` 函数
- `EvaluateComprehensiveQuiz()` 函数

**Acceptance Criteria**:
- [ ] 综合测验包含5-7道题
- [ ] 覆盖多个子主题
- [ ] 评分功能正常

---

### T011: 实现 search.go (搜索索引)
**Priority**: P2  
**Dependencies**: T003-T009  
**Effort**: S

实现搜索索引注册：
- 注册各子主题的关键词
- 中英文关键词支持

**Acceptance Criteria**:
- [ ] 所有子主题可通过关键词搜索
- [ ] 支持模糊匹配

---

## Phase 3: CLI 集成

### T012: 实现 cli/menu.go
**Priority**: P1  
**Dependencies**: T001-T009  
**Effort**: M

实现 CLI 菜单：
- `DisplayMenu()` 函数
- 子主题内容展示
- 测验交互
- 检索功能
- 打印提纲

**Acceptance Criteria**:
- [ ] 菜单导航正常
- [ ] 内容展示清晰
- [ ] 测验交互完整
- [ ] 错误处理友好

---

### T013: 集成到主菜单
**Priority**: P1  
**Dependencies**: T012  
**Effort**: S

将 Properties 章节添加到主菜单：
- 修改 `main.go` 或主菜单文件
- 添加菜单项

**Acceptance Criteria**:
- [ ] 主菜单可见新选项
- [ ] 可正常进入子菜单

---

## Phase 4: HTTP 集成

### T014: 实现 http/handlers.go
**Priority**: P2  
**Dependencies**: T001-T011  
**Effort**: M

实现 HTTP 处理器：
- 菜单处理器
- 内容处理器
- 测验处理器
- 搜索处理器

**Acceptance Criteria**:
- [ ] 所有接口返回正确格式
- [ ] 错误响应规范
- [ ] 支持 JSON/HTML 格式

---

### T015: 注册 HTTP 路由
**Priority**: P2  
**Dependencies**: T014  
**Effort**: S

在路由器中注册 Properties 路由：
- `/api/v1/topic/properties`
- `/api/v1/topic/properties/:topic`
- `/api/v1/topic/properties/:topic/quiz`
- `/api/v1/topic/properties/search`

**Acceptance Criteria**:
- [ ] 所有路由可访问
- [ ] 路由注册符合规范

---

## Phase 5: 测试

### T016: 单元测试
**Priority**: P1  
**Dependencies**: T001-T011  
**Effort**: M

创建单元测试：
- 内容注册测试
- 加载功能测试
- 测验评分测试
- 搜索功能测试

**Acceptance Criteria**:
- [ ] 覆盖率 >= 80%
- [ ] 所有用例通过

---

### T017: 集成测试
**Priority**: P2  
**Dependencies**: T012-T015  
**Effort**: M

创建集成测试：
- CLI 菜单测试
- HTTP 接口测试

**Acceptance Criteria**:
- [ ] CLI 流程测试通过
- [ ] HTTP 接口测试通过

---

## Phase 6: 文档

### T018: 创建 README.md
**Priority**: P2  
**Dependencies**: T001-T015  
**Effort**: S

创建包文档：
- 功能说明
- 使用示例
- API 说明

**Acceptance Criteria**:
- [ ] 文档清晰完整
- [ ] 示例可运行

---

## Summary

| Phase | Tasks | Priority | Total Effort |
|-------|-------|----------|--------------|
| Phase 1 | T001-T002 | P1 | 2S |
| Phase 2 | T003-T011 | P1-P2 | 7M + 2S |
| Phase 3 | T012-T013 | P1 | 1M + 1S |
| Phase 4 | T014-T015 | P2 | 1M + 1S |
| Phase 5 | T016-T017 | P1-P2 | 2M |
| Phase 6 | T018 | P2 | 1S |

**Critical Path**: T001 → T003-T009 → T012 → T013 → T016
