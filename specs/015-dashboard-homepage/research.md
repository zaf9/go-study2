# Research: Dashboard 首页功能

**Feature**: Dashboard 首页  
**Date**: 2025-12-26  
**Related**: [spec.md](./spec.md) | [plan.md](./plan.md)

## Overview

本文档记录 Dashboard 功能实施前的技术研究和决策过程，解决规格说明中的所有技术不确定性。

## Research Topics

### R1: WebSocket 实现方案

**Research Question**: 如何在 Go + Next.js 架构中实现 WebSocket 实时数据推送？

**Options Evaluated**:

1. **gorilla/websocket** (第三方库)
   - Pros: 成熟稳定，社区活跃，文档完善
   - Cons: 需要额外依赖
   - Performance: 支持高并发，经过大规模生产验证

2. **GoFrame 内置 WebSocket**
   - Pros: 与框架集成，无额外依赖
   - Cons: 需要确认 GoFrame v2.9.5 是否内置 WebSocket 支持
   - Performance: 未知（需要验证）

3. **Server-Sent Events (SSE)**
   - Pros: 单向推送简单，HTTP 协议
   - Cons: 不支持客户端发送消息，功能受限
   - Performance: 适合单向推送场景

4. **轮询 (Polling)**
   - Pros: 实现简单，兼容性好
   - Cons: 性能差，服务器负载高，延迟高
   - Performance: 不适合实时场景

**Decision**: 使用 **gorilla/websocket**

**Rationale**:
- gorilla/websocket 是 Go 生态最成熟的 WebSocket 库
- 经过大规模生产验证，性能和稳定性有保障
- 文档完善，社区活跃，问题容易解决
- GoFrame v2.9.5 的 WebSocket 支持需要进一步确认，优先使用成熟方案

**Implementation Notes**:
```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // 生产环境需要验证 Origin
    },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()
    
    // 处理连接...
}
```

---

### R2: 学习天数计算策略

**Research Question**: 如何准确计算用户的累计学习天数？

**Options Evaluated**:

1. **自然天数（从首次学习到当前）**
   - Pros: 计算简单
   - Cons: 包含未学习的日期，数字虚高，误导用户
   - SQL: `SELECT DATEDIFF(NOW(), MIN(created_at)) FROM learning_progress WHERE user_id = ?`

2. **有学习活动的不同日期数**
   - Pros: 准确反映实际学习投入，激励持续学习
   - Cons: 计算稍复杂
   - SQL: `SELECT COUNT(DISTINCT DATE(last_visited_at)) FROM learning_progress WHERE user_id = ?`

3. **连续学习天数（中断后重置）**
   - Pros: 激励连续学习
   - Cons: 中断后重置可能打击用户积极性
   - SQL: 需要复杂的窗口函数计算

**Decision**: **有学习活动的不同日期数**

**Rationale**:
- 更准确反映用户实际学习投入
- 不会因为中断学习而产生误导性数字
- 激励用户持续学习，但不会因偶尔中断而惩罚用户
- 计算逻辑清晰，易于实现和测试

**Implementation Notes**:
```sql
-- PostgreSQL / MySQL
SELECT COUNT(DISTINCT DATE(last_visited_at)) 
FROM learning_progress 
WHERE user_id = ? AND last_visited_at IS NOT NULL;

-- SQLite
SELECT COUNT(DISTINCT DATE(last_visited_at)) 
FROM learning_progress 
WHERE user_id = ? AND last_visited_at IS NOT NULL;
```

**Performance Considerations**:
- 对于大量学习记录，`COUNT(DISTINCT)` 可能较慢
- 建议在 `last_visited_at` 字段上创建索引
- 考虑缓存结果（5-10 分钟），减少数据库查询

---

### R3: 最后学习记录获取方式

**Research Question**: 如何获取用户最后一次学习的主题和章节？

**Options Evaluated**:

