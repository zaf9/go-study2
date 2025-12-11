"use client";

import dynamic from "next/dynamic";
import useSWR from "swr";
import { Button, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useRef, useState } from "react";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchChapterContent } from "@/lib/learning";
import {
  ChapterContent as ChapterContentType,
  ProgressStatus,
} from "@/types/learning";
import useProgress from "@/hooks/useProgress";
import useScrollPosition from "@/hooks/useScrollPosition";
import ChapterProgress from "@/components/progress/ChapterProgress";
import { topicChapters } from "@/lib/static-routes";

const { Title, Paragraph } = Typography;

const ChapterContent = dynamic(
  () => import("@/components/learning/ChapterContent"),
  {
    loading: () => <Loading />,
    ssr: false,
  },
);

export default function ChapterPageClient({
  params,
}: {
  params: { topic: string; chapter: string };
}) {
  const router = useRouter();
  const { topic, chapter } = params;
  const { data, error, isLoading } = useSWR<ChapterContentType>(
    ["chapter", topic, chapter],
    () => fetchChapterContent(topic, chapter),
  );
  const { chapters, recordProgress } = useProgress(topic);
  const scroll = useScrollPosition();
  const lastPosRef = useRef(0);
  const lastTickRef = useRef<number>(Date.now());
  const [restored, setRestored] = useState(false);
  const estimatedSeconds = 600;

  useEffect(() => {
    lastPosRef.current = scroll.scrollY;
  }, [scroll.scrollY]);

  useEffect(() => {
    if (!data) return;
    void recordProgress({
      topic,
      chapter,
      readDuration: 0,
      scrollProgress: calcScrollPercent(),
      lastPosition: JSON.stringify({ scroll: lastPosRef.current }),
      estimatedSeconds,
    });
    return () => {
      const now = Date.now();
      const elapsed = Math.max(0, Math.round((now - lastTickRef.current) / 1000));
      void recordProgress({
        topic,
        chapter,
        readDuration: elapsed,
        scrollProgress: calcScrollPercent(),
        lastPosition: JSON.stringify({ scroll: lastPosRef.current }),
        estimatedSeconds,
        forceSync: true,
      });
    };
  }, [data, topic, chapter, recordProgress]);

  useEffect(() => {
    const interval = window.setInterval(() => {
      const now = Date.now();
      const elapsed = Math.max(0, Math.round((now - lastTickRef.current) / 1000));
      lastTickRef.current = now;
      void recordProgress({
        topic,
        chapter,
        readDuration: elapsed,
        scrollProgress: calcScrollPercent(),
        lastPosition: JSON.stringify({ scroll: lastPosRef.current }),
        estimatedSeconds,
      });
    }, 10000);
    return () => window.clearInterval(interval);
  }, [topic, chapter, recordProgress]);

  const currentStatus: ProgressStatus = useMemo(
    () =>
      chapters.find((p) => p.chapter === chapter)?.status ?? "not_started",
    [chapters, chapter],
  );

  const chapterProgress = useMemo(
    () => chapters.find((p) => p.chapter === chapter),
    [chapters, chapter],
  );

  useEffect(() => {
    if (restored || !chapterProgress?.lastPosition) {
      return;
    }
    try {
      const parsed = JSON.parse(chapterProgress.lastPosition);
      if (typeof parsed?.scroll === "number" && parsed.scroll > 0) {
        window.scrollTo({ top: parsed.scroll, behavior: "smooth" });
        setRestored(true);
      }
    } catch {
      // ignore parse error
    }
  }, [chapterProgress?.lastPosition, restored]);

  if (isLoading) {
    return <Loading />;
  }

  if (error) {
    return (
      <ErrorMessage message="加载章节内容失败" description={error.message} />
    );
  }

  const content = data;
  const topicOrder = topicChapters[topic as keyof typeof topicChapters] ?? [];
  const currentIdx = topicOrder.indexOf(chapter);
  const prevChapter = currentIdx > 0 ? topicOrder[currentIdx - 1] : null;
  const nextChapter =
    currentIdx >= 0 && currentIdx < topicOrder.length - 1
      ? topicOrder[currentIdx + 1]
      : null;

  return (
    <Space direction="vertical" size="middle" className="w-full">
      <div className="flex items-center justify-between bg-white p-4 rounded shadow-sm w-full">
        <div>
          <Title level={3} className="mb-1">
            {content?.title || chapter}
          </Title>
          <Paragraph type="secondary">主题：{topic}</Paragraph>
          <ChapterProgress
            title="章节进度"
            status={currentStatus}
            scrollProgress={chapterProgress?.scrollProgress ?? 0}
            readDuration={chapterProgress?.readDuration ?? 0}
            estimatedSeconds={estimatedSeconds}
            lastVisitAt={chapterProgress?.lastVisitAt}
            onResume={
              restored
                ? undefined
                : () => {
                    let target = 0;
                    try {
                      const parsed = chapterProgress?.lastPosition
                        ? JSON.parse(chapterProgress.lastPosition)
                        : null;
                      if (typeof parsed?.scroll === "number") {
                        target = parsed.scroll;
                      }
                    } catch {
                      target = 0;
                    }
                    if (target > 0) {
                      window.scrollTo({ top: target, behavior: "smooth" });
                      setRestored(true);
                    }
                  }
            }
          />
        </div>
        <Space>
          <Button onClick={() => router.push(`/topics/${topic}`)}>
            返回章节列表
          </Button>
          <Button type="primary" onClick={() => router.push("/topics")}>
            返回主题
          </Button>
        </Space>
      </div>
      {content ? (
        <ChapterContent content={content} />
      ) : (
        <ErrorMessage message="暂无内容" />
      )}
      <div className="flex items-center justify-between bg-white p-3 rounded shadow-sm">
        <Space>
          <Button
            disabled={!prevChapter}
            onClick={() =>
              prevChapter && router.push(`/topics/${topic}/${prevChapter}`)
            }
          >
            上一章
          </Button>
          <Button onClick={() => router.push(`/topics/${topic}`)}>
            返回列表
          </Button>
        </Space>
        <Space>
          <Button
            onClick={() => router.push(`/quiz/${topic}/${chapter}`)}
            type="default"
          >
            开始测验
          </Button>
          <Button
            type="primary"
            disabled={!nextChapter}
            onClick={() =>
              nextChapter && router.push(`/topics/${topic}/${nextChapter}`)
            }
          >
            下一章
          </Button>
        </Space>
      </div>
    </Space>
  );
}

function calcScrollPercent(): number {
  if (typeof window === "undefined") return 0;
  const doc = document.documentElement;
  const total = doc.scrollHeight - window.innerHeight;
  if (total <= 0) return 100;
  return Math.min(100, Math.max(0, Math.round((window.scrollY / total) * 100)));
}
