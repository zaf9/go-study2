# Implementation Plan: 后端目录重构与前端预留

**Branch**: `008-backend-structure-refactor` | **Date**: 2025-12-10 | **Spec**: `specs/008-backend-structure-refactor/spec.md`  
**Input**: Feature specification from `/specs/008-backend-structure-refactor/spec.md`

## Summary

将现有 go-study2 后端源码与配置迁移到仓库根目录的`backend/`下，保持根目录 README.md 为原始版本，仅更新其中涉及目录变更的描述；参考根目录示例中的其他目录（`.github`、`common`、`docker`、`docs`、`frontend`、`helm`、`idl`、`scripts`等）若当前仓库不存在且当前工作无关，则不创建。`backend/README.md` 将重写以聚焦后端架构与 API 说明，并确保迁移后构建与测试保持通过。

## Technical Context

**Language/Version**: Go（以现有 go.mod 版本为准）  
**Primary Dependencies**: 以当前 go.mod 列出的依赖为准，不新增依赖  
**Storage**: 保持现状（未指定新增存储）  
**Testing**: `go test ./...`（迁移后需通过），必要时更新路径相关用例  
**Target Platform**: 维持迁移前后端部署目标一致  
**Project Type**: 后端服务，预留前端空间  
**Performance Goals**: 迁移后构建与运行无性能回退  
**Constraints**: 遵循仓库现有构建入口（优先执行`./build.bat`，若不存在则按通用构建方式）；路径调整需兼容现有脚本  
**Scale/Scope**: 以当前项目规模为主，迁移不引入额外子项目

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Simplicity):** 是，保持目录扁平、职责清晰。  
- **Principle II (Comments):** 是，更新脚本与文档说明，层次清晰。  
- **Principle III (Language):** 是，文档与注释使用中文。  
- **Principle IV (Nesting):** 是，迁移方案保持浅逻辑。  
- **Principle V (YAGNI):** 是，不新增无关目录或依赖。  
- **Principle VI (Testing):** 是，迁移后`go test ./...`与既有构建全部通过，覆盖率维持>=80%。

## Project Structure

### Documentation (this feature)

```text
specs/008-backend-structure-refactor/
├── plan.md              # 本文件
├── research.md          # Phase 0 输出
├── data-model.md        # Phase 1 输出
├── quickstart.md        # Phase 1 输出
├── contracts/           # Phase 1 输出
└── tasks.md             # Phase 2 (/speckit.tasks 生成)
```

### Source Code (repository root)

```text
./                        # 根目录，保留原 go-study2 README.md
├── .github/              # 若已存在则保留
├── backend/              # 迁移后的后端代码主目录（go.mod、go.sum、main.go 等）
│   ├── api/
│   ├── application/
│   ├── conf/
│   ├── crossdomain/
│   ├── domain/
│   ├── infra/
│   ├── internal/
│   ├── pkg/
│   ├── script/
│   ├── types/
│   ├── Dockerfile
│   ├── build.sh
│   ├── go.mod
│   ├── go.sum
│   └── README.md         # 重写：后端架构与 API 说明
├── common/               # 仅在仓库已有且当前需要时保留，不新建
├── docker/               # 同上
├── docs/                 # 同上，更新目录结构说明
├── frontend/             # 预留/保留占位，不影响后端构建
├── helm/                 # 若已存在则保留
├── idl/                  # 若已存在则保留
├── scripts/              # 若已存在则保留
├── .gitattributes/.gitignore/.nvmrc 等配置文件（保持现状）
└── README.md             # 根README保留原 go-study2 内容，仅更新目录变更描述
```

**Structure Decision**: 采用“backend/ 前后端分层”结构；只迁移后端到`backend/`，根目录其他参考目录仅在已存在或当前需用时保留，不主动新建；预留`frontend/`占位但不新增前端实现。

## Complexity Tracking

无额外复杂度引入，未发现需豁免的宪章违规项。
