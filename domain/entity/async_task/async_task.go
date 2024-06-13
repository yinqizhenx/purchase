package async_task

import (
	"context"
	"time"

	"purchase/domain/vo"
)

type AsyncTask struct {
	ID        int64
	TaskID    string
	TaskType  vo.AsyncTaskType
	TaskName  string
	EntityID  string // 聚合根id
	TaskData  string
	State     vo.AsyncTaskState
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *AsyncTask) IsEvent(ctx context.Context) bool {
	return m.TaskType == vo.AsyncTaskType_Event
}

func (m *AsyncTask) IsTask(ctx context.Context) bool {
	return m.TaskType == vo.AsyncTaskType_Task
}
