package main

import (
	"io"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/davecgh/go-spew/spew"
)

type blogPage struct {
	title    string
	path     string
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
		case "esc":
			return b, func() tea.Msg {
				return toggleStateMsg{}
			}
		}
	case updateBlogPageMsg:
		fileContentBytes, err := os.ReadFile(message.path)
		if err != nil {
			if b.dump != nil {
				spew.Fdump(b.dump, message.path)
				spew.Fdump(b.dump, "quitting due to err %s", err.Error())
			}
			return b, func() tea.Msg {
				return fatalErrorMsg{}
			}
		}
		rendered, err := b.renderer.Render(string(fileContentBytes))
		if err != nil {
			if b.dump != nil {
				spew.Fdump(b.dump, "quitting due to err %s", err.Error())
			}
			return b, func() tea.Msg {
				return fatalErrorMsg{}
			}
		}
		b.viewport.SetContent(rendered)

	}
	return b, nil
}

func (b blogPage) View() string {
	return b.viewport.View()
}

func newBlogPage(content string, dump io.Writer) (*blogPage, error) {
	blogPage := blogPage{}
	blogPage.dump = dump

	vp := viewport.New(70, 20)
	vp.Style = style

	blogPage.viewport = vp
	const glamourGutter = 2
	glamourRenderWidth := 70 - vp.Style.GetHorizontalFrameSize() - glamourGutter
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
