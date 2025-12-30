"use client";

import { Tag } from "antd";
import {
  getChapterStatus,
  getChapterStatusText,
  getChapterStatusIcon,
  getChapterStatusColor,
  ChapterProgressData,
} from "@/lib/chapter-status";

interface ChapterStatusBadgeProps {
  /**
   * 章节进度数据
   */
  progressData: ChapterProgressData | null | undefined;
  /**
   * 是否显示图标，默认为 true
   */
  showIcon?: boolean;
  /**
   * 是否显示文本，默认为 true
   */
  showText?: boolean;
}

/**
 * 章节状态徽章组件
 * 显示章节的学习状态（未开始/学习中/已完成）
 */
export default function ChapterStatusBadge({
  progressData,
  showIcon = true,
  showText = true,
}: ChapterStatusBadgeProps) {
  const status = getChapterStatus(progressData);
  const text = getChapterStatusText(status);
  const icon = getChapterStatusIcon(status);
  const color = getChapterStatusColor(status);

  return (
    <Tag color={color}>
      {showIcon && <span style={{ marginRight: showText ? 4 : 0 }}>{icon}</span>}
      {showText && text}
    </Tag>
  );
}
