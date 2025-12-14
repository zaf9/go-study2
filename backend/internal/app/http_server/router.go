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

		// 认证路由（无需 JWT 验证）
		group.POST("/auth/login", h.Login)
		group.POST("/auth/refresh", h.RefreshToken)

		// 需要认证的路由
		group.Group("/", func(authGroup *ghttp.RouterGroup) {
			authGroup.Middleware(middleware.Auth)
			authGroup.Middleware(middleware.ForceChangePassword)
			authGroup.POST("/auth/register", h.Register)
			authGroup.GET("/auth/profile", h.GetProfile)
			authGroup.POST("/auth/logout", h.Logout)
			authGroup.POST("/auth/change-password", h.ChangePassword)

			// 学习进度
			authGroup.GET("/progress", h.GetAllProgress)
			authGroup.GET("/progress/:topic", h.GetTopicProgress)
			authGroup.POST("/progress", h.SaveProgress)

			// 测验
			authGroup.GET("/quiz/:topic/:chapter", h.GetQuiz)
			authGroup.GET("/quiz/:topic/:chapter/stats", h.GetQuizStats)
			authGroup.POST("/quiz/submit", h.SubmitQuiz)
			authGroup.GET("/quiz/history", h.GetQuizHistory)
			authGroup.GET("/quiz/history/:topic", h.GetQuizHistory)
		})

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

		// Variables 菜单
		group.ALL("/topic/variables", h.GetVariablesMenu)
		// Variables 内容
		group.ALL("/topic/variables/:subtopic", h.GetVariableContent)

		// Types 菜单
		group.ALL("/topic/types", h.GetTypesMenu)
		// Types 内容
		group.ALL("/topic/types/:subtopic", h.GetTypesContent)
		// Types 提纲
		group.ALL("/topic/types/outline", h.GetTypesOutline)
		// Types 测验提交
		group.ALL("/topic/types/quiz/submit", h.SubmitTypesQuiz)
		// Types 搜索
		group.ALL("/topic/types/search", h.SearchTypes)

		// 后续路由将在其他 User Story 中添加
	})
}
