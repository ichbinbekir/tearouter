package tearouter

import (
	"errors"

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

type ErrorMsg struct {
	Err error
}

func newErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{Err: err}
	}
}

func Redirect(typ RedirectType, target ...string) tea.Cmd {
	if typ != Pop && len(target) < 1 {
		return newErrorCmd(errors.New("redirect target can't be empty"))
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
