package dtx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"purchase/infra/utils"

	"purchase/pkg/retry"

	"purchase/infra/logx"
)

const (
	defaultTimeout = 10 * time.Second
)

func NewSaga() *Saga {
	s := &Saga{}
	return s
}

type Saga struct {
	id    string
	root  *Step
	steps []*Step
	// order   map[string][]string
	state   atomic.Int32 // 0 - 执行中， 1 - 失败， 2 - 成功
	errCh   chan error   // 正向执行的错误channel，容量为1
	storage TransStorage
	timeout time.Duration
	done    chan struct{}
}

func (s *Saga) Exec(ctx context.Context) error {
	s.AsyncExec(ctx)
	return <-s.errCh
}

func (s *Saga) AsyncExec(ctx context.Context) {
	utils.SafeGo(ctx, func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		defer cancel()

		err := s.sync(ctx)
		if err != nil {
			s.errCh <- err
			return
		}

		utils.SafeGo(ctx, func() {
			s.root.runAction(ctx)
		})
		utils.SafeGo(ctx, func() {
			s.root.runCompensate(ctx)
		})
		for _, step := range s.steps {
			stp := step // 重新赋值，防止step引用变化
			utils.SafeGo(ctx, func() {
				stp.runAction(ctx)
			})
			utils.SafeGo(ctx, func() {
				stp.runCompensate(ctx)
			})
		}
		s.root.actionCh <- struct{}{}

		select {
		case <-ctx.Done():
			fmt.Println("context done 退出AsyncExec")
		case <-s.done:
			fmt.Println("执行完成 退出AsyncExec")
		}
	})
}

func (s *Saga) sync(ctx context.Context) error {
	tx := &Trans{
		TransID:    s.id,
		Name:       "test",
		State:      "pending",
		FinishedAt: time.Now(),
		CreatedAt:  time.Now(),
		CreatedBy:  "yinqizhen",
		UpdatedBy:  "dd",
		UpdatedAt:  time.Now(),
	}
	branchList := make([]*Branch, 0)
	for _, stp := range s.steps {
		actionDepend := make([]string, 0)
		compensateDepend := make([]string, 0)
		for _, p := range stp.previous {
			actionDepend = append(actionDepend, p.name)
		}
		for _, n := range stp.compensatePrevious {
			compensateDepend = append(compensateDepend, n.name)
		}

		branch := &Branch{
			BranchID:         stp.id,
			TransID:          stp.saga.id,
			Type:             "action",
			State:            "pending",
			Name:             stp.name,
			Action:           stp.action.name,
			Compensate:       stp.compensate.name,
			Payload:          "stp.action.payload",
			ActionDepend:     strings.Join(actionDepend, ","),
			CompensateDepend: strings.Join(compensateDepend, ","),
			FinishedAt:       time.Now(),
			IsDead:           false,
			CreatedAt:        time.Now(),
			CreatedBy:        "uy",
			UpdatedBy:        "dd",
			UpdatedAt:        time.Now(),
		}
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
	// 修改全局事务成功，防止并发同时close(s.errCh)
	if !s.state.CompareAndSwap(0, 2) {
		return
	}
	if err := s.syncStateChange(ctx); err != nil {
		logx.Errorf(ctx, "update branch state success: %v", err)
	}
	s.close()
}

func (s *Saga) syncStateChange(ctx context.Context) error {
	var newState string
	switch s.state.Load() {
	case 0:
		newState = "pending"
	case 1:
		newState = "fail"
	case 2:
		newState = "success"
	default:
		return errors.New(fmt.Sprintf("unknown saga state : %d", s.state.Load()))
	}
	return s.storage.UpdateTransState(ctx, s.id, newState)
}

func (s *Saga) isFailed() bool {
	return s.state.Load() == 1
}

func (s *Saga) close() {
	for _, stp := range s.steps {
		close(stp.closed)
	}
	close(s.root.closed)
	close(s.errCh)
	close(s.done)
}

type Step struct {
	id                    string
	saga                  *Saga
	action                Caller
	compensate            Caller
	state                 StepStatus
	name                  string
	previous              []*Step
	next                  []*Step
	compensatePrevious    []*Step // 回滚依赖
	compensateNext        []*Step // 回滚下一step
	actionCh              chan struct{}
	compensateCh          chan struct{}
	closed                chan struct{}
	mu                    sync.Mutex
	actionNotifyCount     int
	compensateNotifyCount int
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
			logx.Errorf(ctx, "saga tx context done: %v", s.name)
			return
		case <-s.closed:
			fmt.Println(fmt.Sprintf("trans closed action准备退出: %s", s.name))
			return
		case <-s.actionCh:
			fmt.Println(fmt.Sprintf("收到action: %s", s.name))
			s.actionNotifyCount++

			// 已经失败了则不执行
			if s.saga.isFailed() {
				fmt.Println("已经失败不执行" + s.name)
				return
			}

			if s.actionNotifyCount < len(s.previous) {
				fmt.Println(fmt.Sprintf("当前action count %d, name %s", s.actionNotifyCount, s.name))
				continue
			}

			s.handleAction(ctx)
		}
	}
}

