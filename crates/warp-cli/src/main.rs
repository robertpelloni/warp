use std::io::{self, Write};



// Core Abstractions Ported from codex-rs/core/src/tools/sandboxing.rs

#[derive(Debug, Clone, PartialEq)]
pub enum SandboxPreference {
    Unsandboxed,
    Sandboxed,
}

#[derive(Debug, Clone)]
pub struct ToolCtx {
    pub call_id: String,
    pub tool_name: String,
    pub session_id: String,
    pub turn_id: String,
}

#[derive(Debug)]
pub enum ToolError {
    Rejected(String),
    ExecutionError(String),
}

pub struct SandboxAttempt {
    pub is_sandboxed: bool,
    pub cwd: String,
}

pub trait Sandboxable {
    fn sandbox_preference(&self) -> SandboxPreference;
    fn escalate_on_failure(&self) -> bool { true }
}

pub trait Approvable<Req> {
    type ApprovalKey;
    fn approval_keys(&self, req: &Req) -> Vec<Self::ApprovalKey>;
    fn requires_approval(&self, _req: &Req) -> bool { true }
}

pub trait ToolRuntime<Req, Out>: Approvable<Req> + Sandboxable {
    fn run(
        &self,
        req: &Req,
        attempt: &SandboxAttempt,
        ctx: &ToolCtx,
    ) -> Result<Out, ToolError>;
}

// Example Implementations

struct ShellCommandRuntime;

impl Sandboxable for ShellCommandRuntime {
    fn sandbox_preference(&self) -> SandboxPreference {
        SandboxPreference::Sandboxed
    }
}

impl Approvable<String> for ShellCommandRuntime {
    type ApprovalKey = String;
    fn approval_keys(&self, req: &String) -> Vec<Self::ApprovalKey> {
        vec![req.clone()]
    }
}

impl ToolRuntime<String, String> for ShellCommandRuntime {
    fn run(
        &self,
        req: &String,
        attempt: &SandboxAttempt,
        ctx: &ToolCtx,
    ) -> Result<String, ToolError> {
        let mode = if attempt.is_sandboxed { "Sandboxed" } else { "Unsandboxed" };
        Ok(format!(
            "[{}] Executed shell command '{}' in {} mode (CallID: {})",
            ctx.tool_name, req, mode, ctx.call_id
        ))
    }
}

struct ToolOrchestrator;

impl ToolOrchestrator {
    fn new() -> Self {
        Self
    }

    fn execute_tool<R, Req, Out>(
        &self,
        runtime: R,
        req: &Req,
        ctx: &ToolCtx,
    ) -> Result<Out, ToolError>
    where
        R: ToolRuntime<Req, Out>,
    {
        // 1. Approval Phase
        let requires_approval = runtime.requires_approval(req);
        if requires_approval {
            println!("[Orchestrator] Requesting tool approval for '{}'...", ctx.tool_name);
            // Assuming approved for now
        }

        // 2. Sandbox Selection
        let attempt = SandboxAttempt {
            is_sandboxed: runtime.sandbox_preference() == SandboxPreference::Sandboxed,
            cwd: "/workspace".to_string(),
        };

        // 3. Execution
        runtime.run(req, &attempt, ctx)
    }
}

fn main() {
    println!("Welcome to Warp CLI (Rust Edition) - Inspired by just-every-code");
    println!("Type '/help' for commands, or 'quit' to close.");

    let orchestrator = ToolOrchestrator::new();
    let mut input = String::new();

    loop {
        print!("warp> ");
        io::stdout().flush().unwrap();

        input.clear();
        if io::stdin().read_line(&mut input).unwrap() == 0 {
            break; // EOF
        }

        let trimmed = input.trim();
        if trimmed.is_empty() {
            continue;
        }

        if trimmed.eq_ignore_ascii_case("quit") || trimmed.eq_ignore_ascii_case("/quit") {
            break;
        }

        handle_command(trimmed, &orchestrator);
    }
}

fn handle_command(input: &str, orchestrator: &ToolOrchestrator) {
    if input.starts_with('/') {
        let cmd_parts: Vec<&str> = input[1..].splitn(2, ' ').collect();
        let cmd = cmd_parts[0].to_lowercase();
        let args = if cmd_parts.len() > 1 { cmd_parts[1] } else { "" };

        match cmd.as_str() {
            "help" => {
                println!("Available commands:");
                println!("  /help        - Show this help message");
                println!("  /shell <cmd> - Run a command through the ToolOrchestrator");
                println!("  quit         - Quit the application");
            }
            "shell" => {
                if args.is_empty() {
                    println!("[Error] /shell requires a command.");
                } else {
                    let ctx = ToolCtx {
                        call_id: "call_abc123".into(),
                        tool_name: "shell".into(),
                        session_id: "sess_1".into(),
                        turn_id: "turn_1".into(),
                    };

                    let runtime = ShellCommandRuntime;
                    let req = args.to_string();

                    match orchestrator.execute_tool(runtime, &req, &ctx) {
                        Ok(out) => println!("[Result] {}", out),
                        Err(e) => println!("[Error] {:?}", e),
                    }
                }
            }
            _ => {
                println!("Unknown command: /{}", cmd);
            }
        }
    } else {
        println!("[Agent] Echoing input: {}", input);
    }
}
