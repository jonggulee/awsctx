package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonggulee/awsctx/internal/aws"
)

type Model struct {
	Profiles []aws.Profile
	Current  string
	cursor   int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.Profiles)-1 {
				m.cursor++
			}
		case "enter":
			selected := m.Profiles[m.cursor].Name
			aws.SaveCurrentContext(selected)
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := "AWS Profiles:\n\n"
	for i, p := range m.Profiles {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		active := "  "
		if p.Name == m.Current {
			active = "* "
		}

		s += fmt.Sprintf("%s%s%-30s %s\n", cursor, active, p.Name, p.AccountID)
	}
	s += "\n↑↓: 이동  q: 종료"
	return s
}
