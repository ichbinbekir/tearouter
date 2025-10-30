package tearouter

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

type MiddlewareFunc func(string) string

type Model struct {
	Routes       []Route
	InitialRoute string
	Middleware   MiddlewareFunc
	routeStack   []tea.Model
}

func (m Model) Init() tea.Cmd {
	return Redirect(Go, m.InitialRoute)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(RedirectMsg); ok {
		switch msg.Type {
		case Go:
			if m.Middleware != nil {
				if target := m.Middleware(msg.Target); target != msg.Target {
					Redirect(msg.Type, target)
					break
				}
			}
			for _, route := range m.Routes {
				if route.Path == msg.Target {
					m.routeStack = []tea.Model{route.Builder()}
				}
			}
		case Push:
			if m.Middleware != nil {
				if target := m.Middleware(msg.Target); target != msg.Target {
					Redirect(msg.Type, target)
					break
				}
			}
			for _, route := range m.Routes {
				if route.Path == msg.Target {
					m.routeStack = append(m.routeStack, route.Builder())
				}
			}
		case Replace:
			if m.Middleware != nil {
				if target := m.Middleware(msg.Target); target != msg.Target {
					Redirect(msg.Type, target)
					break
				}
			}
			for _, route := range m.Routes {
				if route.Path == msg.Target {
					if length := len(m.routeStack); length > 0 {
						m.routeStack[length-1] = route.Builder()
					} else {
						return m, errCmd(errors.New("replace error"))
					}
				}
			}
		case Pop:
			if length := len(m.routeStack); length > 1 {
				m.routeStack = m.routeStack[:length-1]
			} else {
				cmd = errCmd(errors.New("pop error"))
			}
		}
	}

	if length := len(m.routeStack); length > 0 {
		var cmdx tea.Cmd
		m.routeStack[length-1], cmdx = m.routeStack[length-1].Update(msg)
		cmd = tea.Batch(cmd, cmdx)
	} else {
		return m, tea.Batch(cmd, errCmd(errors.New("router cant find route, you should go redirect anyware")))
	}
	return m, cmd
}

func (m Model) View() string {
	if length := len(m.routeStack); length > 0 {
		return m.routeStack[length-1].View()
	}
	return "TEA ROUTER STACK CAN'T BE EMPTY, YOU SHOULD GO REDIRECT ANYWARE"
}
