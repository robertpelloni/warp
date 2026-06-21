using System;

namespace WarpCsharp
{
    public enum SandboxPreference
    {
        Unsandboxed,
        Sandboxed
    }

    public class ToolCtx
    {
        public string CallId { get; set; } = string.Empty;
        public string ToolName { get; set; } = string.Empty;
        public string SessionId { get; set; } = string.Empty;
        public string TurnId { get; set; } = string.Empty;
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
            return $"[{ctx.ToolName}] Executed shell command '{req}' in {mode} mode (CallID: {ctx.CallId})";
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
                Cwd = "/workspace"
            };

            // 3. Execution
            return runtime.Run(req, attempt, ctx);
        }
    }

    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Welcome to Warp CLI (C# Edition) - Inspired by just-every-code");
            Console.WriteLine("Type '/help' for commands, or 'quit' to close.");

            var orchestrator = new ToolOrchestrator();

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
                    break;
                }

                HandleCommand(input, orchestrator);
            }
        }

        static void HandleCommand(string input, ToolOrchestrator orchestrator)
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
                        Console.WriteLine("  /shell <cmd> - Run a command through the ToolOrchestrator");
                        Console.WriteLine("  quit         - Quit the application");
                        break;
                    case "shell":
                        if (string.IsNullOrEmpty(args))
                        {
                            Console.WriteLine("[Error] /shell requires a command.");
                        }
                        else
                        {
                            var ctx = new ToolCtx
                            {
                                CallId = "call_abc123",
                                ToolName = "shell",
                                SessionId = "sess_1",
                                TurnId = "turn_1"
                            };
                            var runtime = new ShellCommandRuntime();

                            try
                            {
                                string result = orchestrator.ExecuteTool(runtime, args, ctx);
                                Console.WriteLine($"[Result] {result}");
                            }
                            catch (Exception ex)
                            {
                                Console.WriteLine($"[Error] {ex.Message}");
                            }
                        }
                        break;
                    default:
                        Console.WriteLine($"Unknown command: /{cmd}");
                        break;
                }
            }
            else
            {
                Console.WriteLine($"[Agent] Echoing input: {input}");
            }
        }
    }
}
