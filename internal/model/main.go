package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Main struct {
}

func NewMain() Main {
	return Main{}
}

func (m Main) Init() tea.Cmd {
	return nil
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Main) View() string {
	return "main view"
}
