"use client";

import useSWR from "swr";
import { Button, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import ChapterList from "@/components/learning/ChapterList";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchChapters } from "@/lib/learning";
import { ChapterSummary } from "@/types/learning";
import useProgress from "@/hooks/useProgress";
import ProgressBar from "@/components/progress/ProgressBar";
import { topicChapters } from "@/lib/static-routes";

const { Title, Paragraph } = Typography;

export default function TopicDetailClient({
  params,
}: {
  params: { topic: string };
}) {
  const router = useRouter();
  const topicKey = params.topic;
  const { data, error, isLoading } = useSWR<ChapterSummary[]>(
    ["chapters", topicKey],
    () => fetchChapters(topicKey),
  );
  const {
    chapters: chapterProgress,
    isLoading: progressLoading,
    topicDetail,
  } = useProgress(topicKey);
  const nextChapter =
    chapterProgress.find((c) => c.status !== "completed")?.chapter ??
    topicChapters[topicKey as keyof typeof topicChapters]?.[0];

  if (isLoading) {
    return <Loading />;
  }

  if (error) {
    return (
      <ErrorMessage message="加载章节列表失败" description={error.message} />
    );
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
          <ProgressBar
            percent={topicDetail?.progress ?? 0}
            status={
              (topicDetail?.progress ?? 0) >= 100
                ? "completed"
                : "in_progress"
            }
            segments={10}
            label="主题进度"
          />
        </div>
        <Space>
          {nextChapter && (
            <Button
              type="primary"
              onClick={() => router.push(`/topics/${topicKey}/${nextChapter}`)}
            >
              继续学习 {nextChapter}
            </Button>
          )}
          <Button onClick={() => router.push("/topics")}>返回主题</Button>
        </Space>
      </div>
      <Space direction="vertical" className="w-full">
        <ChapterList
          topicKey={topicKey}
          chapters={chapters}
          progress={progressLoading ? [] : chapterProgress}
        />
      </Space>
    </div>
  );
}
