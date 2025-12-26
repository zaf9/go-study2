'use client'

/**
 * WelcomeHeader 组件
 * 显示欢迎信息和累计学习天数
 */

import { Typography } from 'antd'

const { Title, Text } = Typography

interface WelcomeHeaderProps {
	/** 用户名 */
	username: string

	/** 累计学习天数 */
	studyDays: number
}

/**
 * WelcomeHeader 组件
 * 展示欢迎信息和累计学习天数
 */
export const WelcomeHeader: React.FC<WelcomeHeaderProps> = ({
	username,
	studyDays,
}) => {
	return (
		<div className="mb-6">
			<Title level={2} className="mb-2">
				欢迎回来，{username}！
			</Title>
			<Text type="secondary" className="text-base">
				您已累计学习 {studyDays} 天
			</Text>
		</div>
	)
}

