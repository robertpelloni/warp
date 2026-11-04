# Changelog

All notable changes to this project will be documented in this file.
Version: sourced from `VERSION.md`.

## [v0.2.0] - 2026-06-14

### Added
- **Unified Dashboard Architecture**: Re-architected `pkg/app/app_windows.go` into a singular dashboard consolidating the Terminal Output, IDE Editor pane (`hEditor`), Command Blocks sidebar, Prompt label, and Status Bar into one cohesive window.
- **System Tray Lifecycle & Server Management**: Full system tray integration via `Shell_NotifyIconW`, providing quick restoration, minimizing to tray, and a server termination shutdown option (`Quit Warp Servers & App`).
- **PTY Platform Abstraction**: Introduced `platformPTY` interface in `pkg/pty` cleanly separating ConPTY Windows implementations from POSIX Unix fallback pipes.
- **Cross-Platform App Stub**: Added `pkg/app/app_unix.go` enabling headless execution and clean test execution in Linux/Unix environments.
- **Integrated IDE Editor Command**: Added `/edit <file>` and `/open <file>` command support in `pkg/command/engine.go` and wired it to load files directly into the dashboard editor buffer.
- **Automated Presubmit Suite**: Added `./script/presubmit` cross-platform validation script running formatting, tests, and Windows cross-compilation checks.
- **Global Documentation Standards**: Established comprehensive governance documentation: `VERSION.md`, `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `ROADMAP.md`, `TODO.md`, `HANDOFF.md`, and `LLM_INSTRUCTIONS.md`.

### Fixed
- Fixed build failures in `pkg/pty` when building on non-Windows platforms.
- Fixed `winPTY.Close()` method signature to return an `error` matching `platformPTY`.
- Fixed missing Win32 definitions and menu constants in `pkg/win32/win32.go`.
