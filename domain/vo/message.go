package vo

type AsyncTaskState string

const (
	AsyncTaskState_Pending AsyncTaskState = "pending"
	AsyncTaskState_Done                   = "done"
)

type AsyncTaskType string

const (
	AsyncTaskType_Task  AsyncTaskType = "task"
	AsyncTaskType_Event AsyncTaskType = "event"
)
