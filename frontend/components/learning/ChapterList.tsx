"use client";

import { List, Button, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import { ChapterSummary } from "@/types/learning";
import { ChapterProgress } from "@/types/learning";
import ChapterStatusBadge from "./ChapterStatusBadge";

interface ChapterListProps {
  topicKey: string;
  chapters: ChapterSummary[];
  progress?: ChapterProgress[];
}

const { Text } = Typography;

export default function ChapterList({
  topicKey,
  chapters,
  progress,
}: ChapterListProps) {
  const router = useRouter();
  const progressMap =
    progress?.reduce<Record<string, ChapterProgress>>((map, item) => {
      map[item.chapter] = item;
      return map;
    }, {}) ?? {};

  return (
    <List
      itemLayout="horizontal"
      dataSource={chapters}
      bordered
      renderItem={(item) => {
        const prog = progressMap[item.id];
        return (
          <List.Item
            actions={[
              <Button
                key="view"
                type="link"
                onClick={() => router.push(`/topics/${topicKey}/${item.id}`)}
              >
                查看
              </Button>,
            ]}
          >
            <List.Item.Meta
              title={item.title}
              description={
                <Space size="small">
                  <Text type="secondary">{item.summary || "章节内容"}</Text>
                  {typeof item.order === "number" && (
                    <Text type="secondary">序号：{item.order + 1}</Text>
                  )}
                  <ChapterStatusBadge
                    progressData={prog}
                  />
                </Space>
              }
            />
          </List.Item>
        );
      }}
    />
  );
}
