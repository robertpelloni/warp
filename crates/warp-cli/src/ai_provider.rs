use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::{mpsc, RwLock};
use tokio::time::{sleep, Duration};

#[derive(Debug, Clone, PartialEq)]
pub enum MessageRole {
    System,
    User,
    Assistant,
    Tool,
}

#[derive(Debug, Clone, PartialEq)]
pub enum ContentPartType {
    Text,
    Image,
    ToolCall,
}

#[derive(Debug, Clone)]
pub struct ContentPart {
    pub part_type: ContentPartType,
    pub text: String,
    pub call_id: String,
    pub name: String,
    pub args: String,
}

#[derive(Debug, Clone)]
pub struct Message {
    pub role: MessageRole,
    pub content: Vec<ContentPart>,
}

#[derive(Debug, Clone)]
pub struct Context {
    pub system_prompt: String,
    pub messages: Vec<Message>,
    pub tools: Vec<String>,
}

#[derive(Debug, Clone)]
pub struct Model {
    pub provider_id: String,
    pub model_id: String,
    pub context_window: usize,
    pub max_tokens: usize,
    pub supports_images: bool,
    pub supports_tools: bool,
}

#[derive(Debug, Clone)]
pub enum EventType {
    Start,
    TextDelta,
    ToolCallStart,
    ToolCallDelta,
    Done,
    Error,
}

#[derive(Debug, Clone)]
pub struct AssistantMessageEvent {
    pub event_type: EventType,
    pub data: String,
    pub error: Option<String>,
}

#[async_trait::async_trait]
pub trait ApiProvider: Send + Sync {
    async fn stream(&self, model: Model, ctx: Context) -> mpsc::Receiver<AssistantMessageEvent>;
}

pub struct DummyProvider {
    pub name: String,
}

#[async_trait::async_trait]
impl ApiProvider for DummyProvider {
    async fn stream(&self, model: Model, ctx: Context) -> mpsc::Receiver<AssistantMessageEvent> {
        let (tx, rx) = mpsc::channel(100);
        let name = self.name.clone();

        tokio::spawn(async move {
            let _ = tx.send(AssistantMessageEvent { event_type: EventType::Start, data: "".to_string(), error: None }).await;

            let phrase = format!("This is a simulated streaming response from {} using model {}. ", name, model.model_id);
            let words: Vec<&str> = phrase.split_whitespace().collect();

            for word in words {
                let _ = tx.send(AssistantMessageEvent {
                    event_type: EventType::TextDelta,
                    data: format!("{} ", word),
                    error: None,
                }).await;
                sleep(Duration::from_millis(10)).await;
            }

            if !ctx.tools.is_empty() {
                let _ = tx.send(AssistantMessageEvent { event_type: EventType::ToolCallStart, data: "mcp_shell".to_string(), error: None }).await;
                let _ = tx.send(AssistantMessageEvent { event_type: EventType::ToolCallDelta, data: "{\"cmd\": \"echo hello\"}".to_string(), error: None }).await;
            }

            let _ = tx.send(AssistantMessageEvent { event_type: EventType::Done, data: "".to_string(), error: None }).await;
        });

        rx
    }
}

pub struct ApiProviderRegistry {
    providers: Arc<RwLock<HashMap<String, Arc<dyn ApiProvider>>>>,
    models: Arc<RwLock<HashMap<String, Model>>>,
}

impl ApiProviderRegistry {
    pub fn new() -> Self {
        Self {
            providers: Arc::new(RwLock::new(HashMap::new())),
            models: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn register_provider(&self, id: String, provider: Arc<dyn ApiProvider>) {
        let mut p = self.providers.write().await;
        p.insert(id.clone(), provider);
        println!("[ApiRegistry] Registered provider: {}", id);
    }

    pub async fn register_model(&self, model: Model) {
        let key = format!("{}:{}", model.provider_id, model.model_id);
        let mut m = self.models.write().await;
        m.insert(key.clone(), model);
        println!("[ApiRegistry] Registered model: {}", key);
    }

    pub async fn execute_stream(&self, provider_id: &str, model_id: &str, ctx: Context) -> Result<mpsc::Receiver<AssistantMessageEvent>, String> {
        let key = format!("{}:{}", provider_id, model_id);

        let models = self.models.read().await;
        let model = models.get(&key).cloned().ok_or_else(|| format!("Model {} not found", key))?;
        drop(models);

        let providers = self.providers.read().await;
        let provider = providers.get(provider_id).cloned().ok_or_else(|| format!("Provider {} not found", provider_id))?;
        drop(providers);

        Ok(provider.stream(model, ctx).await)
    }
}
