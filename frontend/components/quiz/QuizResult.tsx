"use client";

import { Alert, Typography } from "antd";
import { QuizSubmitResult } from "@/types/quiz";

interface QuizResultProps {
  result: QuizSubmitResult;
}

const { Title } = Typography;

export default function QuizResultView({ result }: QuizResultProps) {
  return (
    <div className="space-y-3">
      <Alert
        message={`得分：${result.score} / 100`}
        description={`正确题数：${result.correct_answers}/${result.total_questions}`}
        type={result.passed ? "success" : "warning"}
        showIcon
      />
      <Title level={5}>答题详情</Title>
      <div className="text-sm text-gray-600">
        通过标准：60 分；本次 {result.passed ? "已通过" : "未通过"}。
      </div>
    </div>
  );
}
