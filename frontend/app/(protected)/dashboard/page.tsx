'use client'

import { useEffect, useState, lazy, Suspense } from 'react'
import { Result, Spin } from 'antd'
import { WelcomeHeader } from './components/WelcomeHeader'
import { StatsCards } from './components/StatsCards'
import { QuickContinue } from './components/QuickContinue'
import { fetchDashboardStats, fetchLastLearning, fetchTopicProgress, fetchRecentQuizzes } from '@/lib/dashboard'
import type { DashboardStats, ProgressUpdatedEventData, QuizCompletedEventData, LastLearningRecord, TopicProgressSummary, RecentQuizSummary } from '@/types/dashboard'
import useAuth from '@/hooks/useAuth'

// 懒加载非关键组件以实现代码分割
const TopicProgress = lazy(() => import('./components/TopicProgress').then(module => ({ default: module.TopicProgress })))
const RecentQuizzes = lazy(() => import('./components/RecentQuizzes').then(module => ({ default: module.RecentQuizzes })))

export default function DashboardPage() {
	const { user } = useAuth()
	const [stats, setStats] = useState<DashboardStats | null>(null)
	const [lastLearning, setLastLearning] = useState<LastLearningRecord | null>(null)
	const [topicProgress, setTopicProgress] = useState<TopicProgressSummary[]>([])
	const [recentQuizzes, setRecentQuizzes] = useState<RecentQuizSummary[]>([])
	const [loading, setLoading] = useState(true)
	const [error, setError] = useState<Error | null>(null)

	useEffect(() => {
		if (!user?.id) {
			return
		}

		async function loadData(skipLoading = false) {
			try {
				if (!skipLoading) {
					setLoading(true)
				}
				setError(null)

				// 并行加载统计数据、最后学习记录、主题进度和最近测验
				const [statsData, lastLearningData, topicProgressData, recentQuizzesData] = await Promise.all([
					fetchDashboardStats(),
					fetchLastLearning().catch(() => null), // 如果获取失败，返回 null
					fetchTopicProgress().catch(() => []), // 如果获取失败，返回空数组
					fetchRecentQuizzes(5).catch(() => []), // 如果获取失败，返回空数组
				])

				setStats(statsData)
				setLastLearning(lastLearningData)
				setTopicProgress(topicProgressData)
				setRecentQuizzes(recentQuizzesData)
			} catch (err) {
				console.error('加载 Dashboard 数据失败:', err)
				setError(err as Error)
			} finally {
				if (!skipLoading) {
					setLoading(false)
				}
			}
		}

		loadData()

		const handleMessage = (event: MessageEvent) => {
			try {
				const message = JSON.parse(event.data) as { event: string; data: ProgressUpdatedEventData | QuizCompletedEventData }
				
				if (message.event === 'progress_updated') {
					console.log('[Dashboard] 收到进度更新:', message.data)
					// 重新加载所有数据以更新主题进度（跳过 loading 状态）
					loadData(true)
				} else if (message.event === 'quiz_completed') {
					console.log('[Dashboard] 收到测验完成事件:', message.data)
					// 重新加载最近测验数据（跳过 loading 状态）
					fetchRecentQuizzes(5)
						.then((quizzes) => {
							setRecentQuizzes(quizzes)
						})
						.catch((err) => {
							console.error('[Dashboard] 刷新最近测验失败:', err)
						})
				}
			} catch (err) {
				console.error('[Dashboard] 解析 WebSocket 消息失败:', err)
			}
		}

		window.addEventListener('websocket-message', handleMessage as EventListener)

		return () => {
			window.removeEventListener('websocket-message', handleMessage as EventListener)
		}
	}, [user?.id])

	if (loading) {
		return (
			<div className="flex items-center justify-center min-h-screen">
				<div className="text-center">
					<div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-gray-300 border-t-blue-500" />
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
		<div className="p-3 sm:p-4 md:p-6 max-w-7xl mx-auto">
			<WelcomeHeader
				username={user.username || '用户'}
				studyDays={stats?.studyDays || 0}
			/>
			<QuickContinue lastLearning={lastLearning} />
			{stats && (
				<StatsCards
					overallProgress={stats.progressPercentage}
					completedChapters={stats.completedChapters}
					totalChapters={stats.totalChapters}
					weeklyActivity={stats.weeklyActivity}
				/>
			)}
			<Suspense fallback={<div className="mb-6"><Spin size="large" /></div>}>
				<TopicProgress topics={topicProgress} />
			</Suspense>
			<Suspense fallback={<div className="mb-6"><Spin size="large" /></div>}>
				<RecentQuizzes quizzes={recentQuizzes} />
			</Suspense>
		</div>
	)
}
