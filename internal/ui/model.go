package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	Switched string
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#25A065")).MarginBottom(1)
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#E06C75")).Bold(true)
	activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).Background(lipgloss.Color("#25A065"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

func SwitchedMessage(profile string) string {
	return fmt.Sprintf("\n Switched to profile: %s\n\n", RenderProfileName(profile))
}

func RenderProfileName(name string) string {
	return activeStyle.Render(name)
}

func NewModel(profiles []aws.Profile, current string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	cursor := 0
	for i, p := range profiles {
		if p.Name == current {
			cursor = i
			break
		}
	}

	return Model{
		Profiles: profiles,
		Current:  current,
		spinner:  s,
		loading:  len(profiles),
		cursor:   cursor,
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
			m.Switched = selected
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := titleStyle.Render("AWS Profiles") + "\n\n"
	for i, p := range m.Profiles {
		cursor := "  "
		if m.cursor == i {
			cursor = cursorStyle.Render("> ")
		}

		active := "  "
		if p.Name == m.Current {
			active = activeStyle.Render("* ")
		}

		accountID := p.AccountID
		if accountID == "" {
			accountID = dimStyle.Render(m.spinner.View())
		}

		line := fmt.Sprintf("%-20s %-20s %s", p.Name, p.Region, accountID)
		if m.cursor == i {
			line = selectedStyle.Render(line)
		}
		s += fmt.Sprintf("%s%s%s\n", cursor, active, line)
	}
	s += helpStyle.Render("\n↑↓/jk: move  enter: select  q: quit")
	return s
}
