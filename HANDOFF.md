# HANDOFF.md

Read VISION.md, ROADMAP.md, and TODO.md before making changes. Always update CHANGELOG.md and VERSION.md.

**Handoff notes from Jules:**
- Created initial set of missing `md` documentation files (`ROADMAP`, `TODO`, `CHANGELOG`, `VISION`, `MEMORY`, `DEPLOY`, `HANDOFF`, `VERSION`, `CLAUDE`, `AGENTS`, `GEMINI`, `GPT`, `copilot-instructions`, `LIBRARIES`).
- Investigated `TODO.md` file contents logic and implemented a quick-fix on `computer_use/src/lib.rs` and `computer_use/src/windows/screenshot.rs` to address an existing TODO around bounds validation `ScreenshotRegion::validate`.
- Handled numerous missing local system dependencies during cargo tests like `protobuf-compiler` and `libasound2-dev`. Tests run mostly completely on core targets but we cannot do full integration testing here easily.
- Awaiting the next AI model to pick up another item from the `TODO.md` backlog or iterate on `ROADMAP.md` tasks!
- Further addressed TODOs in `ui_components` regarding text, buttons, handle checks, and cleaned up unused tests around `evict_size`.
- Finished integrating dependency solutions and ensured the entire test suite `cargo test` and `presubmit` passes fully with all system dependencies.
- Synced with upstream/master and cleanly resolved everything.
- Verified that everything builds properly after the upstream merge by running cargo check/test on ai, warpui_core, warp_core, computer_use, and the full workspace.
- Addressed SVG size calculation documentation in image_cache.
- Synced with upstream/master and cleanly resolved everything.
- Addressed further minor documentation TODOs across `warpui_core` files including accessibility, components, events, and image_cache.
- Removed stale and obsolete `TODO` comments from `warpui_core` such as random arbitrary font selection, button borders, and obsolete refactoring notes for `View` and `Entity`.
- Addressed further minor documentation TODOs across `warpui_core` files including accessibility, components, events, and image_cache.
- Removed stale and obsolete `TODO` comments from `warpui_core` such as random arbitrary font selection, button borders, and obsolete refactoring notes for `View` and `Entity`.
- Removed stale and obsolete `TODO` comments from `warpui_core` such as scene z-index notes, image_cache types, and text layout font sizing.
- All remaining TODOs in `warpui_core` have been analyzed and updated correctly where necessary.
- Verified all structural updates completed and submitted successfully over previous phases.
- Verified all structural updates completed and submitted successfully over previous phases.
