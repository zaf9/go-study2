/**
 * Dashboard API 调用
 *
 * 提供与 Dashboard 相关的 API 调用函数
 */

import api from './api'
import { DashboardStats, LastLearningRecord, TopicProgressSummary } from '@/types/dashboard'

/**
 * 获取 Dashboard 统计数据
 * 从现有的 /api/v1/progress 端点获取并转换为 DashboardStats 格式
 */
export async function fetchDashboardStats(): Promise<DashboardStats> {
	try {
		const response = await api.get<any>('/api/v1/progress')

		if (response.code !== 0 || !response.data) {
			throw new Error(response.message || '获取进度数据失败')
		}

		const { overall } = response.data as { overall: { Progress: number; CompletedChapters: number; TotalChapters: number; StudyDays: number } }

		return {
			studyDays: overall.StudyDays || 0,
			totalChapters: overall.TotalChapters || 0,
			completedChapters: overall.CompletedChapters || 0,
			progressPercentage: overall.Progress || 0,
			weeklyActivity: 0, // 暂时设置为 0，后续可以添加专门的 API
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

		if (response.code !== 0) {
			throw new Error(response.message || '获取最后学习记录失败')
		}

		return response.data || null
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
		const response = await api.get<any>('/api/v1/progress')

		if (response.code !== 0 || !response.data) {
			throw new Error(response.message || '获取进度数据失败')
		}

		const { topics } = response.data as {
			topics: Array<{
				id: string
				name: string
				completedChapters: number
				totalChapters: number
				progress: number
			}>
		}

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

