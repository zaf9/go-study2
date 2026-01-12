#!/bin/bash

# Go-Study2 Linux 打包脚本
# 功能：将项目打包成 tar.gz 安装包（前端集成到后端）
#
# 使用方法:
#   ./scripts/package.sh           # 默认：完整流程（包含测试）
#   ./scripts/package.sh --skip-tests  # 跳过测试（快速打包）
#   ./scripts/package.sh -h        # 显示帮助信息

set -e

# 默认参数
SKIP_TESTS=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        -h|--help)
            echo "用法: $0 [选项]"
            echo ""
            echo "选项:"
            echo "  --skip-tests    跳过前后端单元测试，加快打包速度"
            echo "  -h, --help      显示此帮助信息"
            echo ""
            echo "示例:"
            echo "  $0                # 完整打包（包含测试）"
            echo "  $0 --skip-tests   # 快速打包（跳过测试）"
            exit 0
            ;;
        *)
            echo "未知参数: $1"
            echo "使用 -h 或 --help 查看帮助信息"
            exit 1
            ;;
    esac
done

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志输出函数
printlog() {
    echo ""$(date +%Y%m%d-%H%M%S)" - [$(basename $0)]: $1"
}

log_info() {
    echo -e "${GREEN}$(printlog "[INFO] $*")${NC}"
}

log_warn() {
    echo -e "${YELLOW}$(printlog "[WARN] $*")${NC}"
}

log_error() {
    echo -e "${RED}$(printlog "[ERROR] $*")${NC}"
}

log_step() {
    echo -e "${BLUE}$(printlog "[STEP] $*")${NC}"
}

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_DIR="$PROJECT_ROOT/build"
PACKAGE_DIR="$BUILD_DIR/go-study2"

# 加载版本配置
VERSION_FILE="$PROJECT_ROOT/VERSION"
if [ -f "$VERSION_FILE" ]; then
    source "$VERSION_FILE"
else
    VERSION=${VERSION:-"1.0.0"}
fi

BUILD_TIME=$(date +"%Y%m%d_%H%M%S")
PACKAGE_NAME="go-study2-${VERSION}-${BUILD_TIME}-linux-amd64"
PACKAGE_FILE="$BUILD_DIR/${PACKAGE_NAME}.tar.gz"

# 服务名称
SERVICE_NAME="go-study2"

log_info "开始打包 Go-Study2 v${VERSION}"
log_info "构建目录: $BUILD_DIR"

# 显示构建模式
if [ "$SKIP_TESTS" = true ]; then
    log_warn "构建模式: 快速打包（跳过测试）"
else
    log_info "构建模式: 完整打包（包含测试）"
fi

# 清理并创建临时目录
log_step "初始化构建环境..."
rm -rf "$BUILD_DIR"
mkdir -p "$PACKAGE_DIR"
log_info "✓ 构建目录已创建: $PACKAGE_DIR"

# ========== 步骤1：构建前端 ==========
log_step "步骤 1/7: 构建前端静态文件..."
cd "$PROJECT_ROOT/frontend"
if ! command -v npm &> /dev/null; then
    log_error "npm 未找到，请安装 Node.js 18+"
    exit 1
fi

log_info "清理前端旧构建产物..."
rm -rf .next out
log_info "✓ 前端旧构建产物已清理"

log_info "安装前端依赖..."
npm install --silent 2>&1 | while IFS= read -r line; do log_info "  npm: $line"; done
log_info "✓ 前端依赖安装完成"

log_info "running lint..."
npm run lint 2>&1 | while IFS= read -r line; do log_info "  lint: $line"; done
log_info "✓ lint 检查通过"

if [ "$SKIP_TESTS" = true ]; then
    log_warn "跳过前端测试（--skip-tests 模式）"
else
    log_info "running tests with coverage..."
    npm run test -- --coverage 2>&1 | while IFS= read -r line; do log_info "  test: $line"; done
    log_info "✓ 测试完成"
fi

log_info "构建前端..."
# 强制使用相对路径（生产环境）：覆盖 .env.local 中的 localhost 配置
# 这样前端会向同源的 /api/v1 发送请求，而非固定的 localhost:8080
NEXT_PUBLIC_API_URL="" npm run build 2>&1 | while IFS= read -r line; do log_info "  build: $line"; done

# Next.js 14 使用 output: 'export' 时会自动生成静态文件到 out 目录
if [ ! -d "out" ]; then
    log_error "前端构建失败，out 目录不存在"
    log_info "提示: Next.js 14 配置了 output: 'export' 后会自动生成 out 目录"
    exit 1
fi

FRONTEND_SIZE=$(du -sh out | cut -f1)
log_info "✓ 前端构建成功，大小: $FRONTEND_SIZE"

# ========== 步骤2：构建后端 ==========
log_step "步骤 2/7: 构建后端二进制文件..."
cd "$PROJECT_ROOT/backend"
if ! command -v go &> /dev/null; then
    log_error "go 未找到，请安装 Go 1.24.5+"
    exit 1
