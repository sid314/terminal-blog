package main

import (
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

type sessionState uint

const (
	listView sessionState = iota
	contentView
)

var (
	style        = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle = lipgloss.NewStyle().Margin(1, 2).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("62"))
)

type baseModel struct {
	state    sessionState
	postList postList
	blogPage blogPage
	focused  tea.Model
	dump     io.Writer
}

func (b baseModel) Init() tea.Cmd {
	return nil
}

func (b baseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if b.dump != nil {
		spew.Fdump(b.dump, "from basemodel %s", msg)
	}
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
		s = lipgloss.JoinHorizontal(lipgloss.Top, focusedStyle.Render(b.postList.View()), style.Render(b.blogPage.View()))

	case contentView:
		s = lipgloss.JoinHorizontal(lipgloss.Top, style.Render(b.postList.View()), focusedStyle.Render(b.blogPage.View()))
	}
	return s
}

func initialBaseModel(dump io.Writer) (*baseModel, error) {
	filebuf, err := os.ReadFile("./posts/test1.md")
	if err != nil {
		return nil, err
	}
	initBlog, err := newBlogPage(string(filebuf), dump)
	if err != nil {
		return nil, err
	}
	return &baseModel{
		state:    listView,
		postList: initialList(dump),
		blogPage: *initBlog,
		dump:     dump,
	}, nil
}
