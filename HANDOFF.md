# HANDOFF.md

Read VISION.md, ROADMAP.md, and TODO.md before making changes. Always update CHANGELOG.md and VERSION.md.

**Handoff notes from Jules:**
- Created initial set of missing `md` documentation files (`ROADMAP`, `TODO`, `CHANGELOG`, `VISION`, `MEMORY`, `DEPLOY`, `HANDOFF`, `VERSION`, `CLAUDE`, `AGENTS`, `GEMINI`, `GPT`, `copilot-instructions`, `LIBRARIES`).
- Investigated `TODO.md` file contents logic and implemented a quick-fix on `computer_use/src/lib.rs` and `computer_use/src/windows/screenshot.rs` to address an existing TODO around bounds validation `ScreenshotRegion::validate`.
- Handled numerous missing local system dependencies during cargo tests like `protobuf-compiler` and `libasound2-dev`. Tests run mostly completely on core targets but we cannot do full integration testing here easily.
- Awaiting the next AI model to pick up another item from the `TODO.md` backlog or iterate on `ROADMAP.md` tasks!
