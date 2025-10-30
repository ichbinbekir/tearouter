package tearouter

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RedirectType int

const (
	Go RedirectType = iota
	Push
	Replace
	Pop
)

type RedirectMsg struct {
	Type   RedirectType
	Target string
}

func Redirect(typ RedirectType, target ...string) tea.Cmd {
	if typ != Pop && len(target) < 1 {
		// TODO err
	}
	if typ == Pop {
		return func() tea.Msg {
			return RedirectMsg{Type: typ}
		}
	}
	return func() tea.Msg {
		return RedirectMsg{
			Type:   typ,
			Target: target[0],
		}
	}
}
