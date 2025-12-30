# 进度追踪 API 规范

**版本**: v1  
**基础路径**: `/api/v1/progress`  
**认证**: 所有接口需 JWT Bearer Token

---

## 1. 获取进度概览

### 接口信息
- **路径**: `GET /api/v1/progress/overview`
- **描述**: 获取当前用户的全局学习进度概览
- **认证**: 必需

### 请求

**Headers**:
```
Authorization: Bearer <access_token>
```

**Query Parameters**: 无

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "totalChapters": 50,
    "completedChapters": 12,
    "inProgressChapters": 5,
    "completionRate": 24.0,
    "topics": [
      {
        "topic": "lexical_elements",
        "totalChapters": 11,
        "completedChapters": 8,
        "inProgressChapters": 2
      },
      {
        "topic": "constants",
        "totalChapters": 12,
        "completedChapters": 3,
        "inProgressChapters": 2
      },
      {
        "topic": "variables",
        "totalChapters": 4,
        "completedChapters": 1,
        "inProgressChapters": 1
      },
      {
        "topic": "types",
        "totalChapters": 14,
        "completedChapters": 0,
        "inProgressChapters": 0
      }
    ],
    "next": {
      "topic": "constants",
      "chapter": "boolean",
      "title": "布尔常量"
    }
  }
}
```

**字段说明**:
- `totalChapters`: 所有主题的章节总数(从static-routes.ts计算)
- `completedChapters`: 已完成的章节数(status=completed且quiz_passed=true)
- `inProgressChapters`: 学习中的章节数(status=in_progress)
- `completionRate`: 完成率,计算公式: `(completedChapters / totalChapters) * 100`
- `topics`: 各主题的进度统计
- `next`: 建议下一个学习的章节(基于last_visit_at排序,选择学习中但未完成的最早章节)

**错误响应**:

401 Unauthorized:
```json
{
  "code": 40100,
  "message": "未授权访问,请先登录",
  "data": null
}
```

500 Internal Server Error:
```json
{
  "code": 50000,
  "message": "查询进度失败: [错误详情]",
  "data": null
}
```

---

## 2. 获取主题进度详情

### 接口信息
- **路径**: `GET /api/v1/progress/topic/:topic`
- **描述**: 获取指定主题的章节级别进度详情
- **认证**: 必需

### 请求

**Headers**:
```
Authorization: Bearer <access_token>
```

**Path Parameters**:
- `topic` (string, required): 主题标识,取值范围: `lexical_elements`, `constants`, `variables`, `types`

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "success",
  "data": {
    "topic": "constants",
    "totalChapters": 12,
    "chapters": [
      {
        "chapter": "boolean",
        "status": "completed",
        "quizScore": 8,
        "quizPassed": true,
        "lastVisitAt": "2025-12-28T10:30:00Z",
        "completedAt": "2025-12-28T10:45:00Z"
      },
      {
        "chapter": "rune",
        "status": "in_progress",
        "quizScore": 0,
        "quizPassed": false,
        "lastVisitAt": "2025-12-29T15:20:00Z",
        "completedAt": null
      },
      {
        "chapter": "integer",
        "status": "not_started",
        "quizScore": 0,
        "quizPassed": false,
        "lastVisitAt": null,
        "completedAt": null
      }
    ]
  }
}
```

**字段说明**:
- `totalChapters`: 该主题的章节总数(从static-routes.ts获取)
- `chapters`: 章节进度列表,包含该主题的所有章节
  - 有学习记录的章节: 从learning_progress表查询
  - 无学习记录的章节: 填充默认值(status=not_started)
- `status`: 章节状态 (`not_started` | `in_progress` | `completed`)
- `lastVisitAt`: 最后访问时间(ISO 8601格式)
- `completedAt`: 完成时间,仅当通过测验时有值

**错误响应**:

400 Bad Request (无效的topic):
```json
{
  "code": 40000,
  "message": "不支持的主题: invalid_topic",
  "data": null
}
```

404 Not Found (主题不存在):
```json
{
  "code": 40400,
  "message": "主题不存在: nonexistent",
  "data": null
}
```

