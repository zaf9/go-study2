"use client";

import { Button, Space, Statistic } from "antd";

/** 测验导航属性：当前题序/总数、作答数与导航/提交回调。 */
interface QuizNavigationProps {
  current: number;
  total: number;
  answered: number;
  onPrev: () => void;
  onNext: () => void;
  onSubmit: () => void;
  submitting?: boolean;
}

export default function QuizNavigation({
  current,
  total,
  answered,
  onPrev,
  onNext,
  onSubmit,
  submitting,
}: QuizNavigationProps) {
  return (
    <div className="flex items-center justify-between bg-white p-3 rounded shadow-sm">
      <Space>
        <Statistic title="当前题目" value={`${current + 1}/${total}`} />
        <Statistic title="已作答" value={answered} />
      </Space>
      <Space>
        <Button onClick={onPrev} disabled={current <= 0}>
          上一题
        </Button>
        <Button onClick={onNext} disabled={current >= total - 1}>
          下一题
        </Button>
        <Button type="primary" onClick={onSubmit} loading={submitting}>
          提交测验
        </Button>
      </Space>
    </div>
  );
}

