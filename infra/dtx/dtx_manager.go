package dtx

import (
	"context"

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

func (d *DistributeTxManager) NewSagaTx(ctx context.Context, steps []*SagaStep) *Saga {
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
				fn:      d.handlers[s.Action.Name],
				payload: s.Action.Payload,
			},
			compensate: Caller{
				fn:      d.handlers[s.Compensate.Name],
				payload: s.Compensate.Payload,
			},
			actionCh:     make(chan struct{}),
			compensateCh: make(chan struct{}),
			closed:       make(chan struct{}),
		}
		stepMap[s.Name] = stp
	}
	for i := 0; i < len(steps); i++ {
		if steps[i].Action.isNoDepend() {
			stepMap[steps[i].Name].previous = append(stepMap[steps[i].Name].previous, trans.head)
		}
		for j := 0; j != i && j < len(steps); j++ {
			if steps[j].Action.isDependOn(steps[i].Action) {
				stepMap[steps[j].Name].previous = append(stepMap[steps[j].Name].previous, stepMap[steps[i].Name])
				stepMap[steps[i].Name].next = append(stepMap[steps[i].Name].next, stepMap[steps[j].Name])
			}
			if steps[j].Compensate.isDependOn(steps[i].Compensate) {
				stepMap[steps[j].Name].compensatePrevious = append(stepMap[steps[j].Name].compensatePrevious, stepMap[steps[i].Name])
				stepMap[steps[i].Name].compensateNext = append(stepMap[steps[i].Name].compensateNext, stepMap[steps[j].Name])
			}
		}
	}

	return trans
}

func (d *DistributeTxManager) RegisterHandler(ctx context.Context) error {
	return nil
}

func (h *TransHandler) isDependOn(t TransHandler) bool {
	for _, d := range h.depend {
		if d == t.Name {
			return true
		}
	}
	return false
}

func (h *TransHandler) isNoDepend() bool {
	return len(h.depend) == 0
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
	Action     TransHandler
	Compensate TransHandler
}

func isActionCircleDepend(ctx context.Context, s string, exist map[string]struct{}, m map[string]*SagaStep) bool {
	start := m[s]
	for _, d := range start.Action.depend {
		if _, ok := exist[d]; ok {
			logx.Error(ctx, "存在循环依赖")
			return true
		}
		exist[d] = struct{}{}
		return isActionCircleDepend(ctx, d, exist, m)
	}
	return false
}

func (s *Step) isCircleDepend() bool {
	exist := make(map[string]struct{})

	var isCircle func(p *Step) bool
	isCircle = func(p *Step) bool {
		if _, ok := exist[p.name]; ok {
			return true
		}
		exist[p.name] = struct{}{}
		for _, stp := range p.next {
			if isCircle(stp) {
				return true
			}
		}
		delete(exist, p.name)
		return false
	}

	return isCircle(s)
}
