package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonggulee/awsctx/internal/aws"
	"github.com/jonggulee/awsctx/internal/ui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-c" {
		current := aws.LoadCurrentContext()
		if current == "" {
			fmt.Println("No profile selected")
		} else {
			fmt.Printf("\n Current profile: %s\n\n", ui.RenderProfileName(current))
		}
		return
	}

	profiles, err := aws.LoadProfiles()
	if err != nil {
		log.Fatalf("failed to load profiles: %v", err)
	}

	if len(profiles) == 0 {
		fmt.Println("No AWS profiles found in ~/.aws/config")
		os.Exit(1)
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
