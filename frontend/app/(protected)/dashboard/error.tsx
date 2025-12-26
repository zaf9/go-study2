'use client'

/**
 * Dashboard 错误边界
 * 处理 Dashboard 页面的错误状态，提供重试功能
 */

import { useEffect } from 'react'
import { Result, Button } from 'antd'

/**
 * Error Boundary 参数
 */
interface ErrorProps {
	error: Error
	reset: () => void
}

/**
 * Dashboard 错误边界组件
 * 处理页面错误，提供重试功能（FR-021）
 */
export default function DashboardError({ error, reset }: ErrorProps) {
	// 自动重试（当遇到临时性错误时）
	useEffect(() => {
		const isTemporaryError =
			error.message?.includes('网络') ||
			error.message?.includes('超时') ||
			error.message?.includes('Failed to fetch')

		if (isTemporaryError) {
			const timer = setTimeout(() => {
				console.log('[Dashboard Error] 自动重试...')
				reset()
			}, 2000) // 2 秒后自动重试

			return () => clearTimeout(timer)
		}
	}, [error, reset])

	// 判断错误类型
	const isNetworkError =
		error.message?.includes('网络') ||
		error.message?.includes('fetch') ||
		error.message?.includes('Failed to fetch')

	const isServerError =
		error.message?.includes('服务器') ||
		error.message?.includes('Internal Server Error')

	const status = isServerError ? '500' : isNetworkError ? 'error' : 'warning'
	const title = isServerError ? '服务器错误' : isNetworkError ? '网络错误' : '加载失败'

	return (
		<div className="flex min-h-screen items-center justify-center bg-gray-50 p-4">
				<Result
						status={status as any}
						title={title}
						subTitle={
								<div className="max-w-md text-center">
										<p className="text-base mb-2">
												{error.message || '无法加载 Dashboard 数据，请检查网络连接或稍后重试'}
										</p>
										{isNetworkError && (
												<p className="text-sm text-gray-500">
														提示：请检查您的网络连接，或稍后再试
												</p>
										)}
								</div>
						}
						extra={
								<div className="flex flex-col gap-3">
										<Button
												type="primary"
												size="large"
												onClick={() => {
														console.log('[Dashboard Error] 用户手动重试')
														reset()
												}}
												className="w-full sm:w-auto"
										>
												重试
										</Button>
										<Button
												size="large"
												onClick={() => {
														console.log('[Dashboard Error] 返回首页')
														window.location.href = '/'
												}}
												className="w-full sm:w-auto"
										>
												返回首页
										</Button>
								</div>
						}
				/>
		</div>
	)
}

