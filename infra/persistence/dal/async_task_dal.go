package dal

import (
	"context"
	"time"

	eTask "purchase/infra/persistence/dal/db/ent/asynctask"

	"purchase/domain/entity/async_task"
	"purchase/domain/vo"
	"purchase/infra/persistence/convertor"
	"purchase/infra/persistence/dal/db/ent"
	"purchase/infra/persistence/tx"
)

type AsyncTaskDal struct {
	db        *ent.Client
	convertor *convertor.Convertor
}

func NewAsyncTaskDal(cli *ent.Client, c *convertor.Convertor) *AsyncTaskDal {
	return &AsyncTaskDal{
		db:        cli,
		convertor: c,
	}
}

// getClient 确保context在事务中时，数据库操作在事务中执行
func (dal *AsyncTaskDal) getClient(ctx context.Context) *ent.AsyncTaskClient {
	txCtx, ok := ctx.(*tx.TransactionContext)
	if ok {
		return txCtx.Tx().AsyncTask
	}
	return dal.db.AsyncTask
}

func (dal *AsyncTaskDal) AddTask(ctx context.Context, task *async_task.AsyncTask) error {
	err := dal.getClient(ctx).Create().
		SetTaskID(task.TaskID).
		SetTaskType(string(task.TaskType)).
		SetTaskGroup(string(task.TaskGroup)).
		SetTaskData(task.TaskData).
		Exec(ctx)
	return err
}

func (dal *AsyncTaskDal) BatchAddTask(ctx context.Context, taskList ...*async_task.AsyncTask) error {
	taskAddList := make([]*ent.AsyncTaskCreate, 0, len(taskList))
	for _, task := range taskList {
		c := dal.getClient(ctx).Create().
			SetTaskID(task.TaskID).
			SetTaskName(task.TaskName).
			SetState(string(vo.AsyncTaskStatePending)).
			SetBizID(task.EntityID).
			SetTaskType(string(task.TaskType)).
			SetTaskData(task.TaskData)
		taskAddList = append(taskAddList, c)
	}
	return dal.db.AsyncTask.CreateBulk(taskAddList...).Exec(ctx)
}

func (dal *AsyncTaskDal) FindOneNoNil(ctx context.Context, taskID string) (*async_task.AsyncTask, error) {
	res, err := dal.getClient(ctx).Query().Where(eTask.TaskID(taskID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return dal.convertor.ConvertAsyncTaskPoToDo(res), nil
}

func (dal *AsyncTaskDal) FindAll(ctx context.Context, taskIDList []string) ([]*async_task.AsyncTask, error) {
	if len(taskIDList) == 0 {
		return nil, nil
	}
	res, err := dal.getClient(ctx).Query().Where(eTask.TaskIDIn(taskIDList...)).All(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*async_task.AsyncTask, 0, len(res))
	for _, r := range res {
		list = append(list, dal.convertor.ConvertAsyncTaskPoToDo(r))
	}
	return list, nil
}

func (dal *AsyncTaskDal) FindAllPending(ctx context.Context, taskIDList []string) ([]*async_task.AsyncTask, error) {
	if len(taskIDList) == 0 {
		return nil, nil
	}
	res, err := dal.getClient(ctx).Query().
		Where(eTask.TaskIDIn(taskIDList...)).
		Where(eTask.State(string(vo.AsyncTaskStatePending))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*async_task.AsyncTask, 0, len(res))
	for _, r := range res {
		list = append(list, dal.convertor.ConvertAsyncTaskPoToDo(r))
	}
	return list, nil
}

func (dal *AsyncTaskDal) FindAllPendingWithLimit(ctx context.Context, n int) ([]*async_task.AsyncTask, error) {
	// default limit 5
	res, err := dal.getClient(ctx).Query().Where(eTask.State(string(vo.AsyncTaskStatePending))).Order(ent.Asc(eTask.FieldCreatedAt)).Limit(n).All(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*async_task.AsyncTask, 0, len(res))
	for _, r := range res {
		list = append(list, dal.convertor.ConvertAsyncTaskPoToDo(r))
	}
	return list, nil
}

func (dal *AsyncTaskDal) UpdateExecutingTaskSuccess(ctx context.Context, taskIDs ...string) error {
	_, err := dal.getClient(ctx).Update().SetState(vo.AsyncTaskStateSuccess.String()).
		Where(eTask.TaskIDIn(taskIDs...)).
		Where(eTask.StateEQ(vo.AsyncTaskStateExecuting.String())).
		Save(ctx)
	return err
}

func (dal *AsyncTaskDal) UpdateExecutingTaskFail(ctx context.Context, taskIDs ...string) error {
	_, err := dal.getClient(ctx).Update().SetState(vo.AsyncTaskStateFail.String()).
		Where(eTask.TaskIDIn(taskIDs...)).
		Where(eTask.StateEQ(vo.AsyncTaskStateExecuting.String())).
		Save(ctx)
	return err
}

func (dal *AsyncTaskDal) UpdatePendingTaskExecuting(ctx context.Context, taskIDs ...string) (int, error) {
	return dal.getClient(ctx).Update().SetState(vo.AsyncTaskStateExecuting.String()).
		Where(eTask.TaskIDIn(taskIDs...)).
		Where(eTask.StateEQ(vo.AsyncTaskStatePending.String())).
		Save(ctx)
}

// FindStuckExecutingTasks 查找长时间处于 executing 状态的任务
func (dal *AsyncTaskDal) FindStuckExecutingTasks(ctx context.Context, stuckDuration time.Duration, limit int) ([]*async_task.AsyncTask, error) {
	threshold := time.Now().Add(-stuckDuration)
	res, err := dal.getClient(ctx).Query().
		Where(eTask.StateEQ(vo.AsyncTaskStateExecuting.String())).
		Where(eTask.UpdatedAtLT(threshold)).
		Order(ent.Asc(eTask.FieldUpdatedAt)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]*async_task.AsyncTask, 0, len(res))
	for _, r := range res {
		list = append(list, dal.convertor.ConvertAsyncTaskPoToDo(r))
	}
	return list, nil
}

// ResetTaskToPendingWithRetry 将失败或执行中的任务重置为 pending 状态，并增加重试次数（CAS 操作）
func (dal *AsyncTaskDal) ResetTaskToPendingWithRetry(ctx context.Context, taskID string, currentState vo.AsyncTaskState) (int, error) {
	return dal.getClient(ctx).Update().
		SetState(vo.AsyncTaskStatePending.String()).
		AddRetryCount(1).
		Where(eTask.TaskIDEQ(taskID)).
		Where(eTask.StateEQ(currentState.String())).
		Save(ctx)
}
