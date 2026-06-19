package main

import "testing"

func TestAuthenticateUser(t *testing.T) {
	if !AuthenticateUser("admin", "admin") {
		t.Error("Expected true for admin/admin")
	}
	if AuthenticateUser("user", "pass") {
		t.Error("Expected false for user/pass")
	}
}
