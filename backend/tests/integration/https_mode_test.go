package integration

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

	"go-study2/internal/app/http_server"
	"go-study2/internal/config"

	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
)

// T033: 自签名证书握手集成测试；T034: 启用 HTTPS 时禁用 HTTP 端口
func TestHttpsModeWithSelfSignedCert(t *testing.T) {
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

	s, err := http_server.NewServer(cfg, "https-mode-test")
	if err != nil {
		t.Fatalf("failed to create https server: %v", err)
	}
	s.SetAccessLogEnabled(false)
	s.Start()
	defer s.Shutdown()

	client := newHttpsClient(t, caPEM)
	baseURL := fmt.Sprintf("https://127.0.0.1:%d/api/v1", cfg.Https.Port)
	resp, err := client.Post(gctx.New(), baseURL+"/topics")
	if err != nil {
		t.Fatalf("https request failed: %v", err)
	}
	defer resp.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// HTTP 明文请求应失败（未监听）
	httpClient := gclient.New()
	httpResp, httpErr := httpClient.Post(gctx.New(), fmt.Sprintf("http://127.0.0.1:%d/api/topics", port))
	if httpResp != nil {
		defer httpResp.Close()
	}
	if httpErr == nil && httpResp != nil && httpResp.StatusCode < 400 {
		t.Fatalf("expected http request to fail when https enabled, got %d", httpResp.StatusCode)
	}
}

// T044: 协议切换耗时测量，要求 ≤ 30 秒
func TestProtocolSwitchDuration(t *testing.T) {
	start := time.Now()
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

	s, err := http_server.NewServer(cfg, "https-switch-measure")
	if err != nil {
		t.Fatalf("failed to create https server: %v", err)
	}
	s.SetAccessLogEnabled(false)
	s.Start()
	defer s.Shutdown()

	client := newHttpsClient(t, caPEM)
	baseURL := fmt.Sprintf("https://127.0.0.1:%d/api/v1", port)
	resp, err := client.Post(gctx.New(), baseURL+"/topics")
	if err != nil {
		t.Fatalf("https request failed: %v", err)
	}
	defer resp.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	elapsed := time.Since(start)
	if elapsed > 30*time.Second {
		t.Fatalf("protocol switch took too long: %s (limit 30s)", elapsed)
	}
}

func mustCreateTempCertKey(t *testing.T) (string, string, []byte) {
	t.Helper()
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "server.crt")
	keyFile := filepath.Join(tempDir, "server.key")
	certPEM, keyPEM := generateSelfSignedCert(t)
	if err := os.WriteFile(certFile, certPEM, 0o600); err != nil {
		t.Fatalf("write cert failed: %v", err)
	}
	if err := os.WriteFile(keyFile, keyPEM, 0o600); err != nil {
		t.Fatalf("write key failed: %v", err)
	}
	return certFile, keyFile, certPEM
}

func mustGetFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen temp port failed: %v", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func newHttpsClient(t *testing.T, caPEM []byte) *gclient.Client {
	t.Helper()
	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(caPEM); !ok {
		t.Fatalf("append ca pem failed")
	}
	client := gclient.New()
	if err := client.SetTLSConfig(&tls.Config{
		RootCAs:            caPool,
		InsecureSkipVerify: false,
	}); err != nil {
		t.Fatalf("set tls config failed: %v", err)
	}
	return client
}

func generateSelfSignedCert(t *testing.T) ([]byte, []byte) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key failed: %v", err)
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatalf("serial gen failed: %v", err)
	}

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
	if err != nil {
		t.Fatalf("create cert failed: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	return certPEM, keyPEM
}
