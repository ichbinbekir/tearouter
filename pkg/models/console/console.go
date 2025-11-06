package console

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	out      string
	viewport viewport.Model
	input    textinput.Model
}

func New() Model {
	out := "Welcome to the console!\n"

	viewport := viewport.New(100, 1)
	viewport.SetContent(out)
	viewport.Style = viewport.Style.MarginLeft(2)

	input := textinput.New()
	input.SetSuggestions([]string{"clear", "exit"})
	input.ShowSuggestions = true
	input.Focus()

	return Model{input: input, viewport: viewport, out: out}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 8
		m.input.CharLimit = msg.Width - 10
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.out += m.input.Value() + "\n"
			m.viewport.SetContent(m.out)
			m.input.Reset()
			m.viewport.Height = strings.Count(m.out, "\n")
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		m.viewport.View(),
		m.input.View(),
	)
}
