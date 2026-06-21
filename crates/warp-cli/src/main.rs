use std::io::{self, Write};
use std::thread;
use std::time::Duration;

struct ToolOrchestrator {
    sandbox_enabled: bool,
}

impl ToolOrchestrator {
    fn new() -> Self {
        Self { sandbox_enabled: true }
    }

    fn execute_task(&self, task: &str) -> String {
        println!("[Orchestrator] Requesting approval for task: '{}'", task);
        // Simulate approval/execution loop
        thread::sleep(Duration::from_millis(500));

        if self.sandbox_enabled {
            println!("[Orchestrator] Running in sandbox mode...");
        }

        // Simulate multi-agent steps
        let steps = vec!["Plan", "Code", "Review"];
        for step in steps {
            println!("[Agent: {}] Processing...", step);
            thread::sleep(Duration::from_millis(400));
        }

        format!("Task '{}' completed successfully by Auto Drive.", task)
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
                println!("  /auto <task> - Start Auto Drive orchestration for a task");
                println!("  quit         - Quit the application");
            }
            "auto" => {
                if args.is_empty() {
                    println!("[Error] /auto requires a task description.");
                } else {
                    println!("[Orchestrator: Auto Drive] Starting autonomous loop.");
                    let result = orchestrator.execute_task(args);
                    println!("[Result] {}", result);
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
