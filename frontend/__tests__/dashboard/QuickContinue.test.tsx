import '@testing-library/jest-dom'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QuickContinue } from '@/app/(protected)/dashboard/components/QuickContinue'
import type { LastLearningRecord } from '@/types/dashboard'

// Mock next/navigation
const mockPush = jest.fn()
jest.mock('next/navigation', () => ({
	useRouter: () => ({
		push: mockPush,
	}),
}))

describe('QuickContinue', () => {
	beforeEach(() => {
		mockPush.mockClear()
	})

	it('应正确显示最后学习记录', () => {
		const lastLearning: LastLearningRecord = {
			topicId: 'variables',
			topicName: 'Variables',
			topicDisplayName: '变量',
			chapterId: 'storage',
			chapterName: 'Storage',
			chapterDisplayName: '存储',
			lastVisitedAt: '2025-12-26T10:30:00+08:00',
		}

		render(<QuickContinue lastLearning={lastLearning} />)

		// 有多个"继续学习"文本（标题和按钮），使用 getAllByText
		expect(screen.getAllByText('继续学习').length).toBeGreaterThan(0)
		expect(screen.getAllByText((content, element) => element?.textContent === '变量 - 存储')[0]).toBeInTheDocument()
		expect(screen.getByText(/最后访问：/)).toBeInTheDocument()
		// 按钮的可访问名称包含图标，所以使用包含匹配
		expect(screen.getByRole('button', { name: /继续学习/ })).toBeInTheDocument()
	})

	it('点击继续学习按钮应跳转到正确的章节页面', async () => {
		const lastLearning: LastLearningRecord = {
			topicId: 'constants',
			topicName: 'Constants',
			topicDisplayName: '常量',
			chapterId: 'iota',
			chapterName: 'Iota',
			chapterDisplayName: 'Iota',
			lastVisitedAt: '2025-12-26T10:30:00+08:00',
		}

		render(<QuickContinue lastLearning={lastLearning} />)

		// 按钮的可访问名称包含图标，所以使用包含匹配
		const continueButton = screen.getByRole('button', { name: /继续学习/ })
		await userEvent.click(continueButton)

		expect(mockPush).toHaveBeenCalledWith('/topics/constants/iota')
	})

	it('无学习记录时应显示空状态', () => {
		render(<QuickContinue lastLearning={null} />)

		expect(screen.getByText('继续学习')).toBeInTheDocument()
		expect(screen.getByText(/您还没有学习记录，快去开始学习吧！/)).toBeInTheDocument()
		expect(screen.getByRole('button', { name: '开始学习' })).toBeInTheDocument()
	})

	it('无学习记录时点击开始学习按钮应跳转到主题列表', async () => {
		render(<QuickContinue lastLearning={null} />)

		const startButton = screen.getByRole('button', { name: '开始学习' })
		await userEvent.click(startButton)

		expect(mockPush).toHaveBeenCalledWith('/topics')
	})

	it('应正确格式化最后访问时间', () => {
		const lastLearning: LastLearningRecord = {
			topicId: 'types',
			topicName: 'Types',
			topicDisplayName: '类型',
			chapterId: 'array',
			chapterName: 'Array',
			chapterDisplayName: '数组',
			lastVisitedAt: '2025-12-26T10:30:00+08:00',
		}

		render(<QuickContinue lastLearning={lastLearning} />)

		expect(screen.getByText(/最后访问：/)).toBeInTheDocument()
	})
})

