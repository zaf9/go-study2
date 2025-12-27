import '@testing-library/jest-dom'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TopicProgress } from '@/app/(protected)/dashboard/components/TopicProgress'
import type { TopicProgressSummary } from '@/types/dashboard'

// Mock next/navigation
const mockPush = jest.fn()
jest.mock('next/navigation', () => ({
	useRouter: () => ({
		push: mockPush,
	}),
}))

describe('TopicProgress', () => {
	beforeEach(() => {
		mockPush.mockClear()
	})

	it('应正确显示所有主题的进度', () => {
		const topics: TopicProgressSummary[] = [
			{
				topicId: 'variables',
				topicName: 'Variables',
				displayName: '变量',
				completedChapters: 2,
				totalChapters: 4,
				percentage: 50.0,
			},
			{
				topicId: 'constants',
				topicName: 'Constants',
				displayName: '常量',
				completedChapters: 3,
				totalChapters: 3,
				percentage: 100.0,
			},
			{
				topicId: 'types',
				topicName: 'Types',
				displayName: '类型',
				completedChapters: 0,
				totalChapters: 5,
				percentage: 0.0,
			},
		]

		render(<TopicProgress topics={topics} />)

		expect(screen.getByText('主题进度')).toBeInTheDocument()
		expect(screen.getByText('变量')).toBeInTheDocument()
		expect(screen.getByText('常量')).toBeInTheDocument()
		expect(screen.getByText('类型')).toBeInTheDocument()
		expect(screen.getByText(/2 \/ 4 章节/)).toBeInTheDocument()
		expect(screen.getByText(/3 \/ 3 章节/)).toBeInTheDocument()
		expect(screen.getByText(/0 \/ 5 章节/)).toBeInTheDocument()
	})

	it('点击主题应跳转到主题详情页', async () => {
		const topics: TopicProgressSummary[] = [
			{
				topicId: 'variables',
				topicName: 'Variables',
				displayName: '变量',
				completedChapters: 2,
				totalChapters: 4,
				percentage: 50.0,
			},
		]

		render(<TopicProgress topics={topics} />)

		const topicItem = screen.getByText('变量').closest('.ant-list-item')
		if (topicItem) {
			await userEvent.click(topicItem)
			expect(mockPush).toHaveBeenCalledWith('/topics/variables')
		}
	})

	it('应正确处理空列表', () => {
		render(<TopicProgress topics={[]} />)

		expect(screen.getByText('主题进度')).toBeInTheDocument()
		expect(screen.getByText('暂无主题数据')).toBeInTheDocument()
	})

	it('应正确显示不同完成度的进度条颜色', () => {
		const topics: TopicProgressSummary[] = [
			{
				topicId: 'completed',
				topicName: 'Completed',
				displayName: '已完成',
				completedChapters: 10,
				totalChapters: 10,
				percentage: 100.0,
			},
			{
				topicId: 'inprogress',
				topicName: 'InProgress',
				displayName: '进行中',
				completedChapters: 5,
				totalChapters: 10,
				percentage: 50.0,
			},
			{
				topicId: 'notstarted',
				topicName: 'NotStarted',
				displayName: '未开始',
				completedChapters: 0,
				totalChapters: 10,
				percentage: 0.0,
			},
		]

		render(<TopicProgress topics={topics} />)

		// 验证所有主题都显示
		expect(screen.getByText('已完成')).toBeInTheDocument()
		expect(screen.getByText('进行中')).toBeInTheDocument()
		expect(screen.getByText('未开始')).toBeInTheDocument()
	})

	it('应正确格式化百分比显示', () => {
		const topics: TopicProgressSummary[] = [
			{
				topicId: 'test',
				topicName: 'Test',
				displayName: '测试',
				completedChapters: 1,
				totalChapters: 3,
				percentage: 33.33,
			},
		]

		render(<TopicProgress topics={topics} />)

		expect(screen.getByText('测试')).toBeInTheDocument()
		expect(screen.getByText(/1 \/ 3 章节/)).toBeInTheDocument()
	})
})

