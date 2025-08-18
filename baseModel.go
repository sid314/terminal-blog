package main

import (
	"io"
	"os"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

type sessionState uint

const (
	listView sessionState = iota
	contentView
)

const (
	minWidth  = 100
	minHeight = 30
)

var (
	style        = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle = lipgloss.NewStyle().Margin(1, 2).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("62"))
)

type baseModel struct {
	state           sessionState
	postList        postList
	blogPage        blogPage
	focused         tea.Model
	dump            io.Writer
	fatalErrorState bool
	tooSmall        bool
}
type updateBlogPageMsg struct {
	path string
}
type (
	toggleStateMsg struct{}
	fatalErrorMsg  struct{}
)

func (b baseModel) sendBlogPageUpdate() tea.Cmd {
	if b.dump != nil {
		spew.Fdump(b.dump, "path", b.postList.focused.path)
	}
	path := b.postList.focused.path
	return func() tea.Msg {
		return updateBlogPageMsg{
			path: path,
		}
	}
}

func (b baseModel) Init() tea.Cmd {
	return nil
}

func (b baseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var all, active, passToList, passToPage bool
	// if all is true then the msg is passed down to both child models
	// otherwise if active is true then it is passed down only to the focused model
	// passToPage and passToList send the messages to page and list respectively
	if b.dump != nil {
		spew.Fdump(b.dump, "from baseModel %s", reflect.TypeOf(msg))
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch message := msg.(type) {

	// TODO: Handle fatalErrorMsg

	case listUpdateNeededMsg:
		passToList = true
	case tea.KeyMsg:
		active = true
		switch message.String() {
		case "q":
			cmd = tea.Quit
			cmds = append(cmds, cmd)
		case "b":
			b.state = contentView
		case "l":
			b.state = listView
		}
	case blogPageUpdateNeededMsg:
		all = true
		cmd = b.sendBlogPageUpdate()
		cmds = append(cmds, cmd)

	case toggleStateMsg:
		if b.state == contentView {
			b.state = listView
		} else {
			b.state = contentView
		}
	case updateBlogPageMsg:
		fileContentBytes, err := os.ReadFile(message.path)
		if err != nil {
			if b.dump != nil {
				spew.Fdump(b.dump, message.path)
				spew.Fdump(b.dump, "quitting due to err %s", err.Error())
			}
			cmds = append(cmds, func() tea.Msg {
				return fatalErrorMsg{}
			})
		}
		rendered, err := b.blogPage.renderer.Render(string(fileContentBytes))
		if err != nil {
			if b.dump != nil {
				spew.Fdump(b.dump, "quitting due to err %s", err.Error())
			}
			cmds = append(cmds, func() tea.Msg {
				return fatalErrorMsg{}
			})
		}
		b.blogPage.viewport.SetContent(rendered)
		passToList = true
	case tea.WindowSizeMsg:
		spew.Fdump(b.dump, message.Width, message.Height)
		b1 := message.Width < minWidth
		b2 := message.Height < minHeight
		spew.Fdump(b.dump, b1, b2)
		if b1 || b2 {
			b.tooSmall = true
		} else {
			b.tooSmall = false
		}

		all = true
	case requestForNewRendererMsg, newRendererMsg:
		passToPage = true

	}
	if all {

		outModel, cmd := b.postList.Update(msg)
		b.postList = outModel.(postList)
		cmds = append(cmds, cmd)
		outModel, cmd = b.blogPage.Update(msg)
		b.blogPage = outModel.(blogPage)
		cmds = append(cmds, cmd)
	} else if active {
		switch b.state {
		case listView:
			outModel, cmd := b.postList.Update(msg)
			b.postList = outModel.(postList)
			cmds = append(cmds, cmd)
		case contentView:
			outModel, cmd := b.blogPage.Update(msg)
			b.blogPage = outModel.(blogPage)
			cmds = append(cmds, cmd)
		}
	} else if passToList {

		outModel, cmd := b.postList.Update(msg)
		b.postList = outModel.(postList)
		cmds = append(cmds, cmd)
	} else if passToPage {

		outModel, cmd := b.blogPage.Update(msg)
		b.blogPage = outModel.(blogPage)
		cmds = append(cmds, cmd)
	}
	return b, tea.Batch(cmds...)
}

func (b baseModel) View() string {
	var s string
	if !b.tooSmall {
		switch b.state {
		case listView:
			s = lipgloss.JoinHorizontal(lipgloss.Top, focusedStyle.Render(b.postList.View()), style.Render(b.blogPage.View()))

		case contentView:
			s = lipgloss.JoinHorizontal(lipgloss.Top, style.Render(b.postList.View()), focusedStyle.Render(b.blogPage.View()))
		}
		return s
	} else {
		return "Window Size too small"
	}
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
