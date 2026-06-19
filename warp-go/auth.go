package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// In a real application, this would be retrieved from a database.
const dummyHash = "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918" // sha256 of "admin"

func AuthenticateUser(username, password string) bool {
	fmt.Printf("Warp Go: Authenticating user %s\n", username)

	hasher := sha256.New()
	hasher.Write([]byte(password))
	passwordHash := hex.EncodeToString(hasher.Sum(nil))

	return username == "admin" && passwordHash == dummyHash
}
