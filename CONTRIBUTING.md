# Contributing to Warp

We welcome contributions to Warp! As this project is currently in the **Alpha Phase**, things are moving quickly and breaking changes are expected.

## Getting Started

1. **Fork the Repository:** Create your own fork and branch off `master`.
2. **Environment Setup:** Ensure you have the necessary SDKs (`go`, `rustc`, `dotnet`, `java`, `node`). See `README.md` and `DEPLOY.md`.
3. **Run the Tests:** Before writing new code, confirm your environment runs `test_all.sh` cleanly.

## The Four-Pillar Philosophy

Warp maintains strict feature parity across Rust, Go, C#, and Java. When adding a new core feature (e.g., an endpoint like `/nerf`), you **must** implement it across all four backends.

## Submitting Pull Requests

1. Include a comprehensive test plan for your feature.
2. If changing the Web UI, ensure no jitter regressions are introduced. Run the React `vitest` suite.
3. Update `CHANGELOG.md` and `TODO.md` as appropriate.
4. Open your PR against `master`.

## Bug Reports & Feature Requests

Please use GitHub Issues. Provide exact reproduction steps and terminal output when applicable.
