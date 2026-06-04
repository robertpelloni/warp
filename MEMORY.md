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
