"use client";

import { Badge, List, Space, Tag, Typography } from "antd";
import { QuizAnswerDetail } from "@/types/quiz";

/** 答案解析属性：传入判分明细，逐题展示正确答案与说明。 */
interface AnswerExplanationProps {
  details: QuizAnswerDetail[];
}

const { Text } = Typography;

export default function AnswerExplanation({ details }: AnswerExplanationProps) {
  return (
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
                <span>题目 {item.question_id}</span>
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
              <div className="space-y-1">
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
  );
}