fi

log_info "格式化 Go 代码..."
gofmt -w . 2>&1 | while IFS= read -r line; do log_info "  gofmt: $line"; done
log_info "✓ 代码格式化完成"

log_info "代码检查..."
if ! go vet ./... 2>&1 | tee /tmp/vet.log | while IFS= read -r line; do log_warn "  vet: $line"; done; then
    log_error "go vet 发现严重问题，终止构建"
    cat /tmp/vet.log
    exit 1
fi
log_info "✓ 代码检查通过"

if [ "$SKIP_TESTS" = true ]; then
    log_warn "跳过后端测试（--skip-tests 模式）"
else
    log_info "运行测试..."
    if ! go test -cover ./... 2>&1 | tee /tmp/test.log | while IFS= read -r line; do log_info "  test: $line"; done; then
        log_error "测试失败，终止构建"
        cat /tmp/test.log
        exit 1
    fi
fi

log_info "编译后端（Linux x86_64，完全静态链接，禁用 CGO）..."
mkdir -p "$PACKAGE_DIR/bin"

# 显式设置环境变量以禁用 CGO（兼容 CentOS 7.6 的 GLIBC 2.17）
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# 验证 CGO 已禁用
log_info "验证构建环境..."
log_info "  CGO_ENABLED=$CGO_ENABLED"
log_info "  GOOS=$GOOS"
log_info "  GOARCH=$GOARCH"

if [ "$CGO_ENABLED" != "0" ]; then
    log_error "CGO 未禁用，生成的二进制将无法在 CentOS 7.6 上运行"
    exit 1
fi

# 构建静态链接的二进制（移除 extldflags '-static' 因为纯 Go 构建不需要）
go build \
    -ldflags="-w -s -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" \
    -o "$PACKAGE_DIR/bin/go-study2" \
    main.go 2>&1 | tee /tmp/build.log | while IFS= read -r line; do log_info "  build: $line"; done

if [ ! -f "$PACKAGE_DIR/bin/go-study2" ]; then
    log_error "后端编译失败"
    cat /tmp/build.log
    exit 1
fi

chmod +x "$PACKAGE_DIR/bin/go-study2"
BACKEND_SIZE=$(du -h "$PACKAGE_DIR/bin/go-study2" | cut -f1)
log_info "✓ 后端编译成功，大小: $BACKEND_SIZE"

# 验证二进制是否为静态链接
log_info "验证二进制链接类型..."
if command -v file &> /dev/null; then
    FILE_INFO=$(file "$PACKAGE_DIR/bin/go-study2")
    log_info "  $FILE_INFO"

    # 检查是否为动态链接（不应包含 "dynamically linked"）
    if echo "$FILE_INFO" | grep -q "dynamically linked"; then
        log_error "二进制为动态链接，无法在 CentOS 7.6 上运行！"
        log_error "请确保 CGO_ENABLED=0 并且没有使用 CGO 依赖"
        exit 1
    fi

    # 检查是否为静态链接或纯 Go 二进制
    if echo "$FILE_INFO" | grep -q "statically linked" || echo "$FILE_INFO" | grep -q "not stripped"; then
        log_info "✓ 二进制链接类型验证通过（静态或纯 Go）"
    fi
fi

# 检查动态库依赖（应该没有）
if command -v ldd &> /dev/null; then
    LDD_OUTPUT=$(ldd "$PACKAGE_DIR/bin/go-study2" 2>&1 || true)
    if [ -n "$LDD_OUTPUT" ]; then
        log_warn "检测到动态库依赖："
        echo "$LDD_OUTPUT" | while IFS= read -r line; do log_warn "  $line"; done
        log_warn "虽然可能有动态依赖，但纯 Go 二进制通常会显示 'not a dynamic executable'"
    else
        log_info "✓ 无动态库依赖（纯 Go 二进制）"
    fi
fi

# ========== 步骤3：复制配置文件 ==========
log_step "步骤 3/7: 复制配置文件..."
mkdir -p "$PACKAGE_DIR/configs"
mkdir -p "$PACKAGE_DIR/configs/certs"

cp "$PROJECT_ROOT/backend/configs/config.yaml" "$PACKAGE_DIR/configs/"
log_info "✓ config.yaml 已复制"

cp "$PROJECT_ROOT/backend/configs/logger.yaml" "$PACKAGE_DIR/configs/"
log_info "✓ logger.yaml 已复制"

# 复制 SSL/TLS 证书文件（HTTPS 部署必需）
if [ -f "$PROJECT_ROOT/backend/configs/certs/server.crt" ] && [ -f "$PROJECT_ROOT/backend/configs/certs/server.key" ]; then
    cp "$PROJECT_ROOT/backend/configs/certs/server.crt" "$PACKAGE_DIR/configs/certs/"
    cp "$PROJECT_ROOT/backend/configs/certs/server.key" "$PACKAGE_DIR/configs/certs/"
    chmod 644 "$PACKAGE_DIR/configs/certs/server.crt"
    chmod 600 "$PACKAGE_DIR/configs/certs/server.key"
    log_info "✓ SSL 证书文件已复制"
