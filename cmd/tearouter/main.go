package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/tearouter/internal/model"
)

func main() {
	p := tea.NewProgram(model.NewMain(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
