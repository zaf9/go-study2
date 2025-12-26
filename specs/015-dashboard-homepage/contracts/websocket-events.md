# WebSocket Contract: Dashboard Real-time Events

**Feature**: Dashboard 首页 - 实时数据推送  
**Date**: 2025-12-26  
**Version**: v1  
**Related**: [spec.md](../spec.md) | [data-model.md](../data-model.md)

## Overview

Dashboard 使用 WebSocket 实现实时数据推送，当用户完成学习或测验时，自动更新 Dashboard 显示的数据，无需手动刷新页面。

## WebSocket Endpoint

```
WS /api/v1/ws/dashboard
```

## Authentication

**Required**: Yes  
**Method**: JWT Token (通过 URL 参数或首次握手消息)  
**Connection URL**: `ws://localhost:8080/api/v1/ws/dashboard?token=<jwt_token>`

## Connection Lifecycle

### 1. Connection Establishment (建立连接)

**Client Request**:
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/dashboard?token=' + token)

ws.onopen = () => {
  console.log('WebSocket connected')
}
```

**Server Response**:
```json
{
  "type": "connection",
  "message": "连接成功",
  "timestamp": "2025-12-26T14:00:00+08:00"
}
```

### 2. Heartbeat (心跳保活)

**Client → Server** (每 30 秒):
```json
{
  "type": "ping"
}
```

**Server → Client**:
```json
{
  "type": "pong",
  "timestamp": "2025-12-26T14:00:30+08:00"
}
```

### 3. Connection Close (关闭连接)

**Client**:
```javascript
ws.close()
```

**Server** (主动关闭):
```json
{
  "type": "close",
  "reason": "服务器维护",
  "timestamp": "2025-12-26T14:00:00+08:00"
}
```

## Event Types

### Event 1: progress_updated (学习进度更新)

**触发时机**: 用户完成一个章节的学习

**Server → Client**:
```json
{
  "event": "progress_updated",
  "data": {
    "user_id": 123,
    "topic_id": "lexical-elements",
    "topic_display_name": "词法元素",
    "chapter_id": "identifiers",
    "chapter_display_name": "标识符",
    "completed": true,
    "timestamp": "2025-12-26T14:05:30+08:00"
  }
}
```

**Data Fields**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `event` | string | Yes | 事件类型，固定值 `progress_updated` |
| `data.user_id` | integer | Yes | 用户 ID |
| `data.topic_id` | string | Yes | 主题 ID |
| `data.topic_display_name` | string | Yes | 主题中文显示名称 |
| `data.chapter_id` | string | Yes | 章节 ID |
| `data.chapter_display_name` | string | Yes | 章节中文显示名称 |
| `data.completed` | boolean | Yes | 是否完成（true/false） |
| `data.timestamp` | string | Yes | 事件时间戳（ISO 8601 格式） |

**Frontend Handling**:
```typescript
ws.onmessage = (event) => {
  const message = JSON.parse(event.data)
  
  if (message.event === 'progress_updated') {
    // 更新统计卡片
    setStats(prev => ({
      ...prev,
      completedChapters: prev.completedChapters + (message.data.completed ? 1 : 0)
    }))
    
    // 更新主题进度
    setTopicProgress(prev => 
      prev.map(topic => 
        topic.topicId === message.data.topic_id
          ? { ...topic, completedChapters: topic.completedChapters + 1 }
          : topic
      )
    )
    
    // 显示通知
    notification.success({
      message: '学习进度已更新',
      description: `完成了《${message.data.topic_display_name}》的《${message.data.chapter_display_name}》`
    })
  }
}
```

---

### Event 2: quiz_completed (测验完成)

**触发时机**: 用户完成一次测验

**Server → Client**:
```json
{
  "event": "quiz_completed",
  "data": {
    "user_id": 123,
    "quiz_id": 456,
    "topic_id": "constants",
    "topic_display_name": "常量",
    "chapter_id": "boolean-constants",
    "chapter_display_name": "布尔常量",
    "score": 85,
    "total_questions": 100,
    "passed": true,
    "timestamp": "2025-12-26T14:10:00+08:00"
  }
}
```

**Data Fields**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `event` | string | Yes | 事件类型，固定值 `quiz_completed` |
| `data.user_id` | integer | Yes | 用户 ID |
| `data.quiz_id` | integer | Yes | 测验记录 ID |
| `data.topic_id` | string | Yes | 主题 ID |
| `data.topic_display_name` | string | Yes | 主题中文显示名称 |
| `data.chapter_id` | string | Yes | 章节 ID |
| `data.chapter_display_name` | string | Yes | 章节中文显示名称 |
| `data.score` | integer | Yes | 得分 |
| `data.total_questions` | integer | Yes | 总题数 |
| `data.passed` | boolean | Yes | 是否通过（≥60%） |
| `data.timestamp` | string | Yes | 事件时间戳（ISO 8601 格式） |

**Frontend Handling**:
```typescript
ws.onmessage = (event) => {
  const message = JSON.parse(event.data)
  
  if (message.event === 'quiz_completed') {
    // 更新最近测验列表（添加到顶部，保留最多5条）
    setRecentQuizzes(prev => [
      {
        id: message.data.quiz_id,
        topicName: message.data.topic_display_name,
        chapterName: message.data.chapter_display_name,
        score: message.data.score,
        totalQuestions: message.data.total_questions,
        completedAt: message.data.timestamp
      },
      ...prev.slice(0, 4)
    ])
    
    // 显示通知
    notification.info({
      message: '测验已完成',
      description: `《${message.data.topic_display_name}》测验得分: ${message.data.score}/${message.data.total_questions}`
    })
  }
}
```

---

## Error Handling

### Connection Errors

**Error Event**:
```json
{
  "type": "error",
  "code": "AUTH_FAILED",
  "message": "认证失败，请重新登录",
  "timestamp": "2025-12-26T14:00:00+08:00"
}
```

**Error Codes**:

| Code | Description | Action |
|------|-------------|--------|
| `AUTH_FAILED` | 认证失败 | 重新登录获取新 token |
| `INVALID_TOKEN` | Token 无效 | 重新登录 |
| `TOKEN_EXPIRED` | Token 过期 | 刷新 token 或重新登录 |
| `SERVER_ERROR` | 服务器内部错误 | 尝试重连 |
| `RATE_LIMIT` | 请求频率过高 | 延迟重连 |

### Reconnection Strategy (重连策略)

**指数退避算法**:
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
  
  private connect() {
    this.ws = new WebSocket(this.url)
    
    this.ws.onopen = () => {
      console.log('WebSocket 连接成功')
      this.reconnectAttempts = 0 // 重置重连计数
    }
    
    this.ws.onclose = (event) => {
      console.log('WebSocket 连接关闭', event.code, event.reason)
      this.reconnect()
    }
    
    this.ws.onerror = (error) => {
      console.error('WebSocket 错误', error)
    }
  }
}
```