1. **localStorage 客户端记录**
   - Pros: 无需后端支持，实现简单
   - Cons: 数据可能不准确，跨设备不同步，用户清除缓存后丢失
   - Reliability: 低

2. **从现有进度列表前端计算**
   - Pros: 无需新增 API
   - Cons: 增加前端复杂度，性能较差（需要加载所有进度数据）
   - Reliability: 中

3. **新增后端 API 接口**
   - Pros: 数据准确，跨设备同步，性能好
   - Cons: 需要后端开发
   - Reliability: 高

**Decision**: **新增 `/api/v1/progress/last` 接口**

**Rationale**:
- 后端计算更可靠，数据准确性有保障
- 跨设备同步，用户体验更好
- 性能优于前端计算（仅返回一条记录）
- 符合前后端分离架构原则

**API Design**:
```
GET /api/v1/progress/last
Response: {
  "code": 0,
  "message": "success",
  "data": {
    "topic_id": "lexical-elements",
    "topic_display_name": "词法元素",
    "chapter_id": "identifiers",
    "chapter_display_name": "标识符",
    "last_visited_at": "2025-12-26T10:30:00+08:00"
  }
}
```

---

### R4: 身份验证集成方案

**Research Question**: Dashboard 如何与现有身份验证系统集成？

**Options Evaluated**:

1. **独立身份验证逻辑**
   - Pros: 完全控制
   - Cons: 重复开发，安全风险高，用户体验不一致
   - Security: 低

2. **复用现有身份验证机制**
   - Pros: 保持系统一致性，避免重复开发，降低安全风险
   - Cons: 需要了解现有机制
   - Security: 高

3. **仅依赖前端路由守卫**
   - Pros: 实现简单
   - Cons: 不安全，后端 API 无保护
   - Security: 极低

**Decision**: **复用现有身份验证机制（JWT 或 Session）**

**Rationale**:
- 保持系统一致性，用户无需重复登录
- 避免重复开发，降低维护成本
- 降低安全风险，复用已验证的安全机制
- 符合最佳实践

**Implementation Notes**:
```go
// 使用现有中间件保护 Dashboard API
router.Group("/api/v1", func(group *ghttp.RouterGroup) {
    group.Middleware(middleware.Auth) // 现有认证中间件
    
    group.GET("/progress/last", controller.Progress.GetLastLearning)
})

// WebSocket 连接也需要认证
func handleWebSocket(r *ghttp.Request) {
    // 从 URL 参数或 Cookie 中获取 token
    token := r.Get("token").String()
    
    // 验证 token
    userID, err := middleware.ValidateToken(token)
    if err != nil {
        r.Response.WriteStatus(401)
        return
    }
    
    // 建立 WebSocket 连接
    // ...
}
```

---

### R5: 时间格式化策略

**Research Question**: 测验记录的完成时间应该如何显示？

**Options Evaluated**:

1. **仅使用相对时间**
   - Pros: 直观易懂
   - Cons: 时间久远后不够精确（如"30 天前"）
   - Example: "2 小时前"、"3 天前"、"1 个月前"

2. **仅使用绝对时间**
   - Pros: 精确
   - Cons: 近期时间不够直观
   - Example: "2025-12-26 10:30"

3. **混合格式（24 小时内相对 + 超过 24 小时绝对）**
   - Pros: 兼顾直观性和精确性
   - Cons: 实现稍复杂
   - Example: "2 小时前" 或 "2025-12-26 10:30"

**Decision**: **混合格式（24 小时内相对 + 超过 24 小时绝对）**

**Rationale**:
- 符合主流应用习惯（GitHub、Twitter、微信等）
- 近期活动用相对时间更直观
- 较早记录用绝对时间更清晰
- 提供最佳用户体验

