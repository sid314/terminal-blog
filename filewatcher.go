package main

import (
	"github.com/fsnotify/fsnotify"
)

type (
	errMsg              struct{ err error }
	listUpdateNeededMsg struct{}
	nilMsg              struct{}
)

func checkFolderUpdates(a *app) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		a.Send(errMsg{err: err})
	}
	defer watcher.Close()
	err = watcher.Add("./posts/")
	if err != nil {
		a.Send(errMsg{err: err})
	}

	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				a.Send(nilMsg{})
			}
			a.Send(listUpdateNeededMsg{})
		case err, ok := <-watcher.Errors:
			if !ok {
				a.Send(nilMsg{})
			}

			a.Send(errMsg{err: err})
		}
	}
}
