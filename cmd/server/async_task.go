package server

import (
	"context"
	"encoding/json"
	"fmt"

	"purchase/infra/async_task"
	"purchase/infra/dlock"
	"purchase/infra/mq"
	"purchase/infra/persistence/dal"
	"purchase/infra/persistence/tx"
	"purchase/pkg/chanx"
)

func NewAsyncTaskServer(pub mq.Publisher, dal *dal.AsyncTaskDal, txm *tx.TransactionManager, ch *chanx.UnboundedChan[string], lockBuilder dlock.LockBuilder) *async_task.AsyncTaskMux {
	s := async_task.NewAsyncTaskMux(pub, dal, txm, ch, lockBuilder)
	s.RegisterHandler("task_name", func(ctx context.Context, payload []byte) error {
		// need to unmarshal args from payload correctly
		args := make([]string, 0)
		err := json.Unmarshal(payload, &args)
		if err != nil {
			return err
		}
		fmt.Println("here is an example", len(args))
		return nil
	})
	return s
}
