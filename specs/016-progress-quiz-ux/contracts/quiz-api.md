# 测验功能 API 规范

**版本**: v1  
**基础路径**: `/api/v1/quiz`  
**认证**: 所有接口需 JWT Bearer Token

---

## 1. 获取或创建测验会话

### 接口信息
- **路径**: `GET /api/v1/quiz/:topic/:chapter`
- **描述**: 获取或创建指定章节的测验会话,返回测验题目(不含答案)
- **认证**: 必需

### 请求

**Headers**:
```
Authorization: Bearer <access_token>
```

**Path Parameters**:
- `topic` (string, required): 主题标识,取值: `lexical_elements`, `constants`, `variables`, `types`
- `chapter` (string, required): 章节标识,例如: `boolean`, `rune`, `integer`

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "sessionId": "a1b2c3d4e5f6g7h8i9j0",
    "topic": "constants",
    "chapter": "boolean",
    "questions": [
      {
        "id": 1,
        "type": "single",
        "difficulty": "easy",
        "question": "Go语言中,布尔类型的零值是什么?",
        "options": [
          {"id": "A", "label": "true"},
          {"id": "B", "label": "false"},
          {"id": "C", "label": "nil"},
          {"id": "D", "label": "0"}
        ],
        "codeSnippet": null
      },
      {
        "id": 2,
        "type": "multiple",
        "difficulty": "medium",
        "question": "以下哪些是Go语言中的预定义布尔常量?",
        "options": [
          {"id": "A", "label": "true"},
          {"id": "B", "label": "false"},
          {"id": "C", "label": "yes"},
          {"id": "D", "label": "no"}
        ],
        "codeSnippet": null
      },
      {
        "id": 3,
        "type": "single",
        "difficulty": "hard",
        "question": "以下代码的输出是什么?",
        "options": [
          {"id": "A", "label": "true"},
          {"id": "B", "label": "false"},
          {"id": "C", "label": "编译错误"},
          {"id": "D", "label": "运行时错误"}
        ],
        "codeSnippet": "package main\n\nfunc main() {\n\tvar b bool\n\tprintln(b)\n}"
      }
    ]
  }
}
```

**字段说明**:
- `sessionId`: 会话唯一标识符(MD5哈希),用于提交答案时关联
- `questions`: 题目列表
  - `id`: 题目ID(整数)
  - `type`: 题型 (`single`: 单选, `multiple`: 多选)
  - `difficulty`: 难度 (`easy`, `medium`, `hard`)
  - `question`: 题干
  - `options`: 选项列表
    - `id`: 选项ID(通常为A/B/C/D)
    - `label`: 选项文本
  - `codeSnippet`: 代码片段(可选),用于代码分析题

**业务逻辑**:
1. 后端检查是否存在该用户在该章节的活跃会话(24小时内)
2. 如存在,直接返回该会话及题目
3. 如不存在,创建新会话:
   - 从`quiz_data/:topic/:chapter.yaml`加载题目
   - 生成唯一`sessionId`
   - 保存到`quiz_sessions`表
   - 返回题目(不含答案)

**错误响应**:

400 Bad Request (无效的topic或chapter):
```json
{
  "code": 40000,
  "message": "不支持的主题或章节: invalid_topic/invalid_chapter",
  "data": null
}
```

404 Not Found (该章节无测验题目):
```json
{
  "code": 40400,
  "message": "章节 constants/boolean 暂无测验题目",
  "data": null
}
```

500 Internal Server Error (题目加载失败):
```json
{
  "code": 50000,
  "message": "加载题目失败: 文件不存在或格式错误",
  "data": null
}
```

---

## 2. 提交测验答案

### 接口信息
- **路径**: `POST /api/v1/quiz/submit`
- **描述**: 提交用户的测验答案,返回评分结果和正确答案
- **认证**: 必需

### 请求

**Headers**:
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Body**:
```json
{
  "sessionId": "a1b2c3d4e5f6g7h8i9j0",
  "topic": "constants",
  "chapter": "boolean",
  "durationMs": 120000,
  "answers": [
    {
      "questionId": 1,
      "userAnswers": ["B"]
    },
    {
      "questionId": 2,
      "userAnswers": ["A", "B"]
    },
    {
      "questionId": 3,
      "userAnswers": ["A"]
    }
  ]
}
```

**字段说明**:
- `sessionId` (string, required): 会话ID,从`GET /quiz/:topic/:chapter`获取
- `topic` (string, required): 主题标识
- `chapter` (string, required): 章节标识
- `durationMs` (integer, optional): 答题耗时(毫秒),默认0
- `answers` (array, required): 答案列表
  - `questionId` (integer, required): 题目ID
  - `userAnswers` (array of string, required): 用户选择的答案(单选题数组长度为1,多选题可为多个)

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "提交成功",
  "data": {
    "sessionId": "a1b2c3d4e5f6g7h8i9j0",
    "score": 66,
    "total": 3,
    "correctCount": 2,
    "passed": true,
    "details": [
      {
        "questionId": 1,
        "question": "Go语言中,布尔类型的零值是什么?",
        "userAnswers": ["B"],
        "correctAnswers": ["B"],
        "isCorrect": true
      },
      {
        "questionId": 2,
        "question": "以下哪些是Go语言中的预定义布尔常量?",
        "userAnswers": ["A", "B"],
        "correctAnswers": ["A", "B"],
        "isCorrect": true
      },
      {
        "questionId": 3,
        "question": "以下代码的输出是什么?",
        "userAnswers": ["A"],
        "correctAnswers": ["B"],
        "isCorrect": false
      }
    ]
  }
}
```

