package main

import (
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
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
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

func initialBaseModel() (*baseModel, error) {
	filebuf, err := os.ReadFile("./posts/test1.md")
	if err != nil {
		return nil, err
	}
	initBlog, err := newBlogPage(string(filebuf))
	if err != nil {
		return nil, err
	}
	return &baseModel{
		state:    listView,
		postList: initialList(),
		blogPage: *initBlog,
	}, nil
}
