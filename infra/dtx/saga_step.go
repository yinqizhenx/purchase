package dtx

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"purchase/infra/idempotent"
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

func (s *Step) isRoot() bool {
	return s.name == rootStepName
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

			s.handleAction(ctx)
		}
	}
}

func (s *Step) handleAction(ctx context.Context) {
	if s.state != StepPending {
		logx.Errorf(ctx, "handleAction Step状态异常（%s）", s.state)
		return
	}

	s.changeState(ctx, StepInAction)
	err := s.executeActionWithIdempotent(ctx)
	if err != nil {
		logx.Errorf(ctx, "step[%s] run action fail: %v", s.name, err)
		s.onActionFail(ctx, err)
	} else {
		s.onActionSuccess(ctx)
	}
}

// executeActionWithIdempotent 执行action, 保证幂等性
// 方式1: 插入一条幂等数据，成功则执行业务逻辑，整个操作在事务里执行，业务逻辑成功才提交数据库事务
// 极端场景：如果action业务逻辑执行成功（包含rpc请求），但是在提交数据库事务之前，服务挂了，这个时候没有幂等键插入，但是实际已经执行了，会导致重复执行
// 方式2: 先插入一条幂等数据状态pending，然后执行业务逻辑，成功则修改状态为done, 失败则删除幂等key， 整个操作都是独立的，不在同一个事务里
// 业务逻辑开始执行后，服务挂了，无法判断业务逻辑是否执行成功，服务重启后会看作已执行，需要人工介入
// 若修改状态为done失败或删除幂等key失败， 则告警，人工介入， 可避免重复执行的问题
// 采用方式2
func (s *Step) executeActionWithIdempotent(ctx context.Context) error {
	key := fmt.Sprintf("dtx_%s", s.id)
	exist, err := s.saga.idempotentSrv.SetKeyPendingWithDDL(ctx, key, 0)
	if err != nil {
		return err
	}
	if exist {
		state, err := s.saga.idempotentSrv.GetKeyState(ctx, key)
		if err != nil {
			return err
		}
		if state == idempotent.Pending {
			logx.Errorf(ctx, "已存在待处理的重复Action[%s]", key)
		}
		return nil
	}
	err = s.action.run(ctx)
	if err != nil {
		return s.saga.idempotentSrv.RemoveFailKey(ctx, key)
	}
	return s.saga.idempotentSrv.UpdateKeyDone(ctx, key)
}

func (s *Step) handleCompensate(ctx context.Context) {
	switch s.state {
	case StepActionSuccess:
		s.changeState(ctx, StepInCompensate)
		err := retry.Run(func() error {
			return s.executeCompensateWithIdempotent(ctx)
		}, 2)
		if err != nil {
			logx.Errorf(ctx, "step[%s] run compensate fail after retry 2 times: %v", s.name, err)
			s.onCompensateFail(ctx)
		} else {
			s.onCompensateSuccess(ctx)
		}
	case StepActionFail:
		// action 执行失败了，默认看作不需要回滚
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
	logx.Infof(ctx, "step[%s] action fail状态变更完成，开始回滚", s.name)
	s.compensateCh <- s.name // 开始回滚
}

func (s *Step) syncStateChange(ctx context.Context) {
	err := s.saga.storage.UpdateBranchState(ctx, s.id, s.state.String())
	if err != nil {
		logx.Errorf(ctx, "sync branch state to db fail, err:%v", err)
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

// executeCompensateWithIdempotent 执行Compensate, 保证幂等性
func (s *Step) executeCompensateWithIdempotent(ctx context.Context) error {
	key := fmt.Sprintf("dtx_%s", s.id)
	ok, err := s.saga.idempotentSrv.SetKeyPendingWithDDL(ctx, key, 0)
	if err != nil {
		return err
	}
	if ok {
		state, err := s.saga.idempotentSrv.GetKeyState(ctx, key)
		if err != nil {
			return err
		}
		if state == idempotent.Pending {
			logx.Errorf(ctx, "已存在待处理的重复Compensate[%s]", key)
		}
		return nil
	}
	err = s.compensate.run(ctx)
	if err != nil {
		return s.saga.idempotentSrv.RemoveFailKey(ctx, key)
	}
	return s.saga.idempotentSrv.UpdateKeyDone(ctx, key)
}

func (s *Step) onCompensateSuccess(ctx context.Context) {
	s.changeState(ctx, StepCompensateSuccess)
	logx.Infof(ctx, "step[%s] compensate success状态变更完成，通知上有回滚", s.name)
	for _, stp := range s.compensateNext {
		stp.compensateCh <- s.name
	}
}

func (s *Step) onCompensateFail(ctx context.Context) {
	s.changeState(ctx, StepCompensateFail)
	logx.Errorf(ctx, "step[%s] compensate fail状态变更完成，阻塞住了", s.name)
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