**字段说明**:
- `score`: 得分百分比(0-100)
- `total`: 总题数
- `correctCount`: 答对的题数
- `passed`: 是否通过(score >= 60)
- `details`: 每道题的详细结果
  - `questionId`: 题目ID
  - `question`: 题干
  - `userAnswers`: 用户的答案
  - `correctAnswers`: 正确答案
  - `isCorrect`: 是否答对

**业务逻辑**:
1. 验证`sessionId`有效性(存在且未提交)
2. 校验答案格式(题目ID存在,答案非空)
3. 评分:
   - 单选题: 完全正确得1分,否则0分
   - 多选题: 完全正确得1分,部分正确得0.5分,错误0分
4. 计算总分: `score = (correctCount / total) * 100`
5. 判断通过: `passed = score >= 60`
6. 保存记录:
   - 更新`quiz_sessions.submitted_at`(防重复提交)
   - 插入`quiz_attempts`(每题的答题详情)
   - 调用进度服务更新章节状态(如通过,标记为completed)
7. 返回结果和正确答案

**错误响应**:

400 Bad Request (参数错误):
```json
{
  "code": 40000,
  "message": "参数错误: sessionId不能为空",
  "data": null
}
```

404 Not Found (会话不存在):
```json
{
  "code": 40400,
  "message": "测验会话不存在或已过期",
  "data": null
}
```

409 Conflict (重复提交):
```json
{
  "code": 40900,
  "message": "该测验已提交,请勿重复操作",
  "data": null
}
```

422 Unprocessable Entity (答案格式错误):
```json
{
  "code": 42200,
  "message": "答案格式错误: 题目ID 999 不存在",
  "data": null
}
```

500 Internal Server Error:
```json
{
  "code": 50000,
  "message": "提交失败: 数据库错误",
  "data": null
}
```

---

## 3. 获取测验历史记录

### 接口信息
- **路径**: `GET /api/v1/quiz/history`
- **描述**: 获取用户的测验历史记录
- **认证**: 必需

### 请求

**Headers**:
```
Authorization: Bearer <access_token>
```

