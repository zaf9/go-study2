# Quickstart: HTTPS 协议支持

**Feature**: 005-https-protocol-support

## 快速启用 HTTPS

### 1. 生成自签名证书

```bash
# 创建证书目录
mkdir -p configs/certs

# 生成私钥和自签名证书（有效期 365 天）
openssl req -x509 -newkey rsa:4096 -keyout configs/certs/server.key -out configs/certs/server.crt -days 365 -nodes -subj "/CN=localhost"
```

### 2. 修改配置文件

编辑 `configs/config.yaml`，添加 HTTPS 配置：

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
```

### 3. 启动服务

```bash
go run main.go
```

### 4. 验证 HTTPS

```bash
# 使用 curl 测试（-k 跳过证书验证，仅用于自签名证书）
curl -k https://localhost:8443/api/topics
```

---

## 切换回 HTTP 模式

修改 `configs/config.yaml`：

```yaml
http:
  port: 8080

https:
  enabled: false
```

或直接删除 `https` 配置节中的 `enabled: true`。

---

## 常见问题

### Q: 启动时提示证书文件不存在？

确认 `certFile` 和 `keyFile` 路径正确。支持相对路径（基于工作目录）和绝对路径。默认证书存放在 `configs/certs/` 目录。

### Q: 如何使用生产证书？

将 CA 签发的证书和私钥替换 `configs/certs/` 目录下的文件，或修改配置指向实际路径。

### Q: TLS 版本要求？

服务最低支持 TLS 1.2，同时支持 TLS 1.3。
