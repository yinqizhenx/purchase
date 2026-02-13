package async_task

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"purchase/domain/entity/async_task"
	"purchase/domain/vo"
	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/tx"
	"purchase/pkg/chanx"
)

type GroupWorker struct {
	pub mq.Publisher
	// handles 涉及到数据库变动的，需要在事务里执行
	handlers    map[string]Handler
	ch          *chanx.UnboundedChan[*async_task.AsyncTask]
	dal         *dal.AsyncTaskDal
	txm         *tx.TransactionManager
	sem         chan struct{}
	concurrency int
	maxRetry    int // 任务失败自动重试最大次数
	cancel      func()
	cancelChan  func() // 关闭 unbounded channel
	mdw         []Middleware
	mu          sync.Mutex
	taskGroup   vo.AsyncTaskGroup   // 任务组标识
	onTaskDone  func(taskID string) // 任务执行完成后的回调（用于清理去重记录）

	// 优雅关闭相关
	runningTasks sync.WaitGroup // 跟踪正在执行的任务
	stopping     int32          // 原子标记是否正在关闭
	listenWg     sync.WaitGroup // 等待 Listen goroutine 退出
}

func NewGroupWorker(pub mq.Publisher, dal *dal.AsyncTaskDal, txm *tx.TransactionManager, concurrency int, mdw []Middleware, maxRetry int) *GroupWorker {
	ch, cancelChan := NewTaskChan()
	w := &GroupWorker{
		pub:         pub,
		handlers:    make(map[string]Handler),
		ch:          ch,
		dal:         dal,
		txm:         txm,
		concurrency: concurrency,
		maxRetry:    maxRetry,
		cancelChan:  cancelChan,
		mdw:         mdw,
	}
	w.sem = make(chan struct{}, w.concurrency)
	return w
}

// RegisterHandler 注册任务处理器到当前 GroupWorker
func (w *GroupWorker) RegisterHandler(key string, handler Handler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.handlers[key] = handler
}

// SetGroup 设置任务组标识
func (w *GroupWorker) SetGroup(group vo.AsyncTaskGroup) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.taskGroup = group
}

// GetGroup 获取任务组标识
func (w *GroupWorker) GetGroup() vo.AsyncTaskGroup {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.taskGroup
}

func (w *GroupWorker) Start(ctx context.Context) {
	nctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	w.listenWg.Add(1)
	go func() {
		defer w.listenWg.Done()
		w.Listen(nctx)
	}()
}

// Stop 优雅关闭 GroupWorker
// 1. 停止接收新任务
// 2. 等待正在执行的任务完成（带超时）
// 3. 关闭资源
func (w *GroupWorker) Stop(ctx context.Context) error {
	// 标记正在关闭
	if !atomic.CompareAndSwapInt32(&w.stopping, 0, 1) {
		// 已经在关闭中
		return nil
	}

	logx.Info(ctx, "starting graceful shutdown of group worker",
		slog.String("group", string(w.GetGroup())))

	// 1. 停止接收新任务
	if w.cancel != nil {
		w.cancel()
	}

	// 等待 Listen goroutine 退出
	listenDone := make(chan struct{})
	go func() {
		w.listenWg.Wait()
		close(listenDone)
	}()

	// 2. 等待正在执行的任务完成
	taskDone := make(chan struct{})
	go func() {
		w.runningTasks.Wait()
		close(taskDone)
	}()

	// 等待所有组件退出，或超时
	// 需要同时等待 listenDone 和 taskDone 都完成
	allDone := make(chan struct{})
	go func() {
		// 等待 Listen goroutine 退出
		<-listenDone
		logx.Info(ctx, "listen goroutine stopped",
			slog.String("group", string(w.GetGroup())))

		// 等待所有正在执行的任务完成
		<-taskDone
		logx.Info(ctx, "all running tasks completed",
			slog.String("group", string(w.GetGroup())))

		close(allDone)
	}()

	select {
	case <-ctx.Done():
		logx.Warn(ctx, "shutdown timeout, some tasks may still be running",
			slog.String("group", string(w.GetGroup())))
		return ctx.Err()
	case <-allDone:
		logx.Info(ctx, "all components stopped successfully",
			slog.String("group", string(w.GetGroup())))
	}

	// 3. 关闭 channel
	if w.cancelChan != nil {
		w.cancelChan()
	}

	logx.Info(ctx, "group worker stopped gracefully",
		slog.String("group", string(w.GetGroup())))
	return nil
}

// ForceStop 强制停止，不等待任务完成
func (w *GroupWorker) ForceStop() {
	atomic.StoreInt32(&w.stopping, 1)
	if w.cancel != nil {
		w.cancel()
	}
	if w.cancelChan != nil {
		w.cancelChan()
	}
}

func (w *GroupWorker) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logx.Info(ctx, "listen context cancelled, stopping listener",
				slog.String("group", string(w.GetGroup())))
			return
		case task, ok := <-w.ch.Out:
			if !ok {
				logx.Info(ctx, "task channel closed, stopping listener",
					slog.String("group", string(w.GetGroup())))
				return
			}

			// 检查是否正在关闭
			if atomic.LoadInt32(&w.stopping) == 1 {
				logx.Info(ctx, "worker is stopping, rejecting new task",
					slog.String("group", string(w.GetGroup())),
					slog.String("task", task.TaskID))
				// 将任务重新放回队列或标记为 pending
				// 这里可以选择将任务状态重置为 pending，让下次重新处理
				continue
			}

			select {
			case <-ctx.Done():
				return
			case w.sem <- struct{}{}:
				w.runningTasks.Add(1)
				go func() {
					defer func() {
						<-w.sem
						w.runningTasks.Done()
					}()
					ok, err := w.tryLockTask(ctx, task)
					if err != nil || !ok {
						return
					}
					w.handleTask(ctx, task)
				}()
			}
		}
	}
}

