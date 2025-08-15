package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

type (
	errMsg              struct{ err error }
	listUpdateNeededMsg struct{}
	nilMsg              struct{}
)

func checkFolderUpdates(p *tea.Program) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		p.Send(errMsg{err: err})
	}
	defer watcher.Close()
	err = watcher.Add("./posts/")
	if err != nil {
		p.Send(errMsg{err: err})
	}

	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				p.Send(nilMsg{})
			}
			p.Send(listUpdateNeededMsg{})
		case err, ok := <-watcher.Errors:
			if !ok {
				p.Send(nilMsg{})
			}

			p.Send(errMsg{err: err})
		}
	}
}
