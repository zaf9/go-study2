"use client";

import { Table, Tag } from "antd";
import { QuizHistoryItem } from "@/types/quiz";

interface QuizHistoryProps {
  items: QuizHistoryItem[];
  loading?: boolean;
}

export default function QuizHistory({ items, loading }: QuizHistoryProps) {
  const columns = [
    { title: "主题", dataIndex: "topic", key: "topic" },
    {
      title: "章节",
      dataIndex: "chapter",
      key: "chapter",
      render: (v: string | null) => v || "-",
    },
    {
      title: "得分",
      dataIndex: "score",
      key: "score",
      render: (_: unknown, row: QuizHistoryItem) =>
        row.totalQuestions
          ? `${row.score}/${row.totalQuestions}`
          : `${row.score}`,
    },
    {
      title: "完成时间",
      dataIndex: "completedAt",
      key: "completedAt",
      render: (v: string) => (v ? new Date(v).toLocaleString() : "-"),
    },
    {
      title: "状态",
      key: "status",
      render: (_: unknown, row: QuizHistoryItem) => (
        <Tag color={row.passed ? "green" : "blue"}>
          {row.passed ? "通过" : "未通过"}
        </Tag>
      ),
    },
  ];

  return (
    <Table
      rowKey="id"
      columns={columns}
      dataSource={items}
      loading={loading}
      pagination={false}
    />
  );
}
