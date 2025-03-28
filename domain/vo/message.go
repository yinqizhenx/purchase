package vo

type AsyncTaskState string

func (s AsyncTaskState) String() string {
	return string(s)
}

const (
	AsyncTaskStatePending   AsyncTaskState = "pending"
	AsyncTaskStateExecuting AsyncTaskState = "executing"
	AsyncTaskStateSuccess   AsyncTaskState = "success"
	AsyncTaskStateFail      AsyncTaskState = "fail"
)

type AsyncTaskType string

const (
	AsyncTaskTypeTask  AsyncTaskType = "task"
	AsyncTaskTypeEvent AsyncTaskType = "event"
)

type AsyncTaskGroup string

const (
	GroupDefault      AsyncTaskGroup = "default"
	GroupAsyncMessage AsyncTaskGroup = "async_message"
)
