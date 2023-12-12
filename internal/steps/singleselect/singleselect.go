package singleselect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylermeekel/sheer/internal/steps/style"
)

type styles struct {
	screenStyle   lipgloss.Style
	selectedStyle lipgloss.Style
	titleStyle    lipgloss.Style
}

type singleSelectStep struct {
	title    string
	options  []string
	selected int
	result   string

	styles styles
}

func New(title string, opts []string) *singleSelectStep {
	screenStyle := lipgloss.NewStyle().
		Padding(1, 4).Border(lipgloss.RoundedBorder()).BorderForeground(style.AccentColor)

	selectedStyle := lipgloss.NewStyle().
		Foreground(style.AccentColor).
		Bold(true)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		MarginBottom(1).
		Underline(true)

	style := styles{
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

func (s *singleSelectStep) Next() {
	if s.selected < len(s.options)-1 {
		s.selected++
	} else {
		s.selected = 0
	}
}

func (s *singleSelectStep) Prev() {
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
			s.Next()
		case tea.KeyUp:
			s.Prev()
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
