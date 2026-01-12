.PHONY: help build clean package test-all deploy-local info version

# 默认目标
help:
	@echo "Go-Study2 Makefile 命令:"
	@echo ""
	@echo "  make build         - 构建前后端（本地开发）"
	@echo "  make package       - 打包成 tar.gz（Linux 部署包）"
	@echo "  make test-all      - 运行所有测试"
	@echo "  make clean         - 清理构建产物"
	@echo "  make deploy-local  - 本地模拟部署"
	@echo "  make info          - 显示构建信息"
	@echo "  make version       - 显示版本信息"
	@echo ""
	@echo "示例:"
	@echo "  make package VERSION=1.0.0"
	@echo "  make deploy-local VERSION=1.0.0"

# 构建版本（优先从 VERSION 文件读取）
VERSION?=$(shell if [ -f VERSION ]; then \
		source VERSION && echo $$VERSION; \
	else \
		git describe --tags --always --dirty 2>/dev/null || echo "1.0.0"; \
	fi)
BUILD_TIME=$(shell date +"%Y%m%d_%H%M%S")
SERVICE_NAME=go-study2

# 构建前后端（本地开发）
build: build-frontend build-backend
	@echo "✓ 构建完成"

build-frontend:
	@echo "[前端] 开始构建..."
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: 清理旧构建产物"
	cd frontend && rm -rf .next out
	cd frontend && npm install
	cd frontend && npm run build
	cd frontend && npm run export
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: ✓ 前端构建完成: frontend/out/"

build-backend:
	@echo "[后端] 开始构建..."
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: 格式化代码"
	cd backend && gofmt -w .
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: 代码检查"
	cd backend && go vet ./...
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: 运行测试"
	cd backend && go test -cover ./...
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: 编译二进制（本地开发，启用 CGO）"
	cd backend && go build -o ../bin/go-study2 main.go
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: ✓ 后端构建完成: bin/go-study2"

# 运行所有测试
test-all: test-frontend test-backend
	@echo "✓ 所有测试通过"

test-frontend:
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: [前端] 运行测试..."
	cd frontend && npm test -- --coverage

test-backend:
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: [后端] 运行测试..."
	cd backend && go test -cover ./...

# 打包成 tar.gz
package:
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: [打包] 开始打包 v$(VERSION)..."
	@VERSION=$(VERSION) BUILD_TIME=$(BUILD_TIME) ./scripts/package.sh
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: ✓ 打包完成: build/"

# 清理构建产物
clean:
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: [清理] 清理构建产物..."
	rm -rf build/
	rm -rf frontend/.next
	rm -rf frontend/out
	rm -rf bin/go-study2 bin/go-study2.exe
	cd backend && rm -f coverage.out
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: ✓ 清理完成"

# 本地模拟部署（用于测试）
deploy-local: package
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: [本地部署] 模拟部署流程..."
	@rm -rf /tmp/go-study2-test
	@mkdir -p /tmp/go-study2-test
	@tar -xzf build/go-study2-*.tar.gz -C /tmp/go-study2-test
	@echo ""$(date +%Y%m%d-%H%M%S)" - [Makefile]: ✓ 解压完成: /tmp/go-study2-test/"
	@echo ""
	@echo "目录结构:"
	@ls -la /tmp/go-study2-test/
	@echo ""
	@echo "项目文件:"
	@ls -la /tmp/go-study2-test/go-study2/
	@echo ""
	@echo "手动安装测试:"
	@echo "  cd /tmp/go-study2-test/go-study2 && sudo ./scripts/install.sh"
	@echo ""
	@echo "直接启动测试:"
	@echo "  cd /tmp/go-study2-test/go-study2 && ./scripts/start.sh"

# 显示构建信息
info:
	@echo "构建信息:"
	@echo "  项目名称: $(SERVICE_NAME)"
	@echo "  版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo "  Go 版本: $$(go version 2>/dev/null | awk '{print $$3}')"
	@echo "  Node 版本: $$(node --version 2>/dev/null)"
	@echo "  npm 版本: $$(npm --version 2>/dev/null)"
	@echo "  系统: $$(uname -s) $$(uname -m)"
	@echo "  工作目录: $$(pwd)"
	@echo "  构建目录: $$(pwd)/build"
	@echo "  二进制: bin/go-study2"
	@echo "  服务名: $(SERVICE_NAME)"

# 显示版本信息
version:
	@if [ -f VERSION ]; then \
		echo "版本信息:"; \
		cat VERSION; \
	else \
		echo "VERSION 文件不存在"; \
		echo "当前版本: $(VERSION)"; \
	fi
