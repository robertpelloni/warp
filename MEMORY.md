# Architectural Memory & Observations

## Codebase Traits
- **Go-First Architecture**: The module `github.com/robertpelloni/warp` targets Go 1.25.0+ and prioritizes native execution.
- **Platform Separation**: GUI and OS-level operations are strictly separated using Go build constraints (e.g., `//go:build windows` vs `//go:build !windows`).
- **Win32 Dominance (Current)**: The primary frontend implementation (`pkg/app/app_windows.go`, `pkg/win32`) relies heavily on direct Win32 API calls via `golang.org/x/sys/windows`. This provides excellent performance on Windows but necessitates parallel implementations for other OSs.
- **UNIX Stubbing**: Non-Windows platforms currently rely on stub files (like `pkg/app/app_unix.go`) to ensure the module compiles universally, even if the GUI is non-functional on those platforms yet.
- **Component Abstraction**: Core logic components like `pkg/command` (engine), `pkg/editor` (buffer management), and `pkg/session` (tab management) are completely decoupled from the rendering layer, making future UI ports significantly easier.

## Design Preferences
- **Security**: Strict adherence to secure practices. No hardcoded secrets. Use constant-time comparisons.
- **Commenting**: In-depth comments explaining the "why", structural side effects, and alternate methods attempted. Trivial comments are discouraged.
- **Error Handling**: Bubble up errors gracefully rather than panicking, especially in the PTY and command execution pipelines to prevent full IDE crashes on single command failures.
