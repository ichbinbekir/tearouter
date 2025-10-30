package tearouter

import tea "github.com/charmbracelet/bubbletea"

type Route struct {
	Path    string
	Builder func() tea.Model
}
