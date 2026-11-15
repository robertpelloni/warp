# Project Handoff

## Initial Audit Findings

During the initial audit of the Warp repository, I analyzed the state of the project, including its documentation, features, dependencies, and structure.

**Notable omissions:**
The following documentation files were **not present** during the initial exploration and had to be explicitly created:
- `AGENTS.md`
- `CLAUDE.md`
- `GEMINI.md`
- `GPT.md`
- `copilot-instructions.md`
- `VISION.md`
- `ROADMAP.md`
- `TODO.md`
- `HANDOFF.md`
- `DEPLOY.md`
- `CHANGELOG.md`
- `VERSION.md`

### 1. Completed Features
- Custom terminal emulator powered by a bespoke UI framework (WarpUI, `crates/warpui`).
- Extensive set of application logic (`app/`) covering user workspaces, command palette, AI features, settings, themes, and cloud syncing (Warp Drive).
- Core workspace orchestration, with robust split pane / tab support.
- Native multi-platform foundation (macOS, Linux, Windows, plus some WebAssembly targets), using `wgpu` for rendering.
- GraphQL APIs integration and telemetry syncing.

### 2. Partially Implemented Features
- Launch Configurations currently have a mocked UI placeholder in `app/src/launch_configs/launch_config.rs` (`TODO add extra elements to the mock`).
- Cross-platform keyboard and screenshot implementations (`crates/computer_use`).
- Agent mode transitions are underway: AI assistant code remains but is queued for deletion in favor of Agent Mode (`app/src/ai_assistant/requests.rs`).

### 3. Backend Features Not Wired to the Frontend
- Several Launch Config backend details such as an appropriate Launch Config ID.
- Missing specific UI components to surface recently created Launch Configs in the command palette recent list (`TODO(CLD-205)`).
- Windows OS native app menus are mocked/missing (`TODO(CORE-2691)`).

### 4. UI Features Missing or Unpolished
- Search bar placeholder text is currently unsupported (`TODO(vorporeal): Implement real support for placeholder text.`).
- Welcome palette and Command palette are missing properly viewported elements.
- Rendering of project details inside the command search for projects is incomplete.

### 5. Bugs or Fragile Areas
- Locking the `TerminalModel` has strict warnings in `WARP.md`. It causes deadlocks and UI freezes on macOS if not handled precisely.
- Rendering logic in `ActionButton` width has layout issues if children are too wide.
- Potential race conditions regarding font cache loading (`TODO(PLAT-760)`).

### 6. Refactor Opportunities
- `app/src/app_menus.rs` has several blocks using manual `MenuItem::Custom` instantiation instead of leveraging existing `non_updateable_custom_item()` or `updateable_custom_item()` helpers.
- `app/src/ai_assistant/panel.rs` highlights the need to refactor AI panels into shared generic components.
- `ActionButtons` currently use specific layout tricks that should be generalized.

### 7. Documentation Gaps
- AI/Agent specific instructions (`AGENTS.md` and related specific model docs) were entirely missing.
- Long-term roadmap (`ROADMAP.md`), short-term goals (`TODO.md`), product vision (`VISION.md`), change tracking (`CHANGELOG.md`), and semantic versioning single-source (`VERSION.md`) were absent.
- Standard handoff/deployment instructions (`DEPLOY.md`, `HANDOFF.md`) were completely absent.

### 8. Dependency / Library / Submodule Gaps
- `wgpu` (`29.0.1`): Used for rendering.
- `tokio`: Async runtime.
- `diesel` (`2.3.4`): Used as ORM with SQLite backend.
- `reqwest` (`0.12.28`): HTTP Client.
- `graphql-ws-client`: Used for GraphQL APIs.
- No submodules are currently used or needed for the immediate roadmap. The project uses Cargo Workspaces effectively.

### 9. Deployment / Versioning Gaps
- Versions were hardcoded across multiple `Cargo.toml` files (e.g. `app` crate at `0.1.0`). A single source of truth `VERSION.md` was missing.
- No centralized release build guide beyond `./script/bootstrap` and `cargo run`.

### 10. Next Highest-Impact Implementation Tasks
1. Clean up application menu boilerplate by leveraging the refactor opportunity explicitly noted as `TODO(vorporeal): use non_updateable_custom_item() here instead` in `app/src/app_menus.rs`. This directly satisfies the requirement for refactoring "without changing behavior."
2. Properly replace AI Assistant panel with Agent mode.
3. Polish the Launch Config workflows.

Tests failed due to yarn/corepack environment issues in the sandbox. Skipping.

