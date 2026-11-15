# Product Vision

Warp's long-term goal is to be the premier agentic development environment, born out of the terminal.
It aims to seamlessly blend a high-performance, GPU-accelerated interface with rich AI interactions to completely augment the typical developer workflow. By treating the terminal as an application platform—where terminal outputs, editors, notebooks, and web views all live side-by-side—Warp empowers users to write, build, run, and debug code faster.

**Core Pillars:**
1. **Performance:** GPU-accelerated rendering and a lightweight footprint.
2. **Intelligence:** Integrated agents that have full context of the user's codebase, environment, and history.
3. **Collaboration:** First-class support for sharing workflows, terminals, and insights (Warp Drive).
4. **Extensibility:** Support for bringing any agent (Claude, GPT, etc.) or integrating any workflow seamlessly.

# Warp Vision & Architecture

## The Ultimate Vision
Warp is the modern, agentic development environment and intelligent terminal emulator designed for high-performance software engineering teams. By transforming traditional terminal workflows into an integrated, collaborative, and AI-native harness, Warp bridges the gap between low-level terminal execution and modern IDE capabilities.

## Foundational Concepts
1. **Unified Dashboard Experience**:
   - Eliminate scattered windows and excessive subpages. All critical functions—terminal output, command block history, IDE buffer editor, and system status—are consolidated into a singular, highly responsive interface.
2. **First-Class Agentic Symphony**:
   - AI assistance is not an afterthought or external chat sidebar; it is embedded directly in command evaluation (`/ai`), file manipulation (`/edit`), and workflow orchestration.
3. **Cross-Platform Parity with Native Performance**:
   - Warp provides high-fidelity pseudo-terminal operations (ConPTY on Windows, POSIX PTY pipes on Unix/Linux). Native platform controls deliver crisp typography, zero rendering latency, and low memory footprints.
4. **Resilient Background Execution & System Tray**:
   - Dev servers and long-running agent processes can live in the background, controlled cleanly via the system tray, with clear shutdown protocols to terminate background processes predictably.

## User-Satisfaction Architecture
- **Command Blocks**: Terminal output is organized into distinct, navigable command blocks with timing metadata and exit codes.
- **Instant Palette & Keyboard-First Navigation**: Full keyboard control (`Ctrl+Enter`, `Ctrl+K`, `Ctrl+L`, `Ctrl+T`, `Ctrl+Shift+T`).
- **Zero-Friction File Interaction**: Inspect and edit code buffers directly from the terminal prompt via `/edit <path>`.

# Warp: Vision & Architecture Document

## Core Foundational Concepts
The ultimate goal of Warp is to create the **Ultimate Agentic Coding Harness**. It aims to achieve absolute feature-parity with an extensive list of the most sophisticated AI developer tools in existence, encompassing browser extensions, CLI utilities, rich GUI applications, and local model orchestration tools.

Warp does not seek to be a wrapper; it seeks to be the single, unified harness that can natively run an autonomous swarm of models (Architects, Executors, Auditors) to handle end-to-end software delivery.

## Architectural Traits & Design Preferences
1. **Four-Pillar Redundancy:** Warp must maintain exact functionality logic across four robust, statically-typed programming languages: Go, C#, Java, and Rust. No feature can exist in only one language.
2. **Zero-Friction Context:** Agents must have instant state-restoration capabilities (`HANDOFF.md`, compaction, thread persistent stores) allowing for context offloading and model-handoffs without losing memory.
3. **Multi-Model Orchestration:** The `ApiProvider` abstractions must remain totally generic. Warp should seamlessly pipeline context streams out of local models (via Ollama/Llamafile) and cloud models (Anthropic, OpenAI, Bedrock).
4. **Agent Symphony over Solitary Agents:** We strictly implement separation of concerns. One model queries architecture, one model executes code changes, and an Auditor model validates test regressions.

## Documentation Governance
* `VISION.md`: (This file) Extreme-detail outline of ultimate goals.
* `MEMORY.md`: Core system design parameters and structural shifts.
* `DEPLOY.md`: Cross-language build pipelines.
* `CHANGELOG.md`: Strict tracking of updates mapped against `VERSION.md`.
* `IDEAS.md`: Free-form innovation logic.
