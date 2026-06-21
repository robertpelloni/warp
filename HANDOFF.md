[PROJECT_MEMORY]

# Ultimate Agentic Coding Harness Rebuild - Handoff

## Current Status
We are processing the `just-every-code` submodule. I have added it and begun porting the foundational CLI REPL structure across Rust, Go, C#, and Java.

## Action Required for Successor
The code reviewer correctly noted that we are only providing "fake" stubbed functionality instead of functionally porting the core logic (e.g., the true multi-agent orchestration, auto-drive loop, and MCP server).

**Your Immediate Goal:**
1. You must implement *actual* logic ported from `just-every-code/codex-rs` into the four languages.
2. Begin by porting the `Auto Drive` orchestration loop logic or the multi-agent `mcp-server` integration, rather than just mocking the command outputs.
3. Remove `pi-mono` from the staging area so we process *one* submodule at a time.
4. Ensure no build artifacts are committed.
