package async_task

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"

	"purchase/domain/entity/async_task"
	"purchase/domain/vo"
	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/tx"
	"purchase/infra/utils"
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
	cancel      func()
	cancel2     func()
	mdw         []Middleware
	mu          sync.Mutex
	taskGroup   vo.AsyncTaskGroup // 任务组标识

	// 优雅关闭相关
	runningTasks sync.WaitGroup // 跟踪正在执行的任务
	stopping     int32          // 原子标记是否正在关闭
	listenWg     sync.WaitGroup // 等待 Listen goroutine 退出
}

func NewGroupWorker(pub mq.Publisher, dal *dal.AsyncTaskDal, txm *tx.TransactionManager) *GroupWorker {
	ch, cancel := NewTaskChan()
	w := &GroupWorker{
		pub:         pub,
		handlers:    make(map[string]Handler),
		ch:          ch,
		dal:         dal,
		txm:         txm,
		concurrency: defaultConcurrency, // default 5 worker max
		cancel2:     cancel,
	}
	w.sem = make(chan struct{}, w.concurrency)
	return w
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

	// 等待任务完成或超时
	select {
	case <-ctx.Done():
		logx.Warn(ctx, "shutdown timeout, some tasks may still be running",
			slog.String("group", string(w.GetGroup())))
		return ctx.Err()
	case <-listenDone:
		logx.Info(ctx, "listen goroutine stopped",
			slog.String("group", string(w.GetGroup())))
	case <-taskDone:
		logx.Info(ctx, "all running tasks completed",
			slog.String("group", string(w.GetGroup())))
	}

	// 确保 Listen goroutine 已退出
	select {
	case <-ctx.Done():
		logx.Warn(ctx, "final wait timeout",
			slog.String("group", string(w.GetGroup())))
		return ctx.Err()
	case <-listenDone:
		// 已退出
	}

	// 3. 关闭 channel
	if w.cancel2 != nil {
		w.cancel2()
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
	if w.cancel2 != nil {
		w.cancel2()
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

func (w *GroupWorker) onHandleFail(ctx context.Context, task *async_task.AsyncTask) error {
	err := w.dal.UpdateExecutingTaskFail(ctx, task.TaskID)
	if err != nil {
		logx.Error(ctx, "update task state fail fail", slog.String("task", task.TaskID), slog.Any("error", err))
	}
	return err
}

func (w *GroupWorker) handleTask(ctx context.Context, taskList ...*async_task.AsyncTask) {
	wg := &sync.WaitGroup{}
	for _, task := range taskList {
		wg.Add(1)
		t := task
		fn := func(context.Context) (err error) {
			defer func() {
				wg.Done()
			}()
			err = w.Handle(ctx, t)
			if err != nil {
				logx.Error(ctx, "handle task fail", slog.String("task", t.TaskID), slog.Any("error", err))
				return w.onHandleFail(ctx, t)
			}
			return w.onHandleSuccess(ctx, t)
		}
		utils.SafeGo(ctx, func() {
			err := w.txm.Transaction(ctx, fn)
			if err != nil {
				logx.Error(ctx, "handle task fail", slog.String("task", t.TaskID), slog.Any("error", err))
			}
		})
	}
	wg.Wait()
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
	if handler, ok := w.handlers[task.TaskName]; ok {
		h := handler
		for i := len(w.mdw) - 1; i >= 0; i-- {
			h = w.mdw[i](h)
		}
		return h(ctx, []byte(task.TaskData))
	}
	return nil
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
