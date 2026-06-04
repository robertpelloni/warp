# CHANGELOG

## [0.6.0] - 2024-05-21
### Added
- Staging deployment infrastructure (`rebuild/script/deploy_staging.sh`).
- Automated deployment bundle generation (tar.gz).
- Pre-deployment binary verification logic.

## [0.5.0] - 2024-05-21
### Added
- Deployment infrastructure with cross-platform build script (`rebuild/script/build.sh`).
- Cross-compilation support for Linux, macOS (amd64/arm64), and Windows.
- Final test verification and package preparation for deployment.

## [0.4.0] - 2024-05-21
### Added
- Comprehensive integration testing for the Agent, Terminal, and Harness packages.
- Core Simulation Engine in `pkg/simulation` based on AI Game Engine analysis.
- Pilot Autonomous Game Simulation demonstrating AI agent interaction with a simulation loop.

## [0.3.0] - 2024-05-21
### Added
- Functional PTY handling using `creack/pty`.
- Interactive terminal harness mode in Go with raw terminal support.
- Agent toolkit and loop structures from Pi-coding-agent.
- Circuit breaker logic for AI agent reliability from Antigravity 2.0.
- Agent service and health check infrastructure from Claude-mem.
- Memory management and provider orchestration from Hermes Agent.

## [0.2.0] - 2024-05-21
### Added
- Ported core terminal session and tab management from Tabby.
- Integrated persistent terminal block structures inspired by Waveterm.
- Implemented LLM provider registry and model info from Hyperharness.
- Added comprehensive unit tests for terminal and harness packages.

## [0.1.0] - 2024-05-21
### Added
- Initial project documentation (`VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`).
- Roadmap and Todo tracking.
- Started Warp Rebuild project structure.
