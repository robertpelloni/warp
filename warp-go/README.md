# Warp Go Backend Server

The `warp-go` directory contains the core backend server orchestrator for the Warp Agent ecosystem. It serves as the bridge between client applications (like the React WebUI) and the underlying autonomous agents and filesystem tooling.

## Architecture

The server exposes a series of HTTP REST endpoints, a WebSocket interface for streaming, and a JSON-RPC interface for browser extension control.

### Key Endpoints:
- `GET /status`: Health check endpoint returning JSON server metadata.
- `GET /ws`: WebSocket endpoint for real-time bi-directional streaming.
- `POST /rpc`: JSON-RPC 2.0 compliant endpoint for invoking specific agent commands.
- `POST /agent`: Direct agent messaging endpoint.
- `POST /fs`: Local filesystem reading utility.
- `POST /git`: Automated Git command execution wrapper.
- `POST /ast`: Deep Static Dependency Tree (AST) analysis triggers.
- `POST /cdp`: Chromium DevTools Protocol (CDP) triggers for network activity scraping.

## Prerequisites
- Go 1.20 or newer.

## Installation & Setup

1. **Clone and Initialize:** Ensure you are in the `warp-go` directory. Run the module tidy to pull any dependencies.
   ```bash
   go mod tidy
   ```

2. **Build:** You can build the binary output to the current directory.
   ```bash
   go build -o warp-server ./cmd/warp-server/...
   ```

3. **Run Server:** Start the server on port `8080`.
   ```bash
   PORT=8080 ./warp-server
   ```
   Or via `go run`:
   ```bash
   PORT=8080 go run ./cmd/warp-server/...
   ```

## Development and Contribution

- All HTTP route handlers are configured inside `cmd/warp-server/main.go`.
- If introducing new tools, they should be broken out into dedicated Go files (e.g. `fs.go`, `git.go`) keeping the mux routing clean.
- Ensure cross-origin resource sharing (CORS) remains configured correctly via the `corsMiddleware` function in `main.go` to support WebUI connectivity.
