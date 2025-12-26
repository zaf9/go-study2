'use client'

/**
 * StatsCards 组件
 * 显示整体学习进度统计数据
 */

import { Card, Col, Row, Statistic } from 'antd'

interface StatsCardsProps {
	/** 整体完成百分比（0-100） */
	overallProgress: number

	/** 已完成章节数 */
	completedChapters: number

	/** 总章节数 */
	totalChapters: number

	/** 本周学习活动次数 */
	weeklyActivity: number
}

/**
 * StatsCards 组件
 * 展示整体学习进度统计数据
 */
export const StatsCards: React.FC<StatsCardsProps> = ({
	overallProgress,
	completedChapters,
	totalChapters,
	weeklyActivity,
}) => {
	// 确定进度条颜色
	const getProgressColor = (percent: number): string => {
		if (percent >= 100) {
			return '#52c41a'
		}
		if (percent >= 75) {
			return '#1890ff'
		}
		if (percent >= 50) {
			return '#faad14'
		}
		return '#ff4d4f'
	}

	const progressColor = getProgressColor(overallProgress)

	return (
		<Row gutter={[16, 16]}>
			<Col xs={24} sm={12} lg={8}>
				<Card bordered={false}>
					<Statistic
						title="整体进度"
						value={overallProgress}
						suffix="%"
						valueStyle={{
							color: progressColor,
							fontWeight: 600,
						}}
						prefix={
							<div
								style={{
									width: 24,
									height: 24,
									borderRadius: '50%',
									backgroundColor: progressColor,
									display: 'flex',
									alignItems: 'center',
									justifyContent: 'center',
								}}
							>
								<span
									style={{
										color: '#fff',
										fontSize: 14,
										fontWeight: 700,
									}}
								>
									%
								</span>
							</div>
						}
					/>
				</Card>
			</Col>
			<Col xs={24} sm={12} lg={8}>
				<Card bordered={false}>
					<Statistic
						title="完成章节"
						value={completedChapters}
						suffix={`/ ${totalChapters}`}
						valueStyle={{
							color: '#1890ff',
							fontWeight: 600,
						}}
					/>
				</Card>
			</Col>
			<Col xs={24} sm={12} lg={8}>
				<Card bordered={false}>
					<Statistic
						title="本周活跃"
						value={weeklyActivity}
						suffix=" 次"
						valueStyle={{
							color: '#52c41a',
							fontWeight: 600,
						}}
					/>
				</Card>
			</Col>
		</Row>
	)
}