**Implementation Notes**:
```typescript
function formatTime(timestamp: string): string {
  const now = new Date()
  const time = new Date(timestamp)
  const diffMs = now.getTime() - time.getTime()
  const diffHours = diffMs / (1000 * 60 * 60)
  
  if (diffHours < 24) {
    // 24 小时内，使用相对时间
    if (diffHours < 1) {
      const diffMinutes = Math.floor(diffMs / (1000 * 60))
      return `${diffMinutes} 分钟前`
    }
    return `${Math.floor(diffHours)} 小时前`
  } else {
    // 超过 24 小时，使用绝对时间
    return time.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}
```

---

### R6: WebSocket 重连策略

**Research Question**: WebSocket 连接断开后如何重连？

**Options Evaluated**:

1. **固定间隔重连（无限次）**
   - Pros: 简单
   - Cons: 可能导致服务器过载，无法处理永久性故障
   - Example: 每 5 秒重试一次

2. **指数退避（有限次）**
   - Pros: 避免服务器过载，给予足够重连机会
   - Cons: 实现稍复杂
   - Example: 1s, 2s, 4s, 8s, 16s, 30s（最大 30s，最多 5 次）

3. **立即重连（无延迟）**
   - Pros: 最快恢复连接
   - Cons: 可能导致服务器过载，浪费资源
   - Example: 连接断开立即重连

**Decision**: **指数退避（初始 1 秒，最大 30 秒，最多 5 次）**

**Rationale**:
- 业界标准做法，经过大规模验证
- 避免服务器过载，保护系统稳定性
- 给予足够重连机会，处理临时性网络波动
- 对永久性故障有合理的退出机制

**Implementation Notes**:
```typescript
class WebSocketClient {
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private baseDelay = 1000 // 1 秒
  private maxDelay = 30000 // 30 秒
  
  private calculateDelay(): number {
    const delay = Math.min(
      this.baseDelay * Math.pow(2, this.reconnectAttempts),
      this.maxDelay
    )
    return delay
  }
  
  private reconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('WebSocket 重连失败，已达到最大重试次数')
      this.showErrorNotification()
      return
    }
    
    const delay = this.calculateDelay()
    console.log(`WebSocket 将在 ${delay}ms 后重连（第 ${this.reconnectAttempts + 1} 次）`)
    
    setTimeout(() => {
      this.reconnectAttempts++
      this.connect()
    }, delay)
  }
}
```

---

### R7: 现有 API 复用确认

**Research Question**: 哪些现有 API 可以复用？数据格式是否满足需求？

**Assumptions** (需要实施前确认):

1. **`GET /api/v1/progress`** - 获取学习进度统计
   - Assumed Response:
   ```json
   {
     "code": 0,
     "data": {
       "total_chapters": 50,
       "completed_chapters": 20,
       "progress_percentage": 40.0,
       "topic_progress": [
         {
           "topic_id": "lexical-elements",
           "topic_name": "词法元素",
           "completed": 10,
           "total": 15,
           "percentage": 66.7
         }
       ]
     }
   }
   ```
   - **Action Required**: 确认此 API 是否存在，返回数据格式是否如上

2. **`GET /api/v1/quiz/history`** - 获取测验历史
   - Assumed Response:
   ```json
   {
     "code": 0,
     "data": {
       "records": [
         {
           "id": 123,
           "topic_name": "常量",
           "chapter_name": "布尔常量",
           "score": 85,
           "total_questions": 100,
           "completed_at": "2025-12-26T10:00:00+08:00"
         }
       ]
     }
   }
   ```
   - **Action Required**: 确认此 API 是否存在，是否支持 `limit` 参数

3. **`GET /api/v1/topics`** - 获取主题列表
   - Assumed Response:
   ```json
   {
     "code": 0,
     "data": {
       "topics": [
         {
           "id": "lexical-elements",
           "name": "Lexical Elements",
           "display_name": "词法元素",
           "chapter_count": 15
         }
       ]
     }
   }
   ```
   - **Action Required**: 确认此 API 是否存在，返回数据格式是否如上

