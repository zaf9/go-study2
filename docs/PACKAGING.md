# Linux 打包部署方案

## 📦 打包方案概述

本方案将 Go-Study2 项目打包成一个独立的 `tar.gz` 安装包，实现**一键部署**到 Linux 服务器。

### 设计特点

- ✅ **单文件部署**: 一个 tar.gz 包含所有运行所需文件
- ✅ **前后端集成**: 前端静态文件已编译并嵌入后端目录
- ✅ **无外部依赖**: 二进制静态链接，无需安装 Go/Node.js
- ✅ **自动安装**: 提供 systemd 服务管理脚本
- ✅ **开箱即用**: 解压后运行安装脚本即可使用

### 架构说明

```
┌─────────────────────────────────────────┐
│         gostudy-1.0.0.tar.gz           │
├─────────────────────────────────────────┤
│  gostudy/                              │
│  ├── gostudy              ◄── 主程序   │ (35MB)
│  ├── configs/             ◄── 配置     │ (10KB)
│  ├── static/out/          ◄── 前端     │ (5MB)
│  ├── quiz_data/           ◄── 题库     │ (1MB)
│  ├── data/                ◄── 数据库   │ (运行时生成)
│  ├── scripts/             ◄── 脚本     │
│  │   ├── install.sh       ◄── 安装     │
│  │   └── uninstall.sh     ◄── 卸载     │
│  └── README.txt           ◄── 文档     │
└─────────────────────────────────────────┘
```

## 🚀 使用方法

### 本地打包

```bash
# 1. 进入项目根目录
cd /path/to/go-study2

# 2. 执行打包脚本
./scripts/package.sh

# 3. 输出文件
# build/gostudy-1.0.0-20250109_120000-linux-amd64.tar.gz
# build/gostudy-1.0.0-20250109_120000-linux-amd64.tar.gz.sha256
```

### 服务器部署

```bash
# 1. 上传安装包到服务器
scp build/gostudy-*.tar.gz user@server:/tmp/

# 2. SSH 登录服务器
ssh user@server

# 3. 解压安装包
cd /tmp
tar -xzf gostudy-*.tar.gz
cd gostudy

# 4. 运行安装脚本（需要 root 权限）
sudo ./scripts/install.sh

# 5. 启动服务
sudo systemctl start gostudy

# 6. 设置开机自启
sudo systemctl enable gostudy

# 7. 查看运行状态
sudo systemctl status gostudy
```

### 访问应用

```
URL: http://your-server-ip:8080
默认管理员: admin / GoStudy@123
```

## 📋 环境要求

### 服务器要求

- **操作系统**: Linux (x86_64)
  - ✅ Ubuntu 18.04+
  - ✅ CentOS 7+
  - ✅ Debian 9+
  - ✅ Amazon Linux 2
- **权限**: root 或 sudo
- **端口**: 8080 (可配置)
- **磁盘空间**: ≥200MB
- **内存**: ≥512MB

### 编译环境要求（仅打包时需要）

- **Go**: 1.24.5+
- **Node.js**: 18+
- **npm**: 最新版

## 🔧 配置说明

### 环境变量

| 变量 | 说明 | 必需 | 默认值 |
|------|------|------|--------|
| `JWT_SECRET` | JWT 签名密钥（≥32 字符） | ❌ | 内置不安全密钥 |
| `INSTALL_DIR` | 安装目录 | ❌ | `/opt/gostudy` |

**重要**: 生产环境务必设置 `JWT_SECRET`！

```bash
# 方式1: 环境变量
export JWT_SECRET="your-super-secret-key-at-least-32-chars-long"
./scripts/install.sh

# 方式2: 修改配置文件
vim /opt/gostudy/configs/config.yaml
# 找到 jwt.secret 并修改

# 方式3: 启动时传入
sudo systemctl start gostudy
# 编辑 /etc/systemd/system/gostudy.service
# 添加: Environment=JWT_SECRET=your-secret
```

### 配置文件

主配置文件: `/opt/gostudy/configs/config.yaml`

```yaml
http:
  port: 8080          # 修改监听端口

server:
  host: "0.0.0.0"     # 监听所有网卡（默认 127.0.0.1）

database:
  path: "./data/gostudy.db"

jwt:
  secret: "YOUR_SECRET_HERE"  # 修改为安全的密钥

static:
  enabled: true
  path: "./static/out"        # 前端静态文件路径
  spaFallback: true
```

