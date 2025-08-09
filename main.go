package main

import (
	"log"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
)

var style = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	posts []post
	list  list.Model
}

func (m model) Init() tea.Cmd {
	return checkFolderUpdates
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		return m, cmd
	case errMsg:
		return m, tea.Quit
	case updateNeededMsg:
		items, err := addPosts()
		if err != nil {
			return m, tea.Quit
		} else {
			m.list.SetItems(items)
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, tea.Batch(cmd, checkFolderUpdates)
}

func (m model) View() string {
	return style.Render(m.list.View())
}

type post struct {
	title, description string
}

func (p post) Title() string {
	return p.title
}

func (p post) Description() string {
	return p.description
}

func (p post) FilterValue() string {
	return p.title
}

func initialModel() model {
	items, err := addPosts()
	if err != nil {
		log.Fatal(err)
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Posts"
	return model{list: list}
}

type (
	errMsg          struct{ err error }
	updateNeededMsg struct{}
)

func checkFolderUpdates() tea.Msg {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errMsg{err: err}
	}
	defer watcher.Close()
	err = watcher.Add("./posts/")
	if err != nil {
		return errMsg{err: err}
	}

	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			return updateNeededMsg{}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}

			return errMsg{err: err}
		}
	}
}

func addPosts() ([]list.Item, error) {
	var items []list.Item
	files, err := os.ReadDir("./posts/")
	for _, v := range files {
		if name := v.Name(); path.Ext(name) == ".md" {
			items = append(items, post{
				title:       v.Name(),
				description: "test description",
			})
		}
	}
	return items, err
}

func main() {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
