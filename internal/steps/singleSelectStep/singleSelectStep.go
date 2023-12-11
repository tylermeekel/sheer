package singleselectstep

import (
	tea "github.com/charmbracelet/bubbletea"
)

type singleSelectStep struct {
	title    string
	options  []string
	selected int
}

func New(title string, opts []string) *singleSelectStep {
	s := singleSelectStep{
		title:    title,
		options:  opts,
		selected: 0,
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

func (s *singleSelectStep) Update(msg tea.Msg) string {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			s.Next()
		case tea.KeyUp:
			s.Prev()
		case tea.KeyEnter:
			return s.options[s.selected]
		}
	}

	return ""
}

func (s *singleSelectStep) Render() string {
	str := s.title

	for i, option := range s.options {
		str += "\n"
		if s.selected == i {
			str += "[x] " + option
		} else {
			str += "[ ] " + option
		}
	}

	return str
}
