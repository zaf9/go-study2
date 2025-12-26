/**
 * Dashboard 相关类型定义
 *
 * 定义了 Dashboard 页面所需的所有数据结构类型
 */

/**
 * Dashboard 统计数据
 * 包含用户学习状态的汇总信息
 */
export interface DashboardStats {
  /** 累计学习天数（有学习活动的不同日期数） */
  studyDays: number

  /** 总章节数 */
  totalChapters: number

  /** 已完成章节数 */
  completedChapters: number

  /** 整体完成百分比（保留一位小数） */
  progressPercentage: number

  /** 本周学习活动次数（本周完成章节数） */
  weeklyActivity: number
}

/**
 * 最后学习记录
 * 用户最后一次学习的主题和章节信息
 */
export interface LastLearningRecord {
  /** 主题 ID */
  topicId: string

  /** 主题名称 */
  topicName: string

  /** 主题显示名称（中文） */
  topicDisplayName: string

  /** 章节 ID */
  chapterId: string

  /** 章节名称 */
  chapterName: string

  /** 章节显示名称（中文） */
  chapterDisplayName: string

  /** 最后访问时间（ISO 8601 格式） */
  lastVisitedAt: string
}

/**
 * 主题进度汇总
 * 每个主题的学习进度统计信息
 */
export interface TopicProgressSummary {
  /** 主题 ID */
  topicId: string

  /** 主题名称 */
  topicName: string

  /** 主题显示名称（中文） */
  displayName: string

  /** 已完成章节数 */
  completedChapters: number

  /** 总章节数 */
  totalChapters: number

  /** 完成百分比（保留一位小数） */
  percentage: number
}

/**
 * 最近测验汇总
 * 最近测验记录的展示数据
 */
export interface RecentQuizSummary {
  /** 测验记录 ID */
  id: number

  /** 主题名称 */
  topicName: string

  /** 章节名称 */
  chapterName: string

  /** 得分 */
  score: number

  /** 总题数 */
  totalQuestions: number

  /** 是否通过 */
  passed: boolean

  /** 完成时间（ISO 8601 格式） */
  completedAt: string

  /** 格式化后的时间（相对或绝对） */
  displayTime?: string
}

/**
 * WebSocket 进度更新事件数据
 */
export interface ProgressUpdatedEventData {
  /** 用户 ID */
  userId: number

  /** 主题 ID */
  topicId: string

  /** 章节 ID */
  chapterId: string

  /** 是否完成 */
  completed: boolean

  /** 时间戳（ISO 8601 格式） */
  timestamp: string
}

/**
 * WebSocket 测验完成事件数据
 */
export interface QuizCompletedEventData {
  /** 用户 ID */
  userId: number

  /** 测验 ID */
  quizId: number

  /** 得分 */
  score: number

  /** 总题数 */
  totalQuestions: number

  /** 是否通过 */
  passed: boolean

  /** 时间戳（ISO 8601 格式） */
  timestamp: string
}

/**
 * WebSocket 事件类型
 */
export type WebSocketEventType = 'progress_updated' | 'quiz_completed'

/**
 * WebSocket 消息结构
 */
export interface WebSocketMessage<T = any> {
  /** 事件类型 */
  event: WebSocketEventType

  /** 事件数据 */
  data: T
}

/**
 * Dashboard 页面数据
 * 服务端渲染时获取的初始数据
 */
export interface DashboardData {
  /** 用户信息 */
  user: {
    id: number
    username: string
    email: string
  }

  /** 统计数据 */
  stats: DashboardStats

  /** 最后学习记录（可能为 null） */
  lastLearning: LastLearningRecord | null

  /** 主题进度汇总列表 */
  topicProgress: TopicProgressSummary[]

  /** 最近测验记录列表（最多 5 条） */
  recentQuizzes: RecentQuizSummary[]
}

/**
 * Dashboard 错误状态
 */
export interface DashboardError {
  /** 错误消息 */
  message: string

  /** 错误类型 */
  type: 'network' | 'server' | 'data' | 'unknown'

  /** 是否可重试 */
  retryable: boolean
}

