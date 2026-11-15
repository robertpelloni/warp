package workspace

import "testing"

func TestWorkspaceManager(t *testing.T) {
	mgr := NewManager()
	caps := []string{"session_replay", "submodule_intelligence"}
	w := mgr.CreateWorkspace("w1", caps)

	if w.ID != "w1" {
		t.Error("ID mismatch")
	}

	ret := mgr.GetWorkspace("w1")
	if len(ret.Capabilities) != 2 {
		t.Errorf("Expected 2 caps, got %d", len(ret.Capabilities))
	}
}
