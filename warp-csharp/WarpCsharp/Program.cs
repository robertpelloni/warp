using System;
using System.Collections.Generic;

namespace WarpCsharp
{
    // --- MCP Protocol Definitions ---

    public class JsonRpcRequest
    {
        public string Id { get; set; } = string.Empty;
        public string Method { get; set; } = string.Empty;
        public string Params { get; set; } = string.Empty;
    }

    public class JsonRpcResponse
    {
        public string Id { get; set; } = string.Empty;
        public string Result { get; set; } = string.Empty;
        public string Error { get; set; } = string.Empty;
    }

    // --- Orchestration and Context ---

    public enum SandboxPreference
    {
        Unsandboxed,
        Sandboxed
    }

    public enum AgentStatus
    {
        PendingInit,
        Running,
        Completed,
        Interrupted,
        Errored,
        Shutdown
    }

    public class TurnContext
    {
        public string TurnId { get; set; } = string.Empty;
        public string SessionId { get; set; } = string.Empty;
        public string Model { get; set; } = string.Empty;
        public string WorkingDir { get; set; } = string.Empty;
        public string PermissionsProfile { get; set; } = string.Empty;
    }

    public class ToolCtx
    {
        public string CallId { get; set; } = string.Empty;
        public string ToolName { get; set; } = string.Empty;
        public TurnContext TurnContext { get; set; } = new TurnContext();
    }

    public class SandboxAttempt
    {
        public bool IsSandboxed { get; set; }
        public string Cwd { get; set; } = string.Empty;
    }

    public interface IToolRuntime<TReq, TOut>
    {
        SandboxPreference SandboxPreference { get; }
        bool EscalateOnFailure { get; }
        bool RequiresApproval(TReq req);
        TOut Run(TReq req, SandboxAttempt attempt, ToolCtx ctx);
    }

    public class ShellCommandRuntime : IToolRuntime<string, string>
    {
        public SandboxPreference SandboxPreference => SandboxPreference.Sandboxed;

        public bool EscalateOnFailure => true;

        public bool RequiresApproval(string req) => true;

        public string Run(string req, SandboxAttempt attempt, ToolCtx ctx)
        {
            string mode = attempt.IsSandboxed ? "Sandboxed" : "Unsandboxed";
            return $"[{ctx.ToolName}] Executed shell command '{req}' in {mode} mode (CallID: {ctx.CallId}, TurnID: {ctx.TurnContext.TurnId})";
        }
    }

    public class ToolOrchestrator
    {
        private Dictionary<string, string> _mcpTools = new Dictionary<string, string>();
        private readonly object _lock = new object();

        public void RegisterMcpTool(string name, string description)
        {
            lock (_lock)
            {
                _mcpTools[name] = description;
                Console.WriteLine($"[Orchestrator] Registered dynamic MCP tool: {name}");
            }
        }

        public TOut ExecuteTool<TReq, TOut>(IToolRuntime<TReq, TOut> runtime, TReq req, ToolCtx ctx)
        {
            if (runtime.RequiresApproval(req))
            {
                Console.WriteLine($"[Orchestrator] Requesting tool approval for '{ctx.ToolName}'...");
            }

            var attempt = new SandboxAttempt
            {
                IsSandboxed = runtime.SandboxPreference == SandboxPreference.Sandboxed,
                Cwd = ctx.TurnContext.WorkingDir
            };

            return runtime.Run(req, attempt, ctx);
        }
    }

    // --- MCP Server Implementation ---

    public class MessageProcessor
    {
        private ToolOrchestrator _orchestrator;

        public MessageProcessor(ToolOrchestrator orchestrator)
        {
            _orchestrator = orchestrator;
        }

        public JsonRpcResponse ProcessRequest(JsonRpcRequest req)
        {
            Console.WriteLine($"[MCP Server] Processing JSON-RPC method: {req.Method}");

            switch (req.Method)
            {
                case "initialize":
                    _orchestrator.RegisterMcpTool("mcp_shell", "Execute commands via MCP");
                    return new JsonRpcResponse { Id = req.Id, Result = "initialized" };
                case "tools/call":
                    Console.WriteLine($"[MCP Server] Dispatched to tool execution via orchestrator: {req.Params}");
                    return new JsonRpcResponse { Id = req.Id, Result = $"Executed MCP tool call with args: {req.Params}" };
                default:
                    return new JsonRpcResponse { Id = req.Id, Error = "Method not found" };
            }
        }
    }

    // --- Agent Implementation ---

    public class AgentSession
    {
        public string SessionId { get; private set; }
        public AgentStatus Status { get; set; }
        private ToolOrchestrator _orchestrator;
        private MessageProcessor _mcpProcessor;
        private int _turnCounter;

        public AgentSession(string sessionId)
        {
            SessionId = sessionId;
            Status = AgentStatus.PendingInit;
            _orchestrator = new ToolOrchestrator();
            _mcpProcessor = new MessageProcessor(_orchestrator);
            _turnCounter = 0;
        }

        public void InitializeMcp()
        {
            var req = new JsonRpcRequest { Id = "0", Method = "initialize" };
            _mcpProcessor.ProcessRequest(req);
        }

