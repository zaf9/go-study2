# 后端架构与运行说明

## 目录结构

- `go.mod` / `go.sum`：模块与依赖管理  
- `main.go` / `main_test.go`：入口与顶层测试  
- `configs/`：`config.yaml`、证书等配置  
- `internal/`：配置加载、HTTP 服务、学习内容（lexical_elements、constants 等）  
- `src/`：CLI/HTTP 复用的章节内容与测验逻辑  
- `tests/`：unit / integration / contract 测试  
- `scripts/`：`check-go.ps1`（gofmt、go vet、golint 统一入口）  
- `doc/`、`docs/`：文档材料与章节 quickstart  

## 构建与运行

```bash
# 在仓库根执行
cd backend

# 构建/运行
go build ./...         # 或 go run main.go
go run main.go -d      # 启动 HTTP 模式（默认 127.0.0.1:8080）

# 测试
go test ./...
```

## 配置

- 默认配置位于 `configs/config.yaml`，HTTPS 证书放置于 `configs/certs/`。  
- 配置加载通过 `internal/config`，工作目录以 `backend/` 为基准。  

## 默认管理员

- 初始账号：`admin` / `GoStudy@123`。  
- 首次登录会被强制修改密码；改密后旧口令与旧令牌全部失效。

## API 速览

- 主题列表：`GET /api/v1/topics?format=json|html`
- 词法元素菜单：`GET /api/v1/topic/lexical_elements`
- 词法元素子主题：`GET /api/v1/topic/lexical_elements/{chapter}`
- 常量菜单：`GET /api/v1/topic/constants`
- 常量子主题：`GET /api/v1/topic/constants/{subtopic}`
- **测验题库**：`GET /api/v1/quiz/{topic}/{chapter}/start`（随机抽题）
- **测验统计**：`GET /api/v1/quiz/{topic}/{chapter}/stats`（题库统计信息）

## 测验题库系统

### 概述

后端集成了智能测验题库系统，支持41个Go语言章节的随机抽题测验：

- **题库规模**：41个章节，1230-2050道题目
- **题型支持**：单选题、多选题（50%比例）
- **难度分布**：简单40%、中等40%、困难20%
- **智能抽题**：每次测验随机抽取3-5道单选题和3-5道多选题
- **选项打乱**：防止答案位置规律，提高测验公平性

### 题库文件组织

题库文件位于 `quiz_data/` 目录，按主题和章节组织：

```
quiz_data/
├── lexical_elements/     # 词法元素（11章节）
├── constants/           # 常量（12章节）
├── variables/           # 变量（4章节）
└── types/               # 类型（14章节）
```

### 质量保证

- **结构验证**：启动时自动验证所有题目格式完整性
- **答案正确性**：人工审核确保答案准确、解析详细
- **错误定位**：验证失败时提供文件名、题目ID、行号信息
- **测试覆盖**：单元测试覆盖率≥80%，包含格式验证和抽题逻辑

### 配置选项

在 `configs/config.yaml` 中配置：

```yaml
quiz:
  dataPath: "quiz_data"          # 题库文件根目录
  questionCount:
    single: 4                    # 单选题数量
    multiple: 4                  # 多选题数量
  difficultyDistribution:
    easy: 40                     # 简单题占比%
    medium: 40                   # 中等题占比%
    hard: 20                     # 困难题占比%
```

### 维护工具

- **质量检查**：`go run scripts/tools/quiz_quality_check.go`（批量验证题目质量）
- **数据验证**：`go run scripts/tools/validate/validate_quiz_data.go`（验证题库文件）

### 示例：查询章节统计

通过命令行查看某章节统计信息（需有效 access token）：

```bash
curl -H "Authorization: Bearer $ACCESS_TOKEN" \
  "http://127.0.0.1:8080/api/v1/quiz/constants/boolean/stats"
```

响应为统一格式（成功 code=20000），data 字段包含 total / byType / byDifficulty。

## 开发辅助

- 代码检查：`powershell -ExecutionPolicy Bypass -File backend/scripts/check-go.ps1`
- 覆盖率示例：`go test -cover ./...`

