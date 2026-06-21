using System;

namespace WarpCsharp
{
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
        public TOut ExecuteTool<TReq, TOut>(IToolRuntime<TReq, TOut> runtime, TReq req, ToolCtx ctx)
        {
            // 1. Approval Phase
            if (runtime.RequiresApproval(req))
            {
                Console.WriteLine($"[Orchestrator] Requesting tool approval for '{ctx.ToolName}'...");
                // Assuming approved
            }

            // 2. Sandbox Selection
            var attempt = new SandboxAttempt
            {
                IsSandboxed = runtime.SandboxPreference == SandboxPreference.Sandboxed,
                Cwd = ctx.TurnContext.WorkingDir
            };

            // 3. Execution
            return runtime.Run(req, attempt, ctx);
        }
    }

    public class AgentSession
    {
        public string SessionId { get; private set; }
        public AgentStatus Status { get; set; }
        private ToolOrchestrator _orchestrator;
        private int _turnCounter;

        public AgentSession(string sessionId)
        {
            SessionId = sessionId;
            Status = AgentStatus.PendingInit;
            _orchestrator = new ToolOrchestrator();
            _turnCounter = 0;
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
    }
}
