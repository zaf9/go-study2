/**
 * 进度计算工具函数
 *
 * 提供进度计算和相关的统计功能
 */

import { TopicProgressSummary } from '@/types/dashboard'

/**
 * 计算百分比
 *
 * @param completed - 已完成数量
 * @param total - 总数量
 * @param decimals - 保留小数位数，默认为 1
 * @returns 百分比值（如 50.5）
 */
export function calculatePercentage(
  completed: number,
  total: number,
  decimals: number = 1
): number {
  if (total === 0) {
    return 0
  }

  const percentage = (completed / total) * 100

  const multiplier = Math.pow(10, decimals)

  return Math.round(percentage * multiplier) / multiplier
}

/**
 * 判断进度是否为已完成（100%）
 *
 * @param percentage - 百分比值
 * @returns 是否已完成
 */
export function isCompleted(percentage: number): boolean {
  return percentage >= 100
}

/**
 * 判断进度是否为进行中（> 0 且 < 100%）
 *
 * @param percentage - 百分比值
 * @returns 是否进行中
 */
export function isInProgress(percentage: number): boolean {
  return percentage > 0 && percentage < 100
}

/**
 * 判断进度是否为未开始（0%）
 *
 * @param percentage - 百分比值
 * @returns 是否未开始
 */
export function isNotStarted(percentage: number): boolean {
  return percentage === 0
}

/**
 * 获取进度状态标签
 *
 * @param percentage - 百分比值
 * @returns 状态标签（"已完成"、"进行中"、"未开始"）
 */
export function getProgressStatusLabel(percentage: number): string {
  if (isCompleted(percentage)) {
    return '已完成'
  }

  if (isInProgress(percentage)) {
    return '进行中'
  }

  return '未开始'
}

/**
 * 获取进度状态颜色
 *
 * @param percentage - 百分比值
 * @returns 颜色类名（success、processing、default）
 */
export function getProgressStatusColor(percentage: number): string {
  if (isCompleted(percentage)) {
    return 'success'
  }

  if (isInProgress(percentage)) {
    return 'processing'
  }

  return 'default'
}

/**
 * 计算整体完成进度
 *
 * @param topicProgress - 主题进度列表
 * @returns 整体进度百分比
 */
export function calculateOverallProgress(
  topicProgress: TopicProgressSummary[]
): number {
  if (topicProgress.length === 0) {
    return 0
  }

  let totalChapters = 0

  let completedChapters = 0

  for (const topic of topicProgress) {
    totalChapters += topic.totalChapters
    completedChapters += topic.completedChapters
  }

  return calculatePercentage(completedChapters, totalChapters)
}

/**
 * 统计已完成的主题数量
 *
 * @param topicProgress - 主题进度列表
 * @returns 已完成的主题数量
 */
export function countCompletedTopics(topicProgress: TopicProgressSummary[]): number {
  return topicProgress.filter((topic) => isCompleted(topic.percentage)).length
}

/**
 * 统计进行中的主题数量
 *
 * @param topicProgress - 主题进度列表
 * @returns 进行中的主题数量
 */
export function countInProgressTopics(topicProgress: TopicProgressSummary[]): number {
  return topicProgress.filter((topic) => isInProgress(topic.percentage)).length
}

/**
 * 统计未开始的主题数量
 *
 * @param topicProgress - 主题进度列表
 * @returns 未开始的主题数量
 */
export function countNotStartedTopics(topicProgress: TopicProgressSummary[]): number {
  return topicProgress.filter((topic) => isNotStarted(topic.percentage)).length
}

/**
 * 按完成度排序主题
 *
 * @param topicProgress - 主题进度列表
 * @param order - 排序顺序（'asc' 升序，'desc' 降序）
 * @returns 排序后的主题进度列表
 */
export function sortTopicsByProgress(
  topicProgress: TopicProgressSummary[],
  order: 'asc' | 'desc' = 'desc'
): TopicProgressSummary[] {
  const sorted = [...topicProgress]

  sorted.sort((a, b) => {
    const diff = a.percentage - b.percentage
    return order === 'desc' ? -diff : diff
  })

  return sorted
}

/**
 * 按完成状态分组主题
 *
 * @param topicProgress - 主题进度列表
 * @returns 分组后的主题
 */
export function groupTopicsByStatus(topicProgress: TopicProgressSummary[]): {
  completed: TopicProgressSummary[]
  inProgress: TopicProgressSummary[]
  notStarted: TopicProgressSummary[]
} {
  return {
    completed: topicProgress.filter((topic) => isCompleted(topic.percentage)),
    inProgress: topicProgress.filter((topic) => isInProgress(topic.percentage)),
    notStarted: topicProgress.filter((topic) => isNotStarted(topic.percentage)),
  }
}

/**
 * 更新主题进度列表
 * 当收到进度更新事件时，更新对应主题的进度
 *
 * @param topicProgress - 当前主题进度列表
 * @param topicId - 要更新的主题 ID
 * @param completedChapters - 新的已完成章节数
 * @returns 更新后的主题进度列表
 */
export function updateTopicProgress(
  topicProgress: TopicProgressSummary[],
  topicId: string,
  completedChapters: number
): TopicProgressSummary[] {
  return topicProgress.map((topic) => {
    if (topic.topicId === topicId) {
      const newPercentage = calculatePercentage(
        completedChapters,
        topic.totalChapters
      )

      return {
        ...topic,
        completedChapters,
        percentage: newPercentage,
      }
    }

    return topic
  })
}

/**
 * 格式化进度百分比显示
 *
 * @param percentage - 百分比值
 * @param decimals - 保留小数位数，默认为 1
 * @returns 格式化后的字符串（如 "50.5%"）
 */
export function formatPercentage(
  percentage: number,
  decimals: number = 1
): string {
  return `${percentage.toFixed(decimals)}%`
}

/**
 * 计算本周活跃度
 * 根据本周完成章节数计算活跃度等级
 *
 * @param weeklyChaptersCompleted - 本周完成章节数
 * @returns 活跃度等级（"低"、"中"、"高"、"很高"）
 */
export function getActivityLevel(weeklyChaptersCompleted: number): string {
  if (weeklyChaptersCompleted === 0) {
    return '低'
  }

  if (weeklyChaptersCompleted < 5) {
    return '中'
  }

  if (weeklyChaptersCompleted < 10) {
    return '高'
  }

  return '很高'
}

/**
 * 获取活跃度颜色
 *
 * @param level - 活跃度等级
 * @returns 颜色类名
 */
export function getActivityLevelColor(level: string): string {
  switch (level) {
    case '低':
      return 'default'
    case '中':
      return 'processing'
    case '高':
      return 'warning'
    case '很高':
      return 'success'
    default:
      return 'default'
  }
}

