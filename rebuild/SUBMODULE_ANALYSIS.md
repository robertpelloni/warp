# SUBMODULE ANALYSIS: Warp Rebuild Feature Tracking

This document tracks the analysis and architectural implementation status of features from external repositories being integrated into the Warp Go rebuild.

**Note on Status:** "Verified" indicates that the core architectural pattern and a foundational implementation have been ported and verified via integration tests. Full feature parity is an iterative process.

## Integration Progress (Architectural Parity)

| Tool | Status | Features Identified | Implemented in Go (Foundation) |
| :--- | :--- | :--- | :--- |
| **Tabby** | Verified | Tab management, Session abstractions | `pkg/terminal/tabs.go`, `pkg/terminal/session.go` |
| **Hyperharness** | Verified | LLM Provider Registry, Model Metadata | `pkg/harness/registry.go`, `pkg/harness/openai.go` |
| **Waveterm** | Verified | Persistent Terminal Blocks | `pkg/terminal/blocks.go`, `pkg/terminal/ansi.go` |
| **Pi-coding-agent**| Verified | Autonomous Agent Loop | `pkg/agent/agent.go` |
| **Antigravity 2.0** | Verified | Circuit Breaker Reliability | `pkg/agent/agent.go` (CircuitBreaker) |
| **Claude-mem** | Verified | Agent Service & Health Checks | `pkg/agent/service/service.go` |
| **Hermes Agent** | Verified | Memory Provider Orchestration | `pkg/agent/memory.go` |
| **OpenCode** | Verified | Supervisor Routing Pattern | `pkg/agent/supervisors/` |
| **Superpowers** | Verified | Skill Management System | `pkg/agent/skills/` |

## Next Phase: Deep Feature Porting

### Hyper (vercel/hyper fork)
- [x] Plugin architecture (Go-based plugin system).
- [ ] Cross-platform Electron-style UI patterns in TUI.

### Waveterm (Advanced)
- [ ] Multi-block persistent state.
- [ ] Custom UI widgets within blocks.

### Codex / Gemini / Hermes Desktop
- [x] Native Gemini integration patterns (WebAI-to-API).
- [x] Self-hosted workspace and capability management (Jules-autopilot).
- [ ] Specialized CLI tools for specific LLMs.
