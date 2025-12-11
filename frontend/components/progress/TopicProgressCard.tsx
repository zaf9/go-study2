"use client";

import { Button, Card, Collapse, List, Space, Typography } from "antd";
import ChapterStatusIcon from "./ChapterStatusIcon";
import ProgressBar from "./ProgressBar";
import { ChapterProgress, TopicProgressDetail } from "@/types/learning";

/** 主题进度卡片属性：含主题摘要与章节列表，支持继续学习跳转。 */
interface TopicProgressCardProps {
  topic: TopicProgressDetail;
  onContinue?: (chapter: ChapterProgress) => void;
}

const { Text, Title } = Typography;

export default function TopicProgressCard({
  topic,
  onContinue,
}: TopicProgressCardProps) {
  const nextChapter =
    topic.chapters?.find((c) => c.status !== "completed") ?? null;

  return (
    <Card>
      <Space direction="vertical" className="w-full">
        <div className="flex items-center justify-between">
          <div>
            <Title level={4} className="mb-1">
              {topic.name} ({topic.id})
            </Title>
            <Text type="secondary">
              权重 {topic.weight} · 已完成 {topic.completedChapters} /
              {topic.totalChapters}
            </Text>
          </div>
          {nextChapter && onContinue && (
            <Button
              type="primary"
              onClick={() => onContinue(nextChapter)}
              size="small"
            >
              继续学习：{nextChapter.chapter}
            </Button>
          )}
        </div>
        <ProgressBar
          percent={topic.progress}
          status={topic.progress >= 100 ? "completed" : "in_progress"}
          segments={10}
          label="主题完成度"
        />
        <Collapse ghost size="small" items={[
          {
            key: "chapters",
            label: "章节状态",
            children: (
              <List
                dataSource={topic.chapters}
                renderItem={(item) => (
                  <List.Item
                    actions={
                      onContinue
                        ? [
                            <Button
                              type="link"
                              size="small"
                              key="goto"
                              onClick={() => onContinue(item)}
                            >
                              前往
                            </Button>,
                          ]
                        : []
                    }
                  >
                    <List.Item.Meta
                      avatar={<ChapterStatusIcon status={item.status} />}
                      title={item.chapter}
                      description={
                        <Text type="secondary">
                          阅读 {item.scrollProgress ?? 0}% ·{" "}
                          {item.quizPassed ? "测验通过" : "测验未通过"}
                        </Text>
                      }
                    />
                  </List.Item>
                )}
              />
            ),
          },
        ]} />
      </Space>
    </Card>
  );
}

