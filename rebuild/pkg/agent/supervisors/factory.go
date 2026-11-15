package supervisors

import (
	"fmt"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

type Supervisor interface {
	Route(prompt string) (string, error)
}

type DefaultSupervisor struct {
	Provider harness.Provider
}

func (s *DefaultSupervisor) Route(prompt string) (string, error) {
	resp, err := s.Provider.Chat(nil, []harness.Message{{Role: "user", Content: prompt}}, nil)
	if err != nil { return "", err }
	return resp.Content, nil
}

func CreateSupervisor(t string, p harness.Provider) (Supervisor, error) {
	if t == "default" { return &DefaultSupervisor{Provider: p}, nil }
	return nil, fmt.Errorf("bad type")
}
