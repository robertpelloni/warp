use std::io::{self, Write};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

// Core Abstractions Ported from codex-rs/core/src/tools/sandboxing.rs
// and codex-rs/core/src/session/turn_context.rs
// Now includes MCP integration structures derived from codex-rs/mcp-server/

pub mod browser_integration;
pub mod mcp_protocol {
    #[derive(Debug, Clone)]
    pub struct JsonRpcRequest {
        pub id: String,
        pub method: String,
        pub params: Option<String>,
    }

    #[derive(Debug, Clone)]
    pub struct JsonRpcResponse {
        pub id: String,
        pub result: Option<String>,
        pub error: Option<String>,
    }

    #[derive(Debug, Clone)]
    pub struct CallToolRequestParams {
        pub name: String,
        pub arguments: String,
    }

    #[derive(Debug, Clone)]
    pub struct CallToolResult {
        pub content: String,
        pub is_error: bool,
    }
}

use mcp_protocol::*;

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

struct ShellCommandRuntime;

impl Sandboxable for ShellCommandRuntime {
    fn sandbox_preference(&self) -> SandboxPreference {
        SandboxPreference::Unsandboxed
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
        if !attempt.is_sandboxed {
            println!("[ShellRuntime] Executing real command in unsandboxed mode: {}", req);
            let output = std::process::Command::new("sh")
                .arg("-c")
                .arg(req)
                .current_dir(&attempt.cwd)
                .output()
                .map_err(|e| ToolError::ExecutionError(e.to_string()))?;

            let stdout = String::from_utf8_lossy(&output.stdout);
            let stderr = String::from_utf8_lossy(&output.stderr);
            Ok(format!("{}\n{}", stdout, stderr).trim().to_string())
        } else {
            Ok(format!(
                "[{}] Executed shell command '{}' in Sandboxed mode (CallID: {}, TurnID: {})",
                ctx.tool_name, req, ctx.call_id, ctx.turn_context.turn_id
            ))
        }
    }
}

struct ToolOrchestrator {
    mcp_tools: Arc<Mutex<HashMap<String, String>>>, // registered MCP capabilities
}

