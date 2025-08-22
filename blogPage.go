package main

import (
	"io"
	"reflect"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

type blogPage struct {
	viewport viewport.Model
	dump     io.Writer
	renderer *glamour.TermRenderer
	content  string
}

type (
	renderFailedMsg  struct{}
	renderSuccessMsg string
)

func (b blogPage) Init() tea.Cmd {
	return nil
}

func (b blogPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if b.dump != nil {
		spew.Fdump(b.dump, "from blogPage %s", reflect.TypeOf(msg))
	}
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "q", "ctrl+c":
			return b, tea.Quit
		case "esc", "left", "h":
			return b, func() tea.Msg {
				return toggleStateMsg{}
			}
		case "j", "down":
			b.viewport.ScrollDown(1)
		case "k", "up":
			b.viewport.ScrollUp(1)
		case "G":
			b.viewport.GotoBottom()
		case "g":
			b.viewport.GotoTop()
		}
	case tea.WindowSizeMsg:
		b.viewport.Height = message.Height - 7
		b.viewport.Width = message.Width - message.Width/3 - b.viewport.Style.GetHorizontalFrameSize() - 10

		return b, renderWithGlamour(b.content)

	case renderSuccessMsg:
		b.setContent(string(message))
	case renderFailedMsg:
		b.setContent("Render Failed")

	}
	return b, nil
}

func (b *blogPage) setContent(s string) {
	b.viewport.SetContent(s)
}

func (b blogPage) View() string {
	return b.viewport.View()
}

func newViewPort(width, height int) viewport.Model {
	style = lipgloss.NewStyle().Margin(1, 2)
	vp := viewport.New(width, height)
	vp.Style = style
	return vp
}

// this part is heavily inspired by glow
func renderWithGlamour(content string) tea.Cmd {
	return func() tea.Msg {
		s, err := render(content)
		if err != nil {
			return renderFailedMsg{}
		}
		return renderSuccessMsg(s)
	}
}

func render(content string) (string, error) {
	r, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
	if err != nil {
		return "", err
	}
	out, err := r.Render(content)
	if err != nil {
		return "", err
	}
	return out, nil
}

func newBlogPage(content string, dump io.Writer) (*blogPage, error) {
	blogPage := blogPage{}
	blogPage.dump = dump
	blogPage.content = content

	vp := newViewPort(70, 27)

	blogPage.viewport = vp
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(vp.Width-vp.Style.GetHorizontalFrameSize()))
	if err != nil {
		return nil, err
	}

	blogPage.renderer = renderer
	str, err := renderer.Render(blogPage.content)
	if err != nil {
		return nil, err
	}
	vp.SetContent(str)
	return &blogPage, err
}
