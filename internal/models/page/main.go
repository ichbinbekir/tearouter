package page

import (
	tea "github.com/charmbracelet/bubbletea"
)

type main struct {
}

func (m main) Init() tea.Cmd {
	return nil
}

func (m main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	}

	return m, nil
}

func (m main) View() string {
	return "asdasd"
}
