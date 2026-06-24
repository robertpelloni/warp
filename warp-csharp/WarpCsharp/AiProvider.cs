using System;
using System.Collections.Generic;
using System.Threading.Channels;

namespace WarpCsharp
{
    public enum MessageRole { System, User, Assistant, Tool }

    public enum ContentPartType { Text, Image, ToolCall }

    public class ContentPart
    {
        public ContentPartType Type { get; set; }
        public string Text { get; set; } = "";
        public string CallID { get; set; } = "";
        public string Name { get; set; } = "";
        public string Args { get; set; } = "";
    }

    public class Message
    {
        public MessageRole Role { get; set; }
        public List<ContentPart> Content { get; set; } = new();
    }

    public class Context
    {
        public string SystemPrompt { get; set; } = "";
        public List<Message> Messages { get; set; } = new();
        public List<string> Tools { get; set; } = new();
    }

    public class Model
    {
        public string ProviderID { get; set; } = "";
        public string ModelID { get; set; } = "";
        public int ContextWindow { get; set; }
        public int MaxTokens { get; set; }
        public bool SupportsImages { get; set; }
        public bool SupportsTools { get; set; }
    }

    public enum EventType { Start, TextDelta, ToolCallStart, ToolCallDelta, Done, Error }

    public class AssistantMessageEvent
    {
        public EventType Type { get; set; }
        public string Data { get; set; } = "";
        public Exception? Error { get; set; }
    }

    public interface IApiProvider
    {
        ChannelReader<AssistantMessageEvent> Stream(Model model, Context ctx);
    }

    public class DummyProvider : IApiProvider
    {
        public string Name { get; set; } = "Dummy";

        public ChannelReader<AssistantMessageEvent> Stream(Model model, Context ctx)
        {
            var channel = Channel.CreateUnbounded<AssistantMessageEvent>();
            var writer = channel.Writer;

            _ = System.Threading.Tasks.Task.Run(async () =>
            {
                await writer.WriteAsync(new AssistantMessageEvent { Type = EventType.Start });

                string[] words = $"This is a simulated streaming response from {Name} using model {model.ModelID}. ".Split(' ');
                foreach (var w in words)
                {
                    await writer.WriteAsync(new AssistantMessageEvent { Type = EventType.TextDelta, Data = w + " " });
                    await System.Threading.Tasks.Task.Delay(10);
                }

                if (ctx.Tools.Count > 0)
                {
                    await writer.WriteAsync(new AssistantMessageEvent { Type = EventType.ToolCallStart, Data = "mcp_shell" });
                    await writer.WriteAsync(new AssistantMessageEvent { Type = EventType.ToolCallDelta, Data = "{\"cmd\": \"echo hello\"}" });
                }

                await writer.WriteAsync(new AssistantMessageEvent { Type = EventType.Done });
                writer.Complete();
            });

            return channel.Reader;
        }
    }

    public class ApiProviderRegistry
    {
        private readonly Dictionary<string, IApiProvider> _providers = new();
        private readonly Dictionary<string, Model> _models = new();

        public void RegisterProvider(string id, IApiProvider provider)
        {
            _providers[id] = provider;
            Console.WriteLine($"[ApiRegistry] Registered provider: {id}");
        }

        public void RegisterModel(Model model)
        {
            string key = $"{model.ProviderID}:{model.ModelID}";
            _models[key] = model;
            Console.WriteLine($"[ApiRegistry] Registered model: {key}");
        }

        public ChannelReader<AssistantMessageEvent> ExecuteStream(string providerID, string modelID, Context ctx)
        {
            string key = $"{providerID}:{modelID}";
            if (!_models.TryGetValue(key, out var model))
            {
                throw new Exception($"Model {key} not found");
            }

            if (!_providers.TryGetValue(providerID, out var provider))
            {
                throw new Exception($"Provider {providerID} not found");
            }

            return provider.Stream(model, ctx);
        }
    }
}