func (w *GroupWorker) Handle(ctx context.Context, msg *async_task.AsyncTask) error {
	if msg.IsTask(ctx) {
		return w.ProcessTask(ctx, msg)
	}
	if msg.IsEvent(ctx) {
		return w.SendMessage(ctx, msg)
	}
	return errors.New("unknown task type")
}

// onHandleSuccess 这个里面是业务逻辑处理成功后的其他逻辑，一般不会出错，所以要求一定成功，没有抛错误出去，
func (w *GroupWorker) onHandleSuccess(ctx context.Context, task *async_task.AsyncTask) error {
	err := w.dal.UpdateExecutingTaskSuccess(ctx, task.TaskID)
	if err != nil {
		logx.Error(ctx, "update task state success fail", slog.String("task", task.TaskID), slog.Any("error", err))
	}
	return err
}

// onHandleFailWithRetry 业务 handler 失败后的统一处理（在事务外执行）
// 如果未达到最大重试次数，直接 executing → pending（retry_count+1），省去 fail 中间态
// 如果已达到最大重试次数，executing → fail，标记为最终失败
func (w *GroupWorker) onHandleFailWithRetry(ctx context.Context, task *async_task.AsyncTask) {
	if task.RetryCount+1 < w.maxRetry {
		// 未达到重试上限，直接从 executing 重置为 pending，一次 DB 操作完成
		n, err := w.dal.ResetTaskToPendingWithRetry(ctx, task.TaskID, vo.AsyncTaskStateExecuting)
		if err != nil {
			logx.Error(ctx, "auto retry: reset executing task to pending fail",
				slog.String("task_id", task.TaskID), slog.Any("error", err))
			return
		}
		if n > 0 {
			logx.Warn(ctx, "task failed, auto retrying",
				slog.String("task_id", task.TaskID),
				slog.String("task_name", task.TaskName),
				slog.Int("current_retry", task.RetryCount+1),
				slog.Int("max_retry", w.maxRetry))
		}
		return
	}

	// 达到重试上限，标记为最终失败
	logx.Error(ctx, "task failed, max retry reached",
		slog.String("task_id", task.TaskID),
		slog.String("task_name", task.TaskName),
		slog.Int("retry_count", task.RetryCount),
		slog.Int("max_retry", w.maxRetry))
	if err := w.dal.UpdateExecutingTaskFail(ctx, task.TaskID); err != nil {
		logx.Error(ctx, "mark task as final fail error, task will be stuck until stuck checker",
			slog.String("task_id", task.TaskID), slog.Any("error", err))
	}
}

func (w *GroupWorker) handleTask(ctx context.Context, task *async_task.AsyncTask) {
	defer func() {
		// 任务执行完成后清理去重记录
		if w.onTaskDone != nil {
			w.onTaskDone(task.TaskID)
		}
	}()

	fn := func(txCtx context.Context) error {
		if err := w.Handle(txCtx, task); err != nil {
			logx.Error(txCtx, "handle task fail", slog.String("task", task.TaskID), slog.Any("error", err))
			return err
		}
		return w.onHandleSuccess(txCtx, task)
	}
	if err := w.txm.Transaction(ctx, fn); err != nil {
		// 无论是业务 handler 失败还是 onHandleSuccess/事务提交失败，事务都已回滚
		// 业务数据未提交，任务仍为 executing，统一走失败重试逻辑
		logx.Error(ctx, "handle task transaction fail", slog.String("task", task.TaskID), slog.Any("error", err))
		w.onHandleFailWithRetry(ctx, task)
	}
}

// tryLockTask 先锁再查
func (w *GroupWorker) tryLockTask(ctx context.Context, task *async_task.AsyncTask) (bool, error) {
	n, err := w.dal.UpdatePendingTaskExecuting(ctx, task.TaskID)
	if err != nil {
		logx.Errorf(ctx, "UpdatePendingTaskExecuting fail, task:%s, err:%v", task.TaskID, err)
		return false, err
	}
	if n == 0 {
		return false, nil
	}
	return true, nil
}

func (w *GroupWorker) ProcessTask(ctx context.Context, task *async_task.AsyncTask) error {
	handler, ok := w.handlers[task.TaskName]
	if !ok {
		return fmt.Errorf("no handler registered for task: %s", task.TaskName)
	}
	h := handler
	for i := len(w.mdw) - 1; i >= 0; i-- {
		h = w.mdw[i](h)
	}
	return h(ctx, []byte(task.TaskData))
}

func (w *GroupWorker) SendMessage(ctx context.Context, task *async_task.AsyncTask) error {
	msg := &mq.Message{
		Body: []byte(task.TaskData),
	}
	msg.SetBizCode(task.EntityID)
	msg.SetEventName(task.TaskName)
	return w.pub.Publish(ctx, msg)
}

func (w *GroupWorker) PutTask(task *async_task.AsyncTask) {
	w.ch.In <- task
}
