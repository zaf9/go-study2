# 打包方案总结

## 方案概述

为 Go-Study2 项目设计了**单文件部署方案**，将前后端集成打包成一个 `tar.gz` 文件，实现 Linux 服务器的一键部署。

---

## 核心设计

### 1. 架构决策

`★ Insight ─────────────────────────────────────`
**为什么选择 tar.gz 单文件部署？**
- **零依赖**: 二进制静态链接 + SQLite 嵌入式，无需安装 Go/Node/数据库
- **简化运维**: 一个文件包含所有运行时文件，降低部署复杂度
- **快速回滚**: 保留旧版本压缩包，出现问题可快速回退
- **版本管理**: 文件名包含版本和时间戳，便于追踪
`─────────────────────────────────────────────────`

### 2. 打包流程

```
┌──────────────────────────────────────────────┐
│ 1. 前端构建                                   │
│    npm install → build → export              │
│    产出: frontend/out/ (~5MB)                │
└──────────────────┬───────────────────────────┘
                   │
┌──────────────────▼───────────────────────────┐
│ 2. 后端编译                                   │
│    gofmt → vet → test → build               │
│    产出: gostudy 二进制 (~35MB)             │
└──────────────────┬───────────────────────────┘
                   │
┌──────────────────▼───────────────────────────┐
│ 3. 文件组织                                   │
│    - 复制配置文件                             │
│    - 复制前端静态文件                         │
│    - 复制题库数据                             │
│    - 修改配置路径为生产路径                   │
└──────────────────┬───────────────────────────┘
                   │
┌──────────────────▼───────────────────────────┐
│ 4. 脚本生成                                   │
│    - install.sh (安装脚本)                   │
│    - uninstall.sh (卸载脚本)                 │
│    - README.txt (使用说明)                   │
└──────────────────┬───────────────────────────┘
                   │
┌──────────────────▼───────────────────────────┐
│ 5. 打包压缩                                   │
│    tar -czf → gostudy-VERSION.tar.gz        │
│    sha256sum → *.tar.gz.sha256              │
└──────────────────────────────────────────────┘
```

### 3. 文件清单

| 文件/目录 | 大小 | 说明 |
|-----------|------|------|
| `gostudy` | 35MB | 主程序（静态链接 ELF） |
| `configs/` | 10KB | 配置文件（config.yaml, logger.yaml） |
| `static/out/` | 5MB | Next.js 静态导出 |
| `quiz_data/` | 1MB | YAML 题库数据 |
| `scripts/` | 5KB | 安装/卸载脚本 |
| `data/` | 0B | 数据库目录（运行时生成） |
| `logs/` | 0B | 日志目录（运行时生成） |
| **总计** | **~41MB** | 压缩后 **10-15MB** |

---

## 使用方式

### 本地打包

```bash
# 方式1: 直接执行脚本
./scripts/package.sh

# 方式2: 使用 Makefile
make package VERSION=1.0.0

# 输出文件
# build/gostudy-1.0.0-20250109_120000-linux-amd64.tar.gz
# build/gostudy-1.0.0-20250109_120000-linux-amd64.tar.gz.sha256
```

### 服务器部署

```bash
# 1. 上传安装包
scp build/gostudy-*.tar.gz user@server:/tmp/

# 2. SSH 登录
ssh user@server

# 3. 解压并安装
cd /tmp
tar -xzf gostudy-*.tar.gz
cd gostudy
sudo ./scripts/install.sh

# 4. 启动服务
sudo systemctl start gostudy
sudo systemctl enable gostudy

# 5. 验证
curl http://localhost:8080/api/v1/topics
```

---

## 脚本说明

### scripts/package.sh

**功能**: 打包构建脚本

**关键步骤**:
1. 构建前端静态文件
2. 编译后端二进制（Linux x86_64，静态链接）
3. 复制配置文件并修改路径
4. 复制前端静态文件和题库数据
5. 生成安装/卸载脚本
6. 打包成 tar.gz
7. 生成 SHA256 校验和

**版本注入**:
```bash
-ldflags="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME"
```

### scripts/install.sh

**功能**: 服务器安装脚本

