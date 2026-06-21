use std::io::{self, Write};

fn main() {
    println!("Welcome to Warp CLI (Rust Edition) - Inspired by just-every-code");
    println!("Type '/help' for commands, or 'quit' to close.");

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

        handle_command(trimmed);
    }
}

fn handle_command(input: &str) {
    if input.starts_with('/') {
        let cmd = &input[1..].to_lowercase();
        match cmd.as_str() {
            "help" => {
                println!("Available commands:");
                println!("  /help     - Show this help message");
                println!("  /plan     - (Stub) Coordinate planning agent");
                println!("  /code     - (Stub) Coordinate coding agent");
                println!("  /auto     - (Stub) Start Auto Drive orchestration");
                println!("  quit      - Quit the application");
            }
            "plan" => {
                println!("[Agent: Planner] Acknowledged. Ready to plan task.");
            }
            "code" => {
                println!("[Agent: Coder] Acknowledged. Ready to write code.");
            }
            "auto" => {
                println!("[Orchestrator: Auto Drive] Starting autonomous loop (Mock).");
            }
            _ => {
                println!("Unknown command: /{}", cmd);
            }
        }
    } else {
        println!("[Agent] Echoing input: {}", input);
    }
}
