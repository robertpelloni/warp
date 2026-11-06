# Warp Roadmap: Major Milestones

- [x] **Milestone 1: Pure Go Win32 Core Architecture & Port**
  - Implement PTY abstraction, ANSI parser, theme engine, command classification, and basic Win32 window loop.
- [x] **Milestone 2: Unified Single-Dashboard UI & System Tray Integration**
  - Merge terminal output, command blocks, and IDE editor buffer into a single responsive window.
  - Integrate system tray icon with background server shutdown capability.
  - Enable `/edit` command integration.
- [ ] **Milestone 3: Advanced ANSI Terminal Emulation & Virtual Screen Buffer**
  - Complete VT100/Xterm 256-color escape code parsing, cursor positioning, and alternate screen buffer support (for `vim`, `htop`, `less`).
- [x] **Milestone 4: WebGPU Hardware-Accelerated Renderer Scaffolding**
  - Implement texture atlas glyph rendering via DirectX/Vulkan/WebGPU backend to support ultra-fast 120 FPS scrolling.
- [ ] **Milestone 5: Multi-Agent Symphony & Remote Collaboration Bus**
  - Connect terminal sessions to local/remote LLM agent orchestration networks with streaming diff application.
