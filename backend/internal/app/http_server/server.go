package http_server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go-study2/internal/app/http_server/middleware"
	"go-study2/internal/config"
	"go-study2/internal/domain/user"
	"go-study2/internal/infrastructure/database"
	"go-study2/internal/infrastructure/repository"
	appjwt "go-study2/internal/pkg/jwt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

// NewServer 创建并配置 HTTP/HTTPS 服务器
func NewServer(cfg *config.Config, names ...string) (*ghttp.Server, error) {
	var s *ghttp.Server
	if len(names) > 0 {
		s = g.Server(names[0])
	} else {
		s = g.Server()
	}

	// 基础配置
	s.SetGraceful(true) // 开启优雅关闭

	// 注册全局中间件
	s.Use(middleware.Logger)
	s.Use(middleware.Cors)
	s.Use(middleware.PanicRecovery) // Panic recovery should be before access logging
	s.Use(middleware.AccessLog)

	// 注册路由
	RegisterRoutes(s)

	// 初始化 WebSocket Hub
	InitWebSocketHub()

	// 启动前确保默认管理员存在（幂等）
	if err := ensureDefaultAdmin(cfg); err != nil {
		return nil, err
	}

	if cfg.Https.Enabled {
		tlsCfg, err := buildTLSConfig(cfg.Https)
		if err != nil {
			return nil, err
		}
		// 禁用 HTTP，确保仅监听 HTTPS
		s.SetAddr("")
		_ = s.SetConfigWithMap(map[string]any{
			"Address": "",
		})
		s.SetTLSConfig(tlsCfg)
		s.EnableHTTPS(cfg.Https.CertFile, cfg.Https.KeyFile, tlsCfg)
		s.SetHTTPSAddr(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Https.Port))
		g.Log().Infof(gctx.New(), "HTTPS 模式已启用，监听地址: %s:%d，HTTP 已禁用", cfg.Server.Host, cfg.Https.Port)
	} else {
		s.SetAddr(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Http.Port))
		g.Log().Infof(gctx.New(), "HTTP 模式已启用，监听地址: %s:%d", cfg.Server.Host, cfg.Http.Port)
	}

	if cfg.Static.Enabled && cfg.Static.Path != "" {
		// 静态资源目录清理，确保路径可复用
		staticPath := filepath.Clean(cfg.Static.Path)
		s.SetServerRoot(staticPath)
		s.AddStaticPath("/", staticPath)

		// SPA 回退：非 /api/ 路径优先尝试静态文件，否则返回 index.html
		if cfg.Static.SpaFallback {
			s.BindHandler("/*", func(r *ghttp.Request) {
				if strings.HasPrefix(r.URL.Path, "/api/") {
					r.Middleware.Next()
					return
				}
				requested := strings.TrimPrefix(r.URL.Path, "/")
				candidate := filepath.Join(staticPath, requested)

				if requested != "" && gfile.Exists(candidate) && !gfile.IsDir(candidate) {
					r.Response.ServeFile(candidate)
					return
				}
				r.Response.ServeFile(filepath.Join(staticPath, "index.html"))
			})
		}
	}

	return s, nil
}

// buildTLSConfig 构建 TLS 配置，设定最小版本和可选的 CA/跳过校验
func buildTLSConfig(cfg config.HttpsConfig) (*tls.Config, error) {
	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}

	if cfg.CaFile != "" {
		caData, err := os.ReadFile(cfg.CaFile)
		if err != nil {
			return nil, fmt.Errorf("读取 CA 证书失败: %w", err)
		}
		caPool := x509.NewCertPool()
		if ok := caPool.AppendCertsFromPEM(caData); !ok {
			return nil, fmt.Errorf("CA 证书格式无效: %s", cfg.CaFile)
		}
		tlsCfg.ClientCAs = caPool
		tlsCfg.RootCAs = caPool
	}

	tlsCfg.InsecureSkipVerify = cfg.InsecureSkipVerify
	if cfg.InsecureSkipVerify {
		g.Log().Warning(gctx.New(), "已启用 InsecureSkipVerify（仅用于测试/开发环境，生产环境请关闭以确保证书校验）")
	}

	return tlsCfg, nil
}

func ensureDefaultAdmin(cfg *config.Config) error {
	if cfg.Database.Path == "" {
		return nil
	}
	db := database.Default()
	if db == nil {
		var err error
		db, err = database.Init(gctx.New(), cfg.Database)
		if err != nil {
			return err
		}
	}
	svc := user.NewService(repository.NewUserRepository(db), appjwt.AccessTokenTTL(), appjwt.RefreshTokenTTL())
	return svc.EnsureDefaultAdmin(gctx.New())
}
