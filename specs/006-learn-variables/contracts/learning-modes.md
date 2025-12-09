# Contracts: CLI 与 HTTP 学习模式

## CLI 契约
- 命令入口：`go run ./cmd/cli --chapter variables`
- 主要交互：
  - `list` 显示子主题：变量存储/静态类型/动态类型/零值
  - `show <topic>` 输出对应内容与示例代码
  - `quiz` 进入测验，按编号答题并返回正确/错误与解析
- 输出要求：
  - 中文说明；示例代码段与期望输出；测验结果包含正确答案与解析
  - 错误输入需提示可用命令并返回上级菜单（符合分层菜单原则）

## HTTP 契约
- `GET /api/variables/content?topic=<storage|static|dynamic|zero>`
  - Response 200: `{ "topic": "storage", "summary": "...", "examples": [...], "snippet": "..." }`
- `GET /api/variables/quiz`
  - Response 200: `{ "items": [ { "id": "q1", "topic": "storage", "stem": "...", "options": ["A","B"], "answer": "A" } ] }`
- `POST /api/variables/quiz/submit`
  - Request: `{ "answers": [ { "id": "q1", "choice": "A" } ] }`
  - Response 200: `{ "score": 5, "total": 5, "details": [ { "id": "q1", "correct": true, "explanation": "..." } ] }`
- 通用约束：
  - Content-Type 为 application/json；中文字段说明；错误时返回400并给出可用 topic/格式提示
  - CLI 与 HTTP 返回的题目与解析需保持一致