package vo

type AsyncTaskState string

const (
	AsyncTaskStatePending AsyncTaskState = "pending"
	AsyncTaskStateDone                   = "done"
)

type AsyncTaskType string

const (
	AsyncTaskTypeTask  AsyncTaskType = "task"
	AsyncTaskTypeEvent AsyncTaskType = "event"
)
