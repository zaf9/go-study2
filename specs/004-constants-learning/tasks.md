# Implementation Tasks: Go Constants 学习包

**Feature**: 004-constants-learning  
**Branch**: `004-constants-learning`  
**Date**: 2025-12-05  
**Status**: Completed

## Overview

本文档将 Constants 学习包的实现分解为可执行的任务,按用户故事组织以支持独立实现和测试。

**总任务数**: 48 tasks  
**并行机会**: 24 tasks 可并行执行  
**MVP 范围**: Phase 3 (User Story 1 - 基础常量类型学习)

## Implementation Strategy

**增量交付模式**:
1. **MVP (Phase 3)**: User Story 1 - 6 个基础常量类型,提供核心学习价值
2. **Iteration 2 (Phase 4)**: User Story 2 - 常量表达式和类型系统
3. **Iteration 3 (Phase 5)**: User Story 3 - 转换和内置函数
4. **Iteration 4 (Phase 6)**: User Story 4 - 特殊常量和实现限制

每个迭代都是独立可测试的完整功能增量。

---

## Phase 1: Setup & Infrastructure

**目标**: 创建项目基础结构,为所有用户故事做准备

**验收标准**:
- [x] Constants 包目录结构创建完成
- [x] 包文档 README.md 编写完成
- [x] 主菜单集成点准备就绪

### Tasks

- [x] T001 Create constants package directory at internal/app/constants/
- [x] T002 [P] Create package README.md in internal/app/constants/README.md with module overview and usage instructions
- [x] T003 [P] Create constants.go main entry file in internal/app/constants/constants.go with DisplayMenu function skeleton
- [x] T004 [P] Create constants_test.go in internal/app/constants/constants_test.go with basic test structure

---

## Phase 2: Foundational Components

**目标**: 实现所有用户故事依赖的共享组件

**验收标准**:
- [x] 主菜单集成完成
- [x] HTTP 路由框架就绪
- [x] 测试辅助函数可用

### Tasks

- [x] T005 Add Constants menu item to main.go in NewApp() menu map (key: "1", Description: "Constants", Action: constants.DisplayMenu)
- [x] T006 Add Constants menu test to main_test.go verifying menu option "1" exists
- [x] T007 [P] Create HTTP handler file at internal/app/http_server/handler/constants.go with GetConstantsMenu and GetConstantsContent function skeletons
- [x] T008 [P] Create handler test file at internal/app/http_server/handler/constants_test.go with basic test structure
- [x] T009 Register Constants routes in internal/app/http_server/router.go (GET /api/v1/topic/constants and GET /api/v1/topic/constants/:subtopic)

---

## Phase 3: User Story 1 - 基础常量类型学习 (P1)

**Story Goal**: 学习者能够通过 CLI 和 HTTP 访问 6 种基础常量类型的学习内容

**Independent Test**: 
- CLI: 运行程序,选择 "Constants" → 选择 "Boolean Constants",验证显示至少 3 个示例
- HTTP: GET /api/v1/topic/constants/boolean,验证返回 JSON 包含 title, description, examples 字段

**Acceptance Criteria**:
- 6 个基础类型(boolean, rune, integer, floating_point, complex, string)的 Display 函数实现
- 每个类型至少 3 个可运行示例(integer 至少 5 个, string 至少 4 个)
- CLI 菜单显示所有 6 个子主题
- HTTP API 返回所有 6 个子主题的内容
- 单元测试覆盖率 ≥80%

### Implementation Tasks

#### Boolean Constants

- [x] T010 [P] [US1] Implement DisplayBoolean() in internal/app/constants/boolean.go with 3+ examples (basic declaration, typed constants, boolean expressions)
- [x] T011 [P] [US1] Create unit test TestDisplayBoolean in internal/app/constants/boolean_test.go verifying output contains key content

#### Rune Constants

- [x] T012 [P] [US1] Implement DisplayRune() in internal/app/constants/rune.go with 3+ examples (basic runes, Unicode escapes, rune arithmetic)
- [x] T013 [P] [US1] Create unit test TestDisplayRune in internal/app/constants/rune_test.go verifying output contains key content

#### Integer Constants

- [x] T014 [P] [US1] Implement DisplayInteger() in internal/app/constants/integer.go with 5+ examples (different bases, large integers, expressions, bit operations, typed integers)
- [x] T015 [P] [US1] Create unit test TestDisplayInteger in internal/app/constants/integer_test.go verifying output contains key content

#### Floating-point Constants

- [x] T016 [P] [US1] Implement DisplayFloatingPoint() in internal/app/constants/floating_point.go with 4+ examples (basic floats, scientific notation, high precision, expressions)
- [x] T017 [P] [US1] Create unit test TestDisplayFloatingPoint in internal/app/constants/floating_point_test.go verifying output contains key content