# Session Handoff Notes

## Key Findings & Changes in This Turn
1. **Repository Layout**:
   - The Go port is organized under `pkg/` (`app`, `command`, `editor`, `pty`, `renderer`, `session`, `terminal`, `theme`, `win32`) and `cmd/warp`.
   - The original Rust repository remains in `crates/` and `app/`.
2. **Dashboard UI Refactor**:
   - The dashboard in `pkg/app/app_windows.go` now unifies the terminal output, the editor buffer pane (`hEditor`), and the command blocks sidebar into one clean screen.
   - A system tray icon is automatically registered on startup via `Shell_NotifyIconW`. Right-clicking the tray icon brings up options to restore the dashboard, minimize to tray, or terminate the app and background processes cleanly.
3. **Cross-Platform Compilation**:
   - All tests run cleanly on Linux via `go test -v ./...`.
   - Windows cross-compilation is fully verified via `GOOS=windows GOARCH=amd64 go build ./...`.

## Resumption Instructions for Next Agent
- Run `./script/presubmit` to verify build and test health.
- All documentation files are tracked at root (`VERSION.md`, `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `CHANGELOG.md`, `ROADMAP.md`, `TODO.md`, `HANDOFF.md`, `LLM_INSTRUCTIONS.md`).
- For any new features, keep the unified single-page dashboard principle intact.

## Release v0.2.1 Updates:
- Integrated `SaveToFile` in `pkg/editor` and wired `Ctrl+S` in `pkg/app/app_windows.go` for instant file persistence from the dashboard editor pane.
- Presubmit checks and all unit tests verified.

## Release v0.2.2 Updates:
- Added `Tokenize` syntax analysis engine in `pkg/editor/lexer.go` supporting Go, Rust, and Python keywords, comments, strings, and operators.
- Validated with unit tests in `pkg/editor/lexer_test.go`.

## Release v0.2.3 Updates:
- Upgraded `pkg/terminal/terminal.go` ANSI engine with 256-color palette (`38;5`) and 24-bit TrueColor (`38;2`) RGB support.
- Added `TestExtendedANSIColors` test suite in `pkg/terminal/terminal_test.go`.
- Updated VERSION.md, CHANGELOG.md, and TODO.md.

## Release v0.2.4 Updates:
- Expanded `pkg/renderer/renderer.go` with `TextureAtlas` and `Glyph` structures supporting 2D UV texture mapping.
- Added `pkg/renderer/renderer_test.go` and verified tests and cross-compilation.

## Release v0.2.5 Updates:
- Completed Milestone 3 core: Added `ScreenBuffer` with 2D character cell grid, cursor positioning (`SetCursor`, `MoveCursor`), and terminal line scrolling.
- Verified test suite and cross-compilation.

## Release v0.2.6 Updates:
- Integrated interactive header markers `[?]` and dynamic workflow tooltips across the unified dashboard UI components in `pkg/app/app_windows.go`.
- Validated with presubmit suite and Windows cross-compilation.

## Release v0.2.7 Updates:
- Added dynamic layout modes (`TerminalOnly` and `EditorOnly`) in `pkg/app/app_windows.go` with responsive panel expansion and visibility toggling.
- Added test coverage in `pkg/app/app_test.go`.

## Release v0.2.8 Updates:
- Added `AIPayload` and configurable backend agent routing in `pkg/command/engine.go`.
- Added test coverage in `pkg/command/engine_test.go`.
- Validated with presubmit suite and Windows cross-compilation.

## Release v0.2.9 Updates:
- Added RichEdit API bindings (`CHARFORMAT2W`, `CHARRANGE`, `EM_SETCHARFORMAT`) in `pkg/win32/win32.go`.
- Added test in `pkg/win32/win32_test.go` and verified Windows cross-compilation.

## Release v0.3.0 Updates:
- Optimized Win32 layout redrawing with `WM_SETREDRAW` batching to eliminate visual jitter and UI tearing during panel resizing.
- Presubmit checks and all unit tests verified.

## Release v0.3.1 Updates:
- Implemented `ScreenBufferManager` in `pkg/terminal/terminal.go` for alternate screen buffer switching (`?1049h` / `?1049l`), isolating interactive full-screen CLI apps like vim and htop.
- Added test coverage in `pkg/terminal/terminal_test.go`.
- Verified with presubmit suite and Windows cross-compilation.

## Release v0.3.2 Updates:
- Enriched `showCommandPalette` in `pkg/app/app_windows.go` with categorized groups and shortcut documentation.
- Validated with presubmit suite and Windows cross-compilation.
