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
