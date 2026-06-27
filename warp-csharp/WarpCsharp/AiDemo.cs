using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace WarpCsharp
{
    public static class AiDemo
    {
        public static async Task TriggerAiDemo()
        {
            var registry = new ApiProviderRegistry();

            string? openaiKey = Environment.GetEnvironmentVariable("OPENAI_API_KEY");
            if (string.IsNullOrEmpty(openaiKey))
            {
                Console.WriteLine("[Info] No OPENAI_API_KEY found. Falling back to DummyProvider for OpenAI.");
                registry.RegisterProvider("openai", new DummyProvider { Name = "OpenAI-Mock" });
            }
            else
            {
                Console.WriteLine("[Info] Found OPENAI_API_KEY. Registering actual OpenAiProvider.");
                registry.RegisterProvider("openai", new OpenAiProvider(openaiKey));
            }

            registry.RegisterModel(new Model { ProviderID = "openai", ModelID = "gpt-4", ContextWindow = 8192, SupportsTools = true });

            registry.RegisterProvider("anthropic", new DummyProvider { Name = "Anthropic-Mock" });
            registry.RegisterModel(new Model { ProviderID = "anthropic", ModelID = "claude-3-opus", ContextWindow = 200000, SupportsTools = true });

            var ctx = new Context
            {
                SystemPrompt = "You are a helpful coding assistant.",
                Messages = new List<Message> { new Message { Role = MessageRole.User, Content = new List<ContentPart> { new ContentPart { Type = ContentPartType.Text, Text = "Write a short haiku about coding." } } } },
                Tools = new List<string>()
            };

            Console.WriteLine("=== Executing AI Stream (OpenAI) ===");
            var stream1 = registry.ExecuteStream("openai", "gpt-4", ctx);
            await foreach (var ev in stream1.ReadAllAsync())
            {
                if (ev.Type == EventType.TextDelta) Console.Write(ev.Data);
                else if (ev.Type == EventType.Error) Console.WriteLine($"\n[Error] {ev.Error?.Message}");
                else if (ev.Type == EventType.Done) Console.WriteLine("\n[Done]");
            }

            Console.WriteLine("=== Executing AI Stream (Anthropic) ===");
            var stream2 = registry.ExecuteStream("anthropic", "claude-3-opus", ctx);
            await foreach (var ev in stream2.ReadAllAsync())
            {
                if (ev.Type == EventType.TextDelta) Console.Write(ev.Data);
                else if (ev.Type == EventType.Error) Console.WriteLine($"\n[Error] {ev.Error?.Message}");
                else if (ev.Type == EventType.Done) Console.WriteLine("\n[Done]");
            }
        }
    }
}
