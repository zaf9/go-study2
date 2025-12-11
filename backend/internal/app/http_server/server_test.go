package http_server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-study2/internal/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestStaticFallbackAndApiPriority(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 构造临时静态目录，包含 index.html 与静态资源文件
		staticDir := t.TempDir()
		indexContent := "<html><body>index</body></html>"
		staticAsset := "console.log('static');"
		err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte(indexContent), 0o644)
		t.AssertNil(err)
		err = os.WriteFile(filepath.Join(staticDir, "app.js"), []byte(staticAsset), 0o644)
		t.AssertNil(err)

		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
			},
			Http: config.HttpConfig{
				Port: 0,
			},
			Static: config.StaticConfig{
				Enabled:     true,
				Path:        staticDir,
				SpaFallback: true,
			},
		}

		s, err := NewServer(cfg, "static-priority-test")
		t.AssertNil(err)
		t.AssertNE(s, nil)

		s.SetAccessLogEnabled(false)
		s.Start()
		defer s.Shutdown()

		base := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(base)

		// 静态资源应直接返回文件内容
		respStatic, err := client.Get(gctx.New(), "/app.js")
		t.AssertNil(err)
		defer respStatic.Close()
		bodyStatic := respStatic.ReadAllString()
		t.Assert(respStatic.StatusCode, 200)
		t.AssertIN("console.log('static')", bodyStatic)

		// 未命中的路径应回退到 index.html
		respFallback, err := client.Get(gctx.New(), "/unknown/path")
		t.AssertNil(err)
		defer respFallback.Close()
		bodyFallback := respFallback.ReadAllString()
		t.Assert(respFallback.StatusCode, 200)
		t.AssertIN("index", bodyFallback)

		// API 路由应保持优先级，不被静态回退覆盖
		respAPI, err := client.Get(gctx.New(), "/api/v1/topics?format=json")
		t.AssertNil(err)
		defer respAPI.Close()
		t.Assert(respAPI.StatusCode, 200)
	})
}

func TestNewServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
			},
			Http: config.HttpConfig{
				Port: 8080,
			},
		}

		// 测试不带名称创建服务器
		s, err := NewServer(cfg)
		t.AssertNil(err)
		t.AssertNE(s, nil)

		// 测试带名称创建服务器
		s2, err := NewServer(cfg, "test-server")
		t.AssertNil(err)
		t.AssertNE(s2, nil)
	})
}

func TestNewServerWithGracefulShutdown(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
			},
			Http: config.HttpConfig{
				Port: 0, // 使用随机端口
			},
		}

		s, err := NewServer(cfg)
		t.AssertNil(err)
		t.AssertNE(s, nil)
	})
}

func TestNewServerWithHTTPS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		certFile, keyFile, caPEM := mustCreateTempCertKey(t)
		port := mustGetFreePort(t)

		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
			},
			Http: config.HttpConfig{
				Port: 0,
			},
			Https: config.HttpsConfig{
				Enabled:  true,
				Port:     port,
				CertFile: certFile,
				KeyFile:  keyFile,
			},
		}

		s, err := NewServer(cfg, "https-unit-test")
		t.AssertNil(err)
		t.AssertNE(s, nil)

		s.SetAccessLogEnabled(false)
		s.Start()
		defer s.Shutdown()

		client := gtestClientWithCert(t, caPEM)
		url := fmt.Sprintf("https://127.0.0.1:%d/api/v1", s.GetListenedPort())
		resp, err := client.Get(gctx.New(), url+"/topics")
		t.AssertNil(err)
		defer resp.Close()
		t.Assert(resp.StatusCode, 200)
	})
}

func TestNewServerPortOccupiedHTTP(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		t.AssertNil(err)
		defer ln.Close()
		port := ln.Addr().(*net.TCPAddr).Port

		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "127.0.0.1",
			},
			Http: config.HttpConfig{
				Port: port,
			},
		}
		s, err := NewServer(cfg, "port-occupied-test")
		t.AssertNil(err)
		t.AssertNE(s, nil)

		// 再次占用同端口应返回“address already in use”
		_, err = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		t.AssertNE(err, nil)
		if err != nil {
			t.AssertIN("address", err.Error())
		}
	})
}

func mustCreateTempCertKey(t *gtest.T) (string, string, []byte) {
	t.Helper()
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")
	certPEM, keyPEM := generateSelfSignedCert(t)
	err := os.WriteFile(certFile, certPEM, 0o600)
	t.AssertNil(err)
	err = os.WriteFile(keyFile, keyPEM, 0o600)
	t.AssertNil(err)
	return certFile, keyFile, certPEM
}

func mustGetFreePort(t *gtest.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	t.AssertNil(err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func gtestClientWithCert(t *gtest.T, caPEM []byte) *gclient.Client {
	t.Helper()
	caPool := x509.NewCertPool()
	ok := caPool.AppendCertsFromPEM(caPEM)
	t.Assert(ok, true)
	client := g.Client()
	err := client.SetTLSConfig(&tls.Config{
		RootCAs:            caPool,
		InsecureSkipVerify: false,
	})
	t.AssertNil(err)
	return client
}

func generateSelfSignedCert(t *gtest.T) ([]byte, []byte) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	t.AssertNil(err)

	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	t.AssertNil(err)

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	t.AssertNil(err)

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	return certPEM, keyPEM
}
