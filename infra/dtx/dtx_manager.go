package dtx

import (
	"context"
	"errors"
	"fmt"

	"purchase/infra/logx"
)

func NewDistributeTxManager() *DistributeTxManager {
	return &DistributeTxManager{}
}

type DistributeTxManager struct {
	steps    []Step
	storage  TransStorage
	handlers map[string]func(context.Context, []byte) error
}

func (d *DistributeTxManager) Start(ctx context.Context) error {
	return nil
}

func (d *DistributeTxManager) NewTx(ctx context.Context) *Saga {
	return nil
}

func (d *DistributeTxManager) NewSagaTx(ctx context.Context, steps []*SagaStep) (*Saga, error) {
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
				fn:      d.handlers[s.Action],
				payload: s.Payload,
			},
			compensate: Caller{
				fn:      d.handlers[s.Compensate],
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

	for i := 0; i < len(steps); i++ {
		if steps[i].isNoDepend() {
			trans.head.next = append(trans.head.next, stepMap[steps[i].Name])
			stepMap[steps[i].Name].previous = append(stepMap[steps[i].Name].previous, trans.head)
		}
		for j := 0; j != i && j < len(steps); j++ {
			if steps[j].isDependOn(steps[i]) {
				stepMap[steps[j].Name].previous = append(stepMap[steps[j].Name].previous, stepMap[steps[i].Name])
				stepMap[steps[i].Name].next = append(stepMap[steps[i].Name].next, stepMap[steps[j].Name])
			}
			if steps[j].isDependOn(steps[i]) {
				stepMap[steps[j].Name].compensatePrevious = append(stepMap[steps[j].Name].compensatePrevious, stepMap[steps[i].Name])
				stepMap[steps[i].Name].compensateNext = append(stepMap[steps[i].Name].compensateNext, stepMap[steps[j].Name])
			}
		}
	}

	if trans.head.isCircleDepend() {
		return nil, errors.New("exist circle depend")
	}
	return trans, nil
}

func (d *DistributeTxManager) RegisterHandler(ctx context.Context) error {
	return nil
}

func (h *SagaStep) isDependOn(t *SagaStep) bool {
	for _, d := range h.Depend {
		if d == t.Name {
			return true
		}
	}
	return false
}

func (h *SagaStep) isNoDepend() bool {
	return len(h.Depend) == 0
}

type TransHandler struct {
	Name    string
	Payload []byte
	depend  []string
}

func NewTransHandler(name string, payload []byte, depends ...string) TransHandler {
	return TransHandler{
		Name:    name,
		Payload: payload,
		depend:  depends,
	}
}

type SagaStep struct {
	Name       string
	Action     string
	Compensate string
	Payload    []byte
	Depend     []string
}

func isActionCircleDepend(ctx context.Context, s string, exist map[string]struct{}, m map[string]*SagaStep) bool {
	start := m[s]
	for _, d := range start.Depend {
		if _, ok := exist[d]; ok {
			logx.Error(ctx, "存在循环依赖")
			return true
		}
		exist[d] = struct{}{}
		return isActionCircleDepend(ctx, d, exist, m)
	}
	return false
}
