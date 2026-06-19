package dev.warp;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

public class Auth {
    public static boolean authenticateUser(String username, String password) {
        System.out.println("Warp Java: Authenticating user " + username);

        String expectedUser = System.getenv("WARP_ADMIN_USER");
        String expectedHash = System.getenv("WARP_ADMIN_HASH");

        if (expectedUser == null || expectedHash == null) {
            System.out.println("Warning: Authentication not configured");
            return false;
        }

        try {
            MessageDigest md = MessageDigest.getInstance("SHA-256");
            byte[] hash = md.digest(password.getBytes());
            StringBuilder hexString = new StringBuilder();

            for (byte b : hash) {
                String hex = Integer.toHexString(0xff & b);
                if (hex.length() == 1) hexString.append('0');
                hexString.append(hex);
            }

            String passwordHash = hexString.toString();

            boolean userMatch = MessageDigest.isEqual(username.getBytes(), expectedUser.getBytes());
            boolean hashMatch = MessageDigest.isEqual(passwordHash.getBytes(), expectedHash.getBytes());

            return userMatch && hashMatch;

        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
            return false;
        }
    }
}
