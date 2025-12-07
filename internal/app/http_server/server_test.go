package http_server

import (
	"testing"

	"go-study2/internal/config"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestNewServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
				Port: 8080,
			},
		}

		// 测试不带名称创建服务器
		s := NewServer(cfg)
		t.AssertNE(s, nil)

		// 测试带名称创建服务器
		s2 := NewServer(cfg, "test-server")
		t.AssertNE(s2, nil)
	})
}

func TestNewServerWithGracefulShutdown(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
				Port: 0, // 使用随机端口
			},
		}

		s := NewServer(cfg)
		t.AssertNE(s, nil)
	})
}

