'use client'

/**
 * TruncatedText 组件
 * 对过长文本进行截断，并在悬停时通过Tooltip显示完整文本
 * 实现FR-023：文本截断功能
 */

import { Tooltip } from 'antd'

interface TruncatedTextProps {
	/**
	 * 要显示的文本内容
	 */
	text: string

	/**
	 * 最大显示字符数，超过此长度将截断
	 * 默认值为20
	 */
	maxLength?: number

	/**
	 * 自定义样式类名
	 */
	className?: string

	/**
	 * 是否始终显示Tooltip（即使文本未截断）
	 * 默认值为false，只在文本被截断时显示
	 */
	alwaysShowTooltip?: boolean
}

/**
 * TruncatedText 组件
 * 实现长文本截断和Tooltip显示功能
 * 鼠标悬停300ms后显示完整文本（通过Ant Design Tooltip的默认延迟实现）
 */
export const TruncatedText: React.FC<TruncatedTextProps> = ({
	text,
	maxLength = 20,
	className = '',
	alwaysShowTooltip = false,
}) => {
	// 处理空值或 undefined
	const safeText = text || ''
	
	// 判断文本是否需要截断
	const shouldTruncate = safeText.length > maxLength
	const displayText = shouldTruncate ? `${safeText.slice(0, maxLength)}...` : safeText

	// 如果需要截断或设置了alwaysShowTooltip，则显示Tooltip
	const showTooltip = shouldTruncate || alwaysShowTooltip

	if (showTooltip) {
		return (
			<Tooltip title={safeText} mouseEnterDelay={0.3} placement="top">
				<span className={`truncate ${className}`} style={{ display: 'inline-block', maxWidth: '100%' }}>
					{displayText}
				</span>
			</Tooltip>
		)
	}

	return <span className={className}>{displayText}</span>
}

