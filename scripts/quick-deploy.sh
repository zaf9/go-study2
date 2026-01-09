#!/bin/bash

# Go-Study2 快速部署脚本
# 用法: curl -fsSL https://your-domain.com/deploy.sh | sudo bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $*"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# 配置
INSTALL_DIR="${INSTALL_DIR:-/opt/gostudy}"
DOWNLOAD_URL="${DOWNLOAD_URL:-https://github.com/yourusername/go-study2/releases/latest/download}"
PACKAGE_NAME="gostudy-linux-amd64.tar.gz"
SERVICE_NAME="gostudy"

# 检查 root
if [ "$EUID" -ne 0 ]; then
    log_error "请使用 root 权限运行"
    exit 1
fi

log_info "=== Go-Study2 快速部署 ==="

# 检测系统
if [ ! -f /etc/os-release ]; then
    log_error "不支持的操作系统"
    exit 1
fi

source /etc/os-release
log_info "检测到系统: $PRETTY_NAME"

# 检查架构
ARCH=$(uname -m)
if [ "$ARCH" != "x86_64" ]; then
    log_error "不支持的架构: $ARCH (仅支持 x86_64)"
    exit 1
fi

# 安装依赖
log_info "检查系统依赖..."
if ! command -v curl &> /dev/null; then
    log_info "安装 curl..."
    apt-get update -qq
    apt-get install -y curl
fi

# 下载安装包
log_info "下载安装包..."
cd /tmp
curl -fsSL "$DOWNLOAD_URL/$PACKAGE_NAME" -o "$PACKAGE_NAME"

# 校验文件（如果有 sha256）
if curl -fsSL "$DOWNLOAD_URL/$PACKAGE_NAME.sha256" -o "$PACKAGE_NAME.sha256"; then
    log_info "校验文件完整性..."
    if ! sha256sum -c "$PACKAGE_NAME.sha256"; then
        log_error "文件校验失败"
        exit 1
    fi
fi

# 解压
log_info "解压安装包..."
tar -xzf "$PACKAGE_NAME"
cd gostudy

# 运行安装
log_info "运行安装脚本..."
bash ./scripts/install.sh

# 完成
log_info "=== 部署完成 ==="
echo ""
echo "访问: http://localhost:8080"
echo "用户: admin"
echo "密码: GoStudy@123"
echo ""
echo "管理命令:"
echo "  启动: systemctl start $SERVICE_NAME"
echo "  停止: systemctl stop $SERVICE_NAME"
echo "  状态: systemctl status $SERVICE_NAME"
echo "  日志: journalctl -u $SERVICE_NAME -f"
echo ""
