# Deployment & Setup

## Prerequisites
- Rust toolchain (`rust-toolchain.toml` specifies the version).
- Node.js (if required for frontend aspects) or cargo-build tools.
- A local Warp server for specific data features (optional).

## Building from source
```bash
./script/bootstrap   # installs platform dependencies
./script/run         # build and run Warp locally
```

## Environment Variables / Secrets
Warp's frontend codebase relies on some API connections. Do not hardcode secrets or keys.
Placeholder files for environments:
- `.env.example` should contain placeholders like `OPENAI_API_KEY=YOUR_KEY_HERE`.

For local development with the backend server:
```bash
SERVER_ROOT_URL=http://localhost:8080 WS_SERVER_URL=ws://localhost:8080/graphql/v2 cargo run --features with_local_server
```

## Testing and Verification
Run all necessary pre-submit checks before pushing code:
```bash
./script/presubmit
cargo nextest run --no-fail-fast --workspace --exclude command-signatures-v2

# Deployment & Environment Setup

## Environment Requirements
- **Go**: 1.25+ or 1.26+
- **Operating Systems**:
  - Windows 10/11 (AMD64/ARM64) for native Win32 GUI & ConPTY
  - Linux (Ubuntu/Debian, Fedora, Arch) for headless CLI / server runners
  - macOS (Darwin) for headless CLI / server runners

## Building Warp Go
### Local Windows Build
Run the provided batch automation script:
```cmd
start.bat build
start.bat run
```
Or directly using Go:
```cmd
go build -ldflags="-s -w -H=windowsgui" -o warp-go.exe ./cmd/warp
```

### Cross-Compilation for Windows (from Linux / macOS)
```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o warp-go.exe ./cmd/warp
```

### Running Presubmit Verification
```bash
./script/presubmit
```

### Testing
```bash
go test -v ./...
```

# Warp: Deployment & Environment Setup

## Rust (`crates/warp-cli`)
- Dependency Engine: `cargo`
- Build: `cd crates/warp-cli && cargo build --release`
- Run: `cd crates/warp-cli && cargo run`
- Testing: `cargo test`

## Go (`warp-go`)
- Target: Go 1.26.0 (or appropriate workspace module)
- Build: `cd warp-go && go build ./cmd/warp-cli`
- Run: `./warp-cli`

## C# (`warp-csharp`)
- Target: .NET 8.0 SDK
- Build: `cd warp-csharp/WarpCsharp && dotnet build`
- Run: `cd warp-csharp/WarpCsharp && dotnet run`

## Java (`warp-java`)
- Target: JDK 21+
- Build System: Maven
- Build: `cd warp-java && mvn clean package`
- Run: `java -cp target/classes dev.warp.Main`

# DEPLOY: Deployment & Environment Setup

## Development Environment
- **Language:** Go 1.24+
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
