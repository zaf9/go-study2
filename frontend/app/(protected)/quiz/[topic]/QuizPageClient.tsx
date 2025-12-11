"use client";

import { useEffect, useState } from "react";
import useSWR from "swr";
import { Divider, Select, Space, Typography, message } from "antd";
import QuizQuestionCard from "@/components/quiz/QuizQuestion";
import QuizResultView from "@/components/quiz/QuizResult";
import QuizNavigation from "@/components/quiz/QuizNavigation";
import AnswerExplanation from "@/components/quiz/AnswerExplanation";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { fetchChapters } from "@/lib/learning";
import useQuiz from "@/hooks/useQuiz";

const { Title, Paragraph } = Typography;

export default function QuizPageClient({
  params,
}: {
  params: { topic: string; chapter?: string };
}) {
  const topic = params.topic;
  const {
    data: chapters,
    error: chapterError,
    isLoading: chapterLoading,
  } = useSWR(["quiz-chapters", topic], () => fetchChapters(topic));
  const [selectedChapter, setSelectedChapter] = useState(params.chapter ?? "");
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    if (params.chapter) {
      setSelectedChapter(params.chapter);
      return;
    }
    if (chapters && chapters.length > 0 && !selectedChapter) {
      setSelectedChapter(chapters[0].id);
    }
  }, [chapters, selectedChapter, params.chapter]);

  const {
    questions,
    answers,
    selectAnswer,
    submit,
    result,
    isLoading: quizLoading,
    submitting,
    reset,
    answeredCount,
  } = useQuiz(topic, selectedChapter);

  if (chapterLoading) return <Loading />;
  if (chapterError)
    return (
      <ErrorMessage message="加载章节失败" description={chapterError.message} />
    );
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
            <Paragraph type="secondary">
              选择章节后完成测验，提交即可查看结果。
            </Paragraph>
          </div>
          {!params.chapter && (
            <Select
              style={{ minWidth: 220 }}
              value={selectedChapter}
              onChange={(v) => {
                setSelectedChapter(v);
                reset();
              }}
              options={(chapters ?? []).map((c) => ({
                label: c.title,
                value: c.id,
              }))}
            />
          )}
        </div>
        {quizLoading ? (
          <Loading />
        ) : (
          <Space direction="vertical" className="w-full" size="middle">
            {questions.length > 0 && (
              <>
                <QuizQuestionCard
                  question={{
                    id: Number(questions[currentIndex].id),
                    type: questions[currentIndex].type ?? "single",
                    difficulty: questions[currentIndex].difficulty ?? "easy",
                    question: questions[currentIndex].stem,
                    options: questions[currentIndex].options,
                  }}
                  value={answers[Number(questions[currentIndex].id)] ?? []}
                  onChange={(val) => selectAnswer(questions[currentIndex].id, val)}
                />
                <QuizNavigation
                  current={currentIndex}
                  total={questions.length}
                  answered={answeredCount}
                  onPrev={() => setCurrentIndex(Math.max(0, currentIndex - 1))}
                  onNext={() =>
                    setCurrentIndex(
                      Math.min(questions.length - 1, currentIndex + 1),
                    )
                  }
                  onSubmit={handleSubmit}
                  submitting={submitting}
                />
              </>
            )}
            {questions.length === 0 && (
              <ErrorMessage message="当前章节暂无测验" />
            )}
            {result && (
              <>
                <Divider />
                <QuizResultView result={result} />
                {result.details && result.details.length > 0 && (
                  <AnswerExplanation details={result.details} />
                )}
              </>
            )}
          </Space>
        )}
      </Space>
    </div>
  );
}
