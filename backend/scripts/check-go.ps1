# 运行 gofmt、go vet、golint 的统一检查脚本
# 用法：
#   powershell -ExecutionPolicy Bypass -File scripts/check-go.ps1
# 说明：
#   - 自动切换到后端目录 (backend)
#   - 仅在目标目录存在时执行 gofmt
#   - golint 不存在时给出提示但不中断

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

$targets = @("internal", "src", "tests", "main.go")
$existingTargets = @()
foreach ($path in $targets) {
    if (Test-Path $path) {
        $resolved = Resolve-Path $path
        $existingTargets += $resolved
    }
}

if ($existingTargets.Count -gt 0) {
    Write-Host "run gofmt ..." -ForegroundColor Cyan
    gofmt -w $existingTargets
} else {
    Write-Host "skip gofmt: no target dirs" -ForegroundColor Yellow
}

Write-Host "run go vet ..." -ForegroundColor Cyan
go vet ./...

$golint = Get-Command golint -ErrorAction SilentlyContinue
if ($golint) {
    Write-Host "run golint ..." -ForegroundColor Cyan
    golint ./...
} else {
    Write-Warning "golint not found, run 'go install golang.org/x/lint/golint@latest' to enable."
}

Write-Host "check done" -ForegroundColor Green

