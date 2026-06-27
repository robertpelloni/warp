<a href="https://www.warp.dev">
    <img width="1024" alt="Warp Agentic Development Environment product preview" src="https://github.com/user-attachments/assets/9976b2da-2edd-4604-a36c-8fd53719c6d4" />
</a>

<p align="center">
  <a href="https://www.warp.dev">Website</a>
  ·
  <a href="https://www.warp.dev/code">Code</a>
  ·
  <a href="https://www.warp.dev/agents">Agents</a>
  ·
  <a href="https://www.warp.dev/terminal">Terminal</a>
  ·
  <a href="https://www.warp.dev/drive">Drive</a>
  ·
  <a href="https://docs.warp.dev">Docs</a>
  ·
  <a href="https://www.warp.dev/blog/how-warp-works">How Warp Works</a>
</p>

> [!NOTE]
> OpenAI is the founding sponsor of the new, open-source Warp repository, and the new agentic management workflows are powered by GPT models.

<h1></h1>

## About Warp

**Warp is the Ultimate Agentic Coding Harness.** It is designed to achieve full feature parity with top-tier AI developer tools (Tabby, Claude Desktop, Cursor, Codex CLI) while extending capabilities through a multi-model orchestration engine. Warp seamlessly connects Architects, Executors, and Auditor agents across multiple execution environments.

Unique to Warp is its **Four-Pillar Redundancy**: The core execution engine is maintained simultaneously in **Rust**, **Go**, **C#**, and **Java**, ensuring that agent workflows, ThreadStore persistency, context compaction, and LLM ApiProviders behave identically regardless of backend infrastructure.

## Key Features (Alpha Release v0.1.1)

*   **Multi-Provider LLM Streaming**: Native `ApiProviderRegistry` mapping generic contexts to live HTTP Server-Sent Event (SSE) streams for OpenAI and Anthropic models.
*   **MCP Server Protocol**: Full adherence to the Model Context Protocol (MCP) using native `JsonRpcRequest` structures for dynamic tool invocation.
*   **Dynamic Shell Execution**: The `mcp_shell` orchestrator tool actively invokes host system commands natively across language backends for true agentic execution.
*   **Zero-Friction Context**: Includes `ThreadStore` persistence databases and `SessionManager` contexts for model handoffs without memory loss.
*   **Aider Submodule Integration**: Scaffolding for Tree-sitter AST syntax mapping across full codebases using NetworkX-style PageRank rendering.

## Setup & Installation

Warp is structured as a mono-repo encompassing multiple build targets. Ensure you have the appropriate toolchains installed:
- **Rust**: `cargo` (1.70+)
- **Go**: `go` (1.26.0+)
- **C#**: `.NET 8.0 SDK`
- **Java**: `JDK 21+` & `Maven`

### Environment Variables
To enable cloud models, configure your environment:
```bash
export OPENAI_API_KEY="your-key-here"
```

### Building & Running

**Go Backend (Active Server Target):**
```bash
cd warp-go
go build ./cmd/warp-cli
./warp-cli
```

**Rust CLI:**
```bash
cd crates/warp-cli
cargo run --release
```

**C# Runtime:**
```bash
cd warp-csharp/WarpCsharp
dotnet run
```

**Java Implementation:**
```bash
cd warp-java
mvn clean compile exec:java -Dexec.mainClass="dev.warp.Main"
```

## Basic Usage

When launching any of the CLI tools, you enter the Interactive Agent REPL. Use the following commands:

*   `/prompt <msg>`: Send a direct instructional payload to the agent for evaluation.
*   `/ai`: Executes a demonstration streaming sequence to live LLM providers using predefined Context blocks.
*   `/browser`: Triggers the Headless CDP browser automation interface test.
*   `/pimono`: Validates the `SessionManager` context history tracking.
*   `/compact`: Demonstrates memory context compaction logic.
*   `/store`: Tests the `ThreadStore` memory-to-disk persistent writes.

## External Submodules
Make sure to initialize external integration submodules upon cloning:
```bash
git submodule update --init --recursive
```
*`submodules/aider`: For AST map processing.*
*`submodules/pi-mono`: For compaction abstraction.*

## Open Source & Contributing

Warp's client codebase is open source and lives in this repository. We welcome community contributions and have designed a lightweight workflow to help new contributors get started. For the full contribution flow, read our [CONTRIBUTING.md](CONTRIBUTING.md) guide.

### Issue to PR

Before filing, search existing issues for your bug or feature request. If nothing exists, file an issue using our templates. Security vulnerabilities should be reported privately as described in [SECURITY.md](SECURITY.md).

Once filed, a Warp maintainer reviews the issue and may apply a readiness label. Anyone can pick up a labeled issue — mention **@oss-maintainers** on an issue if you'd like it considered for a readiness label.

## Licensing

Warp's UI framework (the `warpui_core` and `warpui` crates) are licensed under the [MIT license](LICENSE-MIT). The rest of the code in this repository is licensed under the [AGPL v3](LICENSE-AGPL).
