# Warp: Immediate Task List

## High Priority
x 1. Implement the actual HTTP REST/Streaming clients for the `ApiProvider` interfaces across Rust, Go, Java, and C#. (Currently using Dummy Mocks).
x 2. Create the unified `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `CHANGELOG.md`, and `VERSION.md` standard documentation files at the root level.
x 3. Migrate the `MCP Server` abstractions from the generic "hello world" JSON-RPC mock into fully compliant Context/Tool passing mechanics matching the MCP specification.
4. Clone the `aider` repository as a submodule to begin evaluating and porting its repository map (Tree-sitter AST) mechanics.

## Medium Priority
1. Establish unit testing suites for the `ThreadStore` persistent layer.
2. Ensure consistent logging verbosity across the 4 languages.
3. Enhance the `mcp_shell` dynamic tool inside the Rust orchestrator to actually invoke `std::process::Command` when running in unsandboxed environments.
