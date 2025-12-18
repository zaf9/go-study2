"use client";

import { Tag } from "antd";

interface AnswerIndicatorProps {
  userChoice: string;
  correctChoice: string;
  isCorrect: boolean;
}

/**
 * AnswerIndicator 组件用于显示用户答案与正确答案的对比标识
 */
export default function AnswerIndicator({
  userChoice,
  correctChoice,
  isCorrect,
}: AnswerIndicatorProps) {
  return (
    <div className="flex gap-2 items-center">
      {userChoice && (
        <Tag color={isCorrect ? "green" : "blue"}>
          你的答案: {userChoice}
        </Tag>
      )}
      {!isCorrect && correctChoice && (
        <Tag color="green">正确答案: {correctChoice}</Tag>
      )}
      {isCorrect && (
        <Tag color="success">回答正确</Tag>
      )}
    </div>
  );
}

