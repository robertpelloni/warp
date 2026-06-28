# Warp: Autonomous Development Agent

Warp is an advanced autonomous coding harness designed to replicate and exceed the functionality of existing AI development tools. It offers full feature parity with platforms like Tabby, Wave, Codex Desktop, and Claude Desktop, all implemented across four core languages: **Rust, Go, C#, and Java**.

## Core Features
*   **Multi-Language Architecture:** Deep redundancy with complete implementations in Rust, Go, C#, and Java.
*   **Persistent State Management:** Built-in `ThreadStore` mechanics to persist context across sessions.
*   **Agent Orchestration:** Real-time routing for Architect, Executor, and Auditor agents.
*   **Headless CDP Integration:** Direct control over local browsers for autonomous web testing.
*   **Extensible Tooling:** Comprehensive integration for local filesystem, git operations, and static code analysis.
*   **MCP Compliant:** Full Context/Tool mechanics compliant with the Model Context Protocol.

## Installation & Setup

### Prerequisites
*   Make sure you have installed standard development tools for whichever backend you prefer:
    *   **Rust:** `cargo build`
    *   **Go:** `go build`
    *   **C#:** `dotnet build`
    *   **Java:** `mvn compile`
*   **Node.js (for UI):** `npm run build` inside the `warp-ui` directory.

### Quick Start
To launch the core server (Go version):

```bash
cd warp-go
go run ./cmd/warp-server
```

To run the WebUI:

```bash
cd warp-ui
npm install
npm run dev &
```

For full environment deployment details, see `DEPLOY.md`.
