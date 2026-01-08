#!/bin/bash

set -e

# Change to script directory
cd "$(dirname "$0")"

# Set Go module mode explicitly
export GO111MODULE=on

# Function to log messages with timestamp
log() {
    local ts
    ts=$(date +"%Y-%m-%d %H:%M:%S.%3N" 2>/dev/null || date +"%Y-%m-%d %H:%M:%S")
    echo "[$ts] $*"
}

# Function to handle errors
error_handler() {
    log "[backend] build failed. See errors above."
    exit 1
}

trap error_handler ERR

log "[backend] start build pipeline..."

log "[backend] cleaning go cache and old artifacts..."
go clean -cache -testcache
rm -f "../bin/gostudy"
rm -f "../bin/gostudy.exe"
rm -f "coverage.out"

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    log "[backend] Go not found. Please install Go and add PATH."
    exit 1
fi

log "[backend] formatting Go source..."
gofmt -w .

log "[backend] running go vet..."
go vet ./...

log "[backend] running tests with coverage..."
go test -coverprofile=coverage.out ./...

log "[backend] coverage report generated: $(pwd)/coverage.out"
go tool cover -func=coverage.out

# Create bin directory if it doesn't exist
mkdir -p "../bin"

# Determine binary name based on OS
BINARY_NAME="../bin/gostudy"
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    BINARY_NAME="../bin/gostudy.exe"
fi

log "[backend] building binary..."
go build -o "$BINARY_NAME" main.go

log "[backend] build success: $BINARY_NAME"
exit 0
