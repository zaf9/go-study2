# HTTPS 域名访问配置指南

## 方案说明

本项目使用域名方式访问 HTTPS 服务，避免 IP 变化时需要重新生成证书。

**域名**: `go.study.org`

**优势**:
- 一次生成证书，永久有效（365天）
- IP 变化时只需修改客户端 hosts 文件
- 无需重新生成证书或重启服务

## 快速配置

### 步骤 1：生成证书（服务器端）

在项目根目录执行：

```bash
# 使用默认域名 go.study.org 生成证书
./scripts/generate-cert.sh

# 或强制覆盖现有证书
./scripts/generate-cert.sh -f
```

### 步骤 2：重启服务（服务器端）

```bash
# 如果使用 systemd
sudo systemctl restart go-study2

# 或直接重启
cd /opt/hadoop/go-study2
./scripts/stop.sh
./scripts/start.sh
```

### 步骤 3：配置客户端 hosts

在**客户端**机器的 hosts 文件中添加：

**Linux/Mac**: `/etc/hosts`
**Windows**: `C:\Windows\System32\drivers\etc\hosts`

```
# Go-Study2 访问域名
172.28.186.65  go.study.org
```

### 步骤 4：访问服务

在浏览器中访问：

```
https://go.study.org:8443/login
```

**首次访问提示**：
- 浏览器会显示"不安全"警告（因为是自签名证书）
- 点击"高级"→"继续访问 go.study.org（不安全）"
- 可以选择将证书添加到浏览器的信任存储以避免重复提示

## IP 地址变化处理

如果服务器 IP 地址发生变化（例如从 172.28.186.65 变为 192.168.1.100）：

**只需修改客户端 hosts 文件**：

```
# 修改前
172.28.186.65  go.study.org

# 修改后
192.168.1.100  go.study.org
```

**无需重新生成证书或重启服务！**

## 证书验证

查看证书内容：

```bash
openssl x509 -in backend/configs/certs/server.crt -text -noout | grep -A 5 "Subject Alternative Name"
```

预期输出：

```
X509v3 Subject Alternative Name:
    DNS:go.study.org, DNS:localhost, DNS:*.local
```

## 多客户端配置

如果有多台客户端需要访问，每台客户端都需要：

1. 在 hosts 文件中添加域名映射
2. 导入并信任证书（可选，避免每次提示）

### Windows 导入证书

1. 双击 `backend/configs/certs/server.crt`
2. 选择"安装证书"
3. 选择"本地计算机"→"受信任的根证书颁发机构"
4. 完成导入

### Mac 导入证书

```bash
# 添加到系统钥匙串
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain backend/configs/certs/server.crt
```

### Linux 导入证书

```bash
# 复制证书到系统信任库
sudo cp backend/configs/certs/server.crt /usr/local/share/ca-certificates/go-study2.crt
sudo update-ca-certificates
```

## 生产环境建议

对于生产环境，建议使用正式域名和 CA 签发的证书：

1. 购买域名（如 `go-study.example.com`）
2. 使用 Let's Encrypt 获取免费证书
3. 或使用云服务商提供的证书

Let's Encrypt 示例：

```bash
# 安装 certbot
sudo apt-get install certbot

# 获取证书
sudo certbot certonly --standalone -d go-study.example.com

# 证书位置
# /etc/letsencrypt/live/go-study.example.com/fullchain.pem
# /etc/letsencrypt/live/go-study.example.com/privkey.pem
```

然后修改 `backend/configs/config.yaml`：

```yaml
https:
  enabled: true
  port: 8443
  certFile: "/etc/letsencrypt/live/go-study.example.com/fullchain.pem"
  keyFile: "/etc/letsencrypt/live/go-study.example.com/privkey.pem"
  insecureSkipVerify: false
```

## 常见问题

### Q: 为什么不使用 IP 地址生成证书？

A: 因为 IP 地址可能会变化，每次变化都需要重新生成证书并重启服务。使用域名方案，只需修改客户端 hosts 文件即可。

### Q: 浏览器显示"不安全"警告怎么办？

A: 这是正常现象，因为使用的是自签名证书。可以：
1. 点击"高级"→"继续访问"
2. 将证书导入到系统的信任存储

### Q: 如何信任自签名证书？

A: 参考"多客户端配置"部分的证书导入步骤。

### Q: 证书有效期多久？

A: 默认 365 天。到期后需要重新生成证书。

### Q: 可以自定义域名吗？

A: 可以。使用 `./scripts/generate-cert.sh -d your.domain` 指定自定义域名。
