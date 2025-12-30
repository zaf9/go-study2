'use client'

/**
 * RecentQuizzes 组件
 * 显示用户最近的测验记录
 */

import { Card, List, Typography, Tag, Empty } from 'antd'
import { useRouter } from 'next/navigation'
import type { RecentQuizSummary } from '@/types/dashboard'
import { formatTime } from '@/lib/utils/time'
import { TruncatedText } from './TruncatedText'

const { Title, Text } = Typography

interface RecentQuizzesProps {
    /** 最近测验记录列表 */
    quizzes: RecentQuizSummary[]
}

/**
 * RecentQuizzes 组件
 * 展示用户最近的测验记录，包括主题、章节、得分、完成时间等信息
 */
export const RecentQuizzes: React.FC<RecentQuizzesProps> = ({ quizzes }) => {
    const router = useRouter()

    // 处理点击，跳转到测验详情页
    const handleQuizClick = (quiz: RecentQuizSummary) => {
        // 使用查询参数传递 sessionId，避免静态导出时的路径参数限制
        router.push(`/quiz/review?sessionId=${quiz.sessionId}`)
    }

    // 根据得分获取颜色标签
    const getScoreColor = (score: number, total: number): string => {
        const percentage = (score / total) * 100
        if (percentage >= 80) {
            return 'success' // 绿色
        }
        if (percentage >= 60) {
            return 'warning' // 橙色
        }
        return 'error' // 红色
    }

    if (quizzes.length === 0) {
        return (
            <Card bordered={false} className="mb-6">
                <div className="flex justify-between items-center mb-4">
                    <Title level={4} className="mb-0">
                        最近测验
                    </Title>
                </div>
                <Empty description="暂无测验记录" />
            </Card>
        )
    }

    return (
        <Card bordered={false} className="mb-6">
            <div className="flex justify-between items-center mb-4">
                <Title level={4} className="mb-0">
                    最近测验
                </Title>
                <a 
                    onClick={() => router.push('/quiz-center')}
                    className="text-blue-500 hover:text-blue-700 cursor-pointer"
                >
                    查看全部 →
                </a>
            </div>
            <List
                dataSource={quizzes}
                renderItem={(quiz) => (
                    <List.Item
                        className="cursor-pointer hover:bg-gray-50 transition-colors"
                        onClick={() => handleQuizClick(quiz)}
                    >
                        <div className="w-full">
                            <div className="flex items-center justify-between mb-2 gap-2">
                                <div className="flex items-center gap-2 flex-1 min-w-0">
                                    <Text strong className="text-base">
                                        <TruncatedText text={quiz.topicName} maxLength={12} /> - <TruncatedText text={quiz.chapterName} maxLength={12} />
                                    </Text>
                                    {quiz.passed ? (
                                        <Tag color="success" className="flex-shrink-0">通过</Tag>
                                    ) : (
                                        <Tag color="error" className="flex-shrink-0">未通过</Tag>
                                    )}
                                </div>
                                <Text type="secondary" className="text-sm flex-shrink-0">
                                    {formatTime(quiz.completedAt)}
                                </Text>
                            </div>
                            <div className="flex items-center gap-2">
                                <Tag color={getScoreColor(quiz.score, quiz.totalQuestions)}>
                                    {quiz.score} / {quiz.totalQuestions}
                                </Tag>
                                <Text type="secondary" className="text-sm">
                                    {Math.round((quiz.score / quiz.totalQuestions) * 100)}%
                                </Text>
                            </div>
                        </div>
                    </List.Item>
                )}
            />
        </Card>
    )
}