**Query Parameters**:
- `topic` (string, optional): 主题标识,用于筛选特定主题的测验记录
- `from` (string, optional): 起始时间,ISO 8601格式,例如: `2025-12-01T00:00:00Z`
- `to` (string, optional): 结束时间,ISO 8601格式
- `limit` (integer, optional): 返回记录数量上限,默认50,最大100
- `offset` (integer, optional): 分页偏移量,默认0

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "total": 15,
    "records": [
      {
        "sessionId": "a1b2c3d4e5f6g7h8i9j0",
        "topic": "constants",
        "chapter": "boolean",
        "score": 80,
        "total": 5,
        "correctCount": 4,
        "passed": true,
        "durationMs": 180000,
        "createdAt": "2025-12-28T10:45:00Z"
      },
      {
        "sessionId": "z9y8x7w6v5u4t3s2r1q0",
        "topic": "constants",
        "chapter": "rune",
        "score": 50,
        "total": 4,
        "correctCount": 2,
        "passed": false,
        "durationMs": 150000,
        "createdAt": "2025-12-27T15:30:00Z"
      }
    ]
  }
}
```

**字段说明**:
- `total`: 符合条件的总记录数
- `records`: 测验记录列表(按`createdAt`降序)
  - `sessionId`: 会话ID,可用于查询详情
  - `topic`: 主题标识
  - `chapter`: 章节标识
  - `score`: 得分(0-100)
  - `total`: 总题数
  - `correctCount`: 答对题数
  - `passed`: 是否通过
  - `durationMs`: 答题耗时(毫秒)
  - `createdAt`: 测验完成时间

**错误响应**:

400 Bad Request (时间格式错误):
```json
{
  "code": 40000,
  "message": "参数错误: from时间格式无效",
  "data": null
}
```

---

## 4. 获取测验详情

### 接口信息
- **路径**: `GET /api/v1/quiz/detail/:sessionId`
- **描述**: 获取某次测验的详细答题情况
- **认证**: 必需

### 请求

**Headers**:
```
Authorization: Bearer <access_token>
```

**Path Parameters**:
- `sessionId` (string, required): 会话ID

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "sessionId": "a1b2c3d4e5f6g7h8i9j0",
    "topic": "constants",
    "chapter": "boolean",
    "score": 66,
    "total": 3,
    "correctCount": 2,
    "passed": true,
    "durationMs": 120000,
    "createdAt": "2025-12-28T10:45:00Z",
    "details": [
      {
        "questionId": 1,
        "question": "Go语言中,布尔类型的零值是什么?",
        "type": "single",
        "difficulty": "easy",
        "userAnswers": ["B"],
        "correctAnswers": ["B"],
        "isCorrect": true,
        "options": [
          {"id": "A", "label": "true"},
          {"id": "B", "label": "false"},
          {"id": "C", "label": "nil"},
          {"id": "D", "label": "0"}
        ]
      },
      {
        "questionId": 2,
        "question": "以下哪些是Go语言中的预定义布尔常量?",
        "type": "multiple",
        "difficulty": "medium",
        "userAnswers": ["A", "B"],
        "correctAnswers": ["A", "B"],
        "isCorrect": true,
        "options": [
          {"id": "A", "label": "true"},
          {"id": "B", "label": "false"},
          {"id": "C", "label": "yes"},
          {"id": "D", "label": "no"}
        ]
      },
      {
        "questionId": 3,
        "question": "以下代码的输出是什么?",
        "type": "single",
        "difficulty": "hard",
        "userAnswers": ["A"],
        "correctAnswers": ["B"],
        "isCorrect": false,
        "options": [
          {"id": "A", "label": "true"},
          {"id": "B", "label": "false"},
          {"id": "C", "label": "编译错误"},
          {"id": "D", "label": "运行时错误"}
        ]
      }
    ]
  }
}
```

**错误响应**:

404 Not Found (会话不存在):
```json
{
  "code": 40400,
  "message": "测验记录不存在",
  "data": null
}
```

403 Forbidden (无权访问):
```json
{
  "code": 40300,
  "message": "无权访问该测验记录",
  "data": null
}
```

---

## 5. 删除测验记录(重新测验)

### 接口信息
- **路径**: `DELETE /api/v1/quiz/:sessionId`
- **描述**: 删除指定的测验记录,允许用户重新测验
- **认证**: 必需

### 请求

**Headers**:
```
Authorization: Bearer <access_token>
```

**Path Parameters**:
- `sessionId` (string, required): 会话ID

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "测验记录已删除,您可以重新开始测验",
  "data": null
}
```

**说明**:
- 删除`quiz_sessions`记录(级联删除`quiz_attempts`)
- 不影响`learning_progress`表的其他字段
- 如果该测验通过导致章节标记为completed,需手动调用进度API重置状态

**错误响应**:

404 Not Found:
```json
{
  "code": 40400,
  "message": "测验记录不存在",
  "data": null
}
```

---

## 业务逻辑说明

### 会话生命周期

1. **创建会话**: 用户访问`/quiz/:topic/:chapter`
   - 检查24小时内是否有活跃会话
   - 无 → 创建新会话,加载题目
   - 有 → 返回现有会话(避免重复创建)

2. **答题**: 前端保存用户选择(仅在内存中)

3. **提交**: 用户点击"提交"
   - 调用`POST /api/v1/quiz/submit`
   - 后端评分并保存结果
   - 更新`quiz_sessions.submitted_at`防止重复提交
   - 如通过,更新章节进度为completed

4. **查看历史**: 用户访问测验中心
   - 调用`GET /api/v1/quiz/history`获取列表
   - 点击记录查看详情: `GET /api/v1/quiz/detail/:sessionId`

### 评分规则

**单选题**:
- 完全正确: 1分
- 错误: 0分

**多选题**:
- 完全正确: 1分
- 部分正确(选了正确答案但也选了错误答案): 0分
- 漏选(只选了部分正确答案): 0.5分
- 完全错误: 0分

**总分计算**:
```
score = (获得分数总和 / 题目总数) * 100
```

**通过标准**:
```
passed = score >= 60
```

### 防重复提交机制

**前端**:
1. 使用`useRef`和`state`双重锁定
2. 提交按钮disabled状态
3. 提交成功后不重置锁定状态

**后端**:
1. 检查`quiz_sessions.submitted_at`字段
2. 非NULL → 返回409 Conflict
3. NULL → 执行提交,更新字段为当前时间
4. 使用数据库事务保证原子性

### 题目加载逻辑

**文件路径规范**:
```
quiz_data/
├── lexical_elements/
│   ├── comments.yaml
│   ├── tokens.yaml
│   └── ...
├── constants/
│   ├── boolean.yaml
│   ├── rune.yaml
│   └── ...
├── variables/
│   └── ...
└── types/
    └── ...
