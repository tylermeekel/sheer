package cli

import tea "github.com/charmbracelet/bubbletea"

type Step interface {
	Render() string
}

type mainModel struct {
	steps       map[string]Step
	currentStep Step
}

func New(steps map[string]Step) mainModel {
	return mainModel{
		steps: steps,
	}
}

func (m *mainModel) Init() tea.Cmd {
	return nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	return m, nil
}

func (m *mainModel) View() string {
	return m.currentStep.Render()
}

func RunCLI() {
	
}
