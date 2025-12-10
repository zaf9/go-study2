'use client';

import useSWR from "swr";
import { Button, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import ChapterContent from "@/components/learning/ChapterContent";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchChapterContent } from "@/lib/learning";
import { ChapterContent as ChapterContentType } from "@/types/learning";

const { Title, Paragraph } = Typography;

export default function ChapterPage({ params }: { params: { topic: string; chapter: string } }) {
  const router = useRouter();
  const { topic, chapter } = params;
  const { data, error, isLoading } = useSWR<ChapterContentType>(["chapter", topic, chapter], () =>
    fetchChapterContent(topic, chapter)
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


