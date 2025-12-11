"use client";

import { Checkbox, Radio, Space, Typography } from "antd";
import { QuizItem as QuizItemType } from "@/types/quiz";

interface QuizItemProps {
  question: QuizItemType;
  value: string[];
  onChange: (choices: string[]) => void;
}

export default function QuizItem({ question, value, onChange }: QuizItemProps) {
  const { Title, Paragraph } = Typography;
  const options = question.options.map((opt) => ({
    label: opt.label,
    value: opt.id,
  }));

  return (
    <div className="bg-white p-4 rounded shadow-sm space-y-3">
      <Title level={5} className="mb-1">
        {question.stem}
      </Title>
      <Paragraph type="secondary" className="mb-2">
        {question.multi ? "多选题" : "单选题"}
      </Paragraph>
      {question.multi ? (
        <Checkbox.Group
          options={options}
          value={value}
          onChange={(val) => onChange((val as string[]) ?? [])}
        />
      ) : (
        <Radio.Group
          value={value?.[0]}
          onChange={(e) => onChange([e.target.value])}
        >
          <Space direction="vertical">
            {options.map((opt) => (
              <Radio key={opt.value} value={opt.value}>
                {opt.label}
              </Radio>
            ))}
          </Space>
        </Radio.Group>
      )}
    </div>
  );
}
