package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/robertpelloni/warp-rebuild/pkg/terminal"
)

type model struct {
	content string
	view    string
	session terminal.Session
}

func initialModel(s terminal.Session) model {
	return model{content: "Warp Ultimate LLM Harness Rebuild", view: "terminal", session: s}
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
		case "v":
			if m.view == "terminal" {
				m.view = "agent"
			} else {
				m.view = "terminal"
			}
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
	var body string
	if m.view == "terminal" {
		body = " [ TERMINAL VIEW ]\n" + m.content
	} else {
		body = " [ AGENT VIEW ]\n" + m.content
	}
	return fmt.Sprintf("\n%s\n\nPress 'v' to toggle views, 'q' to exit.\n", style.Render(body))
}
