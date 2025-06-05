package dtx

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"purchase/infra/utils"

	"purchase/infra/logx"
)

const (
	defaultTimeout = 3 * time.Second
)

func NewTransSaga() *TransSaga {
	t := &TransSaga{}
	return t
}

type TransSaga struct {
	name             string
	id               int
	root             *Step
	steps            []*Step
	state            atomic.Value // executing - 执行中， success - 成功， fail - 失败
	errCh            chan error   // 正向执行的错误channel，容量为1
	storage          TransStorage
	done             chan struct{}
	isFromDB         bool
	strictCompensate bool // 回滚执行失败时，是否进行下游回滚，true代表不进行下游回滚，直接退出
	dtm              *DistributeTxManager
}

func (t *TransSaga) Exec(ctx context.Context) error {
	t.AsyncExec()
	return <-t.errCh
}

func (t *TransSaga) AsyncExec() {
	ctx := context.Background()
	utils.SafeGo(ctx, func() {
		t.state.Store(SagaStateExecuting)

		if !t.isFromDB {
			err := t.sync(ctx)
			if err != nil {
				t.errCh <- err
				return
			}
		}

		utils.SafeGo(ctx, func() {
			t.root.runAction(ctx)
		})
		utils.SafeGo(ctx, func() {
			t.root.runCompensate(ctx)
		})
		for _, step := range t.steps {
			stp := step // 重新赋值，防止step引用变化
			utils.SafeGo(ctx, func() {
				stp.runAction(ctx)
			})
			utils.SafeGo(ctx, func() {
				stp.runCompensate(ctx)
			})
		}
		t.root.actionCh <- ""

		select {
		case <-t.done:
			logx.Info(ctx, "执行完成 退出AsyncExec")
		}
	})
}

func (t *TransSaga) sync(ctx context.Context) error {
	tx := &Trans{
		Name:       t.name,
		State:      string(SagaStateExecuting),
		FinishedAt: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:  time.Now(),
		CreatedBy:  utils.GetCurrentUser(ctx),
		UpdatedBy:  utils.GetCurrentUser(ctx),
		UpdatedAt:  time.Now(),
	}

	tranID, err := t.storage.SaveTrans(ctx, tx)
	if err != nil {
		return err
	}

	t.id = tranID

	branchList := make([]*Branch, 0)
	for _, stp := range t.steps {
		actionDepend := make([]string, 0)
		compensateDepend := make([]string, 0)
		for _, p := range stp.previous {
			name := p.name
			if p.name == rootStepName {
				name = ""
			}
			actionDepend = append(actionDepend, name)
		}
		for _, n := range stp.compensatePrevious {
			compensateDepend = append(compensateDepend, n.name)
		}

		stepID := fmt.Sprintf("%d_%s", tranID, stp.name)
		stp.id = stepID
		branch := &Branch{
			Code:              stepID,
			TransID:           tranID,
			State:             StepPending,
			Name:              stp.name,
			Action:            stp.action.name,
			Compensate:        stp.compensate.name,
			ActionPayload:     stp.action.payload,
			CompensatePayload: stp.compensate.payload,
			ActionDepend:      actionDepend,
			CompensateDepend:  compensateDepend,
			IsDead:            false,
			CreatedAt:         time.Now(),
			CreatedBy:         utils.GetCurrentUser(ctx),
			UpdatedBy:         utils.GetCurrentUser(ctx),
			UpdatedAt:         time.Now(),
		}
		branchList = append(branchList, branch)
	}

	err = t.storage.SaveBranch(ctx, branchList)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransSaga) tryUpdateSuccess(ctx context.Context) {
	for _, stp := range t.steps {
		if !stp.isSuccess() {
			return
		}
	}
	// 修改全局事务成功，防止并发同时close(s.errCh)
	if !t.state.CompareAndSwap(SagaStateExecuting, SagaStateSuccess) {
		return
	}
	if err := t.syncStateChange(ctx); err != nil {
		logx.Errorf(ctx, "update branch state fail: %v", err)
	}
	if err := t.changeExecuteStateFinished(ctx); err != nil {
		logx.Errorf(ctx, "update branch execute state fail, when success: %v", err)
	}

	t.close()
}

