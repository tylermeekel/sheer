package mainmodel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tylermeekel/sheer/internal/steps/fileselect"
	"github.com/tylermeekel/sheer/internal/steps/singleselect"
	"github.com/tylermeekel/sheer/internal/steps/step"
)

type mainModel struct {
	steps       map[step.StepTag]step.Step
	currentStep step.StepTag
}

func (em *mainModel) registerStep(tag step.StepTag, step step.Step) {
	em.steps[tag] = step
}

func New() *mainModel {
	em := mainModel{
		steps: make(map[step.StepTag]step.Step),
	}

	//Main Step
	em.registerStep(step.MainMenu, singleselect.New("Sheer", []string{"Send a File", "Receive a File", "Configuration"}))
	
	//Steps for Sending
	em.registerStep(step.FileSelect, fileselect.New("Select a File to Send"))
	//em.registerStep(send, )
	
	//Steps for Receiving
	em.registerStep(step.Receive, singleselect.New("Receive a File", []string{"This is actually a test.", "This is a test", "Remember when I said this was a test"}))
	
	//Steps for config
	em.registerStep(step.Config, singleselect.New("Configuration", []string{"Change option 1", "Change option 2"}))

	em.currentStep = step.MainMenu
	return &em
}

func (em *mainModel) NextStep() step.StepTag { //? Maybe each step should signal what the next step should be
	switch em.currentStep {
	case step.MainMenu:
		if em.steps[step.MainMenu].Result() == "Send a File" {
			return step.FileSelect
		} else if em.steps[step.MainMenu].Result() == "Receive a File" {
			return step.Receive
		} else if em.steps[step.MainMenu].Result() == "Configuration" {
			return step.Config
		}
	case step.FileSelect:
		return step.Send
	}

	return -1
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
			if isStepDone := em.steps[em.currentStep].Update(msg); isStepDone {
				next := em.NextStep()
				if next == -1 {
					return em, tea.Quit
				} else {
					em.currentStep = next
					return em, nil
				}
			}
		}
	}

	return em, nil
}

func (em *mainModel) View() string {
	return em.steps[em.currentStep].Render()
}
