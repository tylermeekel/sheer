package teststep

type TestStep struct {
	stepdata string
}

func (t *TestStep) Render() string {
	return t.stepdata
}