else
    log_warn "⚠ SSL 证书文件未找到，HTTPS 功能将不可用"
fi

# 修改配置文件中的路径为相对路径（关键修改）
sed -i 's|path: "../frontend/out"|path: "./static/out"|g' "$PACKAGE_DIR/configs/config.yaml"
sed -i 's|path: "./data/gostudy.db"|path: "./data/gostudy.db"|g' "$PACKAGE_DIR/configs/config.yaml"
log_info "✓ 配置文件路径已调整为相对路径"

# 启用 HTTPS 模式（生产环境默认启用）
# 同时启用 insecureSkipVerify 以支持自签名证书
awk '
/^https:/ { in_https = 1 }
/^http:/ { in_https = 0 }
/^server:/ { in_https = 0 }
in_https && /^  enabled: false/ { print "  enabled: true"; replaced = 1; next }
in_https && /^  insecureSkipVerify: false/ { print "  insecureSkipVerify: true"; next }
{ print }
END { exit(replaced ? 0 : 1) }
' "$PACKAGE_DIR/configs/config.yaml" > "$PACKAGE_DIR/configs/config.yaml.tmp" || {
    log_error "HTTPS 模式替换失败，未能找到 https.enabled 配置项"
    rm -f "$PACKAGE_DIR/configs/config.yaml.tmp"
    exit 1
}

# 验证替换是否成功（检查整个文件中的配置）
if ! grep "^  enabled: true" "$PACKAGE_DIR/configs/config.yaml.tmp" | grep -q "true"; then
    log_error "HTTPS 模式验证失败，enabled 未正确设置"
    rm -f "$PACKAGE_DIR/configs/config.yaml.tmp"
    exit 1
fi

if ! grep "^  insecureSkipVerify: true" "$PACKAGE_DIR/configs/config.yaml.tmp" | grep -q "true"; then
    log_error "HTTPS 模式验证失败，insecureSkipVerify 未正确设置"
    rm -f "$PACKAGE_DIR/configs/config.yaml.tmp"
    exit 1
fi

mv "$PACKAGE_DIR/configs/config.yaml.tmp" "$PACKAGE_DIR/configs/config.yaml"
log_info "✓ HTTPS 模式已在生产配置中启用（含自签名证书支持）"

log_info "✓ 配置文件复制完成"

# ========== 步骤 3.5：生成 TLS 证书 ==========
log_step "步骤 3.5/7: 生成 TLS 证书..."

# 检查是否需要生成证书
if [ ! -f "$PROJECT_ROOT/backend/configs/certs/server.crt" ] || [ ! -f "$PROJECT_ROOT/backend/configs/certs/server.key" ]; then
    log_warn "证书文件不存在，将生成默认证书"
    chmod +x "$PROJECT_ROOT/scripts/generate-cert.sh"
    cd "$PROJECT_ROOT"
    ./scripts/generate-cert.sh -f
    cd "$SCRIPT_DIR/.."
fi

# 复制证书到打包目录
if [ -f "$PROJECT_ROOT/backend/configs/certs/server.crt" ] && [ -f "$PROJECT_ROOT/backend/configs/certs/server.key" ]; then
    cp "$PROJECT_ROOT/backend/configs/certs/server.crt" "$PACKAGE_DIR/configs/certs/"
    cp "$PROJECT_ROOT/backend/configs/certs/server.key" "$PACKAGE_DIR/configs/certs/"
    chmod 644 "$PACKAGE_DIR/configs/certs/server.crt"
    chmod 600 "$PACKAGE_DIR/configs/certs/server.key"
    log_info "✓ TLS 证书已生成并复制到打包目录"
else
    log_error "TLS 证书生成失败，HTTPS 功能将不可用"
fi

# ========== 步骤4：复制静态文件 ==========
log_step "步骤 4/7: 复制前端静态文件..."
mkdir -p "$PACKAGE_DIR/static"
cp -r "$PROJECT_ROOT/frontend/out" "$PACKAGE_DIR/static/"
STATIC_SIZE=$(du -sh "$PACKAGE_DIR/static/out" | cut -f1)
log_info "✓ 前端静态文件已复制，大小: $STATIC_SIZE"

# ========== 步骤5：复制题库数据 ==========
log_step "步骤 5/7: 复制题库数据..."
cp -r "$PROJECT_ROOT/backend/quiz_data" "$PACKAGE_DIR/"
QUIZ_SIZE=$(du -sh "$PACKAGE_DIR/quiz_data" | cut -f1)
log_info "✓ 题库数据已复制，大小: $QUIZ_SIZE"

# ========== 步骤6：创建版本配置 ==========
log_step "步骤 6/7: 创建版本配置..."

