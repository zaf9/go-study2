package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestLoad(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试正常加载配置
		cfg, err := Load()
		t.AssertNil(err)
		t.AssertNE(cfg, nil)
		t.Assert(cfg.Server.Host, "127.0.0.1")
		t.Assert(cfg.Http.Port, 8080)
		t.Assert(cfg.Logger.Level, "INFO")
		t.Assert(cfg.Https.Enabled, false)
	})
}

func TestValidate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试 Host 为空
		cfg := &Config{
			Server: ServerConfig{
				Host: "",
			},
			Http: HttpConfig{
				Port: 8080,
			},
		}
		err := Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("server.host", err.Error())

		// 测试 HTTP 端口缺失
		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
			},
			Http: HttpConfig{
				Port: 0,
			},
		}
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("http.port", err.Error())

		// 测试 HTTP 端口范围
		cfg.Http.Port = 70000
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("1-65535", err.Error())

		// 测试 HTTPS 启用但缺少端口
		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
			},
			Http: HttpConfig{
				Port: 8080,
			},
			Https: HttpsConfig{
				Enabled: true,
			},
		}
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("https.port", err.Error())

		// 测试 HTTPS 缺少证书路径
		cfg.Https.Port = 8443
		cfg.Https.CertFile = ""
		cfg.Https.KeyFile = ""
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("https.certFile", err.Error())

		// 测试 HTTPS 缺少私钥路径
		cfg.Https.CertFile = "./configs/certs/server.crt"
		cfg.Https.KeyFile = ""
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("https.keyFile", err.Error())

		// 测试 HTTPS 证书不存在
		cfg.Https.KeyFile = "./configs/certs/server.key"
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("证书文件不存在", err.Error())

		// 构造临时证书文件用于通过验证
		tempDir := t.TempDir()
		certFile := filepath.Join(tempDir, "server.crt")
		keyFile := filepath.Join(tempDir, "server.key")
		certPEM, keyPEM := generateSelfSignedCert(t)
		t.AssertNil(os.WriteFile(certFile, certPEM, 0o600))
		t.AssertNil(os.WriteFile(keyFile, keyPEM, 0o600))

		cfg = &Config{
			Server: ServerConfig{
				Host: "127.0.0.1",
			},
			Http: HttpConfig{
				Port: 8080,
			},
			Https: HttpsConfig{
				Enabled:  true,
				Port:     8443,
				CertFile: certFile,
				KeyFile:  keyFile,
			},
		}
		err = Validate(cfg)
		t.AssertNil(err)

		// 测试证书与私钥不匹配
		otherCertFile := filepath.Join(tempDir, "other.crt")
		otherKeyFile := filepath.Join(tempDir, "other.key")
		certPEM2, _ := generateSelfSignedCert(t)
		t.AssertNil(os.WriteFile(otherCertFile, certPEM2, 0o600))
		// 写入不匹配的私钥
		t.AssertNil(os.WriteFile(otherKeyFile, keyPEM, 0o600))
		cfg.Https.CertFile = otherCertFile
		cfg.Https.KeyFile = otherKeyFile
		err = Validate(cfg)
		t.AssertNE(err, nil)
		t.AssertIN("不匹配", err.Error())

		// 测试相对路径解析
		relDir := t.TempDir()
		relCert := "rel-server.crt"
		relKey := "rel-server.key"
		t.AssertNil(os.WriteFile(filepath.Join(relDir, relCert), certPEM, 0o600))
		t.AssertNil(os.WriteFile(filepath.Join(relDir, relKey), keyPEM, 0o600))
		origWd, _ := os.Getwd()
		_ = os.Chdir(relDir)
		defer os.Chdir(origWd)
		cfg = &Config{
			Server: ServerConfig{Host: "127.0.0.1"},
			Http:   HttpConfig{Port: 8080},
			Https: HttpsConfig{
				Enabled:  true,
				Port:     8443,
				CertFile: relCert,
				KeyFile:  relKey,
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
		t.Assert(cfg.Http.Port > 0, true)
		t.AssertNE(cfg.Logger.Level, "")
	})
}

// generateSelfSignedCert 生成自签名证书
func generateSelfSignedCert(t *gtest.T) ([]byte, []byte) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	t.AssertNil(err)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		DNSNames:    []string{"localhost"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	t.AssertNil(err)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return certPEM, keyPEM
}

func keyPEMFile(t *gtest.T, dir string, pemBytes []byte) string {
	t.Helper()
	path := filepath.Join(dir, "mismatch.key")
	t.AssertNil(os.WriteFile(path, pemBytes, 0o600))
	return path
}
