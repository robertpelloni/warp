# TEAM HANDOVER: Warp Rebuild (Go)

## Overview
This document provides technical depth for the Go-based rebuild of the Warp autonomous harness, ready for Phase 5 development.

## Technical Architecture

### 1. Terminal Engine (`pkg/terminal`)
- **PTY Handling:** Uses `github.com/creack/pty` for native Unix PTY orchestration. The `LocalSession` manages the lifecycle of the child process (shell).
- **Platform Awareness:** Utilizes `session_unix.go` and `session_windows.go` with modern `//go:build` constraints for OS-specific logic (e.g., CWD retrieval and signal handling).
- **Extensibility:** The `plugins` package provides a registry for terminal enhancements, while `blocks` and `tabs` manage the visual hierarchy.

### 2. AI Agent Harness (`pkg/agent`)
- **Autonomous Loop:** The `Agent.RunLoop` implements a multi-turn Thought->Action->Observation cycle.
- **Orchestration:**
  - `supervisors`: Handles intelligent routing between different LLM providers and task oversight.
  - `skills`: A modular registry for high-level procedures (e.g., systematic debugging).
  - `workspace`: Managed capability tracking for agent sessions.
- **Reliability:** Built-in `CircuitBreaker` pattern prevents agent thrashing during LLM outages.

### 3. IPC Bridge (`app/src/ai/go_harness_bridge.rs`)
- **Communication:** Standard HTTP/JSON interface between the Warp Rust core and the Go sidecar.
- **Dual-Mode CLI:** The Go binary (`rebuild/cmd/warp`) can run as a standalone terminal harness or an integrated service (`-port` flag).

## Deployment & RELEASE
- **Build Automation:** `rebuild/script/build.sh` handles cross-compilation for Linux, macOS (amd64/arm64), and Windows.
- **Staging:** `rebuild/script/deploy_staging.sh` prepares release-ready tarballs.
- **Verification:** All components are verified by a comprehensive suite of unit, integration, and E2E tests.

## Future Development (Phase 5)
- Deep VT100 escape sequence parsing.
- Interactive TUI development using `bubbletea`.
- Production LLM provider connectivity.
- Secure credential management.
