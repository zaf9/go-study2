'use client';

import { Card, Space, Tag, Typography } from "antd";
import { useRouter } from "next/navigation";
import { TopicSummary } from "@/types/learning";

interface TopicCardProps {
  topic: TopicSummary;
}

const { Title, Paragraph, Text } = Typography;

export default function TopicCard({ topic }: TopicCardProps) {
  const router = useRouter();

  return (
    <Card
      hoverable
      onClick={() => router.push(`/topics/${topic.key}`)}
      className="h-full"
      bodyStyle={{ height: "100%" }}
    >
      <Space direction="vertical" size="middle" className="h-full w-full">
        <Space align="center" size="small">
          <Tag color="blue">主题</Tag>
          <Text type="secondary">{topic.key}</Text>
        </Space>
        <Title level={4} className="mb-0">
          {topic.title}
        </Title>
        <Paragraph ellipsis={{ rows: 2 }}>{topic.summary || "快速开始该主题的学习。"}</Paragraph>
        <Text type="secondary">章节数：{topic.chapterCount ?? 0}</Text>
      </Space>
    </Card>
  );
}


