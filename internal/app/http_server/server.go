package http_server

import (
	"go-study2/internal/app/http_server/middleware"
	"go-study2/internal/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// NewServer 创建并配置 HTTP 服务器
func NewServer(cfg *config.Config, names ...string) *ghttp.Server {
	var s *ghttp.Server
	if len(names) > 0 {
		s = g.Server(names[0])
	} else {
		s = g.Server()
	}

	// 基础配置
	s.SetPort(cfg.Server.Port)
	s.SetGraceful(true) // 开启优雅关闭

	// 注册全局中间件
	s.Use(middleware.Logger)

	// 注册路由
	RegisterRoutes(s)

	return s
}
