'use client'

/**
 * WebSocket Provider
 * 提供全局的 WebSocket 连接和事件处理
 */

import React, {
	createContext,
	useContext,
	useState,
	useEffect,
	useCallback,
	type ReactNode,
} from 'react'
import {
	WebSocketClient,
	type ProgressUpdatedEventData,
	type QuizCompletedEventData,
} from '@/lib/websocket'

/**
 * WebSocket 上下文类型
 */
interface WebSocketContextValue {
	// WebSocket 客户端实例
	client: WebSocketClient | null

	// 是否已连接
	isConnected: boolean

	// 是否在重连
	isReconnecting: boolean

	// 连接错误
	error: Error | null

	// 连接方法
	connect: () => void

	// 断开连接方法
	disconnect: () => void
}

// 创建 WebSocket 上下文
const WebSocketContext = createContext<WebSocketContextValue | null>(null)

/**
 * Provider Props
 */
interface WebSocketProviderProps {
	children: ReactNode
}

/**
 * WebSocket Provider 组件
 * 管理 WebSocket 连接的生命周期
 */
export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({ children }) => {
	const [client, setClient] = useState<WebSocketClient | null>(null)
	const [isConnected, setIsConnected] = useState(false)
	const [isReconnecting, setIsReconnecting] = useState(false)
	const [error, setError] = useState<Error | null>(null)

	// 连接 WebSocket
	const connect = useCallback(() => {
		// 从 localStorage 获取 token
		const token = localStorage.getItem('token')
		if (!token) {
			setError(new Error('未找到认证令牌'))
			return
		}

		// 获取当前 API 基础 URL
		const apiBaseUrl = window.location.origin

		// 创建 WebSocket 客户端
		const wsClient = new WebSocketClient({
			url: `${apiBaseUrl}/api/v1/ws/dashboard`,
			token,
			autoReconnect: true,
			maxReconnectAttempts: 5,
			baseDelay: 1000,
			maxDelay: 30000,

			// 连接打开回调
			onOpen: () => {
				console.log('[WebSocket] 连接已打开')
				setIsConnected(true)
				setIsReconnecting(false)
				setError(null)
			},

			// 连接关闭回调
			onClose: (event) => {
				console.log('[WebSocket] 连接已关闭:', event)
				setIsConnected(false)
				if (event.code !== 1000 && client?.isConnected()) {
					// 非正常关闭，可能是重连中
					setIsReconnecting(true)
				}
			},

			// 错误回调
			onError: (event) => {
				console.error('[WebSocket] 连接错误:', event)
				setError(new Error('WebSocket 连接失败'))
				setIsConnected(false)
			},

			// 消息回调
			onMessage: (message) => {
				console.log('[WebSocket] 收到消息:', message)
				// 消息处理逻辑可以由使用 Provider 的组件实现
			},
		})

		// 连接
		wsClient.connect()

		// 保存客户端实例
		setClient(wsClient)
	}, [])

	// 断开 WebSocket 连接
	const disconnect = useCallback(() => {
		if (client) {
			client.disconnect()
			setClient(null)
			setIsConnected(false)
			setError(null)
		}
	}, [client])

	// 组件挂载时自动连接
	useEffect(() => {
		connect()

		// 组件卸载时断开连接
		return () => {
			disconnect()
		}
	}, [connect, disconnect])

	// 上下文值
	const contextValue: WebSocketContextValue = {
		client,
		isConnected,
		isReconnecting,
		error,
		connect,
		disconnect,
	}

	return (
		<WebSocketContext.Provider value={contextValue}>
			{children}
		</WebSocketContext.Provider>
	)
}

/**
 * 使用 WebSocket 上下文
 */
export const useWebSocket = (): WebSocketContextValue => {
	const context = useContext(WebSocketContext)
	if (!context) {
		throw new Error('useWebSocket 必须在 WebSocketProvider 内部使用')
	}
	return context
}

/**
 * 使用 WebSocket 进度更新事件
 */
export const useProgressUpdated = (
	callback: (data: ProgressUpdatedEventData) => void
) => {
	const { client } = useWebSocket()

	useEffect(() => {
		if (!client) {
			return
		}

		// 监听进度更新事件
		const handler = (message: any) => {
			if (message.event === 'progress_updated') {
				callback(message.data as ProgressUpdatedEventData)
			}
		}

		// 这里需要扩展 WebSocketClient 来支持添加事件监听器
		// 暂时留空，后续完善

		return () => {
			// 清理监听器
		}
	}, [client, callback])
}

/**
 * 使用 WebSocket 测验完成事件
 */
export const useQuizCompleted = (
	callback: (data: QuizCompletedEventData) => void
) => {
	const { client } = useWebSocket()

	useEffect(() => {
		if (!client) {
			return
		}

		// 监听测验完成事件
		const handler = (message: any) => {
			if (message.event === 'quiz_completed') {
				callback(message.data as QuizCompletedEventData)
			}
		}

		// 这里需要扩展 WebSocketClient 来支持添加事件监听器
		// 暂时留空，后续完善

		return () => {
			// 清理监听器
		}
	}, [client, callback])
}

