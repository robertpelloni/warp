package terminal

import "fmt"

type Pane struct {
	ID      string
	Session Session
}

type Layout struct {
	Panes []*Pane
}

func (l *Layout) SplitHorizontal(parentID string, newSession Session) (*Pane, error) {
	newPane := &Pane{ID: fmt.Sprintf("pane-%d", len(l.Panes)), Session: newSession}
	l.Panes = append(l.Panes, newPane)
	return newPane, nil
}

func (l *Layout) SplitVertical(parentID string, newSession Session) (*Pane, error) {
	newPane := &Pane{ID: fmt.Sprintf("pane-%d", len(l.Panes)), Session: newSession}
	l.Panes = append(l.Panes, newPane)
	return newPane, nil
}
