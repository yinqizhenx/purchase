package repo_impl

import (
	"context"

	"purchase/domain/entity/async_task"
	"purchase/domain/event"
	"purchase/domain/repo"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/tx"
	"purchase/pkg/chanx"
)

type EventRepository struct {
	dal *dal.AsyncTaskDal
	ch  *chanx.UnboundedChan[string]
}

func NewEventRepository(ch *chanx.UnboundedChan[string], dal *dal.AsyncTaskDal) repo.EventRepo {
	return &EventRepository{
		ch:  ch,
		dal: dal,
	}
}

func (r *EventRepository) Save(ctx context.Context, events ...event.Event) error {
	taskList := make([]*async_task.AsyncTask, 0, len(events))
	for _, e := range events {
		task, err := e.ToAsyncTask()
		if err != nil {
			return err
		}
		taskList = append(taskList, task)
	}
	if err := r.dal.BatchAddTask(ctx, taskList...); err != nil {
		return err
	}

	// 在事务提交后，再通知ch
	return tx.RunAfterTxCommit(ctx, func(ctx context.Context) error {
		for _, task := range taskList {
			r.ch.In <- task.TaskID
		}
		return nil
	})
}
