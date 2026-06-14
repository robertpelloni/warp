@echo off
setlocal enabledelayedexpansion
title Warp Go

cd /d "%~dp0"
set "BINARY=warp-go.exe"
set "ENTRY=./cmd/warp"

set "CMD=%1"
if "%CMD%"=="" set "CMD=run"

if /i "%CMD%"=="run" goto :run
if /i "%CMD%"=="build" goto :build
if /i "%CMD%"=="test" goto :test
if /i "%CMD%"=="clean" goto :clean
if /i "%CMD%"=="help" goto :help
echo Unknown command: %CMD%
goto :help

:build
echo.
echo [Warp Go] Building...
go mod download
if errorlevel 1 ( echo [FAIL] Dependency download & exit /b 1 )
go build -buildvcs=false -ldflags="-s -w -H=windowsgui" -o %BINARY% %ENTRY%
if errorlevel 1 ( echo [FAIL] Build failed & exit /b 1 )
for %%f in (%BINARY%) do echo [OK] %%~zf bytes
goto :end

:run
if not exist %BINARY% call :build
if errorlevel 1 exit /b 1
echo.
echo [Warp Go] Starting agentic development environment...
echo.
%BINARY%
goto :end

:test
echo [Warp Go] Running tests...
go test ./pkg/... ./cmd/... -v -count=1 -timeout 120s
goto :end

:clean
del /q %BINARY% 2>nul
go clean
echo [Warp Go] Cleaned.
goto :end

:help
echo.
echo Warp Go - Usage: start.bat [command]
echo.
echo Commands:
echo   run     Build and run (default)
echo   build   Build binary only
echo   test    Run tests
echo   clean   Remove binary
echo   help    Show this help
echo.
echo Packages:
echo   pkg/pty       ConPTY management (Windows)
echo   pkg/terminal  Terminal emulation and ANSI parsing
echo   pkg/command   Command engine, aliases, AI integration
echo   pkg/session   Terminal session/tab manager
echo   pkg/editor    IDE editing capabilities
echo   pkg/theme     7 built-in color themes
echo   pkg/win32     Pure-Go Win32 API bindings
echo   pkg/app       Main GUI application
echo.
echo Keyboard Shortcuts:
echo   Enter          Execute command
echo   Ctrl+L         Clear output
echo   Ctrl+K         Command palette
echo   Ctrl+T         New session
echo   Ctrl+Shift+T   Cycle theme
echo   Ctrl+Up/Down   Command history
echo.
goto :end

:end
endlocal
