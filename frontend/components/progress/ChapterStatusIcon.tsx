"use client";

import {
  CheckCircleTwoTone,
  ClockCircleTwoTone,
  ExclamationCircleTwoTone,
  QuestionCircleTwoTone,
} from "@ant-design/icons";
import { Tooltip } from "antd";
import { ProgressStatus } from "@/types/learning";
import ProgressStatuses from "@/lib/progressStatus";

interface ChapterStatusIconProps {
  status: ProgressStatus;
}

function iconForStatus(status: ProgressStatus) {
  if (status === ProgressStatuses.Completed) {
    return <CheckCircleTwoTone twoToneColor="#52c41a" />;
  }
  if (status === ProgressStatuses.Tested) {
    return <ExclamationCircleTwoTone twoToneColor="#fa8c16" />;
  }
  if (status === ProgressStatuses.InProgress) {
    return <ClockCircleTwoTone twoToneColor="#1677ff" />;
  }
  return <QuestionCircleTwoTone twoToneColor="#bfbfbf" />;
}

function label(status: ProgressStatus) {
  if (status === ProgressStatuses.Completed) return "已完成";
  if (status === ProgressStatuses.Tested) return "已测验";
  if (status === ProgressStatuses.InProgress) return "学习中";
  return "未开始";
}

export default function ChapterStatusIcon({ status }: ChapterStatusIconProps) {
  return <Tooltip title={label(status)}>{iconForStatus(status)}</Tooltip>;
}

