package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type blogPage struct {
	viewport viewport.Model
	renderer *glamour.TermRenderer
	content  string
}

func (b blogPage) Init() tea.Cmd {
	return nil
}

func (b blogPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		str, err := b.renderer.Render(b.content)
		if err != nil {
			return b, func() tea.Msg {
				return fatalErrorMsg{}
			}
		}
		b.viewport.SetContent(str)

		b.viewport, _ = b.viewport.Update(message)
		return b, func() tea.Msg {
			return blogPageUpdateNeededMsg{}
		}

	}
	return b, nil
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

func newRenderer(renderWidth int) (renderer *glamour.TermRenderer, err error) {
	return glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(renderWidth))
}

func newBlogPage(content string) (*blogPage, error) {
	blogPage := blogPage{}
	blogPage.content = content

	vp := newViewPort(70, 27)

	blogPage.viewport = vp
	renderer, err := newRenderer(vp.Width - vp.Style.GetHorizontalFrameSize())
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
