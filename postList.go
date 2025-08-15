package main

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type postList struct {
	posts   []post
	list    list.Model
	dump    io.Writer
	focused post
	index   int
}

type blogPageUpdateNeededMsg struct{}

func (m postList) Init() tea.Cmd {
	return nil
}

func (m postList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, "from postlist %s", msg)
	}
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
		case "j", "down":
			if m.index != len(m.posts)-1 {
				m.index++
				m.focused = m.posts[m.index]
			}

		}
	case tea.WindowSizeMsg:
		h, v := style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-10)
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
			spew.Fdump(m.dump, "updated the items")
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	spew.Fdump(m.dump, "updated the list")
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

func initialList(dump io.Writer) postList {
	var postList postList
	postList.dump = dump
	items, posts, err := addPosts()
	if err != nil {
		log.Fatal(err)
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Posts"
	list.FilterInput.Width = 10
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
