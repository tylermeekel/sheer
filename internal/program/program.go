package program

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tylermeekel/sheer/internal/steps/main"
)

func RunCLI() {
	p := tea.NewProgram(mainmodel.New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
