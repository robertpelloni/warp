# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to Semantic Versioning (managed via `VERSION.md`).

## [0.1.0] - 2026-06-14
### Added
- Initial port of Warp core logic to Go (targeting Go 1.25.0).
- `pkg/pty`: Cross-platform pseudo-terminal interface (ConPTY on Windows, pipes on Unix).
- `pkg/session`: Terminal tab and session management.
- `pkg/command`: Warp-style command engine with alias support and AI stubs.
- `pkg/editor`: Basic buffer management and language detection.
- `pkg/app`: Win32 GUI implementation and UNIX stubs for cross-platform compilation.
- Core documentation files: `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `CHANGELOG.md`, `ROADMAP.md`, `TODO.md`, `HANDOFF.md`.
