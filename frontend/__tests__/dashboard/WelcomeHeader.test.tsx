import '@testing-library/jest-dom'
import { render, screen } from '@testing-library/react'
import { WelcomeHeader } from '@/app/(protected)/dashboard/components/WelcomeHeader'

describe('WelcomeHeader', () => {
	it('应正确显示用户名和累计学习天数', () => {
		render(<WelcomeHeader username="测试用户" studyDays={15} />)

		expect(screen.getByText(/欢迎回来，测试用户！/)).toBeInTheDocument()
		expect(screen.getByText(/您已累计学习 15 天/)).toBeInTheDocument()
	})

	it('应正确处理学习天数为 0 的情况', () => {
		render(<WelcomeHeader username="新用户" studyDays={0} />)

		expect(screen.getByText(/欢迎回来，新用户！/)).toBeInTheDocument()
		expect(screen.getByText(/您已累计学习 0 天/)).toBeInTheDocument()
	})

	it('应正确处理较长的用户名', () => {
		const longUsername = '这是一个非常长的用户名用于测试显示效果'
		render(<WelcomeHeader username={longUsername} studyDays={30} />)

		expect(screen.getByText(new RegExp(`欢迎回来，${longUsername}！`))).toBeInTheDocument()
		expect(screen.getByText(/您已累计学习 30 天/)).toBeInTheDocument()
	})

	it('应正确处理较大的学习天数', () => {
		render(<WelcomeHeader username="老用户" studyDays={365} />)

		expect(screen.getByText(/欢迎回来，老用户！/)).toBeInTheDocument()
		expect(screen.getByText(/您已累计学习 365 天/)).toBeInTheDocument()
	})
})

