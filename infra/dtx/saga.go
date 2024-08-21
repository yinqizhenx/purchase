package dtx

import (
	"context"
	"sync/atomic"

	"purchase/infra/logx"
)

type Saga struct {
	head  *Step
	steps []*Step
	order map[string][]string
	state atomic.Int32 // 0 - 执行中， 1 - 失败， 3 - 成功
	errCh chan error
}

func (s *Saga) Exec(ctx context.Context) error {
	s.AsyncExec(ctx)
	return <-s.errCh
}

func (s *Saga) AsyncExec(ctx context.Context) {
	for _, step := range s.steps {
		go step.runAction(ctx)
	}
	s.head.actionCh <- struct{}{}
}

//
// func (s *Saga) Rollback(ctx context.Context) error {
// 	for _, step := range s.steps {
// 		err := retry.Run(func() error {
// 			return step.runCompensate(ctx)
// 		}, 2)
// 		if err != nil {
// 			logx.Errorf(ctx, "saga tx rollback fail: %v", err)
// 		}
// 	}
// 	return nil
// }

type Step struct {
	*Saga
	action       Caller
	compensate   Caller
	status       StepStatus
	name         string
	previous     []*Step
	next         []*Step
	actionCh     chan struct{}
	compensateCh chan struct{}
	closed       chan struct{}
}

type StepStatus string

const (
	StepStatusRunning           = "running"
	StepStatusReadyToCompensate = "ready_to_compensate"
	StepStatusDone              = "done"
)

func (s *Step) runAction(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logx.Errorf(ctx, "saga tx context done: %v", ctx.Done())
			return
		case <-s.closed:
			return
		case <-s.actionCh:
			err := s.action.run(ctx)
			if err != nil {
				logx.Errorf(ctx, "step[%s] run action fail: %v", s.name, err)
				s.onActionFail(ctx)
				return
			}
			s.onActionSuccess(ctx)
			return
		}
	}
}

func (s *Step) onActionSuccess(ctx context.Context) {
	for _, stp := range s.next {
		stp.actionCh <- struct{}{}
	}
}

func (s *Step) onActionFail(ctx context.Context) {
	// 控制多个同时fail
	if !s.state.CompareAndSwap(0, 1) {
		return
	}
	s.runCompensate(ctx)
	for _, stp := range s.previous {
		stp.compensateCh <- struct{}{}
	}
	// 未执行的step退出
	for _, stp := range s.next {
		close(stp.closed)
	}
}

func (s *Step) runCompensate(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logx.Errorf(ctx, "saga tx context done: %v", ctx.Done())
			return
		case <-s.closed:
			return
		case <-s.compensateCh:
			err := s.compensate.run(ctx)
			if err != nil {
				logx.Errorf(ctx, "step[%s] run compensate fail: %v", s.name, err)
				s.onCompensateFail(ctx)
				return
			}
			s.onCompensateSuccess(ctx)
			return
		}
	}
}

func (s *Step) onCompensateSuccess(ctx context.Context) {
	for _, stp := range s.previous {
		stp.compensateCh <- struct{}{}
	}
}

func (s *Step) onCompensateFail(ctx context.Context) {
	return
}

type Caller struct {
	fn      func(ctx context.Context, payload []byte) error
	name    string
	payload []byte
	signal  chan struct{}
}

func (c *Caller) run(ctx context.Context) error {
	return c.fn(ctx, c.payload)
}
