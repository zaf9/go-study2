'use client'

import { useEffect, useState } from 'react'
import { Result } from 'antd'
import { WelcomeHeader } from './components/WelcomeHeader'
import { StatsCards } from './components/StatsCards'
import { fetchDashboardStats } from '@/lib/dashboard'
import type { DashboardStats, ProgressUpdatedEventData } from '@/types/dashboard'
import { useAuth } from '@/hooks/useAuth'
import { useWebSocket } from '@/components/providers/WebSocketProvider'

export default function DashboardPage() {
	const { user } = useAuth()
	const { isConnected: wsConnected } = useWebSocket()
	const [stats, setStats] = useState<DashboardStats | null>(null)
	const [loading, setLoading] = useState(true)
	const [error, setError] = useState<Error | null>(null)

	useEffect(() => {
		if (!user?.id) {
			return
		}

		async function loadStats() {
			try {
				setLoading(true)
				setError(null)

				const data = await fetchDashboardStats()
				setStats(data)
			} catch (err) {
				console.error('加载 Dashboard 数据失败:', err)
				setError(err as Error)
			} finally {
				setLoading(false)
			}
		}

		loadStats()
	}, [user?.id])

	useEffect(() => {
		const handleMessage = (event: MessageEvent) => {
			try {
				const message = JSON.parse(event.data) as { event: string; data: ProgressUpdatedEventData }
				
				if (message.event === 'progress_updated' && stats) {
					console.log('[Dashboard] 收到进度更新:', message.data)
					loadStats()
				}
			} catch (err) {
				console.error('[Dashboard] 解析 WebSocket 消息失败:', err)
			}
		}

		window.addEventListener('websocket-message', handleMessage as EventListener)

		return () => {
			window.removeEventListener('websocket-message', handleMessage as EventListener)
		}
	}, [stats])

	if (loading) {
		return (
			<div className="flex items-center justify-center min-h-screen">
				<div className="text-center">
						<div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-gray-300 border-t-blue-500" />
					</div>
					<p className="mt-4 text-gray-500">加载中...</p>
				</div>
			</div>
		)
	}

	if (error) {
		return (
			<div className="flex items-center justify-center min-h-screen p-4">
				<Result
					status="error"
					title="加载失败"
					subTitle={error.message || '无法加载 Dashboard 数据，请稍后重试'}
					extra={
							<button
								onClick={() => {
										window.location.reload()
								}}
								className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-2 rounded"
							>
								重试
							</button>
					}
				/>
			</div>
		)
	}

	if (!user?.id) {
		return (
			<div className="flex items-center justify-center min-h-screen">
					<Result
						status="info"
						title="请先登录"
						subTitle="您需要先登录才能访问 Dashboard"
					/>
			</div>
		)
	}

	return (
		<div className="p-6">
			<WelcomeHeader
						username={user.username || '用户'}
						studyDays={stats?.studyDays || 0}
				/>
			{stats && (
					<StatsCards
							overallProgress={stats.progressPercentage}
							completedChapters={stats.completedChapters}
							totalChapters={stats.totalChapters}
							weeklyActivity={stats.weeklyActivity}
					/>
			)}
		</div>
	)
}
