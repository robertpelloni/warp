package dev.warp;

import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

public class AuthTest {
    @Test
    public void testAuthNotConfigured() {
        assertFalse(Auth.authenticateUser("user", "pass"), "Should fail when not configured");
    }
}
