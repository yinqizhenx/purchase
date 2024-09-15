package dtx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	rootStepName = "root"
)

func NewDistributeTxManager(s TransStorage) *DistributeTxManager {
	txm := &DistributeTxManager{
		storage:  s,
		handlers: make(map[string]func(context.Context, []byte) error),
	}
	registerForTest(txm)
	return txm
}

type DistributeTxManager struct {
	trans    []*Saga
	storage  TransStorage
	handlers map[string]func(context.Context, []byte) error
}

func (txm *DistributeTxManager) Start(ctx context.Context) error {
	trans, err := txm.loadPendingTrans(ctx)
	if err != nil {
		return err
	}
	for _, t := range trans {
		t.AsyncExec(ctx)
	}
	return nil
}

func (txm *DistributeTxManager) Stop(ctx context.Context) error {
	for _, t := range txm.trans {
		t.close()
	}
	return nil
}

func (txm *DistributeTxManager) NewTx(ctx context.Context) *Saga {
	return nil
}

func (txm *DistributeTxManager) NewSagaTx(ctx context.Context, steps []*SagaStep, opts ...Option) (*Saga, error) {
	trans := &Saga{
		id:      uuid.NewString(),
		storage: txm.storage,
		errCh:   make(chan error),
		done:    make(chan struct{}),
	}
	root := &Step{
		id:   uuid.NewString(),
		saga: trans,
		name: rootStepName,
		action: Caller{
			fn: func(context.Context, []byte) error { return nil },
		},
		compensate: Caller{
			fn: func(context.Context, []byte) error {
				trans.close()
				return nil
			},
		},
		actionCh:       make(chan struct{}),
		compensateCh:   make(chan struct{}),
		actionDone:     make(chan error),
		compensateDone: make(chan error),
	}
	trans.root = root
	stepMap := make(map[string]*Step)
	for _, s := range steps {
		stp := &Step{
			id:   uuid.NewString(),
			saga: trans,
			name: s.Name,
			action: Caller{
				name:    s.Action,
				fn:      txm.handlers[s.Action],
				payload: s.Payload,
			},
			compensate: Caller{
				name:    s.Compensate,
				fn:      txm.handlers[s.Compensate],
				payload: s.Payload,
			},
			actionCh:       make(chan struct{}),
			compensateCh:   make(chan struct{}),
			closed:         make(chan struct{}),
			actionDone:     make(chan error),
			compensateDone: make(chan error),
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
			trans.root.next = append(trans.root.next, stepMap[step.Name])
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, trans.root)
			trans.root.compensatePrevious = append(trans.root.compensatePrevious, stepMap[step.Name])
			stepMap[step.Name].compensateNext = append(stepMap[step.Name].compensateNext, trans.root)
		}
		for _, d := range step.Depend {
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, stepMap[d])
			stepMap[d].next = append(stepMap[d].next, stepMap[step.Name])
			stepMap[step.Name].compensateNext = append(stepMap[step.Name].compensateNext, stepMap[d])
			stepMap[d].compensatePrevious = append(stepMap[d].compensatePrevious, stepMap[step.Name])
		}
	}

	if len(trans.root.next) == 0 {
		return nil, errors.New("exist circle depend, root next can not be empty")
	}

	for _, stp := range stepMap {
		if stp.isCircleDepend() {
			return nil, errors.New("exist circle depend")
		}
	}

	for _, stp := range stepMap {
		trans.steps = append(trans.steps, stp)
	}

	for _, opt := range opts {
		opt(trans)
	}

	if trans.timeout == 0 {
		trans.timeout = defaultTimeout
	}
	return trans, nil
}

func (txm *DistributeTxManager) loadPendingTrans(ctx context.Context) ([]*Saga, error) {
	transMap, err := txm.storage.GetPendingTrans(ctx)
	if err != nil {
		return nil, err
	}

	transIDList := make([]string, 0, len(transMap))
	for id := range transMap {
		transIDList = append(transIDList, id)
	}

	branchesMap, err := txm.storage.MustGetBranchesByTransIDList(ctx, transIDList)
	if err != nil {
		return nil, err
	}

	result := make([]*Saga, 0, len(transMap))
	for id, t := range transMap {
		branches := branchesMap[id]
		saga := txm.buildTransSaga(t, branches)
		result = append(result, saga)
	}
	return result, nil
}

