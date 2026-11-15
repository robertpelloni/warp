use crate::ai_provider::{ApiProviderRegistry, Context, DummyProvider, Message, MessageRole, ContentPart, ContentPartType, Model};
use std::sync::Arc;

pub async fn trigger_ai_demo() {
    let registry = ApiProviderRegistry::new();

    registry.register_provider("openai".to_string(), Arc::new(DummyProvider { name: "OpenAI-Mock".to_string() })).await;
    registry.register_model(Model {
        provider_id: "openai".to_string(),
        model_id: "gpt-4".to_string(),
        context_window: 8192,
        max_tokens: 4096,
        supports_images: true,
        supports_tools: true,
    }).await;

    registry.register_provider("anthropic".to_string(), Arc::new(DummyProvider { name: "Anthropic-Mock".to_string() })).await;
    registry.register_model(Model {
        provider_id: "anthropic".to_string(),
        model_id: "claude-3-opus".to_string(),
        context_window: 200000,
        max_tokens: 4096,
        supports_images: true,
        supports_tools: true,
    }).await;

    let ctx = Context {
        system_prompt: "You are a helpful coding assistant.".to_string(),
        messages: vec![Message {
            role: MessageRole::User,
            content: vec![ContentPart {
                part_type: ContentPartType::Text,
                text: "Write a function.".to_string(),
                call_id: "".to_string(),
                name: "".to_string(),
                args: "".to_string(),
            }],
        }],
        tools: vec!["mcp_shell".to_string()],
    };

    println!("=== Executing AI Stream (OpenAI) ===");
    if let Ok(mut stream1) = registry.execute_stream("openai", "gpt-4", ctx.clone()).await {
        while let Some(event) = stream1.recv().await {
            println!("Event: {:?} - Data: {}", event.event_type, event.data);
        }
    }

    println!("=== Executing AI Stream (Anthropic) ===");
    if let Ok(mut stream2) = registry.execute_stream("anthropic", "claude-3-opus", ctx.clone()).await {
        while let Some(event) = stream2.recv().await {
            println!("Event: {:?} - Data: {}", event.event_type, event.data);
        }
    }
}
