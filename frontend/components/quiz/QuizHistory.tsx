"use client";

import { Table, Tag, Button } from "antd";
import { useRouter } from "next/navigation";
import { QuizHistoryItem } from "@/types/quiz";

interface QuizHistoryProps {
  items: QuizHistoryItem[];
  loading?: boolean;
}

export default function QuizHistory({ items, loading }: QuizHistoryProps) {
  const router = useRouter();

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
    {
      title: "操作",
      key: "action",
      render: (_: unknown, row: QuizHistoryItem) => {
        const sessionId = row.sessionId || String(row.id);
        return (
          <Button
            type="link"
            onClick={() => router.push(`/quiz/history/${sessionId}`)}
          >
            查看详情
          </Button>
        );
      },
    },
  ];

  return (
    <Table
      rowKey={(record) => record.sessionId || String(record.id)}
      columns={columns}
      dataSource={items}
      loading={loading}
      pagination={false}
    />
  );
}
