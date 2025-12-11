"use client";

import {
  CheckCircleTwoTone,
  ClockCircleTwoTone,
  ExclamationCircleTwoTone,
  QuestionCircleTwoTone,
} from "@ant-design/icons";
import { Tooltip } from "antd";
import { ProgressStatus } from "@/types/learning";

interface ChapterStatusIconProps {
  status: ProgressStatus;
}

function iconForStatus(status: ProgressStatus) {
  if (status === "completed") {
    return <CheckCircleTwoTone twoToneColor="#52c41a" />;
  }
  if (status === "tested") {
    return <ExclamationCircleTwoTone twoToneColor="#fa8c16" />;
  }
  if (status === "in_progress") {
    return <ClockCircleTwoTone twoToneColor="#1677ff" />;
  }
  return <QuestionCircleTwoTone twoToneColor="#bfbfbf" />;
}

function label(status: ProgressStatus) {
  if (status === "completed") return "已完成";
  if (status === "tested") return "已测验";
  if (status === "in_progress") return "学习中";
  return "未开始";
}

export default function ChapterStatusIcon({ status }: ChapterStatusIconProps) {
  return <Tooltip title={label(status)}>{iconForStatus(status)}</Tooltip>;
}

