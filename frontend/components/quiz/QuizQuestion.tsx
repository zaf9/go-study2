"use client";

import { Checkbox, Radio, Space, Typography } from "antd";
import { QuizQuestion } from "@/types/quiz";

/** 题目卡片属性：传入题目信息与当前选项，回调返回选择结果。 */
interface QuizQuestionProps {
  question: QuizQuestion;
  value: string[];
  onChange: (choices: string[]) => void;
}

const { Paragraph, Text, Title } = Typography;

export default function QuizQuestionCard({
  question,
  value,
  onChange,
}: QuizQuestionProps) {
  const isMultiple =
    question.type === "multiple" || question.type === "code_correction";
  const isSingle =
    question.type === "single" ||
    question.type === "truefalse" ||
    question.type === "code_output";

  const options = question.options.map((opt) => ({
    label: opt.label,
    value: opt.id,
  }));

  return (
    <div className="bg-white p-4 rounded shadow-sm space-y-3">
      <Title level={5} className="mb-1">
        {question.question}
      </Title>
      <Paragraph type="secondary" className="mb-2">
        {labelByType(question.type)} · 难度：{question.difficulty}
      </Paragraph>
      {question.codeSnippet && (
        <pre className="bg-gray-50 p-3 rounded border">{question.codeSnippet}</pre>
      )}
      {isMultiple && (
        <Checkbox.Group
          options={options}
          value={value}
          onChange={(val) => onChange((val as string[]) ?? [])}
        />
      )}
      {isSingle && (
        <Radio.Group
          value={value?.[0]}
          onChange={(e) => onChange([e.target.value])}
        >
          <Space direction="vertical">
            {options.map((opt) => (
              <Radio key={opt.value} value={opt.value}>
                <Text>{opt.label}</Text>
              </Radio>
            ))}
          </Space>
        </Radio.Group>
      )}
    </div>
  );
}

function labelByType(t: string) {
  switch (t) {
    case "multiple":
      return "多选题";
    case "truefalse":
      return "判断题";
    case "code_output":
      return "代码输出题";
    case "code_correction":
      return "改错题";
    default:
      return "单选题";
  }
}

