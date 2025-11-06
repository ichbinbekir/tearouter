package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ichbinbekir/tearouter/pkg/models/console"
)

type base struct {
	debug   bool
	console console.Model
	router  tea.Model
}

var width int

func (m base) Init() tea.Cmd {
	return tea.Batch(m.router.Init())
}

func (m base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	if m.debug {
		m.console, cmd = m.console.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width = msg.Width - 4
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "\"":
			m.debug = !m.debug
			cmds = append(cmds, tea.WindowSize())
		}
	}

	m.router, cmd = m.router.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("#FFD700"))

func (m base) View() string {
	out := style.Render(m.router.View()) + "\n"

	if m.debug {
		out += lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).MarginLeft(2).PaddingLeft(1).Width(width).BorderForeground(lipgloss.Color("#FFD700")).Render(m.console.View())
	}

	return out
}
