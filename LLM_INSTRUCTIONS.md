# LLM Instructions & Repository Principles

## Architecture & Code Quality Guidelines
1. **Unified Dashboard First**:
   - Always condense features onto the singular main dashboard page. Avoid creating fragmented popups or disconnected subpages.
   - Group related tools together, prioritize high-value workflows (terminal, editor, blocks), and provide clear status and tooltips.
2. **Cross-Platform Compatibility**:
   - Keep platform-specific code gated (`//go:build windows` vs `//go:build !windows`).
   - Any API modifications to Windows abstractions must have compatible stubs or interface implementations for Linux/macOS.
   - Never introduce CGO dependencies unless strictly necessary and guarded.
3. **Documentation Governance**:
   - Sourced versioning from `VERSION.md`.
   - Update `CHANGELOG.md`, `TODO.md`, and `ROADMAP.md` as features evolve.
   - Before ending a session, record all structural shifts and discoveries in `HANDOFF.md`.
4. **Verification & Testing**:
   - Always verify Go changes with both `go test -v ./...` and `GOOS=windows GOARCH=amd64 go build ./...`.
   - Run `./script/presubmit` before committing.
