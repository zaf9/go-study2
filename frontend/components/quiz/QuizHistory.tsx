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
      render: (_: unknown, row: QuizHistoryItem) => `${row.score}/${row.total}`,
    },
    { title: "用时(ms)", dataIndex: "durationMs", key: "durationMs" },
    {
      title: "提交时间",
      dataIndex: "createdAt",
      key: "createdAt",
      render: (v: string) => new Date(v).toLocaleString(),
    },
    {
      title: "状态",
      key: "status",
      render: (_: unknown, row: QuizHistoryItem) => (
        <Tag color={row.score == row.total ? "green" : "blue"}>
          {row.score == row.total ? "全对" : "已提交"}
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
