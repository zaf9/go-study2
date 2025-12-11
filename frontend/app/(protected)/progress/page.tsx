"use client";

import { Table, Tag, Typography } from "antd";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import useProgress from "@/hooks/useProgress";

const { Title } = Typography;

export default function ProgressPage() {
  const { progress, isLoading, error } = useProgress();

  if (isLoading) return <Loading />;
  if (error)
    return <ErrorMessage message="加载进度失败" description={error.message} />;

  const columns = [
    { title: "主题", dataIndex: "topic", key: "topic" },
    { title: "章节", dataIndex: "chapter", key: "chapter" },
    {
      title: "状态",
      dataIndex: "status",
      key: "status",
      render: (v: string) => (
        <Tag color={v === "done" ? "green" : "blue"}>{statusLabel(v)}</Tag>
      ),
    },
    {
      title: "最近访问",
      dataIndex: "lastVisit",
      key: "lastVisit",
      render: (v: string) => new Date(v).toLocaleString(),
    },
  ];

  return (
    <div className="space-y-4">
      <Title level={3}>学习进度</Title>
      <Table
        rowKey={(row) => `${row.topic}-${row.chapter}`}
        columns={columns}
        dataSource={progress}
        pagination={false}
      />
    </div>
  );
}

function statusLabel(status: string) {
  if (status === "done") return "已完成";
  if (status === "in_progress") return "学习中";
  return "未开始";
}
