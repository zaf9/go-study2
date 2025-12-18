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
- **测验提交**：`POST /api/v1/quiz/submit`（提交答案并记录结果）
- **测验历史**：`GET /api/v1/quiz/history?topic=...`（查询用户测验历史列表）
- **测验回顾**：`GET /api/v1/quiz/history/{sessionId}`（查询指定会话的详细回顾数据）

### 测验 API 详细说明

#### 1. 获取测验题目 - `GET /api/v1/quiz/{topic}/{chapter}/start`

**功能**：为指定主题和章节的用户开启新的测验会话，返回随机打乱的题目列表

**请求参数**：
- `topic`（路径）：章节主题，如 `constants`、`variables` 等
- `chapter`（路径）：子章节名称，如 `boolean`、`numeric` 等

**响应示例**：
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "topic": "constants",
    "chapter": "boolean",
    "sessionId": "uuid-string-xxx",
    "questions": [
      {
        "id": "q-123",
        "type": "single",
        "difficulty": "easy",
        "question": "iota 初始值是多少？",
        "options": [
          { "id": "opt-1", "label": "0" },
          { "id": "opt-2", "label": "1" },
          { "id": "opt-3", "label": "-1" },
          { "id": "opt-4", "label": "undefined" }
        ]
      }
    ]
  }
}
```

#### 2. 提交测验答案 - `POST /api/v1/quiz/submit`

**功能**：提交测验答案，计算百分制得分，判定是否通过，保存会话和答题记录

**请求体**：
```json
{
  "sessionId": "uuid-string-xxx",
  "topic": "constants",
  "chapter": "boolean",
  "answers": [
    { "questionId": "q-123", "choices": ["opt-1"] },
    { "questionId": "q-124", "choices": ["opt-2", "opt-3"] }
  ],
  "durationMs": 45000
}
```

**响应示例**：
```json
{
  "code": 20000,
  "message": "提交成功",
  "data": {
    "score": 85,
    "total_questions": 10,
    "correct_answers": 8,
    "passed": true,
    "details": [
      {
        "question_id": 1,
        "is_correct": true,
        "correct_answers": ["opt-1"],
        "explanation": "iota 在常量块中从 0 开始..."
      }
    ],
    "submittedAt": "2025-12-18T10:30:45Z"
  }
}
```

**百分制得分计算**：`score = (correct_answers / total_questions) * 100`，四舍五入到整数

**通过标准**：`score >= 60` 时 `passed=true`

#### 3. 获取测验历史 - `GET /api/v1/quiz/history?topic=...`

**功能**：查询当前用户的测验历史列表，按时间倒序返回最近 20 条记录

**请求参数**：
- `topic`（查询参数，可选）：按主题过滤，留空则返回全部历史

**响应示例**：
```json
{
  "code": 20000,
  "message": "success",
  "data": [
    {
      "id": 1,
      "sessionId": "uuid-1",
      "topic": "constants",
      "chapter": "boolean",
      "score": 85,
      "total_questions": 10,
      "correct_answers": 8,
      "passed": true,
      "completedAt": "2025-12-18T10:30:45Z"
    },
    {
      "id": 2,
      "sessionId": "uuid-2",
      "topic": "variables",
      "chapter": "declaration",
      "score": 60,
      "total_questions": 8,
      "correct_answers": 5,
      "passed": true,
      "completedAt": "2025-12-17T14:20:30Z"
    }
  ]
}
```

#### 4. 获取测验详情回顾 - `GET /api/v1/quiz/history/{sessionId}`

**功能**：返回指定会话的完整回顾数据，包括所有题目、用户答案、正确答案和解析

**请求参数**：
- `sessionId`（路径）：测验会话 ID

**响应示例**：
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "meta": {
      "sessionId": "uuid-1",
      "topic": "constants",
      "chapter": "boolean",
      "score": 85,
      "passed": true,
      "completedAt": "2025-12-18T10:30:45Z"
    },
    "questions": [
      {
        "questionId": 1,
        "stem": "iota 初始值是多少？",
        "options": ["0", "1", "-1", "undefined"],
        "userChoice": "0",
        "correctChoice": "0",
        "isCorrect": true,
        "explanation": "iota 在常量块中从 0 开始，后续递增..."
      },
      {
        "questionId": 2,
        "stem": "以下哪些是预定义常量？",
        "options": ["true", "false", "iota", "nil"],
        "userChoice": "true;false;iota",
        "correctChoice": "true;false;iota",
        "isCorrect": true,
        "explanation": "预定义常量包括 true、false、iota 和 nil..."
      }
    ]
  }
}
```

#### 5. 获取题库统计 - `GET /api/v1/quiz/{topic}/{chapter}/stats`

**功能**：返回指定章节的题库统计信息，包括总题量、按题型分布、按难度分布

**请求参数**：
- `topic`（路径）：章节主题
- `chapter`（路径）：子章节名称

**响应示例**：
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "total": 42,
    "byType": {
      "single": 21,
      "multiple": 21
    },
    "byDifficulty": {
      "easy": 16,
      "medium": 17,
      "hard": 9
    }
  }
}
```

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

