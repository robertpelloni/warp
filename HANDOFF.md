[PROJECT_MEMORY]

# Ultimate Agentic Coding Harness Rebuild - Handoff

## Current Status
We are processing the `just-every-code` submodule. I have added it and begun porting the foundational CLI REPL structure across Rust, Go, C#, and Java, specifically setting up the `ToolOrchestrator` structures to handle tasks.

## Action Required for Successor
The code reviewer correctly noted that we are still only providing "fake" stubbed functionality (using `Thread.Sleep`) instead of functionally porting the core logic (e.g., actual agent states, tool execution mapping, and API dispatch).

**Your Immediate Goal:**
1. You must implement *actual* functional logic ported from `just-every-code/codex-rs` into the four languages.
2. Stop using simulated delays. You must port the actual definitions of `ToolCtx`, `ToolRuntime`, and the network/sandbox execution abstractions derived from `codex-rs/core/src/tools/orchestrator.rs` and `codex-rs/core/src/thread_manager.rs`.
3. Implement the real parsing and state machine for the agents, rather than string matching against `"/auto"`.
4. Ensure no build artifacts are committed.
