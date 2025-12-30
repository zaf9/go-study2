"use client";

import { Card, Space, Tag, Typography } from "antd";
import { useRouter } from "next/navigation";
import { TopicSummary } from "@/types/learning";
import { memo } from "react";

interface TopicCardProps {
  topic: TopicSummary;
}

const { Title, Paragraph, Text } = Typography;

/**
 * TopicCard 组件
 * 
 * 使用 React.memo 优化渲染性能 (T137)
 * 仅在 topic 数据变化时重新渲染
 */
const TopicCard = memo(function TopicCard({ topic }: TopicCardProps) {
  const router = useRouter();

  return (
    <Card
      hoverable
      onClick={() => router.push(`/topics/${topic.key}`)}
      className="h-full"
      styles={{ body: { height: "100%" } }}
    >
      <Space direction="vertical" size="middle" className="h-full w-full">
        <Space align="center" size="small">
          <Tag color="blue">主题</Tag>
          <Text type="secondary">{topic.key}</Text>
        </Space>
        <Title level={4} className="mb-0">
          {topic.title}
        </Title>
        <Paragraph ellipsis={{ rows: 2 }}>
          {topic.summary || "快速开始该主题的学习。"}
        </Paragraph>
        <Text type="secondary">
          {topic.chapterCount > 0
            ? `章节数：${topic.chapterCount}`
            : "暂无章节"}
        </Text>
      </Space>
    </Card>
  );
});

export default TopicCard;