修改配置后需重启服务:

```bash
sudo systemctl restart gostudy
```

## 📂 目录结构

安装后的目录结构:

```
/opt/gostudy/
├── gostudy                  # 主程序（可执行）
├── configs/
│   ├── config.yaml          # 主配置
│   └── logger.yaml          # 日志配置
├── static/
│   └── out/                 # Next.js 静态导出
│       ├── _next/           # Next.js 资源
│       ├── topics/          # 学习页面
│       ├── quiz-center/     # 测验页面
│       └── index.html       # 首页
├── quiz_data/               # 题库数据（YAML）
│   ├── index.yaml
│   ├── lexical_elements/
│   ├── constants/
│   ├── variables/
│   └── types/
├── data/                    # 数据库目录
│   └── gostudy.db           # SQLite 数据库
├── logs/                    # 日志目录
│   ├── gostudy.log          # 应用日志
│   └── error.log            # 错误日志
├── scripts/
│   ├── install.sh           # 安装脚本
│   └── uninstall.sh         # 卸载脚本
├── README.txt               # 使用说明
└── VERSION.txt              # 版本信息
```

## 🛠️ 服务管理

### systemd 服务命令

```bash
# 启动服务
sudo systemctl start gostudy

# 停止服务
sudo systemctl stop gostudy

# 重启服务
sudo systemctl restart gostudy

# 查看状态
sudo systemctl status gostudy

# 开机自启
sudo systemctl enable gostudy

# 禁用自启
sudo systemctl disable gostudy

# 查看日志（实时）
sudo journalctl -u gostudy -f

# 查看最近 100 条日志
sudo journalctl -u gostudy -n 100

# 查看文件日志
tail -f /opt/gostudy/logs/gostudy.log
```

### 手动启动（调试用）

```bash
cd /opt/gostudy
./gostudy -d
```

## 🔍 故障排查

### 常见问题

#### 1. 端口被占用

```bash
# 查看端口占用
sudo netstat -tuln | grep 8080
# 或
sudo lsof -i :8080

# 解决方案：
# 方案1: 停止占用进程
sudo kill <PID>

# 方案2: 修改配置文件端口
vim /opt/gostudy/configs/config.yaml
# http.port: 8081
sudo systemctl restart gostudy
```

#### 2. 权限问题

```bash
# 检查文件权限
ls -la /opt/gostudy/gostudy

# 修复权限
sudo chmod +x /opt/gostudy/gostudy
sudo chown -R root:root /opt/gostudy
```

#### 3. JWT 校验失败

```bash
# 检查 JWT_SECRET
grep -A 5 "jwt:" /opt/gostudy/configs/config.yaml

# 设置环境变量
export JWT_SECRET="your-secure-secret-at-le-32-chars"
sudo systemctl restart gostudy
```

#### 4. 静态文件 404

```bash
# 检查前端文件
ls -la /opt/gostudy/static/out/

# 检查配置
grep -A 3 "static:" /opt/gostudy/configs/config.yaml
# 确保: static.enabled=true, static.path=./static/out
```

#### 5. 数据库错误

```bash
# 检查数据库文件
ls -la /opt/gostudy/data/

# 检查数据库权限
sudo chmod 644 /opt/gostudy/data/gostudy.db
sudo chmod 755 /opt/gostudy/data/
```

### 日志查看

```bash
# systemd 日志
sudo journalctl -u gostudy -f

# 应用日志
tail -f /opt/gostudy/logs/gostudy.log

# 错误日志
tail -f /opt/gostudy/logs/error.log
```

## 🔄 升级指南

### 升级步骤

```bash
# 1. 备份数据库
sudo cp /opt/gostudy/data/gostudy.db /opt/gostudy/data/gostudy.db.backup

# 2. 停止服务
sudo systemctl stop gostudy

# 3. 上传新版本并解压
cd /tmp
tar -xzf gostudy-2.0.0.tar.gz
cd gostudy

# 4. 替换程序文件（保留数据和配置）
sudo cp gostudy /opt/gostudy/gostudy
sudo cp -r static/out /opt/gostudy/static/
sudo cp -r quiz_data /opt/gostudy/
sudo chmod +x /opt/gostudy/gostudy

# 5. 重启服务
sudo systemctl start gostudy

# 6. 验证
sudo systemctl status gostudy
```