```

**YAML格式示例**:
```yaml
- id: 1
  type: single
  difficulty: easy
  question: "Go语言中,布尔类型的零值是什么?"
  options:
    - id: A
      label: "true"
    - id: B
      label: "false"
    - id: C
      label: "nil"
    - id: D
      label: "0"
  answer:
    - B
  codeSnippet: null

- id: 2
  type: multiple
  difficulty: medium
  question: "以下哪些是Go语言中的预定义布尔常量?"
  options:
    - id: A
      label: "true"
    - id: B
      label: "false"
    - id: C
      label: "yes"
    - id: D
      label: "no"
  answer:
    - A
    - B
  codeSnippet: null
```

**加载流程**:
1. 根据`topic`和`chapter`拼接文件路径
2. 检查文件是否存在,不存在返回空切片(非错误)
3. 读取YAML文件
4. 解析为`[]QuizQuestion`
5. **返回给前端时移除`answer`字段**(安全考虑)

---

## 错误处理规范

### 错误码

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 20000 | 200 | 成功 |
| 40000 | 400 | 请求参数错误 |
| 40100 | 401 | 未授权 |
| 40300 | 403 | 禁止访问 |
| 40400 | 404 | 资源不存在 |
| 40900 | 409 | 冲突(重复提交) |
| 42200 | 422 | 语义错误 |
| 50000 | 500 | 服务器内部错误 |

### 日志记录

**关键操作日志**:
1. 会话创建: `INFO: 用户{userID}创建测验会话 {topic}/{chapter}, sessionId={sessionId}`
2. 题目加载失败: `ERROR: 加载题目失败 {topic}/{chapter}: {error}`
3. 重复提交尝试: `WARN: 用户{userID}尝试重复提交sessionId={sessionId}`
4. 提交成功: `INFO: 用户{userID}提交测验 {sessionId}, 得分={score}, 通过={passed}`

---

## 性能优化

1. **题目缓存**: 
   - 后端缓存已加载的YAML题目(内存或Redis)
   - TTL 30分钟,减少磁盘I/O

2. **会话复用**:
   - 24小时内复用活跃会话,避免重复创建
   - 清理过期会话(定时任务,删除24小时前的未提交会话)

3. **批量查询优化**:
   - 测验历史查询使用索引: `(user_id, created_at DESC)`
   - 限制返回记录数(最多100条)

4. **事务优化**:
   - 提交测验时使用事务,确保原子性
   - 批量插入`quiz_attempts`减少往返

---

## 安全考虑

1. **答案保护**: 
   - 返回题目时**绝对不能**包含`answer`字段
   - 后端评分时从服务器端题目数据获取答案,不信任前端

2. **权限验证**:
   - 所有接口校验JWT token
   - 查询/删除测验记录时验证`user_id`匹配

3. **输入校验**:
   - 严格校验`topic`和`chapter`枚举值
   - 校验`questionId`存在性
   - 校验`userAnswers`非空且格式正确

4. **防止时间攻击**:
   - 评分逻辑统一处理时间,避免暴露答案信息

---

## 版本历史

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0 | 2025-12-30 | 初始版本,定义测验核心API |

---

## 附录: 前端集成示例

```typescript
// frontend/services/quizService.ts
import api from '@/lib/api';

export async function fetchQuizSession(topic: string, chapter: string) {
  const response = await api.get(`/api/v1/quiz/${topic}/${chapter}`);
  return response.data;
}

export async function submitQuiz(payload: QuizSubmitPayload) {
  const response = await api.post('/api/v1/quiz/submit', payload);
  return response.data;
}

export async function fetchQuizHistory(topic?: string) {
  const url = topic 
    ? `/api/v1/quiz/history?topic=${topic}` 
    : '/api/v1/quiz/history';
  const response = await api.get(url);
  return response.data.records;
}

// 使用示例
import { useQuiz } from '@/hooks/useQuiz';

function QuizPage({ topic, chapter }) {
  const { session, questions, submit, result } = useQuiz(topic, chapter);
  
  const handleSubmit = async () => {
    try {
      await submit();
      // 跳转到结果页或刷新进度
    } catch (error) {
      if (error.code === 409) {
        alert('已提交,请勿重复操作');
      }
    }
  };
  
  return (
    <div>
      {/* 渲染题目和选项 */}
      <button onClick={handleSubmit}>提交测验</button>
    </div>
  );
}
```
