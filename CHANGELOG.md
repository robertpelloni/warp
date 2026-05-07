# Changelog

## [0.1.7] - 2026-05-07
- Addressed TODO in `crates/warpui_core/src/image_cache.rs` clarifying that calculating memory size for a parsed SVG tree without traversal is inaccurate.

## [0.1.3] - 2026-05-06
- Addressed TODO(AGENT-2283) in `crates/computer_use/src/lib.rs` to make `Screenshot` data an `Arc<[u8]>` instead of `Vec<u8>` so it's cheap to clone.

## [0.1.2] - 2026-05-06
- Implemented fix for ScreenshotRegion validation TODO in `crates/computer_use`.
- Updated HANDOFF notes.

## [0.1.1] - 2026-05-06
- Generated comprehensive documentation files.
- Prepared project for AI agent handoff.
