package config

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestLoad(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试正常加载配置
		cfg, err := Load()
		t.AssertNil(err)
		t.AssertNE(cfg, nil)
		t.Assert(cfg.Server.Host, "127.0.0.1")
		t.Assert(cfg.Server.Port, 8080)
		t.Assert(cfg.Logger.Level, "INFO")
	})
}

func TestValidate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试 Host 为空
		cfg := &Config{
			Server: ServerConfig{
				Host: "",
				Port: 8080,
			},
		}
		err := Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("server.host", err.Error())

		// 测试 Port 为 0
		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
				Port: 0,
			},
		}
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("server.port", err.Error())

		// 测试 Port 小于 1
		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
				Port: 0,
			},
		}
		err = Validate(cfg)
		t.AssertNE(err, nil)

		// 测试 Port 大于 65535
		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
				Port: 65536,
			},
		}
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("1-65535", err.Error())

		// 测试有效配置
		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
				Port: 8080,
			},
		}
		err = Validate(cfg)
		t.AssertNil(err)

		// 测试边界值：Port = 1
		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
				Port: 1,
			},
		}
		err = Validate(cfg)
		t.AssertNil(err)

		// 测试边界值：Port = 65535
		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
				Port: 65535,
			},
		}
		err = Validate(cfg)
		t.AssertNil(err)
	})
}

func TestLoadWithValidConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试加载有效配置
		cfg, err := Load()
		t.AssertNil(err)
		t.AssertNE(cfg, nil)
		
		// 验证配置结构完整性
		t.AssertNE(cfg.Server.Host, "")
		t.Assert(cfg.Server.Port > 0, true)
		t.AssertNE(cfg.Logger.Level, "")
	})
}

