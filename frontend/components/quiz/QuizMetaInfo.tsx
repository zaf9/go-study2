"use client";

import { Card, Space, Statistic, Tag, Typography } from "antd";
import { ClockCircleOutlined, QuestionCircleOutlined } from "@ant-design/icons";
import { QuizStats } from "@/types/quiz";

const { Text } = Typography;

interface QuizMetaInfoProps {
  stats: QuizStats | null;
  loading?: boolean;
}

/**
 * QuizMetaInfo 组件用于显示题库元数据信息
 * 包括总题量、预计用时、难度分布
 */
export default function QuizMetaInfo({ stats, loading }: QuizMetaInfoProps) {
  if (loading || !stats) {
    return (
      <Card>
        <Space direction="vertical" className="w-full">
          <Text type="secondary">加载中...</Text>
        </Space>
      </Card>
    );
  }

  // 计算预计用时（假设每题约 1-2 分钟，根据难度调整）
  const calculateEstimatedTime = (stats: QuizStats): number => {
    const easyTime = (stats.byDifficulty.easy || 0) * 1;
    const mediumTime = (stats.byDifficulty.medium || 0) * 1.5;
    const hardTime = (stats.byDifficulty.hard || 0) * 2;
    return Math.ceil(easyTime + mediumTime + hardTime);
  };

  const estimatedMinutes = calculateEstimatedTime(stats);

  // 难度标签映射
  const difficultyLabels: Record<string, string> = {
    easy: "简单",
    medium: "中等",
    hard: "困难",
  };

  // 题型标签映射
  const typeLabels: Record<string, string> = {
    single: "单选题",
    multiple: "多选题",
    truefalse: "判断题",
    code_output: "代码输出",
    code_correction: "代码改错",
  };

  return (
    <Card title="题库信息" className="mb-4">
      <Space direction="vertical" className="w-full" size="middle">
        <Space wrap>
          <Statistic
            title="总题量"
            value={stats.total}
            prefix={<QuestionCircleOutlined />}
            valueStyle={{ fontSize: "20px" }}
          />
          <Statistic
            title="预计用时"
            value={estimatedMinutes}
            prefix={<ClockCircleOutlined />}
            suffix="分钟"
            valueStyle={{ fontSize: "20px" }}
          />
        </Space>

        {Object.keys(stats.byDifficulty).length > 0 && (
          <div>
            <Text strong className="block mb-2">
              难度分布:
            </Text>
            <Space wrap>
              {Object.entries(stats.byDifficulty).map(([difficulty, count]) => (
                <Tag key={difficulty} color={getDifficultyColor(difficulty)}>
                  {difficultyLabels[difficulty] || difficulty}: {count} 题
                </Tag>
              ))}
            </Space>
          </div>
        )}

        {Object.keys(stats.byType).length > 0 && (
          <div>
            <Text strong className="block mb-2">
              题型分布:
            </Text>
            <Space wrap>
              {Object.entries(stats.byType).map(([type, count]) => (
                <Tag key={type} color="blue">
                  {typeLabels[type] || type}: {count} 题
                </Tag>
              ))}
            </Space>
          </div>
        )}
      </Space>
    </Card>
  );
}

/**
 * 根据难度返回对应的颜色
 */
function getDifficultyColor(difficulty: string): string {
  switch (difficulty) {
    case "easy":
      return "green";
    case "medium":
      return "orange";
    case "hard":
      return "red";
    default:
      return "default";
  }
}

