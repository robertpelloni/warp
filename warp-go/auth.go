package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func AuthenticateUser(username, password string) bool {
	fmt.Printf("Warp Go: Authenticating user %s\n", username)

	expectedUser := os.Getenv("WARP_ADMIN_USER")
	expectedHash := os.Getenv("WARP_ADMIN_HASH")

	if expectedUser == "" || expectedHash == "" {
		fmt.Println("Warning: Authentication not configured")
		return false
	}

	hasher := sha256.New()
	hasher.Write([]byte(password))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))

	return username == expectedUser && passwordHash == expectedHash
}
