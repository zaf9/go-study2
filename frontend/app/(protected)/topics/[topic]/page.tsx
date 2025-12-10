'use client';

import useSWR from "swr";
import { Button, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import ChapterList from "@/components/learning/ChapterList";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchChapters } from "@/lib/learning";
import { ChapterSummary } from "@/types/learning";
import useProgress from "@/hooks/useProgress";

const { Title, Paragraph } = Typography;

export default function TopicDetailPage({ params }: { params: { topic: string } }) {
  const router = useRouter();
  const topicKey = params.topic;
  const { data, error, isLoading } = useSWR<ChapterSummary[]>(["chapters", topicKey], () =>
    fetchChapters(topicKey)
  );
  const { progress, isLoading: progressLoading } = useProgress(topicKey);

  if (isLoading) {
    return <Loading />;
  }

  if (error) {
    return <ErrorMessage message="加载章节列表失败" description={error.message} />;
  }

  const chapters = data ?? [];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between bg-white p-4 rounded shadow-sm">
        <div>
          <Title level={3} className="mb-1">
            {topicKey} 章节列表
          </Title>
          <Paragraph type="secondary">选择章节进入详情学习。</Paragraph>
        </div>
        <Button onClick={() => router.push("/topics")}>返回主题</Button>
      </div>
      <Space direction="vertical" className="w-full">
        <ChapterList topicKey={topicKey} chapters={chapters} progress={progressLoading ? [] : progress} />
      </Space>
    </div>
  );
}


