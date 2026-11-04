@echo off
setlocal enabledelayedexpansion
title Warp Go - Agentic Development Environment

cd /d "%~dp0"
set "BINARY=warp-go.exe"
set "ENTRY=./cmd/warp"

set "CMD=%1"
if "%CMD%"=="" set "CMD=run"

if /i "%CMD%"=="run" goto :run
if /i "%CMD%"=="build" goto :build
if /i "%CMD%"=="test" goto :test
if /i "%CMD%"=="presubmit" goto :presubmit
if /i "%CMD%"=="clean" goto :clean
if /i "%CMD%"=="help" goto :help
echo Unknown command: %CMD%
goto :help

:build
echo.
echo [Warp Go] Building unified dashboard binary...
go mod download
if errorlevel 1 ( echo [FAIL] Dependency download & exit /b 1 )
go build -buildvcs=false -ldflags="-s -w -H=windowsgui" -o %BINARY% %ENTRY%
if errorlevel 1 ( echo [FAIL] Build failed & exit /b 1 )
for %%f in (%BINARY%) do echo [OK] %%~zf bytes created (%BINARY%)
goto :end

:run
if not exist %BINARY% call :build
if errorlevel 1 exit /b 1
echo.
echo [Warp Go] Starting unified agentic development environment...
echo.
%BINARY%
goto :end

:test
echo [Warp Go] Running tests...
go test ./pkg/... ./cmd/... -v -count=1 -timeout 120s
goto :end

:presubmit
echo [Warp Go] Running full presubmit checks...
gofmt -w ./pkg ./cmd
go test ./... -v -count=1
go build -buildvcs=false ./...
echo [OK] Presubmit verification complete.
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
echo   run        Build and launch unified dashboard (default)
echo   build      Build Windows GUI binary only
echo   test       Run test suite
echo   presubmit  Format code and run all verification checks
echo   clean      Remove binary and build artifacts
echo   help       Show this help
echo.
echo Key Components:
echo   pkg/pty       ConPTY management and POSIX fallback
echo   pkg/terminal  Terminal emulation and ANSI parsing
echo   pkg/command   Command engine, /edit, aliases, AI queries
echo   pkg/session   Terminal session/tab manager
echo   pkg/editor    IDE editing engine
echo   pkg/theme     7 built-in themes
echo   pkg/win32     Win32 API & System Tray bindings
echo   pkg/app       Unified Dashboard GUI & System Tray lifecycle
echo.
goto :end

:end
endlocal