# 创建版本信息文件
cat > "$PACKAGE_DIR/VERSION" << EOF
# Go-Study2 版本配置
# 语义化版本: MAJOR.MINOR.PATCH
VERSION="${VERSION}"
BUILD_TIME="${BUILD_TIME}"
PLATFORM="linux-amd64"

# 构建信息
GO_VERSION="$(go version 2>/dev/null | awk '{print $3}' || echo "unknown")"
NODE_VERSION="$(node --version 2>/dev/null || echo "unknown")"
NPM_VERSION="$(npm --version 2>/dev/null || echo "unknown")"

# 项目信息
PROJECT_NAME="go-study2"
SERVICE_NAME="go-study2"
DESCRIPTION="Go语言学习平台"

# 更新日志
CHANGELOG="\
- 初始发布版本
- 支持Go语言学习模块
- 集成测验系统
- 前后端一体化部署
"
EOF

log_info "✓ 版本配置文件已创建: VERSION"

# ========== 步骤7：创建脚本和文档 ==========
log_step "步骤 7/7: 创建安装脚本和文档..."

# 创建必要的目录
mkdir -p "$PACKAGE_DIR/data"
mkdir -p "$PACKAGE_DIR/logs"
mkdir -p "$PACKAGE_DIR/scripts"

# 创建统一的日志函数库（供所有脚本使用）
cat > "$PACKAGE_DIR/scripts/logging.sh" << 'EOF'
#!/bin/bash

# 日志输出函数（统一格式）
printlog() {
    echo ""$(date +%Y%m%d-%H%M%S)" - [$(basename $0)]: $1"
}

log_info() {
    echo "$(printlog "[INFO] $1")"
}

log_warn() {
    echo "$(printlog "[WARN] $1")"
}

log_error() {
    echo "$(printlog "[ERROR] $1")"
}

log_debug() {
    if [ "$DEBUG" = "true" ]; then
        echo "$(printlog "[DEBUG] $1")"
    fi
}
EOF

chmod +x "$PACKAGE_DIR/scripts/logging.sh"
log_info "✓ 日志函数库已创建: scripts/logging.sh"

# 创建安装脚本
cat > "$PACKAGE_DIR/scripts/install.sh" << EOF
#!/bin/bash

# Go-Study2 安装脚本
# 支持解压到任意目录，使用相对路径运行

set -e

# 加载日志函数
SCRIPT_DIR="\$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
source "\$SCRIPT_DIR/logging.sh"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 服务名称
SERVICE_NAME="go-study2"

# 获取项目根目录（自动检测）
PROJECT_ROOT="\$(cd "\$SCRIPT_DIR/.." && pwd)"
log_info "项目根目录: \$PROJECT_ROOT"

# 检测安装路径类型并设置安全级别
detect_installation_path() {
    local install_path="\$1"

    if [[ "\$install_path" == /tmp* ]]; then
        log_warn "检测到安装路径在 /tmp 下" >&2
        log_warn "将调整安全设置以兼容 systemd" >&2
        echo "tmp"
    elif [[ "\$install_path" == /opt* ]] || [[ "\$install_path" == /usr/local* ]]; then
        log_info "推荐安装路径: \$install_path" >&2
        echo "recommended"
    elif [[ "\$install_path" == /home* ]]; then
        log_info "用户目录安装: \$install_path" >&2
        echo "homedir"
    else
        log_info "自定义安装路径: \$install_path" >&2
        echo "custom"
    fi
}

INSTALLATION_TYPE=\$(detect_installation_path "\$PROJECT_ROOT")

# 显示欢迎信息
echo -e "\${GREEN}=== Go-Study2 安装程序 ===\${NC}"
echo ""

log_info "开始安装 Go-Study2..."
log_info "安装目录: \$PROJECT_ROOT"

# 检查是否为 root
if [ "\$EUID" -ne 0 ]; then
    log_error "请使用 root 权限运行此脚本"
    echo "使用: sudo ./scripts/install.sh"
    exit 1
fi

log_info "Root 权限检查通过"

# 检查必要的命令
log_info "检查系统依赖..."
for cmd in systemctl tar; do
    if ! command -v \$cmd &> /dev/null; then
        log_error "缺少必要命令: \$cmd"
        exit 1
    fi
    log_info "✓ 找到命令: \$cmd"
done

# 检查端口是否被占用
log_info "检查端口 8080..."
if netstat -tuln 2>/dev/null | grep -q ":8080 "; then
    log_warn "端口 8080 已被占用"
    read -p "是否继续安装? (y/N): " -n 1 -r
    echo
    if [[ ! \$REPLY =~ ^[Yy]\$ ]]; then
        log_error "用户取消安装"
        exit 1
    fi
    log_info "用户选择继续安装"
else
    log_info "✓ 端口 8080 可用"
fi

