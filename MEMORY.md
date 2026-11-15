# Internal Architectural Memory & Codebase Traits

## Key System Discoveries & Decisions
- **Cross-Platform PTY Abstraction**:
  - `pkg/pty` defines `platformPTY` interface (`Write`, `Read`, `Resize`, `Close() error`).
  - Windows leverages ConPTY through raw kernel32 syscalls (`procCreatePseudoConsole`, `procInitializeProcThreadAttributeList`).
  - Non-Windows environments use multi-reader piped subprocesses (`startUnix`) with graceful fallback.
- **Pure Go Win32 Dashboard UI**:
  - Windows GUI in `pkg/app/app_windows.go` uses standard Win32 edit controls, custom brushes, and fonts without CGO dependencies.
  - Headless/CLI fallback stub is provided in `pkg/app/app_unix.go` guarded by `//go:build !windows` for Linux/CI execution.
  - Native System Tray integration uses `Shell_NotifyIconW` with custom message `WM_APP_TRAYMSG` and `TrackPopupMenu` for process and server management.
- **Unified Single-Dashboard Layout**:
  - Terminal output, IDE editor pane, and command blocks sidebar are displayed side-by-side with responsive viewport collapsing.
- **Command Processing Engine**:
  - `pkg/command` parses internal commands (`/help`, `/warp`, `/ai`, `/edit`, `/open`), custom shell aliases (`ll`, `gs`, `gd`, `gl`), and passes normal commands to the active PTY session.

# Warp: System Memories & Architectural Observations

## Discovered Traits
- The codebase enforces a rigid 4-language parity rule. Porting features must target `crates/warp-cli`, `warp-go`, `warp-csharp`, and `warp-java`.
- Build artifacts MUST be rigorously `.gitignore`d to prevent commit bloat.
- AI LLM API abstractions revolve around `ApiProviderRegistry`, decoupling the system prompt/message context generation from the streaming transport.
- MCP (Model Context Protocol) is currently mocked as a generic JSON-RPC transport wrapper but needs deep implementation for tool injection.

## Design Preferences
- Heavy reliance on asynchronous streaming: `tokio::sync::mpsc` (Rust), `<-chan` (Go), `System.Threading.Channels` (C#), and `java.util.concurrent.Flow` (Java).
- Console UX uses a unified command dispatcher (e.g., `/browser`, `/compact`, `/pimono`, `/ai`, `/prompt`).

# MEMORY: Architectural Observations & Design Preferences

## Initial State (Rust Implementation)
- **UI Framework:** Custom `warpui` based on an Entity-Component-Handle pattern. Highly modular but complex.
- **Terminal:** Based on Alacritty's core logic but heavily modified for "blocks" and AI integration.
- **Terminal Model (`TerminalModel`):**
    - Managed via `BlockList` (history) and `AltScreen` (fullscreen).
    - Tracks shell/SSH launch state, title stack, and color overrides.
    - Handles complex output streams: in-band command output, completions, and terminal-specific image protocols (ITerm, Kitty).
    - Event-driven via `ChannelEventListener`.
- **AI Integration:** Uses a sophisticated "Agent Mode" with multi-turn conversations and tool use.
- **AI Agent Actions (`AIAgentActionType`):**
    - Extensive toolset: `RequestCommandOutput`, `WriteToLongRunningShellCommand`, `ReadFiles`, `SearchCodebase`, `Grep`, `FileGlob`.
    - Support for MCP (Model Context Protocol): `ReadMCPResource`, `CallMCPTool`.
    - Capability for computer use: `UseComputer`.
    - Multi-agent support: `StartAgent`, `SendMessageToAgent`.
- **Persistence:** SQLite via Diesel.
- **Communication:** GraphQL for client-server sync.

## Rebuild Preferences (Go Implementation)
- **Concurrency:** Utilize Go routines and channels for terminal IO and AI streaming.
- **UI:** Explore Go-native UI libraries (e.g., Gio, Fyne) or a hybrid approach if necessary to maintain performance.
- **Modularity:** Ensure a clean separation between the terminal engine, the UI framework, and the AI agent harness.
- **Tooling:** Prefer standard Go patterns (interfaces, composition) over complex inheritance or ECS if possible, for simplicity and maintainability.
