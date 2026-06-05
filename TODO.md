# TODO

## Phase 1: Foundation (Complete)
- [x] Initialize `rebuild/` directory and PTY handling.
- [x] Implement autonomous agent loop with tool-use.
- [x] Integrate Go sidecar with Rust core via IPC Bridge.
- [x] Verify 100% test pass rate and staging readiness.

## Phase 2: Parity & Expansion (Current)
- [ ] Implement VT100/Xterm terminal parsing in Go.
- [ ] Build out the TUI using `bubbletea`.
- [ ] Integrate actual LLM calls (Anthropic/OpenAI) into the Agent Loop.
