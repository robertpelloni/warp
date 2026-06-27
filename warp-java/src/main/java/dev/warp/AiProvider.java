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
        public AssistantMessageEvent(EventType type, Exception error) { this.type = type; this.error = error; }
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
    public static class OpenAiProvider implements ApiProvider {
        private final String apiKey;
        private final java.net.http.HttpClient client;

        public OpenAiProvider(String apiKey) {
            this.apiKey = apiKey;
            this.client = java.net.http.HttpClient.newHttpClient();
        }

        @Override
        public Flow.Publisher<AssistantMessageEvent> stream(Model model, Context ctx) {
            SubmissionPublisher<AssistantMessageEvent> publisher = new SubmissionPublisher<>();

            CompletableFuture.runAsync(() -> {
                try {
                    publisher.submit(new AssistantMessageEvent(EventType.START));

                    StringBuilder messagesJson = new StringBuilder("[");
                    if (ctx.systemPrompt != null && !ctx.systemPrompt.isEmpty()) {
                        messagesJson.append(String.format("{\"role\": \"system\", \"content\": \"%s\"}", escapeJson(ctx.systemPrompt)));
                    }
                    for (int i = 0; i < ctx.messages.size(); i++) {
                        if (messagesJson.length() > 1) messagesJson.append(",");
                        Message msg = ctx.messages.get(i);
                        StringBuilder contentStr = new StringBuilder();
                        for (ContentPart part : msg.content) {
                            if (part.type == ContentPartType.TEXT) contentStr.append(part.text);
                        }
                        messagesJson.append(String.format("{\"role\": \"%s\", \"content\": \"%s\"}", msg.role.name().toLowerCase(), escapeJson(contentStr.toString())));
                    }
                    messagesJson.append("]");

                    String body = String.format("{\"model\": \"%s\", \"stream\": true, \"messages\": %s}", model.modelID, messagesJson.toString());

                    java.net.http.HttpRequest request = java.net.http.HttpRequest.newBuilder()
                            .uri(java.net.URI.create("https://api.openai.com/v1/chat/completions"))
                            .header("Authorization", "Bearer " + apiKey)
                            .header("Content-Type", "application/json")
                            .POST(java.net.http.HttpRequest.BodyPublishers.ofString(body))
                            .build();

                    java.net.http.HttpResponse<java.io.InputStream> response = client.send(request, java.net.http.HttpResponse.BodyHandlers.ofInputStream());

                    if (response.statusCode() != 200) {
                        publisher.submit(new AssistantMessageEvent(EventType.ERROR, new Exception("OpenAI API error: " + response.statusCode())));
                        publisher.close();
                        return;
                    }

                    try (java.io.BufferedReader reader = new java.io.BufferedReader(new java.io.InputStreamReader(response.body()))) {
                        String line;
                        while ((line = reader.readLine()) != null) {
                            if (line.startsWith("data: ")) {
                                String data = line.substring(6).trim();
                                if (data.equals("[DONE]")) break;

                                String searchString = "\"content\":\"";
                                int contentIndex = data.indexOf(searchString);
                                if (contentIndex != -1) {
                                    int endQuoteIndex = data.indexOf("\"", contentIndex + searchString.length());
                                    if (endQuoteIndex != -1) {
                                        String deltaContent = data.substring(contentIndex + searchString.length(), endQuoteIndex);
                                        deltaContent = deltaContent.replace("\\n", "\n").replace("\\\"", "\"").replace("\\\\", "\\");
                                        publisher.submit(new AssistantMessageEvent(EventType.TEXT_DELTA, deltaContent));
                                    }
                                }
                            }
                        }
                    }
                    publisher.submit(new AssistantMessageEvent(EventType.DONE));
                } catch (Exception e) {
                    publisher.submit(new AssistantMessageEvent(EventType.ERROR, e));
                } finally {
                    publisher.close();
                }
            });

            return publisher;
        }

        private String escapeJson(String input) {
            return input.replace("\\", "\\\\").replace("\"", "\\\"").replace("\n", "\\n");
        }
    }
}