func (t *TransSaga) syncStateChange(ctx context.Context) error {
	return t.storage.UpdateTransState(ctx, t.id, string(t.state.Load().(SagaState)))
}

func (t *TransSaga) isFailed() bool {
	return t.state.Load() == SagaStateFail
}

func (t *TransSaga) isSuccess() bool {
	return t.state.Load() == SagaStateSuccess
}

func (t *TransSaga) build(steps []*Branch, handlers map[string]func(context.Context, []byte) error, opts ...Option) (*TransSaga, error) {
	t.root = t.buildRootStep()
	stepMap := make(map[string]*Step)
	for _, s := range steps {
		stp := &Step{
			id:    s.Code,
			saga:  t,
			name:  s.Name,
			state: s.State,
			action: Caller{
				name:    s.Action,
				fn:      handlers[s.Action],
				payload: s.ActionPayload,
			},
			compensate: Caller{
				name:    s.Compensate,
				fn:      handlers[s.Compensate],
				payload: s.CompensatePayload,
			},
			compensateCh:     make(chan string),
			actionCh:         make(chan string),
			actionNotify:     make(map[string]struct{}),
			compensateNotify: make(map[string]struct{}),
			closed:           make(chan struct{}),
		}
		stepMap[s.Name] = stp
	}

	for _, stp := range steps {
		for _, dp := range stp.ActionDepend {
			if _, ok := stepMap[dp]; !ok {
				return nil, errors.New(fmt.Sprintf("depend not exist: %s", dp))
			}
		}
	}

	for _, step := range steps {
		if len(step.ActionDepend) == 0 {
			t.root.next = append(t.root.next, stepMap[step.Name])
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, t.root)
			t.root.compensatePrevious = append(t.root.compensatePrevious, stepMap[step.Name])
			stepMap[step.Name].compensateNext = append(stepMap[step.Name].compensateNext, t.root)
			continue
		}
		for _, d := range step.ActionDepend {
			stepMap[step.Name].previous = append(stepMap[step.Name].previous, stepMap[d])
			stepMap[d].next = append(stepMap[d].next, stepMap[step.Name])
			stepMap[step.Name].compensateNext = append(stepMap[step.Name].compensateNext, stepMap[d])
			stepMap[d].compensatePrevious = append(stepMap[d].compensatePrevious, stepMap[step.Name])
		}
	}

	if len(t.root.next) == 0 {
		return nil, errors.New("exist circle depend, root next can not be empty")
	}

	for _, stp := range stepMap {
		if stp.isCircleDepend() {
			return nil, errors.New("exist circle depend")
		}
	}

	for _, stp := range stepMap {
		t.steps = append(t.steps, stp)
	}

	for _, opt := range opts {
		opt(t)
	}

	return t, nil
}

func (t *TransSaga) buildRootStep(opts ...StepOption) *Step {
	root := &Step{
		saga:  t,
		name:  rootStepName,
		state: StepPending,
		action: Caller{
			fn: func(context.Context, []byte) error { return nil },
		},
		compensate: Caller{
			fn: func(context.Context, []byte) error {
				ctx := context.Background()
				if err := t.changeExecuteStateFinished(ctx); err != nil {
					logx.Errorf(ctx, "changeExecuteStateFinished fail: %v", err)
				}
				t.close()
				return nil
			},
		},
		actionCh:         make(chan string),
		compensateCh:     make(chan string),
		closed:           make(chan struct{}),
		actionNotify:     make(map[string]struct{}),
		compensateNotify: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(root)
	}
	return root
}

func (t *TransSaga) changeExecuteStateFinished(ctx context.Context) error {
	return t.storage.UpdateTransExecuteStateFinished(ctx, t.id)
}

func (t *TransSaga) close() {
	for _, stp := range t.steps {
		close(stp.closed)
	}
	close(t.root.closed)
	close(t.errCh)
	close(t.done)
}