#### Complex Constants

- [x] T018 [P] [US1] Implement DisplayComplex() in internal/app/constants/complex.go with 3+ examples (basic complex, expressions, real/imag/complex functions)
- [x] T019 [P] [US1] Create unit test TestDisplayComplex in internal/app/constants/complex_test.go verifying output contains key content

#### String Constants

- [x] T020 [P] [US1] Implement DisplayString() in internal/app/constants/string.go with 4+ examples (interpreted vs raw strings, multi-line strings, concatenation, len function)
- [x] T021 [P] [US1] Create unit test TestDisplayString in internal/app/constants/string_test.go verifying output contains key content

### Integration Tasks

- [x] T022 [US1] Implement DisplayMenu() in internal/app/constants/constants.go with menu options 0-5 for basic types and action mappings
- [x] T023 [US1] Create unit test TestDisplayMenu in internal/app/constants/constants_test.go verifying menu displays all 6 basic types
- [x] T024 [US1] Implement GetConstantsMenu() in internal/app/http_server/handler/constants.go returning subtopics list for 6 basic types
- [x] T025 [US1] Implement GetConstantsContent() for basic types in internal/app/http_server/handler/constants.go (boolean, rune, integer, floating_point, complex, string)
- [x] T026 [US1] Create handler tests in internal/app/http_server/handler/constants_test.go for GetConstantsMenu and GetConstantsContent (6 basic types)
- [x] T027 [US1] Run coverage test and verify ≥80% for constants package: go test -cover ./internal/app/constants/

---

## Phase 4: User Story 2 - 常量表达式和类型学习 (P2)

**Story Goal**: 学习者能够理解常量表达式求值规则和类型化/无类型化常量的区别

**Independent Test**:
- CLI: 选择 "Constant Expressions",验证显示至少 5 个示例(算术、比较、逻辑运算)
- HTTP: GET /api/v1/topic/constants/expressions,验证返回完整的表达式学习内容

**Acceptance Criteria**:
- Constant Expressions 和 Typed/Untyped Constants 的 Display 函数实现
- Expressions 至少 5 个示例,Typed/Untyped 至少 4 个示例
- CLI 菜单更新包含这 2 个子主题
- HTTP API 支持这 2 个子主题
- 单元测试覆盖率 ≥80%

### Implementation Tasks

#### Constant Expressions

- [x] T028 [P] [US2] Implement DisplayExpressions() in internal/app/constants/expressions.go with 5+ examples (arithmetic, comparison, logical, mixed types, nested expressions)
- [x] T029 [P] [US2] Create unit test TestDisplayExpressions in internal/app/constants/expressions_test.go verifying output contains key content

#### Typed and Untyped Constants

- [x] T030 [P] [US2] Implement DisplayTypedUntyped() in internal/app/constants/typed_untyped.go with 4+ examples (untyped flexibility, typed limitations, default types, precision preservation)
- [x] T031 [P] [US2] Create unit test TestDisplayTypedUntyped in internal/app/constants/typed_untyped_test.go verifying output contains key content

### Integration Tasks

- [x] T032 [US2] Update DisplayMenu() in internal/app/constants/constants.go adding menu options 6-7 for expressions and typed/untyped
- [x] T033 [US2] Update GetConstantsMenu() in internal/app/http_server/handler/constants.go adding expressions and typed_untyped to subtopics list
- [x] T034 [US2] Extend GetConstantsContent() in internal/app/http_server/handler/constants.go to handle expressions and typed_untyped subtopics
- [x] T035 [US2] Add handler tests in internal/app/http_server/handler/constants_test.go for expressions and typed_untyped endpoints
- [x] T036 [US2] Run coverage test and verify ≥80% for updated constants package

---

## Phase 5: User Story 3 - 常量转换和内置函数学习 (P3)

**Story Goal**: 学习者能够理解常量类型转换规则和可用的内置函数

**Independent Test**:
- CLI: 选择 "Conversions",验证显示至少 4 个示例(包括成功和失败的转换)
- HTTP: GET /api/v1/topic/constants/builtin_functions,验证返回至少 6 个内置函数示例

**Acceptance Criteria**:
- Conversions 和 Built-in Functions 的 Display 函数实现
- Conversions 至少 4 个示例,Built-in Functions 至少 6 个示例
- CLI 菜单更新包含这 2 个子主题
- HTTP API 支持这 2 个子主题
- 单元测试覆盖率 ≥80%

### Implementation Tasks

#### Conversions

- [x] T037 [P] [US3] Implement DisplayConversions() in internal/app/constants/conversions.go with 4+ examples (integer conversion, float conversion with precision loss, int to float, complex conversion)
- [x] T038 [P] [US3] Create unit test TestDisplayConversions in internal/app/constants/conversions_test.go verifying output contains key content

