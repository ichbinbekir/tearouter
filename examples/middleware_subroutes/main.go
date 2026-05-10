package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ichbinbekir/tearouter"
)

// --- State and Middleware ---

var isLoggedIn = false

func authMiddleware(targetPath string) string {
	// Protect all routes starting with /admin
	if strings.HasPrefix(targetPath, "/admin") && !isLoggedIn {
		fmt.Printf("\n[Middleware] Unauthorized access to %s, redirecting to /login\n", targetPath)
		return "/login"
	}
	return ""
}

// --- Page Models ---

type Page struct {
	title string
	body  string
}

func (p Page) Init() tea.Cmd { return nil }
func (p Page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "q", "ctrl+c":
			return p, tea.Quit
		case "1":
			return p, tearouter.Redirect(tearouter.Go, "/public")
		case "2":
			return p, tearouter.Redirect(tearouter.Go, "/admin/dashboard")
		case "3":
			return p, tearouter.Redirect(tearouter.Go, "/admin/settings")
		case "l":
			isLoggedIn = !isLoggedIn
			return p, nil
		case "b":
			return p, tearouter.Redirect(tearouter.Pop)
		}
	}
	return p, nil
}

func (p Page) View() string {
	status := "Logged OUT"
	if isLoggedIn {
		status = "Logged IN"
	}

	return fmt.Sprintf(
		"Current Page: %s\nStatus: %s\n\n"+
			"Press 1: Go to /public\n"+
			"Press 2: Go to /admin/dashboard (Protected)\n"+
			"Press 3: Go to /admin/settings (Protected)\n"+
			"Press l: Toggle Login Status\n"+
			"Press b: Back (Pop)\n"+
			"Press q: Quit\n\n"+
			"%s",
		p.title, status, p.body,
	)
}

func main() {
	model := tearouter.Model{
		InitialRoute: "/public",
		Middleware:   authMiddleware,
		Routes: []tearouter.Route{
			{
				Path: "/public",
				Builder: func() tea.Model {
					return Page{title: "Public Area", body: "Welcome! This area is accessible to everyone."}
				},
			},
			{
				Path: "/login",
				Builder: func() tea.Model {
					return Page{title: "Login Page", body: "Please 'Log In' by pressing 'l' to access admin area."}
				},
			},
			{
				Path: "/admin",
				Builder: func() tea.Model {
					return Page{title: "Admin Shell", body: "This is the parent Admin layout."}
				},
				Children: []tearouter.Route{
					{
						Path: "dashboard",
						Builder: func() tea.Model {
							return Page{title: "Admin Dashboard", body: "Sensitive Data: 42"}
						},
					},
					{
						Path: "settings",
						Builder: func() tea.Model {
							return Page{title: "Admin Settings", body: "Change system behavior here."}
						},
					},
				},
			},
		},
	}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		panic(err)
	}
}
