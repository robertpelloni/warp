# Warp: Immediate Task List

## High Priority
x 1. Implement the actual HTTP REST/Streaming clients for the `ApiProvider` interfaces across Rust, Go, Java, and C#. (Currently using Dummy Mocks).
x 2. Create the unified `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `CHANGELOG.md`, and `VERSION.md` standard documentation files at the root level.
x 3. Migrate the `MCP Server` abstractions from the generic "hello world" JSON-RPC mock into fully compliant Context/Tool passing mechanics matching the MCP specification.
x 4. Clone the `aider` repository as a submodule to begin evaluating and porting its repository map (Tree-sitter AST) mechanics.

## Medium Priority
x 1. Establish unit testing suites for the `ThreadStore` persistent layer.
2. Ensure consistent logging verbosity across the 4 languages.
x 3. Enhance the `mcp_shell` dynamic tool inside the Rust orchestrator to actually invoke `std::process::Command` when running in unsandboxed environments.

## New Tasks Generated from Roadmap
1. Build a basic remote Server Host for WebUI interactions in Go.
2. Build an initial WebUI to interface with the Go Server Host (React/Vite).
3. Integrate the React/Vite WebUI with the Go Server Host using WebSocket connections for real-time text streaming.

## Completed Tasks
- Added basic structual features for `tabby` in Go, Rust, C#, and Java.

5. Clone the `codex-cli` repository as a submodule to analyze and port its features.
6. Clone the `gemini-cli` repository as a submodule to analyze and port its features.

7. Port remaining `aider` core features to Go, C#, Java, and Rust.

8. Document and port core mechanics from `codex-cli` and `gemini-cli`.
- Completed documentation of `codex-cli` mechanics.

9. Build browser extension mechanics and integrate JSON-RPC interface for browser extension control.

10. Deploy and verify browser extensions for compatibility across major browsers.

11. Implement multi-agent routing.

12. Implement git automation commands and CI pipeline auto-fixes.

13. Implement comprehensive deep dependency tree static analysis functions.

14. Deploy fully functional mobile app that wraps the web interface for mobile environments.

15. Complete Phase 3 deployment.

x 16. Add detailed UI testing suite for WebUI React components.

x 17. Improve documentation coverage for `warp-go` endpoints.

x 18. Improve test coverage for Go modules.

19. Ensure consistent logging verbosity across the 4 languages.

x x 20. Verify UI jitter fix.

21. Ensure cross-platform build pipelines for binaries.

22. Final UI Polish before v1.0.

23. Add cross-language websocket tests.

x 24. Implement Neural Radiance Fields (NeRF) algorithm in Go.

x 25. Implement comprehensive unnerf algorithm including ray generation, stratified sampling, MLP evaluation, and volume rendering in Go.

30. Beta testing.