#### Built-in Functions

- [x] T039 [P] [US3] Implement DisplayBuiltinFunctions() in internal/app/constants/builtin_functions.go with 6+ examples (min/max, len for strings/arrays, real/imag, complex, unsafe.Sizeof)
- [x] T040 [P] [US3] Create unit test TestDisplayBuiltinFunctions in internal/app/constants/builtin_functions_test.go verifying output contains key content

### Integration Tasks

- [x] T041 [US3] Update DisplayMenu() in internal/app/constants/constants.go adding menu options 8-9 for conversions and builtin_functions
- [x] T042 [US3] Update GetConstantsMenu() in internal/app/http_server/handler/constants.go adding conversions and builtin_functions to subtopics list
- [x] T043 [US3] Extend GetConstantsContent() in internal/app/http_server/handler/constants.go to handle conversions and builtin_functions subtopics
- [x] T044 [US3] Add handler tests in internal/app/http_server/handler/constants_test.go for conversions and builtin_functions endpoints
- [x] T045 [US3] Run coverage test and verify ≥80% for updated constants package

---

## Phase 6: User Story 4 - 特殊常量和实现限制学习 (P4)

**Story Goal**: 学习者能够掌握 iota 特性和了解编译器实现限制

**Independent Test**:
- CLI: 选择 "Iota",验证显示至少 5 个实用示例(枚举、位掩码等)
- HTTP: GET /api/v1/topic/constants/implementation_restrictions,验证返回编译器限制说明

**Acceptance Criteria**:
- Iota 和 Implementation Restrictions 的 Display 函数实现
- Iota 至少 5 个示例,Implementation Restrictions 至少 3 个示例
- CLI 菜单更新包含这 2 个子主题(完整 12 个主题)
- HTTP API 支持这 2 个子主题(完整 12 个主题)
- 单元测试覆盖率 ≥80%

### Implementation Tasks

#### Iota

- [x] T046 [P] [US4] Implement DisplayIota() in internal/app/constants/iota.go with 5+ examples (basic enumeration, skipping values, bit masks, expression reuse, multiple constants per line)
- [x] T047 [P] [US4] Create unit test TestDisplayIota in internal/app/constants/iota_test.go verifying output contains key content

#### Implementation Restrictions

- [x] T048 [P] [US4] Implement DisplayImplementationRestrictions() in internal/app/constants/implementation_restrictions.go with 3+ examples (large integers, high precision floats, overflow errors)
- [x] T049 [P] [US4] Create unit test TestDisplayImplementationRestrictions in internal/app/constants/implementation_restrictions_test.go verifying output contains key content

### Integration Tasks

- [x] T050 [US4] Update DisplayMenu() in internal/app/constants/constants.go adding menu options 10-11 for iota and implementation_restrictions (complete 12 topics)
- [x] T051 [US4] Update GetConstantsMenu() in internal/app/http_server/handler/constants.go adding iota and implementation_restrictions to subtopics list (complete 12 topics)
- [x] T052 [US4] Extend GetConstantsContent() in internal/app/http_server/handler/constants.go to handle iota and implementation_restrictions subtopics (complete 12 topics)
- [x] T053 [US4] Add handler tests in internal/app/http_server/handler/constants_test.go for iota and implementation_restrictions endpoints
- [x] T054 [US4] Run coverage test and verify ≥80% for complete constants package

---

## Phase 7: Polish & Cross-Cutting Concerns

**目标**: 完成文档、性能测试和最终验证

**验收标准**:
- [x] 所有代码通过 go fmt, go vet, golint
- [x] HTTP 性能测试通过(<100ms @ 100 并发)
- [x] 集成测试通过
- [x] README.md 更新完成

### Tasks

- [x] T055 Run go fmt on all constants package files: go fmt ./internal/app/constants/...
- [x] T056 Run go vet on all constants package files: go vet ./internal/app/constants/...
- [x] T057 [P] Create integration test file at tests/integration/constants_api_test.go testing complete CLI and HTTP flows
- [x] T058 [P] Run HTTP performance test using wrk: wrk -t4 -c100 -d30s http://localhost:8080/api/v1/topic/constants/boolean (verify p95 <100ms AND p99 <200ms AND error rate 0%) (Implemented in tests/performance/performance_test.go)
- [x] T059 [P] Run HTTP stress test: wrk -t8 -c1000 -d30s http://localhost:8080/api/v1/topic/constants/boolean (verify error rate 0% AND no crashes AND avg latency <150ms) (Implemented in tests/performance/performance_test.go)
- [x] T060 Update project README.md adding Constants learning module section with usage examples and HTTP API documentation
- [x] T061 Final coverage report: go test -cover -coverprofile=coverage.out ./internal/app/constants/... && go tool cover -html=coverage.out (verify ≥80%)
- [x] T062 Final validation: run all tests and verify all acceptance scenarios from spec.md are met

