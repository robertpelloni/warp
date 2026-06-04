package simulation

import (
	"fmt"
	"sync"
	"time"
)

// Entity represents an object in the simulation.
type Entity struct {
	ID   int
	Type string
	X, Y float64
	VX, VY float64
	Health float64
}

// Engine handles the simulation loop and entity updates.
type Engine struct {
	Entities []*Entity
	mu       sync.Mutex
}

func NewEngine() *Engine {
	return &Engine{
		Entities: make([]*Entity, 0),
	}
}

func (e *Engine) AddEntity(entity *Entity) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Entities = append(e.Entities, entity)
}

func (e *Engine) Step(dt float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, ent := range e.Entities {
		ent.X += ent.VX * dt
		ent.Y += ent.VY * dt
	}
}

func (e *Engine) Run(steps int, dt float64) {
	for i := 0; i < steps; i++ {
		e.Step(dt)
		e.mu.Lock()
		if len(e.Entities) > 0 {
			fmt.Printf("Step %d: Entity[0] at (%.2f, %.2f)\n", i, e.Entities[0].X, e.Entities[0].Y)
		}
		e.mu.Unlock()
		time.Sleep(time.Duration(dt*1000) * time.Millisecond)
	}
}
