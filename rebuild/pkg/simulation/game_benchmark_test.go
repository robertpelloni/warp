package simulation

import (
	"context"
	"testing"
)

func BenchmarkAutonomousGamePlay(b *testing.B) {
	engine := NewEngine()
	engine.AddEntity(&Entity{ID: 1, X: 0, Y: 0, VX: 1, VY: 1})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.RunBenchmark(ctx, 1)
	}
}
