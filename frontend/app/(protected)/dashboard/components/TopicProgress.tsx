'use client'

/**
 * TopicProgress 组件
 * 显示各主题的学习进度，使用进度条可视化展示完成度
 */

import { Card, Progress, Typography, List } from 'antd'
import { useRouter } from 'next/navigation'
import type { TopicProgressSummary } from '@/types/dashboard'
import { TruncatedText } from './TruncatedText'

const { Title, Text } = Typography

interface TopicProgressProps {
	/** 主题进度汇总列表 */
	topics: TopicProgressSummary[]
}

/**
 * TopicProgress 组件
 * 展示所有主题的学习进度，每个主题显示进度条和完成百分比
 */
export const TopicProgress: React.FC<TopicProgressProps> = ({ topics }) => {
	const router = useRouter()

	// 处理主题点击，跳转到主题详情页
	const handleTopicClick = (topicId: string) => {
		router.push(`/topics/${topicId}`)
	}

	// 根据完成百分比确定进度条颜色
	const getProgressColor = (percentage: number): string => {
		if (percentage >= 100) {
			return '#52c41a' // 绿色 - 已完成
		}
		if (percentage >= 50) {
			return '#1890ff' // 蓝色 - 进行中
		}
		if (percentage > 0) {
			return '#faad14' // 橙色 - 刚开始
		}
		return '#d9d9d9' // 灰色 - 未开始
	}

	if (topics.length === 0) {
		return (
			<Card bordered={false} className="mb-6">
				<Title level={4} className="mb-4">
					主题进度
				</Title>
				<Text type="secondary">暂无主题数据</Text>
			</Card>
		)
	}

	return (
		<Card bordered={false} className="mb-6">
			<Title level={4} className="mb-4">
				主题进度
			</Title>
			<List
				dataSource={topics}
				renderItem={(topic) => (
					<List.Item
						className="cursor-pointer hover:bg-gray-50 transition-colors"
						onClick={() => handleTopicClick(topic.topicId)}
					>
					<div className="w-full">
						<div className="flex items-center justify-between mb-2 gap-2">
							<Text strong className="text-base flex-1 min-w-0">
								<TruncatedText text={topic.displayName} maxLength={20} />
							</Text>
							<Text type="secondary" className="text-sm flex-shrink-0">
								{topic.completedChapters} / {topic.totalChapters} 章节
							</Text>
						</div>
							<Progress
								percent={topic.percentage}
								strokeColor={getProgressColor(topic.percentage)}
								showInfo={true}
								format={(percent) => `${percent?.toFixed(1)}%`}
							/>
						</div>
					</List.Item>
				)}
			/>
		</Card>
	)
}

