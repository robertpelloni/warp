package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"fmt"
)

type model struct {
	content string
}

func initialModel() model {
	return model{content: "Warp Ultimate LLM Harness Rebuild\nInitializing TUI..."}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Padding(1, 2)

func (m model) View() string {
	return fmt.Sprintf("\n%s\n\nPress 'q' to exit back to terminal.\n", style.Render(m.content))
}
