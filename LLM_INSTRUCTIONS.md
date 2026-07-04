# LLM Instructions

When interacting with this repository, LLM agents MUST adhere to the following rules:

1.  **Architecture:** The project is a cross-platform (Go-based) Agentic Development Environment (Warp).
2.  **State Management:** Utilize `HANDOFF.md` for session continuations.
3.  **Documentation:** Keep `ROADMAP.md`, `TODO.md`, `VISION.md`, and `CHANGELOG.md` updated.
4.  **UI Design:** Windows UI (`pkg/app/app_windows.go`) must be unified into a single dashboard. Avoid multi-window fragmentation.
5.  **Concurrency:** Background tasks (like Shadow Pilot) must never block the main UI thread. Use thread-safe communication.
6.  **Cross-Compilation:** Any platform-specific code (e.g., Win32 API calls) MUST be protected by build tags and abstracted via interfaces (e.g., `platformPTY`) with safe UNIX stubs (e.g., `app_unix.go`).
7.  **Testing:** All tests must be dynamically OS-aware (using `runtime.GOOS` instead of hardcoding `cmd.exe` or `sh`).
8.  **Dependencies:** Keep `go.mod` clean and use standard Go conventions.
9.  **Security:** Never hardcode secrets.

Failure to follow these instructions will result in rejected patches.