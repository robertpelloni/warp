package dev.warp;

import java.util.HashMap;
import java.util.Map;
import java.util.Scanner;

// --- MCP Protocol Definitions ---

class JsonRpcRequest {
    public String id;
    public String method;
    public String params;

    public JsonRpcRequest(String id, String method, String params) {
        this.id = id;
        this.method = method;
        this.params = params;
    }
}

class JsonRpcResponse {
    public String id;
    public String result;
    public String error;

    public JsonRpcResponse(String id, String result, String error) {
        this.id = id;
        this.result = result;
        this.error = error;
    }
}

// --- Orchestration and Context ---

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
    private final Map<String, String> mcpTools = new HashMap<>();

    public synchronized void registerMcpTool(String name, String description) {
        mcpTools.put(name, description);
        System.out.printf("[Orchestrator] Registered dynamic MCP tool: %s%n", name);
    }

    public <Req, Out> Out executeTool(ToolRuntime<Req, Out> runtime, Req req, ToolCtx ctx) throws Exception {
        if (runtime.requiresApproval(req)) {
            System.out.printf("[Orchestrator] Requesting tool approval for '%s'...%n", ctx.toolName);
        }

        SandboxAttempt attempt = new SandboxAttempt(
                runtime.getSandboxPreference() == SandboxPreference.Sandboxed,
                ctx.turnContext.workingDir
        );

        return runtime.run(req, attempt, ctx);
    }
}

// --- MCP Server Implementation ---

class MessageProcessor {
    private ToolOrchestrator orchestrator;

    public MessageProcessor(ToolOrchestrator orchestrator) {
        this.orchestrator = orchestrator;
    }

    public JsonRpcResponse processRequest(JsonRpcRequest req) {
        System.out.printf("[MCP Server] Processing JSON-RPC method: %s%n", req.method);

        switch (req.method) {
            case "initialize":
                orchestrator.registerMcpTool("mcp_shell", "Execute commands via MCP");
                return new JsonRpcResponse(req.id, "initialized", null);
            case "tools/call":
                System.out.printf("[MCP Server] Dispatched to tool execution via orchestrator: %s%n", req.params);
                return new JsonRpcResponse(req.id, "Executed MCP tool call with args: " + req.params, null);
            default:
                return new JsonRpcResponse(req.id, null, "Method not found");
        }
    }
}

// --- Agent Implementation ---

class AgentSession {
    public String sessionId;
    public AgentStatus status;
    private ToolOrchestrator orchestrator;
    private MessageProcessor mcpProcessor;
    private int turnCounter;

    public AgentSession(String sessionId) {
        this.sessionId = sessionId;
        this.status = AgentStatus.PendingInit;
        this.orchestrator = new ToolOrchestrator();
        this.mcpProcessor = new MessageProcessor(orchestrator);
        this.turnCounter = 0;
    }

