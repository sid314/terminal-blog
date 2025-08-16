package main

import (
	"io"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/davecgh/go-spew/spew"
)

type blogPage struct {
	viewport viewport.Model
	dump     io.Writer
	renderer *glamour.TermRenderer
}

func (b blogPage) Init() tea.Cmd {
	return nil
}

func (b blogPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if b.dump != nil {
		spew.Fdump(b.dump, "from blogPage %s", msg)
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
	}
	return b, nil
}

func (b blogPage) View() string {
	return b.viewport.View()
}

func newBlogPage(content string, dump io.Writer) (*blogPage, error) {
	blogPage := blogPage{}
	blogPage.dump = dump

	vp := viewport.New(80, 32)
	vp.Style = style

	blogPage.viewport = vp
	glamourRenderWidth := vp.Width
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth))
	if err != nil {
		return nil, err
	}

	blogPage.renderer = renderer
	str, err := renderer.Render(content)
	if err != nil {
		return nil, err
	}
	vp.SetContent(str)
	return &blogPage, err
}
