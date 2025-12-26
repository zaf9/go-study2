/**
 * Dashboard 加载中页面
 * 数据加载时的占位页面
 */

export default function Loading() {
	return (
		<div className="flex min-h-screen items-center justify-center bg-gray-50">
			<div className="flex flex-col items-center">
					<div className="h-12 w-12 animate-spin rounded-full border-4 border-solid border-gray-300 border-t-blue-500" />
					<p className="mt-4 text-gray-500 text-sm">加载 Dashboard 数据中...</p>
			</div>
		</div>
	)
}

