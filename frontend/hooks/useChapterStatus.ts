/**
 * useChapterStatus Hook
 * 封装章节状态的获取和计算逻辑
 */

import { useMemo } from 'react';
import { ChapterProgress } from '@/types/learning';
import { getChapterStatus, ChapterStatus } from '@/lib/chapter-status';

export interface ChapterStatusMap {
  [chapterId: string]: ChapterStatus;
}

interface UseChapterStatusResult {
  /**
   * 章节状态映射表
   */
  statusMap: ChapterStatusMap;
  /**
   * 根据章节ID获取状态
   */
  getStatus: (chapterId: string) => ChapterStatus;
  /**
   * 检查章节是否已完成
   */
  isCompleted: (chapterId: string) => boolean;
  /**
   * 检查章节是否在学习中
   */
  isInProgress: (chapterId: string) => boolean;
  /**
   * 检查章节是否未开始
   */
  isNotStarted: (chapterId: string) => boolean;
}

/**
 * 章节状态管理 Hook
 * @param chaptersProgress 章节进度列表
 * @returns 章节状态相关的工具方法
 */
export default function useChapterStatus(
  chaptersProgress: ChapterProgress[] = []
): UseChapterStatusResult {
  // 构建章节状态映射表
  const statusMap = useMemo<ChapterStatusMap>(() => {
    const map: ChapterStatusMap = {};

    chaptersProgress.forEach((progress) => {
      map[progress.chapter] = getChapterStatus(progress);
    });

    return map;
  }, [chaptersProgress]);

  /**
   * 根据章节ID获取状态
   */
  const getStatus = (chapterId: string): ChapterStatus => {
    return statusMap[chapterId] || 'not_started';
  };

  /**
   * 检查章节是否已完成
   */
  const isCompleted = (chapterId: string): boolean => {
    return getStatus(chapterId) === 'completed';
  };

  /**
   * 检查章节是否在学习中
   */
  const isInProgress = (chapterId: string): boolean => {
    return getStatus(chapterId) === 'in_progress';
  };

  /**
   * 检查章节是否未开始
   */
  const isNotStarted = (chapterId: string): boolean => {
    return getStatus(chapterId) === 'not_started';
  };

  return {
    statusMap,
    getStatus,
    isCompleted,
    isInProgress,
    isNotStarted,
  };
}
