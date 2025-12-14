"use client";

import { Progress, Tag } from "antd";
import { ProgressStatus } from "@/types/learning";

interface ProgressBarProps {
  status: ProgressStatus;
  percent?: number;
}

const statusColor: Record<ProgressStatus, string> = {
  not_started: "default",
  in_progress: "blue",
  completed: "green",
  tested: "orange",
};

export default function ProgressBar({ status, percent }: ProgressBarProps) {
  const value =
    percent ??
    (status === "completed"
      ? 100
      : status === "tested"
        ? 70
        : status === "in_progress"
          ? 50
          : 0);
  // derive display status from percent when percent provided
  const displayStatus: ProgressStatus =
    typeof percent === "number"
      ? value >= 100
        ? "completed"
        : value > 0
          ? "in_progress"
          : "not_started"
      : status;
  return (
    <div className="flex items-center gap-3">
      <Tag color={statusColor[displayStatus]}>{statusLabel(displayStatus)}</Tag>
      <div className="flex-1">
        <Progress
          percent={value}
          size="small"
          status={displayStatus === "completed" ? "success" : "active"}
        />
      </div>
    </div>
  );
}

function statusLabel(status: ProgressStatus) {
  switch (status) {
    case "completed":
      return "已完成";
    case "tested":
      return "已测验";
    case "in_progress":
      return "学习中";
    default:
      return "未开始";
  }
}
