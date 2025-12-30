/**
 * 章节状态计算工具
 * 根据学习进度数据计算章节的学习状态
 */

export type ChapterStatus = 'not_started' | 'in_progress' | 'completed';

export interface ChapterProgressData {
  status?: string;
  completed_at?: string | null;
  first_visit_at?: string | null;
}

/**
 * 计算章节的学习状态
 * @param progressData 章节进度数据,如果为 null 表示未开始学习
 * @returns 章节状态: 未开始/学习中/已完成
 */
export function getChapterStatus(
  progressData: ChapterProgressData | null | undefined
): ChapterStatus {
  // 如果没有进度数据,表示未开始
  if (!progressData) {
    return 'not_started';
  }

  // 如果有 completed_at 或状态为 completed,表示已完成
  if (progressData.completed_at || progressData.status === 'completed') {
    return 'completed';
  }

  // 如果有 first_visit_at 或状态为 in_progress,表示学习中
  if (progressData.first_visit_at || progressData.status === 'in_progress') {
    return 'in_progress';
  }

  // 默认为未开始
  return 'not_started';
}

/**
 * 获取章节状态的显示文本
 * @param status 章节状态
 * @returns 中文显示文本
 */
export function getChapterStatusText(status: ChapterStatus): string {
  const textMap: Record<ChapterStatus, string> = {
    not_started: '未开始',
    in_progress: '学习中',
    completed: '已完成',
  };
  return textMap[status];
}

/**
 * 获取章节状态的图标
 * @param status 章节状态
 * @returns 图标emoji
 */
export function getChapterStatusIcon(status: ChapterStatus): string {
  const iconMap: Record<ChapterStatus, string> = {
    not_started: '🔘',
    in_progress: '📖',
    completed: '✅',
  };
  return iconMap[status];
}

/**
 * 获取章节状态的颜色
 * @param status 章节状态
 * @returns Ant Design 颜色名称
 */
export function getChapterStatusColor(
  status: ChapterStatus
): 'default' | 'processing' | 'success' {
  const colorMap: Record<ChapterStatus, 'default' | 'processing' | 'success'> = {
    not_started: 'default',
    in_progress: 'processing',
    completed: 'success',
  };
  return colorMap[status];
}
