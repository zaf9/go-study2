'use client'

/**
 * QuickContinue 组件
 * 显示最后一次学习的主题和章节，提供一键继续学习功能
 */

import { Card, Button, Typography, Space } from 'antd'
import { RightOutlined } from '@ant-design/icons'
import { useRouter } from 'next/navigation'
import type { LastLearningRecord } from '@/types/dashboard'
import { TruncatedText } from './TruncatedText'

const { Title, Text } = Typography

interface QuickContinueProps {
	/** 最后学习记录，为 null 时显示空状态 */
	lastLearning: LastLearningRecord | null
}

/**
 * QuickContinue 组件
 * 展示最后学习记录并提供继续学习按钮
 */
export const QuickContinue: React.FC<QuickContinueProps> = ({ lastLearning }) => {
	const router = useRouter()

	// 处理继续学习按钮点击
	const handleContinue = () => {
		if (lastLearning) {
			router.push(`/topics/${lastLearning.topicId}/${lastLearning.chapterId}`)
		}
	}

	// 空状态：无学习记录
	if (!lastLearning) {
		return (
			<Card bordered={false} className="mb-6">
				<Space direction="vertical" size="small" className="w-full">
					<Title level={4} className="mb-0">
						继续学习
					</Title>
					<Text type="secondary">您还没有学习记录，快去开始学习吧！</Text>
					<Button type="primary" onClick={() => router.push('/topics')}>
						开始学习
					</Button>
				</Space>
			</Card>
		)
	}

	// 有学习记录：显示最后学习的主题和章节
	return (
		<Card bordered={false} className="mb-6">
			<Space direction="vertical" size="small" className="w-full">
				<Title level={4} className="mb-0">
					继续学习
				</Title>
				<Space direction="vertical" size={4} className="w-full">
					<Text strong className="text-base">
						<TruncatedText text={lastLearning.topicDisplayName} maxLength={15} /> - <TruncatedText text={lastLearning.chapterDisplayName} maxLength={15} />
					</Text>
					<Text type="secondary" className="text-sm">
						最后访问：{new Date(lastLearning.lastVisitedAt).toLocaleString('zh-CN')}
					</Text>
				</Space>
				<Button
					type="primary"
					size="large"
					icon={<RightOutlined />}
					onClick={handleContinue}
					className="mt-2"
				>
					继续学习
				</Button>
			</Space>
		</Card>
	)
}

