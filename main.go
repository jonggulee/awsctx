package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonggulee/awsctx/internal/aws"
	"github.com/jonggulee/awsctx/internal/ui"
)

func main() {
	profiles, err := aws.LoadProfiles()
	if err != nil {
		log.Fatalf("failed to load profiles: %v", err)
	}

	m := ui.NewModel(profiles, aws.LoadCurrentContext())
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		log.Fatalf("failed to run program: %v", err)
	}

	m, ok := result.(ui.Model)
	if ok && m.Switched != "" {
		fmt.Print(ui.SwitchedMessage(m.Switched))
	}
}
