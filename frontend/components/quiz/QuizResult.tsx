'use client';

import { Alert, List, Typography } from "antd";
import { QuizResult } from "@/types/quiz";

interface QuizResultProps {
  result: QuizResult;
}

export default function QuizResultView({ result }: QuizResultProps) {
  const { Title, Text } = Typography;
  return (
    <div className="space-y-3">
      <Alert
        message={`得分：${result.score}/${result.total}`}
        description={`用时：${result.durationMs || 0} ms`}
        type="success"
        showIcon
      />
      <Title level={5}>答题详情</Title>
      <List
        dataSource={[...result.correctIds.map((id) => ({ id, correct: true })), ...result.wrongIds.map((id) => ({ id, correct: false }))]}
        renderItem={(item) => (
          <List.Item>
            <Text type={item.correct ? "success" : "danger"}>
              {item.id} - {item.correct ? "正确" : "错误"}
            </Text>
          </List.Item>
        )}
      />
    </div>
  );
}

