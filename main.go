package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m, err := initialBaseModel()
	if err != nil {
		log.Fatal(err)
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	// p := tea.NewProgram(initialList(), tea.WithAltScreen())
	go func() {
		checkFolderUpdates(p)
	}()
	p.Run()
}