**执行步骤**:
1. 检查 root 权限
2. 检查端口占用
3. 复制文件到 `/opt/gostudy`
4. 创建 systemd 服务
5. 重新加载 systemd

**生成服务**: `/etc/systemd/system/gostudy.service`

### scripts/uninstall.sh

**功能**: 服务器卸载脚本

**执行步骤**:
1. 停止服务
2. 禁用开机自启
3. 删除 systemd 服务
4. 询问是否删除数据目录

### scripts/quick-deploy.sh

**功能**: 快速部署脚本（curl 一键安装）

**使用场景**: 服务器上直接执行远程脚本

```bash
curl -fsSL https://your-domain.com/deploy.sh | sudo bash
```

---

## CI/CD 集成

### GitHub Actions (.github/workflows/release.yml)

**触发条件**:
- 推送 tag: `v*`
- 手动触发（指定版本）

**工作流程**:
1. 检出代码
2. 设置 Go + Node.js 环境
3. 安装依赖并构建前端
4. 运行测试
5. 编译后端（Linux x86_64）
6. 组织打包文件
7. 生成 tar.gz 和 SHA256
8. 创建 GitHub Release
9. 上传构建产物

**使用方式**:
```bash
# 推送 tag 触发构建
git tag v1.0.0
git push origin v1.0.0

# 或在 GitHub Actions 页面手动触发
```

---

## 配置调整

### 生产环境配置

安装后的配置文件: `/opt/gostudy/configs/config.yaml`

**关键配置**:

```yaml
# 1. 监听地址和端口
http:
  port: 8080
server:
  host: "0.0.0.0"  # 允许外部访问

# 2. JWT 密钥（重要！）
jwt:
  secret: "your-secure-secret-at-least-32-characters-long"

# 3. 数据库路径
database:
  path: "./data/gostudy.db"

# 4. 前端静态文件
static:
  enabled: true
  path: "./static/out"
  spaFallback: true
```

### 环境变量

| 变量 | 说明 | 推荐方式 |
|------|------|----------|
| `JWT_SECRET` | JWT 签名密钥 | systemd Environment |
| `INSTALL_DIR` | 安装目录 | install.sh 参数 |

---

## 故障排查

### 常见问题

#### 1. 端口占用

```bash
# 检查端口
sudo netstat -tuln | grep 8080

# 解决方案1: 停止占用进程
sudo kill <PID>

# 解决方案2: 修改配置端口
sudo vim /opt/gostudy/configs/config.yaml
# http.port: 8081
sudo systemctl restart gostudy
```

#### 2. 权限问题

```bash
# 检查文件权限
ls -la /opt/gostudy/gostudy

# 修复权限
sudo chmod +x /opt/gostudy/gostudy
```

#### 3. 服务启动失败

```bash
# 查看详细日志
sudo journalctl -u gostudy -n 100

# 查看配置文件错误
sudo gostudy -t  # 测试配置（如果支持）
```

#### 4. 前端 404

```bash
# 检查静态文件
ls -la /opt/gostudy/static/out/

# 检查配置
grep -A 3 "static:" /opt/gostudy/configs/config.yaml
```

---

## 安全建议

### 生产环境检查清单

- [ ] **修改默认密码**: 首次登录后修改 admin 密码
- [ ] **设置强 JWT_SECRET**: 使用 `openssl rand -base64 32` 生成
- [ ] **配置 HTTPS**: 使用 Nginx/Caddy 反向代理
- [ ] **防火墙规则**: 仅开放必要端口
- [ ] **定期备份**: 数据库文件 `/opt/gostudy/data/gostudy.db`
- [ ] **日志轮转**: 配置 logrotate 防止日志过大
- [ ] **系统更新**: 定期更新系统补丁

### 生成安全密钥

```bash
# 使用 OpenSSL
openssl rand -base64 32

# 使用 /dev/urandom
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
```

### 反向代理配置

**Nginx**:
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Caddy** (自动 HTTPS):
```
your-domain.com {
    reverse_proxy localhost:8080
}
```

---

## 版本升级

### 升级流程

