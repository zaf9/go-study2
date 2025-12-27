/**
 * Dashboard API 调用
 *
 * 提供与 Dashboard 相关的 API 调用函数
 */

import api from './api'
import { DashboardStats, LastLearningRecord, TopicProgressSummary, RecentQuizSummary } from '@/types/dashboard'

/**
 * 获取 Dashboard 统计数据
 * 从现有的 /api/v1/progress 端点获取并转换为 DashboardStats 格式
 */
export async function fetchDashboardStats(): Promise<DashboardStats> {
	try {
		const response = await api.get<{
			code: number
			message: string
			data: {
				overall: {
					progress: number
					completedChapters: number
					totalChapters: number
					studyDays: number
				}
			}
		}>('/api/v1/progress')

		if (response.data.code !== 0 || !response.data.data) {
			throw new Error(response.data.message || '获取进度数据失败')
		}

		const { overall } = response.data.data

		return {
			studyDays: overall.studyDays || 0,
			totalChapters: overall.totalChapters || 0,
			completedChapters: overall.completedChapters || 0,
			progressPercentage: overall.progress || 0,
			weeklyActivity: 0, // 暂时设置为 0,后续可以添加专门的 API
		}
	} catch (error) {
		console.error('获取 Dashboard 统计数据失败:', error)
		throw error
	}
}

/**
 * 获取最后学习记录
 * 从 /api/v1/progress/last 端点获取用户最后一次学习的主题和章节信息
 */
export async function fetchLastLearning(): Promise<LastLearningRecord | null> {
	try {
		const response = await api.get<{
			code: number
			message: string
			data: LastLearningRecord | null
		}>('/api/v1/progress/last')

		if (response.data.code !== 0) {
			throw new Error(response.data.message || '获取最后学习记录失败')
		}

		return response.data.data || null
	} catch (error) {
		console.error('获取最后学习记录失败:', error)
		throw error
	}
}

/**
 * 获取主题进度汇总列表
 * 从 /api/v1/progress 端点获取并转换为 TopicProgressSummary 格式
 */
export async function fetchTopicProgress(): Promise<TopicProgressSummary[]> {
	try {
		const response = await api.get<{
			code: number
			message: string
			data: {
				topics: Array<{
					id: string
					name: string
					completedChapters: number
					totalChapters: number
					progress: number
				}>
			}
		}>('/api/v1/progress')

		if (response.data.code !== 0 || !response.data.data) {
			throw new Error(response.data.message || '获取进度数据失败')
		}

		const { topics } = response.data.data

		// 转换为 TopicProgressSummary 格式
		// 注意：后端的 progress 是 0-100 的整数，需要转换为浮点数并保留一位小数
		return (topics || []).map((topic) => {
			// 计算百分比（保留一位小数）
			let percentage = 0.0
			if (topic.totalChapters > 0) {
				percentage = (topic.completedChapters / topic.totalChapters) * 100
				percentage = Math.round(percentage * 10) / 10 // 四舍五入到一位小数
			}

			return {
				topicId: topic.id,
				topicName: topic.name,
				displayName: topic.name,
				completedChapters: topic.completedChapters || 0,
				totalChapters: topic.totalChapters || 0,
				percentage,
			}
		})
	} catch (error) {
		console.error('获取主题进度汇总失败:', error)
		throw error
	}
}

/**
 * 获取最近测验记录
 * 从 /api/v1/quiz/recent 端点获取用户最近的测验记录
 */
export async function fetchRecentQuizzes(limit: number = 5): Promise<RecentQuizSummary[]> {
	try {
		const response = await api.get<{
			code: number
			message: string
			data: Array<{
				id: number
				topic_name: string
				chapter_name: string
				score: number
				total_questions: number
				passed: boolean
				completed_at: string
			}>
		}>('/api/v1/quiz/recent', { params: { limit } })

		if (response.data.code !== 0) {
			throw new Error(response.data.message || '获取最近测验记录失败')
		}

		// 转换为前端类型格式
		return (response.data.data || []).map((item) => ({
			id: item.id,
			topicName: item.topic_name,
			chapterName: item.chapter_name,
			score: item.score,
			totalQuestions: item.total_questions,
			passed: item.passed,
			completedAt: item.completed_at,
		}))
	} catch (error) {
		console.error('获取最近测验记录失败:', error)
		throw error
	}
}

