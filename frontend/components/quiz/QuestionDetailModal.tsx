"use client";

import { Modal, Typography, Tag, Space, Alert, Radio, Checkbox } from "antd";
import { QuizQuestion } from "@/types/quiz";

const { Title, Text } = Typography;

interface QuestionDetailModalProps {
  open: boolean;
  onClose: () => void;
  question: QuizQuestion | null;
  userAnswers?: string[];
  correctAnswers?: string[];
  isCorrect?: boolean;
}

export default function QuestionDetailModal({
  open,
  onClose,
  question,
  userAnswers = [],
  correctAnswers = [],
  isCorrect = false,
}: QuestionDetailModalProps) {
  if (!question) return null;

  const isMultiple = question.type === "multiple" || question.type === "code_correction";

  // 获取选项标签
  const getOptionLabel = (index: number): string => {
    if (index < 26) {
      return String.fromCharCode(65 + index);
    } else {
      let label = '';
      let num = index - 26;
      while (num >= 0) {
        label = String.fromCharCode(65 + (num % 26)) + label;
        num = Math.floor(num / 26) - 1;
        if (num < 0) break;
      }
      if (label.length === 1) {
        label = 'A' + label;
      }
      return label;
    }
  };

  // 检查选项是否被用户选中
  const isUserSelected = (optionId: string) => {
    return userAnswers.includes(optionId);
  };

  // 检查选项是否是正确答案
  const isCorrectOption = (optionId: string) => {
    return correctAnswers.includes(optionId);
  };

  // 获取选项状态样式
  const getOptionStatus = (optionId: string) => {
    const userSelected = isUserSelected(optionId);
    const correct = isCorrectOption(optionId);

    if (userSelected && correct) {
      return { borderColor: '#52c41a', backgroundColor: '#f6ffed' }; // 绿色 - 用户选的且正确
    } else if (userSelected && !correct) {
      return { borderColor: '#ff4d4f', backgroundColor: '#fff1f0' }; // 红色 - 用户选的但错误
    } else if (!userSelected && correct) {
      return { borderColor: '#faad14', backgroundColor: '#fffbe6' }; // 橙色 - 正确但用户没选
    }
    return {}; // 默认样式
  };

  return (
    <Modal
      title={
        <Space>
          <span>题目详情</span>
          <Tag color={isCorrect ? "success" : "error"}>
            {isCorrect ? "回答正确" : "回答错误"}
          </Tag>
        </Space>
      }
      open={open}
      onCancel={onClose}
      width={800}
      footer={null}
    >
      <div className="space-y-4">
        {/* 题干 */}
        <div>
          <Title level={4}>{question.question}</Title>
        </div>

        {/* 代码片段 */}
        {question.codeSnippet && (
          <div>
            <pre className="bg-gray-800 text-gray-100 p-4 rounded-md overflow-x-auto font-mono text-sm">
              <code>{question.codeSnippet}</code>
            </pre>
          </div>
        )}

        {/* 选项 */}
        <div>
          <Space direction="vertical" style={{ width: "100%" }} size="middle">
            {question.options.map((opt, index) => {
              const label = getOptionLabel(index);
              const statusStyle = getOptionStatus(opt.id);
              const userSelected = isUserSelected(opt.id);
              const correct = isCorrectOption(opt.id);

              return (
                <div
                  key={opt.id}
                  style={{
                    padding: '12px 16px',
                    border: '2px solid #d9d9d9',
                    borderRadius: '8px',
                    backgroundColor: '#fff',
                    transition: 'all 0.3s',
                    ...statusStyle,
                  }}
                >
                  <Space direction="vertical" style={{ width: "100%" }}>
                    <div style={{ display: 'flex', alignItems: 'flex-start' }}>
                      {/* 复选框/单选框图标 */}
                      <span style={{ marginRight: 12, marginTop: 2 }}>
                        {isMultiple ? (
                          <Checkbox checked={userSelected} disabled>
                            <Text strong style={{ marginRight: 8 }}>
                              {label}.
                            </Text>
                          </Checkbox>
                        ) : (
                          <Radio checked={userSelected} disabled>
                            <Text strong style={{ marginRight: 8 }}>
                              {label}.
                            </Text>
                          </Radio>
                        )}
                      </span>

                      {/* 选项内容 */}
                      <Text style={{ flex: 1 }}>{opt.label}</Text>

                      {/* 状态标记 */}
                      <Space size={4}>
                        {userSelected && (
                          <Tag color={correct ? "success" : "error"}>
                            {correct ? "✓ 你的答案" : "✗ 你的答案"}
                          </Tag>
                        )}
                        {correct && !userSelected && (
                          <Tag color="warning">
                            正确答案 (未选)
                          </Tag>
                        )}
                      </Space>
                    </div>
                  </Space>
                </div>
              );
            })}
          </Space>
        </div>

        {/* 答案总结 */}
        <Alert
          message="答案总结"
          description={
            <div>
              <div className="mb-2">
                <Text strong>你的答案：</Text>
                <Text className={isCorrect ? "text-green-600" : "text-red-600"}>
                  {userAnswers.length > 0 ? userAnswers.join(", ") : "未作答"}
                </Text>
              </div>
              <div>
                <Text strong>正确答案：</Text>
                <Text type="success">
                  {correctAnswers.join(", ")}
                </Text>
              </div>
            </div>
          }
          type={isCorrect ? "success" : "error"}
          showIcon
        />

        {/* 题目信息 */}
        <div>
          <Space>
            <Text type="secondary">题型：</Text>
            <Tag color={question.type === 'single' ? 'blue' : 'purple'}>
              {question.type === 'single' ? '单选题' :
               question.type === 'multiple' ? '多选题' :
               question.type === 'truefalse' ? '判断题' :
               question.type === 'code_output' ? '代码输出题' :
               question.type === 'code_correction' ? '代码改错题' : '未知'}
            </Tag>
            <Text type="secondary">难度：</Text>
            <Tag color={
              question.difficulty === 'easy' ? 'green' :
              question.difficulty === 'medium' ? 'orange' : 'red'
            }>
              {question.difficulty === 'easy' ? '简单' :
               question.difficulty === 'medium' ? '中等' : '困难'}
            </Tag>
          </Space>
        </div>
      </div>
    </Modal>
  );
}
