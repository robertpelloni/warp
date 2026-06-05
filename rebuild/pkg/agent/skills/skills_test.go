package skills

import (
	"strings"
	"testing"
)

func TestSkills(t *testing.T) {
	m := NewSkillManager()
	m.Register(&Skill{Name: "debug", Description: "fix stuff"})
	p := m.BuildPrompt()
	if !strings.Contains(p, "debug") { t.Error("missing debug skill") }
}
