package dtx

import (
	"context"
	"sync"
	"sync/atomic"

	"purchase/pkg/retry"

	"purchase/infra/logx"
)

func NewSaga() *Saga {
	s := &Saga{}
	return s
}

type Saga struct {
	head  *Step
	steps []*Step
	//order   map[string][]string
	state   atomic.Int32 // 0 - 执行中， 1 - 失败， 2 - 成功
	errCh   chan error   // 正向执行的错误channel，容量为1
	storage TransStorage
}

func (s *Saga) Exec(ctx context.Context) error {
	s.AsyncExec(ctx)
	return <-s.errCh
}

func (s *Saga) AsyncExec(ctx context.Context) {
	err := s.sync(ctx)
	if err != nil {
		s.errCh <- err
	}
	for _, step := range s.steps {
		go step.runAction(ctx)
		go step.runCompensate(ctx)
	}
	s.head.actionCh <- struct{}{}
}

func (s *Saga) sync(ctx context.Context) error {
	tx := NewTrans()
	branchList := make([]*Branch, 0)
	for _, s := range s.steps {
		branch := NewBranch(s.name)
		branchList = append(branchList, branch)
	}
	err := s.storage.SaveTrans(ctx, tx)
	if err != nil {
		return err
	}
	err = s.storage.SaveBranch(ctx, branchList)
	if err != nil {
		return err
	}
	return nil
}

func (s *Saga) tryUpdateSuccess(ctx context.Context) {
	for _, stp := range s.steps {
		if !stp.isSuccess() {
			return
		}
	}
	if err := s.storage.UpdateTransState(ctx, NewTrans(), "tx_success"); err != nil {
		logx.Errorf(ctx, "update branch state fail: %v", err)
	}
	close(s.errCh)
}

func (s *Saga) close() {
	for _, stp := range s.steps {
		close(stp.closed)
	}
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
	StepPending           StepStatus = 0
	StepInAction          StepStatus = 1
	StepActionFail        StepStatus = 2
	StepActionSuccess     StepStatus = 3
	StepInCompensate      StepStatus = 4
	StepCompensateSuccess StepStatus = 5
	StepCompensateFail    StepStatus = 6
	// StepSuccess           StepStatus = 7
	// StepFailed            StepStatus = 8
)

func (s *Step) isSuccess() bool {
	return s.state == StepActionSuccess
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

			s.mu.Lock()
			if !s.shouldRunAction() {
				s.mu.Unlock()
				continue
			}
			s.state = StepInAction
			s.mu.Unlock()

			err := s.action.run(ctx)
			if err != nil {
				logx.Errorf(ctx, "step[%s] run action fail: %v", s.name, err)
				s.onActionFail(ctx, err)
				return
			}
			s.onActionSuccess(ctx)
			return
		}
	}
}

func (s *Step) shouldRunAction() bool {
	if s.saga.state.Load() != 0 {
		return false
	}
	return s.state == StepPending
}

func (s *Step) onActionSuccess(ctx context.Context) {
	s.mu.Lock()
	s.state = StepActionSuccess
	s.mu.Unlock()
	err := s.saga.storage.UpdateBranchState(ctx, nil, "action_success")
	if err != nil {
		logx.Errorf(ctx, "update branch state fail: %v", err)
		// s.saga.errCh <- err
	}
	s.saga.tryUpdateSuccess(ctx)
	for _, stp := range s.next {
		stp.actionCh <- struct{}{}
	}
}

func (s *Step) onActionFail(ctx context.Context, err error) {
	s.mu.Lock()
	s.state = StepActionFail
	s.mu.Unlock()
	// 修改全局事务失败，控制多个同时fail
	if !s.saga.state.CompareAndSwap(0, 1) {
		return
	}
	s.saga.errCh <- err // 返回第一个失败的错误, 容量为1，不阻塞

	if err := s.saga.storage.UpdateBranchState(ctx, NewBranch(s.name), "action_fail"); err != nil {
		// todo 这里错误的处理，
		// todo 要放在事务里吗
		logx.Errorf(ctx, "update branch state fail: %v", err)
	}
	if err := s.saga.storage.UpdateTransState(ctx, NewTrans(), "tx_fail"); err != nil {
		// todo 这里错误的处理，
		// todo 要放在事务里吗
		logx.Errorf(ctx, "update branch state fail: %v", err)
	}
	s.compensateCh <- struct{}{} // 开始回滚
}

// isRunActionFinished 正向执行是否结束（成功或失败）
func (s *Step) isRunActionFinished() bool {
	return s.state == StepActionSuccess || s.state == StepActionFail
}

func (s *Step) needCompensate() bool {
	return s.state == StepActionSuccess || s.state == StepInAction || s.state == StepActionFail
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
			// 前序依赖step需要全部回滚完成或待执行
			for _, stp := range s.compensatePrevious {
				if stp.needCompensate() {
					stp.compensateCh <- struct{}{}
					continue
				}
			}

			if !s.needCompensate() {
				continue
			}

			// 阻塞，直到当前action执行结束
			for !s.isRunActionFinished() {
			}

			s.mu.Lock()
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
			s.onCompensateSuccess(ctx)
		}
	}
}

func (s *Step) onCompensateSuccess(ctx context.Context) {
	s.mu.Lock()
	s.state = StepCompensateSuccess
	s.mu.Unlock()
	if err := s.saga.storage.UpdateBranchState(ctx, NewBranch(s.name), "compensate_success"); err != nil {
		// todo 这里错误的处理，
		// todo 要放在事务里吗
		logx.Errorf(ctx, "update branch state fail: %v", err)
	}
	for _, stp := range s.concurrent {
		stp.compensateCh <- struct{}{}
	}
	for _, stp := range s.compensateNext {
		stp.compensateCh <- struct{}{}
	}
}

func (s *Step) onCompensateFail(ctx context.Context) {
	s.mu.Lock()
	s.state = StepCompensateFail
	s.mu.Unlock()
	if err := s.saga.storage.UpdateBranchState(ctx, NewBranch(s.name), "compensate_fail"); err != nil {
		// todo 这里错误的处理，
		// todo 要放在事务里吗
		logx.Errorf(ctx, "update branch state fail: %v", err)
	}
	// todo 更新db任务状态，人工介入, 其他回滚是否要继续执行
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
