package dev.warp;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Flow;
import java.util.concurrent.SubmissionPublisher;
import java.util.concurrent.CompletableFuture;

public class AiProvider {

    public enum MessageRole { SYSTEM, USER, ASSISTANT, TOOL }

    public enum ContentPartType { TEXT, IMAGE, TOOL_CALL }

    public static class ContentPart {
        public ContentPartType type;
        public String text = "";
        public String callID = "";
        public String name = "";
        public String args = "";
    }

    public static class Message {
        public MessageRole role;
        public List<ContentPart> content = new ArrayList<>();
    }

    public static class Context {
        public String systemPrompt = "";
        public List<Message> messages = new ArrayList<>();
        public List<String> tools = new ArrayList<>();
    }

    public static class Model {
        public String providerID = "";
        public String modelID = "";
        public int contextWindow;
        public int maxTokens;
        public boolean supportsImages;
        public boolean supportsTools;
    }

    public enum EventType { START, TEXT_DELTA, TOOL_CALL_START, TOOL_CALL_DELTA, DONE, ERROR }

    public static class AssistantMessageEvent {
        public EventType type;
        public String data = "";
        public Exception error;

        public AssistantMessageEvent(EventType type) { this.type = type; }
        public AssistantMessageEvent(EventType type, String data) { this.type = type; this.data = data; }
    }

    public interface ApiProvider {
        Flow.Publisher<AssistantMessageEvent> stream(Model model, Context ctx);
    }

    public static class DummyProvider implements ApiProvider {
        public String name = "Dummy";

        @Override
        public Flow.Publisher<AssistantMessageEvent> stream(Model model, Context ctx) {
            SubmissionPublisher<AssistantMessageEvent> publisher = new SubmissionPublisher<>();

            CompletableFuture.runAsync(() -> {
                publisher.submit(new AssistantMessageEvent(EventType.START));

                String[] words = ("This is a simulated streaming response from " + name + " using model " + model.modelID + ". ").split(" ");
                for (String w : words) {
                    publisher.submit(new AssistantMessageEvent(EventType.TEXT_DELTA, w + " "));
                    try { Thread.sleep(10); } catch (InterruptedException ignored) {}
                }

                if (!ctx.tools.isEmpty()) {
                    publisher.submit(new AssistantMessageEvent(EventType.TOOL_CALL_START, "mcp_shell"));
                    publisher.submit(new AssistantMessageEvent(EventType.TOOL_CALL_DELTA, "{\"cmd\": \"echo hello\"}"));
                }

                publisher.submit(new AssistantMessageEvent(EventType.DONE));
                publisher.close();
            });

            return publisher;
        }
    }

    public static class ApiProviderRegistry {
        private final Map<String, ApiProvider> providers = new ConcurrentHashMap<>();
        private final Map<String, Model> models = new ConcurrentHashMap<>();

        public void registerProvider(String id, ApiProvider provider) {
            providers.put(id, provider);
            System.out.println("[ApiRegistry] Registered provider: " + id);
        }

        public void registerModel(Model model) {
            String key = model.providerID + ":" + model.modelID;
            models.put(key, model);
            System.out.println("[ApiRegistry] Registered model: " + key);
        }

        public Flow.Publisher<AssistantMessageEvent> executeStream(String providerID, String modelID, Context ctx) throws Exception {
            String key = providerID + ":" + modelID;
            Model model = models.get(key);
            if (model == null) {
                throw new Exception("Model " + key + " not found");
            }

            ApiProvider provider = providers.get(providerID);
            if (provider == null) {
                throw new Exception("Provider " + providerID + " not found");
            }

            return provider.stream(model, ctx);
        }
    }
}
