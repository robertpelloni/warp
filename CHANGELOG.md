# Changelog

All notable changes to this project will be documented in this file.
Version: sourced from `VERSION.md`.

## [v0.2.5] - 2026-06-14

### Added
- **Virtual Screen Buffer (`ScreenBuffer`)**: Added 2D character grid emulation with `Cell` color and styling attributes in `pkg/terminal/terminal.go`.
- **Cursor Positioning & Scrolling**: Implemented `WriteRune`, `SetCursor`, `MoveCursor`, line-wrapping, and vertical scroll buffer shifts (`scrollUp`).
- **Screen Buffer Test Suite**: Added `TestScreenBufferOperations` in `pkg/terminal/terminal_test.go`.

## [v0.2.4] - 2026-06-14

### Added
- **Texture Atlas Glyph Caching (`TextureAtlas`)**: Implemented dynamic 2D glyph packing and UV coordinate calculation in `pkg/renderer/renderer.go`.
- **Fast Glyph Retrieval**: Thread-safe glyph query caching via `Get(rune)` and automatic row-wrapping allocation.
- **Renderer Test Suite**: Added test coverage in `pkg/renderer/renderer_test.go` verifying allocation, caching, and boundary overflow handling.

## [v0.2.3] - 2026-06-14

### Added
- **256-Color & 24-bit TrueColor ANSI Processing**: Extended `parseANSIColor` and `ParseANSI` in `pkg/terminal/terminal.go` to support extended SGR parameters (`38;5;n`, `48;5;n`, `38;2;r;g;b`, `48;2;r;g;b`).
- **ANSI Color Cube & Grayscale Calculation**: Added 6x6x6 color cube and grayscale ramp interpolation for xterm-256 color lookup.
- **Unit Test Coverage for Extended ANSI**: Added `TestExtendedANSIColors` in `pkg/terminal/terminal_test.go`.

## [v0.2.2] - 2026-06-14

### Added
- **Syntax Highlighting Lexer (`Tokenize`)**: Added multi-language lexical analysis in `pkg/editor/lexer.go` supporting Go, Rust, and Python.
- **Language Token Categorization**: Classifies tokens into keywords, strings, comments, numbers, operators, and identifiers with character offset tracking.
- **Lexer Test Suite**: Added comprehensive test coverage in `pkg/editor/lexer_test.go`.

## [v0.2.1] - 2026-06-14

### Added
- **Editor File Persistence (`SaveToFile`)**: Implemented disk writing and dirty-state tracking for active buffers in `pkg/editor`.
- **Global Save Shortcut (`Ctrl+S`)**: Intercepted `Ctrl+S` in both the Win32 `hEditor` pane and input edit controls to seamlessly save file changes directly to disk.
- **Unit Tests for Editor Disk Operations**: Added `TestSaveToFile` verifying persistence and content integrity.

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
