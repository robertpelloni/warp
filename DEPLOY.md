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
