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
```
