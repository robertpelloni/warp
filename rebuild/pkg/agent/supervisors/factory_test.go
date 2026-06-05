package supervisors

import (
	"testing"
	"github.com/robertpelloni/warp-rebuild/pkg/harness"
)

func TestSupervisor(t *testing.T) {
	mock := &harness.MockProvider{Responses: []*harness.LLMResponse{{Content: "ok"}}}
	s, _ := CreateSupervisor("default", mock)
	res, _ := s.Route("hi")
	if res != "ok" { t.Error("route fail") }
}
