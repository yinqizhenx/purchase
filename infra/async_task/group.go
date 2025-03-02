package async_task

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"purchase/domain/entity/async_task"
	"purchase/infra/logx"
	"purchase/infra/mq"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/tx"
	"purchase/infra/utils"
	"purchase/pkg/chanx"
)

type GroupWorker struct {
	// msg          *message.Task
	pub mq.Publisher
	// handles 涉及到数据库变动的，需要在事务里执行
	handlers map[string]Handler
	ch       *chanx.UnboundedChan[*async_task.AsyncTask]
	// cron        *cron.Cron
	dal         *dal.AsyncTaskDal
	txm         *tx.TransactionManager
	sem         chan struct{}
	concurrency int
	// maxTaskLoad int
	// lockBuilder dlock.LockBuilder
	// locks       map[string]dlock.DistributeLock
	cancel  func()
	cancel2 func()
	mdw     []Middleware
	mu      sync.Mutex
}

func NewGroupWorker(pub mq.Publisher, dal *dal.AsyncTaskDal, txm *tx.TransactionManager) *GroupWorker {
	ch, cancel := NewTaskChan()
	h := &GroupWorker{
		pub:      pub,
		handlers: make(map[string]Handler),
		ch:       ch,
		// cron:        cron.New(),
		dal:         dal,
		txm:         txm,
		concurrency: defaultConcurrency, // default 5 worker max
		cancel2:     cancel,
		// maxTaskLoad: defaultMaxTaskLoad,
	}
	h.sem = make(chan struct{}, h.concurrency)
	return h
}

func (m *GroupWorker) Start(ctx context.Context) {
	nctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	// err := m.RunCron(nctx)
	// if err != nil {
	// 	logx.Error(ctx, "launch task cron handle fail", slog.Any("error", err))
	// 	return err
	// }
	m.Listen(nctx)
}

func (m *GroupWorker) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
}

func (m *GroupWorker) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-m.ch.Out:
			select {
			case <-ctx.Done():
				return
			case m.sem <- struct{}{}:
				go func() {
					defer func() {
						<-m.sem
					}()
					ok, err := m.tryLockTask(ctx, task)
					if err != nil || !ok {
						return
					}
					m.handleTask(ctx, task)
				}()
			}
		}
	}
}

func (m *GroupWorker) Handle(ctx context.Context, msg *async_task.AsyncTask) error {
	if msg.IsTask(ctx) {
		return m.ProcessTask(ctx, msg)
	}
	if msg.IsEvent(ctx) {
		return m.SendMessage(ctx, msg)
	}
	return errors.New("unknown task type")
}

// onHandleSuccess 这个里面是业务逻辑处理成功后的其他逻辑，一般不会出错，所以要求一定成功，没有抛错误出去，
func (m *GroupWorker) onHandleSuccess(ctx context.Context, task *async_task.AsyncTask) error {
	err := m.dal.UpdateExecutingTaskSuccess(ctx, task.TaskID)
	if err != nil {
		logx.Error(ctx, "update task state success fail", slog.String("task", task.TaskID), slog.Any("error", err))
	}
	return err
}

func (m *GroupWorker) onHandleFail(ctx context.Context, task *async_task.AsyncTask) error {
	err := m.dal.UpdateExecutingTaskFail(ctx, task.TaskID)
	if err != nil {
		logx.Error(ctx, "update task state fail fail", slog.String("task", task.TaskID), slog.Any("error", err))
	}
	return err
}

func (m *GroupWorker) handleTask(ctx context.Context, taskList ...*async_task.AsyncTask) {
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

// tryLockTask 先锁再查
func (m *GroupWorker) tryLockTask(ctx context.Context, task *async_task.AsyncTask) (bool, error) {
	n, err := m.dal.UpdatePendingTaskExecuting(ctx, task.TaskID)
	if err != nil {
		logx.Errorf(ctx, "UpdatePendingTaskExecuting fail, task:%s, err:%v", task.TaskID, err)
		return false, err
	}
	if n == 0 {
		return false, nil
	}
	return true, nil
}

func (m *GroupWorker) ProcessTask(ctx context.Context, task *async_task.AsyncTask) error {
	if handler, ok := m.handlers[task.TaskName]; ok {
		h := handler
		for i := len(m.mdw) - 1; i >= 0; i-- {
			h = m.mdw[i](h)
		}
		return h(ctx, []byte(task.TaskData))
	}
	return nil
}

func (m *GroupWorker) SendMessage(ctx context.Context, task *async_task.AsyncTask) error {
	msg := &mq.Message{
		Body: []byte(task.TaskData),
	}
	msg.SetBizCode(task.EntityID)
	msg.SetEventName(task.TaskName)
	return m.pub.Publish(ctx, msg)
}

func (m *GroupWorker) PutTask(task *async_task.AsyncTask) {
	m.ch.In <- task
}

// func (m *GroupWorker) RegisterHandler(key string, handler Handler) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.handlers[key] = handler
// }
