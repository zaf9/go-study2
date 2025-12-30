/**
 * T120: QuizHistoryCard组件
 * 展示单个测验历史记录的卡片
 */

import { QuizHistoryItem } from "@/types/quiz";
import { Card, Tag, Typography, Space } from "antd";
import { 
  CheckCircleOutlined, 
  CloseCircleOutlined,
  TrophyOutlined 
} from "@ant-design/icons";
import { formatDistanceToNow } from "date-fns";
import { zhCN } from "date-fns/locale";

const { Text } = Typography;

interface QuizHistoryCardProps {
  item: QuizHistoryItem;
  onClick?: (sessionId: string) => void;
}

export default function QuizHistoryCard({ item, onClick }: QuizHistoryCardProps) {
  // 计算正确率
  const correctnessRate = (item.totalQuestions ?? 0) > 0 
    ? Math.round(((item.correctAnswers ?? 0) / (item.totalQuestions ?? 1)) * 100) 
    : 0;

  // 判断是否完美通过
  const isPerfect = item.score === 100;

  // 格式化时间
  const timeAgo = item.completedAt 
    ? formatDistanceToNow(new Date(item.completedAt), { 
        addSuffix: true, 
        locale: zhCN 
      })
    : "";

  return (
    <Card
      hoverable
      className="quiz-history-card cursor-pointer"
      onClick={() => item.sessionId && onClick?.(item.sessionId)}
      data-testid="quiz-history-card"
    >
      <Space direction="vertical" size="small" className="w-full">
        {/* 标题行: 章节名称 + 主题标签 */}
        <div className="flex justify-between items-center">
          <Text strong className="text-base">
            {item.chapter || "测验"}
          </Text>
          <Tag color="blue">{item.topic}</Tag>
        </div>

        {/* 分数和正确率 */}
        <div className="flex gap-4">
          <Space>
            <TrophyOutlined className="text-yellow-500" />
            <Text>{item.score}分</Text>
          </Space>
          <Space>
            <Text type="secondary">正确率:</Text>
            <Text>{correctnessRate}%</Text>
          </Space>
        </div>

        {/* 题目数量 */}
        <div>
          <Text type="secondary">
            {item.correctAnswers ?? 0}/{item.totalQuestions ?? 0}题
          </Text>
        </div>

        {/* 通过状态 */}
        <div className="flex justify-between items-center">
          {item.passed ? (
            <Tag 
              icon={<CheckCircleOutlined />} 
              color="success"
              className="quiz-passed-badge"
            >
              {isPerfect ? "完美通过" : "已通过"}
            </Tag>
          ) : (
            <Tag 
              icon={<CloseCircleOutlined />} 
              color="error"
              className="quiz-failed-badge"
            >
              未通过
            </Tag>
          )}

          {/* 完成时间 */}
          {timeAgo && (
            <Text type="secondary" className="text-sm">
              {timeAgo}
            </Text>
          )}
        </div>
      </Space>
    </Card>
  );
}
