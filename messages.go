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
	Middleware
)

type RedirectMsg struct {
	Type   RedirectType
	Target string
}

func Redirect(typ RedirectType, target ...string) tea.Cmd {
	return func() tea.Msg {
		if len(target) < 1 {
			return ErrMsg{Err: errors.New("redirect, must define a target")}
		}
		return RedirectMsg{
			Type:   typ,
			Target: target[0],
		}
	}
}

type ErrMsg struct {
	Err error
}

func errCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrMsg{Err: err}
	}
}
