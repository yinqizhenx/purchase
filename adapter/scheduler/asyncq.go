package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/wire"
	"github.com/hibiken/asynq"
)

// const redisAddr = "127.0.0.1:6379"

// func main() {
// 	srv := asynq.NewServer(
// 		asynq.RedisClientOpt{Addr: "127.0.0.1:6379", DB: 1},
// 		asynq.Config{
// 			// Specify how many concurrent workers to use
// 			Concurrency: 10,
// 			// Optionally specify multiple queues with different priority.
// 			Queues: map[string]int{
// 				"critical": 6,
// 				"default":  3,
// 				"low":      1,
// 			},
// 			// See the godoc for other configuration options
// 		},
// 	)
//
// 	// mux maps a type to a handler
// 	mux := asynq.NewServeMux()
// 	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
// 	mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
// 	// ...register other handlers...
//
// 	if err := srv.Run(mux); err != nil {
// 		log.Fatalf("could not run server: %v", err)
// 	}
// }

var ProviderSet = wire.NewSet(NewAsyncQueueServer)

type AsyncQueueServer struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewAsyncQueueServer() *AsyncQueueServer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "127.0.0.1:6379", DB: 1},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)
	s := &AsyncQueueServer{
		server: srv,
		mux:    asynq.NewServeMux(),
	}
	return s
}

func (s *AsyncQueueServer) Run() error {
	// mux := asynq.NewServeMux()
	// mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	// mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
	// ...register other handlers...
	s.registerTask()
	if err := s.server.Run(s.mux); err != nil {
		log.Fatalf("could not run server: %v", err)
		return err
	}
	return nil
}

func (s *AsyncQueueServer) registerTask() {
	// mux := asynq.NewServeMux()
	s.mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	// mux.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
	// ...register other handlers...

	// if err := srv.Run(mux); err != nil {
	// 	log.Fatalf("could not run server: %v", err)
	// }
}

const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// Email delivery code ...
	return nil
}
