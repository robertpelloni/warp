# Granular Tasks & Short-Term TODO

- [x] Abstract `platformPTY` interface in `pkg/pty/pty.go`.
- [x] Create non-Windows app stub `pkg/app/app_unix.go`.
- [x] Add system tray API bindings in `pkg/win32/win32.go`.
- [x] Add system tray icon and context menu in `pkg/app/app_windows.go`.
- [x] Wire IDE Editor pane (`hEditor`) into unified dashboard layout.
- [x] Implement `/edit <file>` command in `pkg/command/engine.go`.
- [x] Write unit tests for `pkg/pty` and `/edit` command execution.
- [x] Create governance documentation (`VERSION.md`, `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `CHANGELOG.md`, `ROADMAP.md`, `TODO.md`, `HANDOFF.md`, `LLM_INSTRUCTIONS.md`).
- [ ] Enhance ANSI escape sequence parsing for complex cursor movements.
- [ ] Add syntax highlighting lexer for Go, Rust, and Python in the IDE editor pane.
- [x] Implement file save listener (`Ctrl+S`) in the editor pane.
