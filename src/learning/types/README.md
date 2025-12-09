# Go Types 章节

版本：基于 Go 1.24.5，CLI/HTTP 双模式共用同一内容源与测验。

## 章节结构（计划）
- `types.go`：类型主题注册器与通用数据结构。
- `overview.go`：类型概览与提纲。
- `boolean.go` / `numeric.go` / `string_type.go`：基础类型内容。
- `array.go` / `slice.go` / `struct_type.go` / `pointer.go` / `function_type.go` / `map_type.go` / `channel_type.go`：复合类型内容。
- `interface_*.go`：接口与类型集相关内容。
- `search.go`：关键词检索索引。
- `quiz.go`：测验题库与评分逻辑。
- `cli/menu.go`：CLI 菜单与交互。
- `http/handlers.go`：HTTP 路由处理。

## 协同约定
- CLI 与 HTTP 共用统一的数据结构与注册表，避免内容漂移。
- 所有注释与输出保持中文，便于教学与回归测试。
- 后续填充内容后需保持 `go test ./...` 通过且覆盖率≥80%。

