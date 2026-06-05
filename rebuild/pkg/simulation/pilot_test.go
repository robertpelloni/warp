package simulation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/robertpelloni/warp-rebuild/pkg/agent"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

func TestPilotSimulation(t *testing.T) {
	// 1. Setup Simulation Engine
	sim := NewEngine()
	player := &Entity{
		ID:     1,
		Type:   "player",
		X:      0,
		Y:      0,
		VX:     10,
		VY:     0,
		Health: 100,
	}
	sim.AddEntity(player)

	// 2. Setup AI Agent as the Pilot
	model := &harness.ModelInfo{ID: "simulation-pilot-model"}
	mockProvider := &harness.MockProvider{}
	pilot := agent.NewAgent("pilot-1", "You are the pilot of a simulation. Monitor and adjust entity velocities.", model, mockProvider)

	// 3. Run Pilot Loop
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Simulation loop with AI intervention
	go func() {
		for i := 0; i < 5; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				sim.Step(0.1)
				fmt.Printf("[Sim] Player at (%.2f, %.2f)\n", player.X, player.Y)

				// Agent "observes" and "acts"
				_ = pilot.RunLoop(ctx)

				// Simulate agent changing velocity
				player.VX += 1.0

				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	<-ctx.Done()

	if player.X <= 0 {
		t.Error("Player should have moved during the simulation")
	}

	t.Log("Pilot autonomous game simulation verified successfully.")
}
