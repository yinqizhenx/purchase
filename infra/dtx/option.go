package dtx

import "time"

type Option func(*TransSaga)

func WithTimeout(d time.Duration) Option {
	return func(t *TransSaga) {
		t.timeout = d
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
