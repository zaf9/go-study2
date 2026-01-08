"use client";

import { Select, Space, Typography } from "antd";
import { useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import Loading from "@/components/common/Loading";
import ErrorMessage from "@/components/common/ErrorMessage";
import ProgressOverview from "@/components/progress/ProgressOverview";
import TopicProgressCard from "@/components/progress/TopicProgressCard";
import { topicChapters } from "@/lib/static-routes";
import { TopicProgressDetail } from "@/types/learning";
import useProgress from "@/hooks/useProgress";
import useTopicProgressDetail from "@/hooks/useTopicProgressDetail";
import PageContainer from "@/components/layout/PageContainer";

const { Title } = Typography;

export default function ProgressPage() {
    const router = useRouter();
    const { overview, next, isLoading, error } = useProgress();
    const [selectedTopic, setSelectedTopic] = useState<string | undefined>();

    // 为每个主题获取完整的章节详情
    const topicIds = useMemo(
        () => (overview?.topics ?? [])
            .map(t => t.id)
            .filter((id) => !!id && typeof id === 'string'), // 过滤掉undefined和空字符串
        [overview?.topics]
    );

    const topicDetailsMap = useTopicProgressDetail(topicIds);

    const topics = useMemo(
        () =>
            (overview?.topics ?? [])
                .map<TopicProgressDetail>((item) => ({
                    id: item.id,
                    name: item.name,
                    weight: item.weight ?? 0,
                    progress: item.progress ?? 0,
                    totalChapters: item.totalChapters,
                    completedChapters: item.completedChapters,
                    chapters: topicDetailsMap[item.id]?.chapters ?? [],
                }))
                .sort((a, b) => b.progress - a.progress),
        [overview?.topics, topicDetailsMap],
    );

    const filteredTopics = useMemo(() => {
        if (!selectedTopic) return topics;
        return topics.filter((t) => t.id === selectedTopic);
    }, [selectedTopic, topics]);

    if (isLoading) return <Loading />;
    if (error)
        return <ErrorMessage message="加载进度失败" description={error.message} />;
    if (!overview) return null;

    // 处理新用户空数据场景
    const hasNoProgress = !overview.topics || overview.topics.length === 0;
    // 只有当没有进度数据且没有总章节数据时才认为是全新用户
    const isNewUser = hasNoProgress && overview.overall.completedChapters === 0 && overview.overall.totalChapters === 0;

    if (isNewUser) {
        return (
            <PageContainer>
                <Space direction="vertical" className="w-full items-center justify-center min-h-[400px]">
                    <div className="text-center">
                        <div className="text-6xl mb-4">📚</div>
                        <Title level={3}>欢迎开始学习之旅！</Title>
                        <p className="text-gray-600 mb-6">
                            您还没有任何学习记录，从主题列表选择感兴趣的内容开始学习吧
                        </p>
                        <button
                            onClick={() => router.push('/topics')}
                            className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                        >
                            开始学习
                        </button>
                    </div>
                </Space>
            </PageContainer>
        );
    }

    return (
        <PageContainer>
            <Space direction="vertical" className="w-full" size="large">
                <Title level={3}>学习进度</Title>
                <ProgressOverview
                    overall={overview.overall}
                    next={next}
                    onContinue={(hint) =>
                        router.push(`/topics/${hint.topic}/${hint.chapter}`)
                    }
                />
                <div>
                    <Title level={4} className="mb-4">
                        主题进度
                    </Title>
                    <div className="flex items-center justify-between mb-4">
                        <Select
                            allowClear
                            placeholder="筛选主题"
                            value={selectedTopic}
                            onChange={(v) => setSelectedTopic(v)}
                            options={topics.map((t) => ({ label: t.name, value: t.id }))}
                            className="min-w-[200px]"
                        />
                    </div>
                    <Space direction="vertical" className="w-full" size="middle">
                        {filteredTopics.map((topic) => (
                            <TopicProgressCard
                                key={topic.id}
                                topic={topic}
                                onContinue={(chapter) => {
                                    const target =
                                        chapter?.chapter ??
                                        topicChapters[topic.id as keyof typeof topicChapters]?.[0];
                                    if (target) {
                                        router.push(`/topics/${topic.id}/${target}`);
                                    }
                                }}
                            />
                        ))}
                    </Space>
                </div>
            </Space>
        </PageContainer>
    );
}
