package orchestration

import "testing"

func TestRunner(t *testing.T) {
	r := NewRunner()
	if r.Status != "idle" {
		t.Error("Initial status should be idle")
	}

	err := r.RunWorkflow("test-flow")
	if err != nil {
		t.Fatalf("Workflow failed: %v", err)
	}

	if r.Status != "completed" {
		t.Error("Status should be completed")
	}
}
