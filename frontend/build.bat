@echo off
setlocal
chcp 65001 >nul

cd /d "%~dp0"

call :log [frontend] start build pipeline...

call :log [frontend] cleaning previous build outputs...
if exist ".next" rmdir /s /q ".next"
if exist "coverage" rmdir /s /q "coverage"

where node >nul 2>nul || (
  call :log Node.js 18+ not found. Please install and add PATH.
  goto :error
)

where npm >nul 2>nul || (
  call :log npm not found. Please verify Node.js installation.
  goto :error
)

call :log [frontend] installing dependencies...
call npm install
if errorlevel 1 goto :error

call :log [frontend] running lint...
call npm run lint
if errorlevel 1 goto :error

call :log [frontend] running tests with coverage...
call npm run test -- --coverage
if errorlevel 1 goto :error

call :log [frontend] building production bundle...
call npm run build
if errorlevel 1 goto :error

call :log [frontend] build success. Output: out
goto :end

:error
call :log [frontend] build failed. See errors above.
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
