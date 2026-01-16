"use client";

import { Badge, List, Space, Tag, Typography, Button } from "antd";
import { QuizAnswerDetail, QuizQuestion } from "@/types/quiz";
import { useState } from "react";
import DeepLearningDrawer from "./DeepLearningDrawer";
import QuestionDetailModal from "./QuestionDetailModal";

/** 答案解析属性：传入判分明细，逐题展示正确答案与说明。 */
interface AnswerExplanationProps {
  details: QuizAnswerDetail[];
  topic?: string;
  chapter?: string;
  chapterTitle?: string;
  questions?: QuizQuestion[];
  userAnswers?: Record<string, string[]>;
}

const { Text, Link } = Typography;

export default function AnswerExplanation({
  details,
  topic,
  chapter,
  chapterTitle,
  questions = [],
  userAnswers = {}
}: AnswerExplanationProps) {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [selectedQuestion, setSelectedQuestion] = useState<QuizQuestion | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const handleDeepLearn = () => {
    setDrawerOpen(true);
  };

  const handleQuestionClick = (detail: QuizAnswerDetail) => {
    // 通过 question_id 找到对应的完整题目
    // 注意：questions中的id是string类型，detail.question_id现在也是string类型（后端修复）
    const question = questions.find(q => q.id === detail.question_id);
    if (question) {
      setSelectedQuestion(question);
      setModalOpen(true);
    }
  };

  const canDeepLearn = topic && chapter;

  return (
    <>
      <List
        dataSource={details}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              avatar={
                <Badge
                  status={item.is_correct ? "success" : "error"}
                  text={item.is_correct ? "正确" : "错误"}
                />
              }
              title={
                <Space>
                  <Link
                    strong
                    onClick={() => handleQuestionClick(item)}
                    style={{ cursor: 'pointer' }}
                  >
                    题目 {item.question_id}
                  </Link>
                  <Tag color={item.type === 'single' ? 'blue' : 'purple'}>
                    {item.type === 'single' ? '单选题' : '多选题'}
                  </Tag>
                  <Tag color={
                    item.difficulty === 'easy' ? 'green' :
                    item.difficulty === 'medium' ? 'orange' : 'red'
                  }>
                    {item.difficulty === 'easy' ? '简单' :
                     item.difficulty === 'medium' ? '中等' : '困难'}
                  </Tag>
                </Space>
              }
              description={
                <div className="space-y-2">
                  <Text type="secondary">
                    正确答案：{item.correct_answers.join(", ")}
                  </Text>
                  <div>{item.explanation}</div>
                </div>
              }
            />
          </List.Item>
        )}
      />

      {/* 深入学习按钮 - 仅在第一项后显示一次 */}
      {canDeepLearn && details.length > 0 && (
        <div className="mt-4 pt-4 border-t border-gray-200">
          <div className="flex items-center justify-between bg-blue-50 p-4 rounded-lg border border-blue-200">
            <div className="flex-1">
              <Text strong className="text-blue-900 block mb-1">
                💡 想要更深入地了解这些知识点吗？
              </Text>
              <Text type="secondary" className="text-sm">
                查看完整章节内容，包含详细的概念讲解、代码示例和最佳实践
              </Text>
            </div>
            <Button
              type="primary"
              size="large"
              onClick={handleDeepLearn}
              className="ml-4"
            >
              深入学习 →
            </Button>
          </div>
        </div>
      )}

      <DeepLearningDrawer
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        topic={topic || ""}
        chapter={chapter || ""}
        chapterTitle={chapterTitle}
      />

      <QuestionDetailModal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        question={selectedQuestion}
        userAnswers={selectedQuestion ? (userAnswers[selectedQuestion.id] || []) : []}
        correctAnswers={selectedQuestion ? details.find(d => d.question_id === selectedQuestion.id)?.correct_answers || [] : []}
        isCorrect={selectedQuestion ? details.find(d => d.question_id === selectedQuestion.id)?.is_correct || false : false}
      />
    </>
  );
}

