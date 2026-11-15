# TODO

## High Priority
1. **App Menus Refactor:** Replace verbose `MenuItem::Custom` blocks with `non_updateable_custom_item()` helpers in `app/src/app_menus.rs` to clean up codebase.
2. **Launch Config IDs:** Add proper Launch Config IDs to data sources.
3. **Agent Mode cleanup:** Remove deprecated AI assistant codebase elements that conflict with Agent Mode.

## Medium Priority
4. **Search Bar Placeholders:** Implement real support for placeholder text in the search bar.
5. **UI Rendering:** Address layout issues in `ActionButton` when children exceed maximum widths.
6. **Viewporting:** Fix viewported elements inside the welcome and command palettes.

## Low Priority
7. **Native App Menus:** Implement native Windows OS app menus.
8. **Font Caching:** Implement retry mechanisms for partially loaded font families.

# Granular Tasks & Short-Term TODO

- [x] Abstract `platformPTY` interface in `pkg/pty/pty.go`.
- [x] Create non-Windows app stub `pkg/app/app_unix.go`.
- [x] Add system tray API bindings in `pkg/win32/win32.go`.
- [x] Add system tray icon and context menu in `pkg/app/app_windows.go`.
- [x] Wire IDE Editor pane (`hEditor`) into unified dashboard layout.
- [x] Implement `/edit <file>` command in `pkg/command/engine.go`.
- [x] Write unit tests for `pkg/pty` and `/edit` command execution.
- [x] Create governance documentation (`VERSION.md`, `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `CHANGELOG.md`, `ROADMAP.md`, `TODO.md`, `HANDOFF.md`, `LLM_INSTRUCTIONS.md`).
- [x] Enhance ANSI escape sequence parsing for complex cursor movements, 256 colors, and TrueColor.
- [x] Add syntax highlighting lexer for Go, Rust, and Python in the IDE editor pane.
- [x] Implement file save listener (`Ctrl+S`) in the editor pane.

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

x 37. Verify integration with `npm run test`.

38. Implement initial C# translation framework.

x 39. Resolve gitignore tracking for `warp-csharp/Server`.

40. Add proper .NET specific tools for compilation.

41. Add setup scripts for streamlined contributor onboarding.

x 49. Define main package entry point and primary types/interfaces.

50. Refine CLI architecture.
