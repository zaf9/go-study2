# Contracts: Types 学习模块（CLI / HTTP）

## CLI 契约
- 入口：`go run main.go`，主菜单编号 `3` 为 Types。编号从 0 开始，输入 `q` 返回上一层。  
- 子菜单：按顺序列出 `boolean|numeric|string|array|slice|struct|pointer|function|interface|map|channel`。显示格式：`<index>. <title>`。  
- 内容展示：选择子主题后输出标题、摘要、规则列表（含规则编号/说明）、正反例代码块（含预期输出）、测验题目与答案解析。  
- 测验：显示题干与编号选项，用户输入选项序号后返回得分、正确答案与解析，可提示重做。  
- 检索：输入命令 `search <keyword>` 返回匹配的规则摘要与正反例锚点（菜单编号+HTTP 路径）。  
- 错误处理：非法编号或关键字时提示可用选项，并保持中文错误信息；任何失败均不退出主菜单。

## HTTP 契约
- 通用：支持 `format=json|html`（默认 json），错误时返回 `code` 非 0 的 JSON，HTML 模式保持页面内提示。

### GET `/api/v1/topic/types`
- 描述：返回 Types 菜单。  
- Response 200 (json)：`{ "code":0, "message":"OK", "data": { "items":[ { "id":0, "title":"Boolean", "name":"boolean" }, ... ] } }`

### GET `/api/v1/topic/types/{subtopic}`
- Path 参数：`subtopic` ∈ `boolean|numeric|string|array|slice|struct|pointer|function|interface|map|channel`.  
- Response 200 (json)：  
  ```json
  {
    "code": 0,
    "message": "OK",
    "data": {
      "content": {
        "topic": "array",
        "title": "数组与递归限制",
        "summary": "...",
        "rules": ["长度为类型的一部分", "禁止数组/struct 纯递归嵌套"],
        "details": ["...", "..."],
        "snippet": "code",
        "examples": [
          { "id":"ex1", "title":"合法数组", "code":"...", "output":"...", "isValid":true, "ruleRef":"TR-ARRAY-LEN" },
          { "id":"ex2", "title":"非法递归", "code":"...", "output":"compile error", "isValid":false, "ruleRef":"TR-RECURSION" }
        ]
      },
      "quiz": [
        { "id":"q1", "stem":"数组长度与类型关系？", "options":["A","B","C"], "answer":"A", "explanation":"...", "ruleRef":"TR-ARRAY-LEN" }
      ]
    }
  }
  ```
- Response 404：当 `subtopic` 不支持或内容缺失，`{ "code":404, "message":"Subtopic not found" }`；HTML 显示错误页并附返回链接。

### POST `/api/v1/topic/types/quiz/submit`
- Request Body (json)：`{ "answers": [ { "id":"q1", "choice":"A" }, ... ] }`。  
- Response 200：`{ "code":0, "message":"OK", "data":{ "score":5, "total":5, "details":[ { "id":"q1", "correct":true, "answer":"A", "explanation":"...", "ruleRef":"TR-ARRAY-LEN" } ] } }`

### GET `/api/v1/topic/types/search`
- Query：`keyword=<string>`，允许多关键词用逗号分隔。  
- Response 200：`{ "code":0, "message":"OK", "data":{ "results":[ { "keyword":"map key", "conceptId":"map", "summary":"map 键需可比较", "positiveExampleId":"ex_map_ok", "negativeExampleId":"ex_map_func", "anchors":{"http":"/api/v1/topic/types/map","cli":"3 > map"} } ] } }`  
- Response 400：keyword 为空或无匹配时返回可用关键词提示。

