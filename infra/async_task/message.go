package async_task

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"purchase/domain/entity/async_task"
	"purchase/infra/dlock"
	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/tx"
	"purchase/infra/utils"
	"purchase/pkg/chanx"
)

const defaultMaxTaskLoad = 10

const defaultConcurrency = 5

type Handler func(ctx context.Context, payload []byte) error

type AsyncTaskMux struct {
	// msg          *message.Task
	pub mq.Publisher
	// handles 涉及到数据库变动的，需要在事务里执行
	handlers    map[string]Handler
	ch          *chanx.UnboundedChan[string]
	cron        *cron.Cron
	dal         *dal.AsyncTaskDal
	txm         *tx.TransactionManager
	sem         chan struct{}
	concurrency int
	maxTaskLoad int
	lockBuilder dlock.LockBuilder
	locks       map[string]dlock.DistributeLock
	cancel      func()
	mdw         []Middleware
	mu          sync.Mutex
}

func NewAsyncTaskMux(pub mq.Publisher, dal *dal.AsyncTaskDal, txm *tx.TransactionManager, ch *chanx.UnboundedChan[string], lockBuilder dlock.LockBuilder, opts ...Option) *AsyncTaskMux {
	h := &AsyncTaskMux{
		pub:         pub,
		handlers:    make(map[string]Handler),
		ch:          ch,
		cron:        cron.New(),
		dal:         dal,
		txm:         txm,
		concurrency: defaultConcurrency, // default 5 worker max
		maxTaskLoad: defaultMaxTaskLoad,
		lockBuilder: lockBuilder,
		locks:       make(map[string]dlock.DistributeLock),
	}
	for _, opt := range opts {
		opt(h)
	}
	h.sem = make(chan struct{}, h.concurrency)
	return h
}

func (m *AsyncTaskMux) Start(ctx context.Context) error {
	nctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	err := m.RunCron(nctx)
	if err != nil {
		logx.Error(ctx, "launch task cron handle fail", slog.Any("error", err))
		return err
	}
	m.Listen(nctx)
	return nil
}

func (m *AsyncTaskMux) Stop(ctx context.Context) error {
	if m.cancel != nil {
		m.cancel()
	}
	m.cron.Stop()
	return nil
}

func (m *AsyncTaskMux) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case taskID := <-m.ch.Out:
			select {
			case <-ctx.Done():
				return
			case m.sem <- struct{}{}:
				go func() {
					task, err := m.lockAndLoadTask(ctx, taskID)
					if err != nil {
						logx.Error(ctx, "lockAndLoadTask fail", slog.String("task", taskID), slog.Any("error", err))
						return
					}
					defer func() {
						err := m.unLockTask(ctx, taskID)
						if err != nil {
							logx.Error(ctx, "unlock task fail", slog.String("task", taskID), slog.Any("error", err))
						}
						<-m.sem
						if p := recover(); p != nil {
							logx.Error(ctx, "handle async task from listen panic", slog.Any("error", p))
						}
					}()
					if len(task) == 1 {
						err = m.Handle(ctx, task[0])
						if err != nil {
							err = m.onHandleFail(ctx, task[0])
						} else {
							err = m.onHandleSuccess(ctx, task[0])
						}
						if err != nil {
							logx.Error(ctx, "handle task fail", slog.String("task", taskID), slog.Any("error", err))
						}
					} else {
						logx.Errorf(ctx, "load task expect only 1, get %v, task id: %s", len(task), taskID)
					}
				}()
			}
		}
	}
}

func (m *AsyncTaskMux) Handle(ctx context.Context, msg *async_task.AsyncTask) error {
	if msg.IsTask(ctx) {
		return m.ProcessTask(ctx, msg)
	}
	if msg.IsEvent(ctx) {
		return m.SendMessage(ctx, msg)
	}
	return errors.New("unknown task type")
}

func (m *AsyncTaskMux) RunCron(ctx context.Context) error {
	// run every 5 second
	_, err := m.cron.AddFunc("@every 5s", m.buildCronHandler(ctx))
	if err != nil {
		return err
	}
	go m.cron.Run()
	return nil
}

func (m *AsyncTaskMux) onHandleSuccess(ctx context.Context, task *async_task.AsyncTask) error {
	err := m.dal.UpdateDone(ctx, task.TaskID)
	if err != nil {
		logx.Error(ctx, "update task state done fail", slog.String("task", task.TaskID), slog.Any("error", err))
	}
	return err
}

// onHandleFail retry in 1s
func (m *AsyncTaskMux) onHandleFail(ctx context.Context, task *async_task.AsyncTask) error {
	time.Sleep(500 * time.Millisecond)
	err := m.Handle(ctx, task)
	if err != nil {
		logx.Error(ctx, "handle task fail", slog.String("task", task.TaskID), slog.Any("error", err))
		return err
	}
	return m.onHandleSuccess(ctx, task)
}

