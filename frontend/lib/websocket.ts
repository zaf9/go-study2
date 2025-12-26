/**
 * WebSocket 客户端封装
 *
 * 提供带重连机制的 WebSocket 连接管理
 */

import {
  WebSocketMessage,
  ProgressUpdatedEventData,
  QuizCompletedEventData,
} from '@/types/dashboard'

/**
 * WebSocket 客户端配置选项
 */
export interface WebSocketClientOptions {
  /** WebSocket 服务器 URL */
  url: string

  /** 认证 token */
  token?: string

  /** 连接打开回调 */
  onOpen?: () => void

  /** 连接关闭回调 */
  onClose?: (event: CloseEvent) => void

  /** 错误回调 */
  onError?: (error: Event) => void

  /** 消息回调 */
  onMessage?: (message: WebSocketMessage) => void

  /** 最大重连次数 */
  maxReconnectAttempts?: number

  /** 初始重连延迟（毫秒） */
  baseDelay?: number

  /** 最大重连延迟（毫秒） */
  maxDelay?: number

  /** 是否自动重连 */
  autoReconnect?: boolean
}

/**
 * WebSocket 客户端类
 */
export class WebSocketClient {
  private ws: WebSocket | null = null

  private options: Required<WebSocketClientOptions>

  private reconnectAttempts = 0

  private reconnectTimeoutId: number | null = null

  private isManualClose = false

  constructor(options: WebSocketClientOptions) {
    this.options = {
      url: options.url,
      token: options.token,
      onOpen: options.onOpen || (() => {}),
      onClose: options.onClose || (() => {}),
      onError: options.onError || (() => {}),
      onMessage: options.onMessage || (() => {}),
      maxReconnectAttempts: options.maxReconnectAttempts || 5,
      baseDelay: options.baseDelay || 1000, // 1 秒
      maxDelay: options.maxDelay || 30000, // 30 秒
      autoReconnect: options.autoReconnect !== false,
    }
  }

  /**
   * 连接到 WebSocket 服务器
   */
  public connect(): void {
    if (this.ws && (this.ws.readyState === WebSocket.CONNECTING || this.ws.readyState === WebSocket.OPEN)) {
      console.warn('WebSocket 已经连接或正在连接中')
      return
    }

    try {
      // 构建带 token 的 URL
      const url = this.options.token
        ? `${this.options.url}?token=${encodeURIComponent(this.options.token)}`
        : this.options.url

      console.log(`WebSocket 正在连接到: ${this.options.url}`)

      this.ws = new WebSocket(url)

      this.ws.onopen = this.handleOpen.bind(this)

      this.ws.onclose = this.handleClose.bind(this)

      this.ws.onerror = this.handleError.bind(this)

      this.ws.onmessage = this.handleMessage.bind(this)
    } catch (error) {
      console.error('WebSocket 连接失败:', error)
      this.options.onError(error as Event)

      if (this.options.autoReconnect) {
        this.scheduleReconnect()
      }
    }
  }

  /**
   * 断开 WebSocket 连接
   */
  public disconnect(): void {
    this.isManualClose = true

    if (this.reconnectTimeoutId) {
      clearTimeout(this.reconnectTimeoutId)
      this.reconnectTimeoutId = null
    }

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  /**
   * 发送消息
   */
  public send(data: any): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket 未连接，无法发送消息')
      return
    }

    try {
      this.ws.send(JSON.stringify(data))
    } catch (error) {
      console.error('发送消息失败:', error)
    }
  }

  /**
   * 获取连接状态
   */
  public getReadyState(): number {
    return this.ws ? this.ws.readyState : WebSocket.CLOSED
  }

  /**
   * 是否已连接
   */
  public isConnected(): boolean {
    return this.ws ? this.ws.readyState === WebSocket.OPEN : false
  }

  /**
   * 处理连接打开事件
   */
  private handleOpen(): void {
    console.log('WebSocket 连接已打开')

    this.reconnectAttempts = 0

    this.options.onOpen()
  }

  /**
   * 处理连接关闭事件
   */
  private handleClose(event: CloseEvent): void {
    console.log(`WebSocket 连接已关闭: code=${event.code}, reason=${event.reason}`)

    this.ws = null

    this.options.onClose(event)

    // 如果不是手动关闭且启用了自动重连，则尝试重连
    if (!this.isManualClose && this.options.autoReconnect) {
      this.scheduleReconnect()
    }
  }

  /**
   * 处理错误事件
   */
  private handleError(event: Event): void {
    console.error('WebSocket 错误:', event)
    this.options.onError(event)
  }

  /**
   * 处理消息事件
   */
  private handleMessage(event: MessageEvent): void {
    try {
      const message = JSON.parse(event.data) as WebSocketMessage
      this.options.onMessage(message)
    } catch (error) {
      console.error('解析 WebSocket 消息失败:', error)
    }
  }

  /**
   * 计算重连延迟（指数退避）
   */
  private calculateDelay(): number {
    const delay = Math.min(
      this.options.baseDelay * Math.pow(2, this.reconnectAttempts),
      this.options.maxDelay
    )
    return delay
  }

  /**
   * 安排重连
   */
  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.options.maxReconnectAttempts) {
      console.error(`WebSocket 重连失败，已达到最大重试次数: ${this.options.maxReconnectAttempts}`)
      return
    }

    const delay = this.calculateDelay()

    console.log(`WebSocket 将在 ${delay}ms 后重连（第 ${this.reconnectAttempts + 1} 次）`)

    this.reconnectTimeoutId = setTimeout(() => {
      this.reconnectAttempts++
      this.connect()
    }, delay) as unknown as number
  }
}

/**
 * 创建 WebSocket 客户端实例
 */
export function createWebSocketClient(options: WebSocketClientOptions): WebSocketClient {
  return new WebSocketClient(options)
}

/**
 * 创建用于 Dashboard 的 WebSocket 客户端
 */
export function createDashboardWebSocket(
  token: string,
  handlers: {
    onOpen?: () => void
    onClose?: (event: CloseEvent) => void
    onError?: (error: Event) => void
    onProgressUpdated?: (data: ProgressUpdatedEventData) => void
    onQuizCompleted?: (data: QuizCompletedEventData) => void
  }
): WebSocketClient {
  return createWebSocketClient({
    url: `${window.location.origin}/api/v1/ws/dashboard`,
    token,
    autoReconnect: true,
    maxReconnectAttempts: 5,
    baseDelay: 1000,
    maxDelay: 30000,
    onOpen: handlers.onOpen,
    onClose: handlers.onClose,
    onError: handlers.onError,
    onMessage: (message) => {
      switch (message.event) {
        case 'progress_updated':
          if (handlers.onProgressUpdated) {
            handlers.onProgressUpdated(message.data as ProgressUpdatedEventData)
          }
          break
        case 'quiz_completed':
          if (handlers.onQuizCompleted) {
            handlers.onQuizCompleted(message.data as QuizCompletedEventData)
          }
          break
        default:
          console.warn('未知的 WebSocket 事件类型:', (message as any).event)
      }
    },
  })
}
