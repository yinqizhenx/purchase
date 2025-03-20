package dtx

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"purchase/infra/utils"

	"purchase/pkg/retry"

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
	id    string
	root  *Step
	steps []*Step
	// order   map[string][]string
	state   atomic.Int32 // 0 - 执行中， 1 - 失败， 2 - 成功
	errCh   chan error   // 正向执行的错误channel，容量为1
	storage TransStorage
	// timeout  time.Duration
	done     chan struct{}
	isFromDB bool
}

func (t *TransSaga) Exec(ctx context.Context) error {
	t.AsyncExec()
	return <-t.errCh
}

func (t *TransSaga) AsyncExec() {
	ctx := context.Background()
	utils.SafeGo(ctx, func() {
		// ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
		// defer cancel()

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
		t.root.actionCh <- struct{}{}

		select {
		// case <-ctx.Done():
		// 	fmt.Println("context done 退出AsyncExec")
		case <-t.done:
			fmt.Println("执行完成 退出AsyncExec")
		}
	})
}

func (t *TransSaga) sync(ctx context.Context) error {
	tx := &Trans{
		TransID:    t.id,
		Name:       "test",
		State:      "pending",
		FinishedAt: time.Now(),
		CreatedAt:  time.Now(),
		CreatedBy:  "yinqizhen",
		UpdatedBy:  "dd",
		UpdatedAt:  time.Now(),
	}
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

		branch := &Branch{
			BranchID:         stp.id,
			TransID:          stp.saga.id,
			Type:             "action",
			State:            "pending",
			Name:             stp.name,
			Action:           stp.action.name,
			Compensate:       stp.compensate.name,
			Payload:          "stp.action.payload",
			ActionDepend:     actionDepend,
			CompensateDepend: compensateDepend,
			FinishedAt:       time.Now(),
			IsDead:           false,
			CreatedAt:        time.Now(),
			CreatedBy:        "uy",
			UpdatedBy:        "dd",
			UpdatedAt:        time.Now(),
		}
		branchList = append(branchList, branch)
	}
	err := t.storage.SaveTrans(ctx, tx)
	if err != nil {
		return err
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
	if !t.state.CompareAndSwap(0, 2) {
		return
	}
	if err := t.syncStateChange(ctx); err != nil {
		logx.Errorf(ctx, "update branch state success: %v", err)
	}
	t.close()
}

func (t *TransSaga) syncStateChange(ctx context.Context) error {
	var newState string
	switch t.state.Load() {
	case 0:
		newState = "pending"
	case 1:
		newState = "fail"
	case 2:
		newState = "success"
	default:
		return errors.New(fmt.Sprintf("unknown saga state : %d", t.state.Load()))
	}
	return t.storage.UpdateTransState(ctx, t.id, newState)
}

func (t *TransSaga) isFailed() bool {
	return t.state.Load() == 1
}

func (t *TransSaga) build(steps []*Branch, handlers map[string]func(context.Context, []byte) error, opts ...Option) (*TransSaga, error) {
	t.root = t.buildRootStep()
	stepMap := make(map[string]*Step)
	for _, s := range steps {
		stp := &Step{
			id:    s.BranchID,
			saga:  t,
			name:  s.Name,
			state: s.State,
			action: Caller{
				name:    s.Action,
				fn:      handlers[s.Action],
				payload: []byte(s.Payload),
			},
			compensate: Caller{
				name:    s.Compensate,
				fn:      handlers[s.Compensate],
				payload: []byte(s.Payload),
			},
			actionCh:     make(chan struct{}),
			compensateCh: make(chan struct{}),
			closed:       make(chan struct{}),
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

	// if t.timeout == 0 {
	// 	t.timeout = defaultTimeout
	// }
	return t, nil
}

func (t *TransSaga) buildRootStep(opts ...StepOption) *Step {
	root := &Step{
		id:    uuid.NewString(),
		saga:  t,
		name:  rootStepName,
		state: StepPending,
		action: Caller{
			fn: func(context.Context, []byte) error { return nil },
		},
		compensate: Caller{
			fn: func(context.Context, []byte) error {
				t.close()
				return nil
			},
		},
		actionCh:     make(chan struct{}),
		compensateCh: make(chan struct{}),
		closed:       make(chan struct{}),
	}
	for _, opt := range opts {
		opt(root)
	}
	return root
}

func (t *TransSaga) close() {
	for _, stp := range t.steps {
		close(stp.closed)
	}
	close(t.root.closed)
	close(t.errCh)
	close(t.done)
}

type Step struct {
	id                    string
	saga                  *TransSaga
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

type StepStatus string

func (s StepStatus) String() string {
	return string(s)
}

const (
	StepPending           StepStatus = "pending"
	StepInAction          StepStatus = "in action"
	StepActionFail        StepStatus = "action fail"
	StepActionSuccess     StepStatus = "action success"
	StepInCompensate      StepStatus = "in compensate"
	StepCompensateSuccess StepStatus = "compensate success"
	StepCompensateFail    StepStatus = "compensate fail"
	// StepSuccess           StepStatus = 7
	// StepFailed            StepStatus = 8
)

func (s *Step) isSuccess() bool {
	return s.state == StepActionSuccess
}

func (s *Step) runAction(ctx context.Context) {
	for {
		select {
		// case <-ctx.Done():
		// 	logx.Errorf(ctx, "saga tx context done: %v", s.name)
		// 	return
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
	// fmt.Println(fmt.Sprintf("step[%s] action done success，状态变更完成， 开始通知下游", s.name))
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
	err := s.saga.storage.UpdateBranchState(ctx, s.id, s.state.String())
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
		// case <-ctx.Done():
		// 	logx.Errorf(ctx, "saga tx context done: %v", ctx.Done())
		// 	return
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
