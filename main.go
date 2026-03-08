package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		refreshTick(),
		initialLoad(),
	)
}

func main() {
	ctx := context.Background()

	m, err := initialModel(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
