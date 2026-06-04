# TODO

## Immediate Tasks
- [x] Initialize `rebuild/` directory with `go mod init`.
- [x] Implement PTY handling with `creack/pty`.
- [x] Integrate agent loop, circuit breaker, and memory manager.
- [x] Run comprehensive test suite and verify 100% pass rate.
- [x] Verify staging deployment integrity.

## Next Phase
- [ ] Implement VT100/Xterm terminal parsing in Go.
- [ ] Build out the TUI using `bubbletea`.
- [ ] Integrate actual LLM calls (Anthropic/OpenAI) into the Agent Loop.
