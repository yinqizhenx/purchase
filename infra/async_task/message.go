package async_task

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/robfig/cron/v3"

	"purchase/domain/entity/async_task"
	"purchase/domain/vo"
	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/tx"
	"purchase/pkg/chanx"
)

const defaultMaxTaskLoad = 100

const defaultConcurrency = 5

type Handler func(ctx context.Context, payload []byte) error

type AsyncTaskMux struct {
	pub mq.Publisher
	// handles 涉及到数据库变动的，需要在事务里执行
	handlers     map[string]Handler
	ch           *chanx.UnboundedChan[string]
	cron         *cron.Cron
	dal          *dal.AsyncTaskDal
	txm          *tx.TransactionManager
	concurrency  int
	maxTaskLoad  int
	cancel       func()
	mdw          []Middleware
	mu           sync.Mutex
	groupWorkers map[vo.AsyncTaskGroup]*GroupWorker
}

func NewAsyncTaskMux(pub mq.Publisher, dal *dal.AsyncTaskDal, txm *tx.TransactionManager, ch *chanx.UnboundedChan[string], opts ...Option) *AsyncTaskMux {
	h := &AsyncTaskMux{
		pub:          pub,
		handlers:     make(map[string]Handler),
		ch:           ch,
		cron:         cron.New(),
		dal:          dal,
		txm:          txm,
		concurrency:  defaultConcurrency, // default 5 worker max
		maxTaskLoad:  defaultMaxTaskLoad,
		groupWorkers: make(map[vo.AsyncTaskGroup]*GroupWorker),
	}
	for _, opt := range opts {
		opt(h)
	}
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
	for _, gw := range m.groupWorkers {
		worker := gw
		go worker.Start(ctx)
	}
	m.Listen(nctx)
	return nil
}

func (m *AsyncTaskMux) Stop(_ context.Context) error {
	if m.cancel != nil {
		m.cancel()
	}
	m.cron.Stop()
	for _, gw := range m.groupWorkers {
		gw.Stop()
	}
	return nil
}

func (m *AsyncTaskMux) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case taskID := <-m.ch.Out:
			go func() {
				task, err := m.mustGetPendingTask(ctx, taskID)
				if err != nil {
					logx.Errorf(ctx, "find one task fail:%s, err :%s", taskID, err)
					return
				}

				m.distribute(ctx, task)
			}()
		}
	}
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

func (m *AsyncTaskMux) buildCronHandler(ctx context.Context) func() {
	return func() {
		taskList, err := m.getPendingTaskWithLimit(ctx, m.maxTaskLoad)
		if err != nil {
			logx.Error(ctx, "batch load task fail", slog.Any("error", err))
			return
		}
		for _, task := range taskList {
			m.distribute(ctx, task)
		}
	}
}

func (m *AsyncTaskMux) distribute(ctx context.Context, task *async_task.AsyncTask) {
	gw := m.groupWorkers[task.TaskGroup]
	if gw == nil {
		logx.Errorf(ctx, "unknown task group: %s", task.TaskGroup)
		return
	}
	gw.PutTask(task)
}

func (m *AsyncTaskMux) getPendingTaskWithLimit(ctx context.Context, limit int) ([]*async_task.AsyncTask, error) {
	list, err := m.dal.FindAllPendingWithLimit(ctx, limit)
	return list, err
}

func (m *AsyncTaskMux) mustGetPendingTask(ctx context.Context, taskID string) (*async_task.AsyncTask, error) {
	task, err := m.dal.FindOneNoNil(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task.State != vo.AsyncTaskStatePending {
		return nil, errors.New("task is not pending state")
	}
	return task, nil
}

func (m *AsyncTaskMux) RegisterHandler(key string, group vo.AsyncTaskGroup, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[key] = handler
	if m.groupWorkers[group] == nil {
		m.groupWorkers[group] = NewGroupWorker(m.pub, m.dal, m.txm)
	}
}
