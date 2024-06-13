package vo

import "github.com/pkg/errors"

var (
	SomeStatusZero = SomeStatus{}
	SomeStatusOne  = SomeStatus{1}
	SomeStatusTwo  = SomeStatus{2}
)

var someStatusValues = []SomeStatus{
	SomeStatusOne,
	SomeStatusTwo,
}

func NewSomeStatus(s int) (SomeStatus, error) {
	for _, item := range someStatusValues {
		if item.s == s {
			return item, nil
		}
	}
	return SomeStatusZero, errors.Errorf("unknown '%d' status", s)
}

type SomeStatus struct {
	s int
}

func (s SomeStatus) IsZero() bool {
	return s.s == 0
}

func (s SomeStatus) Int() int {
	return s.s
}
