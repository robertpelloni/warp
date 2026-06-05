package simulation

import "testing"

func TestEngine(t *testing.T) {
	e := NewEngine()
	ent := &Entity{ID: 1, VX: 10}
	e.AddEntity(ent)
	e.Step(1.0)
	if ent.X != 10 { t.Errorf("expected 10, got %f", ent.X) }
}
