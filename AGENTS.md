# Agent Instructions

This document provides universal guidelines for all AI agents operating within the Warp repository.

1. **Verify your work:** Never assume a command succeeded. Always check file contents or run `cargo check` after edits.
2. **Obey the Lints:** The project relies heavily on `cargo fmt` and `cargo clippy`. Run these before concluding any implementation step.
3. **Respect Arch & Style Guidelines:**
   - Prefer passing down existing references (handles) rather than using direct mutability or unsafe code.
   - Use `FeatureFlag` enums instead of raw `#[cfg]` strings unless absolutely necessary.
   - Minimize locking the `TerminalModel` as it can cause deadlocks on macOS.
4. **Don't hardcode secrets:** Ensure all keys are routed through environment variables.
5. **No unnecessary comments:** Do not document obvious code. Write comments only for non-obvious architecture or business logic.

See model-specific files for particular prompting adjustments if necessary.
