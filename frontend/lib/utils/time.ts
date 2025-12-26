/**
 * 时间格式化工具函数
 *
 * 提供时间显示的格式化功能
 */

/**
 * 时间格式化选项
 */
export interface TimeFormatOptions {
  /** 是否使用 24 小时制 */
  hour24?: boolean

  /** 是否显示秒 */
  includeSeconds?: boolean
}

/**
 * 格式化时间戳为相对时间或绝对时间
 * 24 小时内显示相对时间（如"2 小时前"），超过 24 小时显示绝对时间
 *
 * @param timestamp - ISO 8601 格式的时间戳
 * @param options - 格式化选项
 * @returns 格式化后的时间字符串
 */
export function formatTime(
  timestamp: string,
  options?: TimeFormatOptions
): string {
  const now = new Date()

  const time = new Date(timestamp)

  const diffMs = now.getTime() - time.getTime()

  const diffMinutes = diffMs / (1000 * 60)

  const diffHours = diffMinutes / 60

  const diffDays = diffMinutes / (60 * 24)

  // 24 小时内，使用相对时间
  if (diffHours < 24) {
    if (diffMinutes < 1) {
      return '刚刚'
    }

    if (diffMinutes < 60) {
      const minutes = Math.floor(diffMinutes)
      return `${minutes} 分钟前`
    }

    const hours = Math.floor(diffHours)
    return `${hours} 小时前`
  }

  // 超过 24 小时，使用绝对时间
  return formatAbsoluteTime(time, options)
}

/**
 * 格式化为绝对时间
 *
 * @param date - 日期对象
 * @param options - 格式化选项
 * @returns 格式化后的时间字符串
 */
export function formatAbsoluteTime(
  date: Date,
  options?: TimeFormatOptions
): string {
  const opts = {
    hour12: !(options?.hour24 ?? false),
    hour: '2-digit',
    minute: '2-digit',
    second: options?.includeSeconds ? '2-digit' : undefined,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  } as const

  return date.toLocaleString('zh-CN', opts)
}

/**
 * 格式化为日期（不含时间）
 *
 * @param timestamp - ISO 8601 格式的时间戳
 * @returns 格式化后的日期字符串（如"2025-12-26"）
 */
export function formatDate(timestamp: string): string {
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

/**
 * 格式化为时间（不含日期）
 *
 * @param timestamp - ISO 8601 格式的时间戳
 * @param options - 格式化选项
 * @returns 格式化后的时间字符串（如"10:30"）
 */
export function formatTimeOnly(
  timestamp: string,
  options?: TimeFormatOptions
): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', {
    hour12: !(options?.hour24 ?? false),
    hour: '2-digit',
    minute: '2-digit',
    second: options?.includeSeconds ? '2-digit' : undefined,
  })
}

/**
 * 计算时间差（以可读格式返回）
 *
 * @param timestamp - ISO 8601 格式的时间戳
 * @returns 时间差描述（如"2 天前"）
 */
export function getTimeAgo(timestamp: string): string {
  const now = new Date()

  const time = new Date(timestamp)

  const diffMs = now.getTime() - time.getTime()

  const diffSeconds = diffMs / 1000

  const diffMinutes = diffSeconds / 60

  const diffHours = diffMinutes / 60

  const diffDays = diffHours / 24

  const diffMonths = diffDays / 30

  const diffYears = diffDays / 365

  if (diffSeconds < 60) {
    return '刚刚'
  }

  if (diffMinutes < 60) {
    const minutes = Math.floor(diffMinutes)
    return `${minutes} 分钟前`
  }

  if (diffHours < 24) {
    const hours = Math.floor(diffHours)
    return `${hours} 小时前`
  }

  if (diffDays < 30) {
    const days = Math.floor(diffDays)
    return `${days} 天前`
  }

  if (diffMonths < 12) {
    const months = Math.floor(diffMonths)
    return `${months} 个月前`
  }

  const years = Math.floor(diffYears)
  return `${years} 年前`
}

/**
 * 判断时间是否为今天
 *
 * @param timestamp - ISO 8601 格式的时间戳
 * @returns 是否为今天
 */
export function isToday(timestamp: string): boolean {
  const date = new Date(timestamp)

  const today = new Date()

  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  )
}

/**
 * 判断时间是否为本周
 *
 * @param timestamp - ISO 8601 格式的时间戳
 * @returns 是否为本周
 */
export function isThisWeek(timestamp: string): boolean {
  const date = new Date(timestamp)

  const now = new Date()

  const weekStart = new Date(now)

  weekStart.setDate(now.getDate() - now.getDay()) // 本周周日（0）

  const weekEnd = new Date(weekStart)

  weekEnd.setDate(weekStart.getDate() + 6) // 本周周六

  return date >= weekStart && date <= weekEnd
}

/**
 * 获取友好的日期显示
 * 今天显示"今天"，本周显示"周X"，其他显示日期
 *
 * @param timestamp - ISO 8601 格式的时间戳
 * @returns 友好的日期显示
 */
export function getFriendlyDate(timestamp: string): string {
  if (isToday(timestamp)) {
    return '今天'
  }

  const date = new Date(timestamp)

  const now = new Date()

  // 检查是否是本周
  const weekStart = new Date(now)

  weekStart.setDate(now.getDate() - now.getDay())

  const weekEnd = new Date(weekStart)

  weekEnd.setDate(weekStart.getDate() + 6)

  if (date >= weekStart && date <= weekEnd) {
    const dayNames = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    return dayNames[date.getDay()]
  }

  return formatDate(timestamp)
}

