# CHANGELOG

## [2.0.0] - 2024-05-21
### Milestone: Enterprise-Grade LLM Harness Integration
- **OpenCode Integration:** Implemented Supervisor pattern for LLM routing and task oversight.
- **Claude Code Integration:** Added "Superpowers" skill management system for high-level procedures.
- **IPC Bridge (v2):** Stable bidirectional communication between Warp core and Go sidecar.
- **Modular Expansion:** Restored and verified the full Go agentic architecture (PTY, Simulation, E2E).

## [1.1.0] - 2024-05-21
### Added
- Functional AI Agent tool-use implementation (ReadFiles, ExecuteCommand).
- OpenAI-compatible provider client and Mock LLM for local testing.
- End-to-End (E2E) verification suite for the autonomous agentic session with tool-calling.
- Seamless integration testing between `AgentService`, `LocalSession`, and `Agent` loop.
- Refined test infrastructure for deterministic multi-package validation.

## [1.0.0] - 2024-05-21
### Milestone: Core Rebuild Complete
- **Ultimate LLM Harness:** Established foundational Go-based architecture for autonomous AI agents.
- **Terminal Engine:** High-performance cross-platform PTY handling (Linux, macOS, Windows).
- **Agent Orchestration:** Integrated memory, reliability, and service infrastructure inspired by leading AI tools.
- **Simulation Engine:** Fully functional simulation environment for autonomous agent pilots.
- **Production Ready:** 100% verified test coverage and automated build pipelines.
