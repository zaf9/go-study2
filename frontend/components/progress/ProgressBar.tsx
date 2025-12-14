"use client";

import { Progress, Tag, Tooltip } from "antd";
import { ProgressStatus } from "@/types/learning";
import ProgressStatuses from "@/lib/progressStatus";

interface ProgressBarProps {
  percent: number;
  status?: ProgressStatus;
  segments?: number;
  label?: string;
}

const statusColor: Record<ProgressStatus, string> = {
  not_started: "default",
  in_progress: "blue",
  completed: "green",
  tested: "orange",
};

function statusLabel(status?: ProgressStatus) {
  if (status === ProgressStatuses.Completed) return "已完成";
  if (status === ProgressStatuses.Tested) return "已测验";
  if (status === ProgressStatuses.InProgress) return "学习中";
  return "未开始";
}

export default function ProgressBar({
  percent,
  status = ProgressStatuses.InProgress,
  segments = 0,
  label,
}: ProgressBarProps) {
  const capped = Math.min(100, Math.max(0, Math.round(percent)));
  // derive display status from percent when there's a mismatch
  const displayStatus: ProgressStatus =
    capped >= 100
      ? ProgressStatuses.Completed
      : capped > 0
      ? ProgressStatuses.InProgress
      : ProgressStatuses.NotStarted;
  const effectiveStatus = status ?? displayStatus;
  const progress = (
    <Progress
      percent={capped}
      steps={segments > 0 ? segments : undefined}
      showInfo
      size="small"
      status={effectiveStatus === ProgressStatuses.Completed ? "success" : "active"}
    />
  );

  return (
    <div className="flex items-center gap-3">
      <Tooltip title={label}>
        <Tag color={statusColor[effectiveStatus]}>{statusLabel(effectiveStatus)}</Tag>
      </Tooltip>
      <div className="flex-1">{progress}</div>
    </div>
  );
}

