package handler

// Topic 学习主题
type Topic struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Response 标准响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TopicListResponse 主题列表响应数据
type TopicListResponse struct {
	Topics []Topic `json:"topics"`
}

// LexicalMenuItem 词法元素菜单项
type LexicalMenuItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"` // 英文标识符，用于路由
}

// LexicalMenuResponse 词法元素菜单响应
type LexicalMenuResponse struct {
	Items []LexicalMenuItem `json:"items"`
}

// Handler 处理HTTP请求的控制器
type Handler struct{}

func New() *Handler {
	return &Handler{}
}
