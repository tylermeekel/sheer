package singleselect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylermeekel/sheer/internal/steps/style"
)

type singleSelectStyles struct {
	screenStyle   lipgloss.Style
	selectedStyle lipgloss.Style
	titleStyle    lipgloss.Style
}

type singleSelectStep struct {
	title    string
	options  []string
	selected int
	result   string

	styles singleSelectStyles
}

func New(title string, opts []string) *singleSelectStep {
	screenStyle := style.BaseScreenStyle

	selectedStyle := lipgloss.NewStyle().
		Foreground(style.AccentColor).
		Bold(true)

	titleStyle := style.BaseTitleStyle

	style := singleSelectStyles{
		screenStyle:   screenStyle,
		selectedStyle: selectedStyle,
		titleStyle: titleStyle,
	}

	s := singleSelectStep{
		title:    title,
		options:  opts,
		selected: 0,
		styles:   style,
	}

	return &s
}

func (s *singleSelectStep) next() {
	if s.selected < len(s.options)-1 {
		s.selected++
	} else {
		s.selected = 0
	}
}

func (s *singleSelectStep) prev() {
	if s.selected > 0 {
		s.selected--
	} else {
		s.selected = len(s.options) - 1
	}
}

func (s *singleSelectStep) Update(msg tea.Msg) (done bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			s.next()
		case tea.KeyUp:
			s.prev()
		case tea.KeyEnter:
			s.result = s.options[s.selected]
			done = true
		}
	}

	return
}

func (s *singleSelectStep) Render() string {
	var str string

	str += s.styles.titleStyle.Render(s.title)

	for i, option := range s.options {
		str += "\n"
		if s.selected == i {
			str += s.styles.selectedStyle.Render("> " + option)
		} else {
			str += option
		}
	}

	str = s.styles.screenStyle.Render(str)
	return str
}

func (s *singleSelectStep) Result() string {
	return s.result
}
