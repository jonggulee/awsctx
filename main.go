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
		log.Fatalf("config 파일을 읽을 수 없습니다: %v", err)
	}

	m := ui.NewModel(profiles, aws.LoadCurrentContext())
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("프로그램 실행 중 오류 발생: %v", err)
	}
}