# 设置权限
log_info "设置文件权限..."
chmod +x "\$PROJECT_ROOT/bin/go-study2"
chmod +x "\$PROJECT_ROOT/scripts/"*.sh
# 设置 SSL 证书权限（安全要求）
if [ -f "\$PROJECT_ROOT/configs/certs/server.crt" ]; then
    chmod 644 "\$PROJECT_ROOT/configs/certs/server.crt"
    log_info "✓ 证书文件权限已设置为 644"
fi
if [ -f "\$PROJECT_ROOT/configs/certs/server.key" ]; then
    chmod 600 "\$PROJECT_ROOT/configs/certs/server.key"
    log_info "✓ 私钥文件权限已设置为 600"
fi
log_info "✓ 文件权限设置完成"

# 创建必要的目录
log_info "创建运行时目录..."
mkdir -p "\$PROJECT_ROOT/data"
mkdir -p "\$PROJECT_ROOT/logs"
log_info "✓ 运行时目录已创建"

# 检查环境变量
log_info "检查环境变量..."
if [ -z "\$JWT_SECRET" ]; then
    log_warn "JWT_SECRET 未设置，自动生成随机密钥"
    JWT_SECRET=\$(openssl rand -base64 32 | tr -d '=+/' | cut -c1-32)
    log_info "✓ JWT_SECRET 已自动生成: \${JWT_SECRET:0:8}..."

    # 将生成的密钥保存到配置文件（使用 awk 直接替换 secret 的值）
    if [ -f "\$PROJECT_ROOT/configs/config.yaml" ]; then
        awk -v secret="\$JWT_SECRET" '
        /^jwt:/ { in_jwt = 1; print; next }
        in_jwt && /^  secret:/ {
            print "  secret: \"" secret "\""
            replaced = 1
            next
        }
        in_jwt && /^[a-z]/ { in_jwt = 0 }
        { print }
        END { exit(replaced ? 0 : 1) }
        ' "\$PROJECT_ROOT/configs/config.yaml" > "\$PROJECT_ROOT/configs/config.yaml.tmp" && \
        mv "\$PROJECT_ROOT/configs/config.yaml.tmp" "\$PROJECT_ROOT/configs/config.yaml" && \
        log_info "✓ JWT_SECRET 已保存到配置文件" || {
            log_warn "配置文件更新失败，请手动修改 configs/config.yaml 中的 jwt.secret"
        }
    fi
else
    log_info "✓ JWT_SECRET 已设置"
fi

# /tmp 路径警告
if [[ "\$PROJECT_ROOT" == /tmp* ]]; then
    echo ""
    echo -e "\${YELLOW}⚠️  警告：安装路径在 /tmp 下\${NC}"
    echo "  /tmp 目录在系统重启后会被清空"
    echo "  systemd 的 PrivateTmp 功能会导致命名空间冲突"
    echo ""
    echo "推荐安装目录："
    echo "  - /opt/go-study2       (生产环境)"
    echo "  - /usr/local/go-study2 (生产环境)"
    echo "  - /home/\$USER/go-study2 (开发环境)"
    echo ""
    read -p "是否继续安装? (y/N): " -n 1 -r
    echo
    if [[ ! \$REPLY =~ ^[Yy]\$ ]]; then
        log_error "用户取消安装"
        exit 1
    fi
    log_warn "用户选择继续 /tmp 安装，安全设置已调整"
fi

# 创建 systemd 服务
log_info "创建 systemd 服务..."

# 检测 systemd 版本
SYSTEMD_VERSION=\$(systemctl --version 2>/dev/null | head -1 | awk '{print \$2}' || echo "0")
log_info "检测到 systemd 版本: \${SYSTEMD_VERSION}"

# 根据安装路径类型和 systemd 版本选择安全配置
case "\$INSTALLATION_TYPE" in
    tmp)
        # /tmp 下安装：禁用 PrivateTmp，调整 ProtectSystem
        SECURITY_PRIVATE_TMP="false"
        SECURITY_PROTECT_SYSTEM="full"
        log_warn "⚠️  /tmp 路径安全限制：已禁用 PrivateTmp"
        ;;
    recommended)
        # /opt, /usr/local：根据 systemd 版本选择安全级别
        if [ "\${SYSTEMD_VERSION}" -lt 232 ]; then
            # CentOS 7 的 systemd 版本为 219，不支持 strict 和 ReadWritePaths
            SECURITY_PRIVATE_TMP="true"
            SECURITY_PROTECT_SYSTEM="full"
            USE_READ_WRITE_PATHS="false"
            log_warn "⚠️  旧版 systemd (\${SYSTEMD_VERSION})，使用兼容配置"
        else
            SECURITY_PRIVATE_TMP="true"
            SECURITY_PROTECT_SYSTEM="strict"
            USE_READ_WRITE_PATHS="true"
        fi
        ;;
    homedir|custom)
        # 其他路径：平衡安全性和兼容性
        SECURITY_PRIVATE_TMP="true"
        SECURITY_PROTECT_SYSTEM="full"
        USE_READ_WRITE_PATHS="false"
        ;;
esac

