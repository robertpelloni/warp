package dev.warp;

public class AiDemo {
    public static void triggerAiDemo() {
        AiProvider.ApiProviderRegistry registry = new AiProvider.ApiProviderRegistry();

        String openaiKey = System.getenv("OPENAI_API_KEY");
        if (openaiKey == null || openaiKey.isEmpty()) {
            System.out.println("[Info] No OPENAI_API_KEY found. Falling back to DummyProvider for OpenAI.");
            AiProvider.DummyProvider openai = new AiProvider.DummyProvider();
            openai.name = "OpenAI-Mock";
            registry.registerProvider("openai", openai);
        } else {
            System.out.println("[Info] Found OPENAI_API_KEY. Registering actual OpenAiProvider.");
            registry.registerProvider("openai", new AiProvider.OpenAiProvider(openaiKey));
        }

        AiProvider.Model gpt4 = new AiProvider.Model();
        gpt4.providerID = "openai";
        gpt4.modelID = "gpt-4";
        gpt4.contextWindow = 8192;
        gpt4.supportsTools = true;
        registry.registerModel(gpt4);

        AiProvider.DummyProvider anthropic = new AiProvider.DummyProvider();
        anthropic.name = "Anthropic-Mock";
        registry.registerProvider("anthropic", anthropic);

        AiProvider.Model claude3 = new AiProvider.Model();
        claude3.providerID = "anthropic";
        claude3.modelID = "claude-3-opus";
        claude3.contextWindow = 200000;
        claude3.supportsTools = true;
        registry.registerModel(claude3);

        AiProvider.Context ctx = new AiProvider.Context();
        ctx.systemPrompt = "You are a helpful coding assistant.";
        AiProvider.Message msg = new AiProvider.Message();
        msg.role = AiProvider.MessageRole.USER;
        AiProvider.ContentPart pt = new AiProvider.ContentPart();
        pt.type = AiProvider.ContentPartType.TEXT;
        pt.text = "Write a short haiku about coding.";
        msg.content.add(pt);
        ctx.messages.add(msg);

        try {
            System.out.println("=== Executing AI Stream (OpenAI) ===");
            java.util.concurrent.Flow.Publisher<AiProvider.AssistantMessageEvent> stream1 = registry.executeStream("openai", "gpt-4", ctx);
            java.util.concurrent.CountDownLatch latch1 = new java.util.concurrent.CountDownLatch(1);
            stream1.subscribe(new java.util.concurrent.Flow.Subscriber<>() {
                public void onSubscribe(java.util.concurrent.Flow.Subscription s) { s.request(Long.MAX_VALUE); }
                public void onNext(AiProvider.AssistantMessageEvent ev) {
                    if (ev.type == AiProvider.EventType.TEXT_DELTA) System.out.print(ev.data);
                    else if (ev.type == AiProvider.EventType.ERROR) System.out.println("\n[Error] " + ev.error);
                    else if (ev.type == AiProvider.EventType.DONE) System.out.println("\n[Done]");
                }
                public void onError(Throwable t) { t.printStackTrace(); latch1.countDown(); }
                public void onComplete() { latch1.countDown(); }
            });
            latch1.await();

            System.out.println("=== Executing AI Stream (Anthropic) ===");
            java.util.concurrent.Flow.Publisher<AiProvider.AssistantMessageEvent> stream2 = registry.executeStream("anthropic", "claude-3-opus", ctx);
            java.util.concurrent.CountDownLatch latch2 = new java.util.concurrent.CountDownLatch(1);
            stream2.subscribe(new java.util.concurrent.Flow.Subscriber<>() {
                public void onSubscribe(java.util.concurrent.Flow.Subscription s) { s.request(Long.MAX_VALUE); }
                public void onNext(AiProvider.AssistantMessageEvent ev) {
                    if (ev.type == AiProvider.EventType.TEXT_DELTA) System.out.print(ev.data);
                    else if (ev.type == AiProvider.EventType.ERROR) System.out.println("\n[Error] " + ev.error);
                    else if (ev.type == AiProvider.EventType.DONE) System.out.println("\n[Done]");
                }
                public void onError(Throwable t) { t.printStackTrace(); latch2.countDown(); }
                public void onComplete() { latch2.countDown(); }
            });
            latch2.await();

        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
