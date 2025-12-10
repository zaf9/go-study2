'use client';

import { List, Button, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import { ChapterSummary } from "@/types/learning";

interface ChapterListProps {
  topicKey: string;
  chapters: ChapterSummary[];
}

const { Text } = Typography;

export default function ChapterList({ topicKey, chapters }: ChapterListProps) {
  const router = useRouter();

  return (
    <List
      itemLayout="horizontal"
      dataSource={chapters}
      bordered
      renderItem={(item) => (
        <List.Item
          actions={[
            <Button key="view" type="link" onClick={() => router.push(`/topics/${topicKey}/${item.id}`)}>
              查看
            </Button>,
          ]}
        >
          <List.Item.Meta
            title={item.title}
            description={
              <Space size="small">
                <Text type="secondary">{item.summary || "章节内容"}</Text>
                {typeof item.order === "number" && <Text type="secondary">序号：{item.order + 1}</Text>}
              </Space>
            }
          />
        </List.Item>
      )}
    />
  );
}


