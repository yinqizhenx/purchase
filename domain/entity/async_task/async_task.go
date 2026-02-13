package async_task

import (
	"context"
	"time"

	"purchase/domain/vo"
)

type AsyncTask struct {
	ID          int64
	TaskID      string
	TaskType    vo.AsyncTaskType
	TaskGroup   vo.AsyncTaskGroup
	TaskName    string
	EntityID    string // 聚合根id
	TaskData    string
	State       vo.AsyncTaskState
	RetryCount  int
	ScheduledAt time.Time // 任务最早可执行时间，零值表示立即执行
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *AsyncTask) IsEvent(ctx context.Context) bool {
	return m.TaskType == vo.AsyncTaskTypeEvent
}

func (m *AsyncTask) IsTask(ctx context.Context) bool {
	return m.TaskType == vo.AsyncTaskTypeTask
}
