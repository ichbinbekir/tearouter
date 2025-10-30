package main

import (
	"fmt"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/tearouter"
)

type Model1 struct {
	number int
}
type Model2 struct {
	number int
}

func (m Model1) Init() tea.Cmd {
	return nil
}
func (m Model2) Init() tea.Cmd {
	return nil
}

func (m Model1) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "g":
			return m, tearouter.Redirect(tearouter.Go, "/model2")
		case "p":
			return m, tearouter.Redirect(tearouter.Push, "/model2")
		case "o":
			return m, tearouter.Redirect(tearouter.Pop)
		}
	}
	return m, nil
}
func (m Model2) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "g":
			return m, tearouter.Redirect(tearouter.Go, "/model1")
		case "p":
			return m, tearouter.Redirect(tearouter.Push, "/model1")
		case "o":
			return m, tearouter.Redirect(tearouter.Pop)
		}
	}
	return m, nil
}

func (m Model1) View() string {
	return fmt.Sprintf("Model1 %d", m.number)
}
func (m Model2) View() string {
	return fmt.Sprintf("Model2 %d", m.number)
}

func model() tea.Model {
	return tearouter.Model{
		InitialRoute: "/model1",
		Routes: []tearouter.Route{
			{
				Path:    "/model1",
				Builder: func() tea.Model { return Model1{number: rand.Int()} },
			},
			{
				Path:    "/model2",
				Builder: func() tea.Model { return Model2{number: rand.Int()} },
			},
		},
	}
}

func main() {
	p := tea.NewProgram(model())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
