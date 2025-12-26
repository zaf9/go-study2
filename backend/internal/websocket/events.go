// WebSocket 事件定义
// 定义了系统支持的 WebSocket 事件类型和数据结构

package websocket

/**
 * WebSocket 事件类型
 */
const (
	// 进度更新事件 - 用户完成章节学习时触发
	EventProgressUpdated = "progress_updated"

	// 测验完成事件 - 用户完成测验时触发
	EventQuizCompleted = "quiz_completed"
)

/**
 * 进度更新事件数据
 */
type ProgressUpdatedEventData struct {
	// 用户 ID
	UserID uint `json:"user_id"`

	// 主题 ID
	TopicID string `json:"topic_id"`

	// 章节 ID
	ChapterID string `json:"chapter_id"`

	// 是否完成
	Completed bool `json:"completed"`

	// 时间戳（ISO 8601 格式）
	Timestamp string `json:"timestamp"`
}

/**
 * 测验完成事件数据
 */
type QuizCompletedEventData struct {
	// 用户 ID
	UserID uint `json:"user_id"`

	// 测验 ID
	QuizID uint `json:"quiz_id"`

	// 得分
	Score int `json:"score"`

	// 总题数
	TotalQuestions int `json:"total_questions"`

	// 是否通过
	Passed bool `json:"passed"`

	// 时间戳（ISO 8601 格式）
	Timestamp string `json:"timestamp"`
}

/**
 * WebSocket 消息结构
 */
type WebSocketMessage struct {
	// 事件类型
	Event string `json:"event"`

	// 事件数据（使用 interface{} 支持不同类型的数据）
	Data interface{} `json:"data"`
}

/**
 * 创建进度更新消息
 */
func NewProgressUpdatedMessage(data ProgressUpdatedEventData) WebSocketMessage {
	return WebSocketMessage{
		Event: EventProgressUpdated,
		Data:  data,
	}
}

/**
 * 创建测验完成消息
 */
func NewQuizCompletedMessage(data QuizCompletedEventData) WebSocketMessage {
	return WebSocketMessage{
		Event: EventQuizCompleted,
		Data:  data,
	}
}

