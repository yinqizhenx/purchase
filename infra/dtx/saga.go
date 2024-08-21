package dtx

import (
	"context"
	"purchase/pkg/retry"
	"sync"
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
		go step.runCompensate(ctx)
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
	saga               *Saga
	action             Caller
	compensate         Caller
	state              StepStatus
	name               string
	previous           []*Step
	concurrent         []*Step
	next               []*Step
	compensatePrevious []*Step // 回滚依赖
	compensateNext     []*Step // 回滚下一step
	actionCh           chan struct{}
	compensateCh       chan struct{}
	closed             chan struct{}
	mu                 sync.Mutex
}

type StepStatus int

const (
	StepPending      StepStatus = 0
	StepInAction     StepStatus = 1
	StepInCompensate StepStatus = 2
	StepSuccess      StepStatus = 3
	StepFailed       StepStatus = 4
)

func (s *Step) isSuccess() bool {
	return s.state == StepSuccess
}

func (s *Step) runAction(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logx.Errorf(ctx, "saga tx context done: %v", ctx.Done())
			return
		case <-s.closed:
			return
		case <-s.actionCh:
			// 等待所有依赖step执行成功
			for _, stp := range s.previous {
				if !stp.isSuccess() {
					continue
				}
			}
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
	if !s.saga.state.CompareAndSwap(0, 1) {
		return
	}
	// 当前step回滚，previous next concurrent 都需要回滚
	s.runCompensate(ctx)
	for _, stp := range s.previous {
		stp.compensateCh <- struct{}{}
	}
	for _, stp := range s.concurrent {
		stp.compensateCh <- struct{}{}
	}
	//for _, stp := range s.next {
	//	stp.compensateCh <- struct{}{}
	//}
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
			s.mu.Lock()
			// 只有正在执行的和执行成功的需要回滚
			if s.state != StepInAction && s.state != StepSuccess {
				s.mu.Unlock()
				continue
			}
			// 前序依赖step需要全部回滚完成或待执行
			if s.state == StepSuccess {
				for _, stp := range s.compensatePrevious {
					if stp.state != StepFailed && s.state != StepPending {
						s.mu.Unlock()
						continue
					}
				}
			}
			s.state = StepInCompensate
			s.mu.Unlock()
			err := retry.Run(func() error {
				return s.compensate.run(ctx)
			}, 2)
			if err != nil {
				logx.Errorf(ctx, "step[%s] run compensate fail after retry 2 times: %v", s.name, err)
				s.onCompensateFail(ctx)
				continue
			}
			s.onCompensateSuccess()
		}
	}
}

func (s *Step) onCompensateSuccess() {
	s.mu.Lock()
	s.state = StepFailed
	s.mu.Unlock()
	for _, stp := range s.previous {
		stp.compensateCh <- struct{}{}
	}
}

func (s *Step) onCompensateFail(ctx context.Context) {
	// todo 更新db任务状态，人工介入
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
