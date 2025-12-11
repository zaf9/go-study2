"use client";

import { Select, Space, Typography } from "antd";
import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import ProgressOverview from "@/components/progress/ProgressOverview";
import TopicProgressCard from "@/components/progress/TopicProgressCard";
import { topicChapters } from "@/lib/static-routes";
import { TopicProgressDetail } from "@/types/learning";
import useProgress from "@/hooks/useProgress";

const { Title } = Typography;

export default function ProgressPage() {
  const router = useRouter();
  const { overview, next, isLoading, error } = useProgress();
  const [selectedTopic, setSelectedTopic] = useState<string | undefined>();

  const topics = useMemo(
    () =>
      (overview?.topics ?? [])
        .map<TopicProgressDetail>((item) => ({
          ...item,
          chapters: [],
        }))
        .sort((a, b) => b.progress - a.progress),
    [overview?.topics],
  );

  const filteredTopics = useMemo(() => {
    if (!selectedTopic) return topics;
    return topics.filter((t) => t.id === selectedTopic);
  }, [selectedTopic, topics]);

  if (isLoading) return <Loading />;
  if (error)
    return <ErrorMessage message="加载进度失败" description={error.message} />;
  if (!overview) return null;

  return (
    <Space direction="vertical" className="w-full">
      <Title level={3}>学习进度</Title>
      <ProgressOverview
        overall={overview.overall}
        next={next}
        onContinue={(hint) =>
          router.push(`/topics/${hint.topic}/${hint.chapter}`)
        }
      />
      <div className="flex items-center justify-between">
        <Title level={4} className="mb-0">
          主题进度
        </Title>
        <Select
          allowClear
          placeholder="筛选主题"
          value={selectedTopic}
          onChange={(v) => setSelectedTopic(v)}
          options={topics.map((t) => ({ label: t.name, value: t.id }))}
          style={{ minWidth: 200 }}
        />
      </div>
      <Space direction="vertical" className="w-full">
        {filteredTopics.map((topic) => (
          <TopicProgressCard
            key={topic.id}
            topic={topic}
            onContinue={(chapter) => {
              const target =
                chapter?.chapter ??
                topicChapters[topic.id as keyof typeof topicChapters]?.[0];
              if (target) {
                router.push(`/topics/${topic.id}/${target}`);
              }
            }}
          />
        ))}
      </Space>
    </Space>
  );
}