# 根据 systemd 版本选择日志输出方式
if [ "\${SYSTEMD_VERSION}" -lt 232 ]; then
    # CentOS 7 不支持 append:，使用文件重定向
    LOG_OUTPUT=""
    LOG_ERROR=""
    log_info "使用传统日志方式（systemd < 232）"
else
    # 新版 systemd 支持 append:
    LOG_OUTPUT="StandardOutput=append:\$PROJECT_ROOT/logs/go-study2.log"
    LOG_ERROR="StandardError=append:\$PROJECT_ROOT/logs/error.log"
    log_info "使用 append 日志方式（systemd >= 232）"
fi

cat > "/etc/systemd/system/\${SERVICE_NAME}.service" << EOM
[Unit]
Description=Go-Study2 Learning Platform
Documentation=https://github.com/yourusername/go-study2
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=\$PROJECT_ROOT
ExecStart=\$PROJECT_ROOT/bin/go-study2 -d
ExecReload=/bin/kill -HUP \$MAINPID
Restart=on-failure
RestartSec=5s
TimeoutStartSec=30

# 环境变量
Environment="JWT_SECRET=\${JWT_SECRET:-}"
Environment="INSTALL_DIR=\$PROJECT_ROOT"
Environment="INSTALLATION_TYPE=\${INSTALLATION_TYPE}"

\${LOG_OUTPUT}
\${LOG_ERROR}

# 安全设置（自适应 systemd 版本）
NoNewPrivileges=true
PrivateTmp=\${SECURITY_PRIVATE_TMP}
ProtectSystem=\${SECURITY_PROTECT_SYSTEM}
ProtectHome=true
\${USE_READ_WRITE_PATHS:+ReadWritePaths=}\${USE_READ_WRITE_PATHS:+\$PROJECT_ROOT/data }\${USE_READ_WRITE_PATHS:+\$PROJECT_ROOT/logs}

[Install]
WantedBy=multi-user.target
EOM

log_info "✓ systemd 服务文件已创建: /etc/systemd/system/\${SERVICE_NAME}.service"
log_info "  安全配置: PrivateTmp=\${SECURITY_PRIVATE_TMP}, ProtectSystem=\${SECURITY_PROTECT_SYSTEM}"

# 重新加载 systemd
log_info "重新加载 systemd 配置..."
systemctl daemon-reload
log_info "✓ systemd 配置已重新加载"

# 启用服务（但不启动）
log_info "启用开机自启..."
systemctl enable \$SERVICE_NAME 2>&1 | while IFS= read -r line; do log_info "  systemctl: \$line"; done
log_info "✓ 开机自启已启用"

# 安装完成
echo ""
log_info "================================"
log_info "安装完成！"
log_info "================================"
echo ""
echo "项目信息:"
echo "  安装目录: \$PROJECT_ROOT"
echo "  服务名称: \$SERVICE_NAME"
echo "  主程序: \$PROJECT_ROOT/bin/go-study2"
echo "  配置文件: \$PROJECT_ROOT/configs/config.yaml"
echo ""
echo "快速开始:"
echo "  启动服务: systemctl start \$SERVICE_NAME"
echo "  停止服务: systemctl stop \$SERVICE_NAME"
echo "  重启服务: systemctl restart \$SERVICE_NAME"
echo "  开机自启: systemctl enable \$SERVICE_NAME"
echo "  取消自启: systemctl disable \$SERVICE_NAME"
echo "  查看状态: systemctl status \$SERVICE_NAME"
echo "  查看日志: journalctl -u \$SERVICE_NAME -f"
echo "  文件日志: tail -f \$PROJECT_ROOT/logs/go-study2.log"
echo ""
echo "访问地址:"
echo "  HTTP: http://localhost:8080"
echo "  默认管理员: admin / GoStudy@123"
echo ""
echo "配置说明:"
echo "  编辑配置: vim \$PROJECT_ROOT/configs/config.yaml"
echo "  修改后重启: systemctl restart \$SERVICE_NAME"
echo ""
echo -e "\${YELLOW}注意: 请修改默认管理员密码和 JWT_SECRET！\${NC}"
echo ""

log_info "安装脚本执行完毕"
EOF

chmod +x "$PACKAGE_DIR/scripts/install.sh"
log_info "✓ 安装脚本已创建: scripts/install.sh"

# 创建卸载脚本
cat > "$PACKAGE_DIR/scripts/uninstall.sh" << EOF
#!/bin/bash

# Go-Study2 卸载脚本

set -e

# 加载日志函数
SCRIPT_DIR="\$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
source "\$SCRIPT_DIR/logging.sh"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# 服务名称
SERVICE_NAME="go-study2"

# 获取项目根目录
PROJECT_ROOT="\$(cd "\$SCRIPT_DIR/.." && pwd)"

echo -e "\${RED}=== Go-Study2 卸载程序 ===\${NC}"
echo ""
log_info "项目目录: \$PROJECT_ROOT"

