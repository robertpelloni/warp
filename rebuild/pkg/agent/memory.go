package agent

import (
	"context"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

type MemoryProvider interface {
	Store(ctx context.Context, msg harness.Message) error
	Retrieve(ctx context.Context, query string) ([]harness.Message, error)
}

type LocalMemory struct {
	history []harness.Message
}

func NewLocalMemory() *LocalMemory {
	return &LocalMemory{}
}

func (m *LocalMemory) Store(ctx context.Context, msg harness.Message) error {
	m.history = append(m.history, msg)
	return nil
}

func (m *LocalMemory) Retrieve(ctx context.Context, query string) ([]harness.Message, error) {
	return m.history, nil
}
