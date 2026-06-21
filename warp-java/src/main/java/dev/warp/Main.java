package dev.warp;

import java.util.Scanner;

enum SandboxPreference {
    Unsandboxed,
    Sandboxed
}

enum AgentStatus {
    PendingInit,
    Running,
    Completed,
    Interrupted,
    Errored,
    Shutdown
}

class TurnContext {
    public String turnId;
    public String sessionId;
    public String model;
    public String workingDir;
    public String permissionsProfile;

    public TurnContext(String turnId, String sessionId, String model, String workingDir, String permissionsProfile) {
        this.turnId = turnId;
        this.sessionId = sessionId;
        this.model = model;
        this.workingDir = workingDir;
        this.permissionsProfile = permissionsProfile;
    }
}

class ToolCtx {
    public String callId;
    public String toolName;
    public TurnContext turnContext;

    public ToolCtx(String callId, String toolName, TurnContext turnContext) {
        this.callId = callId;
        this.toolName = toolName;
        this.turnContext = turnContext;
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
        return String.format("[%s] Executed shell command '%s' in %s mode (CallID: %s, TurnID: %s)",
                ctx.toolName, req, mode, ctx.callId, ctx.turnContext.turnId);
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
                ctx.turnContext.workingDir
        );

        // 3. Execution
        return runtime.run(req, attempt, ctx);
    }
}

class AgentSession {
    public String sessionId;
    public AgentStatus status;
    private ToolOrchestrator orchestrator;
    private int turnCounter;

    public AgentSession(String sessionId) {
        this.sessionId = sessionId;
        this.status = AgentStatus.PendingInit;
        this.orchestrator = new ToolOrchestrator();
        this.turnCounter = 0;
    }

    public void steerInput(String input) {
        this.status = AgentStatus.Running;
        this.turnCounter++;
        String turnId = "turn_" + this.turnCounter;

        System.out.printf("[Agent] Received input: '%s'. Generating TurnContext (%s)...%n", input, turnId);

        TurnContext turnContext = new TurnContext(
                turnId,
                this.sessionId,
                "gpt-5.5",
                "/workspace",
                "default"
        );

        String toolReq = "echo '" + input + "'";
        ToolCtx ctx = new ToolCtx("call_" + this.turnCounter, "shell", turnContext);
        ShellCommandRuntime runtime = new ShellCommandRuntime();

        try {
            String result = orchestrator.executeTool(runtime, toolReq, ctx);
            System.out.println("[Agent] Turn executed. Result: " + result);
            this.status = AgentStatus.Completed;
        } catch (Exception e) {
            System.out.println("[Agent] Turn failed. Error: " + e.getMessage());
            this.status = AgentStatus.Errored;
        }
    }
}

public class Main {
    public static void main(String[] args) {
        System.out.println("Welcome to Warp CLI (Java Edition) - Inspired by just-every-code");
        System.out.println("Type '/help' for commands, or 'quit' to close.");

        AgentSession agent = new AgentSession("sess_123");
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
                agent.status = AgentStatus.Shutdown;
                break;
            }

            handleCommand(input, agent);
        }
        scanner.close();
    }

    private static void handleCommand(String input, AgentSession agent) {
        if (input.startsWith("/")) {
            String[] parts = input.substring(1).split(" ", 2);
            String cmd = parts[0].toLowerCase();
            String args = parts.length > 1 ? parts[1].trim() : "";

            switch (cmd) {
                case "help":
                    System.out.println("Available commands:");
                    System.out.println("  /help        - Show this help message");
                    System.out.println("  /prompt <msg>- Send a natural language prompt to the agent");
                    System.out.println("  quit         - Quit the application");
                    break;
                case "prompt":
                    if (args.isEmpty()) {
                        System.out.println("[Error] /prompt requires a message.");
                    } else {
                        agent.steerInput(args);
                    }
                    break;
                default:
                    System.out.println("Unknown command: /" + cmd);
                    break;
            }
        } else {
            agent.steerInput(input);
        }
    }
}
