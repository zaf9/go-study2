#!/bin/bash

set -e

# Change to script directory
cd "$(dirname "$0")"

# Function to log messages with timestamp
log() {
    local ts
    ts=$(date +"%Y-%m-%d %H:%M:%S.%3N" 2>/dev/null || date +"%Y-%m-%d %H:%M:%S")
    echo "[$ts] $*"
}

# Function to handle errors
error_handler() {
    log "[frontend] build failed. See errors above."
    exit 1
}

trap error_handler ERR

log "[frontend] start build pipeline..."

log "[frontend] cleaning previous build outputs..."
rm -rf .next coverage

# Check if Node.js is available
if ! command -v node >/dev/null 2>&1; then
    log "Node.js 18+ not found. Please install and add PATH."
    exit 1
fi

# Check if npm is available
if ! command -v npm >/dev/null 2>&1; then
    log "npm not found. Please verify Node.js installation."
    exit 1
fi

log "[frontend] installing dependencies..."
npm install

log "[frontend] running lint..."
npm run lint

log "[frontend] running tests with coverage..."
npm run test -- --coverage

log "[frontend] building production bundle..."
npm run build

log "[frontend] exporting static files..."
rm -rf out
mkdir -p out

# Copy build output (equivalent to robocopy /E)
cp -r .next/server/app/* out/
cp -r .next/static out/_next

log "[frontend] build success. Output: out"
exit 0
