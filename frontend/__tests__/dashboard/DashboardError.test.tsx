import '@testing-library/jest-dom'
import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import DashboardError from '@/app/(protected)/dashboard/error'

// Mock window.location.reload
const mockReload = jest.fn()
Object.defineProperty(window, 'location', {
	value: {
		reload: mockReload,
		href: '/',
	},
	writable: true,
})

describe('DashboardError', () => {
	beforeEach(() => {
		mockReload.mockClear()
	})

	it('应显示网络错误消息和重试按钮', () => {
		const error = new Error('网络错误：Failed to fetch')
		const reset = jest.fn()

		render(<DashboardError error={error} reset={reset} />)

		expect(screen.getByText('网络错误')).toBeInTheDocument()
		expect(screen.getByText(/网络错误：Failed to fetch/)).toBeInTheDocument()
		expect(screen.getByRole('button', { name: /重\s*试/ })).toBeInTheDocument()
		expect(screen.getByRole('button', { name: '返回首页' })).toBeInTheDocument()
	})

	it('应显示服务器错误消息', () => {
		const error = new Error('服务器错误：Internal Server Error')
		const reset = jest.fn()

		render(<DashboardError error={error} reset={reset} />)

		expect(screen.getByText('服务器错误')).toBeInTheDocument()
		expect(screen.getByText(/服务器错误：Internal Server Error/)).toBeInTheDocument()
	})

	it('应显示通用错误消息', () => {
		const error = new Error('未知错误')
		const reset = jest.fn()

		render(<DashboardError error={error} reset={reset} />)

		expect(screen.getByText('加载失败')).toBeInTheDocument()
		expect(screen.getByText(/未知错误/)).toBeInTheDocument()
	})

	it('点击重试按钮应调用 reset 函数', async () => {
		const error = new Error('网络错误')
		const reset = jest.fn()

		render(<DashboardError error={error} reset={reset} />)

		const retryButton = await screen.findByRole('button', { name: /重\s*试/ })
		
		await userEvent.click(retryButton)

		expect(reset).toHaveBeenCalled()
	})

	it('点击返回首页按钮应跳转到首页', async () => {
		const error = new Error('网络错误')
		const reset = jest.fn()
		
		// Mock window.location.href
		delete (window as any).location
		;(window as any).location = { href: '/' }

		render(<DashboardError error={error} reset={reset} />)

		const homeButton = await screen.findByRole('button', { name: '返回首页' })
		
		await userEvent.click(homeButton)

		expect(window.location.href).toBe('/')
	})

	it('网络错误应在2秒后自动重试', async () => {
		jest.useFakeTimers()
		const error = new Error('网络错误')
		const reset = jest.fn()

		render(<DashboardError error={error} reset={reset} />)

		// 初始不应调用 reset
		expect(reset).not.toHaveBeenCalled()

		// 快进2秒
		act(() => {
			jest.advanceTimersByTime(2000)
		})

		expect(reset).toHaveBeenCalledTimes(1)
		
		jest.useRealTimers()
	})

	it('非网络错误不应自动重试', async () => {
		jest.useFakeTimers()
		const error = new Error('其他错误')
		const reset = jest.fn()

		render(<DashboardError error={error} reset={reset} />)

		// 快进5秒
		jest.advanceTimersByTime(5000)

		// 不应自动调用 reset
		expect(reset).not.toHaveBeenCalled()
		
		jest.useRealTimers()
	})

	it('应显示错误提示信息', () => {
		const error = new Error('网络错误：Failed to fetch')
		const reset = jest.fn()

		render(<DashboardError error={error} reset={reset} />)

		expect(
			screen.getByText(/请检查您的网络连接，或稍后再试/),
		).toBeInTheDocument()
	})
})

