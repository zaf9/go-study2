# Data Model: HTTPS 协议支持

**Feature**: 005-https-protocol-support  
**Date**: 2025-12-08

## 实体定义

### HttpConfig

HTTP 配置结构，包含 HTTP 服务端口配置。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `port` | `int` | 是 | - | HTTP 服务监听端口 |

---

### HttpsConfig

HTTPS 配置结构，包含启用状态、端口和证书路径配置。

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `enabled` | `bool` | 否 | `false` | 是否启用 HTTPS 模式 |
| `port` | `int` | 条件 | - | HTTPS 服务监听端口，当 enabled=true 时必填 |
| `certFile` | `string` | 条件 | - | 证书文件路径（PEM 格式），当 enabled=true 时必填 |
| `keyFile` | `string` | 条件 | - | 私钥文件路径（PEM 格式），当 enabled=true 时必填 |

**验证规则**:
1. 当 `enabled = true` 时，`port`、`certFile` 和 `keyFile` 均为必填
2. 证书文件必须存在且可读
3. 私钥文件必须存在且可读
4. 路径支持相对路径（基于工作目录）和绝对路径
5. 证书文件默认存放在 `configs/certs/` 目录

---

### Config（扩展）

现有配置结构的扩展版本。

| 字段 | 类型 | 说明 |
|------|------|------|
| `http` | `HttpConfig` | **新增** HTTP 配置（替代原 server.port） |
| `https` | `HttpsConfig` | **新增** HTTPS 配置 |
| `server` | `ServerConfig` | 服务器通用配置（host, shutdownTimeout） |
| `logger` | `LoggerConfig` | 日志配置（现有） |

---

## 状态转换

### 协议模式选择

```
配置加载 (configs/config.yaml)
    │
    ├─ https.enabled = false (或未配置)
    │       │
    │       └──> HTTP 模式
    │               └──> 监听 http.port (HTTP)
    │
    └─ https.enabled = true
            │
            ├─ 证书验证失败
            │       └──> 启动失败 + 错误提示
            │
            └─ 证书验证成功
                    └──> HTTPS 模式
                            └──> 监听 https.port (HTTPS, TLS 1.2+)
```

---

## 配置文件示例

**配置文件路径**: `configs/config.yaml`

### HTTP 模式（默认）

```yaml
http:
  port: 8080

server:
  host: "127.0.0.1"
  shutdownTimeout: 10

logger:
  level: "INFO"
  path: "./logs"
  stdout: true

# https 配置省略或设置 enabled: false
```

### HTTPS 模式

```yaml
http:
  port: 8080

https:
  enabled: true
  port: 8443
  certFile: "./configs/certs/server.crt"
  keyFile: "./configs/certs/server.key"

server:
  host: "127.0.0.1"
  shutdownTimeout: 10

logger:
  level: "INFO"
  path: "./logs"
  stdout: true
```

---

## 错误消息

| 场景 | 错误消息 |
|------|---------|
| 缺少 certFile | `配置项 https.certFile 为必填项（当 https.enabled = true 时）` |
| 缺少 keyFile | `配置项 https.keyFile 为必填项（当 https.enabled = true 时）` |
| 证书文件不存在 | `证书文件不存在: {路径}` |
| 私钥文件不存在 | `私钥文件不存在: {路径}` |
