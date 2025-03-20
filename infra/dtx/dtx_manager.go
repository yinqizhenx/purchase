package dtx

import (
	"context"
	"errors"
	"fmt"
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
		t.AsyncExec()
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

func (txm *DistributeTxManager) NewTransSaga(steps []*TransSagaStep, opts ...Option) (*TransSaga, error) {
	trans := &TransSaga{
		id:      uuid.NewString(),
		storage: txm.storage,
		errCh:   make(chan error, 1),
		done:    make(chan struct{}),
	}
	branches := make([]*Branch, 0, len(steps))
	for _, stp := range steps {
		branches = append(branches, &Branch{
			BranchID:     uuid.NewString(),
			State:        StepPending,
			Name:         stp.Name,
			Action:       stp.Action,
			Compensate:   stp.Compensate,
			Payload:      string(stp.Payload),
			ActionDepend: stp.Depend,
		})
	}
	return trans.build(branches, txm.handlers, opts...)
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
		saga, err := txm.buildTransSaga(t, branches)
		if err != nil {
			return nil, err
		}
		result = append(result, saga)
	}
	return result, nil
}

func (txm *DistributeTxManager) buildTransSaga(t *Trans, branches []*Branch, opts ...Option) (*TransSaga, error) {
	trans := &TransSaga{
		id:       t.TransID,
		storage:  txm.storage,
		errCh:    make(chan error, 1),
		done:     make(chan struct{}),
		isFromDB: true,
	}
	return trans.build(branches, txm.handlers, opts...)
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
		// Depend:     []string{"step5"},
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
		// time.Sleep(16 * time.Second)
		fmt.Println("run step3_rollback")
		return nil
	})
	txm.RegisterHandler("step4_action", func(ctx context.Context, payload []byte) error {
		time.Sleep(5 * time.Second)
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
