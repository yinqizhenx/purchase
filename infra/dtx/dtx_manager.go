package dtx

import (
	"context"
	"errors"
	"fmt"
)

func NewDistributeTxManager() *DistributeTxManager {
	return &DistributeTxManager{}
}

type DistributeTxManager struct {
	steps    []Step
	storage  TransStorage
	handlers map[string]func(context.Context, []byte) error
}

func (t *DistributeTxManager) Start(ctx context.Context) error {
	return nil
}

func (t *DistributeTxManager) NewTx(ctx context.Context) *Saga {
	return nil
}

func (t *DistributeTxManager) NewSagaTx(ctx context.Context, steps []*SagaStep) (*Saga, error) {
	trans := &Saga{}
	head := &Step{
		saga:         trans,
		actionCh:     make(chan struct{}),
		compensateCh: make(chan struct{}),
	}
	trans.head = head
	stepMap := make(map[string]*Step)
	for _, s := range steps {
		stp := &Step{
			saga: trans,
			name: s.Name,
			action: Caller{
				fn:      t.handlers[s.Action],
				payload: s.Payload,
			},
			compensate: Caller{
				fn:      t.handlers[s.Compensate],
				payload: s.Payload,
			},
			actionCh:     make(chan struct{}),
			compensateCh: make(chan struct{}),
			closed:       make(chan struct{}),
		}
		stepMap[s.Name] = stp
	}

	for _, stp := range steps {
		for _, dp := range stp.Depend {
			if _, ok := stepMap[dp]; !ok {
				return nil, errors.New(fmt.Sprintf("depend not exist: %s", dp))
			}
		}
	}

	for _, step := range steps {
		if step.hasNoDepend() {
			trans.head.next = append(trans.head.next, stepMap[step.Name])
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, trans.head)
		}
		for _, d := range step.Depend {
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, stepMap[d])
			stepMap[d].next = append(stepMap[d].next, stepMap[step.Name])
			stepMap[step.Name].compensateNext = append(stepMap[step.Name].compensateNext, stepMap[d])
			stepMap[d].compensatePrevious = append(stepMap[d].compensatePrevious, stepMap[step.Name])
		}
	}

	if trans.head.isCircleDepend() {
		return nil, errors.New("exist circle depend")
	}
	return trans, nil
}

func (t *DistributeTxManager) RegisterHandler(ctx context.Context) error {
	return nil
}

func (h *SagaStep) hasNoDepend() bool {
	return len(h.Depend) == 0
}

type SagaStep struct {
	Name       string
	Action     string
	Compensate string
	Payload    []byte
	Depend     []string
}
