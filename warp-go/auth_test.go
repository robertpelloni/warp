package main

import (
	"os"
	"testing"
)

func TestAuthenticateUser(t *testing.T) {
	os.Setenv("WARP_ADMIN_USER", "admin")
	os.Setenv("WARP_ADMIN_HASH", "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918")

	if !AuthenticateUser("admin", "admin") {
		t.Error("Expected true for admin/admin")
	}
	if AuthenticateUser("user", "pass") {
		t.Error("Expected false for user/pass")
	}
}