```bash
# 1. 备份数据
sudo systemctl stop gostudy
sudo cp /opt/gostudy/data/gostudy.db /opt/gostudy/data/gostudy.db.backup

# 2. 下载新版本
cd /tmp
wget https://github.com/xxx/gostudy-2.0.0.tar.gz
tar -xzf gostudy-2.0.0.tar.gz
cd gostudy

# 3. 替换程序文件（保留数据和配置）
sudo cp gostudy /opt/gostudy/
sudo cp -r static/out /opt/gostudy/static/
sudo cp -r quiz_data /opt/gostudy/

# 4. 重启服务
sudo systemctl start gostudy
sudo systemctl status gostudy
```

### 回滚

```bash
# 1. 停止服务
sudo systemctl stop gostudy

# 2. 恢复数据库
sudo cp /opt/gostudy/data/gostudy.db.backup /opt/gostudy/data/gostudy.db

# 3. 恢复旧版本程序
cd /tmp/gostudy-1.0.0
sudo cp gostudy /opt/gostudy/gostudy

# 4. 重启服务
sudo systemctl start gostudy
```

---

## 测试建议

### 本地测试

```bash
# 1. 执行打包
make package VERSION=1.0.0-test

# 2. 模拟部署
mkdir -p /tmp/gostudy-test
tar -xzf build/gostudy-*-test.tar.gz -C /tmp/gostudy-test

# 3. 检查文件结构
ls -la /tmp/gostudy-test/gostudy/

# 4. 手动安装测试
cd /tmp/gostudy-test/gostudy
sudo ./scripts/install.sh

# 5. 启动服务测试
sudo systemctl start gostudy
curl http://localhost:8080/api/v1/topics

# 6. 清理测试环境
sudo systemctl stop gostudy
sudo ./scripts/uninstall.sh
```

### Docker 测试（可选）

```dockerfile
FROM ubuntu:22.04

# 安装依赖
RUN apt-get update && apt-get install -y curl systemd

# 复制安装包
COPY build/gostudy-*.tar.gz /tmp/

# 安装
RUN cd /tmp && \
    tar -xzf gostudy-*.tar.gz && \
    cd gostudy && \
    ./scripts/install.sh

# 启动服务
CMD ["/bin/systemctl", "start", "gostudy"]
```

---

## 文件清单

### 新增文件

```
go-study2/
├── scripts/
│   ├── package.sh           # 打包脚本（核心）
│   ├── quick-deploy.sh      # 快速部署脚本
│   ├── install.sh           # 安装脚本（生成到包内）
│   └── uninstall.sh         # 卸载脚本（生成到包内）
├── docs/
│   ├── PACKAGING.md         # 完整打包文档
│   ├── QUICKREF.md          # 快速参考
│   └── PACKAGING_SUMMARY.md # 本文档
├── .github/workflows/
│   └── release.yml          # CI/CD 配置
└── Makefile                 # Makefile 构建
```

### 输出文件

```
build/
└── gostudy-1.0.0-20250109_120000-linux-amd64/
    ├── gostudy              # 主程序
    ├── configs/             # 配置
    ├── static/out/          # 前端
    ├── quiz_data/           # 题库
    ├── scripts/             # 脚本
    │   ├── install.sh
    │   └── uninstall.sh
    ├── data/                # 数据库（空）
    ├── logs/                # 日志（空）
    ├── README.txt           # 说明
    └── VERSION.txt          # 版本信息
```

---

## 下一步优化

- [ ] **ARM64 支持**: 添加 `GOARCH=arm64` 构建目标
- [ ] **多平台打包**: 支持 macOS、Windows
- [ ] **自动测试**: 集成端到端测试
- [ ] **一键安装**: 提供 curl 安装命令
- [ ] **健康检查**: 添加 `/health` 端点
- [ ] **监控集成**: Prometheus 指标导出
- [ ] **配置管理**: 支持环境变量覆盖
- [ ] **蓝绿部署**: 支持零停机升级

---

## 参考资料

- [完整打包文档](./PACKAGING.md)
- [快速参考](./QUICKREF.md)
- [部署指南](./DEPLOYMENT.md)
- [API 文档](./API.md)
- [GoFrame 官方文档](https://goframe.org)
- [Next.js 静态导出](https://nextjs.org/docs/app/building-your-application/deploying#static-exports)
