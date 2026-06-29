# Changelog

## [0.1.1] - 2026-06-24
### Added
- Multi-provider LLM API abstract ports across Rust, Go, Java, and C#.
- Dual streaming demonstration via `/ai` repl command.
- Asynchronous stream handling (`mpsc`, `Flow.Publisher`, `ChannelReader`, `<-chan`).
- Core documentation initialization (`VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `TODO.md`, `ROADMAP.md`).
- Intelligent repository synchronization scripts mapped to `HANDOFF.md`.

## [0.1.0] - Initial Release
- Implemented `ThreadStore` context layers.
- Ported `pi-mono` Context Compaction mechanics.
- Headless Browser (CDP) skeleton classes.

## [v1.0.0-alpha] - 2026-06-28
### Added
- Core Go agent loop with persistent state tracking in `.pi/`
- Full NeRF pipeline (Ray generation, stratified sampling, MLP evaluation, Volume rendering)
- JSON-RPC abstraction layer
- `git`, `fs`, `ast`, and `cdp` modules
- UI jitter fixes and end-to-end testing hooks
