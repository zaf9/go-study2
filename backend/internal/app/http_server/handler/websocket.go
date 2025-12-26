package handler

import (
	"log"
	"net/http"

	"go-study2/internal/app/http_server/middleware"
	"go-study2/internal/websocket"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
)

/**
 * WebSocket 升级器
 * 将 HTTP 连接升级为 WebSocket 连接
 */
var upgrader = websocket.Upgrader{
	// 检查请求来源，生产环境应该验证 Origin
	CheckOrigin: func(r *http.Request) bool {
		// 开发环境允许所有来源
		// 生产环境需要验证 Origin
		origin := r.Header.Get("Origin")
		return origin != "" // 简化处理，生产环境需要更严格的验证
	},

	// 读取缓冲区大小
	ReadBufferSize: 1024,

	// 写入缓冲区大小
	WriteBufferSize: 1024,
}

/**
 * 全局 WebSocket Hub
 */
var wsHub *websocket.Hub

/**
 * 初始化 WebSocket Hub
 * 应该在应用启动时调用
 */
func InitWebSocketHub() {
	wsHub = websocket.NewHub()
	go wsHub.Run()
	log.Println("[WebSocket] Hub 已启动")
}

/**
 * GetWebSocketHub 获取 WebSocket Hub 实例
 */
func GetWebSocketHub() *websocket.Hub {
	return wsHub
}

/**
 * HandleWebSocket 处理 WebSocket 连接请求
 * URL: /api/v1/ws/dashboard
 * 需要 JWT 认证
 */
func (h *Handler) HandleWebSocket(r *ghttp.Request) {
	// 从上下文中获取用户 ID（由 Auth 中间件设置）
	userID, exists := r.GetCtxVar("user_id").(uint)
	if !exists {
		log.Println("[WebSocket] 未认证的连接请求")
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}

	// 清除响应缓冲区（防止 Upgrade 失败时发送 JSON 响应）
	r.Response.ClearBuffer()

	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(r.Response.Writer.ResponseWriter, r.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] 升级连接失败: 用户ID=%d, 错误=%v", userID, err)
		r.Response.WriteStatus(http.StatusInternalServerError)
		return
	}

	// 创建 WebSocket 客户端
	client := websocket.NewClient(wsHub, userID, conn)

	// 注册客户端到 Hub
	wsHub.RegisterClient(client)

	// 启动读写 pump
	go client.WritePump()
	go client.ReadPump()

	log.Printf("[WebSocket] 新连接已建立: 用户ID=%d", userID)
}

