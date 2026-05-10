package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/tearouter"
)

type Page struct {
	title string
}

func (p Page) Init() tea.Cmd { return nil }
func (p Page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "q", "ctrl+c":
			return p, tea.Quit
		case "1":
			return p, tearouter.Redirect(tearouter.Go, "/main")
		case "2":
			return p, tearouter.Redirect(tearouter.Go, "/main/settings")
		case "3":
			return p, tearouter.Redirect(tearouter.Go, "/main/settings/profile")
		case "b":
			return p, tearouter.Redirect(tearouter.Pop)
		}
	}
	return p, nil
}
func (p Page) View() string {
	return fmt.Sprintf("Current Page: %s\n\nPress 1 for /main (Go)\nPress 2 for /main/settings (Go)\nPress 3 for /main/settings/profile (Go)\nPress b for Back (Pop)\nPress q to Quit", p.title)
}

func main() {
	model := tearouter.Model{
		InitialRoute: "/main",
		Routes: []tearouter.Route{
			{
				Path: "/main",
				Builder: func() tea.Model { return Page{title: "Main Page"} },
				Children: []tearouter.Route{
					{
						Path: "settings",
						Builder: func() tea.Model { return Page{title: "Settings Page"} },
						Children: []tearouter.Route{
							{
								Path: "profile",
								Builder: func() tea.Model { return Page{title: "Profile Page"} },
							},
						},
					},
				},
			},
		},
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
