package simulation

import "sync"

type Entity struct {
	ID int
	X, Y, VX, VY float64
}

type Engine struct {
	Entities []*Entity
	mu sync.Mutex
}

func NewEngine() *Engine { return &Engine{Entities: make([]*Entity, 0)} }
func (e *Engine) AddEntity(ent *Entity) { e.mu.Lock(); defer e.mu.Unlock(); e.Entities = append(e.Entities, ent) }
func (e *Engine) Step(dt float64) {
	e.mu.Lock(); defer e.mu.Unlock()
	for _, ent := range e.Entities { ent.X += ent.VX * dt; ent.Y += ent.VY * dt }
}
