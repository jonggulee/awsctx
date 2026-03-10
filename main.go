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
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h", "--help":
			printHelp()
		case "-c":
			current := aws.LoadCurrentContext()
			if current == "" {
				fmt.Println("No profile selected")
			} else {
				fmt.Printf("\n Current profile: %s\n\n", ui.RenderProfileName(current))
			}
		default:
			profileName := os.Args[1]
			profiles, err := aws.LoadProfiles()
			if err != nil {
				log.Fatalf("failed to load profiles: %v", err)
			}

			for _, p := range profiles {
				if p.Name == profileName {
					err := aws.SaveCurrentContext(profileName)
					if err != nil {
						log.Fatalf("failed to save context: %v", err)
					}
					fmt.Print(ui.SwitchedMessage(profileName))
					return
				}
			}
			fmt.Fprintf(os.Stderr, "error: profile %q not found\n", profileName)
			os.Exit(1)
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

func printHelp() {
	fmt.Println(`Usage: awsctx [flags] [profile]

  awsctx               Launch interactive TUI to select a profile
  awsctx <profile>     Switch directly to the specified profile
  awsctx -c            Show the current active profile
  awsctx -h            Show this help message`)
}
