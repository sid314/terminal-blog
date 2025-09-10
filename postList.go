package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type postList struct {
	posts   []post
	list    list.Model
	focused post
	index   int
}

type blogPageUpdateNeededMsg struct{}

func (m postList) Init() tea.Cmd {
	return nil
}

func (m postList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(func() tea.Msg {
				return blogPageUpdateNeededMsg{}
			}, func() tea.Msg {
				return toggleStateMsg{}
			})
		case "right", "l":
			return m, func() tea.Msg {
				return toggleStateMsg{}
			}

		case "j", "down":
			if m.index != len(m.posts)-1 {
				m.index++
			}

			m.focused = m.posts[m.index]
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, tea.Batch(func() tea.Msg {
				return blogPageUpdateNeededMsg{}
			}, cmd)
		case "k", "up":
			if m.index != 0 {
				m.index--
			}
			m.focused = m.posts[m.index]
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, tea.Batch(func() tea.Msg {
				return blogPageUpdateNeededMsg{}
			}, cmd)
			// case "tab":
			// 	if m.list.Paginator.Page == m.list.Paginator.TotalPages-1 {
			// 		m.list.Paginator.Page = 0
			// 	} else {
			// 		m.list.NextPage()
			// 	}
			// 	m.focused = m.posts[m.list.GlobalIndex()]

		}
	case tea.WindowSizeMsg:
		// h, v := style.GetFrameSize()
		m.list.SetSize(msg.Width/3, msg.Height-5)
		m.list.FilterInput.Width = 10
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		return m, cmd
	case errMsg:
		return m, tea.Quit
	case listUpdateNeededMsg:
		items, posts, err := addPosts()
		if err != nil {
			return m, tea.Quit
		} else {
			m.list.SetItems(items)
			m.posts = posts
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m postList) View() string {
	return style.Render(m.list.View())
}

type post struct {
	title, description, path string
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

func initialList() postList {
	var postList postList
	items, posts, err := addPosts()
	if err != nil {
		log.Fatal(err)
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Posts"
	list.FilterInput.Width = 10
	list.SetShowHelp(false)
	postList.list = list
	postList.posts = posts
	postList.focused = postList.posts[postList.index]
	return postList
}

func addPosts() ([]list.Item, []post, error) {
	var items []list.Item
	var posts []post
	files, err := os.ReadDir("./posts/")
	for _, v := range files {
		if name := v.Name(); path.Ext(name) == ".md" {
			info, _ := v.Info()

			items = append(items, post{
				title:       name,
				description: info.ModTime().String(),
				path:        filepath.Join("./posts/", name),
			})
			posts = append(posts, post{
				title:       name,
				description: info.ModTime().String(),
				path:        filepath.Join("./posts/", name),
			})
		}
	}
	return items, posts, err
}
