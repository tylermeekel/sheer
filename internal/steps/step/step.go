package step

import tea "github.com/charmbracelet/bubbletea"

type Step interface {
	Render() string
	Update(msg tea.Msg) bool
	Result() string
}

type StepTag int

const (
	MainMenu StepTag = iota
	FileSelect
	Send
	Receive
	Config
)
