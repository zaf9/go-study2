// WebSocket Hub
// 管理所有 WebSocket 连接，支持按用户 ID 广播消息

package websocket

import (
	"log"
	"sync"
)

/**
 * WebSocket Hub 结构
 * 管理客户端连接池和消息广播
 */
type Hub struct {
	// 所有注册的客户端，按用户 ID 分组
	// key: 用户 ID, value: 该用户的所有客户端连接
	clients map[uint][]*Client

	// 用于保护 clients map 的互斥锁
	mu sync.RWMutex

	// 用于客户端注册的通道
	register chan *Client

	// 用于客户端注销的通道
	unregister chan *Client

	// 用于广播消息的通道
	broadcast chan *BroadcastMessage
}

/**
 * 广播消息结构
 */
type BroadcastMessage struct {
	// 目标用户 ID，如果为 0 则广播给所有用户
	UserID uint

	// 要广播的 WebSocket 消息
	Message WebSocketMessage
}

/**
 * 创建新的 Hub
 */
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage),
	}
}

/**
 * 运行 Hub
 * 启动一个 goroutine 来处理注册、注销和广播
 */
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

/**
 * 处理客户端注册
 */
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 将客户端添加到对应用户的客户端列表中
	h.clients[client.UserID] = append(h.clients[client.UserID], client)

	log.Printf("[WS Hub] 客户端注册: 用户ID=%d, 总连接数=%d",
		client.UserID, len(h.clients[client.UserID]))
}

/**
 * 处理客户端注销
 */
func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userClients, exists := h.clients[client.UserID]
	if !exists {
		return
	}

	// 从该用户的客户端列表中移除当前客户端
	for i, c := range userClients {
		if c == client {
			h.clients[client.UserID] = append(userClients[:i], userClients[i+1:]...)
			break
		}
	}

	// 如果该用户没有其他客户端了，删除用户记录
	if len(h.clients[client.UserID]) == 0 {
		delete(h.clients, client.UserID)
	}

	log.Printf("[WS Hub] 客户端注销: 用户ID=%d, 剩余连接数=%d",
		client.UserID, len(h.clients[client.UserID]))
}

/**
 * 处理消息广播
 */
func (h *Hub) handleBroadcast(message *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 广播给所有用户
	if message.UserID == 0 {
		for _, clients := range h.clients {
			h.sendToClients(clients, message.Message)
		}
		return
	}

	// 广播给特定用户
	clients, exists := h.clients[message.UserID]
	if !exists {
		return
	}

	h.sendToClients(clients, message.Message)
}

/**
 * 发送消息给客户端列表
 */
func (h *Hub) sendToClients(clients []*Client, message WebSocketMessage) {
	// 序列化消息
	data, err := MarshalWebSocketMessage(message)
	if err != nil {
		log.Printf("[WS Hub] 消息序列化失败: %v", err)
		return
	}

	// 发送给所有客户端
	for _, client := range clients {
		client.Send(data)
	}
}

/**
 * 注册客户端
 */
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

/**
 * 注销客户端
 */
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

/**
 * 广播消息给所有用户
 */
func (h *Hub) Broadcast(message WebSocketMessage) {
	h.broadcast <- &BroadcastMessage{
		UserID:  0, // 0 表示广播给所有用户
		Message: message,
	}
}

/**
 * 广播消息给特定用户
 */
func (h *Hub) BroadcastToUser(userID uint, message WebSocketMessage) {
	h.broadcast <- &BroadcastMessage{
		UserID:  userID,
		Message: message,
	}
}

/**
 * 获取在线用户数量
 */
func (h *Hub) GetOnlineUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

/**
 * 获取特定用户的连接数
 */
func (h *Hub) GetConnectionCount(userID uint) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID])
}

