package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	posts []post
	list  list.Model
}

func (m model) Init() tea.Cmd {
	return nil
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
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
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

func initialModael() model {
	items := []list.Item{
		post{
			title:       "test1",
			description: "test1",
		},
		post{
			title:       "test2",
			description: "test2",
		},
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Posts"
	return model{list: list}
}

func main() {
	p := tea.NewProgram(initialModael())
	_, err := p.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
