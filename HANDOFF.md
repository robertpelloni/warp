[PROJECT_MEMORY]

# Ultimate Agentic Coding Harness Rebuild - Handoff

## Current Status
We are actively porting the `just-every-code` repository architecture identically across Rust, Go, C#, and Java.
In the previous session, we successfully ported the core `ToolCtx`, `ToolRuntime`, `SandboxAttempt`, and `ToolOrchestrator` execution structures into the four CLI targets. We have eliminated the mock `Thread.Sleep` auto-drive loops and implemented actual interfaces that mirror the rust trait pipeline from `just-every-code/codex-rs/core/src/tools/sandboxing.rs`.

## Action Required for Successor
**Your Immediate Goal:**
1. Proceed with porting the actual `Agent` / `CodexThread` state machine loops that consume these orchestrators. Look at how `codex-rs` handles turning natural language into tool calls (e.g., `TurnContext` generation and passing to the orchestrator).
2. Or, if applicable based on the ROADMAP, proceed with porting the `mcp-server` logic identically into all four target languages.
3. Ensure absolute strictness with `.gitignore` - no compiled artifacts may be staged or committed. Keep processing the `just-every-code` submodule deep dive.
