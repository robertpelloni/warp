# Product Spec: Warp Agent `.pi/` Cache Integration

## 1. Problem Statement
The Warp Go port's agentic execution engine relies heavily on frequent state synchronization with the remote agent server to persist token budgets (Extreme Telemetry) and session progress. Writing this state sequentially blocks the main terminal UI thread or standard output pipeline, introducing perceivable I/O latency (jitter) during rapid command bursts and remote server handshakes.

## 2. Proposed Behavior (Product Invariants)
1. **Zero-Block Terminal:** Telemetry and agent session state emissions must never block the user's keystrokes, standard output parsing, or rendering (`WM_PAINT` on Windows).
2. **Invisible Sandbox:** The local agent cache must reside strictly within a `.pi/` directory at the project root. This directory must be ignored by version control (via `.gitignore` and `.warpindexingignore`).
3. **Resiliency:** If the remote agent connection drops during a handshake or execution, the local `.pi/` cache must preserve the session state and flush it gracefully upon reconnection.
4. **Latency Cap:** Handshakes and token emissions must complete locally (in memory) within < 2ms to maintain terminal UI responsiveness.

---

# Tech Spec: Debounced SQLite Orchestrator

## 1. Context
Currently, the `pkg/agent/orchestrator.go` manages the agent lifecycle. To solve the I/O bottleneck, we recently introduced a pure-Go SQLite driver (`github.com/glebarez/sqlite`) to store the state locally inside `.pi/agent.db`.

## 2. Proposed Architecture (Write-Behind Cache)
To integrate the `.pi/` cache safely with the Warp agentic execution engine while minimizing I/O latency:

### Data Flow
1. **Memory-First:** When the terminal hooks intercept an AI command (`/ai`) or token usage update, the `Orchestrator` immediately updates its in-memory `State` struct.
2. **Lock Minimization:** The `mu.Lock()` is held only long enough to update a boolean `dirty` flag and the memory struct (typically < 0.05ms).
3. **Debounced Flush Loop:** A background Goroutine (`flushLoop()`) wakes up every 2 seconds. If `dirty` is true, it acquires the lock, quickly copies the struct, releases the lock, and performs the heavy I/O SQLite write asynchronously.
4. **Immediate Flushing (Handshake):** During critical operations like a remote server handshake (`ConnectRemote()`), the `SaveStateImmediate()` function bypasses the debouncer for a synchronous write to ensure the handshake session ID is safely persisted before proceeding.

### Modules Touched
* `pkg/agent/orchestrator.go` (Implemented the `flushLoop` and `SQLite` initialization).
* `internal/agent/cache_test.go` (Test harnesses validating the sandbox).
* `.gitignore` / `.warpindexingignore` (To strictly enforce the invisible sandbox).

## 3. Testing and Validation
* **Isolation Verification:** `TestAgentModuleInteraction` validates that spinning up the Orchestrator inside a mock Go module does not pollute `git status` and is correctly filtered by `git check-ignore`.
* **Thread Safety:** Standard Go race detector (`go test -race ./...`) will validate that the `flushLoop` goroutine does not conflict with the main terminal rendering loop when reading/writing agent state.
