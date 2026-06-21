package dev.warp;

import java.util.Scanner;

enum SandboxPreference {
    Unsandboxed,
    Sandboxed
}

class ToolCtx {
    public String callId;
    public String toolName;
    public String sessionId;
    public String turnId;

    public ToolCtx(String callId, String toolName, String sessionId, String turnId) {
        this.callId = callId;
        this.toolName = toolName;
        this.sessionId = sessionId;
        this.turnId = turnId;
    }
}

class SandboxAttempt {
    public boolean isSandboxed;
    public String cwd;

    public SandboxAttempt(boolean isSandboxed, String cwd) {
        this.isSandboxed = isSandboxed;
        this.cwd = cwd;
    }
}

interface ToolRuntime<Req, Out> {
    SandboxPreference getSandboxPreference();
    boolean isEscalateOnFailure();
    boolean requiresApproval(Req req);
    Out run(Req req, SandboxAttempt attempt, ToolCtx ctx) throws Exception;
}

class ShellCommandRuntime implements ToolRuntime<String, String> {
    @Override
    public SandboxPreference getSandboxPreference() {
        return SandboxPreference.Sandboxed;
    }

    @Override
    public boolean isEscalateOnFailure() {
        return true;
    }

    @Override
    public boolean requiresApproval(String req) {
        return true;
    }

    @Override
    public String run(String req, SandboxAttempt attempt, ToolCtx ctx) {
        String mode = attempt.isSandboxed ? "Sandboxed" : "Unsandboxed";
        return String.format("[%s] Executed shell command '%s' in %s mode (CallID: %s)",
                ctx.toolName, req, mode, ctx.callId);
    }
}

class ToolOrchestrator {
    public <Req, Out> Out executeTool(ToolRuntime<Req, Out> runtime, Req req, ToolCtx ctx) throws Exception {
        // 1. Approval Phase
        if (runtime.requiresApproval(req)) {
            System.out.printf("[Orchestrator] Requesting tool approval for '%s'...%n", ctx.toolName);
            // Assuming approved
        }

        // 2. Sandbox Selection
        SandboxAttempt attempt = new SandboxAttempt(
                runtime.getSandboxPreference() == SandboxPreference.Sandboxed,
                "/workspace"
        );

        // 3. Execution
        return runtime.run(req, attempt, ctx);
    }
}

public class Main {
    public static void main(String[] args) {
        System.out.println("Welcome to Warp CLI (Java Edition) - Inspired by just-every-code");
        System.out.println("Type '/help' for commands, or 'quit' to close.");

        ToolOrchestrator orchestrator = new ToolOrchestrator();
        Scanner scanner = new Scanner(System.in);

        while (true) {
            System.out.print("warp> ");
            if (!scanner.hasNextLine()) {
                break;
            }

            String input = scanner.nextLine().trim();
            if (input.isEmpty()) {
                continue;
            }

            if (input.equalsIgnoreCase("quit") || input.equalsIgnoreCase("/quit")) {
                break;
            }

            handleCommand(input, orchestrator);
        }
        scanner.close();
    }

    private static void handleCommand(String input, ToolOrchestrator orchestrator) {
        if (input.startsWith("/")) {
            String[] parts = input.substring(1).split(" ", 2);
            String cmd = parts[0].toLowerCase();
            String args = parts.length > 1 ? parts[1].trim() : "";

            switch (cmd) {
                case "help":
                    System.out.println("Available commands:");
                    System.out.println("  /help        - Show this help message");
                    System.out.println("  /shell <cmd> - Run a command through the ToolOrchestrator");
                    System.out.println("  quit         - Quit the application");
                    break;
                case "shell":
                    if (args.isEmpty()) {
                        System.out.println("[Error] /shell requires a command.");
                    } else {
                        ToolCtx ctx = new ToolCtx("call_abc123", "shell", "sess_1", "turn_1");
                        ShellCommandRuntime runtime = new ShellCommandRuntime();

                        try {
                            String result = orchestrator.executeTool(runtime, args, ctx);
                            System.out.println("[Result] " + result);
                        } catch (Exception e) {
                            System.out.println("[Error] " + e.getMessage());
                        }
                    }
                    break;
                default:
                    System.out.println("Unknown command: /" + cmd);
                    break;
            }
        } else {
            System.out.println("[Agent] Echoing input: " + input);
        }
    }
}
