'use client';

import { useEffect, useState } from "react";
import useSWR from "swr";
import { Button, Divider, Select, Space, Typography, message } from "antd";
import QuizItem from "@/components/quiz/QuizItem";
import QuizResultView from "@/components/quiz/QuizResult";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchChapters } from "@/lib/learning";
import useQuiz from "@/hooks/useQuiz";

const { Title, Paragraph } = Typography;

export default function QuizPageClient({ params }: { params: { topic: string } }) {
  const topic = params.topic;
  const { data: chapters, error: chapterError, isLoading: chapterLoading } = useSWR(["quiz-chapters", topic], () =>
    fetchChapters(topic)
  );
  const [selectedChapter, setSelectedChapter] = useState("");

  useEffect(() => {
    if (chapters && chapters.length > 0 && !selectedChapter) {
      setSelectedChapter(chapters[0].id);
    }
  }, [chapters, selectedChapter]);

  const {
    questions,
    answers,
    selectAnswer,
    submit,
    result,
    isLoading: quizLoading,
    submitting,
    reset,
  } = useQuiz(topic, selectedChapter);

  if (chapterLoading) return <Loading />;
  if (chapterError) return <ErrorMessage message="加载章节失败" description={chapterError.message} />;
  if (!selectedChapter) return <ErrorMessage message="暂无章节可用" />;

  const handleSubmit = async () => {
    const res = await submit();
    if (res) {
      message.success("提交成功");
    }
  };

  return (
    <div className="space-y-4">
      <Space direction="vertical" className="w-full">
        <div className="flex items-center justify-between bg-white p-4 rounded shadow-sm">
          <div>
            <Title level={3} className="mb-1">
              {topic} 测验
            </Title>
            <Paragraph type="secondary">选择章节后完成测验，提交即可查看结果。</Paragraph>
          </div>
          <Select
            style={{ minWidth: 220 }}
            value={selectedChapter}
            onChange={(v) => {
              setSelectedChapter(v);
              reset();
            }}
            options={(chapters ?? []).map((c) => ({ label: c.title, value: c.id }))}
          />
        </div>
        {quizLoading ? (
          <Loading />
        ) : (
          <Space direction="vertical" className="w-full" size="middle">
            {questions.map((q) => (
              <QuizItem key={q.id} question={q} value={answers[q.id] ?? []} onChange={(val) => selectAnswer(q.id, val)} />
            ))}
            {questions.length === 0 && <ErrorMessage message="当前章节暂无测验" />}
            <Space>
              <Button type="primary" onClick={handleSubmit} loading={submitting} disabled={questions.length === 0}>
                提交测验
              </Button>
              <Button onClick={reset}>重置</Button>
            </Space>
            {result && (
              <>
                <Divider />
                <QuizResultView result={result} />
              </>
            )}
          </Space>
        )}
      </Space>
    </div>
  );
}