        public void SteerInput(string input)
        {
            Status = AgentStatus.Running;
            _turnCounter++;
            string turnId = $"turn_{_turnCounter}";

            Console.WriteLine($"[Agent] Received input: '{input}'. Generating TurnContext ({turnId})...");

            var turnContext = new TurnContext
            {
                TurnId = turnId,
                SessionId = SessionId,
                Model = "gpt-5.5",
                WorkingDir = "/workspace",
                PermissionsProfile = "default"
            };

            if (input.StartsWith("mcp"))
            {
                var req = new JsonRpcRequest { Id = "1", Method = "tools/call", Params = input };
                var resp = _mcpProcessor.ProcessRequest(req);
                Console.WriteLine($"[Agent] MCP Turn executed. Result: {resp.Result}");
                Status = AgentStatus.Completed;
                return;
            }

            string toolReq = $"echo '{input}'";
            var ctx = new ToolCtx
            {
                CallId = $"call_{_turnCounter}",
                ToolName = "shell",
                TurnContext = turnContext
            };

            var runtime = new ShellCommandRuntime();

            try
            {
                string result = _orchestrator.ExecuteTool(runtime, toolReq, ctx);
                Console.WriteLine($"[Agent] Turn executed. Result: {result}");
                Status = AgentStatus.Completed;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"[Agent] Turn failed. Error: {ex.Message}");
                Status = AgentStatus.Errored;
            }
        }
    }

    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Welcome to Warp CLI (C# Edition) - Inspired by just-every-code");
            Console.WriteLine("Type '/help' for commands, or 'quit' to close.");

            var agent = new AgentSession("sess_123");
            agent.InitializeMcp();

            while (true)
            {
                Console.Write("warp> ");
                string? input = Console.ReadLine();

                if (string.IsNullOrWhiteSpace(input))
                {
                    continue;
                }

                input = input.Trim();

                if (input.Equals("quit", StringComparison.OrdinalIgnoreCase) || input.Equals("/quit", StringComparison.OrdinalIgnoreCase))
                {
                    agent.Status = AgentStatus.Shutdown;
                    break;
                }

                HandleCommand(input, agent);
            }
        }

        static void HandleCommand(string input, AgentSession agent)
        {
            if (input.StartsWith("/"))
            {
                int firstSpaceIndex = input.IndexOf(' ');
                string cmd;
                string args = "";

                if (firstSpaceIndex != -1)
                {
                    cmd = input.Substring(1, firstSpaceIndex - 1).ToLowerInvariant();
                    args = input.Substring(firstSpaceIndex + 1).Trim();
                }
                else
                {
                    cmd = input.Substring(1).ToLowerInvariant();
                }

                switch (cmd)
                {
                    case "pimono":
                        TriggerPiMonoDemo();
                        break;
                    case "store":
                        TriggerThreadStoreDemo();
                        break;
                    case "browser":
                        TriggerBrowserDemo();
                        break;
                    case "help":
                        Console.WriteLine("Available commands:");
                        Console.WriteLine("  /help        - Show this help message");
                        Console.WriteLine("  /prompt <msg>- Send a natural language prompt to the agent");
                        Console.WriteLine("  quit         - Quit the application");
                        break;
                    case "prompt":
                        if (string.IsNullOrEmpty(args))
                        {
                            Console.WriteLine("[Error] /prompt requires a message.");
                        }
                        else
                        {
                            agent.SteerInput(args);
                        }
                        break;
                    default:
                        Console.WriteLine($"Unknown command: /{cmd}");
                        break;
                }
            }
            else
            {
                agent.SteerInput(input);
            }
        }
        public static void TriggerPiMonoDemo()
        {
            var manager = new SessionManager("sess_pi_456", "/workspace", "/tmp/sessions");
            manager.AppendMessageEntry("1", "", "Initial prompt");
            manager.AppendMessageEntry("2", "1", "Assistant response");
            Console.WriteLine($"[SessionManager] Tracked entries: {manager.FileEntries.Count}");
        }

        public static void TriggerThreadStoreDemo()
        {
            var store = new InMemoryThreadStore();
            string threadId = "test-thread-123";

            Console.WriteLine($"[ThreadStore] Creating thread: {threadId}");
            store.CreateThread(threadId, "/workspace");

            Console.WriteLine("[ThreadStore] Appending history items...");
            store.AppendItem(threadId, "{\"event\":\"user_input\", \"content\":\"hello\"}");
            store.AppendItem(threadId, "{\"event\":\"agent_reply\", \"content\":\"hi there!\"}");

            var threadInfo = store.LoadThread(threadId);
            Console.WriteLine($"[ThreadStore] Loaded thread metadata: StoredThread {{ ThreadId: \"{threadInfo.ThreadId}\", Preview: \"{threadInfo.Preview}\", Cwd: \"{threadInfo.Cwd}\" }}");

            var history = store.LoadHistory(threadId);
            Console.WriteLine($"[ThreadStore] Loaded thread history items: {history.Items.Count}");
        }

        public static void TriggerBrowserDemo()
        {
            var config = new BrowserConfig
            {
                Headless = true,
                Viewport = new ViewportConfig
                {
                    Width = 1280,
                    Height = 720,
                    DeviceScaleFactor = 1.0
                }
            };

            var manager = new BrowserManager(config);
            manager.Launch();
            if (manager.ActivePage != null)
            {
                manager.ActivePage.Navigate("https://github.com/just-every/code");
                manager.ActivePage.DispatchMouseEvent(250.0, 300.0, "click");
                manager.ActivePage.DispatchKeyEvent("Hello Warp!");
                manager.ActivePage.CaptureScreenshot();
            }
        }
    }
}
