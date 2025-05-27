package dtx

import (
	"context"
	"errors"
	"sync"
	"time"

	"purchase/infra/logx"
	"purchase/pkg/retry"
)

type Step struct {
	id                 string
	saga               *TransSaga
	action             Caller
	compensate         Caller
	state              StepStatus
	name               string
	previous           []*Step
	next               []*Step
	compensatePrevious []*Step // 回滚依赖
	compensateNext     []*Step // 回滚下一step
	actionCh           chan string
	compensateCh       chan string
	closed             chan struct{}
	mu                 sync.Mutex
	actionNotify       map[string]struct{}
	compensateNotify   map[string]struct{}
}

func (s *Step) isSuccess() bool {
	return s.state == StepActionSuccess
}

// IsReadyToRunAction 是否前序依赖都action完成了
func (s *Step) isReadyToRunAction() bool {
	for _, actionDepend := range s.previous {
		if _, ok := s.actionNotify[actionDepend.name]; !ok {
			return false
		}
	}
	return true
}

// IsReadyToRunCompensate 是否前序依赖都Compensate完成了
func (s *Step) isReadyToRunCompensate() bool {
	for _, compensateDepend := range s.compensatePrevious {
		if _, ok := s.compensateNotify[compensateDepend.name]; !ok {
			return false
		}
	}
	return true
}

func (s *Step) runAction(ctx context.Context) {
	for {
		select {
		case <-s.closed:
			logx.Infof(ctx, "trans closed action准备退出: %s", s.name)
			return
		case stepName := <-s.actionCh:
			logx.Infof(ctx, "收到action: %s", s.name)

			// 已经失败了则不执行
			if s.saga.isFailed() {
				logx.Infof(ctx, "已经失败不执行"+s.name)
				break
			}

			s.actionNotify[stepName] = struct{}{}
			if !s.isReadyToRunAction() {
				break
			}

			if s.saga.isFromDB {
				s.handleDBAction(ctx)
			} else {
				s.handleAction(ctx)
			}
		}
	}
}

func (s *Step) handleAction(ctx context.Context) {
	if s.state != StepPending {
		logx.Errorf(ctx, "handleAction Step状态异常（%s）", s.state)
		return
	}

	s.changeState(ctx, StepInAction)
	err := s.action.run(ctx)
	if err != nil {
		logx.Errorf(ctx, "step[%s] run action fail: %v", s.name, err)
		s.onActionFail(ctx, err)
	} else {
		s.onActionSuccess(ctx)
	}
}

func (s *Step) handleDBAction(ctx context.Context) {
	switch s.state {
	case StepPending:
		s.changeState(ctx, StepInAction)
		err := s.action.run(ctx)
		if err != nil {
			logx.Errorf(ctx, "step[%s] run action fail: %v", s.name, err)
			s.onActionFail(ctx, err)
		} else {
			s.onActionSuccess(ctx)
		}
	case StepInAction: // 需要action保持幂等
		err := s.action.run(ctx)
		if err != nil {
			logx.Errorf(ctx, "step[%s] run action fail: %v", s.name, err)
			s.onActionFail(ctx, err)
		} else {
			s.onActionSuccess(ctx)
		}
	case StepActionSuccess:
		s.onActionSuccess(ctx)
	case StepActionFail:
		s.onActionFail(ctx, errors.New("unknown error"))
	case StepInCompensate: // 需要Compensate保持幂等
		err := retry.Run(func() error {
			return s.compensate.run(ctx)
		}, 2)
		if err != nil {
			logx.Errorf(ctx, "step[%s] run compensate fail after retry 2 times: %v", s.name, err)
			s.onCompensateFail(ctx)
		} else {
			s.onCompensateSuccess(ctx)
		}
	case StepCompensateSuccess:
		s.onCompensateSuccess(ctx)
	case StepCompensateFail:
		s.onCompensateFail(ctx)
	}
}

func (s *Step) handleCompensate(ctx context.Context) {
	switch s.state {
	case StepActionSuccess:
		s.changeState(ctx, StepInCompensate)
		err := retry.Run(func() error {
			return s.compensate.run(ctx)
		}, 2)
		if err != nil {
			logx.Errorf(ctx, "step[%s] run compensate fail after retry 2 times: %v", s.name, err)
			s.onCompensateFail(ctx)
		} else {
			s.onCompensateSuccess(ctx)
		}
	case StepActionFail:
		s.changeState(ctx, StepCompensateSuccess)
		for _, stp := range s.compensateNext {
			stp.compensateCh <- s.name
		}
	default:
		logx.Errorf(ctx, "handleCompensate Step状态异常（%s）", s.state)
	}
}