---

## 3. 记录学习进度

### 接口信息
- **路径**: `POST /api/v1/progress`
- **描述**: 记录或更新用户的章节学习进度
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
  "topic": "constants",
  "chapter": "boolean",
  "readDuration": 120,
  "scrollProgress": 85,
  "lastPosition": "{\"scroll\":1200}",
  "quizScore": 0,
  "quizPassed": false,
  "estimatedSeconds": 180,
  "forceSync": false
}
```

**字段说明**:
- `topic` (string, required): 主题标识
- `chapter` (string, required): 章节标识
- `readDuration` (integer, optional): 阅读时长(秒),默认0
- `scrollProgress` (integer, optional): 滚动进度百分比(0-100),默认0
- `lastPosition` (string, optional): 最后位置(JSON字符串),默认空字符串
- `quizScore` (integer, optional): 测验得分,默认0
- `quizPassed` (boolean, optional): 是否通过测验,默认false
- `estimatedSeconds` (integer, optional): 预计学习时长,默认0
- `forceSync` (boolean, optional): 是否强制同步(用于beforeunload),默认false

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "进度已保存",
  "data": {
    "status": "in_progress",
    "lastVisitAt": "2025-12-30T10:30:00Z"
  }
}
```

**错误响应**:

400 Bad Request (参数校验失败):
```json
{
  "code": 40000,
  "message": "参数错误: topic不能为空",
  "data": null
}
```

422 Unprocessable Entity (无效的topic或chapter):
```json
{
  "code": 42200,
  "message": "不支持的主题或章节: constants/invalid_chapter",
  "data": null
}
```

---

## 4. 更新章节完成状态

### 接口信息
- **路径**: `PUT /api/v1/progress/complete`
- **描述**: 将章节标记为已完成(通常在测验通过后调用)
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
  "topic": "constants",
  "chapter": "boolean",
  "quizScore": 8,
  "quizPassed": true
}
```

**字段说明**:
- `topic` (string, required): 主题标识
- `chapter` (string, required): 章节标识
- `quizScore` (integer, required): 测验得分
- `quizPassed` (boolean, required): 是否通过测验

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "章节已标记为完成",
  "data": {
    "status": "completed",
    "completedAt": "2025-12-30T11:00:00Z"
  }
}
```

**说明**:
- 仅当`quizPassed=true`时,才将`status`更新为`completed`并设置`completedAt`
- 如果`quizPassed=false`,保持`status=in_progress`

---

## 5. 删除学习记录

### 接口信息
- **路径**: `DELETE /api/v1/progress/:topic/:chapter`
- **描述**: 删除指定章节的学习记录(用于重新开始学习)
- **认证**: 必需

### 请求

**Headers**:
```
Authorization: Bearer <access_token>
```

**Path Parameters**:
- `topic` (string, required): 主题标识
- `chapter` (string, required): 章节标识

### 响应

**成功响应** (200 OK):
```json
{
  "code": 20000,
  "message": "学习记录已删除",
  "data": null
}
```

**错误响应**:

404 Not Found (记录不存在):
```json
{
  "code": 40400,
  "message": "未找到该章节的学习记录",
  "data": null
}
```

---

## 业务逻辑说明

### 状态转换规则

```
未开始 (not_started)
    ↓ 用户访问章节内容页
学习中 (in_progress)
    ↓ 完成测验且通过
已完成 (completed)
```

**详细说明**:
1. **未开始 → 学习中**: 
   - 触发: 用户首次访问章节内容页(`/topics/:topic/:chapter`)
   - 操作: 调用`POST /api/v1/progress`,创建记录,设置`status=in_progress`

2. **学习中 → 已完成**:
   - 触发: 用户完成测验且得分≥60%(quiz_passed=true)
   - 操作: 测验提交成功后,后端自动调用进度更新,设置`status=completed`和`completedAt`

3. **已完成 → 学习中**(重新学习):
   - 触发: 用户点击"重新学习"或删除学习记录
   - 操作: 调用`DELETE /api/v1/progress/:topic/:chapter`,删除记录后重新开始

