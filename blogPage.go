package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type blogPage struct {
	title    string
	path     string
	viewport viewport.Model
}

func (b blogPage) Init() tea.Cmd {
	return nil
}

func (b blogPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return b, tea.Quit
		}
	}
	return b, nil
}

func (b blogPage) View() string {
	return b.viewport.View()
}

func newBlogPage(fileContent string) (*blogPage, error) {
	vp := viewport.New(70, 20)
	vp.Style = style
	const glamourGutter = 2
	glamourRenderWidth := 70 - vp.Style.GetHorizontalFrameSize() - glamourGutter
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth))
	if err != nil {
		return nil, err
	}
	str, err := renderer.Render(fileContent)
	if err != nil {
		return nil, err
	}
	vp.SetContent(str)
	return &blogPage{
		viewport: vp,
	}, nil
}