---

## Implementation Notes

### Backend (Go)

```go
// WebSocket Hub (连接管理)
type Hub struct {
    clients    map[uint]*Client // userID -> Client
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

// Client (单个连接)
type Client struct {
    hub    *Hub
    conn   *websocket.Conn
    userID uint
    send   chan []byte
}

// 广播给特定用户的所有连接
func (h *Hub) BroadcastToUser(userID uint, message interface{}) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    client, ok := h.clients[userID]
    if !ok {
        return // 用户未连接
    }
    
    data, err := json.Marshal(message)
    if err != nil {
        return
    }
    
    select {
    case client.send <- data:
    default:
        // 发送失败，关闭连接
        close(client.send)
        delete(h.clients, userID)
    }
}

// 触发学习进度更新事件
func (s *ProgressService) UpdateProgress(userID uint, topicID, chapterID string, completed bool) error {
    // 1. 更新数据库
    err := s.db.Model(&LearningProgress{}).
        Where("user_id = ? AND topic_id = ? AND chapter_id = ?", userID, topicID, chapterID).
        Updates(map[string]interface{}{
            "completed":       completed,
            "last_visited_at": time.Now(),
        }).Error
    
    if err != nil {
        return err
    }
    
    // 2. 获取主题和章节名称
    var topic Topic
    var chapter Chapter
    s.db.First(&topic, "id = ?", topicID)
    s.db.First(&chapter, "id = ?", chapterID)
    
    // 3. 广播 WebSocket 事件
    event := map[string]interface{}{
        "event": "progress_updated",
        "data": map[string]interface{}{
            "user_id":              userID,
            "topic_id":             topicID,
            "topic_display_name":   topic.DisplayName,
            "chapter_id":           chapterID,
            "chapter_display_name": chapter.DisplayName,
            "completed":            completed,
            "timestamp":            time.Now().Format(time.RFC3339),
        },
    }
    
    websocketHub.BroadcastToUser(userID, event)
    
    return nil
}
```

### Frontend (TypeScript)

```typescript
// WebSocket Provider
export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [connected, setConnected] = useState(false)
  
  useEffect(() => {
    const token = getAuthToken()
    if (!token) return
    
    const client = new WebSocketClient(`ws://localhost:8080/api/v1/ws/dashboard?token=${token}`)
    
    client.onConnect(() => setConnected(true))
    client.onDisconnect(() => setConnected(false))
    client.onMessage((message) => {
      // 处理消息
      handleWebSocketMessage(message)
    })
    
    client.connect()
    setWs(client.ws)
    
    return () => client.disconnect()
  }, [])
  
  return (
    <WebSocketContext.Provider value={{ ws, connected }}>
      {children}
    </WebSocketContext.Provider>
  )
}
```

---

## Testing

### Backend Test

```go
func TestWebSocketHub(t *testing.T) {
    hub := NewHub()
    go hub.Run()
    
    // Test case 1: 用户连接
    t.Run("UserConnect", func(t *testing.T) {
        // Simulate WebSocket connection
        // Assert user added to hub
    })
    
    // Test case 2: 广播消息
    t.Run("BroadcastToUser", func(t *testing.T) {
        // Add mock client
        // Broadcast message
        // Assert message received
    })
    
    // Test case 3: 用户断开
    t.Run("UserDisconnect", func(t *testing.T) {
        // Disconnect client
        // Assert user removed from hub
    })
}
```

### Frontend Test

```typescript
describe('WebSocket Integration', () => {
  it('should connect and receive progress update', async () => {
    const mockWs = new MockWebSocket()
    
    mockWs.send({
      event: 'progress_updated',
      data: {
        user_id: 123,
        topic_id: 'lexical-elements',
        completed: true
      }
    })
    
    // Assert state updated
    expect(stats.completedChapters).toBe(1)
  })
})
```

---

## Performance

- **Connection Limit**: 每用户最多 1 个活跃连接
- **Message Size**: 单条消息 < 1KB
- **Latency**: 消息延迟 < 500ms
- **Heartbeat Interval**: 30 秒

## Security

- **Authentication**: 连接时验证 JWT token
- **Authorization**: 只推送用户自己的事件
- **Rate Limiting**: 限制连接频率（每用户每分钟最多 10 次连接尝试）
- **Message Validation**: 验证消息格式和内容

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| v1 | 2025-12-26 | 初始版本，支持 progress_updated 和 quiz_completed 事件 |
