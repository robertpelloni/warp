# Warp Go ROADMAP

## Core Vision
Warp Go aims to provide a fast, cross-platform port of the agentic coding harness from Rust to Go, featuring native speed and direct AI agent integrations inside the shell interface.

## Milestones

### Phase 1: Core Architecture & Terminal Engine (Current)
- [x] Basic PTY integration (ConPTY on Windows, standard pipes on Unix).
- [x] Terminal session management (`pkg/session`).
- [x] Basic command parser (`pkg/command`).
- [ ] Output rendering and ANSI parsing to native GUI (`pkg/app`).

### Phase 2: Cross-Platform GUI
- [ ] Fully native Win32 UI for Windows.
- [ ] Implement an equivalent Linux/macOS GUI layer.
- [ ] Integrate a cross-platform renderer capable of block-style command groupings (like Warp).

### Phase 3: AI and Telemetry Integration
- [ ] Connect the `/ai` command to the backend Agent Symphony (Claude/Gemini/GPT-4o).
- [ ] Enable Extreme Telemetry for workspace budgets.
- [ ] Connect the workspace frontend UI to these backend APIs.

### Phase 4: CI & Git Shadow Pilot
- [ ] Continuous monitoring of git diffs via a background process.
- [ ] Intelligent, autonomous anomaly detection.
