package async_queue

import (
	"github.com/google/wire"
	"github.com/hibiken/asynq"
	// "go-study/asyncq/tasks"
)

const redisAddr = "127.0.0.1:6379"

const db = 1

var ProviderSet = wire.NewSet(NewAsynqClient)

func NewAsynqClient() (*asynq.Client, func()) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, DB: db})
	return client, func() {
		client.Close()
	}
}

// func main() {
// 	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, DB: db})
// 	defer client.Close()
//
// 	// ------------------------------------------------------
// 	// Example 1: Enqueue task to be processed immediately.
// 	//            Use (*Client).Enqueue method.
// 	// ------------------------------------------------------
//
// 	task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
// 	if err != nil {
// 		log.Fatalf("could not create task: %v", err)
// 	}
// 	info, err := client.Enqueue(task)
// 	if err != nil {
// 		log.Fatalf("could not enqueue task: %v", err)
// 	}
// 	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
//
// 	// ------------------------------------------------------------
// 	// Example 2: Schedule task to be processed in the future.
// 	//            Use ProcessIn or ProcessAt option.
// 	// ------------------------------------------------------------
//
// 	info, err = client.Enqueue(task, asynq.ProcessIn(24*time.Hour))
// 	if err != nil {
// 		log.Fatalf("could not schedule task: %v", err)
// 	}
// 	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
//
// 	// ----------------------------------------------------------------------------
// 	// Example 3: Set other options to tune task processing behavior.
// 	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
// 	// ----------------------------------------------------------------------------
//
// 	task, err = tasks.NewImageResizeTask("https://example.com/myassets/image.jpg")
// 	if err != nil {
// 		log.Fatalf("could not create task: %v", err)
// 	}
// 	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
// 	if err != nil {
// 		log.Fatalf("could not enqueue task: %v", err)
// 	}
// 	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
// }
