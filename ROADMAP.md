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
