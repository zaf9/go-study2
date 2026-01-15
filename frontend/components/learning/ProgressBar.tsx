"use client";

import { Progress, Tag, Tooltip } from "antd";
import { ProgressStatus } from "@/types/learning";
import ProgressStatuses from "@/lib/progressStatus";

interface ProgressBarProps {
  status: ProgressStatus;
  percent?: number;
  segments?: number;
  label?: string;
}

const statusColor: Record<ProgressStatus, string> = {
  not_started: "default",
  in_progress: "blue",
  completed: "green",
  tested: "orange",
};

function statusLabel(status: ProgressStatus): string {
  if (status === ProgressStatuses.Completed) return "已完成";
  if (status === ProgressStatuses.Tested) return "已测验";
  if (status === ProgressStatuses.InProgress) return "学习中";
  return "未开始";
}

export default function ProgressBar({
  status,
  percent,
  segments = 0,
  label,
}: ProgressBarProps) {
  const value =
    percent ??
    (status === ProgressStatuses.Completed
      ? 100
      : status === ProgressStatuses.Tested
        ? 70
        : status === ProgressStatuses.InProgress
          ? 50
          : 0);

  const capped = Math.min(100, Math.max(0, Math.round(value)));

  // Derive display status from percent when percent provided
  const displayStatus: ProgressStatus =
    typeof percent === "number"
      ? capped >= 100
        ? ProgressStatuses.Completed
        : capped > 0
          ? ProgressStatuses.InProgress
          : ProgressStatuses.NotStarted
      : status;

  const effectiveStatus = status ?? displayStatus;
  const progress = (
    <Progress
      percent={capped}
      steps={segments > 0 ? segments : undefined}
      showInfo
      size="small"
      status={
        effectiveStatus === ProgressStatuses.Completed ? "success" : "active"
      }
    />
  );

  return (
    <div className="flex items-center gap-3">
      {label && (
        <Tooltip title={label}>
          <Tag color={statusColor[effectiveStatus]}>
            {statusLabel(effectiveStatus)}
          </Tag>
        </Tooltip>
      )}
      {!label && (
        <Tag color={statusColor[effectiveStatus]}>
          {statusLabel(effectiveStatus)}
        </Tag>
      )}
      <div className="flex-1">{progress}</div>
    </div>
  );
}
