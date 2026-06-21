using System;
using System.Threading;

namespace WarpCsharp
{
    class ToolOrchestrator
    {
        public bool SandboxEnabled { get; set; } = true;

        public string ExecuteTask(string task)
        {
            Console.WriteLine($"[Orchestrator] Requesting approval for task: '{task}'");
            Thread.Sleep(500);

            if (SandboxEnabled)
            {
                Console.WriteLine("[Orchestrator] Running in sandbox mode...");
            }

            string[] steps = { "Plan", "Code", "Review" };
            foreach (var step in steps)
            {
                Console.WriteLine($"[Agent: {step}] Processing...");
                Thread.Sleep(400);
            }

            return $"Task '{task}' completed successfully by Auto Drive.";
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
                        Console.WriteLine("  /auto <task> - Start Auto Drive orchestration for a task");
                        Console.WriteLine("  quit         - Quit the application");
                        break;
                    case "auto":
                        if (string.IsNullOrEmpty(args))
                        {
                            Console.WriteLine("[Error] /auto requires a task description.");
                        }
                        else
                        {
                            Console.WriteLine("[Orchestrator: Auto Drive] Starting autonomous loop.");
                            string result = orchestrator.ExecuteTask(args);
                            Console.WriteLine($"[Result] {result}");
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
