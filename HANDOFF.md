# Session Handoff Notes

## Key Findings & Changes in This Turn
1. **Repository Layout**:
   - The Go port is organized under `pkg/` (`app`, `command`, `editor`, `pty`, `renderer`, `session`, `terminal`, `theme`, `win32`) and `cmd/warp`.
   - The original Rust repository remains in `crates/` and `app/`.
2. **Dashboard UI Refactor**:
   - The dashboard in `pkg/app/app_windows.go` now unifies the terminal output, the editor buffer pane (`hEditor`), and the command blocks sidebar into one clean screen.
   - A system tray icon is automatically registered on startup via `Shell_NotifyIconW`. Right-clicking the tray icon brings up options to restore the dashboard, minimize to tray, or terminate the app and background processes cleanly.
3. **Cross-Platform Compilation**:
   - All tests run cleanly on Linux via `go test -v ./...`.
   - Windows cross-compilation is fully verified via `GOOS=windows GOARCH=amd64 go build ./...`.

## Resumption Instructions for Next Agent
- Run `./script/presubmit` to verify build and test health.
- All documentation files are tracked at root (`VERSION.md`, `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `CHANGELOG.md`, `ROADMAP.md`, `TODO.md`, `HANDOFF.md`, `LLM_INSTRUCTIONS.md`).
- For any new features, keep the unified single-page dashboard principle intact.

## Release v0.2.1 Updates:
- Integrated `SaveToFile` in `pkg/editor` and wired `Ctrl+S` in `pkg/app/app_windows.go` for instant file persistence from the dashboard editor pane.
- Presubmit checks and all unit tests verified.
