package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/tearouter/internal/models/layout"
)

func main() {
	if _, err := tea.NewProgram(layout.Base(), tea.WithAltScreen()).Run(); err != nil {
		panic(err)
	}
}
