'use client'

/**
 * RecentQuizzes 组件
 * 显示用户最近的测验记录
 */

import { Card, List, Typography, Tag, Empty } from 'antd'
import { useRouter } from 'next/navigation'
import type { RecentQuizSummary } from '@/types/dashboard'
import { formatTime } from '@/lib/utils/time'

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

	// 处理点击，跳转到测验回顾页面（如果存在）
	const handleQuizClick = (quiz: RecentQuizSummary) => {
		// 跳转到测验回顾页面
		router.push(`/quiz/history/${quiz.id}`)
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
				<Title level={4} className="mb-4">
					最近测验
				</Title>
				<Empty description="暂无测验记录" />
			</Card>
		)
	}

	return (
		<Card bordered={false} className="mb-6">
			<Title level={4} className="mb-4">
				最近测验
			</Title>
			<List
				dataSource={quizzes}
				renderItem={(quiz) => (
					<List.Item
						className="cursor-pointer hover:bg-gray-50 transition-colors"
						onClick={() => handleQuizClick(quiz)}
					>
						<div className="w-full">
							<div className="flex items-center justify-between mb-2">
								<div className="flex items-center gap-2">
									<Text strong className="text-base">
										{quiz.topicName} - {quiz.chapterName}
									</Text>
									{quiz.passed ? (
										<Tag color="success">通过</Tag>
									) : (
										<Tag color="error">未通过</Tag>
									)}
								</div>
								<Text type="secondary" className="text-sm">
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

