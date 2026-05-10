package tearouter

import (
	"errors"
	"fmt"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Middleware func(targetPath string) (newPath string)

type Route struct {
	Path     string
	Builder  func() tea.Model
	Children []Route
}

type Model struct {
	InitialRoute string
	Routes       []Route
	Middleware   Middleware
	modelStack   []tea.Model
	lastSize     tea.WindowSizeMsg
}

func (m Model) Init() tea.Cmd {
	return m.routeInitial()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.lastSize = msg
	case RedirectMsg:
		if m.Middleware != nil && msg.Type != Pop {
			newTarget := m.Middleware(msg.Target)
			if newTarget != "" {
				msg.Type = Go
				msg.Target = newTarget
			}
		}

		switch msg.Type {
		case Go:
			m, cmd = m.gox(msg.Target)
		case Push:
			m, cmd = m.push(msg.Target)
		case Replace:
			m, cmd = m.replace(msg.Target)
		case Pop:
			m, cmd = m.pop()
		}
		cmds = append(cmds, cmd)
	}

	if length := len(m.modelStack); length > 0 {
		var cmdx tea.Cmd
		m.modelStack[length-1], cmdx = m.modelStack[length-1].Update(msg)
		cmds = append(cmds, cmdx)
	} else {
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		cmds = append(cmds, newErrorCmd(errors.New("router stack is empty, no model to update")))
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if length := len(m.modelStack); length > 0 {
		return m.modelStack[length-1].View()
	}
	return "TEA ROUTER STACK CAN'T BE EMPTY, YOU SHOULD GO REDIRECT ANYWARE"
}

// findRoutePath finds all routes in the hierarchy that lead to the target path.
func (m Model) findRoutePath(routes []Route, target string, parentPath string) []Route {
	for _, route := range routes {
		var currentPath string
		if strings.HasPrefix(route.Path, "/") {
			currentPath = path.Clean(route.Path)
		} else {
			currentPath = path.Join(parentPath, route.Path)
		}

		if currentPath == target {
			return []Route{route}
		}

		// Ensure target starts with currentPath followed by a / to avoid partial matches
		// e.g., /main should not match /main2
		prefix := currentPath
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}

		if strings.HasPrefix(target, prefix) {
			if subPath := m.findRoutePath(route.Children, target, currentPath); subPath != nil {
				return append([]Route{route}, subPath...)
			}
		}
	}
	return nil
}

func (m Model) gox(target string) (Model, tea.Cmd) {
	routePath := m.findRoutePath(m.Routes, target, "")
	if len(routePath) > 0 {
		m.modelStack = make([]tea.Model, len(routePath))
		for i, route := range routePath {
			m.modelStack[i] = route.Builder()
		}
		return m, m.initModel(len(m.modelStack) - 1)
	}
	return m, newErrorCmd(fmt.Errorf("route not found: %s", target))
}

func (m Model) push(target string) (Model, tea.Cmd) {
	routePath := m.findRoutePath(m.Routes, target, "")
	if len(routePath) > 0 {
		// If we are pushing a hierarchical route, we should push all missing parts
		// For now, let's simplify: if it's already in the hierarchy, we just push the final target
		// But the user wants Pop to work correctly.
		// A better approach: if we push /a/b/c, the stack should ideally contain a, b, and c if they are missing.
		
		// Let's implement full hierarchical push: replace stack with hierarchical path
		// BUT 'Push' in routers usually means just adding one level.
		// Given the user's feedback, they want the hierarchy to be there.
		
		// Update stack to match the target's full hierarchy
		newStack := make([]tea.Model, len(routePath))
		for i, route := range routePath {
			newStack[i] = route.Builder()
		}
		m.modelStack = newStack
		return m, m.initModel(len(m.modelStack) - 1)
	}
	return m, newErrorCmd(fmt.Errorf("route not found: %s", target))
}

func (m Model) replace(target string) (Model, tea.Cmd) {
	if len(m.modelStack) == 0 {
		return m, newErrorCmd(errors.New("cannot replace on an empty stack"))
	}
	routePath := m.findRoutePath(m.Routes, target, "")
	if len(routePath) > 0 {
		// Replace the last one with the new hierarchy
		newStack := make([]tea.Model, len(routePath))
		for i, route := range routePath {
			newStack[i] = route.Builder()
		}
		m.modelStack = newStack
		return m, m.initModel(len(m.modelStack) - 1)
	}
	return m, newErrorCmd(fmt.Errorf("route not found: %s", target))
}

func (m Model) pop() (Model, tea.Cmd) {
	if length := len(m.modelStack); length > 1 {
		m.modelStack = m.modelStack[:length-1]
		return m, m.initModel(len(m.modelStack) - 1)
	}
	return m, newErrorCmd(errors.New("cannot pop from the root of the stack"))
}

func (m Model) initModel(index int) tea.Cmd {
	cmds := []tea.Cmd{m.modelStack[index].Init()}
	if m.lastSize.Width > 0 || m.lastSize.Height > 0 {
		var cmd tea.Cmd
		m.modelStack[index], cmd = m.modelStack[index].Update(m.lastSize)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (m Model) routeInitial() tea.Cmd {
	if m.InitialRoute == "" {
		m.InitialRoute = "/"
	}
	return Redirect(Go, m.InitialRoute)
}
