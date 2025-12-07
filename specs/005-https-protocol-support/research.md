# Research: HTTPS 协议支持

**Feature**: 005-https-protocol-support  
**Date**: 2025-12-08  
**Status**: Complete

## 研究问题

1. GoFrame 如何启用 HTTPS 模式？
2. 如何配置 TLS 证书和密钥文件路径？
3. 如何设置最低 TLS 版本（TLS 1.2+）？
4. GoFrame 是否支持同时运行 HTTP 和 HTTPS？
5. TLS 配置的最佳实践

---

## 决策 1：HTTPS 启用方式

**Decision**: 使用 GoFrame 的 `EnableHTTPS()` 方法结合 `SetTLSConfig()` 配置自定义 TLS 参数

**Rationale**: 
- GoFrame 原生支持 HTTPS，无需额外依赖
- `EnableHTTPS()` 方法简单直接，一行代码即可启用
- 支持传入自定义 `*tls.Config` 以配置 TLS 版本

**Alternatives Considered**:
| 方案 | 优点 | 缺点 | 评估 |
|------|------|------|------|
| EnableHTTPS 方法 | 简单直接 | 需额外设置 TLS 版本 | ✅ 采用 |
| 反向代理（Nginx） | 统一管理 | 增加部署复杂度 | ❌ 超出范围 |
| 仅配置文件 | 运维友好 | 无法配置 TLS 版本 | ❌ 不够灵活 |

---

## 决策 2：协议模式选择

**Decision**: 单协议模式 - 根据 `https.enabled` 配置项选择 HTTP 或 HTTPS，不同时运行

**Rationale**:
- 符合规格要求：`https.enabled = true` 时禁用 HTTP
- 简化配置和部署，避免端口冲突
- 符合 YAGNI 原则，满足当前需求

**Implementation**:
```go
if cfg.Https.Enabled {
    // 仅启用 HTTPS
    s.EnableHTTPS(cfg.Https.CertFile, cfg.Https.KeyFile, tlsConfig)
    s.SetHTTPSPort(cfg.Server.Port)
} else {
    // 仅启用 HTTP
    s.SetPort(cfg.Server.Port)
}
```

---

## 决策 3：TLS 版本配置

**Decision**: 最低 TLS 1.2，最高 TLS 1.3

**Rationale**:
- TLS 1.2 是当前行业标准最低版本
- TLS 1.3 提供最佳安全性和性能
- 符合澄清会话中确定的 "TLS 1.2+" 要求

**Implementation**:
```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    MaxVersion: tls.VersionTLS13,
}
```

---

## 决策 4：配置结构设计

**Decision**: 在 Config 结构中添加 HttpsConfig 嵌套结构

**Rationale**:
- 与现有 ServerConfig、LoggerConfig 结构保持一致
- 配置项清晰分组，便于管理
- HTTP 和 HTTPS 端口独立配置

**Implementation**:
```go
// HttpConfig HTTP 配置
type HttpConfig struct {
    Port int `json:"port"` // HTTP 服务监听端口
}

// HttpsConfig HTTPS 配置
type HttpsConfig struct {
    Enabled  bool   `json:"enabled"`   // 是否启用 HTTPS
    Port     int    `json:"port"`      // HTTPS 服务监听端口
    CertFile string `json:"certFile"`  // 证书文件路径
    KeyFile  string `json:"keyFile"`   // 私钥文件路径
}

// Config 应用配置结构
type Config struct {
    Http   HttpConfig   `json:"http"`   // 新增
    Https  HttpsConfig  `json:"https"`  // 新增
    Server ServerConfig `json:"server"`
    Logger LoggerConfig `json:"logger"`
}
```

**配置文件路径**: `configs/config.yaml`（从根目录迁移）

---

## 决策 5：证书验证时机

**Decision**: 在配置加载阶段进行证书文件存在性验证

**Rationale**:
- 尽早发现配置错误，提供友好提示
- 避免服务启动后才报错

**Implementation**:
```go
func Validate(cfg *Config) error {
    if cfg.Https.Enabled {
        if cfg.Https.Port == 0 {
            return fmt.Errorf("配置项 https.port 为必填项（当 https.enabled = true 时）")
        }
        if cfg.Https.CertFile == "" {
            return fmt.Errorf("配置项 https.certFile 为必填项（当 https.enabled = true 时）")
        }
        if cfg.Https.KeyFile == "" {
            return fmt.Errorf("配置项 https.keyFile 为必填项（当 https.enabled = true 时）")
        }
        // 验证文件存在性
        if _, err := os.Stat(cfg.Https.CertFile); os.IsNotExist(err) {
            return fmt.Errorf("证书文件不存在: %s", cfg.Https.CertFile)
        }
        if _, err := os.Stat(cfg.Https.KeyFile); os.IsNotExist(err) {
            return fmt.Errorf("私钥文件不存在: %s", cfg.Https.KeyFile)
        }
    }
    return nil
}
```

---

## GoFrame API 参考

| 方法 | 说明 |
|------|------|
| `EnableHTTPS(certFile, keyFile string, tlsConfig ...*tls.Config)` | 启用 HTTPS |
| `SetHTTPSPort(port ...int)` | 设置 HTTPS 端口 |
| `SetTLSConfig(tlsConfig *tls.Config)` | 设置自定义 TLS 配置 |
| `SetPort(port ...int)` | 设置 HTTP 端口 |

---

## 结论

所有技术问题已解决，可进入 Phase 1 设计阶段。
