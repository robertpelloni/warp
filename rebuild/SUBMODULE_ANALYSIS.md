# SUBMODULE ANALYSIS: Warp Rebuild Feature Tracking

This document tracks the analysis and implementation status of features from external repositories being integrated into the Warp Go rebuild.

## Integration Progress

| Tool | Status | Features Identified | Implemented in Go |
| :--- | :--- | :--- | :--- |
| **Tabby** | Verified | Tab management, Session abstractions | `pkg/terminal/tabs.go`, `pkg/terminal/session.go` |
| **Hyperharness** | Verified | LLM Provider Registry, Model Metadata | `pkg/harness/registry.go`, `pkg/harness/openai.go` |
| **Waveterm** | Verified | Persistent Terminal Blocks | `pkg/terminal/blocks.go`, `pkg/simulation/engine.go` |
| **Pi-coding-agent**| Verified | Autonomous Agent Loop | `pkg/agent/agent.go` |
| **Antigravity 2.0** | Verified | Circuit Breaker Reliability | `pkg/agent/circuit_breaker.go` |
| **Claude-mem** | Verified | Agent Service & Health Checks | `pkg/agent/service/service.go` |
| **Hermes Agent** | Verified | Memory Provider Orchestration | `pkg/agent/memory.go` |
| **OpenCode** | Verified | Supervisor Routing Pattern | `pkg/agent/supervisors/` |
| **Superpowers** | Verified | Skill Management System | `pkg/agent/skills/` |

## Next Phase: Deep Feature Parity

### Hyper (vercel/hyper fork)
- [x] Plugin architecture (Go-based plugin system).
- [ ] Cross-platform Electron-style UI patterns in TUI.

### Waveterm (Advanced)
- [ ] Multi-block persistent state.
- [ ] Custom UI widgets within blocks.

### Codex / Gemini / Hermes Desktop
- [x] Native Gemini integration patterns (WebAI-to-API).
- [x] Self-hosted workspace and capability management (Jules-autopilot).
- [x] Multi-agent orchestration and workflow runner (Maestro).
- [x] Smart routing and provider fallbacks (OmniRoute).
- [ ] Specialized CLI tools for specific LLMs.
