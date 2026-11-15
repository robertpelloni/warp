# SUBMODULE ANALYSIS: Warp Rebuild Feature Tracking

This document tracks the analysis and foundational implementation status of features from external repositories being integrated into the Warp Go rebuild.

**Note on Status:** "Verified" indicates that the core architectural pattern and a foundational implementation have been ported and verified via integration tests. Achievement of full feature parity is an ongoing, iterative process across subsequent development phases.

## Integration Progress (Foundational Parity)

| Tool | Status | Features Identified | Implemented in Go (Foundation) |
| :--- | :--- | :--- | :--- |
| **Tabby** | Verified | Tab management, Codebase indexing | `pkg/terminal/tabs.go`, `pkg/agent/tools/search_codebase.go` |
| **Hyperharness** | Verified | LLM Provider Registry, Model Metadata | `pkg/harness/registry.go`, `pkg/harness/openai.go` |
| **Waveterm** | Verified | Persistent Blocks, Agent TUI | `pkg/terminal/blocks.go`, `cmd/warp/tui.go` |
| **Pi-coding-agent**| Verified | Autonomous Agent Loop | `pkg/agent/agent.go` |
| **Antigravity 2.0** | Verified | Circuit Breaker Reliability | `pkg/agent/agent.go` (CircuitBreaker) |
| **Claude-mem** | Verified | Agent Service & Health Checks | `pkg/agent/service/service.go` |
| **Hermes Agent** | Verified | Memory Provider Orchestration | `pkg/agent/memory.go` |
| **OpenCode** | Verified | Supervisor Routing Pattern | `pkg/agent/supervisors/` |
| **Superpowers** | Verified | Skill Management System | `pkg/agent/skills/` |
| **MCP** | Verified | Model Context Protocol support | `pkg/mcp/`, `pkg/agent/tools/mcp_proxy.go` |
| **Windows Terminal**| Verified | Multi-pane Layouts | `pkg/terminal/layout.go` |

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
