[PROJECT_MEMORY]

# Ultimate Agentic Coding Harness Rebuild - Handoff

## Current Status
We are actively porting the `just-every-code` repository architecture identically across Rust, Go, C#, and Java.
In the previous sessions, we successfully ported the core `ToolCtx`, `ToolRuntime`, `SandboxAttempt`, and `ToolOrchestrator` execution structures into the four CLI targets. We then implemented the `AgentSession`, `AgentStatus`, and `TurnContext` structures to functionally map user input to the ToolOrchestrator, moving away from simulated mock executions.

## Action Required for Successor
**Your Immediate Goal:**
1. Now that the `ToolOrchestrator` and `TurnContext` pipelines are structurally mapped across all four languages, we need to begin porting the `mcp-server` integration logic from `just-every-code/codex-rs/mcp-server` and `codex-rs/core/src/mcp.rs`.
2. Add support to the CLIs to dynamically load MCP capabilities into the `ToolOrchestrator`.
3. Ensure absolute strictness with `.gitignore` - no compiled artifacts may be staged or committed. Keep processing the `just-every-code` submodule deep dive.
