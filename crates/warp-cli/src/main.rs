use std::io::{self, Write};

// Core Abstractions Ported from codex-rs/core/src/tools/sandboxing.rs
// and codex-rs/core/src/session/turn_context.rs

#[derive(Debug, Clone, PartialEq)]
pub enum SandboxPreference {
    Unsandboxed,
    Sandboxed,
}

#[derive(Debug, Clone)]
pub enum AgentStatus {
    PendingInit,
    Running,
    Completed(String),
    Interrupted,
    Errored(String),
    Shutdown,
}

#[derive(Debug, Clone)]
pub struct TurnContext {
    pub turn_id: String,
    pub session_id: String,
    pub model: String,
    pub working_dir: String,
    pub permissions_profile: String,
}

#[derive(Debug, Clone)]
pub struct ToolCtx {
    pub call_id: String,
    pub tool_name: String,
    pub turn_context: TurnContext,
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
            "[{}] Executed shell command '{}' in {} mode (CallID: {}, TurnID: {})",
            ctx.tool_name, req, mode, ctx.call_id, ctx.turn_context.turn_id
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
            cwd: ctx.turn_context.working_dir.clone(),
        };

        // 3. Execution
        runtime.run(req, &attempt, ctx)
    }
}

struct AgentSession {
    session_id: String,
    status: AgentStatus,
    orchestrator: ToolOrchestrator,
    turn_counter: u32,
}

impl AgentSession {
    fn new(session_id: String) -> Self {
        Self {
            session_id,
            status: AgentStatus::PendingInit,
            orchestrator: ToolOrchestrator::new(),
            turn_counter: 0,
        }
    }

    fn steer_input(&mut self, input: &str) {
        self.status = AgentStatus::Running;
        self.turn_counter += 1;
        let turn_id = format!("turn_{}", self.turn_counter);

        println!("[Agent] Received input: '{}'. Generating TurnContext ({})...", input, turn_id);

        let turn_context = TurnContext {
            turn_id: turn_id.clone(),
            session_id: self.session_id.clone(),
            model: "gpt-5.5".into(),
            working_dir: "/workspace".into(),
            permissions_profile: "default".into(),
        };

        // For simulation purposes, we map natural language to a shell tool call.
        let tool_req = format!("echo '{}'", input);
        let ctx = ToolCtx {
            call_id: format!("call_{}", self.turn_counter),
            tool_name: "shell".into(),
            turn_context: turn_context.clone(),
        };

        let runtime = ShellCommandRuntime;
        match self.orchestrator.execute_tool(runtime, &tool_req, &ctx) {
            Ok(out) => {
                println!("[Agent] Turn executed. Result: {}", out);
                self.status = AgentStatus::Completed("Success".into());
            }
            Err(e) => {
                println!("[Agent] Turn failed. Error: {:?}", e);
                self.status = AgentStatus::Errored("Execution Failed".into());
            }
        }
    }
}

fn main() {
    println!("Welcome to Warp CLI (Rust Edition) - Inspired by just-every-code");
    println!("Type '/help' for commands, or 'quit' to close.");

    let mut agent = AgentSession::new("sess_123".into());
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
            agent.status = AgentStatus::Shutdown;
            break;
        }

        handle_command(trimmed, &mut agent);
    }
}

fn handle_command(input: &str, agent: &mut AgentSession) {
    if input.starts_with('/') {
        let cmd_parts: Vec<&str> = input[1..].splitn(2, ' ').collect();
        let cmd = cmd_parts[0].to_lowercase();
        let args = if cmd_parts.len() > 1 { cmd_parts[1] } else { "" };

        match cmd.as_str() {
            "help" => {
                println!("Available commands:");
                println!("  /help        - Show this help message");
                println!("  /prompt <msg>- Send a natural language prompt to the agent");
                println!("  quit         - Quit the application");
            }
            "prompt" => {
                if args.is_empty() {
                    println!("[Error] /prompt requires a message.");
                } else {
                    agent.steer_input(args);
                }
            }
            _ => {
                println!("Unknown command: /{}", cmd);
            }
        }
    } else {
        // Default treat as prompt
        agent.steer_input(input);
    }
}
