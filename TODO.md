# Warp Go TODO

## Immediate Tasks
- [x] Connect the `pkg/command` engine to the `pkg/app` GUI input and execution pipeline.
- [x] Render the output buffer from `pkg/terminal` inside the Win32 main window layout.
- [x] Implement syntax highlighting and basic editing features in `pkg/editor`.
- [ ] Flesh out UNIX/Linux stub (`pkg/app/app_unix.go`) with an actual GUI framework (e.g. Fyne, GTK, or Qt bindings) to achieve true cross-platform parity.
- [x] Create missing global documentation files (`VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `CHANGELOG.md`).
- [x] Integrate submodule status checks into a status output command.

## Next Steps
- [x] Finish UI representations of the `pkg/session` active states (e.g., drawing tabs for the tab manager).
- [x] Implement Git Diff Monitoring for background Shadow Pilot anomaly detection.
- [x] Add CI Pipeline Auto-Fix support in the AI Agent command logic.
- [x] Write integration tests for the full terminal lifecycle.
