package tearouter

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Route struct {
	Path    string
	Builder func() tea.Model
}

type Model struct {
	InitialRoute string
	Routes       []Route
	modelStack   []tea.Model
}

func (m Model) Init() tea.Cmd {
	return m.routeInitial()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(RedirectMsg); ok {
		switch msg.Type {
		case Go:
			m.gox(msg.Target)
		case Push:
			m.push(msg.Target)
		case Replace:
			if err := m.replace(msg.Target); err != nil {
				//TODO err
			}
		case Pop:
			if err := m.pop(); err != nil {
				//TODO err
			}
		}
	}

	if length := len(m.modelStack); length > 0 {
		var cmdx tea.Cmd
		m.modelStack[length-1], cmdx = m.modelStack[length-1].Update(msg)
		cmd = tea.Batch(cmd, cmdx)
	} else {
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		//TODO err
	}
	return m, cmd
}

func (m Model) View() string {
	if length := len(m.modelStack); length > 0 {
		return m.modelStack[length-1].View()
	}
	//TODO err
	return "TEA ROUTER STACK CAN'T BE EMPTY, YOU SHOULD GO REDIRECT ANYWARE"
}

func (m *Model) gox(target string) {
	for _, route := range m.Routes {
		if route.Path == target {
			m.modelStack = []tea.Model{route.Builder()}
		}
	}
}

func (m *Model) push(target string) {
	for _, route := range m.Routes {
		if route.Path == target {
			m.modelStack = append(m.modelStack, route.Builder())
		}
	}
}

func (m *Model) replace(target string) error {
	for _, route := range m.Routes {
		if route.Path == target {
			if length := len(m.modelStack); length > 0 {
				m.modelStack[length-1] = route.Builder()
				return nil
			}
			//TODO ret err
		}
	}
	return nil
}

func (m *Model) pop() error {
	if length := len(m.modelStack); length > 1 {
		m.modelStack = m.modelStack[:length-1]
		return nil
	}
	// TODO ret err
	return nil
}

func (m Model) routeInitial() tea.Cmd {
	if m.InitialRoute == "" {
		m.InitialRoute = "/"
	}
	return Redirect(Go, m.InitialRoute)
}
