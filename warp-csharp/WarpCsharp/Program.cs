using System;

namespace WarpCsharp
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Welcome to Warp CLI (C# Edition) - Inspired by just-every-code");
            Console.WriteLine("Type '/help' for commands, or 'quit' to close.");

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

                HandleCommand(input);
            }
        }

        static void HandleCommand(string input)
        {
            if (input.StartsWith("/"))
            {
                string cmd = input.Substring(1).ToLowerInvariant();
                switch (cmd)
                {
                    case "help":
                        Console.WriteLine("Available commands:");
                        Console.WriteLine("  /help     - Show this help message");
                        Console.WriteLine("  /plan     - (Stub) Coordinate planning agent");
                        Console.WriteLine("  /code     - (Stub) Coordinate coding agent");
                        Console.WriteLine("  /auto     - (Stub) Start Auto Drive orchestration");
                        Console.WriteLine("  quit      - Quit the application");
                        break;
                    case "plan":
                        Console.WriteLine("[Agent: Planner] Acknowledged. Ready to plan task.");
                        break;
                    case "code":
                        Console.WriteLine("[Agent: Coder] Acknowledged. Ready to write code.");
                        break;
                    case "auto":
                        Console.WriteLine("[Orchestrator: Auto Drive] Starting autonomous loop (Mock).");
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
