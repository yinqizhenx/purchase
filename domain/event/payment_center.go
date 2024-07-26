package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"purchase/domain/entity/async_task"
)

var _ Event = (*PACreated)(nil)

type PACreated struct {
	EventID   string
	PACode    string
	CreatedBy string    `json:"created_by"` // 创建人
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

func (e *PACreated) ID() string {
	return e.EventID
}

func (e *PACreated) EntityID() string {
	return e.PACode
}

func (e *PACreated) OccurredOn() time.Time {
	return e.CreatedAt
}

func (e *PACreated) OrderID() int64 {
	return 0
}

func (e *PACreated) EventName() string {
	// json marshal
	return "pa_created"
}

func (e *PACreated) Encode() ([]byte, error) {
	return json.Marshal(e)
}

func (e *PACreated) Decode(v []byte) (Event, error) {
	data := &PACreated{}
	err := json.Unmarshal(v, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *PACreated) ToAsyncTask() (*async_task.AsyncTask, error) {
	c, err := e.Encode()
	if err != nil {
		return nil, err
	}
	return &async_task.AsyncTask{
		TaskID:   uuid.New().String(),
		TaskType: "event",
		TaskName: e.EventName(),
		EntityID: e.PACode,
		TaskData: string(c),
		State:    "pending",
	}, nil
}

type ProductInventoryChanged struct {
	EventId     string
	AggregateId int64
	UpdatedAt   time.Time
	OldValue    int // 变更前状态
	NewValue    int // 变更后状态
}
