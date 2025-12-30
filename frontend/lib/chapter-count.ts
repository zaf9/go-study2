/**
 * 章节数量计算工具
 * 从 static-routes.ts 获取准确的章节数量
 */

import { topicChapters } from './static-routes';
import { TopicKey } from '@/types/learning';

/**
 * 获取指定主题的章节数量
 * @param topic 主题标识符
 * @returns 章节数量,如果主题不存在则返回 0
 */
export function getChapterCount(topic: TopicKey): number {
  const chapters = topicChapters[topic];
  return chapters ? chapters.length : 0;
}

/**
 * 获取所有主题的章节数量
 * @returns 每个主题的章节数量映射
 */
export function getAllChapterCounts(): Record<TopicKey, number> {
  const counts: Partial<Record<TopicKey, number>> = {};
  
  for (const topic in topicChapters) {
    counts[topic as TopicKey] = topicChapters[topic as TopicKey].length;
  }
  
  return counts as Record<TopicKey, number>;
}

/**
 * 获取所有章节的总数
 * @returns 所有章节的总数
 */
export function getTotalChapterCount(): number {
  return Object.values(topicChapters).reduce(
    (total, chapters) => total + chapters.length,
    0
  );
}
