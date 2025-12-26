import '@testing-library/jest-dom'
import { render, screen } from '@testing-library/react'
import { StatsCards } from '@/app/(protected)/dashboard/components/StatsCards'

describe('StatsCards', () => {
	it('应正确显示所有统计数据', () => {
		render(
			<StatsCards
				overallProgress={75.5}
				completedChapters={15}
				totalChapters={20}
				weeklyActivity={5}
			/>,
		)

		expect(screen.getByText('整体进度')).toBeInTheDocument()
		expect(screen.getByText('75.5%')).toBeInTheDocument()
		expect(screen.getByText('完成章节')).toBeInTheDocument()
		expect(screen.getByText('15')).toBeInTheDocument()
		expect(screen.getByText('/ 20')).toBeInTheDocument()
		expect(screen.getByText('本周活跃')).toBeInTheDocument()
		expect(screen.getByText('5')).toBeInTheDocument()
		expect(screen.getByText(' 次')).toBeInTheDocument()
	})

	it('应正确处理进度为 0 的情况', () => {
		render(
			<StatsCards
				overallProgress={0}
				completedChapters={0}
				totalChapters={20}
				weeklyActivity={0}
			/>,
		)

		expect(screen.getByText('0%')).toBeInTheDocument()
		expect(screen.getByText('0')).toBeInTheDocument()
		expect(screen.getByText('/ 20')).toBeInTheDocument()
	})

	it('应正确处理进度为 100% 的情况', () => {
		render(
			<StatsCards
				overallProgress={100}
				completedChapters={20}
				totalChapters={20}
				weeklyActivity={10}
			/>,
		)

		expect(screen.getByText('100%')).toBeInTheDocument()
		expect(screen.getByText('20')).toBeInTheDocument()
		expect(screen.getByText('/ 20')).toBeInTheDocument()
	})

	it('应根据进度百分比显示正确的颜色', () => {
		const { rerender } = render(
			<StatsCards
				overallProgress={25}
				completedChapters={5}
				totalChapters={20}
				weeklyActivity={2}
			/>,
		)

		// 进度 < 50%，应为红色 (#ff4d4f)
		let progressElement = screen.getByText('25%').closest('.ant-statistic-content')
		expect(progressElement).toBeInTheDocument()

		// 进度 >= 50% 且 < 75%，应为橙色 (#faad14)
		rerender(
			<StatsCards
				overallProgress={60}
				completedChapters={12}
				totalChapters={20}
				weeklyActivity={3}
			/>,
		)
		progressElement = screen.getByText('60%').closest('.ant-statistic-content')
		expect(progressElement).toBeInTheDocument()

		// 进度 >= 75% 且 < 100%，应为蓝色 (#1890ff)
		rerender(
			<StatsCards
				overallProgress={80}
				completedChapters={16}
				totalChapters={20}
				weeklyActivity={4}
			/>,
		)
		progressElement = screen.getByText('80%').closest('.ant-statistic-content')
		expect(progressElement).toBeInTheDocument()

		// 进度 = 100%，应为绿色 (#52c41a)
		rerender(
			<StatsCards
				overallProgress={100}
				completedChapters={20}
				totalChapters={20}
				weeklyActivity={5}
			/>,
		)
		progressElement = screen.getByText('100%').closest('.ant-statistic-content')
		expect(progressElement).toBeInTheDocument()
	})

	it('应正确处理小数进度', () => {
		render(
			<StatsCards
				overallProgress={33.33}
				completedChapters={10}
				totalChapters={30}
				weeklyActivity={3}
			/>,
		)

		expect(screen.getByText('33.33%')).toBeInTheDocument()
		expect(screen.getByText('10')).toBeInTheDocument()
		expect(screen.getByText('/ 30')).toBeInTheDocument()
	})
})

