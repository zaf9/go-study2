# Configuration Contract: HTTPS 配置

**Feature**: 005-https-protocol-support  
**Date**: 2025-12-08

## 配置文件

**路径**: `configs/config.yaml`（从根目录 `config.yaml` 迁移）

## 配置模式

### config.yaml 结构

```yaml
# HTTP 配置节
http:
  # HTTP 服务监听端口
  # 类型: int
  # 必填: 是
  port: 8080

# HTTPS 配置节
https:
  # 是否启用 HTTPS 模式
  # 类型: boolean
  # 默认值: false
  # 当设置为 true 时，服务将以 HTTPS 模式启动
  # 当设置为 false 或未配置时，服务以 HTTP 模式启动
  enabled: false
  
  # HTTPS 服务监听端口
  # 类型: int
  # 必填: 当 enabled = true 时
  port: 8443
  
  # TLS 证书文件路径
  # 类型: string
  # 必填: 当 enabled = true 时
  # 支持相对路径（基于工作目录）和绝对路径
  # 证书格式: PEM
  certFile: "./configs/certs/server.crt"
  
  # TLS 私钥文件路径
  # 类型: string
  # 必填: 当 enabled = true 时
  # 支持相对路径（基于工作目录）和绝对路径
  # 私钥格式: PEM (PKCS#1 或 PKCS#8)
  keyFile: "./configs/certs/server.key"

# 服务器通用配置
server:
  host: "127.0.0.1"
  shutdownTimeout: 10
```

---

## 行为契约

### HTTP 模式（默认）

**条件**: `https.enabled = false` 或未配置 `https` 节

**行为**:
- 服务监听 `server.host:http.port`
- 使用 HTTP 协议
- 现有功能完全兼容

### HTTPS 模式

**条件**: `https.enabled = true`

**前置条件**:
- `https.port` 已配置
- `https.certFile` 已配置且文件存在
- `https.keyFile` 已配置且文件存在

**行为**:
- 服务监听 `server.host:https.port`
- 使用 HTTPS 协议（TLS 1.2+）
- HTTP 协议不可用

**错误处理**:
- 配置缺失: 服务启动失败，显示缺失配置项的错误消息
- 文件不存在: 服务启动失败，显示文件路径的错误消息

---

## TLS 安全契约

| 属性 | 值 |
|------|-----|
| 最低 TLS 版本 | TLS 1.2 |
| 最高 TLS 版本 | TLS 1.3 |
| 证书格式 | PEM |
| 私钥格式 | PEM (PKCS#1/PKCS#8) |