impl ToolOrchestrator {
    fn new() -> Self {
        Self {
            mcp_tools: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    fn register_mcp_tool(&self, name: String, description: String) {
        let mut tools = self.mcp_tools.lock().unwrap();
        tools.insert(name.clone(), description);
        println!("[Orchestrator] Registered dynamic MCP tool: {}", name);
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
        let requires_approval = runtime.requires_approval(req);
        if requires_approval {
            println!("[Orchestrator] Requesting tool approval for '{}'...", ctx.tool_name);
        }

        let attempt = SandboxAttempt {
            is_sandboxed: runtime.sandbox_preference() == SandboxPreference::Unsandboxed,
            cwd: ctx.turn_context.working_dir.clone(),
        };

        runtime.run(req, &attempt, ctx)
    }
}

struct MessageProcessor {
    orchestrator: Arc<ToolOrchestrator>,
}

impl MessageProcessor {
    fn new(orchestrator: Arc<ToolOrchestrator>) -> Self {
        Self { orchestrator }
    }

    fn process_request(&self, req: JsonRpcRequest) -> JsonRpcResponse {
        println!("[MCP Server] Processing JSON-RPC method: {}", req.method);
        match req.method.as_str() {
            "initialize" => {
                self.orchestrator.register_mcp_tool("mcp_shell".into(), "Execute commands via MCP".into());
                JsonRpcResponse {
                    id: req.id,
                    result: Some("initialized".into()),
                    error: None,
                }
            }
            "tools/call" => {
                let params = req.params.unwrap_or_default();
                println!("[MCP Server] Dispatched to tool execution via orchestrator: {}", params);
                JsonRpcResponse {
                    id: req.id,
                    result: Some(format!("Executed MCP tool call with args: {}", params)),
                    error: None,
                }
            }
            _ => JsonRpcResponse {
                id: req.id,
                result: None,
                error: Some("Method not found".into()),
            }
        }
    }
}

struct AgentSession {
    session_id: String,
    pub status: AgentStatus,
    orchestrator: Arc<ToolOrchestrator>,
    mcp_processor: MessageProcessor,
    turn_counter: u32,
}

impl AgentSession {
    fn new(session_id: String) -> Self {
        let orchestrator = Arc::new(ToolOrchestrator::new());
        let mcp_processor = MessageProcessor::new(Arc::clone(&orchestrator));

        Self {
            session_id,
            status: AgentStatus::PendingInit,
            orchestrator,
            mcp_processor,
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

        if input.starts_with("mcp") {
            let req = JsonRpcRequest {
                id: "1".into(),
                method: "tools/call".into(),
                params: Some(input.to_string()),
            };
            let response = self.mcp_processor.process_request(req);
            println!("[Agent] MCP Turn executed. Result: {:?}", response.result.unwrap_or_default());
            self.status = AgentStatus::Completed("Success".into());
            return;
        }

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

    fn initialize_mcp(&self) {
        let req = JsonRpcRequest {
            id: "0".into(),
            method: "initialize".into(),
            params: None,
        };
        self.mcp_processor.process_request(req);
    }
}

fn main() {
    println!("Welcome to Warp CLI (Rust Edition) - Inspired by just-every-code");
    println!("Type '/help' for commands, or 'quit' to close.");

    let mut agent = AgentSession::new("sess_123".into());
    agent.initialize_mcp();

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
"browser" => { trigger_browser_demo(); }
"compact" => { trigger_compaction_demo(); }
"pimono" => { trigger_pimono_demo(); }
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
        agent.steer_input(input);
    }
}

// Testing browser abstraction inside the CLI REPL
fn trigger_browser_demo() {
    use std::time::Duration;
    let config = browser_integration::BrowserConfig {
        headless: true,
        viewport: browser_integration::ViewportConfig {
            width: 1280,
            height: 720,
            device_scale_factor: 1.0,
        },
        timeout: Duration::from_secs(30),
    };

    let manager = browser_integration::BrowserManager::new(config);
    if let Ok(_) = manager.launch() {
        if let Some(mut page) = manager.active_page.lock().unwrap().as_mut() {
            let _ = page.navigate("https://github.com/just-every/code");
            let _ = page.dispatch_mouse_event(250.0, 300.0, "click");
            let _ = page.dispatch_key_event("Hello Warp!");
            let _ = page.capture_screenshot();
        }
    }
}

// --- Pi-Mono Specific Abstractions ---

#[derive(Debug, Clone, PartialEq)]
pub enum SessionEntryType {
    Message,
    ModelChange,
    Compaction,
    BranchSummary,
}

#[derive(Debug, Clone)]
pub struct SessionEntryBase {
    pub entry_type: SessionEntryType,
    pub id: String,
    pub parent_id: String,
    pub timestamp: String,
}

#[derive(Debug, Clone)]
pub struct SessionMessageEntry {
    pub base: SessionEntryBase,
    pub message: String,
}

pub struct SessionManager {
    pub session_id: String,
    pub cwd: String,
    pub session_dir: String,
    pub file_entries: Vec<SessionMessageEntry>, // Simplified representation
}

impl SessionManager {
    pub fn new(session_id: &str, cwd: &str, session_dir: &str) -> Self {
        Self {
            session_id: session_id.to_string(),
            cwd: cwd.to_string(),
            session_dir: session_dir.to_string(),
            file_entries: Vec::new(),
        }
    }

    pub fn append_message_entry(&mut self, id: &str, parent_id: &str, message: &str) {
        let entry = SessionMessageEntry {
            base: SessionEntryBase {
                entry_type: SessionEntryType::Message,
                id: id.to_string(),
                parent_id: parent_id.to_string(),
                timestamp: "now".to_string(),
            },
            message: message.to_string(),
        };
        self.file_entries.push(entry);
        println!("[SessionManager] Appended Message: {}", message);
    }
}

// Testing Pi-Mono abstraction inside the CLI REPL
fn trigger_pimono_demo() {
    let mut manager = SessionManager::new("sess_pi_456", "/workspace", "/tmp/sessions");
    manager.append_message_entry("1", "", "Initial prompt");
    manager.append_message_entry("2", "1", "Assistant response");
    println!("[SessionManager] Tracked entries: {}", manager.file_entries.len());
}

// --- Context Compaction Abstractions ---

use std::collections::HashSet;

#[derive(Debug, Clone)]
pub struct FileOperations {
    pub read: HashSet<String>,
    pub written: HashSet<String>,
    pub edited: HashSet<String>,
}

impl FileOperations {
    pub fn new() -> Self {
        Self {
            read: HashSet::new(),
            written: HashSet::new(),
            edited: HashSet::new(),
        }
    }
}

#[derive(Debug, Clone)]
pub struct CompactionPreparation {
    pub first_kept_entry_id: String,
    pub messages_to_summarize: Vec<String>,
    pub turn_prefix_messages: Vec<String>,
    pub is_split_turn: bool,
    pub tokens_before: usize,
    pub previous_summary: String,
    pub file_ops: FileOperations,
}

#[derive(Debug, Clone)]
pub struct CompactionResult {
    pub summary: String,
    pub new_entry_id: String,
}

pub struct CompactionService;

impl CompactionService {
    pub fn prepare_compaction(entries: &[SessionMessageEntry]) -> Result<CompactionPreparation, String> {
        println!("[Compaction] Analyzing {} session entries for compaction...", entries.len());

        if entries.is_empty() {
            return Err("no entries to compact".into());
        }

        let mut prep = CompactionPreparation {
            first_kept_entry_id: "entry_xyz".into(),
            messages_to_summarize: vec!["User asked to build feature".into(), "Assistant planned feature".into()],
            turn_prefix_messages: Vec::new(),
            is_split_turn: false,
            tokens_before: 4500,
            previous_summary: "Previous session context...".into(),
            file_ops: FileOperations::new(),
        };

        prep.file_ops.read.insert("/workspace/main.go".into());
        prep.file_ops.edited.insert("/workspace/README.md".into());

        Ok(prep)
    }

    pub fn compact(prep: &CompactionPreparation) -> CompactionResult {
        println!("[Compaction] Compacting {} tokens down to summary...", prep.tokens_before);

        let summary = format!(
            "Summarized {} messages including: {}",
            prep.messages_to_summarize.len(),
            prep.messages_to_summarize.first().unwrap_or(&String::new())
        );

        CompactionResult {
            summary,
            new_entry_id: "compact_abc123".into(),
        }
    }
}

fn trigger_compaction_demo() {
    let mut manager = SessionManager::new("sess_pi_456", "/workspace", "/tmp/sessions");
    manager.append_message_entry("1", "", "Initial prompt");
    manager.append_message_entry("2", "1", "Assistant response");
    manager.append_message_entry("3", "2", "Please refactor the function.");

    match CompactionService::prepare_compaction(&manager.file_entries) {
        Ok(prep) => {
            let result = CompactionService::compact(&prep);
            println!("[Compaction] Successfully generated summary entry {}: '{}'", result.new_entry_id, result.summary);
        }
        Err(e) => println!("[Error] {}", e),
    }
}
