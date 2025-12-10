'use client';

import { Typography } from "antd";
import QuizHistory from "@/components/quiz/QuizHistory";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import { useQuizHistory } from "@/hooks/useQuiz";

const { Title } = Typography;

export default function QuizHistoryPage() {
  const { history, isLoading, error } = useQuizHistory();

  if (isLoading) return <Loading />;
  if (error) return <ErrorMessage message="加载测验历史失败" description={error.message} />;

  return (
    <div className="space-y-4">
      <Title level={3}>测验历史记录</Title>
      <QuizHistory items={history} />
    </div>
  );
}

