#!/bin/bash

# Go-Study2 自签名证书生成脚本（域名版本）
# 生成包含域名的自签名证书，支持 go.study.org

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 默认值
CERT_DIR="backend/configs/certs"
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"
VALID_DAYS=365

# 主域名（固定）
DOMAIN="go.study.org"

# 额外的 DNS 名称
EXTRA_DNS="localhost,*.local"

# 额外的 IP 地址（逗号分隔）
EXTRA_IP="172.28.186.65"

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -d, --domain DOMAIN    主域名 (默认: go.study.org)"
    echo "  -e, --extra-dns LIST   额外的 DNS 名称，逗号分隔 (默认: localhost,*.local)"
    echo "  -i, --extra-ip LIST    额外的 IP 地址，逗号分隔 (默认: 172.28.186.65)"
    echo "  -o, --output DIR       证书输出目录 (默认: backend/configs/certs)"
    echo "  -v, --valid-days DAYS  证书有效期天数 (默认: 365)"
    echo "  -f, --force            强制覆盖现有证书"
    echo "  -h, --help             显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0                                    # 使用默认域名生成证书"
    echo "  $0 -f                                 # 强制覆盖现有证书"
    echo "  $0 -d myapp.local -e api.myapp.local  # 自定义域名"
    echo ""
    echo -e "${BLUE}域名方案优势：${NC}"
    echo "  1. 一次生成，永久有效（无需 IP 变化时重新生成）"
    echo "  2. IP 变化时只需修改客户端 /etc/hosts 文件"
    echo "  3. 无需重启服务"
    echo ""
    echo -e "${BLUE}客户端配置：${NC}"
    echo "  在 /etc/hosts 或 C:\\Windows\\System32\\drivers\\etc\\hosts 添加："
    echo "  172.28.186.65 $DOMAIN"
    echo ""
}

# 解析命令行参数
FORCE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--domain)
            DOMAIN="$2"
            shift 2
            ;;
        -e|--extra-dns)
            EXTRA_DNS="$2"
            shift 2
            ;;
        -i|--extra-ip)
            EXTRA_IP="$2"
            shift 2
            ;;
        -o|--output)
            CERT_DIR="$2"
            CERT_FILE="$CERT_DIR/server.crt"
            KEY_FILE="$CERT_DIR/server.key"
            shift 2
            ;;
        -v|--valid-days)
            VALID_DAYS="$2"
            shift 2
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 检查现有证书
if [ -f "$CERT_FILE" ] && [ "$FORCE" = false ]; then
    echo -e "${YELLOW}证书文件已存在: $CERT_FILE${NC}"
    echo "使用 -f 或 --force 选项覆盖现有证书"
    echo ""
    echo "当前证书信息："
    openssl x509 -in "$CERT_FILE" -text -noout | grep -A 5 "Subject:"
    openssl x509 -in "$CERT_FILE" -text -noout | grep -A 10 "Subject Alternative Name"
    exit 0
fi

# 创建证书目录
mkdir -p "$CERT_DIR"

# 构建配置文件
OPENSSL_CONFIG=$(mktemp)
cat > "$OPENSSL_CONFIG" << EOF
[req]
default_bits = 4096
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = CN
ST = Beijing
L = Beijing
O = Go-Study2
OU = Development
CN = $DOMAIN

[v3_req]
keyUsage = digitalSignature, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = $DOMAIN
EOF

# 添加额外的 DNS 名称
DNS_INDEX=2
echo "$EXTRA_DNS" | tr ',' '\n' | while read -r dns; do
    if [ -n "$dns" ]; then
        echo "DNS.$DNS_INDEX = $dns" >> "$OPENSSL_CONFIG"
        DNS_INDEX=$((DNS_INDEX + 1))
    fi
done

# 添加 IP 地址
IP_INDEX=1
echo "$EXTRA_IP" | tr ',' '\n' | while read -r ip; do
    if [ -n "$ip" ]; then
        echo "IP.$IP_INDEX = $ip" >> "$OPENSSL_CONFIG"
        IP_INDEX=$((IP_INDEX + 1))
    fi
done

# 生成私钥和证书
echo -e "${GREEN}正在生成自签名证书...${NC}"
echo "主域名: $DOMAIN"
echo "额外 DNS: $EXTRA_DNS"
echo "额外 IP: $EXTRA_IP"
echo "有效期: $VALID_DAYS 天"
echo "输出目录: $CERT_DIR"
echo ""

openssl req -x509 -newkey rsa:4096 -keyout "$KEY_FILE" -out "$CERT_FILE" \
    -days "$VALID_DAYS" -nodes \
    -config "$OPENSSL_CONFIG" \
    -extensions v3_req 2>&1 | grep -v "^+"

# 清理临时文件
rm -f "$OPENSSL_CONFIG"

# 设置权限
chmod 600 "$KEY_FILE"
chmod 644 "$CERT_FILE"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}证书生成完成！${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "证书文件: $CERT_FILE"
echo "私钥文件: $KEY_FILE"
echo ""
echo "证书详情："
openssl x509 -in "$CERT_FILE" -text -noout | grep -A 10 "Subject Alternative Name"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}客户端配置说明${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}1. 在客户端机器的 hosts 文件中添加以下内容：${NC}"
echo ""
echo -e "   ${GREEN}Linux/Mac: /etc/hosts${NC}"
echo -e "   ${GREEN}Windows:   C:\\Windows\\System32\\drivers\\etc\\hosts${NC}"
echo ""
echo -e "   ${BLUE}# Go-Study2 访问域名${NC}"
echo -e "   ${BLUE}172.28.186.65  $DOMAIN${NC}"
echo ""
echo -e "${YELLOW}2. 访问地址：${NC}"
echo -e "   ${GREEN}https://$DOMAIN:8443/login${NC}"
echo ""
echo -e "${YELLOW}3. 如果 IP 地址变化：${NC}"
echo "   只需修改客户端 hosts 文件中的 IP 地址，无需重新生成证书！"
echo ""
echo -e "${YELLOW}4. 浏览器安全警告：${NC}"
echo "   首次访问会显示安全警告（因为是自签名证书）"
echo "   点击\"高级\"→\"继续访问\"即可"
echo "   可以选择将证书添加到浏览器的信任存储以避免重复提示"
echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
