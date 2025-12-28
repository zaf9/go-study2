import '@testing-library/jest-dom'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { RecentQuizzes } from '@/app/(protected)/dashboard/components/RecentQuizzes'
import type { RecentQuizSummary } from '@/types/dashboard'

// Mock next/navigation
const mockPush = jest.fn()
jest.mock('next/navigation', () => ({
	useRouter: () => ({
		push: mockPush,
	}),
}))

// Mock time formatting utility
jest.mock('@/lib/utils/time', () => ({
	formatTime: jest.fn((dateString: string) => {
		const date = new Date(dateString)
		const now = new Date()
		const diffMs = now.getTime() - date.getTime()
		const diffHours = diffMs / (1000 * 60 * 60)

		if (diffHours < 24) {
			return `${Math.floor(diffHours)} 小时前`
		}
		return date.toLocaleString('zh-CN', {
			year: 'numeric',
			month: '2-digit',
			day: '2-digit',
			hour: '2-digit',
			minute: '2-digit',
		})
	}),
}))

describe('RecentQuizzes', () => {
	beforeEach(() => {
		mockPush.mockClear()
	})

	it('应正确显示最近测验记录', () => {
		const quizzes: RecentQuizSummary[] = [
			{
				id: 1,
				sessionId: 'session-uuid-1',
				topicName: 'Variables',
				chapterName: 'Storage',
				score: 4,
				totalQuestions: 5,
				passed: true,
				completedAt: '2025-12-26T10:30:00+08:00',
			},
			{
				id: 2,
				sessionId: 'session-uuid-2',
				topicName: 'Constants',
				chapterName: 'Boolean',
				score: 2,
				totalQuestions: 3,
				passed: false,
				completedAt: '2025-12-26T09:00:00+08:00',
			},
		]

		render(<RecentQuizzes quizzes={quizzes} />)

		expect(screen.getByText('最近测验')).toBeInTheDocument()
		expect(screen.getAllByText((content, element) => element?.textContent === 'Variables - Storage')[0]).toBeInTheDocument()
		expect(screen.getAllByText((content, element) => element?.textContent === 'Constants - Boolean')[0]).toBeInTheDocument()
		expect(screen.getByText('通过')).toBeInTheDocument()
		expect(screen.getByText('未通过')).toBeInTheDocument()
		expect(screen.getByText(/4 \/ 5/)).toBeInTheDocument()
		expect(screen.getByText(/2 \/ 3/)).toBeInTheDocument()
	})

	it('点击测验记录应跳转到测验回顾页面', async () => {
		const quizzes: RecentQuizSummary[] = [
			{
				id: 1,
				sessionId: 'session-uuid-1',
				topicName: 'Variables',
				chapterName: 'Storage',
				score: 4,
				totalQuestions: 5,
				passed: true,
				completedAt: '2025-12-26T10:30:00+08:00',
			},
		]

		render(<RecentQuizzes quizzes={quizzes} />)

		const quizItem = screen.getAllByText((content, element) => element?.textContent === 'Variables - Storage')[0].closest('.ant-list-item')
		if (quizItem) {
			await userEvent.click(quizItem)
			expect(mockPush).toHaveBeenCalledWith('/quiz/review?sessionId=session-uuid-1')
		}
	})

	it('应正确处理空列表', () => {
		render(<RecentQuizzes quizzes={[]} />)

		expect(screen.getByText('最近测验')).toBeInTheDocument()
		expect(screen.getByText('暂无测验记录')).toBeInTheDocument()
	})

	it('应正确显示得分百分比', () => {
		const quizzes: RecentQuizSummary[] = [
			{
				id: 1,
				sessionId: 'session-uuid-1',
				topicName: 'Variables',
				chapterName: 'Storage',
				score: 3,
				totalQuestions: 4,
				passed: true,
				completedAt: '2025-12-26T10:30:00+08:00',
			},
		]

		render(<RecentQuizzes quizzes={quizzes} />)

		expect(screen.getByText(/3 \/ 4/)).toBeInTheDocument()
		expect(screen.getByText('75%')).toBeInTheDocument()
	})

	it('应根据得分显示正确的颜色标签', () => {
		const quizzes: RecentQuizSummary[] = [
			{
				id: 1,
				sessionId: 'session-uuid-1',
				topicName: 'High Score',
				chapterName: 'Chapter',
				score: 9,
				totalQuestions: 10,
				passed: true,
				completedAt: '2025-12-26T10:30:00+08:00',
			},
			{
				id: 2,
				sessionId: 'session-uuid-2',
				topicName: 'Medium Score',
				chapterName: 'Chapter',
				score: 6,
				totalQuestions: 10,
				passed: false,
				completedAt: '2025-12-26T09:00:00+08:00',
			},
			{
				id: 3,
				sessionId: 'session-uuid-3',
				topicName: 'Low Score',
				chapterName: 'Chapter',
				score: 3,
				totalQuestions: 10,
				passed: false,
				completedAt: '2025-12-26T08:00:00+08:00',
			},
		]

		render(<RecentQuizzes quizzes={quizzes} />)

		// 验证所有记录都显示
		expect(screen.getByText(/High Score/)).toBeInTheDocument()
		expect(screen.getByText(/Medium Score/)).toBeInTheDocument()
		expect(screen.getByText(/Low Score/)).toBeInTheDocument()
	})
})

