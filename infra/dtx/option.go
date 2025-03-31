package dtx

type Option func(*TransSaga)

func WithStrictCompensate() Option {
	return func(t *TransSaga) {
		t.strictCompensate = true
	}
}

type StepOption func(*Step)

func WithStepID(id string) StepOption {
	return func(s *Step) {
		s.id = id
	}
}

func WithAction(a Caller) StepOption {
	return func(s *Step) {
		s.action = a
	}
}

func WithCompensate(c Caller) StepOption {
	return func(s *Step) {
		s.compensate = c
	}
}

func WithState(state StepStatus) StepOption {
	return func(s *Step) {
		s.state = state
	}
}
