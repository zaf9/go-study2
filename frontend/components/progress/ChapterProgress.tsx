"use client";

import { Alert, Button, Space, Typography } from "antd";
import ProgressBar from "./ProgressBar";
import { ProgressStatus } from "@/types/learning";

interface ChapterProgressProps {
  title: string;
  status: ProgressStatus;
  scrollProgress?: number;
  readDuration?: number;
  estimatedSeconds?: number;
  lastVisitAt?: string;
  onResume?: () => void;
}

const { Text } = Typography;

export default function ChapterProgress({
  title,
  status,
  scrollProgress = 0,
  readDuration = 0,
  estimatedSeconds = 0,
  lastVisitAt,
  onResume,
}: ChapterProgressProps) {
  const remaining =
    estimatedSeconds > 0
      ? Math.max(0, estimatedSeconds - readDuration)
      : undefined;

  return (
    <Space direction="vertical" className="w-full">
      <div className="flex items-center justify-between">
        <div>
          <Text strong>{title}</Text>
          {lastVisitAt && (
            <Text type="secondary" className="ml-3">
              最近访问：{new Date(lastVisitAt).toLocaleString()}
            </Text>
          )}
        </div>
        {onResume && (
          <Button type="link" onClick={onResume}>
            恢复到上次阅读位置
          </Button>
        )}
      </div>
      <ProgressBar
        status={status}
        percent={scrollProgress}
        segments={10}
        label="阅读进度"
      />
      <div className="flex items-center justify-between">
        <Text type="secondary">
          已阅读 {Math.round(readDuration)} 秒
          {remaining !== undefined ? ` / 预计 ${Math.round(estimatedSeconds)} 秒` : ""}
        </Text>
        {remaining !== undefined && (
          <Text type="success">预计剩余 {Math.round(remaining)} 秒</Text>
        )}
      </div>
      {status === "tested" && (
        <Alert
          type="info"
          showIcon
          message="已完成测验，建议补充阅读未完成段落。"
        />
      )}
    </Space>
  );
}

