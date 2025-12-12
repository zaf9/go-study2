param(
    [string]$ServerExe = "..\bin\gostudy.exe",
    [int]$Concurrency = 100,
    [int]$Requests = 1000,
    [int]$StartupWaitSec = 3
)

Set-StrictMode -Version Latest

function Write-Log {
    param([string]$m) 
    $ts = Get-Date -Format o
    Write-Host "[$ts] $m"
}

try {
    $scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
    Push-Location $scriptDir\..

    Write-Log "Starting quickstart validation in: $(Get-Location)"

    $serverProc = $null

    if (Test-Path $ServerExe) {
        Write-Log "Found server executable: $ServerExe. Starting..."
        $serverProc = Start-Process -FilePath $ServerExe -ArgumentList "-d" -PassThru -NoNewWindow
    } else {
        Write-Log "Server executable not found, attempting 'go run main.go -d'..."
        $serverProc = Start-Process -FilePath "go" -ArgumentList "run main.go -d" -WorkingDirectory "./" -PassThru -NoNewWindow
    }

    Start-Sleep -Seconds $StartupWaitSec

    Write-Log "Running stress client (concurrency=$Concurrency, requests=$Requests)"
    $env:GO111MODULE = 'on'
    $args = @(
        'run',
        'scripts/stress_client.go',
        '--url','http://localhost:8080/',
        '--concurrency',$Concurrency.ToString(),
        '--requests',$Requests.ToString()
    )
    Write-Log "Executing: go $($args -join ' ')"
    & go @args

    # collect logs
    $ts = Get-Date -Format "yyyyMMdd-HHmmss"
    $outDir = "logs\validation-$ts"
    Write-Log "Collecting logs to: $outDir"
    New-Item -ItemType Directory -Force -Path $outDir | Out-Null
    Get-ChildItem -Path logs -Recurse -File | ForEach-Object {
        $rel = $_.FullName.Substring((Get-Location).Path.Length).TrimStart('\')
        $dest = Join-Path $outDir $_.Name
        Copy-Item -Path $_.FullName -Destination $dest -Force
    }

    Write-Log "Logs copied. Listing archive:"
    Get-ChildItem $outDir | Select-Object Name,Length | Format-Table

} catch {
    Write-Error "Validation failed: $_"
} finally {
    # stop server if started
    if ($serverProc -ne $null) {
        try {
            Write-Log "Stopping server (Id=$($serverProc.Id))"
            Stop-Process -Id $serverProc.Id -Force -ErrorAction SilentlyContinue
        } catch {
            Write-Warning "Failed to stop server process: $_"
        }
    }
    Pop-Location
}
