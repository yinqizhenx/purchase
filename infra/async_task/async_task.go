package async_task

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
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

// defaultMaxRetry 任务失败自动重试最大次数
const defaultMaxRetry = 3

// defaultStuckThreshold 任务处于 executing 状态超过此时间视为卡住
const defaultStuckThreshold = 5 * time.Minute

type Handler func(ctx context.Context, payload []byte) error

// deduplicateTTL 控制任务去重的有效期，在此时间内相同 taskID 不会重复分发。
// 设为 cron 间隔的 2 倍，既能有效去重，又不会影响 cron 兜底重试。
const deduplicateTTL = 10 * time.Second

type AsyncTaskMux struct {
	pub              mq.Publisher
	ch               *chanx.UnboundedChan[string]
	cron             *cron.Cron
	dal              *dal.AsyncTaskDal
	txm              *tx.TransactionManager
	concurrency      int
	maxTaskLoad      int
	maxRetry         int           // 任务失败自动重试最大次数
	stuckThreshold   time.Duration // 任务卡住告警阈值
	shutdownTimeout  time.Duration
	cancel           func()
	mdw              []Middleware
	mu               sync.Mutex
	groupWorkers     map[vo.AsyncTaskGroup]*GroupWorker
	listenWg         sync.WaitGroup // 等待 Listen goroutine 退出
	recentDispatched sync.Map       // taskID -> dispatchTime，用于分发去重
}

func NewAsyncTaskMux(pub mq.Publisher, dal *dal.AsyncTaskDal, txm *tx.TransactionManager, ch *chanx.UnboundedChan[string], opts ...Option) *AsyncTaskMux {
	h := &AsyncTaskMux{
		pub:             pub,
		ch:              ch,
		cron:            cron.New(),
		dal:             dal,
		txm:             txm,
		concurrency:     defaultConcurrency, // default 5 worker max
		maxTaskLoad:     defaultMaxTaskLoad,
		maxRetry:        defaultMaxRetry,
		stuckThreshold:  defaultStuckThreshold,
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
		go worker.Start(nctx)
	}
	m.listenWg.Add(1)
	go func() {
		defer m.listenWg.Done()
		m.Listen(nctx)
	}()
	// 启动去重记录定期清理
	go m.cleanupDedup(nctx)
	return nil
}

// cleanupDedup 定期清理过期的去重记录，防止内存泄漏
func (m *AsyncTaskMux) cleanupDedup(ctx context.Context) {
	ticker := time.NewTicker(deduplicateTTL)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			m.recentDispatched.Range(func(key, value any) bool {
				if now.Sub(value.(time.Time)) > deduplicateTTL {
					m.recentDispatched.Delete(key)
				}
				return true
			})
		}
	}
}

