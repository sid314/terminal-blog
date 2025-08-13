package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var dump *os.File
	a := 1
	os.Setenv("DEBUG", "true")
	if _, ok := os.LookupEnv("DEBUG"); ok {
		a = 2
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}
	println(a)
	m, err := initialBaseModel(dump)
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
