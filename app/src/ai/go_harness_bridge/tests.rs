#[cfg(test)]
mod tests {
    use crate::ai::go_harness_bridge::GoHarnessBridge;
    use std::process::{Command, Child};
    use std::time::Duration;
    use tokio::time::sleep;

    struct GoSidecar {
        child: Child,
    }

    impl GoSidecar {
        fn start(port: u16) -> Self {
            let child = Command::new("./rebuild/warp")
                .arg("-port")
                .arg(port.to_string())
                .spawn()
                .expect("Failed to start Go sidecar. Make sure it is built.");
            Self { child }
        }
    }

    impl Drop for GoSidecar {
        fn drop(&mut self) {
            let _ = self.child.kill();
        }
    }

    #[tokio::test]
    async fn test_bridge_integration() {
        let port = 10006;
        // Build the sidecar first to ensure it exists
        let build_status = Command::new("go")
            .args(&["build", "-o", "rebuild/warp", "./rebuild/cmd/warp"])
            .status()
            .expect("Failed to build Go sidecar");
        assert!(build_status.success());

        let _sidecar = GoSidecar::start(port);
        let bridge = GoHarnessBridge::new(port);

        // Wait for sidecar to start
        let mut success = false;
        for _ in 0..10 {
            if bridge.check_health().await.is_ok() {
                success = true;
                break;
            }
            sleep(Duration::from_millis(200)).await;
        }
        assert!(success, "Sidecar failed to start or health check failed");

        let resp = bridge.run_prompt("Say hello from Go").await.expect("Bridge call failed");
        assert!(!resp.history.is_empty());
        assert!(resp.history.iter().any(|m| m.role == "assistant"));
    }
}
