package terminal

import (
	"testing"
)

type MockSession struct {
	BaseSession
}

func (m *MockSession) Read(p []byte) (n int, err error)  { return 0, nil }
func (m *MockSession) Write(p []byte) (n int, err error) { return 0, nil }
func (m *MockSession) Close() error                     { return nil }
func (m *MockSession) Resize(cols, rows int) error      { return nil }
func (m *MockSession) GetWorkingDirectory() (string, error) { return "", nil }

func TestTabManager(t *testing.T) {
	tm := NewTabManager()
	session := &MockSession{}

	tab := tm.AddTab("Test Tab", session)
	if tab.Title != "Test Tab" {
		t.Errorf("Expected tab title 'Test Tab', got %s", tab.Title)
	}

	if tm.GetActiveTab() != tab {
		t.Error("Expected active tab to be the one just added")
	}

	tm.AddTab("Second Tab", session)
	if len(tm.GetTabs()) != 2 {
		t.Errorf("Expected 2 tabs, got %d", len(tm.GetTabs()))
	}
}
