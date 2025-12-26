// WebSocket Client
// 处理单个 WebSocket 客户端连接，负责消息读/写和心跳检测

package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

/**
 * WebSocket Client 结构
 * 代表一个 WebSocket 客户端连接
 */
type Client struct {
	// 对应的 Hub
	hub *Hub

	// 用户 ID
	UserID uint

	// WebSocket 连接
	conn *websocket.Conn

	// 发送消息的缓冲通道
	send chan []byte

	// 用于保护 conn 的互斥锁
	mu sync.Mutex
}

/**
 * Ping 间隔（秒）
 */
const pingPeriod = 54 * time.Second

/**
 * 写入超时时间（秒）
 */
const writeWait = 10 * time.Second

/**
 * 读取超时时间（秒）
 */
const readWait = 60 * time.Second

/**
 * 最大消息大小（字节）
 */
const maxMessageSize = 512

/**
 * 创建新的 WebSocket 客户端
 */
func NewClient(hub *Hub, userID uint, conn *websocket.Conn) *Client {
	return &Client{
		hub:    hub,
		UserID:  userID,
		conn:    conn,
		send:    make(chan []byte, 256),
	}
}

/**
 * 读取消息
 * 从 WebSocket 连接中读取消息
 */
func (c *Client) ReadPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	// 设置读取参数
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(readWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(readWait))
		return nil
	})

	for {
		// 读取消息
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("[WS Client] 读取错误: 用户ID=%d, 错误=%v",
					c.UserID, err)
			}
			break
		}

		// 客户端发送的消息可以在这里处理
		log.Printf("[WS Client] 收到消息: 用户ID=%d, 内容=%s",
			c.UserID, string(message))
	}
}

/**
 * 写入消息
 * 从发送通道中写入消息到 WebSocket 连接
 */
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// 发送消息
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("[WS Client] 获取写入器失败: 用户ID=%d, 错误=%v",
					c.UserID, err)
				return
			}

			w.Write(message)

			// 关闭写入器
			if err := w.Close(); err != nil {
				log.Printf("[WS Client] 关闭写入器失败: 用户ID=%d, 错误=%v",
					c.UserID, err)
				return
			}

		case <-ticker.C:
			// 发送 Ping 消息（心跳）
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[WS Client] 发送 Ping 失败: 用户ID=%d, 错误=%v",
					c.UserID, err)
				return
			}
		}
	}
}

/**
 * 发送消息
 * 非阻塞地将消息添加到发送队列
 */
func (c *Client) Send(message []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.send <- message:
		// 消息已加入队列
	default:
		// 发送通道已满，丢弃消息
		log.Printf("[WS Client] 发送通道已满，丢弃消息: 用户ID=%d", c.UserID)
	}
}

/**
 * 关闭连接
 */
func (c *Client) Close() {
	close(c.send)
}

/**
 * 序列化 WebSocket 消息为 JSON
 */
func MarshalWebSocketMessage(message WebSocketMessage) ([]byte, error) {
	return json.Marshal(message)
}

/**
 * 反序列化 JSON 为 WebSocket 消息
 */
func UnmarshalWebSocketMessage(data []byte) (WebSocketMessage, error) {
	var message WebSocketMessage
	err := json.Unmarshal(data, &message)
	return message, err
}

