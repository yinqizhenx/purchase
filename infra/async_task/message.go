package async_task

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"

	"purchase/domain/entity/async_task"
	"purchase/domain/vo"
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
					task, err := m.loadAndLockPendingTask(ctx, taskID)
					if err != nil {
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

					m.handleTask(ctx, task)
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

// onHandleSuccess 这个里面是业务逻辑处理成功后的其他逻辑，一般不会出错，所以要求一定成功，没有抛错误出去，
func (m *AsyncTaskMux) onHandleSuccess(ctx context.Context, task *async_task.AsyncTask) error {
	err := m.dal.UpdateTaskSuccess(ctx, task.TaskID)
	if err != nil {
		logx.Error(ctx, "update task state success fail", slog.String("task", task.TaskID), slog.Any("error", err))
	}
	return err
}

func (m *AsyncTaskMux) onHandleFail(ctx context.Context, task *async_task.AsyncTask) error {
	err := m.dal.UpdateTaskFail(ctx, task.TaskID)
	if err != nil {
		logx.Error(ctx, "update task state fail fail", slog.String("task", task.TaskID), slog.Any("error", err))
	}
	return err
}

func (m *AsyncTaskMux) buildCronHandler(ctx context.Context) func() {
	return func() {
		taskList, err := m.loadAndLockTaskWithLimit(ctx, m.maxTaskLoad)
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
		m.handleTask(ctx, taskList...)
	}
}

func (m *AsyncTaskMux) handleTask(ctx context.Context, taskList ...*async_task.AsyncTask) {
	wg := &sync.WaitGroup{}
	for _, task := range taskList {
		wg.Add(1)
		t := task
		fn := func(context.Context) (err error) {
			defer func() {
				wg.Done()
			}()
			err = m.Handle(ctx, t)
			if err != nil {
				logx.Error(ctx, "handle task fail", slog.String("task", t.TaskID), slog.Any("error", err))
				return m.onHandleFail(ctx, t)
			}
			return m.onHandleSuccess(ctx, t)
		}
		utils.SafeGo(ctx, func() {
			err := m.txm.Transaction(ctx, fn)
			if err != nil {
				logx.Error(ctx, "handle task fail", slog.String("task", t.TaskID), slog.Any("error", err))
			}
		})
	}
	wg.Wait()
}

func (m *AsyncTaskMux) loadAndLockPendingTask(ctx context.Context, taskID string) (*async_task.AsyncTask, error) {
	task, err := m.dal.FindOneNoNil(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task.State != vo.AsyncTaskStatePending {
		return nil, errors.New("task state is not pending")
	}
	err = m.lockTaskAndUpdateState(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *AsyncTaskMux) loadAndLockTaskWithLimit(ctx context.Context, limit int) ([]*async_task.AsyncTask, error) {
	taskList, err := m.loadPendingTaskWithLimit(ctx, limit)
	if err != nil {
		return nil, err
	}
	err = m.lockTaskAndUpdateState(ctx, taskList...)
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

// lockAndLoadTask 先锁再查
func (m *AsyncTaskMux) lockTaskAndUpdateState(ctx context.Context, tasks ...*async_task.AsyncTask) error {
	lockedIDs := make([]string, 0)
	for _, task := range tasks {
		err := m.lockTask(ctx, task.TaskID)
		if err != nil {
			logx.Infof(ctx, "lock task(%s) fail: %v", task.TaskID, err)
			continue
		}
		lockedIDs = append(lockedIDs, task.TaskID)
	}
	if len(lockedIDs) == 0 {
		return nil
	}
	err := m.dal.UpdateTaskExecuting(ctx, lockedIDs...)
	if err != nil {
		if err2 := m.unLockTasks(ctx, lockedIDs...); err2 != nil {
			log.Errorf("un lock task fail, err:%v", err2)
		}
	}
	return err
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
