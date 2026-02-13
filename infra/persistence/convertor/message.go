package convertor

import (
	"purchase/domain/entity/async_task"
	"purchase/domain/vo"
	"purchase/infra/persistence/dal/db/ent"
)

func (c *Convertor) ConvertAsyncTaskPoToDo(m *ent.AsyncTask) *async_task.AsyncTask {
	return &async_task.AsyncTask{
		ID:         m.ID,
		TaskID:     m.TaskID,
		TaskType:   vo.AsyncTaskType(m.TaskType),
		TaskName:   m.TaskName,
		EntityID:   m.BizID,
		TaskData:   m.TaskData,
		State:      vo.AsyncTaskState(m.State),
		RetryCount: m.RetryCount,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
