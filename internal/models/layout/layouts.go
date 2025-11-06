package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/tearouter/internal/router"
	"github.com/ichbinbekir/tearouter/pkg/models/console"
)

func Base() tea.Model {
	return base{console: console.New(), router: router.Model()}
}
