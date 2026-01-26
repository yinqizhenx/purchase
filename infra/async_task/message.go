package async_task

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

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

const defaultShutdownTimeout = 30 * time.Second

type Handler func(ctx context.Context, payload []byte) error

type AsyncTaskMux struct {
	pub mq.Publisher
	// handles 涉及到数据库变动的，需要在事务里执行
	handlers        map[string]Handler
	ch              *chanx.UnboundedChan[string]
	cron            *cron.Cron
	dal             *dal.AsyncTaskDal
	txm             *tx.TransactionManager
	concurrency     int
	maxTaskLoad     int
	shutdownTimeout time.Duration
	cancel          func()
	mdw             []Middleware
	mu              sync.Mutex
	groupWorkers    map[vo.AsyncTaskGroup]*GroupWorker
	listenWg        sync.WaitGroup // 等待 Listen goroutine 退出
	stopping        int32          // 使用 atomic 标记是否正在关闭
}

func NewAsyncTaskMux(pub mq.Publisher, dal *dal.AsyncTaskDal, txm *tx.TransactionManager, ch *chanx.UnboundedChan[string], opts ...Option) *AsyncTaskMux {
	h := &AsyncTaskMux{
		pub:             pub,
		handlers:        make(map[string]Handler),
		ch:              ch,
		cron:            cron.New(),
		dal:             dal,
		txm:             txm,
		concurrency:     defaultConcurrency, // default 5 worker max
		maxTaskLoad:     defaultMaxTaskLoad,
		shutdownTimeout: defaultShutdownTimeout,
		groupWorkers:    make(map[vo.AsyncTaskGroup]*GroupWorker),
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
	m.listenWg.Add(1)
	go func() {
		defer m.listenWg.Done()
		m.Listen(nctx)
	}()
	return nil
}

// Stop 优雅关闭异步任务处理器
// 1. 停止接收新任务（取消 context，停止 cron）
// 2. 等待正在执行的任务完成（带超时）
// 3. 关闭所有资源
func (m *AsyncTaskMux) Stop(ctx context.Context) error {
	logx.Info(ctx, "starting graceful shutdown of async task mux")

	// 1. 停止接收新任务
	if m.cancel != nil {
		m.cancel()
	}

	// 停止 cron 定时任务
	m.cron.Stop()
	logx.Info(ctx, "stopped cron scheduler")

	// 2. 等待 Listen goroutine 退出
	listenDone := make(chan struct{})
	go func() {
		m.listenWg.Wait()
		close(listenDone)
	}()

	// 3. 等待所有 GroupWorker 优雅退出
	shutdownCtx, cancel := context.WithTimeout(ctx, m.shutdownTimeout)
	defer cancel()

	workerDone := make(chan struct{})
	go func() {
		for _, gw := range m.groupWorkers {
			if err := gw.Stop(shutdownCtx); err != nil {
				logx.Error(shutdownCtx, "stop group worker failed",
					slog.String("group", string(gw.GetGroup())),
					slog.Any("error", err))
			}
		}
		close(workerDone)
	}()

	// 等待所有组件退出，或超时
	select {
	case <-shutdownCtx.Done():
		logx.Warn(ctx, "shutdown timeout exceeded, forcing stop",
			slog.Duration("timeout", m.shutdownTimeout))
		// 超时后强制停止
		for _, gw := range m.groupWorkers {
			gw.ForceStop()
		}
		return shutdownCtx.Err()
	case <-listenDone:
		logx.Info(ctx, "listen goroutine stopped")
	case <-workerDone:
		logx.Info(ctx, "all group workers stopped")
	}

	// 确保所有 goroutine 都已退出
	select {
	case <-shutdownCtx.Done():
		logx.Warn(ctx, "final wait timeout exceeded")
		return shutdownCtx.Err()
	case <-listenDone:
		// 已退出
	}

	logx.Info(ctx, "async task mux stopped gracefully")
	return nil
}

func (m *AsyncTaskMux) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logx.Info(ctx, "listen context cancelled, stopping listener")
			return
		case taskID, ok := <-m.ch.Out:
			if !ok {
				logx.Info(ctx, "task channel closed, stopping listener")
				return
			}
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
		gw := NewGroupWorker(m.pub, m.dal, m.txm)
		gw.SetGroup(group)
		m.groupWorkers[group] = gw
	}
}
