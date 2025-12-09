# Quickstart: Go 类型章节学习方案

1) 构建  
- 优先执行 `./build.bat`；若不存在则运行 `go test ./...` 与 `go build ./...`。

2) 运行 CLI 学习模式  
- `go run main.go`，主菜单选择 `3` 进入 Types（布尔/数值/字符串/数组/切片/结构体/指针/函数/接口/map/chan）。  
- 子菜单输入编号查看内容与示例，测验结束后显示得分与解析，可按提示重做。

3) 运行 HTTP 学习模式  
- `go run main.go -d`（默认 127.0.0.1:8080）。  
- `curl "http://localhost:8080/api/v1/topic/types?format=html"` 查看菜单。  
- `curl "http://localhost:8080/api/v1/topic/types/array"` 获取子主题内容与测验。  
- `curl -X POST "http://localhost:8080/api/v1/topic/types/quiz/submit" -d "{\"answers\":[{\"id\":\"q1\",\"choice\":\"A\"}]}"` 体验评分与解析。  
- `curl "http://localhost:8080/api/v1/topic/types/search?keyword=map"` 演示规则检索（JSON/HTML 均可）。

4) 测试  
- `go test ./src/learning/types/... ./tests/...`  
- 目标覆盖率 80%+；重点校验 CLI/HTTP 内容一致、测验评分与搜索索引正确。

5) 文档同步  
- 如有变更同步 `src/learning/types/README.md`、根级 `README.md` 与 `docs/quickstart-variables.md`（新增 Types 入口）。