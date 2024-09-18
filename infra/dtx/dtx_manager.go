package dtx

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
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
	trans    []*TransSaga
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

func (txm *DistributeTxManager) NewTx(ctx context.Context) *TransSaga {
	return nil
}

func (txm *DistributeTxManager) NewTransSagaTx(ctx context.Context, steps []*TransSagaStep, opts ...Option) (*TransSaga, error) {
	trans := &TransSaga{
		id:      uuid.NewString(),
		storage: txm.storage,
		errCh:   make(chan error, 1),
		done:    make(chan struct{}),
	}
	root := &Step{
		id:    uuid.NewString(),
		saga:  trans,
		name:  rootStepName,
		state: StepPending,
		action: Caller{
			fn: func(context.Context, []byte) error { return nil },
		},
		compensate: Caller{
			fn: func(context.Context, []byte) error {
				trans.close()
				return nil
			},
		},
		actionCh:     make(chan struct{}),
		compensateCh: make(chan struct{}),
		closed:       make(chan struct{}),
	}
	trans.root = root
	stepMap := make(map[string]*Step)
	for _, s := range steps {
		stp := &Step{
			id:    uuid.NewString(),
			saga:  trans,
			name:  s.Name,
			state: StepPending,
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

func (txm *DistributeTxManager) loadPendingTrans(ctx context.Context) ([]*TransSaga, error) {
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

	result := make([]*TransSaga, 0, len(transMap))
	for id, t := range transMap {
		branches := branchesMap[id]
		saga := txm.buildTransSaga(t, branches)
		result = append(result, saga)
	}
	return result, nil
}

func (txm *DistributeTxManager) buildTransSaga(t *Trans, branches []*Branch, opts ...Option) *TransSaga {
	trans := &TransSaga{
		id:       t.TransID,
		storage:  txm.storage,
		errCh:    make(chan error, 1),
		done:     make(chan struct{}),
		isFromDB: true,
	}
	root := &Step{
		id:    uuid.NewString(),
		saga:  trans,
		name:  rootStepName,
		state: StepPending,
		action: Caller{
			fn: func(context.Context, []byte) error { return nil },
		},
		compensate: Caller{
			fn: func(context.Context, []byte) error {
				trans.close()
				return nil
			},
		},
		actionCh:     make(chan struct{}),
		compensateCh: make(chan struct{}),
		closed:       make(chan struct{}),
	}
	trans.root = root
	stepMap := make(map[string]*Step)
	for _, s := range branches {
		stp := &Step{
			id:    s.BranchID,
			saga:  trans,
			name:  s.Name,
			state: StepStatus(s.State),
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

			actionCh:     make(chan struct{}),
			compensateCh: make(chan struct{}),
			closed:       make(chan struct{}),
		}
		stepMap[s.Name] = stp
	}

	for _, step := range branches {
		depends := strings.Split(step.ActionDepend, ",")
		if depends[0] == "" {
			trans.root.next = append(trans.root.next, stepMap[step.Name])
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, trans.root)
			trans.root.compensatePrevious = append(trans.root.compensatePrevious, stepMap[step.Name])
			stepMap[step.Name].compensateNext = append(stepMap[step.Name].compensateNext, trans.root)
			continue
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

	for _, opt := range opts {
		opt(trans)
	}

	if trans.timeout == 0 {
		trans.timeout = defaultTimeout
	}

	return trans
}

func (txm *DistributeTxManager) RegisterHandler(name string, handler func(context.Context, []byte) error) {
	if _, ok := txm.handlers[name]; ok {
		panic(fmt.Sprintf("duplicated handler: %s", name))
	}
	txm.handlers[name] = handler
}

func (h *TransSagaStep) hasNoDepend() bool {
	return len(h.Depend) == 0
}

type TransSagaStep struct {
	Name       string
	Action     string
	Compensate string
	Payload    []byte
	Depend     []string
}

var StepTest = []*TransSagaStep{
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
		Depend:     []string{"step1"},
		Payload:    nil,
	},
	{
		Name:       "step3",
		Action:     "step3_action",
		Compensate: "step3_rollback",
		Depend:     []string{"step1"},
		Payload:    nil,
	},
	{
		Name:       "step4",
		Action:     "step4_action",
		Compensate: "step4_rollback",
		Depend:     []string{"step2"},
		Payload:    nil,
	},
	{
		Name:       "step5",
		Action:     "step5_action",
		Compensate: "step5_rollback",
		//Depend:     []string{"step5"},
		Payload: nil,
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
		//time.Sleep(16 * time.Second)
		fmt.Println("run step3_rollback")
		return nil
	})
	txm.RegisterHandler("step4_action", func(ctx context.Context, payload []byte) error {
		//time.Sleep(16 * time.Second)
		fmt.Println("run step4_action")
		return errors.New("ssccc")
	})
	txm.RegisterHandler("step4_rollback", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step4_rollback")
		return nil
	})
	txm.RegisterHandler("step5_action", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step5_action")
		return nil
	})
	txm.RegisterHandler("step5_rollback", func(ctx context.Context, payload []byte) error {
		fmt.Println("run step5_rollback")
		return nil
	})
}
