package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonggulee/awsctx/internal/aws"
)

type accountIDMsg struct {
	index     int
	accountID string
}

type Model struct {
	Profiles []aws.Profile
	Current  string
	cursor   int
	spinner  spinner.Model
	loading  int
}

func NewModel(profiles []aws.Profile, current string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		Profiles: profiles,
		Current:  current,
		spinner:  s,
		loading:  len(profiles),
	}
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick}

	for i, p := range m.Profiles {
		i, name := i, p.Name
		cmds = append(cmds, func() tea.Msg {
			accountID, err := aws.FetchAccountID(name)
			if err != nil {
				accountID = "unknown"
			}
			return accountIDMsg{index: i, accountID: accountID}
		})
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case accountIDMsg:
		m.Profiles[msg.index].AccountID = msg.accountID
		m.loading--
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.loading > 0 {
			return m, nil
		}
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

		accountID := p.AccountID
		if accountID == "" {
			accountID = m.spinner.View()
		}

		s += fmt.Sprintf("%s%s%-30s %s\n", cursor, active, p.Name, accountID)
	}
	s += "\n↑↓/jk: move  enter: select  q: quit"
	return s
}