// clearDispatched 任务执行完成后清除去重记录，允许后续重新分发（如任务失败需要重试）
func (m *AsyncTaskMux) clearDispatched(taskID string) {
	m.recentDispatched.Delete(taskID)
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
	// 需要同时等待 listenDone 和 workerDone 都完成
	allDone := make(chan struct{})
	go func() {
		// 等待 Listen goroutine 退出
		<-listenDone
		logx.Info(ctx, "listen goroutine stopped")

		// 等待所有 GroupWorker 退出
		<-workerDone
		logx.Info(ctx, "all group workers stopped")

		close(allDone)
	}()

	select {
	case <-shutdownCtx.Done():
		logx.Warn(ctx, "shutdown timeout exceeded, forcing stop",
			slog.Duration("timeout", m.shutdownTimeout))
		// 超时后强制停止
		for _, gw := range m.groupWorkers {
			gw.ForceStop()
		}
		return shutdownCtx.Err()
	case <-allDone:
		logx.Info(ctx, "all components stopped successfully")
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
				// ctx 取消后不再查询数据库
				if ctx.Err() != nil {
					return
				}
				// 先检查去重，如果该 taskID 近期已分发过，跳过 DB 查询
				if v, loaded := m.recentDispatched.Load(taskID); loaded {
					if time.Since(v.(time.Time)) < deduplicateTTL {
						return
					}
				}
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
	// 每 60 秒检测一次卡住的任务
	_, err = m.cron.AddFunc("@every 60s", m.buildStuckTaskChecker(ctx))
	if err != nil {
		return err
	}
	go m.cron.Run()
	return nil
}

func (m *AsyncTaskMux) buildCronHandler(ctx context.Context) func() {
	return func() {
		// 随机抖动 0~2s，错开多实例的查询时间，减少同时查同一批数据
		jitter := time.Duration(rand.Int63n(int64(2 * time.Second)))
		time.Sleep(jitter)

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
	// 去重：短时间内相同 taskID 不重复分发
	if v, loaded := m.recentDispatched.LoadOrStore(task.TaskID, time.Now()); loaded {
		dispatchTime := v.(time.Time)
		if time.Since(dispatchTime) < deduplicateTTL {
			return
		}
		// 已过期，更新时间戳并继续分发
		m.recentDispatched.Store(task.TaskID, time.Now())
	}

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

// buildStuckTaskChecker 构建卡住任务检测定时函数
// 定期查询长时间处于 executing 状态的任务，输出告警日志
func (m *AsyncTaskMux) buildStuckTaskChecker(ctx context.Context) func() {
	return func() {
		stuckTasks, err := m.dal.FindStuckExecutingTasks(ctx, m.stuckThreshold, m.maxTaskLoad)
		if err != nil {
			logx.Error(ctx, "check stuck tasks fail", slog.Any("error", err))
			return
		}
		for _, task := range stuckTasks {
			logx.Error(ctx, "task stuck in executing state for too long",
				slog.String("task_id", task.TaskID),
				slog.String("task_name", task.TaskName),
				slog.String("task_group", string(task.TaskGroup)),
				slog.String("entity_id", task.EntityID),
				slog.Int("retry_count", task.RetryCount),
				slog.Time("updated_at", task.UpdatedAt),
				slog.Duration("stuck_duration", time.Since(task.UpdatedAt)),
			)
		}
	}
}

// RetryTask 手动重试指定的异步任务
// 将任务状态重置为 pending，等待下次 cron 调度执行
// 仅允许重试 fail 或 executing 状态的任务
func (m *AsyncTaskMux) RetryTask(ctx context.Context, taskID string) error {
	task, err := m.dal.FindOneNoNil(ctx, taskID)
	if err != nil {
		return fmt.Errorf("find task fail: %w", err)
	}

	switch task.State {
	case vo.AsyncTaskStateFail:
		n, err := m.dal.ResetTaskToPendingWithRetry(ctx, taskID, vo.AsyncTaskStateFail)
		if err != nil {
			return fmt.Errorf("reset failed task to pending fail: %w", err)
		}
		if n == 0 {
			return fmt.Errorf("task %s state has changed, retry aborted", taskID)
		}
		logx.Info(ctx, "manually retried failed task",
			slog.String("task_id", taskID),
			slog.Int("retry_count", task.RetryCount+1))
		// 清除去重记录，允许立即重新分发
		m.clearDispatched(taskID)
		return nil

	case vo.AsyncTaskStateExecuting:
		n, err := m.dal.ResetTaskToPendingWithRetry(ctx, taskID, vo.AsyncTaskStateExecuting)
		if err != nil {
			return fmt.Errorf("reset executing task to pending fail: %w", err)
		}
		if n == 0 {
			return fmt.Errorf("task %s state has changed, retry aborted", taskID)
		}
		logx.Info(ctx, "manually retried stuck executing task",
			slog.String("task_id", taskID),
			slog.Int("retry_count", task.RetryCount+1))
		m.clearDispatched(taskID)
		return nil

	case vo.AsyncTaskStatePending:
		return fmt.Errorf("task %s is already in pending state", taskID)

	case vo.AsyncTaskStateSuccess:
		return fmt.Errorf("task %s has already succeeded", taskID)

	default:
		return fmt.Errorf("task %s is in unknown state: %s", taskID, task.State)
	}
}

func (m *AsyncTaskMux) RegisterHandler(key string, group vo.AsyncTaskGroup, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.groupWorkers[group] == nil {
		gw := NewGroupWorker(m.pub, m.dal, m.txm, m.concurrency, m.mdw, m.maxRetry)
		gw.SetGroup(group)
		gw.onTaskDone = m.clearDispatched
		m.groupWorkers[group] = gw
	}
	m.groupWorkers[group].RegisterHandler(key, handler)
}
