# TODO

## Phase 1: Foundation (Complete)
- [x] Initialize `rebuild/` directory and PTY handling.
- [x] Implement autonomous agent loop with tool-use.
- [x] Integrate Go sidecar with Rust core via IPC Bridge.
- [x] Verify 100% test pass rate and staging readiness.
- [x] Finalize Milestone v3.1.0: Audit, Handover, and Staging complete.

## Phase 2: Parity & Expansion (Complete)
- [x] Implement VT100/Xterm terminal parsing in Go (Foundational grid/cursor).
- [x] Build out the TUI using `bubbletea` (Experimental/Wired).
- [x] Integrate actual LLM calls (Anthropic/OpenAI) with tool-calling support.
- [x] Integrate OpenCode Supervisor pattern.
- [x] Integrate Claude Code Skill Manager.
- [x] Integrate Model Context Protocol (MCP) foundation.
- [x] Implement multi-pane layout management.
