# Research - 后端目录重构与前端预留

## 目录与文件清单

- 当前待迁移的后端文件/目录（根路径）：
  - `go.mod`, `go.sum`, `main.go`, `main_test.go`
  - 目录：`configs/`, `doc/`, `docs/`, `internal/`, `logs/`, `scripts/`, `src/`, `tests/`
- 计划目标：迁移至 `backend/`，保持相对结构不变，根目录仅保留顶层文档与规范目录（.specify、specs、.github 等）。

## 路径引用与脚本依赖

- 已发现需要更新的路径引用：
  - 根 `README.md`：运行命令、项目结构、配置路径均假定后端位于根；需改为进入 `backend/` 运行。
  - `scripts/check-go.ps1`：目标路径列表为 `internal`, `src`, `tests`, `main.go`（根）；迁移后需指向 `backend/*`。
- 待验证的路径引用：
  - 其他脚本/文档是否硬编码相对路径（迁移后需确认）。

## 发现与决策

- 迁移后所有构建与测试从 `backend/` 执行；根无 `build.bat`，采用 `go build ./...` / `go test ./...`。
- 若存在新增路径引用问题，统一修正为以 `backend/` 为工作目录的相对路径。

