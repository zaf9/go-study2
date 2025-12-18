@echo off
setlocal
chcp 65001 >nul

REM ensure running from script dir
cd /d "%~dp0"

REM set Go module mode explicitly
set GO111MODULE=on

call :log [backend] start build pipeline...

call :log [backend] cleaning go cache and old artifacts...
go clean -cache -testcache
if exist "..\bin\gostudy.exe" del /f /q "..\bin\gostudy.exe"
if exist "coverage.out" del /f /q "coverage.out"

where go >nul 2>nul || (
  call :log [backend] Go not found. Please install Go and add PATH.
  goto :error
)

call :log [backend] formatting Go source...
gofmt -w .
if errorlevel 1 goto :error

call :log [backend] running go vet...
go vet ./...
if errorlevel 1 goto :error

call :log [backend] running tests with coverage...
go test -coverprofile=coverage.out ./...
if errorlevel 1 goto :error

call :log [backend] coverage report generated: %CD%\coverage.out
go tool cover -func=coverage.out

if not exist "..\bin" mkdir "..\bin"

call :log [backend] building binary...
go build -o "..\bin\gostudy.exe" main.go
if errorlevel 1 goto :error

call :log [backend] build success: ..\bin\gostudy.exe
goto :end

:error
call :log [backend] build failed. See errors above.
exit /b 1

:end
exit /b 0

:log
call :ts
echo [%TS%] %*
goto :eof

:ts
for /f "usebackq tokens=*" %%i in (`powershell -NoProfile -Command "(Get-Date).ToString('yyyy-MM-dd HH:mm:ss.fff')"`) do set TS=%%i
goto :eof