### 章节总数计算

- **数据源**: 前端`static-routes.ts`中的`topicChapters`映射
- **计算方式**: `topicChapters[topic].length`
- **同步机制**: 后端不维护章节列表,完全依赖前端静态路由定义
- **优点**: 单一数据源,避免前后端数据不一致

### 进度统计规则

1. **已完成章节**: `status = 'completed' AND quiz_passed = true`
2. **学习中章节**: `status = 'in_progress'`
3. **未开始章节**: 在`topicChapters[topic]`中但`learning_progress`表中无记录
4. **完成率**: `(已完成章节数 / 总章节数) * 100`,保留1位小数

### 并发控制

- **幂等性**: `POST /api/v1/progress`使用`INSERT OR REPLACE`(SQLite)或`ON CONFLICT UPDATE`(PostgreSQL),保证幂等
- **乐观锁**: `learning_progress`表使用`updated_at`字段,检测并发更新冲突
- **事务**: 测验提交和进度更新在同一事务中完成,保证原子性

---

## 错误码规范

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 20000 | 200 | 成功 |
| 40000 | 400 | 请求参数错误 |
| 40100 | 401 | 未授权(token无效或缺失) |
| 40400 | 404 | 资源不存在 |
| 42200 | 422 | 请求格式正确但语义错误 |
| 50000 | 500 | 服务器内部错误 |
| 50300 | 503 | 服务暂时不可用 |

---

## 示例场景

### 场景1: 用户首次访问章节

**请求**:
```http
POST /api/v1/progress
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "topic": "constants",
  "chapter": "boolean",
  "readDuration": 0,
  "scrollProgress": 0
}
```

**响应**:
```json
{
  "code": 20000,
  "message": "进度已保存",
  "data": {
    "status": "in_progress",
    "lastVisitAt": "2025-12-30T10:30:00Z"
  }
}
```

### 场景2: 用户完成测验

**流程**:
1. 前端调用测验提交API: `POST /api/v1/quiz/submit`
2. 后端在测验提交成功后,内部调用进度更新逻辑
3. 如果`quiz_passed=true`,更新`status=completed`

**前端无需额外调用进度API**,后端自动完成状态更新。

### 场景3: 查看进度概览

**请求**:
```http
GET /api/v1/progress/overview
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**: 见"接口1"示例

---

## 性能优化建议

1. **批量查询**: 
   - `GET /api/v1/progress/overview`一次性查询所有进度记录,避免N+1查询
   - 使用`IN`查询或`JOIN`优化多主题查询

2. **缓存策略**:
   - 前端使用SWR缓存,缓存时间5分钟
   - 后端可考虑Redis缓存用户进度概览,TTL 1分钟

3. **索引优化**:
   - 确保`learning_progress`表有`(user_id, topic)`复合索引
   - 添加`(user_id, status)`索引加速按状态筛选

4. **数据传输优化**:
   - 仅返回必要字段,避免传输冗余数据
   - 对于大量章节,考虑分页(当前暂不需要)

---

## 版本历史

| 版本 | 日期 | 变更说明 |
|------|------|----------|
| v1.0 | 2025-12-30 | 初始版本,定义进度追踪核心API |

---

## 附录: 前端集成示例

```typescript
// frontend/services/progressService.ts
import api from '@/lib/api';

export async function fetchProgressOverview() {
  const response = await api.get('/api/v1/progress/overview');
  return response.data;
}

export async function fetchTopicProgress(topic: string) {
  const response = await api.get(`/api/v1/progress/topic/${topic}`);
  return response.data;
}

export async function updateProgress(payload: UpdateProgressPayload) {
  const response = await api.post('/api/v1/progress', payload);
  return response.data;
}

// 使用SWR
import useSWR from 'swr';

export function useProgressOverview() {
  return useSWR('progress/overview', fetchProgressOverview, {
    revalidateOnFocus: false,
    dedupingInterval: 60000, // 1分钟去重
  });
}
```
