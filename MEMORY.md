# Warp: System Memories & Architectural Observations

## Discovered Traits
- The codebase enforces a rigid 4-language parity rule. Porting features must target `crates/warp-cli`, `warp-go`, `warp-csharp`, and `warp-java`.
- Build artifacts MUST be rigorously `.gitignore`d to prevent commit bloat.
- AI LLM API abstractions revolve around `ApiProviderRegistry`, decoupling the system prompt/message context generation from the streaming transport.
- MCP (Model Context Protocol) is currently mocked as a generic JSON-RPC transport wrapper but needs deep implementation for tool injection.

## Design Preferences
- Heavy reliance on asynchronous streaming: `tokio::sync::mpsc` (Rust), `<-chan` (Go), `System.Threading.Channels` (C#), and `java.util.concurrent.Flow` (Java).
- Console UX uses a unified command dispatcher (e.g., `/browser`, `/compact`, `/pimono`, `/ai`, `/prompt`).
