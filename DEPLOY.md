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
3. Ensure Go 1.24+ is installed.
4. Run `go mod tidy` to install dependencies.
5. Run `./script/build.sh` to build for all platforms, or `go build ./cmd/warp` for local.

## Deployment
### Staging
Run `./script/deploy_staging.sh` to generate release bundles in the `staging/` directory.

### Production
1. Build production binaries: `./script/build.sh`.
2. Verified binaries are available in `rebuild/build/`.
3. Distributed via automated CI/CD pipelines or manual upload of artifacts.
- Warp Drive sync requires connection to the Warp backend (or a self-hosted equivalent in the future).
