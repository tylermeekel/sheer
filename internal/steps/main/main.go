package mainmodel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tylermeekel/sheer/internal/steps/singleselect"
)

type Step interface {
	Render() string
	Update(msg tea.Msg) bool
	Result() string
}

type StepTag int

const (
	sendReceive StepTag = iota
	fileSelect
	send
	receive
)

type mainModel struct {
	steps       map[StepTag]Step
	currentStep StepTag
}

func (em *mainModel) registerStep(tag StepTag, step Step) {
	em.steps[tag] = step
}

func New() *mainModel {
	em := mainModel{
		steps: make(map[StepTag]Step),
	}

	em.registerStep(sendReceive, singleselect.New("Send or Receive a File?", []string{"Send", "Receive"}))
	em.registerStep(fileSelect, singleselect.New("Select File", []string{"This is actually a test.", "This is a test", "Remember when I said this was a test"}))
	em.registerStep(receive, singleselect.New("Receive a File", []string{"This is actually a test.", "This is a test", "Remember when I said this was a test"}))

	em.currentStep = sendReceive
	return &em
}

func (em *mainModel) NextStep() StepTag {
	switch em.currentStep {
	case sendReceive:
		if em.steps[em.currentStep].Result() == "Send" {
			return fileSelect
		} else if em.steps[em.currentStep].Result() == "Receive" {
			return receive
		}
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
			if isStepDone := em.steps[em.currentStep].Update(msg); isStepDone { // Unnecessary but more readable
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
