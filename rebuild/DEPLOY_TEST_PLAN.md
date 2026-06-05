# DEPLOY TEST PLAN: Warp Rebuild v2.1.0

## Objective
Verify the stability and functionality of the Warp Go binaries across target environments.

## Test Matrix
| Platform | Architecture | Mode | Result |
| :--- | :--- | :--- | :--- |
| Linux | amd64 | CLI (Terminal) | [ ] |
| Linux | amd64 | Service (-port) | [ ] |
| macOS | amd64 | CLI (Terminal) | [ ] |
| macOS | arm64 | CLI (Terminal) | [ ] |
| Windows | amd64 | CLI (Terminal) | [ ] |

## Smoke Tests

### 1. CLI Mode (Terminal Harness)
- **Command:** `./warp`
- **Expected:** Spawns a subshell (bash), allows command entry, and handles `Ctrl+D` exit.
- **Verification:** Ensure raw terminal mode restores properly on exit.

### 2. Service Mode (AI Harness)
- **Command:** `./warp -port 10005`
- **Expected:** Starts the HTTP agent service.
- **Verification:**
  - `curl http://localhost:10005/health` returns "OK".
  - `curl -X POST http://localhost:10005/run -d '{"prompt": "hi"}'` returns a valid JSON history.

### 3. IPC Bridge Test
- **Requirement:** Warp Rust binary and Go `warp` sidecar.
- **Process:** Start Go service, run Rust bridge tests.
- **Expected:** Successful bidirectional communication.

## Regression Checklist
- [ ] No build artifacts committed to Git.
- [ ] VERSION.md matches binary --version output (when implemented).
- [ ] All tests pass on the specific deployment target.
