# Tasks: 后端目录重构与前端预留

**Input**: Design documents from `/specs/008-backend-structure-refactor/`
**Prerequisites**: plan.md, spec.md

**Tests**: 按宪章要求，需保持至少 80% 单测覆盖；每个用户故事均需可独立验证。

**组织方式**: 按用户故事分组，确保可独立实现与测试。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行（不同文件且无依赖）
- **[Story]**: 用户故事标签（US1, US2, US3）
- 描述中需包含明确文件路径

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 为迁移做最小准备，不引入新目录（除计划中的 backend/ 骨架）。

- [x] T001 检查根构建入口 `./build.bat` 是否存在，记录可用构建命令到 `specs/008-backend-structure-refactor/quickstart.md`（若文件后续生成）
- [x] T002 [P] 在 `backend/` 创建目录骨架（api/application/conf/crossdomain/domain/infra/internal/pkg/script/types），仅当不存在时创建

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 梳理需迁移的文件与路径引用，形成可执行清单。

- [x] T003 [P] 清点当前后端源码与配置清单，输出到 `specs/008-backend-structure-refactor/research.md`
- [x] T004 [P] 列出脚本与配置中的路径引用（如 `scripts/*`, `Dockerfile`, `build.sh`），记录到 `specs/008-backend-structure-refactor/research.md`

---

## Phase 3: User Story 1 - 后端代码集中化 (Priority: P1) 🎯 MVP

**Goal**: 后端源码与配置迁移到 `backend/` 并保持构建可用。  
**Independent Test**: 从仓库根执行 `./build.bat` 或等效命令、在 `backend/` 执行 `go test ./...` 均成功。

### Tests for User Story 1 (MANDATORY)

- [x] T005 [P] [US1] 迁移后在 `backend/` 运行 `go test ./...`，记录结果到 `specs/008-backend-structure-refactor/quickstart.md`
- [x] T006 [P] [US1] 迁移后从仓库根执行 `./build.bat`（若不存在则使用通用构建命令），记录结果到 `specs/008-backend-structure-refactor/quickstart.md`

### Implementation for User Story 1

- [x] T007 [US1] 将后端源码、配置、`Dockerfile`、`build.sh` 等迁移至 `backend/`，保持原有相对结构
- [x] T008 [US1] 更新后端脚本与配置中的路径引用（如 `backend/script/*`, `backend/Dockerfile` 等），确保指向新目录
- [x] T009 [US1] 在 `backend/go.mod` 校验模块路径与依赖，执行 `go mod tidy`
- [x] T010 [US1] 清理仓库根遗留的后端源码/配置重复文件，确认根仅保留 `backend/` 与预留的 `frontend/`

**Checkpoint**: User Story 1 可独立构建与测试。

---

## Phase 4: User Story 2 - 文档与路径一致 (Priority: P2)

**Goal**: 文档反映新目录结构，指引可直接复现构建运行。  
**Independent Test**: 按更新后的文档从根执行构建与运行，流程可复现。

### Tests for User Story 2 (MANDATORY)

- [x] T011 [P] [US2] 按更新后的根 `README.md` 指引执行一次后端构建/运行，记录验证到 `specs/008-backend-structure-refactor/quickstart.md`

### Implementation for User Story 2

- [x] T012 [US2] 重写 `backend/README.md`，补充后端架构与主要 API 说明
- [x] T013 [US2] 更新根 `README.md` 中的目录结构与后端路径描述，保持其他内容不变
- [x] T014 [P] [US2] 若 `docs/` 目录存在，更新其中的目录结构示意与后端路径说明

**Checkpoint**: 文档可指导新人完成后端构建与运行。

---

## Phase 5: User Story 3 - 预留前端空间 (Priority: P3)

**Goal**: 预留 `frontend/` 占位且不影响后端构建。  
**Independent Test**: 创建占位后，后端构建与测试仍全通过。

### Tests for User Story 3 (MANDATORY)

- [x] T015 [P] [US3] 创建/确认 `frontend/` 占位（必要时添加 `.gitkeep`），确保构建脚本不会误用
- [x] T016 [US3] 占位后在 `backend/` 运行 `go test ./...`，确认不受影响

### Implementation for User Story 3

- [x] T017 [US3] 检查构建/脚本（如 `./build.bat`、`backend/script/*`）对 `frontend/` 的潜在依赖，确保无耦合

**Checkpoint**: 预留前端空间且后端流程不受影响。

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: 收尾与一致性校验。

- [x] T018 [P] 扫描并修正残留的旧路径引用（`scripts/`, `docs/`, `Dockerfile` 等）  
- [x] T019 [P] 确认新增/更新文档均为中文且同步目录变更（根 `README.md`、`backend/README.md`、若有 `docs/`）  
- [x] T020 全量验证：仓库根执行 `./build.bat`（或等效命令）与 `backend/` 下 `go test ./...` 均通过

---

## Dependencies & Execution Order

- Setup 完成后进入 Foundational；Foundational 完成后方可开始各用户故事。  
- 用户故事按优先级：US1 → US2 → US3；若团队并行，US2/US3 需等待 US1 迁移完成。  
- Polish 在所有目标用户故事完成后执行。

### User Story Dependencies

- US1 无依赖。  
- US2 依赖 US1（文档需基于已迁移结构）。  
- US3 依赖 US1（占位需建立在迁移后结构）。

### Parallel Opportunities

- 标记 [P] 的任务可并行：T002、T003、T004、T005、T006、T011、T014、T015、T018、T019。  
- 不同用户故事可在 US1 完成后并行推进。

## Implementation Strategy

- MVP：完成 US1（迁移与构建通过）后即可获得可演示版本。  
- 增量交付：依次完成 US2（文档一致）、US3（前端占位），最后执行 Polish 全量验证。  
- 每个阶段后运行对应测试任务，确保独立可验证。

