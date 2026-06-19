package dev.warp;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

public class Auth {
    // In a real application, this would be retrieved from a database.
    private static final String DUMMY_HASH = "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918";

    public static boolean authenticateUser(String username, String password) {
        System.out.println("Warp Java: Authenticating user " + username);

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
            return "admin".equals(username) && DUMMY_HASH.equals(passwordHash);

        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
            return false;
        }
    }
}
