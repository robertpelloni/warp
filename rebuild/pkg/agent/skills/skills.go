package skills

import (
	"fmt"
	"strings"
)

type Skill struct {
	Name, Description string
}

type SkillManager struct {
	skills map[string]*Skill
}

func NewSkillManager() *SkillManager { return &SkillManager{skills: make(map[string]*Skill)} }
func (m *SkillManager) Register(s *Skill) { m.skills[s.Name] = s }
func (m *SkillManager) BuildPrompt() string {
	var b strings.Builder
	b.WriteString("Skills:\n")
	for _, s := range m.skills { b.WriteString(fmt.Sprintf("- %s: %s\n", s.Name, s.Description)) }
	return b.String()
}
