package session

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager() returned nil")
	}
	if m.Count() != 0 {
		t.Errorf("Count() = %d, want 0", m.Count())
	}
}

func TestCreateSession(t *testing.T) {
	m := NewManager()
	sess, err := m.Create("Test", "cmd.exe", 80, 24, nil)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if sess == nil {
		t.Fatal("Create() returned nil session")
	}
	if sess.Name != "Test" {
		t.Errorf("Name = %q, want 'Test'", sess.Name)
	}
	if !sess.Active {
		t.Error("Session should be active")
	}
	if m.Count() != 1 {
		t.Errorf("Count() = %d, want 1", m.Count())
	}
}

func TestActiveSession(t *testing.T) {
	m := NewManager()
	m.Create("First", "cmd.exe", 80, 24, nil)
	m.Create("Second", "cmd.exe", 80, 24, nil)

	active := m.Active()
	if active == nil {
		t.Fatal("Active() returned nil")
	}
	if active.Name != "Second" {
		t.Errorf("Active() Name = %q, want 'Second'", active.Name)
	}
}

func TestSetActive(t *testing.T) {
	m := NewManager()
	sess1, _ := m.Create("First", "cmd.exe", 80, 24, nil)
	m.Create("Second", "cmd.exe", 80, 24, nil)

	m.SetActive(sess1.ID)
	active := m.Active()
	if active.Name != "First" {
		t.Errorf("After SetActive, Active() Name = %q, want 'First'", active.Name)
	}
}

func TestRemoveSession(t *testing.T) {
	m := NewManager()
	sess, _ := m.Create("Test", "cmd.exe", 80, 24, nil)
	m.Remove(sess.ID)
	if m.Count() != 0 {
		t.Errorf("Count() after Remove = %d, want 0", m.Count())
	}
}

func TestListSessions(t *testing.T) {
	m := NewManager()
	m.Create("A", "cmd.exe", 80, 24, nil)
	m.Create("B", "cmd.exe", 80, 24, nil)

	list := m.List()
	if len(list) != 2 {
		t.Errorf("List() len = %d, want 2", len(list))
	}
}

func TestCreateDefaultName(t *testing.T) {
	m := NewManager()
	sess, err := m.Create("", "cmd.exe", 80, 24, nil)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	if sess.Name == "" {
		t.Error("Name should not be empty when default is applied")
	}
}
