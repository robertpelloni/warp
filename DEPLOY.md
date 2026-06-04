# DEPLOY: Deployment & Environment Setup

## Development Environment
- **Language:** Go 1.21+
- **Platform-specific Dependencies:**
  - macOS: Xcode tools
  - Linux: x11/wayland dev libs
  - Windows: C++ build tools (for some cgo components)

## Setup
1. Clone the repository.
2. Navigate to `rebuild/`.
3. Run `go mod tidy`.
4. Run `go build ./cmd/warp`.

## Deployment
- Binaries are built for each target platform and bundled with necessary assets.
- Warp Drive sync requires connection to the Warp backend (or a self-hosted equivalent in the future).
