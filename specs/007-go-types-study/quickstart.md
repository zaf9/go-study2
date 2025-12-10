# Quickstart: Go 类型章节学习方案

1) 构建  
- 在仓库根执行：`cd backend && go test ./... && go build ./...`（如存在 `./build.bat` 可优先执行）。

2) 运行 CLI 学习模式  
- `cd backend && go run main.go`，主菜单选择 `3` 进入 Types（布尔/数值/字符串/数组/切片/结构体/指针/函数/接口/map/chan）。  
- 支持命令：输入编号查看内容与测验；`o` 打印提纲；`quiz` 进行综合测验；`search <keyword>` 关键词检索；`q` 返回。

3) 运行 HTTP 学习模式  
- `cd backend && go run main.go -d`（默认 127.0.0.1:8080）。  
- 菜单：`curl "http://localhost:8080/api/v1/topic/types?format=html"`  
- 子主题：`curl "http://localhost:8080/api/v1/topic/types/array"`  
- 提纲：`curl "http://localhost:8080/api/v1/topic/types/outline?format=html"`  
- 测验提交：`curl -X POST "http://localhost:8080/api/v1/topic/types/quiz/submit" -H "Content-Type: application/json" -d "{\"answers\":[{\"id\":\"q-all-1\",\"choice\":\"A\"}]}"`  
- 检索：`curl "http://localhost:8080/api/v1/topic/types/search?keyword=map%20key"`（JSON/HTML 均可）。

4) 搜索关键词示例  
- `map key`、`~int`、`interface nil`、`slice share`、`array length`

5) 测试  
- `cd backend && go test ./src/learning/types/... ./tests/...`  
- 目标覆盖率 80%+；重点校验 CLI/HTTP 内容一致、测验评分与搜索索引正确。

6) 文档同步  
- 如有变更同步 `backend/src/learning/types/README.md`、根级 `README.md` 与 `backend/docs/quickstart-variables.md`（新增 Types 入口）。