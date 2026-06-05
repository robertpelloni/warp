#[cfg(test)]
mod tests {
    use super::super::go_harness_bridge::GoHarnessBridge;
    use std::process::{Command, Child};
    use std::time::Duration;
    use tokio::time::sleep;
    use std::path::PathBuf;

    struct GoServiceGuard {
        child: Child,
    }

    impl Drop for GoServiceGuard {
        fn drop(&mut self) {
            let _ = self.child.kill();
        }
    }

    #[tokio::test]
    async fn test_rust_to_go_bridge() {
        let mut rebuild_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
        rebuild_path.push("rebuild");
        // Adjust path if needed. If CARGO_MANIFEST_DIR is app/, we need to go up.
        if !rebuild_path.exists() {
             rebuild_path = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
             rebuild_path.pop();
             rebuild_path.push("rebuild");
        }

        let port = 10003;
        let child = Command::new("./warp")
            .arg("-port")
            .arg(port.to_string())
            .current_dir(&rebuild_path)
            .spawn()
            .expect("Failed to start Go harness service");

        let _guard = GoServiceGuard { child };

        let bridge = GoHarnessBridge::new(port);
        let mut ready = false;
        for _ in 0..20 {
            if bridge.check_health().await.is_ok() {
                ready = true;
                break;
            }
            sleep(Duration::from_millis(500)).await;
        }

        assert!(ready, "Go harness service did not become ready in time");

        let response = bridge.run_prompt("Test integration from Rust")
            .await
            .expect("Failed to run prompt via bridge");

        assert!(!response.history.is_empty());
        assert!(response.history.iter().any(|m| m.role == "assistant"));
    }
}
