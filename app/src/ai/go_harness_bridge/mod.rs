use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};

#[derive(Serialize)]
struct RunRequest {
    prompt: String,
}

#[derive(Deserialize, Debug)]
pub struct Message {
    pub role: String,
    pub content: String,
}

#[derive(Deserialize, Debug)]
pub struct RunResponse {
    pub history: Vec<Message>,
    pub error: Option<String>,
}

pub struct GoHarnessBridge {
    base_url: String,
    client: reqwest::Client,
}

impl GoHarnessBridge {
    pub fn new(port: u16) -> Self {
        Self {
            base_url: format!("http://localhost:{}", port),
            client: reqwest::Client::new(),
        }
    }

    pub async fn run_prompt(&self, prompt: &str) -> Result<RunResponse> {
        let url = format!("{}/run", self.base_url);
        let req = RunRequest { prompt: prompt.to_string() };
        let resp = self.client.post(&url).json(&req).send().await.context("Go harness call failed")?;
        let run_resp: RunResponse = resp.json().await.context("Go harness decode failed")?;
        Ok(run_resp)
    }

    pub async fn check_health(&self) -> Result<()> {
        let url = format!("{}/health", self.base_url);
        let resp = self.client.get(&url).send().await?;
        if resp.status().is_success() { Ok(()) } else { Err(anyhow::anyhow!("health fail")) }
    }
}

#[cfg(test)]
mod tests;