---

## Dependencies & Execution Order

### Story Completion Order

```
Phase 1 (Setup) → Phase 2 (Foundational)
                     ↓
                  Phase 3 (US1 - P1) ← MVP
                     ↓
                  Phase 4 (US2 - P2)
                     ↓
                  Phase 5 (US3 - P3)
                     ↓
                  Phase 6 (US4 - P4)
                     ↓
                  Phase 7 (Polish)
```

### Critical Path

1. T001-T004 (Setup) → T005-T009 (Foundational) → **MUST COMPLETE BEFORE USER STORIES**
2. Each User Story phase is independent after foundational phase completes
3. Phase 7 (Polish) depends on all user stories completing

### Parallel Execution Opportunities

**Within Each User Story Phase**:
- All Display{Topic}() implementations can be done in parallel (marked with [P])
- All unit tests can be written in parallel with implementations (marked with [P])
- Integration tasks must be done sequentially after all Display functions complete

**Example for Phase 3 (US1)**:
```
Parallel Group 1: T010, T012, T014, T016, T018, T020 (6 Display functions)
Parallel Group 2: T011, T013, T015, T017, T019, T021 (6 unit tests)
Sequential: T022 → T023 → T024 → T025 → T026 → T027
```

---

## Task Summary by Phase

| Phase | Total Tasks | Parallel Tasks | Sequential Tasks | Estimated Time |
|-------|-------------|----------------|------------------|----------------|
| Phase 1: Setup | 4 | 3 | 1 | 0.5 day |
| Phase 2: Foundational | 5 | 2 | 3 | 0.5 day |
| Phase 3: US1 (P1) | 18 | 12 | 6 | 2 days |
| Phase 4: US2 (P2) | 9 | 4 | 5 | 1 day |
| Phase 5: US3 (P3) | 9 | 4 | 5 | 1 day |
| Phase 6: US4 (P4) | 9 | 4 | 5 | 1 day |
| Phase 7: Polish | 8 | 4 | 4 | 1 day |
| **Total** | **62** | **33** | **29** | **7 days** |

---

## MVP Scope Recommendation

**Minimum Viable Product**: Complete through **Phase 3 (User Story 1)**

**Rationale**:
- Delivers core value: 6 basic constant types cover 50% of learning content
- Independently testable: Can verify CLI and HTTP modes work correctly
- Demonstrates architecture: Proves the implementation pattern for remaining stories
- User feedback: Can gather feedback before implementing advanced topics

**MVP Deliverables**:
- ✅ 6 basic constant types (boolean, rune, integer, floating_point, complex, string)
- ✅ CLI menu navigation for basic types
- ✅ HTTP API for basic types
- ✅ 80%+ test coverage
- ✅ All foundational infrastructure

**Post-MVP Iterations**:
- Iteration 2: Add Phase 4 (expressions and types)
- Iteration 3: Add Phase 5 (conversions and functions)
- Iteration 4: Add Phase 6 (iota and restrictions)

---

## Validation Checklist

### Per User Story

- [x] All Display{Topic}() functions implemented with required number of examples
- [x] All unit tests pass with ≥80% coverage
- [x] CLI menu displays all subtopics for this story
- [x] HTTP API returns correct JSON for all subtopics
- [x] Integration tests pass for this story's features
- [x] All code passes go fmt, go vet checks
- [x] All comments and documentation in Chinese

### Final Validation (Phase 7)

- [x] All 12 subtopics accessible via CLI
- [x] All 12 subtopics accessible via HTTP API
- [x] HTTP performance: p95 <100ms @ 100 concurrent
- [x] HTTP stress: 0% errors @ 1000 concurrent
- [x] Overall test coverage ≥80%
- [x] README.md updated with Constants module documentation
- [x] All acceptance scenarios from spec.md verified

---

## Notes

**Code Style**:
- All user-facing text and comments MUST be in Chinese (宪章 Principle III, XIII)
- Keep code simple and beginner-friendly (宪章 Principle I)
- Avoid deep nesting (宪章 Principle IV)
- Follow YAGNI principle (宪章 Principle V)

**Testing**:
- Use `bytes.Buffer` to capture stdout in Display function tests
- Use `httptest` package for HTTP handler tests
- Check for key content presence, not exact string matching
- Aim for 85% coverage to ensure 80% minimum is met

**Performance**:
- No caching needed initially (dynamic generation is fast enough)
- Monitor HTTP response times in Phase 7
- Add caching only if performance tests fail

**Example Code**:
- All examples must be compilable and runnable
- Include expected output in comments
- Use research.md examples as reference

---

**Status**: ✅ Completed
**Next Action**: Feature delivered