func (s *Step) handleAction(ctx context.Context) {
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
	case StepInAction:
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
	case StepInCompensate:
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
	}
}

func (s *Step) onActionSuccess(ctx context.Context) {
	s.changeState(ctx, StepActionSuccess)
	s.saga.tryUpdateSuccess(ctx)
	fmt.Println(fmt.Sprintf("step[%s] action done success，状态变更完成， 开始通知下游", s.name))
	for _, stp := range s.next {
		stp.actionCh <- struct{}{}
	}
}

func (s *Step) onActionFail(ctx context.Context, err error) {
	s.changeState(ctx, StepActionFail)
	// 修改全局事务失败，控制多个同时fail
	if !s.saga.state.CompareAndSwap(0, 1) {
		return
	}

	s.saga.errCh <- err // 返回第一个失败的错误, 容量为1，不阻塞

	if err := s.saga.syncStateChange(ctx); err != nil {
		logx.Errorf(ctx, "update branch state fail: %v", err)
	}
	fmt.Println(fmt.Sprintf("step[%s] action done fail状态变更完成， 开始回滚", s.name))
	s.compensateCh <- struct{}{} // 开始回滚
}

func (s *Step) syncStateChange(ctx context.Context) {
	var newState string
	switch s.state {
	case 0:
		newState = "pending"
	case 1:
		newState = "in action"
	case 2:
		newState = "action fail"
	case 3:
		newState = "action success"
	case 4:
		newState = "in compensate"
	case 5:
		newState = "compensate success"
	case 6:
		newState = "compensate fail"
	default:
		logx.Errorf(ctx, "unknown saga state : %d", s.state)
	}
	err := s.saga.storage.UpdateBranchState(ctx, s.id, newState)
	if err != nil {
		logx.Errorf(ctx, "sync branch state change fail, err:%v", err)
	}
}

// isRunActionFinished 正向执行是否结束（成功或失败）
func (s *Step) isRunActionFinished() bool {
	return s.state == StepActionSuccess || s.state == StepActionFail
}

func (s *Step) shouldCompensate() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state == StepActionSuccess || s.state == StepInAction || s.state == StepActionFail
}

func (s *Step) noNeedCompensate() bool {
	return s.state == StepPending
}

func (s *Step) runCompensate(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logx.Errorf(ctx, "saga tx context done: %v", ctx.Done())
			return
		case <-s.closed:
			fmt.Println(fmt.Sprintf("trans closed compensate准备退出: %s", s.name))
			return
		case <-s.compensateCh:
			// 前序依赖step需要全部回滚完成或待执行
			fmt.Println(fmt.Sprintf("收到compensate: %s", s.name))
			s.compensateNotifyCount++

			isPreviousCompensateDone := true
			for _, stp := range s.compensatePrevious {
				// 此刻处于pending状态的step, 后续一定不会继续往下流转，不需要回滚
				if stp.shouldCompensate() {
					isPreviousCompensateDone = false
					stp.compensateCh <- struct{}{}
				}
			}

			// 前序的补偿完成了，才能往下走
			if !isPreviousCompensateDone {
				continue
			}

			if s.compensateNotifyCount < len(s.compensatePrevious) {
				fmt.Println(fmt.Sprintf("当前compensate count %d, name %s", s.compensateNotifyCount, s.name))
				continue
			}

			// 无需补偿，直接通知下游补偿
			if s.noNeedCompensate() {
				for _, stp := range s.compensateNext {
					stp.compensateCh <- struct{}{}
				}
			}

			// 阻塞，直到当前action执行结束
			for !s.isRunActionFinished() {
			}

			s.handleCompensate(ctx)
		}
	}
}

func (s *Step) onCompensateSuccess(ctx context.Context) {
	s.changeState(ctx, StepCompensateSuccess)
	fmt.Println(fmt.Sprintf("step[%s] compensate done success状态变更完成，通知上有回滚", s.name))
	for _, stp := range s.compensateNext {
		stp.compensateCh <- struct{}{}
	}
}

func (s *Step) onCompensateFail(ctx context.Context) {
	s.changeState(ctx, StepCompensateFail)
	fmt.Println(fmt.Sprintf("step[%s] compensate done fail状态变更完成，阻塞住了", s.name))
	// todo 更新db任务状态，人工介入, 其他回滚是否要继续执行
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
