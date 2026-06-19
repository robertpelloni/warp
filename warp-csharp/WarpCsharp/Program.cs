using System;

namespace WarpCsharp
{
    class Program
    {
        static void Main(string[] args)
        {
            // Run manual tests if requested via hidden flag
            if (args.Length > 0 && args[0] == "--test-auth")
            {
                Tests.AuthTests.RunTests();
                Console.WriteLine("C# Auth Tests Passed");
                return;
            }

            if (args.Length == 0)
            {
                Console.WriteLine("Warp C# Initialized");
                Console.WriteLine("Usage: warp <command> [arguments]");
                return;
            }

            string command = args[0];

            switch (command)
            {
                case "run":
                    string agent = "default";
                    if (args.Length > 2 && args[1] == "--agent")
                    {
                        agent = args[2];
                    }
                    Console.WriteLine($"Warp C# running agent: {agent}");
                    break;
                default:
                    Console.WriteLine($"Unknown command: {command}");
                    break;
            }
        }
    }
}
