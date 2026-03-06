package main

import (
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
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to run program: %v", err)
	}
}