**Decision**: 优先复用现有 API，如数据格式不符，考虑以下方案：
1. 修改现有 API（如果不影响其他功能）
2. 新增专用 API（如果修改影响较大）
3. 前端适配数据格式（如果差异较小）

---

## Technology Stack Validation

### Frontend

| Technology | Version | Status | Notes |
|------------|---------|--------|-------|
| Next.js | 14.2.15 | ✅ Confirmed | 已在项目中使用 |
| React | 18.x | ✅ Confirmed | Next.js 依赖 |
| TypeScript | 5.x | ✅ Confirmed | 已在项目中使用 |
| Ant Design | 5.x | ✅ Confirmed | 已在项目中使用 |
| WebSocket API | Native | ✅ Confirmed | 浏览器原生支持 |

### Backend

| Technology | Version | Status | Notes |
|------------|---------|--------|-------|
| Go | 1.24.5 | ✅ Confirmed | 已在项目中使用 |
| GoFrame | v2.9.5 | ✅ Confirmed | 已在项目中使用 |
| gorilla/websocket | latest | ⏳ To Add | 需要添加依赖 |

**Action Required**:
```bash
# 添加 gorilla/websocket 依赖
cd backend
go get github.com/gorilla/websocket
go mod tidy
```

---

## Performance Research

### WebSocket Performance

**Benchmark Results** (from gorilla/websocket documentation):
- Concurrent Connections: 10,000+
- Message Throughput: 100,000+ msg/s
- Latency: < 1ms (local), < 50ms (network)

**Conclusion**: gorilla/websocket 性能满足需求（预计用户规模 100-1000）

### Database Query Performance

**Test Query** (学习天数计算):
```sql
EXPLAIN ANALYZE
SELECT COUNT(DISTINCT DATE(last_visited_at)) 
FROM learning_progress 
WHERE user_id = 1 AND last_visited_at IS NOT NULL;
```

**Expected Performance**:
- Without Index: ~100ms (1000 records)
- With Index: ~10ms (1000 records)

**Recommendation**: 在 `learning_progress(user_id, last_visited_at)` 上创建复合索引

---

## Security Research

### WebSocket Security

**Best Practices**:
1. **Authentication**: 在连接建立时验证 token
2. **Authorization**: 只推送用户自己的数据
3. **Rate Limiting**: 限制连接频率和消息频率
4. **Input Validation**: 验证所有客户端消息
5. **Origin Check**: 验证 WebSocket 请求来源

**Implementation**:
```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        // 生产环境需要验证 Origin
        return origin == "https://yourdomain.com"
    },
}
```

---

## Open Questions for Implementation

以下问题需要在实施阶段确认：

1. ✅ **WebSocket 库选择**: 使用 gorilla/websocket
2. ✅ **学习天数计算**: 有学习活动的不同日期数
3. ✅ **最后学习记录获取**: 新增 `/api/v1/progress/last` 接口
4. ✅ **身份验证集成**: 复用现有机制
5. ✅ **时间格式化**: 混合格式（24h 内相对 + 超过绝对）
6. ✅ **WebSocket 重连**: 指数退避（1s-30s，最多 5 次）
7. ⏳ **现有 API 确认**: 需要实施前验证 API 存在性和数据格式

---

## Summary

所有关键技术决策已完成，研究结果如下：

| 研究主题 | 决策 | 状态 |
|---------|------|------|
| WebSocket 实现 | gorilla/websocket | ✅ 完成 |
| 学习天数计算 | 活动日期数 | ✅ 完成 |
| 最后学习记录 | 新增 API | ✅ 完成 |
| 身份验证 | 复用现有机制 | ✅ 完成 |
| 时间格式化 | 混合格式 | ✅ 完成 |
| WebSocket 重连 | 指数退避 | ✅ 完成 |
| 现有 API 复用 | 待确认 | ⏳ 待实施 |

**下一步**: 进入 Phase 1 设计阶段，生成数据模型和 API 契约。
