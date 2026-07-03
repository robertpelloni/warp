# Deployment & Environment Setup

## Prerequisites
- Go 1.25.0 or higher.
- `golang.org/x/sys` dependency.

## Building the Application
To build the Warp Go binary, run the following command from the repository root:

```bash
go build -o warp-go ./cmd/warp
```

### Windows
On Windows, this will compile a native Win32 executable leveraging the `pkg/win32` and `pkg/app/app_windows.go` files. It will utilize ConPTY for terminal emulation. No external C compiler (CGO) is required, making the build process entirely self-contained within the Go toolchain.

### macOS / Linux (UNIX)
On UNIX-like systems, the build will succeed but compile the stub implementation (`pkg/app/app_unix.go`). Running the resulting binary will simply print a message indicating that the GUI is currently Windows-only until a cross-platform framework is integrated. The underlying terminal logic (`pkg/pty_unix.go`), however, is fully functional and tested.

## Testing
To run the cross-platform test suite:

```bash
go test ./...
```
