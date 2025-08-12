package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	listView sessionState = iota
	contentView
)

type baseModel struct {
	state    sessionState
	postList postList
	blogPage blogPage
}

func (b baseModel) Init() tea.Cmd {
	return nil
}

func (b baseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			cmd = tea.Quit
			cmds = append(cmds, cmd)
		case "b":
			b.state = contentView
		case "l":
			b.state = listView
		}
	}
	switch b.state {
	case listView:
		_, cmd := b.postList.Update(msg)
		cmds = append(cmds, cmd)
	case contentView:
		_, cmd := b.blogPage.Update(msg)
		cmds = append(cmds, cmd)
	}
	return b, tea.Batch(cmds...)
}

func (b baseModel) View() string {
	var s string
	switch b.state {
	case listView:
		s = b.postList.View()
	case contentView:
		s = b.blogPage.View()
	}
	return s
}

func main() {
	filebuf, err := os.ReadFile("./posts/test1.md")
	if err != nil {
		log.Fatalln(err)
	}
	initBlog, err := newBlogPage(string(filebuf))
	if err != nil {
		log.Fatalln(err)
	}
	p := tea.NewProgram(baseModel{
		state:    listView,
		postList: initialList(),
		blogPage: *initBlog,
	}, tea.WithAltScreen())
	// p := tea.NewProgram(initialList(), tea.WithAltScreen())
	go func() {
		checkFolderUpdates(p)
	}()
	p.Run()
}