func (txm *DistributeTxManager) buildTransSaga(t *Trans, branches []*Branch) *Saga {
	trans := &Saga{
		id:      t.TransID,
		storage: txm.storage,
		errCh:   make(chan error),
		done:    make(chan struct{}),
	}
	root := &Step{
		id:   uuid.NewString(),
		saga: trans,
		name: rootStepName,
		action: Caller{
			fn: func(context.Context, []byte) error { return nil },
		},
		compensate: Caller{
			fn: func(context.Context, []byte) error {
				trans.close()
				return nil
			},
		},
		actionCh:       make(chan struct{}),
		compensateCh:   make(chan struct{}),
		actionDone:     make(chan error),
		compensateDone: make(chan error),
	}
	trans.root = root
	stepMap := make(map[string]*Step)
	for _, s := range branches {
		stp := &Step{
			id:   s.BranchID,
			saga: trans,
			name: s.Name,
			action: Caller{
				name:    s.Action,
				fn:      txm.handlers[s.Action],
				payload: []byte(s.Payload),
			},
			compensate: Caller{
				name:    s.Compensate,
				fn:      txm.handlers[s.Compensate],
				payload: []byte(s.Payload),
			},

			actionCh:       make(chan struct{}),
			compensateCh:   make(chan struct{}),
			closed:         make(chan struct{}),
			actionDone:     make(chan error),
			compensateDone: make(chan error),
		}
		stepMap[s.Name] = stp
	}

	for _, step := range branches {
		depends := strings.Split(step.ActionDepend, ",")
		if len(depends) == 0 {
			trans.root.next = append(trans.root.next, stepMap[step.Name])
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, trans.root)
			trans.root.compensatePrevious = append(trans.root.compensatePrevious, stepMap[step.Name])
			stepMap[step.Name].compensateNext = append(stepMap[step.Name].compensateNext, trans.root)
		}
		for _, d := range depends {
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, stepMap[d])
			stepMap[d].next = append(stepMap[d].next, stepMap[step.Name])
			stepMap[step.Name].compensateNext = append(stepMap[step.Name].compensateNext, stepMap[d])
			stepMap[d].compensatePrevious = append(stepMap[d].compensatePrevious, stepMap[step.Name])
		}
	}

	for _, stp := range stepMap {
		trans.steps = append(trans.steps, stp)
	}

	return trans
}

func (txm *DistributeTxManager) RegisterHandler(name string, handler func(context.Context, []byte) error) {
	if _, ok := txm.handlers[name]; ok {
		panic(fmt.Sprintf("duplicated handler: %s", name))
	}
	txm.handlers[name] = handler
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

var StepTest = []*SagaStep{
	{
		Name:       "step1",
		Action:     "step1_action",
		Compensate: "step1_rollback",
		Depend:     nil,
		Payload:    nil,
	},
	{
		Name:       "step2",
		Action:     "step2_action",
		Compensate: "step2_rollback",
		Depend:     nil,
		Payload:    nil,
	},
	{
		Name:       "step3",
		Action:     "step3_action",
		Compensate: "step3_rollback",
		Depend:     nil,
		Payload:    nil,
	},
	{
		Name:       "step4",
		Action:     "step4_action",
		Compensate: "step4_rollback",
		Depend:     []string{"step1", "step2", "step3"},
		Payload:    nil,
	},
}

func registerForTest(txm *DistributeTxManager) {
	txm.RegisterHandler("step1_action", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step1_action")
		return nil
	})
	txm.RegisterHandler("step1_rollback", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step1_rollback")
		return nil
	})
	txm.RegisterHandler("step2_action", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step2_action")
		return nil
	})
	txm.RegisterHandler("step2_rollback", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step2_rollback")
		return nil
	})
	txm.RegisterHandler("step3_action", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step3_action")
		return nil
	})
	txm.RegisterHandler("step3_rollback", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step3_rollback")
		return nil
	})
	txm.RegisterHandler("step4_action", func(ctx context.Context, payload []byte) error {
		time.Sleep(1 * time.Second)
		fmt.Println("run step4_action")
		return errors.New("ssccc")
	})
	txm.RegisterHandler("step4_rollback", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step4_rollback")
		return nil
	})
}
