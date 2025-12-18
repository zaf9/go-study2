"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Typography, Button, Card, Tag, Space } from "antd";
import { LeftOutlined } from "@ant-design/icons";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import AnswerIndicator from "@/components/quiz/AnswerIndicator";
import { fetchQuizReview } from "@/lib/quiz";
import { QuizReviewDetail } from "@/types/quiz";

const { Title, Paragraph } = Typography;

export default function QuizReviewPage() {
  const params = useParams();
  const router = useRouter();
  const sessionId = params.sessionId as string;
  const [review, setReview] = useState<QuizReviewDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!sessionId) {
      setError(new Error("会话ID不能为空"));
      setLoading(false);
      return;
    }

    const loadReview = async () => {
      try {
        setLoading(true);
        const data = await fetchQuizReview(sessionId);
        setReview(data);
      } catch (err) {
        setError(err instanceof Error ? err : new Error("加载失败"));
      } finally {
        setLoading(false);
      }
    };

    loadReview();
  }, [sessionId]);

  if (loading) return <Loading />;
  if (error) {
    return (
      <ErrorMessage
        message="加载回顾详情失败"
        description={error.message}
      />
    );
  }
  if (!review) {
    return <ErrorMessage message="未找到回顾详情" />;
  }

  const { meta, items } = review;

  return (
    <div className="space-y-4">
      <Space>
        <Button
          icon={<LeftOutlined />}
          onClick={() => router.push("/quiz/history")}
        >
          返回历史记录
        </Button>
        <Title level={3}>测验回顾</Title>
      </Space>

      <Card>
        <Space direction="vertical" className="w-full">
          <div>
            <strong>主题:</strong> {meta.topic}
          </div>
          <div>
            <strong>章节:</strong> {meta.chapter}
          </div>
          <div>
            <strong>得分:</strong> {meta.score} 分
            <Tag
              color={meta.passed ? "green" : "red"}
              className="ml-2"
            >
              {meta.passed ? "通过" : "未通过"}
            </Tag>
          </div>
          {meta.completedAt && (
            <div>
              <strong>完成时间:</strong>{" "}
              {new Date(meta.completedAt).toLocaleString()}
            </div>
          )}
        </Space>
      </Card>

      <div className="space-y-4">
        {items.map((item, index) => (
          <Card key={item.questionId} title={`题目 ${index + 1}`}>
            <Space direction="vertical" className="w-full" size="large">
              <Paragraph className="text-base font-medium">
                {item.stem}
              </Paragraph>

              <div className="space-y-2">
                <div className="font-medium">选项:</div>
                {item.options.map((option, optIdx) => {
                  const label = String.fromCharCode(65 + optIdx);
                  const isUserChoice = item.userChoice
                    .split(",")
                    .includes(option);
                  const isCorrectChoice = item.correctChoice
                    .split(",")
                    .includes(option);

                  return (
                    <div
                      key={optIdx}
                      className={`p-2 rounded ${
                        isUserChoice || isCorrectChoice
                          ? "bg-gray-100"
                          : ""
                      }`}
                    >
                      <Space>
                        <span className="font-medium">{label}.</span>
                        <span>{option}</span>
                        {isUserChoice && (
                          <Tag color={item.isCorrect ? "green" : "blue"}>
                            你的答案
                          </Tag>
                        )}
                        {isCorrectChoice && !isUserChoice && (
                          <Tag color="green">正确答案</Tag>
                        )}
                      </Space>
                    </div>
                  );
                })}
              </div>

              <AnswerIndicator
                userChoice={item.userChoice}
                correctChoice={item.correctChoice}
                isCorrect={item.isCorrect}
              />

              {item.explanation && (
                <div>
                  <div className="font-medium mb-2">解析:</div>
                  <Paragraph className="text-gray-600">
                    {item.explanation}
                  </Paragraph>
                </div>
              )}
            </Space>
          </Card>
        ))}
      </div>
    </div>
  );
}

