package mainmodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	singleselectstep "github.com/tylermeekel/sheer/internal/steps/singleSelectStep"
)

type Step interface {
	Render() string
	Update(msg tea.Msg) string //TODO: Add a Data return value
}

type mainModel struct {
	steps       map[string]Step
	currentStep Step
}

func (em *mainModel) registerStep(tag string, step Step) {
	em.steps[tag] = step
}

func New() *mainModel {
	em := mainModel{
		steps: make(map[string]Step),
	}

	em.registerStep("singleselect", singleselectstep.New("Send or Receive a File?", []string{"Send", "Receive"}))

	em.currentStep = em.steps["singleselect"]
	return &em
}

func (em *mainModel) Init() tea.Cmd {
	return nil
}

func (em *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return em, tea.Quit
		default:
			fmt.Println(msg)
			em.currentStep.Update(msg)
		}
	}

	return em, nil
}

func (em *mainModel) View() string {
	return em.currentStep.Render()
}