### 回滚

```bash
# 1. 停止服务
sudo systemctl stop gostudy

# 2. 恢复数据库
sudo cp /opt/gostudy/data/gostudy.db.backup /opt/gostudy/data/gostudy.db

# 3. 恢复旧版本程序（如有备份）
# sudo cp /path/to/old/gostudy /opt/gostudy/gostudy

# 4. 重启服务
sudo systemctl start gostudy
```

## 🗑️ 卸载

```bash
# 1. 运行卸载脚本
cd /opt/gostudy
sudo ./scripts/uninstall.sh

# 2. 根据提示选择是否删除数据

# 3. 手动清理（如果卸载脚本未执行）
sudo systemctl stop gostudy
sudo systemctl disable gostudy
sudo rm /etc/systemd/system/gostudy.service
sudo systemctl daemon-reload
sudo rm -rf /opt/gostudy
```

## 🔒 安全建议

### 生产环境检查清单

- [ ] 修改默认管理员密码
- [ ] 设置强 JWT_SECRET（≥32 字符，随机生成）
- [ ] 配置 HTTPS（使用反向代理 Nginx/Caddy）
- [ ] 配置防火墙（仅开放必要端口）
- [ ] 定期备份数据库
- [ ] 限制日志文件大小（使用 logrotate）
- [ ] 更新系统补丁
- [ ] 监控服务运行状态

### 生成安全的 JWT_SECRET

```bash
# 方法1: 使用 openssl
openssl rand -base64 32

# 方法2: 使用 /dev/urandom
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1

# 方法3: 使用 uuidgen（需拼接）
echo "$(uuidgen)$(uuidgen)" | cut -c1-32
```

## 🌐 反向代理配置

### Nginx 配置

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Caddy 配置（自动 HTTPS）

```
your-domain.com {
    reverse_proxy localhost:8080
}
```

## 📊 性能优化

### 1. 数据库优化

```yaml
# configs/config.yaml
database:
  pragmas:
    - "journal_mode=WAL"
    - "cache_size=-64000"    # 64MB 缓存
    - "synchronous=NORMAL"
```

### 2. 日志管理

创建 `/etc/logrotate.d/gostudy`:

```
/opt/gostudy/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 root root
    sharedscripts
    postrotate
        systemctl reload gostudy >/dev/null 2>&1 || true
    endscript
}
```

### 3. 系统资源限制

编辑 `/etc/systemd/system/gostudy.service`:

```ini
[Service]
# 限制内存使用
MemoryLimit=512M
# 限制文件描述符
LimitNOFILE=65536
```

## 📝 打包脚本详解

### 脚本流程

```
1. 前端构建
   ├── npm install
   ├── npm run build
   └── npm run export → frontend/out/

2. 后端编译
   ├── gofmt -w .
   ├── go vet ./...
   ├── go test ./...
   └── go build -ldflags="-w -s" → gostudy (静态链接)

3. 文件组织
   ├── 复制配置文件
   ├── 复制前端静态文件
   ├── 复制题库数据
   └── 修改配置路径

4. 创建脚本
   ├── install.sh (安装脚本)
   ├── uninstall.sh (卸载脚本)
   └── README.txt (使用说明)

5. 打包压缩
   ├── 创建目录结构
   └── tar -czf → gostudy-VERSION.tar.gz

6. 生成校验和
   └── sha256sum → *.tar.gz.sha256
```

### 关键技术点

- **静态链接**: `CGO_ENABLED=1` 确保 SQLite 嵌入
- **路径修正**: `sed` 替换配置文件中的相对路径
- **权限设置**: `chmod +x` 确保脚本可执行
- **版本注入**: `-ldflags` 注入版本和构建时间
- **SHA256 校验**: 确保传输完整性

## 🎯 下一步

- [ ] 添加 CI/CD 自动打包（GitHub Actions）
- [ ] 支持 ARM64 架构（树莓派等）
- [ ] 添加健康检查接口
- [ ] 提供一键安装脚本（curl 一行命令）
- [ ] 支持多版本管理
