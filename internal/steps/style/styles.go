package style

import "github.com/charmbracelet/lipgloss"

var BaseScreenStyle = lipgloss.NewStyle().
	Padding(1, 4).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(AccentColor).
	Width(80)

var BaseTitleStyle = lipgloss.NewStyle().
	Bold(true).
	MarginBottom(1).
	Underline(true)
