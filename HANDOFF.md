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
