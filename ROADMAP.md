# Product Roadmap

## Q1 - The Agentic Shift
- Complete the transition from the legacy AI Assistant panel to full Agent Mode (`TODO(roland)`).
- Finalize the new Command Palette and properly viewport elements.
- Clean up the UI frameworks with stabilized primitives (e.g. `ActionButton` layouts).

## Q2 - Refined Workflows
- Launch Configurations: Polish the mock UI and integrate proper Backend IDs.
- Implement robust placeholder and auto-completion support across all input elements.
- Implement proper project details rendering for local and synced projects.

## Q3 - Deep Integration
- Enhance the `computer_use` crate to handle advanced screen and mouse telemetry.
- Support full native OS integrations (Windows app menus, macOS dock interactions).

## Q4 - Ecosystem
- Open up plugin/workflow APIs based on WASM.
- Polish collaborative features, sharing entire workspace environments via Warp Drive.

# Warp Roadmap: Major Milestones

- [x] **Milestone 1: Pure Go Win32 Core Architecture & Port**
  - Implement PTY abstraction, ANSI parser, theme engine, command classification, and basic Win32 window loop.
- [x] **Milestone 2: Unified Single-Dashboard UI & System Tray Integration**
  - Merge terminal output, command blocks, and IDE editor buffer into a single responsive window.
  - Integrate system tray icon with background server shutdown capability.
  - Enable `/edit` command integration.
- [x] **Milestone 3: Advanced ANSI Terminal Emulation & Virtual Screen Buffer**
  - Complete VT100/Xterm 256-color escape code parsing, cursor positioning, and alternate screen buffer support (for `vim`, `htop`, `less`).
- [x] **Milestone 4: WebGPU Hardware-Accelerated Renderer Scaffolding**
  - Implement texture atlas glyph rendering via DirectX/Vulkan/WebGPU backend to support ultra-fast 120 FPS scrolling.
- [x] **Milestone 5: Multi-Agent Symphony & Remote Collaboration Bus Scaffolding**
  - Connect terminal sessions to local/remote LLM agent orchestration networks with streaming diff application.

# Warp: Autonomous Development Agent - Master Roadmap

## Core Directive
To build the ULTIMATE AGENTIC CODING HARNESS achieving full feature parity with Tabby, Warp, Hyper, Wave, Codex Desktop, Antigravity 2.0, Claude Desktop, Claude Code, Codex CLI, Gemini-cli, Opencode, Pi-mono/Pi-coding-agent, Hermes-agent, Hermes Desktop, Adrenaline CLI, Aider CLI, Amazon Q CLI, Amazon Q Developer CLI, Amp Code CLI, Auggie CLI, Azure OpenAI CLI, Bito CLI, Byterover CLI, Codebuff CLI, Codemachine CLI, Copilot CLI, Crush CLI, Dolt CLI, Factory CLI, Goose CLI, Grok CLI, Jules CLI, Kilo Code CLI, Kimi CLI, LLM CLI, LiteLLM CLI, Llamafile CLI, Manus CLI, Mistral Vibe CLI, Ollama CLI, Open Interpreter CLI, Qwen Code CLI, RowboatX CLI, Rovo CLI, Shell Pilot CLI, Smithery CLI, and Trae CLI.

## Phase 1: Architectural Foundation (In Progress)
- [x] Multi-language core CLI scaffolding (Rust, Go, C#, Java).
- [x] ThreadStore Persistent State abstraction ported across all 4 environments.
- [x] Session Manager & Context Compaction mechanics ported (`pi-mono` integration).
- [x] Headless browser CDP integration abstractions (`BrowserManager`, `CdpPage`).
- [x] Generic LLM Provider Registry and stream mechanics (`Model`, `Context`, `EventStream`).

## Phase 2: Agent Orchestration & Tooling
- [x] Incorporate comprehensive sandbox filesystem management tools.
- [x] Add `git` automation commands, automated patch-resolution, and CI pipeline auto-fixes.
- [x] Implement multi-agent routing (Architect, Executor, Auditor).
- [x] Connect LLM API Provider Registry to external HTTP streaming endpoints (OpenAI, Anthropic API wiring).
- [x] Implement deep dependency tree static analysis functions.

## Phase 3: External Platform & IDE Integration
- [x] Implement MCP Server features allowing Warp to act as a toolset within Claude Desktop / Cursor.
- [x] Expose an internal JSON-RPC interface for browser extension control.
- [x] Implement Mobile-friendly WebUI for remote Agent operation.
- [x] Develop native integration with Chromium DevTools to scrape network activity.

## Phase 4: Production & Deployment
- [x] Stabilize all endpoints.
- [x] Remove UNDER CONSTRUCTION banner.

## Phase 9: End-to-End Test and Polish
- [ ] Add integration tests.

## Phase 14: Release v1.0
- [ ] Tag and cut v1.0 release.

# ROADMAP

## Phase 1: Foundation (Complete)
- [x] Initialize Go project structure in `rebuild/`.
- [x] Port core terminal emulation (PTY handling).
- [x] Implement basic terminal harness in Go.
- [x] Implement Simulation Engine and Pilot Autonomous Loop.
- [x] Comprehensive Test Verification and Staging Readiness.
- [x] Final Production Release Prep (v1.0.0).

## Phase 2: Warp Feature Parity (Complete)
- [x] Port Warp Blocks (Foundational persistent state).
- [x] Port Warp AI Agent Mode (Autonomous tool-calling loop).
- [x] Port Warp Drive (Sync) foundation.

## Phase 3: The Ultimate Harness Integration (Complete)
- [x] Integrate Go LLM Harness with Rust Core (IPC Bridge).
- [x] Establish Functional Agentic Foundation (Milestone 1.3.0).

## Phase 4: Feature Expansion (Complete)
- [x] Integrate OpenCode (Supervisors).
- [x] Integrate Claude Code Superpowers (Skill Manager).

## Phase 5: Deep Parity & Polish (Current)
- [ ] Integrate Tabby features.
- [ ] Integrate Hyper features.
- [ ] Integrate Wave features.
- [ ] Integrate Codex Desktop features.
- [ ] Integrate Antigravity 2.0 features.
- [ ] Integrate Claude Desktop/Code features.
- [ ] Integrate Gemini-CLI features.
- [ ] Integrate Hermes features.

## Phase 6: Optimization & Polish
- [ ] Advanced performance tuning.
- [ ] Comprehensive documentation.
- [ ] Public Beta Release.