func (s *Step) onActionSuccess(ctx context.Context) {
	s.changeState(ctx, StepActionSuccess)
	s.saga.tryUpdateSuccess(ctx)
	// fmt.Println(fmt.Sprintf("step[%s] action done success，状态变更完成， 开始通知下游", s.name))
	for _, stp := range s.next {
		stp.actionCh <- s.name
	}
}

func (s *Step) onActionFail(ctx context.Context, err error) {
	s.changeState(ctx, StepActionFail)
	// 修改全局事务失败，控制多个同时fail
	if !s.saga.state.CompareAndSwap(SagaStateExecuting, SagaStateFail) {
		return
	}

	s.saga.errCh <- err // 返回第一个失败的错误, 容量为1，不阻塞

	if err := s.saga.syncStateChange(ctx); err != nil {
		logx.Errorf(ctx, "update branch state fail: %v", err)
	}
	logx.Infof(ctx, "step[%s] action done fail状态变更完成， 开始回滚", s.name)
	s.compensateCh <- s.name // 开始回滚
}

func (s *Step) syncStateChange(ctx context.Context) {
	err := s.saga.storage.UpdateBranchState(ctx, s.id, s.state.String())
	if err != nil {
		logx.Errorf(ctx, "sync branch state change fail, err:%v", err)
	}
}

func (s *Step) noNeedCompensate() bool {
	return s.state == StepPending
}

func (s *Step) runCompensate(ctx context.Context) {
	for {
		select {
		case <-s.closed:
			logx.Infof(ctx, "trans closed compensate准备退出: %s", s.name)
			return
		case stepName := <-s.compensateCh:
			// 前序依赖step需要全部回滚完成或待执行
			logx.Infof(ctx, "收到compensate: %s", s.name)

			// 无需回滚，直接通知下游
			if s.noNeedCompensate() {
				s.changeState(ctx, StepCompensateSuccess)
				for _, stp := range s.compensateNext {
					stp.compensateCh <- s.name
				}
				break
			}

			// 上游依赖回滚还未完成，通知未完成的回滚
			s.compensateNotify[stepName] = struct{}{}
			if !s.isReadyToRunCompensate() {
				for _, stp := range s.compensatePrevious {
					if _, ok := s.compensateNotify[stp.name]; !ok {
						stp.compensateCh <- s.name
					}
				}
				break
			}

			// 阻塞，直到当前action执行结束
			for s.state == StepInAction {
				time.Sleep(10 * time.Millisecond)
			}

			s.handleCompensate(ctx)
		}
	}
}

func (s *Step) onCompensateSuccess(ctx context.Context) {
	s.changeState(ctx, StepCompensateSuccess)
	logx.Infof(ctx, "step[%s] compensate done success状态变更完成，通知上有回滚", s.name)
	for _, stp := range s.compensateNext {
		stp.compensateCh <- s.name
	}
}

func (s *Step) onCompensateFail(ctx context.Context) {
	s.changeState(ctx, StepCompensateFail)
	logx.Errorf(ctx, "step[%s] compensate done fail状态变更完成，阻塞住了", s.name)
	if s.saga.strictCompensate {
		s.saga.close() // 直接退出
		return
	}
	for _, stp := range s.compensateNext {
		stp.compensateCh <- s.name
	}
}

func (s *Step) changeState(ctx context.Context, state StepStatus) {
	s.mu.Lock()
	s.state = state
	s.mu.Unlock()
	s.syncStateChange(ctx)
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

type Caller struct {
	fn      func(ctx context.Context, payload []byte) error
	name    string
	payload []byte
	signal  chan struct{}
}

func (c *Caller) run(ctx context.Context) error {
	return c.fn(ctx, c.payload)
}

type CompensateError struct {
	error
}

func (c CompensateError) Error() string {
	return c.error.Error()
}

func NewCompensateError(e error) error {
	return CompensateError{
		error: e,
	}
}

func IsCompensateError(e error) bool {
	var compensateError CompensateError
	ok := errors.Is(e, &compensateError)
	return ok
}
