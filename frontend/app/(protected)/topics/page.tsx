'use client';

import useSWR from "swr";
import { Button, Col, Row, Typography } from "antd";
import { useRouter } from "next/navigation";
import TopicCard from "@/components/learning/TopicCard";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchTopics } from "@/lib/learning";
import { TopicSummary } from "@/types/learning";
import useProgress from "@/hooks/useProgress";

const { Title, Paragraph } = Typography;

export default function TopicsPage() {
  const router = useRouter();
  const { data, error, isLoading } = useSWR<TopicSummary[]>("topics", fetchTopics);
  const { latest, isLoading: progressLoading } = useProgress();

  if (isLoading) {
    return <Loading />;
  }

  if (error) {
    return <ErrorMessage message="加载主题列表失败" description={error.message} />;
  }

  const topics = data ?? [];

  return (
    <div className="space-y-4">
      <div className="bg-white p-4 rounded shadow-sm">
        <Title level={3} className="mb-1">
          学习主题
        </Title>
        <Paragraph type="secondary">选择主题，进入章节学习并查看代码示例。</Paragraph>
        {latest && !progressLoading && (
          <Button type="primary" onClick={() => router.push(`/topics/${latest.topic}/${latest.chapter}`)}>
            继续学习：{latest.topic} / {latest.chapter}
          </Button>
        )}
      </div>
      <Row gutter={[16, 16]}>
        {topics.map((topic) => (
          <Col key={topic.key} xs={24} sm={12} lg={8}>
            <TopicCard topic={topic} />
          </Col>
        ))}
      </Row>
    </div>
  );
}


