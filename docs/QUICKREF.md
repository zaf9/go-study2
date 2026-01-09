# 打包部署快速参考

## 一分钟上手

### 本地打包

```bash
# 方式1: 使用打包脚本（推荐）
./scripts/package.sh

# 方式2: 使用 Makefile
make package VERSION=1.0.0

# 输出文件
# build/go-study2-1.0.0-20250109_120000-linux-amd64.tar.gz
```

### 服务器部署

```bash
# 1. 上传安装包
scp build/go-study2-*.tar.gz user@server:/tmp/

# 2. 解压到指定目录（注意使用 -C 参数）
tar -xf go-study2-*.tar.gz -C /opt/hadoop

# 3. 运行安装脚本
cd /opt/hadoop/go-study2
sudo ./scripts/install.sh

# 4. 启动服务
sudo systemctl start go-study2

# 5. 查看状态
sudo systemctl status go-study2
```

## 目录结构

```
go-study2/
├── bin/
│   └── go-study2           # 主程序（可执行）
├── configs/                # 配置文件
│   ├── config.yaml
│   └── logger.yaml
├── static/out/             # 前端静态文件
├── quiz_data/              # 题库数据
├── data/                   # 数据库（运行时生成）
├── logs/                   # 日志（运行时生成）
├── scripts/                # 管理脚本
│   ├── install.sh          # 安装（systemd）
│   ├── uninstall.sh        # 卸载
│   ├── start.sh            # 启动（前台）
│   ├── stop.sh             # 停止
│   └── logging.sh          # 日志函数库
├── VERSION                 # 版本配置
└── README.txt              # 使用说明
```

## 命令速查

### 打包相关

| 命令 | 说明 |
|------|------|
| `./scripts/package.sh` | 执行打包 |
| `make package` | 使用 Makefile 打包 |
| `make clean` | 清理构建产物 |
| `make build` | 本地构建（不打包） |
| `make deploy-local` | 本地模拟部署测试 |
| `make info` | 显示构建信息 |
| `make version` | 显示版本信息 |

### 服务管理

| 命令 | 说明 |
|------|------|
| `systemctl start go-study2` | 启动服务 |
| `systemctl stop go-study2` | 停止服务 |
| `systemctl restart go-study2` | 重启服务 |
| `systemctl status go-study2` | 查看状态 |
| `systemctl enable go-study2` | 开机自启 |
| `systemctl disable go-study2` | 禁用自启 |
| `journalctl -u go-study2 -f` | 查看日志 |

### 直接启动（不使用 systemd）

| 命令 | 说明 |
|------|------|
| `./scripts/start.sh` | 前台启动服务 |
| `./scripts/stop.sh` | 停止服务 |

### 故障排查

```bash
# 查看服务状态
sudo systemctl status go-study2

# 查看日志（最近100行）
sudo journalctl -u go-study2 -n 100

# 查看文件日志
tail -f /opt/hadoop/go-study2/logs/go-study2.log

# 检查端口
sudo netstat -tuln | grep 8080

# 检查进程
ps aux | grep go-study2

# 测试 API
curl http://localhost:8080/api/v1/topics
```

## 环境变量

| 变量 | 说明 | 示例 |
|------|------|------|
| `JWT_SECRET` | JWT 密钥（≥32 字符） | `export JWT_SECRET="my-secret-key-32-chars-long"` |

## 配置说明

主配置文件: `configs/config.yaml`

```yaml
# 重要配置项
http:
  port: 8080                  # 监听端口
server:
  host: "0.0.0.0"             # 监听地址（0.0.0.0 允许外部访问）
jwt:
  secret: "YOUR_SECRET"       # JWT 密钥（必须修改！）
database:
  path: "./data/gostudy.db"   # 数据库路径（相对路径）
static:
  enabled: true
  path: "./static/out"        # 前端静态文件路径（相对路径）
  spaFallback: true
```

修改配置后重启服务:
```bash
sudo systemctl restart go-study2
```

## 默认配置

| 项目 | 默认值 |
|------|--------|
| 监听端口 | 8080 |
| 管理员账号 | admin |
| 管理员密码 | GoStudy@123 |
| 数据库 | SQLite (嵌入式) |
| 服务名称 | go-study2 |
| 主程序 | bin/go-study2 |

## 升级流程

```bash
# 1. 停止服务
sudo systemctl stop go-study2

# 2. 备份数据
sudo cp /opt/hadoop/go-study2/data/gostudy.db /opt/hadoop/go-study2/data/gostudy.db.backup

# 3. 解压新版本覆盖文件
# 注意：只覆盖 bin/、static/、quiz_data/，保留 data/、configs/
tar -xf go-study2-2.0.0.tar.gz -C /tmp
cd /tmp/go-study2
sudo cp -r bin /opt/hadoop/go-study2/
sudo cp -r static /opt/hadoop/go-study2/
sudo cp -r quiz_data /opt/hadoop/go-study2/

# 4. 重启服务
sudo systemctl start go-study2
```

## 卸载

```bash
cd /opt/hadoop/go-study2
sudo ./scripts/uninstall.sh
```

## 日志输出格式

所有脚本使用统一的日志格式：

```
YYYYMMDD-HHMMSS - [script_name]: message

示例：
20250109-143022 - [install.sh]: 开始安装 Go-Study2...
20250109-143023 - [install.sh]: [INFO] 项目根目录: /opt/hadoop/go-study2
20250109-143024 - [install.sh]: [ERROR] 请使用 root 权限运行
```

## 版本管理

### 查看版本

```bash
# 方式1: 查看 VERSION 文件
cat /opt/hadoop/go-study2/VERSION

# 方式2: 使用 make 命令
make version

# 方式3: 查看二进制版本（如果支持）
/opt/hadoop/go-study2/bin/go-study2 --version
```

### 修改版本

编辑项目根目录的 `VERSION` 文件：

```bash
# VERSION 文件内容
VERSION="1.0.0"

# 修改为
VERSION="1.1.0"
```

### 语义化版本规则

- **MAJOR** (主版本): 不兼容的 API 变更
- **MINOR** (次版本): 向下兼容的功能新增
- **PATCH** (修订号): 向下兼容的问题修正

示例:
- `1.0.0` → `1.0.1`: 修复 bug
- `1.0.0` → `1.1.0`: 新增功能
- `1.0.0` → `2.0.0`: 重大变更

## 更多文档

- [完整打包文档](./PACKAGING.md)
- [部署指南](./DEPLOYMENT.md)
- [API 文档](./API.md)
