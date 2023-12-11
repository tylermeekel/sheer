package program

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	mainmodel "github.com/tylermeekel/sheer/internal/steps/mainModel"
)

func RunCLI() {
	p := tea.NewProgram(mainmodel.New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
