package tearouter

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Middleware func(targetPath string) (newPath string)

type Route struct {
	Path    string
	Builder func() tea.Model
}

type Model struct {
	InitialRoute string
	Routes       []Route
	Middleware   Middleware
	modelStack   []tea.Model
}

func (m Model) Init() tea.Cmd {
	return m.routeInitial()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(RedirectMsg); ok {
		if m.Middleware != nil && msg.Type != Pop {
			newTarget := m.Middleware(msg.Target)
			if newTarget != "" {
				msg.Type = Go
				msg.Target = newTarget
			}
		}

		switch msg.Type {
		case Go:
			cmd = m.gox(msg.Target)
		case Push:
			cmd = m.push(msg.Target)
		case Replace:
			cmd = m.replace(msg.Target)
		case Pop:
			cmd = m.pop()
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
		cmd = newErrorCmd(errors.New("router stack is empty, no model to update"))
	}
	return m, cmd
}

func (m Model) View() string {
	if length := len(m.modelStack); length > 0 {
		return m.modelStack[length-1].View()
	}
	return "TEA ROUTER STACK CAN'T BE EMPTY, YOU SHOULD GO REDIRECT ANYWARE"
}

func (m *Model) gox(target string) tea.Cmd {
	for _, route := range m.Routes {
		if route.Path == target {
			newModel := route.Builder()
			m.modelStack = []tea.Model{newModel}
			return newModel.Init()
		}
	}
	return newErrorCmd(fmt.Errorf("route not found: %s", target))
}

func (m *Model) push(target string) tea.Cmd {
	for _, route := range m.Routes {
		if route.Path == target {
			newModel := route.Builder()
			m.modelStack = append(m.modelStack, newModel)
			return newModel.Init()
		}
	}
	return newErrorCmd(fmt.Errorf("route not found: %s", target))
}

func (m *Model) replace(target string) tea.Cmd {
	if len(m.modelStack) == 0 {
		return newErrorCmd(errors.New("cannot replace on an empty stack"))
	}
	for _, route := range m.Routes {
		if route.Path == target {
			newModel := route.Builder()
			m.modelStack[len(m.modelStack)-1] = newModel
			return newModel.Init()
		}
	}
	return newErrorCmd(fmt.Errorf("route not found: %s", target))
}

func (m *Model) pop() tea.Cmd {
	if length := len(m.modelStack); length > 1 {
		m.modelStack = m.modelStack[:length-1]
		return nil
	}
	return newErrorCmd(errors.New("cannot pop from the root of the stack"))
}

func (m *Model) routeInitial() tea.Cmd {
	if m.InitialRoute == "" {
		m.InitialRoute = "/"
	}
	return Redirect(Go, m.InitialRoute)
}