    public void initializeMcp() {
        JsonRpcRequest req = new JsonRpcRequest("0", "initialize", null);
        mcpProcessor.processRequest(req);
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

        if (input.startsWith("mcp")) {
            JsonRpcRequest req = new JsonRpcRequest("1", "tools/call", input);
            JsonRpcResponse resp = mcpProcessor.processRequest(req);
            System.out.println("[Agent] MCP Turn executed. Result: " + resp.result);
            this.status = AgentStatus.Completed;
            return;
        }

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
    public static void triggerPiMonoDemo() {

        SessionManager manager = new SessionManager("sess_pi_456", "/workspace", "/tmp/sessions");

        manager.appendMessageEntry("1", "", "Initial prompt");

        manager.appendMessageEntry("2", "1", "Assistant response");

        System.out.printf("[SessionManager] Tracked entries: %d%n", manager.fileEntries.size());

    }

    public static void triggerCompactionDemo() {

        try {

            SessionManager manager = new SessionManager("sess_pi_456", "/workspace", "/tmp/sessions");

            manager.appendMessageEntry("1", "", "Initial prompt");

            manager.appendMessageEntry("2", "1", "Assistant response");

            manager.appendMessageEntry("3", "2", "Please refactor the function.");

            CompactionPreparation prep = CompactionService.prepareCompaction(manager.fileEntries);

            CompactionResult result = CompactionService.compact(prep);

            System.out.printf("[Compaction] Successfully generated summary entry %s: \"%s\"%n", result.newEntryId, result.summary);

        } catch (Exception e) {

            System.out.println("[Error] " + e.getMessage());

        }

    }

    public static void triggerBrowserDemo() {

        BrowserConfig config = new BrowserConfig(

                true,

                new ViewportConfig(1280, 720, 1.0)

        );

        BrowserManager manager = new BrowserManager(config);

        manager.launch();

        if (manager.activePage != null) {

            manager.activePage.navigate("https://github.com/just-every/code");

            manager.activePage.dispatchMouseEvent(250.0, 300.0, "click");

            manager.activePage.dispatchKeyEvent("Hello Warp!");

            manager.activePage.captureScreenshot();

        }

    }
    public static void main(String[] args) {
        System.out.println("Welcome to Warp CLI (Java Edition) - Inspired by just-every-code");
        System.out.println("Type '/help' for commands, or 'quit' to close.");

        AgentSession agent = new AgentSession("sess_123");
        agent.initializeMcp();

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
                case "compact":
                    triggerCompactionDemo();
                    break;
                case "pimono":
                    triggerPiMonoDemo();
                    break;
                case "browser":
                    triggerBrowserDemo();
                    break;
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

// --- Browser Integration Abstractions ---

class ViewportConfig {
    public int width;
    public int height;
    public double deviceScaleFactor;

    public ViewportConfig(int width, int height, double deviceScaleFactor) {
        this.width = width;
        this.height = height;
        this.deviceScaleFactor = deviceScaleFactor;
    }
}

class BrowserConfig {
    public boolean headless;
    public ViewportConfig viewport;

    public BrowserConfig(boolean headless, ViewportConfig viewport) {
        this.headless = headless;
        this.viewport = viewport;
    }
}

class CursorState {
    public double x;
    public double y;

    public CursorState(double x, double y) {
        this.x = x;
        this.y = y;
    }
}

class CdpPage {
    public String url = "about:blank";
    public CursorState cursorState = new CursorState(0.0, 0.0);

    public void navigate(String newUrl) {
        System.out.printf("[Browser:Page] Navigating to: %s%n", newUrl);
        this.url = newUrl;
        try {
            Thread.sleep(600); // Simulate load time
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        System.out.printf("[Browser:Page] Successfully loaded %s%n", newUrl);
    }

    public synchronized void dispatchMouseEvent(double x, double y, String eventType) {
        System.out.printf("[Browser:Page] Dispatching mouse %s at (%f, %f)%n", eventType, x, y);
        this.cursorState.x = x;
        this.cursorState.y = y;
    }

    public void dispatchKeyEvent(String text) {
        System.out.printf("[Browser:Page] Dispatching keystroke: '%s'%n", text);
    }

    public String captureScreenshot() {
        System.out.printf("[Browser:Page] Capturing screenshot for %s%n", this.url);
        return "screenshot_data_base64_for_" + this.url;
    }
}

class BrowserManager {
    public BrowserConfig config;
    public CdpPage activePage;

    public BrowserManager(BrowserConfig config) {
        this.config = config;
    }

    public synchronized void launch() {
        String mode = config.headless ? "Headless" : "Windowed";
        System.out.printf("[Browser:Manager] Launching %s Chrome via CDP...%n", mode);
        this.activePage = new CdpPage();
    }
}

// --- Pi-Mono Specific Abstractions ---

enum SessionEntryType {
    MESSAGE,
    MODEL_CHANGE,
    COMPACTION,
    BRANCH_SUMMARY
}

class SessionEntryBase {
    public SessionEntryType type;
    public String id;
    public String parentId;
    public String timestamp;

    public SessionEntryBase(SessionEntryType type, String id, String parentId, String timestamp) {
        this.type = type;
        this.id = id;
        this.parentId = parentId;
        this.timestamp = timestamp;
    }
}

class SessionMessageEntry extends SessionEntryBase {
    public String message;

    public SessionMessageEntry(String id, String parentId, String timestamp, String message) {
        super(SessionEntryType.MESSAGE, id, parentId, timestamp);
        this.message = message;
    }
}

class SessionManager {
    public String sessionId;
    public String cwd;
    public String sessionDir;
    public java.util.List<Object> fileEntries;

    public SessionManager(String sessionId, String cwd, String sessionDir) {
        this.sessionId = sessionId;
        this.cwd = cwd;
        this.sessionDir = sessionDir;
        this.fileEntries = new java.util.ArrayList<>();
    }

    public void appendMessageEntry(String id, String parentId, String message) {
        SessionMessageEntry entry = new SessionMessageEntry(id, parentId, "now", message);
        this.fileEntries.add(entry);
        System.out.printf("[SessionManager] Appended Message: %s%n", message);
    }
}

// --- Context Compaction Abstractions ---

class FileOperations {
    public java.util.Set<String> read = new java.util.HashSet<>();
    public java.util.Set<String> written = new java.util.HashSet<>();
    public java.util.Set<String> edited = new java.util.HashSet<>();
}

class CompactionPreparation {
    public String firstKeptEntryId;
    public java.util.List<String> messagesToSummarize;
    public java.util.List<String> turnPrefixMessages;
    public boolean isSplitTurn;
    public int tokensBefore;
    public String previousSummary;
    public FileOperations fileOps;

    public CompactionPreparation() {
        this.messagesToSummarize = new java.util.ArrayList<>();
        this.turnPrefixMessages = new java.util.ArrayList<>();
        this.fileOps = new FileOperations();
    }
}

class CompactionResult {
    public String summary;
    public String newEntryId;
}

class CompactionService {
    public static CompactionPreparation prepareCompaction(java.util.List<Object> entries) throws Exception {
        System.out.printf("[Compaction] Analyzing %d session entries for compaction...%n", entries.size());

        if (entries.isEmpty()) {
            throw new Exception("no entries to compact");
        }

        CompactionPreparation prep = new CompactionPreparation();
        prep.firstKeptEntryId = "entry_xyz";
        prep.messagesToSummarize.add("User asked to build feature");
        prep.messagesToSummarize.add("Assistant planned feature");
        prep.tokensBefore = 4500;
        prep.previousSummary = "Previous session context...";

        prep.fileOps.read.add("/workspace/main.go");
        prep.fileOps.edited.add("/workspace/README.md");

        return prep;
    }

    public static CompactionResult compact(CompactionPreparation prep) {
        System.out.printf("[Compaction] Compacting %d tokens down to summary...%n", prep.tokensBefore);

        String summary = String.format("Summarized %d messages including: %s",
                prep.messagesToSummarize.size(),
                prep.messagesToSummarize.get(0));

        CompactionResult result = new CompactionResult();
        result.summary = summary;
        result.newEntryId = "compact_abc123";

        return result;
    }
}
