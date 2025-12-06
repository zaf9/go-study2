package http_server

import (
	"go-study2/internal/app/http_server/handler"
	"go-study2/internal/app/http_server/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// RegisterRoutes 注册路由
func RegisterRoutes(s *ghttp.Server) {
	h := handler.New()

	// API v1 路由组
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		// 应用格式转换中间件
		group.Middleware(middleware.Format)

		// 主题列表
		group.ALL("/topics", h.GetTopics)

		// 词法元素菜单
		group.ALL("/topic/lexical_elements", h.GetLexicalMenu)
		// 词法元素章节内容
		group.ALL("/topic/lexical_elements/:chapter", h.GetLexicalContent)

		// Constants 菜单
		group.ALL("/topic/constants", h.GetConstantsMenu)
		// Constants 内容
		group.ALL("/topic/constants/:subtopic", h.GetConstantsContent)

		// 后续路由将在其他 User Story 中添加
	})
}
