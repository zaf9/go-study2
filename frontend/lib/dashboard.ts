/**
 * Dashboard API 调用
 *
 * 提供与 Dashboard 相关的 API 调用函数
 */

import api from './api'
import { DashboardStats } from '@/types/dashboard'

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