# 检查是否为 root
if [ "\$EUID" -ne 0 ]; then
    log_error "请使用 root 权限运行此脚本"
    echo "使用: sudo ./scripts/uninstall.sh"
    exit 1
fi

# 确认卸载
read -p "确定要卸载 Go-Study2 吗? (y/N): " -n 1 -r
echo
if [[ ! \$REPLY =~ ^[Yy]\$ ]]; then
    log_info "用户取消卸载"
    exit 0
fi

log_info "开始卸载..."

# 停止服务
if systemctl is-active --quiet \$SERVICE_NAME; then
    log_info "停止服务..."
    systemctl stop \$SERVICE_NAME 2>&1 | while IFS= read -r line; do log_info "  \$line"; done
    log_info "✓ 服务已停止"
else
    log_info "服务未运行"
fi

# 禁用服务
if systemctl is-enabled --quiet \$SERVICE_NAME; then
    log_info "禁用开机自启..."
    systemctl disable \$SERVICE_NAME 2>&1 | while IFS= read -r line; do log_info "  \$line"; done
    log_info "✓ 开机自启已禁用"
else
    log_info "服务未启用"
fi

# 删除服务文件
log_info "删除 systemd 服务..."
rm -f "/etc/systemd/system/\${SERVICE_NAME}.service"
systemctl daemon-reload
log_info "✓ systemd 服务已删除"

# 询问是否删除数据
echo ""
read -p "是否删除数据库和日志? (y/N): " -n 1 -r
echo
if [[ \$REPLY =~ ^[Yy]\$ ]]; then
    log_info "删除所有文件..."
    rm -rf "\$PROJECT_ROOT"
    log_info "✓ 所有文件已删除"
else
    log_info "保留数据目录: \$PROJECT_ROOT"
    log_info "如需手动删除，请执行: rm -rf \$PROJECT_ROOT"
fi

echo ""
log_info "卸载完成"
echo ""
EOF

chmod +x "$PACKAGE_DIR/scripts/uninstall.sh"
log_info "✓ 卸载脚本已创建: scripts/uninstall.sh"

# 创建启动脚本（直接启动，不通过 systemd）
cat > "$PACKAGE_DIR/scripts/start.sh" << EOF
#!/bin/bash

# Go-Study2 启动脚本（前台运行）

set -e

# 加载日志函数
SCRIPT_DIR="\$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
source "\$SCRIPT_DIR/logging.sh"

# 获取项目根目录
PROJECT_ROOT="\$(cd "\$SCRIPT_DIR/.." && pwd)"

log_info "启动 Go-Study2..."
log_info "项目目录: \$PROJECT_ROOT"

# 检查端口
if netstat -tuln 2>/dev/null | grep -q ":8080 "; then
    log_error "端口 8080 已被占用"
    exit 1
fi

# 启动服务
log_info "执行: \$PROJECT_ROOT/bin/go-study2 -d"
cd "\$PROJECT_ROOT"
exec ./bin/go-study2 -d
EOF

chmod +x "$PACKAGE_DIR/scripts/start.sh"
log_info "✓ 启动脚本已创建: scripts/start.sh"

# 创建停止脚本
cat > "$PACKAGE_DIR/scripts/stop.sh" << EOF
#!/bin/bash

# Go-Study2 停止脚本

set -e

# 加载日志函数
SCRIPT_DIR="\$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
source "\$SCRIPT_DIR/logging.sh"

log_info "停止 Go-Study2..."

# 查找进程
PID=\$(ps aux | grep "[b]in/go-study2" | awk '{print \$2}')

if [ -z "\$PID" ]; then
    log_info "未找到运行中的进程"
    exit 0
fi

log_info "找到进程 PID: \$PID"
kill \$PID
log_info "✓ 进程已停止"

# 等待进程结束
for i in {1..10}; do
    if ! ps -p \$PID > /dev/null 2>&1; then
        log_info "✓ 进程已完全退出"
        exit 0
    fi
    sleep 1
done

log_warn "进程未在10秒内退出，强制终止..."
kill -9 \$PID
log_info "✓ 进程已强制终止"
EOF

chmod +x "$PACKAGE_DIR/scripts/stop.sh"
log_info "✓ 停止脚本已创建: scripts/stop.sh"

# 创建 README
cat > "$PACKAGE_DIR/README.txt" << EOF
Go-Study2 学习平台 v${VERSION}
============================

快速安装
--------
⚠️  重要：选择合适的安装目录
  推荐目录（生产环境）：
    - /opt/go-study2
    - /usr/local/go-study2
    - /home/youruser/go-study2

  不推荐目录：
    - /tmp/go-study2      (系统重启后数据丢失，systemd 兼容性问题)
    - /var/tmp/go-study2  (清理策略不确定)

方法1: systemd 服务（推荐）
  1. 解压到推荐目录:
     tar -xf go-study2-*.tar.gz -C /opt
     # 结果: /opt/go-study2/

  2. 运行安装脚本:
     cd /opt/go-study2
     sudo ./scripts/install.sh

  3. 启动服务:
     sudo systemctl start go-study2

  4. 查看状态:
     sudo systemctl status go-study2

