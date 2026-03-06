package main

import (
	"log"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonggulee/awsctx/internal/aws"
	"github.com/jonggulee/awsctx/internal/ui"
)

func main() {
	profiles, err := aws.LoadProfiles()
	if err != nil {
		log.Fatalf("config 파일을 읽을 수 없습니다: %v", err)
	}

	var wg sync.WaitGroup
	for i, p := range profiles {
		wg.Add(1)
		go func(i int, name string) {
			defer wg.Done()
			accountID, err := aws.FetchAccountID(name)
			if err != nil {
				accountID = "unknown"
			}
			profiles[i].AccountID = accountID
		}(i, p.Name)
	}
	wg.Wait()

	m := ui.Model{Profiles: profiles, Current: os.Getenv("AWS_PROFILE")}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("프로그램 실행 중 오류 발생: %v", err)
	}
}
