'use client';

import useSWR from "swr";
import { Button, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useRef } from "react";
import ChapterContent from "@/components/learning/ChapterContent";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchChapterContent } from "@/lib/learning";
import { ChapterContent as ChapterContentType, ProgressStatus } from "@/types/learning";
import useProgress from "@/hooks/useProgress";
import useScrollPosition from "@/hooks/useScrollPosition";
import ProgressBar from "@/components/learning/ProgressBar";

const { Title, Paragraph } = Typography;

export default function ChapterPage({ params }: { params: { topic: string; chapter: string } }) {
  const router = useRouter();
  const { topic, chapter } = params;
  const { data, error, isLoading } = useSWR<ChapterContentType>(["chapter", topic, chapter], () =>
    fetchChapterContent(topic, chapter)
  );
  const { progress, recordProgress } = useProgress(topic);
  const scroll = useScrollPosition();
  const lastPosRef = useRef(0);

  useEffect(() => {
    lastPosRef.current = scroll.scrollY;
  }, [scroll.scrollY]);

  useEffect(() => {
    if (!data) return;
    void recordProgress({
      topic,
      chapter,
      status: "in_progress",
      position: JSON.stringify({ scroll: lastPosRef.current }),
    });
    return () => {
      void recordProgress({
        topic,
        chapter,
        status: "done",
        position: JSON.stringify({ scroll: lastPosRef.current }),
      });
    };
  }, [data, topic, chapter, recordProgress]);

  const currentStatus: ProgressStatus =
    useMemo(
      () => progress.find((p) => p.chapter === chapter)?.status ?? "not_started",
      [progress, chapter]
    );

  if (isLoading) {
    return <Loading />;
  }

  if (error) {
    return <ErrorMessage message="加载章节内容失败" description={error.message} />;
  }

  const content = data;

  return (
    <Space direction="vertical" size="middle" className="w-full">
      <div className="flex items-center justify-between bg-white p-4 rounded shadow-sm w-full">
        <div>
          <Title level={3} className="mb-1">
            {content?.title || chapter}
          </Title>
          <Paragraph type="secondary">主题：{topic}</Paragraph>
          <ProgressBar status={currentStatus} />
        </div>
        <Space>
          <Button onClick={() => router.push(`/topics/${topic}`)}>返回章节列表</Button>
          <Button type="primary" onClick={() => router.push("/topics")}>
            返回主题
          </Button>
        </Space>
      </div>
      {content ? <ChapterContent content={content} /> : <ErrorMessage message="暂无内容" />}
    </Space>
  );
}