方法2: 直接启动（测试用）
  1. 解压安装包:
     tar -xf go-study2-*.tar.gz -C /opt

  2. 启动服务:
     cd /opt/go-study2
     ./scripts/start.sh

目录结构
--------
go-study2/
├── bin/go-study2        # 主程序
├── configs/             # 配置文件
│   ├── config.yaml
│   └── logger.yaml
├── static/out/          # 前端静态文件
├── quiz_data/           # 题库数据
├── data/                # 数据库（运行时生成）
├── logs/                # 日志（运行时生成）
├── scripts/             # 管理脚本
│   ├── install.sh       # 安装（systemd）
│   ├── uninstall.sh     # 卸载
│   ├── start.sh         # 启动（前台）
│   ├── stop.sh          # 停止
│   └── logging.sh       # 日志函数库
├── VERSION              # 版本配置
└── README.txt           # 本文件

配置说明
--------
主配置文件: configs/config.yaml
日志配置: configs/logger.yaml

重要配置项:
  http.port: 8080              # 监听端口
  server.host: "0.0.0.0"       # 监听地址（0.0.0.0 允许外部访问）
  jwt.secret: "YOUR_SECRET"    # JWT 密钥（必须修改！）
  database.path: "./data/gostudy.db"
  static.path: "./static/out"

环境变量
--------
  JWT_SECRET              # JWT 签名密钥（≥32 字符，推荐设置）

服务管理（systemd）
--------
  启动: systemctl start go-study2
  停止: systemctl stop go-study2
  重启: systemctl restart go-study2
  开机自启: systemctl enable go-study2
  取消自启: systemctl disable go-study2
  查看状态: systemctl status go-study2
  查看日志: journalctl -u go-study2 -f

访问地址
--------
  HTTP: http://localhost:8080
  默认管理员: admin / GoStudy@123

⚠️  安全提醒
--------
1. 首次登录后立即修改默认密码
2. 设置强 JWT_SECRET（≥32 字符）
3. 生产环境配置 HTTPS
4. 配置防火墙规则
5. 定期备份数据库

升级指南
--------
1. 停止服务: sudo systemctl stop go-study2
2. 备份数据: cp data/gostudy.db data/gostudy.db.backup
3. 解压新版本覆盖文件
4. 启动服务: sudo systemctl start go-study2

卸载
----
sudo ./scripts/uninstall.sh

技术支持
--------
GitHub: https://github.com/yourusername/go-study2
文档: docs/PACKAGING.md
EOF

log_info "✓ README.txt 已创建"

# ========== 创建 tar.gz 包 ==========
log_step "打包成 tar.gz..."
cd "$BUILD_DIR"
tar -czf "$PACKAGE_FILE" "go-study2" 2>&1 | while IFS= read -r line; do log_info "  tar: $line"; done

if [ ! -f "$PACKAGE_FILE" ]; then
    log_error "打包失败"
    exit 1
fi

PACKAGE_SIZE=$(du -h "$PACKAGE_FILE" | cut -f1)
log_info "✓ 打包完成: $PACKAGE_FILE"
log_info "  文件大小: $PACKAGE_SIZE"

# ========== 生成 SHA256 校验和 ==========
log_step "生成 SHA256 校验和..."
cd "$BUILD_DIR"
sha256sum "$PACKAGE_FILE" > "${PACKAGE_FILE}.sha256"
log_info "✓ 校验和: ${PACKAGE_FILE}.sha256"
cat "${PACKAGE_FILE}.sha256" | while read -r hash file; do
    log_info "  SHA256: $hash"
done

# ========== 完成 ==========
echo ""
log_info "================================"
log_info "打包完成！"
log_info "================================"
echo ""
echo "安装包信息:"
echo "  文件名: ${PACKAGE_NAME}.tar.gz"
echo "  文件路径: ${PACKAGE_FILE}"
echo "  文件大小: ${PACKAGE_SIZE}"
echo "  版本: ${VERSION}"
echo "  构建时间: ${BUILD_TIME}"
echo "  服务名称: ${SERVICE_NAME}"
echo ""
echo "校验和文件:"
echo "  ${PACKAGE_NAME}.tar.gz.sha256"
echo ""
echo "部署命令:"
echo "  1. 上传: scp ${PACKAGE_FILE}.tar.gz user@server:/tmp/"
echo "  2. 解压: tar -xf ${PACKAGE_NAME}.tar.gz -C /opt/hadoop"
echo "  3. 安装: cd /opt/hadoop/go-study2 && sudo ./scripts/install.sh"
echo "  4. 启动: sudo systemctl start ${SERVICE_NAME}"
echo ""
echo "查看文档:"
echo "  快速参考: docs/QUICKREF.md"
echo "  完整文档: docs/PACKAGING.md"
echo ""
log_info "打包脚本执行完毕"
log_info "祝部署顺利！"