func (m *AsyncTaskMux) buildCronHandler(ctx context.Context) func() {
	return func() {
		taskList, err := m.lockAndLoadTaskWithLimit(ctx, m.maxTaskLoad)
		if err != nil {
			logx.Error(ctx, "batch load task fail", slog.Any("error", err))
			return
		}
		defer func() {
			for _, task := range taskList {
				err := m.unLockTask(ctx, task.TaskID)
				if err != nil {
					logx.Error(ctx, "unLock task fail", slog.String("task", task.TaskID), slog.Any("error", err))
				}
			}
		}()
		wg := &sync.WaitGroup{}
		for _, task := range taskList {
			wg.Add(1)
			t := task
			fn := func(context.Context) (hErr error) {
				defer func() {
					// task 执行完成及时解锁
					hErr = m.unLockTask(ctx, task.TaskID)
					if p := recover(); p != nil {
						hErr = fmt.Errorf("handle async task from cron panic :%v", p)
					}
					// 放在解锁后面
					wg.Done()
				}()
				hErr = m.Handle(ctx, t)
				if hErr != nil {
					hErr = m.onHandleFail(ctx, t)
				} else {
					hErr = m.onHandleSuccess(ctx, t)
				}
				return
			}
			utils.SafeGo(ctx, func() {
				err = m.txm.Transaction(ctx, fn, tx.PropagationRequiresNew)
				if err != nil {
					logx.Error(ctx, "handle task fail", slog.String("task", t.TaskID), slog.Any("error", err))
				}
			})
		}
		wg.Wait()
	}
}

func (m *AsyncTaskMux) lockAndLoadTaskWithLimit(ctx context.Context, limit int) ([]*async_task.AsyncTask, error) {
	taskList, err := m.loadPendingTaskWithLimit(ctx, limit)
	if err != nil {
		return nil, err
	}
	taskIDs := make([]string, 0, len(taskList))
	for _, task := range taskList {
		taskIDs = append(taskIDs, task.TaskID)
	}
	taskList, err = m.lockAndLoadTask(ctx, taskIDs...)
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

func (m *AsyncTaskMux) lockAndLoadTask(ctx context.Context, taskIDs ...string) ([]*async_task.AsyncTask, error) {
	lockedIDs := make([]string, 0)
	for _, taskID := range taskIDs {
		err := m.lockTask(ctx, taskID)
		if err == nil {
			lockedIDs = append(lockedIDs, taskID)
		}
	}
	tasks, err := m.dal.FindAllPending(ctx, taskIDs)
	if err != nil {
		for _, taskID := range taskIDs {
			if err := m.unLockTask(ctx, taskID); err != nil {
				logx.Errorf(ctx, "unlock task fail, err:%v", err)
			}
		}
	}
	return tasks, err
}

func (m *AsyncTaskMux) lockTask(ctx context.Context, taskID string) error {
	lock := m.lockBuilder.NewLock(taskID)
	err := lock.Lock(ctx)
	if err != nil {
		return err
	}
	m.locks[taskID] = lock
	return nil
}

func (m *AsyncTaskMux) unLockTasks(ctx context.Context, taskIDs ...string) error {
	for _, taskID := range taskIDs {
		err := m.unLockTask(ctx, taskID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *AsyncTaskMux) unLockTask(ctx context.Context, taskID string) error {
	lock := m.locks[taskID]
	if lock != nil {
		m.locks[taskID] = nil
		return lock.Unlock(ctx)
	}
	return nil
}

func (m *AsyncTaskMux) loadTaskByID(ctx context.Context, taskIDs ...string) ([]*async_task.AsyncTask, error) {
	res, err := m.dal.FindAll(ctx, taskIDs)
	return res, err
}

func (m *AsyncTaskMux) loadPendingTaskWithLimit(ctx context.Context, limit int) ([]*async_task.AsyncTask, error) {
	list, err := m.dal.FindAllPendingWithLimit(ctx, limit)
	return list, err
}

func (m *AsyncTaskMux) ProcessTask(ctx context.Context, task *async_task.AsyncTask) error {
	if handler, ok := m.handlers[task.TaskName]; ok {
		h := handler
		for i := len(m.mdw) - 1; i >= 0; i-- {
			h = m.mdw[i](h)
		}
		return h(ctx, []byte(task.TaskData))
	}
	return nil
}

func (m *AsyncTaskMux) SendMessage(ctx context.Context, task *async_task.AsyncTask) error {
	msg := &mq.Message{
		Body: []byte(task.TaskData),
	}
	msg.SetBizCode(task.EntityID)
	msg.SetEventName(task.TaskName)
	return m.pub.Publish(ctx, msg)
}

// func (m *AsyncTask) ParseTaskBody(ctx context.Context) ([]byte, error) {
// 	return []byte(m.msg.MsgContent), nil
// }

func (m *AsyncTaskMux) RegisterHandler(key string, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[key] = handler
}
