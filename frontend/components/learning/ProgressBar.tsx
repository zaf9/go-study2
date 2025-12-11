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
  done: "green",
};

export default function ProgressBar({ status, percent }: ProgressBarProps) {
  const value =
    percent ?? (status === "done" ? 100 : status === "in_progress" ? 50 : 0);
  return (
    <div className="flex items-center gap-3">
      <Tag color={statusColor[status]}>{statusLabel(status)}</Tag>
      <div className="flex-1">
        <Progress
          percent={value}
          size="small"
          status={status === "done" ? "success" : "active"}
        />
      </div>
    </div>
  );
}

function statusLabel(status: ProgressStatus) {
  switch (status) {
    case "done":
      return "已完成";
    case "in_progress":
      return "学习中";
    default:
      return "未开始";
  }
}
